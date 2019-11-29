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
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/installer/internal/pkg/workflow"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

const InstallOperation = "Install cluster"
const UninstallOperation = "Uninstall cluster"

// Operation structure representing an managed operation with its workflow and associated status.
type Operation struct {
	sync.Mutex
	OrganizationID string
	RequestID      string
	OperationName  string
	status         grpc_common_go.OpStatus
	Created        int64
	Params         *workflow.Parameters
	Workflow       *workflow.Workflow
	error          derrors.Error
	workflowState  workflow.WorkflowState
}

// NewOperation creates a new Operation
func NewOperation(organizationID string, requestID string, operationName string) *Operation {
	log.Debug().Str("organizationID", organizationID).Str("requestID", requestID).Str("operationName", operationName).Msg("creating operation")
	return &Operation{
		OrganizationID: organizationID,
		RequestID:      requestID,
		status:         grpc_common_go.OpStatus_INIT,
		Created:        time.Now().Unix(),
		workflowState:  workflow.InitState,
	}
}

func (is *Operation) Clone() *Operation {
	return &Operation{
		OrganizationID: is.OrganizationID,
		RequestID:      is.RequestID,
		OperationName:  is.OperationName,
		status:         is.status,
		Created:        is.Created,
		Params:         is.Params,
		Workflow:       is.Workflow,
		error:          is.error,
		workflowState:  is.workflowState,
	}
}

func (is *Operation) UpdateStatus(newStatus grpc_common_go.OpStatus) {
	is.Lock()
	is.status = newStatus
	is.Unlock()
}

func (is *Operation) GetState() *grpc_common_go.OpStatus {
	is.Lock()
	defer is.Unlock()
	aux := is.status
	return &aux
}

func (is *Operation) UpdateError(error derrors.Error) {
	is.Lock()
	is.error = error
	is.Unlock()
}

func (is *Operation) UpdateWorkflowState(state workflow.WorkflowState) {
	is.Lock()
	is.workflowState = state
	is.Unlock()
}

// ToGRPCOpResponse transforms the information of an install operation in common OpResponse.
func (is *Operation) ToGRPCOpResponse() *grpc_common_go.OpResponse {
	is.Lock()
	rStatus := is.status
	elapsed := time.Now().Unix() - is.Created
	var e string
	if is.error != nil {
		e = is.error.Error()
	}
	is.Unlock()

	return &grpc_common_go.OpResponse{
		OrganizationId: is.OrganizationID,
		RequestId:      is.RequestID,
		OperationName:  is.OperationName,
		ElapsedTime:    elapsed,
		Timestamp:      time.Now().Unix(),
		Status:         rStatus,
		Info:           "",
		Error:          e,
	}
}
