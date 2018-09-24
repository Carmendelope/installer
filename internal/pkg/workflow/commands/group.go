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

// Group command
// Permits to execute a subset of commands sequentially. The command manages the sequential execution of the child
// commands and calls its callback function once all commands are executed.
//
// {"type":"sync", "name": "group", "commands": [{"type":...},{"type":...}]}

package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/nalej/installer/internal/pkg/errors"
	"github.com/rs/zerolog/log"
	"strings"

	"github.com/nalej/derrors"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
	"github.com/nalej/installer/internal/pkg/workflow/handler"

)

// Group structure with the commands to be executed.
type Group struct {
	entities.GenericSyncCommand
	Description        string             `json:"description"`
	Commands           []entities.Command `json:"commands"`
	commandHandler     handler.CommandHandler
	commandResults     map[string]entities.CommandResult
	executionErrors    map[string]derrors.Error
	asyncFinishChannel chan string
	asyncCmdID         string
}

// GroupFromJSON structure with helper RawMessage to parse
// target commands.
type GroupFromJSON struct {
	entities.GenericCommand
	Description string            `json:"description"`
	Commands    []json.RawMessage `json:"commands"`
}

// ToGroup transforms a JSON group into a Group parsing the required commands.
func (gfj *GroupFromJSON) ToGroup() (*Group, derrors.Error) {
	p := NewCmdParser()
	cmds := make([]entities.Command, 0)
	for _, toParse := range gfj.Commands {
		toAdd, err := p.ParseCommand(toParse)
		if err != nil {
			return nil, err
		}
		cmds = append(cmds, *toAdd)
	}
	return NewGroup(gfj.Description, cmds), nil
}

// NewGroup creates a new Group with a given description and associated commands.
func NewGroup(description string, cmds []entities.Command) *Group {
	return &Group{
		*entities.NewSyncCommand(entities.GroupCmd),
		description, cmds,
		handler.GetCommandHandler(),
		make(map[string]entities.CommandResult),
		make(map[string]derrors.Error),
		make(chan string), ""}
}

// NewGroupFromJSON creates a new command from a raw json payload.
func NewGroupFromJSON(raw []byte) (*entities.Command, derrors.Error) {
	gfj := &GroupFromJSON{}
	if err := json.Unmarshal(raw, &gfj); err != nil {
		return nil, derrors.NewOperationError(errors.UnmarshalError, err)
	}
	toGroup, err := gfj.ToGroup()
	if err != nil {
		return nil, err
	}
	var r entities.Command = toGroup
	return &r, nil
}

// Run the current command.
//   returns:
//     The CommandResult
//     An error if the command execution fails
func (g *Group) Run(workflowID string) (*entities.CommandResult, derrors.Error) {

	if len(g.Commands) == 0 {
		return nil, derrors.NewOperationError(errors.CannotExecuteSyncCommand)
	}

	g.commandHandler.AddLogEntry(g.CommandID, fmt.Sprintf("Starting sequential group execution of %d commands", len(g.Commands)))
	log.Info().Str("groupCmdId", g.CommandID).Str("description", g.Description).Msg("Executing sequential group")
	results := make([]entities.CommandResult, 0)
	for _, nextCommand := range g.Commands {
		result, err := g.executeCommand(workflowID, nextCommand)
		if err != nil {
			return nil, err
		}
		if result != nil {
			log.Debug().Str("groupCmdId", g.CommandID).Bool("success", result.Success).Msg("Adding result to group")
			results = append(results, *result)
			g.commandHandler.AddLogEntry(g.CommandID, result.UserString())
			if !result.Success {
				break
			}
		} else {
			log.Warn().Str("groupCmdId", g.CommandID).Str("cmd", nextCommand.String()).Msg("Empty result returned")
		}

	}
	log.Debug().Str("groupCmdId", g.CommandID).Msg("Building final group command result")

	overallSuccess := true
	var overallOutput bytes.Buffer
	var overallError = derrors.NewGenericError(errors.CannotExecuteSyncCommand)
	for _, value := range results {
		overallSuccess = overallSuccess && value.Success
		overallOutput.WriteString("Output\n" + value.Output + "\n")
		if value.Error != nil {
			overallError = overallError.WithParams(value.Error)
		}
	}
	//g.commandHandler.AddLogEntry(g.CommandID, fmt.Sprintf("Group result, success: %t, output: %s", overallSuccess, overallOutput))
	overallOutputString := overallOutput.String()
	log.Debug().Str("groupCmdId", g.CommandID).Bool("success", overallSuccess).Str("outputString", overallOutputString).Msg("Group result")
	if overallSuccess {
		return entities.NewCommandResultNoShow(overallSuccess, overallOutputString, nil), nil
	}
	log.Debug().Str("groupCmdId", g.CommandID).Str("err", overallError.DebugReport()).Msg("Group execution failed")
	return entities.NewCommandResult(overallSuccess, overallOutputString, overallError), nil
}

func (g *Group) executeCommand(workflowID string, cmd entities.Command) (*entities.CommandResult, derrors.Error) {
	err := g.commandHandler.AddCommand(cmd.ID(), g.commandCallback, g.logCallback)
	if err != nil {
		return nil, err
	}
	if cmd.Name() != entities.Logger {
		g.commandHandler.AddLogEntry(g.CommandID, "Executing: "+cmd.String()+" with Id: "+cmd.ID())
	}
	if cmd.Type() == entities.SyncCommandType {
		log.Debug().Str("groupCmdId", g.CommandID).Str("cmdID", cmd.ID()).Str("cmd", cmd.String()).Msg("SYNC")
		result, err := cmd.(entities.SyncCommand).Run(workflowID)
		if err != nil {
			log.Warn().Str("id", cmd.ID()).Str("err", err.DebugReport()).Msg("error executing sync command on sequential group")
			return nil, err
		}
		err = g.commandHandler.FinishCommand(cmd.ID(), result, err)
		return result, nil
	}
	// Async command expected
	log.Debug().Str("groupCmdId", g.CommandID).Str("cmdID", cmd.ID()).Str("cmd", cmd.String()).Msg("ASYNC")
	g.asyncCmdID = cmd.ID()
	err = cmd.(entities.AsyncCommand).Run(workflowID)
	if err != nil {
		log.Warn().Str("id", cmd.ID()).Str("err", err.DebugReport()).Msg("error executing async command on sequential group")
		//If the execution return errors, the executor call to the commandHandler with the error.
		err = g.commandHandler.FinishCommand(cmd.ID(), nil, err)
		if err != nil {
			log.Warn().Str("groupCmdId", g.CommandID).Str("err", err.DebugReport()).Msg("error on async.FinishCommand")
			return nil, err
		}
	} else {
		waitForCmd := <-g.asyncFinishChannel
		log.Debug().Str("groupCmdId", g.CommandID).Str("waitingFor", g.asyncCmdID).Str("finished", waitForCmd).
			Msg("Async command waiting")
		var cmdResult *entities.CommandResult = nil

		if value, exists := g.commandResults[waitForCmd]; exists {
			cmdResult = &value
		}
		if value, exists := g.executionErrors[waitForCmd]; exists {
			return cmdResult, value
		}
		if cmdResult == nil {
			return nil, derrors.NewOperationError(errors.InvalidWorkflowState)
		}
		return cmdResult, nil
	}
	return nil, derrors.NewOperationError(errors.InvalidWorkflowState)

}

func (g *Group) commandCallback(cmdID string, result *entities.CommandResult, error derrors.Error) {
	log.Debug().Str("groupCmdId", g.CommandID).Str("cmdID", cmdID).Msg("received callback from sequential command")
	if result != nil {
		g.commandResults[cmdID] = *result
	}
	if error != nil {
		g.executionErrors[cmdID] = error
	}
	if cmdID == g.asyncCmdID {
		g.asyncFinishChannel <- cmdID
	}

}

func (g *Group) logCallback(cmdID string, logEntry string) {
	//g.commandHandler.AddLogEntry(g.CommandID, fmt.Sprintf("[%s] %s", cmdID, logEntry))
	g.commandHandler.AddLogEntry(g.CommandID, fmt.Sprintf("[%s] %s", g.Description, logEntry))
}

// String obtains a string representation
func (g *Group) String() string {
	cmdNames := make([]string, 0)
	for _, cmd := range g.Commands {
		cmdNames = append(cmdNames, cmd.Name())
	}
	return fmt.Sprintf("SYNC Group %s actions: %s", g.Description, strings.Join(cmdNames, ", "))
}

// PrettyPrint returns a simple space indexed string.
func (g *Group) PrettyPrint(identation int) string {
	cmdNames := make([]string, 0)
	for _, cmd := range g.Commands {
		cmdNames = append(cmdNames, cmd.PrettyPrint(identation+2))
	}
	return fmt.Sprintf("%sSYNC Group %s actions:\n%s\n", strings.Repeat(" ", identation), g.Description, strings.Join(cmdNames, "\n"))
}

// UserString returns a simple string representation of the command for the user.
func (g *Group) UserString() string {
	cmdNames := make([]string, 0)
	for _, cmd := range g.Commands {
		cmdNames = append(cmdNames, cmd.Name())
	}
	return fmt.Sprintf("Sequential group: %s actions: %s", g.Description, strings.Join(cmdNames, ", "))
}
