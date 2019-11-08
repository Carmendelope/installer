/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
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
