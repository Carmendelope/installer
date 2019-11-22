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

package installer

import (
	"github.com/nalej/grpc-common-go"
	"sync"

	"github.com/nalej/derrors"
	"github.com/nalej/grpc-installer-go"
	"github.com/nalej/installer/internal/pkg/server/config"
	"github.com/nalej/installer/internal/pkg/templates"
	"github.com/nalej/installer/internal/pkg/workflow"
	"github.com/rs/zerolog/log"
)

// Manager structure in charge of orchestrating the install/uninstall processes.
type Manager struct {
	// Mutex to manage the access to the operation related structures.
	sync.Mutex
	// Config with the component configuration.
	Config config.Config
	// Paths with the paths specification to extract binaries, assets, etc.
	Paths workflow.Paths
	// ExecHandler with the workflow executor handler.
	ExecHandler workflow.ExecutorHandler
	// Parser to parametrize templates for execution.
	Parser *workflow.Parser
	// Requests managed
	Requests map[string]grpc_installer_go.InstallRequest
	// Operations with the list of ongoing operations.
	Operations map[string]*Operation
}

// NewManager creates a new installer manager.
func NewManager(config config.Config) Manager {
	return Manager{
		Config:      config,
		Paths:       *workflow.NewPaths(config.ComponentsPath, config.BinaryPath, config.TempPath),
		ExecHandler: workflow.GetExecutorHandler(),
		Parser:      workflow.NewParser(),
		Requests:    make(map[string]grpc_installer_go.InstallRequest, 0),
		Operations:  make(map[string]*Operation, 0),
	}
}

func (m *Manager) unsafeExist(installID string) bool {
	_, exists := m.Operations[installID]
	return exists
}

func (m *Manager) unsafeRegister(installRequest grpc_installer_go.InstallRequest) {
	m.Requests[installRequest.InstallId] = installRequest
	m.Operations[installRequest.InstallId] = NewOperation(installRequest.OrganizationId, installRequest.InstallId)
}

func (m *Manager) InstallCluster(installRequest grpc_installer_go.InstallRequest) (*Operation, derrors.Error) {
	var result *Operation
	m.Lock()
	if m.unsafeExist(installRequest.InstallId) {
		m.Unlock()
		return nil, derrors.NewAlreadyExistsError("installID").WithParams(installRequest.InstallId)
	}
	m.unsafeRegister(installRequest)
	status, _ := m.Operations[installRequest.InstallId]
	result = status.Clone()
	m.Unlock()
	go m.launchInstall(installRequest.InstallId)
	return result, nil
}

func (m *Manager) markInstallAsFailed(installID string, error derrors.Error) {
	m.Lock()
	status, _ := m.Operations[installID]
	status.UpdateError(error)
	status.UpdateStatus(grpc_common_go.OpStatus_FAILED)
	m.Unlock()
}

func (m *Manager) launchInstall(installID string) {
	m.Lock()
	request, exitsRequest := m.Requests[installID]
	status, existStatus := m.Operations[installID]
	m.Unlock()

	if !exitsRequest || !existStatus {
		log.Error().Str("installID", installID).Msg("cannot launch the install process")
		return
	}

	// Create Parameters
	params := workflow.NewInstallParameters(
		&request, workflow.Assets{}, m.Paths,
		m.Config.ManagementClusterHost, m.Config.ManagementClusterPort,
		m.Config.DNSClusterHost, m.Config.DNSClusterPort,
		m.Config.Environment.Target,
		true,
		*workflow.EmptyNetworkConfig, m.Config.AuthSecret, m.Config.ClusterCertIssuerCACertPath)

	status.Params = params
	err := status.Params.LoadCredentials()
	if err != nil {
		log.Error().Str("err", err.DebugReport()).Msg("cannot load credentials")
		m.markInstallAsFailed(installID, err)
	}
	err = status.Params.Validate()
	if err != nil {
		log.Error().Str("err", err.DebugReport()).Msg("invalid parameters")
		m.markInstallAsFailed(installID, err)
	}

	// Create Workflow
	workflow, err := m.Parser.ParseWorkflow(installID, templates.InstallManagementCluster, installID, *status.Params)
	if err != nil {
		log.Error().Str("err", err.DebugReport()).Msg("cannot parse workflow")
		m.markInstallAsFailed(installID, err)
	}
	status.Workflow = workflow

	// Launch install process
	exec, err := m.ExecHandler.Add(status.Workflow, m.WorkflowCallback)
	if err != nil {
		log.Error().Str("err", err.DebugReport()).Msg("cannot parse workflow")
		m.markInstallAsFailed(installID, err)
	}
	exec.SetLogListener(m.logListener)
	exec.Exec()
}

func (m *Manager) GetProgress(requestID string) (*Operation, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExist(requestID) {
		return nil, derrors.NewNotFoundError("installID").WithParams(requestID)
	}
	status, _ := m.Operations[requestID]
	log.Debug().Interface("status", status).Msg("GetProgress()")
	return status.Clone(), nil
}

func (m *Manager) WorkflowCallback(
	workflowID string,
	error derrors.Error,
	state workflow.WorkflowState) {
	log.Debug().Str("workflowID", workflowID).Err(error).Interface("state", state).Msg("WorkflowCallback()")

	m.Lock()
	defer m.Unlock()
	status, exist := m.Operations[workflowID]
	if !exist {
		log.Warn().Str("workflowID", workflowID).Msg("received callback for unregistered workflow")
	}
	if error != nil {
		status.UpdateStatus(grpc_common_go.OpStatus_FAILED)
	}
	status.UpdateWorkflowState(state)
	switch state {
	case workflow.InitState:
		log.Warn().Msg("Not expecting init update")
		return
	case workflow.RegisteredState:
		status.UpdateStatus(grpc_common_go.OpStatus_SCHEDULED)
		return
	case workflow.InProgressState:
		status.UpdateStatus(grpc_common_go.OpStatus_INPROGRESS)
		return
	case workflow.FinishedState:
		status.UpdateStatus(grpc_common_go.OpStatus_SUCCESS)
		return
	case workflow.ErrorState:
		status.UpdateStatus(grpc_common_go.OpStatus_FAILED)
	default:
		log.Warn().Interface("state", state).Msg("State not recognized")
	}
}

func (m *Manager) logListener(msg string) {
	// TODO store the information on the install status
	log.Info().Msg(msg)
}

func (m *Manager) RemoveInstall(installID string) derrors.Error {
	m.Lock()
	_, exitsRequest := m.Requests[installID]
	if exitsRequest {
		log.Debug().Str("installID", installID).Msg("Removing request")
		delete(m.Requests, installID)
	}
	_, existStatus := m.Operations[installID]
	m.Unlock()

	if existStatus {
		err := m.ExecHandler.Stop(installID)
		if err != nil {
			return err
		}
		m.Lock()
		delete(m.Operations, installID)
		m.Unlock()
	}

	return nil
}

func (m *Manager) UninstallCluster(request *grpc_installer_go.UninstallClusterRequest) (*Operation, derrors.Error) {
	return nil, derrors.NewUnimplementedError("uninstall not implemented")
}
