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

package k8s

import (
	"encoding/json"
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/installer/internal/pkg/errors"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
	"github.com/rs/zerolog/log"
	"strings"
)

// DeleteService structure with the attributes required to delete a given service from a namespace.
type DeleteService struct {
	// Kubernetes embedded object
	Kubernetes
	// Namespace with the name of the target namespace
	Namespace string `json:"namespace"`
	// ServiceName with the name of the target service.
	ServiceName string `json:"service_name"`
	// FailIfNotExists flag determines if the command fails in case the namespace does not exits.
	FailIfNotExists bool `json:"fail_if_not_exists"`
}

// NewDeleteService creates a new DeleteService command
func NewDeleteService(kubeConfigPath string, namespace string, serviceName string) *DeleteService {
	return &DeleteService{
		Kubernetes: Kubernetes{
			GenericSyncCommand: *entities.NewSyncCommand(entities.DeleteService),
			KubeConfigPath:     kubeConfigPath,
		},
		Namespace:   namespace,
		ServiceName: serviceName,
	}
}

// NewDeleteServiceFromJSON creates a new DeleteService command from a raw JSON representation.
func NewDeleteServiceFromJSON(raw []byte) (*entities.Command, derrors.Error) {
	cmc := &DeleteService{}
	if err := json.Unmarshal(raw, &cmc); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	cmc.CommandID = entities.GenerateCommandID(cmc.Name())
	var r entities.Command = cmc
	return &r, nil
}

// Run the current command returning the result or an error.
func (ds *DeleteService) Run(workflowID string) (*entities.CommandResult, derrors.Error) {
	connectErr := ds.Connect()
	if connectErr != nil {
		return nil, connectErr
	}
	exists, err := ds.ExistsNamespace(ds.Namespace)
	if err != nil {
		return entities.NewCommandResult(false, "cannot determine if the namespace exists", err), nil
	}
	if !exists && ds.FailIfNotExists {
		toReturn := derrors.NewNotFoundError("target namespace not found").WithParams(ds.Namespace)
		return entities.NewCommandResult(false, "target namespace does not exist", toReturn), nil
	}
	log.Debug().Str("namespace", ds.Namespace).Bool("exists", exists).Msg("namespace check")
	exists, err = ds.ExistsEntity(ds.Namespace, "", "v1", "services", ds.ServiceName)
	if err != nil {
		return entities.NewCommandResult(false, "cannot determine if the service exists", err), nil
	}
	log.Debug().Str("serviceName", ds.ServiceName).Bool("exists", exists).Msg("service check")

	if exists {
		err := ds.DeleteEntity(ds.Namespace, "", "v1", "services", ds.ServiceName)
		if err != nil {
			return entities.NewErrCommand("cannot delete service", err), nil
		}
	}
	return entities.NewSuccessCommand([]byte("Service deleted")), nil
}

// String returns a string representation
func (ds *DeleteService) String() string {
	return fmt.Sprintf("SYNC DeleteService %s:%s", ds.Namespace, ds.ServiceName)
}

// PrettyPrint returns a simple space indexed string.
func (ds *DeleteService) PrettyPrint(indentation int) string {
	return strings.Repeat(" ", indentation) + ds.String()
}

// UserString returns a simple string representation of the command for the user.
func (ds *DeleteService) UserString() string {
	return fmt.Sprintf("Deleting service %s from %s", ds.ServiceName, ds.Namespace)
}
