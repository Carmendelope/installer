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

package rke

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/nalej/installer/internal/pkg/errors"
	"github.com/rs/zerolog/log"
	"io"
	"io/ioutil"
	"os/exec"
	"strings"
	"sync"

	"github.com/nalej/derrors"

	"github.com/nalej/installer/internal/pkg/workflow/entities"
	"github.com/nalej/installer/internal/pkg/workflow/handler"
)

// RKERemove structure defining the fields required to uninstall a cluster using RKE.
type RKERemove struct {
	entities.GenericSyncCommand
	RkeBinaryPath string `json:"rkeBinaryPath"`
	ClusterConfig
	installTemplate string
}

// NewRKERemove create a new command with all parameters.
func NewRKERemove(
	rkeBinaryPath string,
	clusterConfig ClusterConfig,
	installTemplate string) *RKERemove {
	return &RKERemove{
		*entities.NewSyncCommand(entities.RKERemove),
		rkeBinaryPath,
		clusterConfig, installTemplate}
}

// NewRKERKERemoveFromJSON creates a RKE Install command from a JSON object.
func NewRKERemoveFromJSON(raw []byte) (*entities.Command, derrors.Error) {
	f := &RKERemove{}
	if err := json.Unmarshal(raw, &f); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	f.CommandID = entities.GenerateCommandID(f.Name())
	var r entities.Command = f
	return &r, nil
}

// getTemplate returns the template to be used for the installation process. If empty, the default one will be used.
func (cmd *RKERemove) getTemplate() string {
	if cmd.installTemplate != "" {
		return cmd.installTemplate
	}
	return ClusterTemplate
}

// CreateClusterConfig generate the RKE cluster.yaml file using the installation parameters.
func (cmd *RKERemove) CreateClusterConfig() (string, derrors.Error) {
	template := NewRKETemplate(cmd.getTemplate())
	config := cmd.ClusterConfig
	yamlString, err := template.ParseTemplate(&config)
	if err != nil {
		return "", err
	}
	clusterFile, createErr := ioutil.TempFile("", "cluster.yaml")
	if createErr != nil {
		return "", derrors.AsError(createErr, errors.IOError)
	}
	if _, writeErr := clusterFile.Write([]byte(yamlString)); writeErr != nil {
		return "", derrors.AsError(writeErr, errors.IOError)
	}
	clusterFile.Close()
	log.Debug().Str("cluster.yaml", clusterFile.Name()).Msg("Temporal cluster.yaml stored")
	return clusterFile.Name(), nil
}

// copyToLog copies a reader output to the associated command handler.
func (cmd *RKERemove) copyToLog(commandHandler handler.CommandHandler, r io.Reader) {
	output := bufio.NewReader(r)
	for {
		line, err := output.ReadString('\n')
		commandHandler.AddLogEntry(cmd.CommandID, strings.TrimSpace(line))
		if err != nil {
			break
		}
	}
}

// Run triggers the execution of the command.
func (cmd *RKERemove) Run(workflowID string) (*entities.CommandResult, derrors.Error) {
	clusterConfigPath, err := cmd.CreateClusterConfig()
	if err != nil {
		return nil, err
	}

	log.Debug().Str("path", cmd.RkeBinaryPath).Msg("RKE binary")
	rke := exec.Command(cmd.RkeBinaryPath, "remove", "--config", clusterConfigPath, "--force")
	rkeOut, pipeErr := rke.StdoutPipe()
	if pipeErr != nil {
		return nil, derrors.AsError(pipeErr, errors.IOError)
	}

	rkeErr, pipeErr := rke.StderrPipe()
	if pipeErr != nil {
		return nil, derrors.AsError(pipeErr, errors.IOError)
	}

	var wg sync.WaitGroup
	commandHandler := handler.GetCommandHandler()
	log.Debug().Msg("Starting rke binary")
	if err := rke.Start(); err != nil {
		return nil, derrors.AsError(err, errors.OpFail)
	}

	wg.Add(2)
	go func() {
		defer wg.Done()
		cmd.copyToLog(commandHandler, rkeOut)
	}()
	go func() {
		defer wg.Done()
		cmd.copyToLog(commandHandler, rkeErr)
	}()

	// Wait for the stdout and stderr pipes to close.
	wg.Wait()
	// Wait for the command itself to close.
	if err := rke.Wait(); err != nil {
		return entities.NewCommandResult(false, "rke failed", derrors.AsError(err, errors.OpFail)), nil
	}
	return entities.NewCommandResult(true, "rke finished successfully", nil), nil
}

// Obtain a string representation
func (cmd *RKERemove) String() string {
	return fmt.Sprintf("SYNC RKE Remove on %s", strings.Join(cmd.TargetNodes, ", "))
}

// PrettyPrint returns a simple space indexed string.
func (cmd *RKERemove) PrettyPrint(indentation int) string {
	binaryPath := strings.Repeat("  ", indentation) + fmt.Sprintf("  RKE binary: %s", cmd.RkeBinaryPath)
	return strings.Repeat(" ", indentation) + fmt.Sprintf("SYNC RKE Remove on %s\n%s",
		strings.Join(cmd.TargetNodes, ", "), binaryPath)
}

// UserString returns a simple string representation of the command for the user.
func (cmd *RKERemove) UserString() string {
	return fmt.Sprintf("Removing Kubernetes on %s ", strings.Join(cmd.TargetNodes, ", "))
}
