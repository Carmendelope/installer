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

package handler

import (
	"github.com/nalej/installer/internal/pkg/errors"
	"sync"

	"github.com/nalej/derrors"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
)

var instanceCommandHandler CommandHandler
var onceCommandHandler sync.Once

// CommandHandler interface definition with all functions exposed related to handling the execution of
// the commands of a workflow.
type CommandHandler interface {
	AddCommand(id string,
		resultCallback func(id string, result *entities.CommandResult, error derrors.Error),
		logCallback func(id string, logEntry string),
	) derrors.Error
	AddLogEntry(id string, logEntry string) derrors.Error
	AttachLogListener(id string, f func(logEntry string))
	FinishCommand(id string, result *entities.CommandResult, error derrors.Error) derrors.Error
}

// GetCommandHandler provides a singleton implementation to retrieve the CommandHandler to be used for all commands
// of a workflow.
func GetCommandHandler() CommandHandler {
	onceCommandHandler.Do(func() {
		instanceCommandHandler = NewCommandHandler()
	})
	return instanceCommandHandler
}

type commandHandler struct {
	sync.Mutex
	resultCallbacks map[string]func(id string, result *entities.CommandResult, error derrors.Error)
	logCallbacks    map[string]func(id string, logEntry string)
	logListeners    map[string]func(logEntry string)
}

// NewCommandHandler creates a new CommandHandler initializing the internal structures.
func NewCommandHandler() CommandHandler {
	return &commandHandler{
		resultCallbacks: make(map[string]func(id string, result *entities.CommandResult, error derrors.Error)),
		logCallbacks:    make(map[string]func(id string, logEntry string)),
		logListeners:    make(map[string]func(logEntry string)),
	}
}

func (h *commandHandler) AddCommand(id string,
	resultCallback func(id string, result *entities.CommandResult, error derrors.Error),
	logCallback func(id string, logEntry string)) derrors.Error {
	h.Lock()
	defer h.Unlock()
	_, exist := h.resultCallbacks[id]
	if exist {
		return derrors.NewAlreadyExistsError(errors.DuplicatedIDCommand).WithParams(id)
	}
	_, exist = h.logCallbacks[id]
	if exist {
		return derrors.NewAlreadyExistsError(errors.DuplicatedIDCommand).WithParams(id)
	}
	h.resultCallbacks[id] = resultCallback
	h.logCallbacks[id] = logCallback
	return nil
}

func (h *commandHandler) AddLogEntry(id string, logEntry string) derrors.Error {
	h.Lock()
	defer h.Unlock()
	callback, exist := h.logCallbacks[id]
	if !exist {
		return derrors.NewNotFoundError(errors.NotExistCommand).WithParams(id)
	}
	go callback(id, logEntry)
	listener, exists := h.logListeners[id]
	if exists {
		go listener(logEntry)
	}
	return nil
}

func (h *commandHandler) AttachLogListener(id string, f func(logEntry string)) {
	h.logListeners[id] = f
}

func (h *commandHandler) FinishCommand(id string, result *entities.CommandResult,
	error derrors.Error) derrors.Error {
	h.Lock()
	defer h.Unlock()
	callback, exist := h.resultCallbacks[id]
	if !exist {
		return derrors.NewNotFoundError(errors.NotExistCommand).WithParams(id)
	}
	_, exist = h.logCallbacks[id]
	if !exist {
		return derrors.NewNotFoundError(errors.NotExistCommand).WithParams(id)
	}
	go callback(id, result, error)
	delete(h.logCallbacks, id)
	delete(h.resultCallbacks, id)
	delete(h.logListeners, id)
	return nil

}
