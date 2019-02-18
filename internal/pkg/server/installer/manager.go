/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package installer

import (
	"sync"

	"github.com/nalej/derrors"
	"github.com/nalej/grpc-installer-go"
	"github.com/nalej/installer/internal/pkg/server/config"
	"github.com/nalej/installer/internal/pkg/templates"
	"github.com/nalej/installer/internal/pkg/workflow"
	"github.com/rs/zerolog/log"
)

type Manager struct {
	sync.Mutex
	Config      config.Config
	Paths       workflow.Paths
	ExecHandler workflow.ExecutorHandler
	Parser      *workflow.Parser
	Requests    map[string]grpc_installer_go.InstallRequest
	Status      map[string]*InstallStatus
}

func NewManager(config config.Config) Manager {
	return Manager{
		Config:      config,
		Paths:       *workflow.NewPaths(config.ComponentsPath, config.BinaryPath, config.TempPath),
		ExecHandler: workflow.GetExecutorHandler(),
		Parser:      workflow.NewParser(),
		Requests:    make(map[string]grpc_installer_go.InstallRequest, 0),
		Status:      make(map[string]*InstallStatus, 0),
	}
}

func (m *Manager) unsafeExist(installID string) bool {
	_, exists := m.Status[installID]
	return exists
}

func (m *Manager) unsafeRegister(installRequest grpc_installer_go.InstallRequest) {
	m.Requests[installRequest.InstallId] = installRequest
	m.Status[installRequest.InstallId] = NewInstallStatus(installRequest.InstallId)
}

func (m *Manager) InstallCluster(installRequest grpc_installer_go.InstallRequest) (*InstallStatus, derrors.Error) {
	var result *InstallStatus
	m.Lock()
	if m.unsafeExist(installRequest.InstallId) {
		m.Unlock()
		return nil, derrors.NewAlreadyExistsError("installID").WithParams(installRequest.InstallId)
	}
	m.unsafeRegister(installRequest)
	status, _ := m.Status[installRequest.InstallId]
	result = status.Clone()
	m.Unlock()
	go m.launchInstall(installRequest.InstallId)
	return result, nil
}

func (m *Manager) markInstallAsFailed(installID string, error derrors.Error) {
	m.Lock()
	status, _ := m.Status[installID]
	status.UpdateError(error)
	status.UpdateState(grpc_installer_go.InstallProgress_ERROR)
	m.Unlock()
}

func (m *Manager) launchInstall(installID string) {
	m.Lock()
	request, exitsRequest := m.Requests[installID]
	status, existStatus := m.Status[installID]
	m.Unlock()

	if !exitsRequest || !existStatus {
		log.Error().Str("installID", installID).Msg("cannot launch the install process")
		return
	}

	registryCredentials := workflow.NewRegistryCredentialsFromEnvironment(m.Config.Environment)

	// Create Parameters
	params := workflow.NewParameters(
		request, workflow.Assets{}, m.Paths,
		m.Config.ManagementClusterHost, m.Config.ManagementClusterPort,
		m.Config.DNSClusterHost, m.Config.DNSClusterPort,
		m.Config.Environment.Target,
		true,
		registryCredentials)
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

func (m *Manager) GetProgress(installID string) (*InstallStatus, derrors.Error) {
	m.Lock()
	defer m.Unlock()
	if !m.unsafeExist(installID) {
		return nil, derrors.NewNotFoundError("installID").WithParams(installID)
	}
	status, _ := m.Status[installID]
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
	status, exist := m.Status[workflowID]
	if !exist {
		log.Warn().Str("workflowID", workflowID).Msg("received callback for unregistered workflow")
	}
	if error != nil {
		status.UpdateState(grpc_installer_go.InstallProgress_ERROR)
	}
	status.UpdateWorkflowState(state)
	switch state {
	case workflow.InitState:
		log.Warn().Msg("Not expecting init update")
		return
	case workflow.RegisteredState:
		status.UpdateState(grpc_installer_go.InstallProgress_REGISTERED)
		return
	case workflow.InProgressState:
		status.UpdateState(grpc_installer_go.InstallProgress_IN_PROGRESS)
		return
	case workflow.FinishedState:
		status.UpdateState(grpc_installer_go.InstallProgress_FINISHED)
		return
	case workflow.ErrorState:
		status.UpdateState(grpc_installer_go.InstallProgress_ERROR)
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
	_, existStatus := m.Status[installID]
	m.Unlock()

	if existStatus {
		err := m.ExecHandler.Stop(installID)
		if err != nil {
			return err
		}
		m.Lock()
		delete(m.Status, installID)
		m.Unlock()
	}

	return nil
}
