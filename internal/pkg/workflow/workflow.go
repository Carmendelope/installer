/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package workflow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/nalej/installer/internal/pkg/errors"
	"github.com/rs/zerolog/log"

	"github.com/nalej/derrors"
	"github.com/nalej/installer/internal/pkg/workflow/commands"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
)

// Prefix for workflow identifiers.
const Prefix = "wf-"

// WorkflowState defines the different states of a workflow.
type WorkflowState string

// InitState is the first state of a workflow.
const InitState WorkflowState = "init"

// RegisteredState represents a workflow that is on the queue to be executed.
const RegisteredState WorkflowState = "registered"

// InProgressState represents a workflow that is currently running.
const InProgressState WorkflowState = "in-progress"

// ErrorState represents a workflow that failed during the execution.
const ErrorState WorkflowState = "error"

// FinishedState represents a workflow that has finished.
const FinishedState WorkflowState = "finished"

// Workflow defines a basic structure for a pipeline workflow definition.
type Workflow struct {
	// WorkflowID contains the workflow identifier.
	WorkflowID string `json:"id"`
	// Name of the workflow.
	Name string `json:"name"`
	// Description of the workflow
	Description string `json:"description"`
	// Commands that are going to be executed.
	Commands []entities.Command `json:"commands"`
}

// NewWorkflow creates a new workflow.
//   params:
//     name The workflow name.
//     description The workflow description.
//     commands The list of commands.
//   returns:
//     A GenericWorkflow.
func NewWorkflow(workflowID string, name string, description string, commands []entities.Command) *Workflow {
	return &Workflow{
		WorkflowID:  workflowID,
		Name:        name,
		Description: description,
		Commands:    commands,
	}

}

// PrettyPrint creates a string with the debug information of this workflow.
func (w *Workflow) PrettyPrint() string {
	var buffer bytes.Buffer
	buffer.WriteString("WorkflowID: " + w.WorkflowID)
	buffer.WriteString("\nName: " + w.Name)
	buffer.WriteString("\nDescription: " + w.Description + "\n")
	for index, cmd := range w.Commands {
		buffer.WriteString(fmt.Sprintf("%d) - %s\n", index, cmd.PrettyPrint(0)))
	}
	return buffer.String()
}

// StatusResponse structure to report on workflow execution.
type StatusResponse struct {
	WorkflowID     string        `json:"workflowId"`
	State          WorkflowState `json:"state"`
	CurrentCommand int           `json:"currentCommand"`
	NumCommands    int           `json:"numCommands"`
}

// ToStatusResponse creates a StatusResponse from a executor.
func ToStatusResponse(exe *Executor) *StatusResponse {
	curr, total := exe.CurrentCommand()
	return NewStatusResponse(exe.WorkflowID, exe.State, curr+1, total)
}

// NewStatusResponse creates a new StatusResponse.
func NewStatusResponse(
	workflowID string,
	state WorkflowState,
	currentCommand int,
	numCommands int) *StatusResponse {
	return &StatusResponse{workflowID, state, currentCommand, numCommands}
}

// LogEntry structure for remote log manipulation.
type LogEntry struct {
	Msg string `json:"msg"`
}

// NewLogEntry creates a new LogEntry with a given message.
func NewLogEntry(msg string) *LogEntry {
	return &LogEntry{msg}
}

// LogResponse structure to respond HTTP requests.
type LogResponse struct {
	WorkflowID string   `json:"workflowId"`
	Log        []string `json:"log"`
}

// NewLogResponse creates a new log response.
func NewLogResponse(workflowID string, log []string) *LogResponse {
	return &LogResponse{workflowID, log}
}

// WorkflowParameter structure for live parameters sets by the commands.
type WorkflowParameter struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// NewWorkflowParameter creates a WorkflowParameter.
func NewWorkflowParameter(key string, value string) *WorkflowParameter {
	return &WorkflowParameter{key, value}
}

// ParseWorkflowRequest is the structure sent by the clients to parse a given template.
type ParseWorkflowRequest struct {
	Name       string     `json:"name"`
	Template   string     `json:"template"`
	Parameters Parameters `json:"parameters"`
}

// NewParseWorkflowRequest creates a new ParseWorkflowRequest.
func NewParseWorkflowRequest(name string, template string, parameters Parameters) *ParseWorkflowRequest {
	return &ParseWorkflowRequest{Name: name, Template: template, Parameters: parameters}
}

// WorkflowFromJSON structure with RawMessages as commands to make possible proper deserialization.
type WorkflowFromJSON struct {
	// WorkflowID contains the workflow identifier.
	WorkflowID string `json:"id"`
	// Name of the workflow.
	Name string `json:"name"`
	// Description of the workflow
	Description string `json:"description"`
	// Commands that are going to be executed.
	Commands []json.RawMessage `json:"commands"`
}

// ToWorkflow transforms the current structure into a workflow by parsing individual parameters.
func (wfj *WorkflowFromJSON) ToWorkflow() (*Workflow, derrors.Error) {

	p := commands.NewCmdParser()
	result := make([]entities.Command, 0)
	for index, raw := range wfj.Commands {
		log.Debug().Int("index", index).Str("raw", string(raw)).Msg("processing raw command")
		var gc entities.GenericCommand
		if err := json.Unmarshal(raw, &gc); err != nil {
			return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
		}

		cmd, err := p.ParseCommand(raw)
		if err != nil {
			return nil, err
		}
		result = append(result, *cmd)
	}

	return &Workflow{
		wfj.WorkflowID,
		wfj.Name,
		wfj.Description,
		result}, nil
}
