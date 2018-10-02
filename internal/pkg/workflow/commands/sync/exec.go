/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

// Exec command
// Executes an arbitrary command.
//
// {"type":"sync", "name": "exec", "cmd": "ls", "args":["-lash", "/tmp/."]}

package sync

import (
	"encoding/json"
	"github.com/nalej/installer/internal/pkg/errors"
	"os/exec"
	"strings"

	"github.com/nalej/derrors"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
)

// Exec command structure with supported parameters.
type Exec struct {
	entities.GenericSyncCommand
	Cmd  string   `json:"cmd"`
	Args []string `json:"args"`
}

// NewExec creates an Exec command from a set of parameters.
func NewExec(cmd string, args []string) *Exec {
	return &Exec{
		*entities.NewSyncCommand(entities.Exec),
		cmd, args}
}

// NewExecFromJSON creates an Exec command from a JSON object.
func NewExecFromJSON(raw []byte) (*entities.Command, derrors.Error) {
	exec := &Exec{}
	if err := json.Unmarshal(raw, &exec); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	exec.CommandID = entities.GenerateCommandID(exec.Name())
	var r entities.Command = exec
	return &r, nil
}

// Run the current command.
//   returns:
//     The CommandResult
//     An error if the command execution fails
func (e *Exec) Run(_ string) (*entities.CommandResult, derrors.Error) {

	// TODO Proper exit code manipulation
	// It seems that a lot of people are struggling with this cause there is not an easy way to determine the exit
	// status of a command. Two things should be added to this method.
	// 1) Check if the command execution fails in we cannot execute the command, not that the command execution fails itself.
	// 2) Return a DaishoError if the command cannot be executed, a Success/Fail command otherwise.
	// Related:
	// https://stackoverflow.com/questions/10385551/get-exit-code-go
	// https://groups.google.com/forum/#!topic/golang-nuts/MI4TyIkQqqg
	// https://groups.google.com/forum/#!msg/golang-nuts/dKbL1oOiCIY/OCfhH2rFp80J

	cmd := exec.Command(e.Cmd, e.Args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return nil, derrors.NewInternalError(errors.CannotExecuteSyncCommand, err).WithParams(e.Cmd, e.Args)
	}

	return entities.NewSuccessCommand(output), nil
}

// String obtains a string representation
func (e *Exec) String() string {
	return "SYNC Exec " + e.Cmd + strings.Join(e.Args, " ")
}

// PrettyPrint returns a simple space indexed string.
func (e *Exec) PrettyPrint(indentation int) string {
	return strings.Repeat(" ", indentation) + e.String()
}

// UserString returns a simple string representation of the command for the user.
func (e *Exec) UserString() string {
	return "Exec " + e.Cmd + strings.Join(e.Args, " ")
}
