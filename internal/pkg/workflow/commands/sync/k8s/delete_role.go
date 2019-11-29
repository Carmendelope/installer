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

// DeleteRole structure with the attributes required to delete a given role from a namespace.
type DeleteRole struct {
	// Kubernetes embedded object
	Kubernetes
	// Namespace with the name of the target namespace
	Namespace string `json:"namespace"`
	// RoleName with the name of the target role.
	RoleName string `json:"role_name"`
	// FailIfNotExists flag determines if the command fails in case the namespace does not exits.
	FailIfNotExists bool `json:"fail_if_not_exists"`
}

// NewDeleteRole creates a new DeleteRole command
func NewDeleteRole(kubeConfigPath string, namespace string, roleName string) *DeleteRole {
	return &DeleteRole{
		Kubernetes: Kubernetes{
			GenericSyncCommand: *entities.NewSyncCommand(entities.DeleteRole),
			KubeConfigPath:     kubeConfigPath,
		},
		Namespace: namespace,
		RoleName:  roleName,
	}
}

// NewDeleteRoleFromJSON creates a new DeleteServiceAccount command from a raw JSON representation.
func NewDeleteRoleFromJSON(raw []byte) (*entities.Command, derrors.Error) {
	cmc := &DeleteRole{}
	if err := json.Unmarshal(raw, &cmc); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	cmc.CommandID = entities.GenerateCommandID(cmc.Name())
	var r entities.Command = cmc
	return &r, nil
}

// Run the current command returning the result or an error.
func (dr *DeleteRole) Run(workflowID string) (*entities.CommandResult, derrors.Error) {
	connectErr := dr.Connect()
	if connectErr != nil {
		return nil, connectErr
	}
	exists, err := dr.ExistsNamespace(dr.Namespace)
	if err != nil {
		return entities.NewCommandResult(false, "cannot determine if the namespace exists", err), nil
	}
	if !exists && dr.FailIfNotExists {
		toReturn := derrors.NewNotFoundError("target namespace not found").WithParams(dr.Namespace)
		return entities.NewCommandResult(false, "target namespace does not exist", toReturn), nil
	}
	log.Debug().Str("namespace", dr.Namespace).Bool("exists", exists).Msg("namespace check")
	exists, err = dr.ExistsEntity(dr.Namespace, "rbac.authorization.k8s.io", "v1", "roles", dr.RoleName)
	if err != nil {
		return entities.NewCommandResult(false, "cannot determine if the role exists", err), nil
	}
	log.Debug().Str("roleName", dr.RoleName).Bool("exists", exists).Msg("role check")

	if exists {
		err := dr.DeleteEntity(dr.Namespace, "rbac.authorization.k8s.io", "v1", "roles", dr.RoleName)
		if err != nil {
			return entities.NewErrCommand("cannot delete role", err), nil
		}
	}
	return entities.NewSuccessCommand([]byte("Role deleted")), nil
}

// String returns a string representation
func (dr *DeleteRole) String() string {
	return fmt.Sprintf("SYNC DeleteRole %s:%s", dr.Namespace, dr.RoleName)
}

// PrettyPrint returns a simple space indexed string.
func (dr *DeleteRole) PrettyPrint(indentation int) string {
	return strings.Repeat(" ", indentation) + dr.String()
}

// UserString returns a simple string representation of the command for the user.
func (dr *DeleteRole) UserString() string {
	return fmt.Sprintf("Deleting role %s from %s", dr.RoleName, dr.Namespace)
}
