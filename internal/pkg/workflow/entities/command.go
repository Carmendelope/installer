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

// This file contains the command definition

package entities

import (
	"fmt"
	"github.com/satori/go.uuid"

	"github.com/nalej/derrors"
)

// CommandType defines the different types of commands in the system.
type CommandType string

// SyncCommandType represents a command that is synchronously executed.
const SyncCommandType CommandType = "sync"

// AsyncCommandType represents a command that is asynchronously executed.
const AsyncCommandType CommandType = "async"

// ValidCommandType checks the type enum to determine if the string belongs to the enumeration.
//   params:
//     commandType The type to be checked
//   returns:
//     Whether it is contained in the enum.
func ValidCommandType(commandType CommandType) bool {
	switch commandType {
	case "":
		return false
	case SyncCommandType:
		return true
	case AsyncCommandType:
		return true
	default:
		return false
	}
}

// Command interface.
type Command interface {
	//ID is the internal command identification.
	ID() string
	// Type returns the CommandType
	Type() CommandType
	// Name of the command to be executed
	Name() string
	// Obtain a string representation
	String() string
	// PrettyPrint returns a simple space indexed string.
	PrettyPrint(identation int) string
	// UserString returns a simple string representation of the command for the user.
	UserString() string
}

// GenericCommand providing a type and name.
type GenericCommand struct {
	// CommandID.
	CommandID string `json:"id"`
	// CommandType.
	CommandType CommandType `json:"type"`
	// CommandName with the command name.a
	CommandName string `json:"name"`
}

//ID is the internal command identification.
func (gc *GenericCommand) ID() string {
	return gc.CommandID
}

// Type returns the CommandType
func (gc *GenericCommand) Type() CommandType {
	return gc.CommandType
}

// Name of the command to be executed
func (gc *GenericCommand) Name() string {
	return gc.CommandName
}

// NewGenericCommand creates a basic GenericCommand.
func NewGenericCommand(commandType CommandType, name string) GenericCommand {
	id := GenerateCommandID(name)
	return GenericCommand{id, commandType, name}
}

// CommandResult structure defines the elements of a command result.
// TODO This may be converted into an interface in the future if we support other result types.
type CommandResult struct {
	// Success returns true if the command was executed successfully.
	Success bool `json:"success"`
	// Output returns the command output in case of success, "" otherwise.
	Output string `json:"output"`
	// Error returns a DaishoError in case of command failure.
	Error      derrors.Error `json:"error"`
	showResult bool
}

// UserString provides a string to be reported to the final user.
func (cr *CommandResult) UserString() string {
	if cr.Error != nil {
		return fmt.Sprintf("Command failed\nOutput:\n%s\nError:\n%s", cr.Output, cr.Error.DebugReport())
	}
	return cr.Output
}

// String obtains a string representation of the current CommandResult.
func (cr *CommandResult) String() string {
	if cr.Error != nil {
		return fmt.Sprintf("Success: %t\nOutput:\n%s\nError:\n%s\n", cr.Success, cr.Output, cr.Error.DebugReport())
	}
	return fmt.Sprintf("Success: %t\nOutput:\n%s\n", cr.Success, cr.Output)
}

// CommandResultFromJSON is structure used to unmarshal command results. Instead of using the error interface as element
// it contains a GenericError for the deserialization.
type CommandResultFromJSON struct {
	Success bool                  `json:"success"`
	Output  string                `json:"output"`
	Error   *derrors.GenericError `json:"error"`
}

// ToCommandResult generates a CommandResult from the current structure.
func (crfj *CommandResultFromJSON) ToCommandResult() *CommandResult {
	if crfj.Error != nil {
		var daishoError derrors.Error = crfj.Error
		return &CommandResult{crfj.Success, crfj.Output, daishoError, true}
	}
	return &CommandResult{crfj.Success, crfj.Output, nil, true}
}

// NewCommandResult creates a new CommandResult.
func NewCommandResult(success bool, output string, err derrors.Error) *CommandResult {
	return &CommandResult{success, output, err, true}
}

// NewCommandResultNoShow creates a new CommandResult whose result will not be reported.
func NewCommandResultNoShow(success bool, output string, err derrors.Error) *CommandResult {
	return &CommandResult{success, output, err, false}
}

// NewSuccessCommand creates a successful command result.
func NewSuccessCommand(output []byte) *CommandResult {
	return &CommandResult{true, string(output), nil, true}
}

// NewErrCommand creates a failed command result.
func NewErrCommand(output string, err derrors.Error) *CommandResult {
	return &CommandResult{false, output, err, true}
}

// HasOutput checks if the command result has output attached to it.
func (cr *CommandResult) HasOutput() bool {
	return cr.Output != ""
}

// ShowResult checks if the command result should be reported to the user.
func (cr *CommandResult) ShowResult() bool {
	return cr.showResult
}

// SyncCommand interface defines the functions synchronous commands need to implement.
type SyncCommand interface {
	// Run the current command.
	//   returns:
	//     The CommandResult
	//     An error if the command execution fails
	Run(workflowID string) (*CommandResult, derrors.Error)
}

// GenericSyncCommand is a basic synchronous command.
type GenericSyncCommand struct {
	GenericCommand
}

// NewSyncCommand creates a GenericSyncCommand.
func NewSyncCommand(name string) *GenericSyncCommand {
	return &GenericSyncCommand{NewGenericCommand(SyncCommandType, name)}
}

// AsyncCommand interfaces defines the functions asynchronous commands need to implement.
type AsyncCommand interface {
	// Actions return the array of actions in the command.
	Actions() []Action
	// Run the current command.
	//   returns:
	//     An error if the command execution fails
	Run(workflowID string) derrors.Error
}

// GenericAsyncCommand structure with the command actions.
type GenericAsyncCommand struct {
	GenericCommand
	CommandActions []Action `json:"actions"`
}

// Actions return the array of actions in the command.
func (cmd *GenericAsyncCommand) Actions() []Action {
	return cmd.CommandActions
}

// NewAsyncCommand creates a new GenericAsyncCommand.
//   params:
//     name The command name.
//     actions The command actions.
//   returns:
//     A GenericAsyncCommand.
func NewAsyncCommand(name string, actions []Action) *GenericAsyncCommand {
	return &GenericAsyncCommand{NewGenericCommand(AsyncCommandType, name), actions}
}

// GenerateCommandID creates a new identifier with a given prefix.
func GenerateCommandID(name string) string {
	genID := uuid.NewV4().String()
	return fmt.Sprintf("cmd-%s-%s", name, genID)
}
