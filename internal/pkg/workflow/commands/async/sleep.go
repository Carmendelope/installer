/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

// Sleep command
// Asynchronous sleep implementation for given set of seconds.
//
// {"type":"async", "name": "sleep", "time": "2"}

package async

import (
	"encoding/json"
	"github.com/nalej/installer/internal/pkg/errors"
	"strconv"
	"strings"
	"time"

	"github.com/nalej/derrors"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
	"github.com/nalej/installer/internal/pkg/workflow/handler"
)

// Sleep structure with the time to sleep.
type Sleep struct {
	entities.GenericAsyncCommand
	Time string `json:"time"`
}

// NewSleep creates a new sleep command.
func NewSleep(time string) *Sleep {
	return &Sleep{*entities.NewAsyncCommand(entities.Sleep, make([]entities.Action, 0)), time}
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
//     An error if the command execution fails
func (s *Sleep) Run(workflowID string) derrors.Error {
	go s.sleepAndNotify(workflowID)
	return nil
}

func (s *Sleep) sleepAndNotify(workflowID string) {
	t, _ := strconv.Atoi(s.Time)
	d := time.Duration(t)
	cmdHandler := handler.GetCommandHandler()
	cmdHandler.AddLogEntry(s.CommandID, "Asynchronous sleep command")
	time.Sleep(time.Second * d)
	result := entities.NewCommandResult(true, "Slept for "+s.Time, nil)
	cmdHandler.FinishCommand(s.CommandID, result, nil)
}

// String obtains a string representation
func (s *Sleep) String() string {
	return "SLEEP: " + s.Time
}

// PrettyPrint returns a simple space indexed string.
func (s *Sleep) PrettyPrint(indentation int) string {
	return strings.Repeat(" ", indentation) + s.String()
}

// UserString returns a simple string representation of the command for the user.
func (s *Sleep) UserString() string {
	return "Sleeping for " + s.Time
}
