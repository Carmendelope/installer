/*
 * Copyright 2018 Nalej
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
 */

package workflow

import (
	"fmt"
	"github.com/nalej/installer/internal/pkg/errors"
	"github.com/rs/zerolog/log"
	"strings"

	"github.com/nalej/installer/internal/pkg/workflow/entities"
	"github.com/nalej/installer/internal/pkg/workflow/handler"


	"github.com/nalej/derrors"
)

// OK constant for commands that succeeded in the execution.
const OK = "OK"

// Fail constant for commands whose execution failed.
const Fail = "Fail"

var executorLogger = log.With().Str("component", "workflow.executor").Logger()

// Executor structure.
type Executor struct {
	*Workflow
	handler        handler.CommandHandler
	currentCommand int
	// ExecutionLog contains the log entries for all commands in the workflow.
	ExecutionLog []string `json:"executionLog"`
	logListener  func(msg string)
	// State contains the workflow state.
	State            WorkflowState `json:"state"`
	workflowCallback func(workflowID string, error derrors.Error, state WorkflowState)
	Parameters       map[string]string `json:"parameters"`
}

// NewWorkflowExecutor creates a new executor
//   params:
//     workflow The workflow to be executed.
//     executionHandler The async commands callback commandHandler.
func NewWorkflowExecutor(workflow *Workflow,
	workflowCallback func(workflowID string, error derrors.Error, state WorkflowState)) *Executor {
	return &Executor{workflow, handler.GetCommandHandler(),
		0, make([]string, 0), nil,
		InitState, workflowCallback, make(map[string]string, 0)}
}

// SetLogListener attaches a given function as the log listener for input log entries.
func (e *Executor) SetLogListener(f func(msg string)) {
	e.logListener = f
}

func (e *Executor) executeCommand(index int) derrors.Error {
	if index >= len(e.Workflow.Commands) {
		return derrors.NewOperationError(errors.InvalidCommandIndex).WithParams(index, e.Workflow)
	}
	e.currentCommand = index
	go func() {
		toExecuted := e.Workflow.Commands[index]
		e.execOnBackground(index, toExecuted)
	}()

	return nil
}

func (e *Executor) execOnBackground(index int, cmd entities.Command) {
	err := e.handler.AddCommand(cmd.ID(), e.commandCallback, e.logCallback)
	if err != nil {
		// If the executor cannot allocate the callback the workflow fails.
		e.failed(err)
		return
	}

	if cmd.Name() != entities.Logger {
		e.AddLogEntry("Executing: " + cmd.UserString())
	}
	if cmd.Type() == entities.SyncCommandType {
		executorLogger.Debug().Str("cmd", cmd.String()).Msg("Executing sync command")
		result, err := cmd.(entities.SyncCommand).Run(e.Workflow.WorkflowID)

		err = e.handler.FinishCommand(cmd.ID(), result, err)
		if err != nil {
			e.failed(err)
		}
	} else {
		executorLogger.Debug().Str("cmd", cmd.String()).Msg("Executing async command")
		err := cmd.(entities.AsyncCommand).Run(e.Workflow.WorkflowID)
		if err != nil {
			//If the execution return errors, the executor call to the commandHandler with the error.
			err = e.handler.FinishCommand(cmd.ID(), nil, err)
			if err != nil {
				e.failed(err)
			}
		}
	}
}

func (e *Executor) commandCallback(cmdID string, result *entities.CommandResult, error derrors.Error) {
	// To support parallel execution of commands, we can implement a barrier command that will make commandCallback
	// not to launch more commands until all pending commands have finished.

	if error != nil {
		// Stop workflow execution
		e.failed(derrors.NewOperationError(errors.WorkflowExecutionFailed).CausedBy(error))
		return
	}

	if result != nil {
		//e.AddLogEntry("Success: " + strconv.FormatBool((*result).Success))
		if (*result).HasOutput() && (*result).ShowResult() {
			if strings.Contains(cmdID, entities.Logger) {
				e.AddLogEntry((*result).Output)
			} else {
				e.AddLogEntry(fmt.Sprintf("Command %s:\n%s", cmdID, (*result).Output))
			}
		}

		if (*result).Success {
			if e.currentCommand == len(e.Workflow.Commands)-1 {
				executorLogger.Debug().Interface("workflowState", e.State).Msg("all commands have been executed")
				e.AddLogEntry("All commands have been executed")
				e.State = FinishedState
				e.workflowCallback(e.Workflow.WorkflowID, nil, e.State)
				return
			}

			err := e.executeCommand(e.currentCommand + 1)
			if err != nil {
				e.failed(err)
			}
		} else {
			log.Warn().Str("workflowID", e.WorkflowID).Msg(result.String())
			e.failed(derrors.NewOperationError(errors.WorkflowExecutionFailed).WithParams(result.String()))
		}
	} else {
		e.failed(derrors.NewOperationError(errors.InvalidWorkflowState))
	}

}

func (e *Executor) logCallback(id string, logEntry string) {
	e.AddLogEntry(logEntry)
}

// Exec starts the execution of the target workflow.
func (e *Executor) Exec() {
	if len(e.Workflow.Commands) > 0 {
		executorLogger.Debug().Str("workflowID", e.WorkflowID).Int("numCommands", len(e.Workflow.Commands)).
			Msg("Executing workflow")
		e.State = InProgressState
		err := e.executeCommand(0)
		if err != nil {
			e.failed(err)
		}
		return
	}
	e.failed(derrors.NewOperationError(errors.WorkflowWithoutCommands))
}

func (e *Executor) failed(reason derrors.Error) {
	e.AddLogEntry(reason.Error())
	e.AddLogEntry(Fail)
	e.State = ErrorState
	e.workflowCallback(e.Workflow.WorkflowID, reason, e.State)
}

func (e *Executor) commandLogListener(logEntry string) {
	e.ExecutionLog = append(e.ExecutionLog, "commandLogListener: "+logEntry)
	if e.logListener != nil {
		e.logListener("commandLogListener: " + logEntry)
	}
}

// AddLogEntry adds a new line to the log.
func (e *Executor) AddLogEntry(line string) {
	e.ExecutionLog = append(e.ExecutionLog, line)
	if e.logListener != nil {
		e.logListener(line)
	}
}

// Log retrieves the execution log of the current workflow.
func (e *Executor) Log() []string {
	logCopy := make([]string, len(e.ExecutionLog))
	copy(logCopy, e.ExecutionLog)
	return logCopy
}

// CurrentCommand returns the index of the command being executed and the total of commands to be executed in
// in the workflow.
func (e *Executor) CurrentCommand() (int, int) {
	return e.currentCommand, len(e.Commands)
}

// ParameterSet upserts a workflow parameter.
func (e *Executor) ParameterSet(key string, value string) {
	e.Parameters[key] = value
}

// ParameterGet retrieves the value of a given key.
func (e *Executor) ParameterGet(key string) (*WorkflowParameter, derrors.Error) {
	value, exists := e.Parameters[key]
	if exists {
		return NewWorkflowParameter(key, value), nil
	}
	return nil, derrors.NewOperationError(errors.ParameterDoesNotExists)
}
