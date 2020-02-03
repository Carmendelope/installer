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
	"github.com/nalej/installer/internal/pkg/entities"
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
	// InstallRequest by request identifier
	InstallRequests map[string]grpc_installer_go.InstallRequest
	// UninstallRequest by request identifier
	UninstallRequests map[string]grpc_installer_go.UninstallClusterRequest
	// Operations with the list of ongoing operations.
	Operations map[string]*Operation
}

// NewManager creates a new installer manager.
func NewManager(config config.Config) Manager {
	return Manager{
		Config:            config,
		Paths:             *workflow.NewPaths(config.ComponentsPath, config.BinaryPath, config.TempPath),
		ExecHandler:       workflow.GetExecutorHandler(),
		Parser:            workflow.NewParser(),
		InstallRequests:   make(map[string]grpc_installer_go.InstallRequest, 0),
		UninstallRequests: make(map[string]grpc_installer_go.UninstallClusterRequest, 0),
		Operations:        make(map[string]*Operation, 0),
	}
}

func (m *Manager) unsafeExist(requestID string) bool {
	_, exists := m.Operations[requestID]
	return exists
}

func (m *Manager) unsafeInstallRegister(installRequest grpc_installer_go.InstallRequest) {
	m.InstallRequests[installRequest.RequestId] = installRequest
	m.Operations[installRequest.RequestId] = NewOperation(installRequest.OrganizationId, installRequest.RequestId, InstallOperation)
}

func (m *Manager) unsafeUninstallRegister(request grpc_installer_go.UninstallClusterRequest) {
	m.UninstallRequests[request.RequestId] = request
	m.Operations[request.RequestId] = NewOperation(request.OrganizationId, request.RequestId, UninstallOperation)
}

func (m *Manager) InstallCluster(installRequest grpc_installer_go.InstallRequest) (*Operation, derrors.Error) {
	var result *Operation
	m.Lock()
	if m.unsafeExist(installRequest.RequestId) {
		m.Unlock()
		return nil, derrors.NewAlreadyExistsError("requestID").WithParams(installRequest.RequestId)
	}
	m.unsafeInstallRegister(installRequest)
	status, _ := m.Operations[installRequest.RequestId]
	result = status.Clone()
	m.Unlock()
	go m.launchInstall(installRequest.RequestId)
	return result, nil
}

func (m *Manager) markOperationAsFailed(requestID string, error derrors.Error) {
	m.Lock()
	status, _ := m.Operations[requestID]
	status.UpdateError(error)
	status.UpdateStatus(grpc_common_go.OpStatus_FAILED)
	m.Unlock()
}

func (m *Manager) launchInstall(requestID string) {
	m.Lock()
	request, exitsRequest := m.InstallRequests[requestID]
	status, existStatus := m.Operations[requestID]
	m.Unlock()

	if !exitsRequest || !existStatus {
		log.Error().Str("requestID", requestID).Msg("cannot launch the install process")
		return
	}

	// The network configuration is taken from the running parameters of the installer service
	networkingConfig := workflow.NetworkConfig{
		NetworkingMode:     entities.NetworkingModeToString[m.Config.NetworkingMode],
		IstioPath:          m.Config.IstioPath,
		ZTPlanetSecretPath: "",
	}

	// Create Parameters
	params := workflow.NewInstallParameters(
		&request, workflow.Assets{}, m.Paths,
		m.Config.ManagementClusterHost, m.Config.ManagementClusterPort,
		m.Config.DNSClusterHost, m.Config.DNSClusterPort,
		m.Config.Environment.Target,
		true,
		networkingConfig, m.Config.AuthSecret, m.Config.ClusterCertIssuerCACertPath)

	status.Params = params
	err := status.Params.LoadCredentials()
	if err != nil {
		log.Error().Str("err", err.DebugReport()).Msg("cannot load credentials")
		m.markOperationAsFailed(requestID, err)
	}
	err = status.Params.Validate()
	if err != nil {
		log.Error().Str("err", err.DebugReport()).Msg("invalid parameters")
		m.markOperationAsFailed(requestID, err)
	}

	// Create Workflow
	workflow, err := m.Parser.ParseWorkflow(requestID, templates.InstallManagementCluster, requestID, *status.Params)
	if err != nil {
		log.Error().Str("err", err.DebugReport()).Msg("cannot parse workflow")
		m.markOperationAsFailed(requestID, err)
	}
	status.Workflow = workflow

	// Launch install process
	exec, err := m.ExecHandler.Add(status.Workflow, m.WorkflowCallback)
	if err != nil {
		log.Error().Str("err", err.DebugReport()).Msg("cannot parse workflow")
		m.markOperationAsFailed(requestID, err)
	}
	exec.SetLogListener(m.logListener)
	exec.Exec()
}

func (m *Manager) GetProgress(requestID string) (*Operation, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExist(requestID) {
		return nil, derrors.NewNotFoundError("requestID").WithParams(requestID)
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

func (m *Manager) RemoveInstall(requestID string) derrors.Error {
	m.Lock()
	// Determine the type of operation
	op, existsOp := m.Operations[requestID]
	if !existsOp {
		return derrors.NewNotFoundError("request is not managed by the installer").WithParams(requestID)
	}

	if op.OperationName == InstallOperation {
		_, exitsRequest := m.InstallRequests[requestID]
		if exitsRequest {
			log.Debug().Str("requestID", requestID).Msg("Removing install request")
			delete(m.InstallRequests, requestID)
		}
	} else if op.OperationName == UninstallOperation {
		_, exitsRequest := m.UninstallRequests[requestID]
		if exitsRequest {
			log.Debug().Str("requestID", requestID).Msg("Removing uninstall request")
			delete(m.InstallRequests, requestID)
		}
	}
	m.Unlock()

	if existsOp {
		err := m.ExecHandler.Stop(requestID)
		if err != nil {
			return err
		}
		m.Lock()
		delete(m.Operations, requestID)
		m.Unlock()
	}

	return nil
}

func (m *Manager) UninstallCluster(request grpc_installer_go.UninstallClusterRequest) (*Operation, derrors.Error) {
	var result *Operation
	m.Lock()
	if m.unsafeExist(request.RequestId) {
		m.Unlock()
		return nil, derrors.NewAlreadyExistsError("requestID").WithParams(request.RequestId)
	}
	m.unsafeUninstallRegister(request)
	status, _ := m.Operations[request.RequestId]
	result = status.Clone()
	m.Unlock()
	go m.launchUninstall(request.RequestId)
	return result, nil
}

func (m *Manager) launchUninstall(requestID string) {
	m.Lock()
	request, exitsRequest := m.UninstallRequests[requestID]
	status, existStatus := m.Operations[requestID]
	m.Unlock()

	if !exitsRequest || !existStatus {
		log.Error().Str("requestID", requestID).Msg("cannot launch the uninstall process")
		return
	}

	params := workflow.NewUninstallParameters(&request, true)

	status.Params = params
	err := status.Params.LoadCredentials()
	if err != nil {
		log.Error().Str("err", err.DebugReport()).Msg("cannot load credentials")
		m.markOperationAsFailed(requestID, err)
	}
	err = status.Params.Validate()
	if err != nil {
		log.Error().Str("err", err.DebugReport()).Msg("invalid parameters")
		m.markOperationAsFailed(requestID, err)
	}

	// Create Workflow
	workflow, err := m.Parser.ParseWorkflow(requestID, templates.UninstallCluster, requestID, *status.Params)
	if err != nil {
		log.Error().Str("err", err.DebugReport()).Msg("cannot parse workflow")
		m.markOperationAsFailed(requestID, err)
	}
	status.Workflow = workflow

	// Launch install process
	exec, err := m.ExecHandler.Add(status.Workflow, m.WorkflowCallback)
	if err != nil {
		log.Error().Str("err", err.DebugReport()).Msg("cannot parse workflow")
		m.markOperationAsFailed(requestID, err)
	}
	exec.SetLogListener(m.logListener)
	exec.Exec()
}
