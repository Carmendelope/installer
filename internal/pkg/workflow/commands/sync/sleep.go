/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
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
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
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
