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
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-installer-go"
	"github.com/nalej/installer/internal/pkg/workflow"
	"sync"
	"time"
)

type InstallStatus struct {
	sync.Mutex
	InstallID     string
	state         grpc_installer_go.InstallProgress
	Started       int64
	Params        *workflow.Parameters
	Workflow      *workflow.Workflow
	error         derrors.Error
	workflowState workflow.WorkflowState
}

func NewInstallStatus(installID string) *InstallStatus {
	return &InstallStatus{
		InstallID:     installID,
		state:         grpc_installer_go.InstallProgress_REGISTERED,
		Started:       time.Now().Unix(),
		workflowState: workflow.InitState,
	}
}

func (is *InstallStatus) Clone() *InstallStatus {
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

func (is *InstallStatus) UpdateState(installProgress grpc_installer_go.InstallProgress) {
	is.Lock()
	is.state = installProgress
	is.Unlock()
}

func (is *InstallStatus) GetState() *grpc_installer_go.InstallProgress {
	is.Lock()
	defer is.Unlock()
	aux := is.state
	return &aux
}

func (is *InstallStatus) UpdateError(error derrors.Error) {
	is.Lock()
	is.error = error
	is.Unlock()
}

func (is *InstallStatus) UpdateWorkflowState(state workflow.WorkflowState) {
	is.Lock()
	is.workflowState = state
	is.Unlock()
}

func (is *InstallStatus) ToGRPCInstallResponse() *grpc_installer_go.InstallResponse {
	is.Lock()
	rState := is.state
	elapsed := time.Now().Unix() - is.Started
	var e string
	if is.error != nil {
		e = is.error.Error()
	}
	is.Unlock()

	return &grpc_installer_go.InstallResponse{
		InstallId:   is.InstallID,
		State:       rState,
		ElapsedTime: elapsed,
		Error:       e,
	}
}
