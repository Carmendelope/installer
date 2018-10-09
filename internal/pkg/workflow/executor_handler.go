/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package workflow

import (
	"github.com/nalej/installer/internal/pkg/errors"
	"github.com/rs/zerolog/log"
	"sync"

	"github.com/nalej/derrors"
)

var instanceExecutorHandler ExecutorHandler
var onceExecutorHandler sync.Once

// ExecutorHandler interface defining the operations available in the handler of the workflow execution.
type ExecutorHandler interface {
	Add(workflow *Workflow,
		workflowCallback func(workflowID string, error derrors.Error,
			state WorkflowState)) (*Executor, derrors.Error)
	Execute(workflowID string) (*Executor, derrors.Error)
	Get(workflowID string) (*Executor, derrors.Error)
	Stop(workflowID string) derrors.Error
}

type executorHandler struct {
	sync.Mutex
	executorMap map[string]*Executor
}

// NewExecutorHandler creates an ExecutorHandler.
func NewExecutorHandler() ExecutorHandler {
	return &executorHandler{sync.Mutex{}, make(map[string]*Executor)}
}

// GetExecutorHandler allows to retrieve the singleton ExecutorHandler instance.
func GetExecutorHandler() ExecutorHandler {
	onceExecutorHandler.Do(func() {
		instanceExecutorHandler = NewExecutorHandler()
	})
	return instanceExecutorHandler
}

func (handler *executorHandler) Add(workflow *Workflow,
	workflowCallback func(workflowID string, error derrors.Error,
		state WorkflowState)) (*Executor, derrors.Error) {
	handler.Lock()
	defer handler.Unlock()
	exe := NewWorkflowExecutor(workflow, workflowCallback)
	_, exist := handler.executorMap[exe.WorkflowID]
	if exist {
		return nil, derrors.NewAlreadyExistsError(errors.WorkflowAlreadyExists).WithParams(workflow.WorkflowID)
	}
	handler.executorMap[exe.WorkflowID] = exe
	return exe, nil
}
func (handler *executorHandler) Get(workflowID string) (*Executor, derrors.Error) {
	handler.Lock()
	defer handler.Unlock()
	exe, exist := handler.executorMap[workflowID]
	if !exist {
		return nil, derrors.NewAlreadyExistsError(errors.WorkflowDoesNotExists).WithParams(workflowID)
	}
	return exe, nil
}

func (handler *executorHandler) Execute(workflowID string) (*Executor, derrors.Error) {
	exe, err := handler.Get(workflowID)
	if err != nil {
		return nil, err
	}
	exe.Exec()
	return exe, nil
}

func (handler *executorHandler) Stop(workflowID string) derrors.Error {
	log.Debug().Str("workflowID", workflowID).Msg("ExecutorHandler stop request")
	exe, err := handler.Get(workflowID)
	if err != nil {
		return err
	}
	exe.Stop()
	delete(handler.executorMap, workflowID)
	return nil
}
