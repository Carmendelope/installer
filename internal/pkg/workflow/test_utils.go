/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
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
