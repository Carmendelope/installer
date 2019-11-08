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
// Fails the execution of the workflow.
//
// {"type":"sync", "name": "fail"}

package sync

import (
	"encoding/json"
	"fmt"
	"github.com/nalej/installer/internal/pkg/errors"
	"strings"

	"github.com/nalej/derrors"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
)

// Fail command structure with supported parameters.
type Fail struct {
	entities.GenericSyncCommand
}

// NewFail creates an Fail command.
func NewFail() *Fail {
	return &Fail{*entities.NewSyncCommand(entities.Fail)}
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
//     The CommandResult
//     An error if the command execution fails
func (f *Fail) Run(workflowID string) (*entities.CommandResult, derrors.Error) {
	return entities.NewErrCommand("fail command - "+workflowID, derrors.NewGenericError("forced failure")), nil
}

// String obtains a string representation
func (f *Fail) String() string {
	return fmt.Sprintf("---FAIL")
}

// PrettyPrint returns a simple space indexed string.
func (f *Fail) PrettyPrint(identation int) string {
	return strings.Repeat(" ", identation) + f.String()
}

// UserString returns a simple string representation of the command for the user.
func (f *Fail) UserString() string {
	return "Fail command"
}
