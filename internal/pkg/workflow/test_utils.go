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

package workflow

import (
	"github.com/nalej/derrors"
)

// WorkflowResult structure for testing the callback of a workflow execution.
type WorkflowResult struct {
	Called bool
	Error  derrors.Error
	State  WorkflowState
}

// NewWorkflowResult creates a WorkflowResult.
func NewWorkflowResult() *WorkflowResult {
	return &WorkflowResult{false, nil, InitState}
}

// Finished returns true if the workflow result received the callback.
func (wr *WorkflowResult) Finished() bool {
	return wr.Called
}

// Callback function.
func (wr *WorkflowResult) Callback(workflowID string, error derrors.Error,
	state WorkflowState) {
	wr.Error = error
	wr.Called = true
	wr.State = state
}
