/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

// Process check command
// This command permits to check if a remote host is running a given process. The functionality is similar to checking
// the output with pgrep on the remote host.
//
// {"type":"sync", "name": "processcheck", "targetHost": "127.0.0.1", "targetPort": "22",
// "credentials":{"username": "username", "password":"passwd"},
// "process":"<name_of_process>", "shouldExists":true}

package sync

import (
	"encoding/json"
	"fmt"
	"github.com/nalej/installer/internal/pkg/errors"
	"github.com/rs/zerolog/log"
	"strings"

	"github.com/nalej/derrors"
	"github.com/nalej/installer/internal/pkg/workflow/commands/sync/connection"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
)

// ProcessCheck command structure with supported parameters.
type ProcessCheck struct {
	entities.GenericSyncCommand
	// Target node
	TargetHost string `json:"targetHost"`
	// Target port
	TargetPort string `json:"targetPort"`
	// Credentials for SSH.
	Credentials entities.Credentials `json:"credentials"`
	// Command to be execute
	Process string `json:"process"`
	// Command arguments
	ShouldExists bool `json:"shouldExists"`
}

// NewProcessCheck creates an ProcessCheck command from a set of parameters.
func NewProcessCheck(targetHost string, targetPort string, credentials entities.Credentials, process string, shouldExists bool) *ProcessCheck {
	return &ProcessCheck{*entities.NewSyncCommand(entities.SSH),
		targetHost,
		targetPort,
		credentials,
		process,
		shouldExists}
}

// NewProcessCheckFromJSON creates an ProcessCheck command from a JSON object.
func NewProcessCheckFromJSON(raw []byte) (*entities.Command, derrors.Error) {
	pc := &ProcessCheck{}
	if err := json.Unmarshal(raw, &pc); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	pc.CommandID = entities.GenerateCommandID(pc.Name())
	var r entities.Command = pc
	return &r, nil
}

func (pc *ProcessCheck) getTargetPort() string {
	if pc.TargetPort != "" {
		return pc.TargetPort
	}
	return DefaultSSHPort
}

// Run the current command.
//   returns:
//     The CommandResult
//     An error if the command execution fails
func (pc *ProcessCheck) Run(_ string) (*entities.CommandResult, derrors.Error) {

	conn, err := connection.NewSSHConnection(
		pc.TargetHost, pc.getTargetPort(),
		pc.Credentials.Username, pc.Credentials.Password, "", pc.Credentials.PrivateKey)
	if err != nil {
		log.Warn().Str("targetHost", pc.TargetHost).Err(err).Msg("Cannot establish connection")
		return nil, derrors.NewInternalError(errors.SSHConnectionError, err)
	}
	cmd := fmt.Sprintf("pgrep %s || echo nf", pc.Process)
	log.Debug().Str("cmd", cmd).Msg("ProcessCheck exec")
	output, err := conn.Execute(cmd)
	if err != nil {
		log.Warn().Str("targetHost", pc.TargetHost).Err(err).Msg("Cannot execute command")
		return nil, derrors.NewInternalError(errors.SSHConnectionError, err)
	}

	processFound := len(output) > 0 && !strings.Contains(string(output), "nf")
	log.Debug().Bool("processFound", processFound).Str("output", string(output)).Msg("")
	if pc.ShouldExists && processFound {
		msg := fmt.Sprintf("Process %s has been found", pc.Process)
		return entities.NewSuccessCommand([]byte(msg)), nil
	}
	if pc.ShouldExists && !processFound {
		msg := fmt.Sprintf("Process %s has not been found and should exist", pc.Process)
		return entities.NewCommandResult(false, msg, nil), nil
	}
	if !pc.ShouldExists && processFound {
		msg := fmt.Sprintf("Process %s has been found and should not exist", pc.Process)
		return entities.NewCommandResult(false, msg, nil), nil
	}
	if !pc.ShouldExists && !processFound {
		msg := fmt.Sprintf("Process %s has not been found", pc.Process)
		return entities.NewSuccessCommand([]byte(msg)), nil
	}

	return entities.NewCommandResult(false, "unexpected combination",
		derrors.NewInternalError(errors.CannotExecuteSyncCommand)), nil
}

// Obtain a string representation
func (pc *ProcessCheck) String() string {
	if pc.ShouldExists {
		return fmt.Sprintf("SYNC ProcessCheck %s %s is running", pc.TargetHost, pc.Process)
	}
	return fmt.Sprintf("SYNC ProcessCheck %s %s is not running ", pc.TargetHost, pc.Process)
}

// PrettyPrint returns a simple space indexed string.
func (pc *ProcessCheck) PrettyPrint(identation int) string {
	return strings.Repeat(" ", identation) + pc.String()
}

// UserString returns a simple string representation of the command for the user.
func (pc *ProcessCheck) UserString() string {
	if pc.ShouldExists {
		return fmt.Sprintf("ProcessCheck %s %s is running", pc.TargetHost, pc.Process)
	}
	return fmt.Sprintf("ProcessCheck %s %s is not running ", pc.TargetHost, pc.Process)
}
