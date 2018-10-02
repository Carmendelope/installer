/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

// Copy command
// Copies a file to a remote host using SCP.
//
// {"type":"sync", "name": "scp", "targetHost": "127.0.0.1", "targetPort": "22",
// "credentials":{"username": "username", "password":"passwd"},
// "source":"script.sh", "destination":"/opt/scripts/."]}

package sync

import (
	"encoding/json"
	"fmt"
	"github.com/nalej/installer/internal/pkg/errors"
	"strings"
	"time"

	"github.com/nalej/derrors"
	"github.com/nalej/installer/internal/pkg/workflow/commands/sync/connection"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
)

// DefaultSSHPort defines the default port for SSH connections.
const DefaultSSHPort = "22"

// SCP command structure with supported fields.
type SCP struct {
	entities.GenericSyncCommand
	// Target node
	TargetHost string `json:"targetHost"`
	// Target port
	TargetPort string `json:"targetPort"`
	// Credentials for SSH.
	Credentials entities.Credentials `json:"credentials"`
	// Source path
	Source string `json:"source"`
	// Destination path
	Destination string `json:"destination"`
}

// NewSCP creates an SCP command from a set of parameters.
func NewSCP(targetHost string, targetPort string, credentials entities.Credentials, source string, destination string) *SCP {
	return &SCP{*entities.NewSyncCommand(entities.SCP),
		targetHost,
		targetPort,
		credentials,
		source,
		destination}
}

// NewSCPFromJSON creates an SCP command from a JSON object.
func NewSCPFromJSON(raw []byte) (*entities.Command, derrors.Error) {
	scp := &SCP{}
	if err := json.Unmarshal(raw, &scp); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	scp.CommandID = entities.GenerateCommandID(scp.Name())
	var r entities.Command = scp
	return &r, nil
}

func (scp *SCP) getTargetPort() string {
	if scp.TargetPort != "" {
		return scp.TargetPort
	}
	return DefaultSSHPort
}

// Run the current command.
//   returns:
//     The CommandResult
//     An error if the command execution fails
func (scp *SCP) Run(_ string) (*entities.CommandResult, derrors.Error) {

	conn, err := connection.NewSSHConnection(
		scp.TargetHost, scp.getTargetPort(),
		scp.Credentials.Username, scp.Credentials.Password, "", scp.Credentials.PrivateKey)
	if err != nil {
		return nil, derrors.NewInternalError(errors.SSHConnectionError, err).WithParams(scp.TargetHost)
	}
	start := time.Now()
	err = conn.Copy(scp.Source, scp.Destination, false)
	if err != nil {
		return nil, derrors.NewInternalError(errors.SSHConnectionError, err).WithParams(scp.TargetHost)
	}

	return entities.NewSuccessCommand([]byte(scp.String() + ": OK " + time.Since(start).String())), nil
}

// String obtains a string representation
func (scp *SCP) String() string {
	return fmt.Sprintf("SYNC SCP %s %s:%s", scp.Source, scp.TargetHost, scp.Destination)
}

// PrettyPrint returns a simple space indexed string.
func (scp *SCP) PrettyPrint(indentation int) string {
	return strings.Repeat(" ", indentation) + scp.String()
}

// UserString returns a simple string representation of the command for the user.
func (scp *SCP) UserString() string {
	return fmt.Sprintf("SCP %s %s:%s", scp.Source, scp.TargetHost, scp.Destination)
}
