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

// This file contains the command parsing facilities to avoid import cycles.

package commands

import (
	"encoding/json"
	"github.com/nalej/installer/internal/pkg/errors"

	"github.com/nalej/derrors"
	"github.com/nalej/installer/internal/pkg/workflow/commands/async"
	"github.com/nalej/installer/internal/pkg/workflow/commands/sync"
	"github.com/nalej/installer/internal/pkg/workflow/commands/sync/rke"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
)

// CmdParser structure for the command parsing.
type CmdParser struct {
}

// NewCmdParser creates a new command parser.
func NewCmdParser() *CmdParser {
	return &CmdParser{}
}

// ParseCommand extracts a command from a raw JSON message.
func (cp *CmdParser) ParseCommand(raw []byte) (*entities.Command, derrors.Error) {
	var gc entities.GenericCommand
	if err := json.Unmarshal(raw, &gc); err != nil {
		return nil, derrors.NewOperationError(errors.UnmarshalError, err)
	}
	cmd, err := cp.parseCommand(gc, raw)
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

func (cp *CmdParser) parseCommand(generic entities.GenericCommand, raw []byte) (*entities.Command, derrors.Error) {
	switch generic.CommandType {
	case entities.SyncCommandType:
		return cp.parseSyncCommand(generic, raw)
	case entities.AsyncCommandType:
		return cp.parseAsyncCommand(generic, raw)
	default:
		return nil, derrors.NewOperationError(errors.UnsupportedCommandType).WithParams(generic)
	}
}

func (cp *CmdParser) parseSyncCommand(generic entities.GenericCommand, raw []byte) (*entities.Command, derrors.Error) {
	switch generic.CommandName {
	case entities.Exec:
		return sync.NewExecFromJSON(raw)
	case entities.SCP:
		return sync.NewSCPFromJSON(raw)
	case entities.SSH:
		return sync.NewSSHFromJSON(raw)
	case entities.Logger:
		return sync.NewLoggerFromJSON(raw)
	case entities.Sleep:
		return sync.NewSleepFromJSON(raw)
	case entities.Fail:
		return sync.NewFailFromJSON(raw)
	case entities.ParallelCmd:
		return NewParallelFromJSON(raw)
	case entities.GroupCmd:
		return NewGroupFromJSON(raw)
	case entities.TryCmd:
		return NewTryFromJSON(raw)
	case entities.ProcessCheck:
		return sync.NewProcessCheckFromJSON(raw)
	case entities.RKEInstall:
		return rke.NewRKEInstallFromJSON(raw)
	case entities.RKERemove:
		return rke.NewRKERemoveFromJSON(raw)
	default:
		return nil, derrors.NewOperationError(errors.UnsupportedCommand).WithParams(generic)
	}
}

func (cp *CmdParser) parseAsyncCommand(generic entities.GenericCommand, raw []byte) (*entities.Command, derrors.Error) {
	switch generic.CommandName {
	case entities.Fail:
		return async.NewFailFromJSON(raw)
	case entities.Sleep:
		return async.NewSleepFromJSON(raw)
	default:
		return nil, derrors.NewOperationError(errors.UnsupportedCommand).WithParams(generic)
	}
}
