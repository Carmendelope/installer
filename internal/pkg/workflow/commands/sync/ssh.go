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

// SSH command
// Executes commands in a remote host.
//
// {"type":"sync", "name": "ssh", "targetHost": "127.0.0.1", "targetPort": "22",
// "credentials":{"username": "username", "password":"passwd"},
// "cmd":"script.sh", "args":["args1", "arg2"]}
//
// For PKI auth, specify privateKey in the credentials object.

package sync

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/nalej/installer/internal/pkg/errors"
	"github.com/rs/zerolog/log"
	"strings"

	"github.com/nalej/derrors"
	"github.com/nalej/installer/internal/pkg/workflow/commands/sync/connection"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
)

// SSH command structure with supported parameters.
type SSH struct {
	entities.GenericSyncCommand
	// Target node
	TargetHost string `json:"targetHost"`
	// Target port
	TargetPort string `json:"targetPort"`
	// Credentials for SSH.
	Credentials entities.Credentials `json:"credentials"`
	// Command to be execute
	Cmd string `json:"cmd"`
	// Command arguments
	Args []string `json:"args"`
}

// NewSSH creates an SSH command from a set of parameters.
func NewSSH(targetHost string, targetPort string, credentials entities.Credentials, cmd string, args []string) *SSH {
	return &SSH{*entities.NewSyncCommand(entities.SSH),
		targetHost,
		targetPort,
		credentials,
		cmd,
		args}
}

// NewSSHFromJSON creates an SSH command from a JSON object.
func NewSSHFromJSON(raw []byte) (*entities.Command, derrors.Error) {
	ssh := &SSH{}
	if err := json.Unmarshal(raw, &ssh); err != nil {
		return nil, derrors.NewOperationError(errors.UnmarshalError, err)
	}
	ssh.CommandID = entities.GenerateCommandID(ssh.Name())
	var r entities.Command = ssh
	return &r, nil
}

func (ssh *SSH) getTargetPort() string {
	if ssh.TargetPort != "" {
		return ssh.TargetPort
	}
	return DefaultSSHPort
}

// Run the current command.
//   returns:
//     The CommandResult
//     An error if the command execution fails
func (ssh *SSH) Run(_ string) (*entities.CommandResult, derrors.Error) {

	conn, err := connection.NewSSHConnection(
		ssh.TargetHost, ssh.getTargetPort(),
		ssh.Credentials.Username, ssh.Credentials.Password, "", ssh.Credentials.PrivateKey)
	if err != nil {
		log.Warn().Str("targetHost", ssh.TargetHost).Err(err).Msg("Cannot establish connection ")
		return nil, derrors.NewConnectionError(errors.SSHConnectionError, err)
	}
	var buffer bytes.Buffer
	buffer.WriteString(ssh.Cmd)
	for _, arg := range ssh.Args {
		buffer.WriteString(" " + arg)
	}
	toExecute := buffer.String()
	log.Debug().Str("toExecute", toExecute).Msg("SSH exec")
	output, err := conn.Execute(toExecute)
	if err != nil {
		log.Warn().Str("targetHost", ssh.TargetHost).Err(err).Msg("Cannot execute command")
		return nil, derrors.NewConnectionError(errors.SSHConnectionError, err)
	}

	return entities.NewSuccessCommand(output), nil
}

// Obtain a string representation
func (ssh *SSH) String() string {
	return fmt.Sprintf("SYNC SSH %s %s %s", ssh.TargetHost, ssh.Cmd, strings.Join(ssh.Args, " "))
}

// PrettyPrint returns a simple space indexed string.
func (ssh *SSH) PrettyPrint(indentation int) string {
	return strings.Repeat(" ", indentation) + ssh.String()
}

// UserString returns a simple string representation of the command for the user.
func (ssh *SSH) UserString() string {
	return fmt.Sprintf("SSH %s %s %s", ssh.TargetHost, ssh.Cmd, strings.Join(ssh.Args, " "))
}
