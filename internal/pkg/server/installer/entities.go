/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package installer

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-installer-go"
	"github.com/nalej/installer/internal/pkg/workflow"
	"sync"
	"time"
)

type InstallStatus struct {
	sync.Mutex
	InstallID string
	state grpc_installer_go.InstallProgress
	Started int64
	Params * workflow.Parameters
	Workflow * workflow.Workflow
	error derrors.Error
	workflowState workflow.WorkflowState
}

func NewInstallStatus(installID string) * InstallStatus {
	return &InstallStatus{
		InstallID: installID,
		state:     grpc_installer_go.InstallProgress_REGISTERED,
		Started:   time.Now().Unix(),
		workflowState: workflow.InitState,
	}
}

func (is * InstallStatus) Clone() * InstallStatus {
	return &InstallStatus{
		InstallID:     is.InstallID,
		state:         is.state,
		Started:       is.Started,
		Params:        is.Params,
		Workflow:      is.Workflow,
		error:         is.error,
		workflowState: is.workflowState,
	}
}

func (is * InstallStatus) UpdateState(installProgress grpc_installer_go.InstallProgress) {
	is.Lock()
	is.state = installProgress
	is.Unlock()
}

func (is * InstallStatus) GetState() * grpc_installer_go.InstallProgress {
	is.Lock()
	defer is.Unlock()
	aux := is.state
	return &aux
}

func (is * InstallStatus) UpdateError(error derrors.Error) {
	is.Lock()
	is.error = error
	is.Unlock()
}

func (is * InstallStatus) UpdateWorkflowState(state workflow.WorkflowState) {
	is.Lock()
	is.workflowState = state
	is.Unlock()
}

func (is * InstallStatus) ToGRPCInstallResponse() *grpc_installer_go.InstallResponse {
	is.Lock()
	rState := is.state
	elapsed := time.Now().Unix() - is.Started
	var e string
	if is.error != nil {
		e = is.error.Error()
	}
	is.Unlock()

	return &grpc_installer_go.InstallResponse{
		InstallId:            is.InstallID,
		State:                rState,
		ElapsedTime:          elapsed,
		Error:                e,
	}
}