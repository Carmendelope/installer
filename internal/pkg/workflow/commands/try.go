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

// The try will execute a command first. If the result is ok, that result will be returned. If the command fails, it
// will execute another command.
//

package commands

import (
	"encoding/json"
	"fmt"
	"github.com/nalej/installer/internal/pkg/errors"
	"github.com/rs/zerolog/log"
	"strings"

	"github.com/nalej/derrors"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
	"github.com/nalej/installer/internal/pkg/workflow/handler"
)

// Try command structure with the command to be executed and the alternative in case of failure.
type Try struct {
	entities.GenericSyncCommand
	Description        string           `json:"description"`
	TryCommand         entities.Command `json:"cmd"`
	OnFailCommand      entities.Command `json:"onFail"`
	commandHandler     handler.CommandHandler
	commandResult      *entities.CommandResult
	executionError     derrors.Error
	asyncFinishChannel chan string
	asyncCmdID         string
}

// NewTry creates a new Try command with all parameters.
func NewTry(description string, tryCommand entities.Command, onFailCommand entities.Command) *Try {
	return &Try{*entities.NewSyncCommand(entities.TryCmd),
		description, tryCommand, onFailCommand,
		handler.GetCommandHandler(), nil, nil,
		make(chan string), ""}
}

// TryFromJSON structure required to be able to parse individual commands.
type TryFromJSON struct {
	entities.GenericCommand
	Description   string          `json:"description"`
	TryCommand    json.RawMessage `json:"cmd"`
	OnFailCommand json.RawMessage `json:"onFail"`
}

// ToTry transforms the raw JSON structure into a Try by parsing individual commands.
func (tfj *TryFromJSON) ToTry() (*Try, derrors.Error) {
	p := NewCmdParser()
	tryCommand, err := p.ParseCommand(tfj.TryCommand)
	if err != nil {
		return nil, err
	}
	onFailCommand, err := p.ParseCommand(tfj.OnFailCommand)
	if err != nil {
		return nil, err
	}
	return NewTry(tfj.Description, *tryCommand, *onFailCommand), nil
}

// NewTryFromJSON creates a command using a raw JSON payload.
func NewTryFromJSON(raw []byte) (*entities.Command, derrors.Error) {
	tfj := &TryFromJSON{}
	if err := json.Unmarshal(raw, &tfj); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	toTry, err := tfj.ToTry()
	if err != nil {
		return nil, err
	}
	var r entities.Command = toTry
	return &r, nil
}

// Run the current command.
//   returns:
//     The CommandResult
//     An error if the command execution fails
func (t *Try) Run(workflowID string) (*entities.CommandResult, derrors.Error) {
	t.commandHandler.AddLogEntry(t.CommandID, fmt.Sprintf("Try %s", t.TryCommand.Name()))
	result, err := t.executeCommand(workflowID, t.TryCommand)
	if err != nil {
		log.Debug().Str("err", err.Error()).Msg("retry on cmd error")
	}
	if err != nil || !result.Success {
		result, err = t.executeCommand(workflowID, t.OnFailCommand)
	}
	if result.Success {
		return entities.NewCommandResultNoShow(true, result.Output, nil), nil
	}
	return result, nil
}

func (t *Try) executeCommand(workflowID string, cmd entities.Command) (*entities.CommandResult, derrors.Error) {
	err := t.commandHandler.AddCommand(cmd.ID(), t.commandCallback, t.logCallback)
	if err != nil {
		return nil, err
	}
	if cmd.Name() != entities.Logger {
		t.commandHandler.AddLogEntry(t.CommandID, "Executing: "+cmd.String()+" with Id: "+cmd.ID())
	}

	if cmd.Type() == entities.SyncCommandType {
		log.Debug().Str("cmd", t.CommandID).Str("cmdID", cmd.ID()).Str("cmd", cmd.String()).Msg("SYNC")
		result, err := cmd.(entities.SyncCommand).Run(workflowID)
		if err != nil {
			log.Warn().Str("cmd", t.CommandID).Str("cmdID", cmd.ID()).Str("err", err.DebugReport()).
				Msg("error executing sync command on sequential group: ")
			return nil, err
		}
		err = t.commandHandler.FinishCommand(cmd.ID(), result, err)
		return result, nil
	}
	// Assume async command.
	log.Debug().Str("cmd", t.CommandID).Str("cmdID", cmd.ID()).Str("cmd", cmd.String()).Msg("ASYNC")
	t.asyncCmdID = cmd.ID()
	err = cmd.(entities.AsyncCommand).Run(workflowID)
	if err != nil {
		log.Warn().Str("cmd", t.CommandID).Str("cmdID", cmd.ID()).Str("err", err.DebugReport()).
			Msg("error executing async command on sequential group: ")
		//If the execution return errors, the executor call to the commandHandler with the error.
		err = t.commandHandler.FinishCommand(cmd.ID(), nil, err)
		if err != nil {
			log.Warn().Str("cmd", t.CommandID).Str("err", err.DebugReport()).
				Msg("error on async.FinishCommand")
			return nil, err
		}
	} else {
		waitForCmd := <-t.asyncFinishChannel
		log.Debug().Str("cmd", t.CommandID).Str("awaitingCmdId", t.asyncCmdID).Str("finished", waitForCmd).
			Msg("Async command waiting")
		var cmdResult = t.commandResult

		if t.executionError != nil {
			return cmdResult, t.executionError
		}

		if cmdResult == nil {
			return nil, derrors.NewInternalError(errors.InvalidWorkflowState)
		}
		return cmdResult, nil
	}
	return nil, derrors.NewInternalError(errors.InvalidWorkflowState)

}

func (t *Try) commandCallback(cmdID string, result *entities.CommandResult, error derrors.Error) {
	log.Debug().Str("cmd", t.CommandID).Str("cmdID", cmdID).Msg("received callback from command")
	if result != nil {
		t.commandResult = result
	}
	if error != nil {
		t.executionError = error
	}
	if cmdID == t.asyncCmdID {
		t.asyncFinishChannel <- cmdID
	}
}

func (t *Try) logCallback(cmdID string, logEntry string) {
	t.commandHandler.AddLogEntry(t.CommandID, fmt.Sprintf("[%s] %s", t.Description, logEntry))
}

// String obtains a string representation
func (t *Try) String() string {
	return fmt.Sprintf("SYNC Try %s execute: %s onFailure: %s", t.Description, t.TryCommand.Name(), t.OnFailCommand.Name())
}

// PrettyPrint returns a simple space indexed string.
func (t *Try) PrettyPrint(indentation int) string {
	return fmt.Sprintf("%sSYNC Try %s execute:\n%s\n%sonFailure:\n%s\n",
		strings.Repeat(" ", indentation), t.Description,
		t.TryCommand.PrettyPrint(indentation+2),
		strings.Repeat(" ", indentation), t.OnFailCommand.PrettyPrint(indentation+2))
}

// UserString returns a simple string representation of the command for the user.
func (t *Try) UserString() string {
	return fmt.Sprintf("Try %s execute: %s onFailure: %s", t.Description, t.TryCommand.Name(), t.OnFailCommand.Name())
}
