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

// Fail command
// Asynchronous implementation of the fail command to simulate and error.
//
// {"type":"async", "name": "fail"}

package async

import (
	"encoding/json"
	"fmt"
	"github.com/nalej/installer/internal/pkg/errors"
	"strings"

	"github.com/nalej/derrors"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
	"github.com/nalej/installer/internal/pkg/workflow/handler"
)

// Fail command structure with supported parameters.
type Fail struct {
	entities.GenericAsyncCommand
}

// NewFail creates an Fail command.
func NewFail() *Fail {
	return &Fail{*entities.NewAsyncCommand(entities.Fail, make([]entities.Action, 0))}
}

// NewFailFromJSON creates an Fail command from a JSON object.
func NewFailFromJSON(raw []byte) (*entities.Command, derrors.Error) {
	f := &Fail{}
	if err := json.Unmarshal(raw, &f); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	f.CommandID = entities.GenerateCommandID(f.Name())
	var r entities.Command = f
	return &r, nil
}

// Run the current command.
//   returns:
//     An error if the command execution fails
func (f *Fail) Run(workflowID string) derrors.Error {
	go f.LogAndFail(workflowID)
	return nil
}

// LogAndFail adds an entry to the log and fails the execution.
func (f *Fail) LogAndFail(workflowID string) {
	cmdHandler := handler.GetCommandHandler()
	cmdHandler.AddLogEntry(f.CommandID, "Asynchronous fail will be triggered")
	result := entities.NewCommandResult(false, "Asynchronous fail", nil)
	cmdHandler.FinishCommand(f.CommandID, result, nil)
}

// String obtains a string representation
func (f *Fail) String() string {
	return fmt.Sprintf("---ASYNC FAIL")
}

// PrettyPrint returns a simple space indexed string.
func (f *Fail) PrettyPrint(indentation int) string {
	return strings.Repeat(" ", indentation) + f.String()
}

// UserString returns a simple string representation of the command for the user.
func (f *Fail) UserString() string {
	return "Fail command"
}
