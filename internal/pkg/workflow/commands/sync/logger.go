/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

// Logger command
// Adds an entry to the workflow log.
//
// {"type":"sync", "name": "logger", "msg": "This is a logging message"}

package sync

import (
	"encoding/json"
	"github.com/nalej/installer/internal/pkg/errors"
	"strings"

	"github.com/nalej/derrors"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
)

// Logger structure with the message to be added.
type Logger struct {
	entities.GenericSyncCommand
	Msg string `json:"msg"`
}

// NewLogger creates a new logger with a message.
func NewLogger(msg string) *Logger {
	return &Logger{*entities.NewSyncCommand(entities.Logger), msg}
}

// NewLoggerFromJSON creates a Logger command from a JSON object.
func NewLoggerFromJSON(raw []byte) (*entities.Command, derrors.Error) {
	logger := &Logger{}
	if err := json.Unmarshal(raw, &logger); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	logger.CommandID = entities.GenerateCommandID(logger.Name())
	var r entities.Command = logger
	return &r, nil
}

// Run the current command.
//   returns:
//     The CommandResult
//     An error if the command execution fails
func (l *Logger) Run(_ string) (*entities.CommandResult, derrors.Error) {
	return entities.NewSuccessCommand([]byte(l.Msg)), nil
}

// String obtains a string representation
func (l *Logger) String() string {
	return "LOG: " + l.Msg
}

// PrettyPrint returns a simple space indexed string.
func (l *Logger) PrettyPrint(identation int) string {
	return strings.Repeat(" ", identation) + l.String()
}

// UserString returns a simple string representation of the command for the user.
func (l *Logger) UserString() string {
	return "adding log entry"
}
