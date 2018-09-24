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

// Sleep command
// Sleeps for given set of seconds.
//
// {"type":"sync", "name": "sleep", "time": "2"}

package sync

import (
	"encoding/json"
	"github.com/nalej/installer/internal/pkg/errors"
	"strconv"
	"strings"
	"time"

	"github.com/nalej/derrors"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
)

// Sleep structure with the time to sleep.
type Sleep struct {
	entities.GenericSyncCommand
	Time string `json:"time"`
}

// NewSleep creates a new sleep command.
func NewSleep(time string) *Sleep {
	return &Sleep{*entities.NewSyncCommand(entities.Sleep), time}
}

// NewSleepFromJSON creates a Sleep command from a JSON object.
func NewSleepFromJSON(raw []byte) (*entities.Command, derrors.Error) {
	sleep := &Sleep{}
	if err := json.Unmarshal(raw, &sleep); err != nil {
		return nil, derrors.NewOperationError(errors.UnmarshalError, err)
	}
	sleep.CommandID = entities.GenerateCommandID(sleep.Name())
	var r entities.Command = sleep
	return &r, nil
}

// Run the current command.
//   returns:
//     The CommandResult
//     An error if the command execution fails
func (s *Sleep) Run(_ string) (*entities.CommandResult, derrors.Error) {
	t, _ := strconv.Atoi(s.Time)
	d := time.Duration(t)
	time.Sleep(time.Second * d)
	return entities.NewSuccessCommand([]byte("slept for " + s.Time)), nil
}

// String obtains a string representation
func (s *Sleep) String() string {
	return "SLEEP: " + s.Time
}

// PrettyPrint returns a simple space indexed string.
func (s *Sleep) PrettyPrint(identation int) string {
	return strings.Repeat(" ", identation) + s.String()
}

// UserString returns a simple string representation of the command for the user.
func (s *Sleep) UserString() string {
	return "Sleeping for " + s.Time
}
