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

// DeleteRoleBinding structure with the attributes required to delete a given role binding from a namespace.
type DeleteRoleBinding struct {
	// Kubernetes embedded object
	Kubernetes
	// Namespace with the name of the target namespace
	Namespace string `json:"namespace"`
	// RoleName with the name of the target role.
	RoleName string `json:"role_name"`
	// FailIfNotExists flag determines if the command fails in case the namespace does not exits.
	FailIfNotExists bool `json:"fail_if_not_exists"`
}

// NewDeleteRoleBinding creates a new DeleteRoleBinding command
func NewDeleteRoleBinding(kubeConfigPath string, namespace string, roleName string) *DeleteRoleBinding {
	return &DeleteRoleBinding{
		Kubernetes: Kubernetes{
			GenericSyncCommand: *entities.NewSyncCommand(entities.DeleteRoleBinding),
			KubeConfigPath:     kubeConfigPath,
		},
		Namespace: namespace,
		RoleName:  roleName,
	}
}

// NewDeleteRoleBindingFromJSON creates a new DeleteRoleBinding command from a raw JSON representation.
func NewDeleteRoleBindingFromJSON(raw []byte) (*entities.Command, derrors.Error) {
	cmc := &DeleteRoleBinding{}
	if err := json.Unmarshal(raw, &cmc); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	cmc.CommandID = entities.GenerateCommandID(cmc.Name())
	var r entities.Command = cmc
	return &r, nil
}

// Run the current command returning the result or an error.
func (drb *DeleteRoleBinding) Run(workflowID string) (*entities.CommandResult, derrors.Error) {
	connectErr := drb.Connect()
	if connectErr != nil {
		return nil, connectErr
	}
	exists, err := drb.ExistsNamespace(drb.Namespace)
	if err != nil {
		return entities.NewCommandResult(false, "cannot determine if the namespace exists", err), nil
	}
	if !exists && drb.FailIfNotExists {
		toReturn := derrors.NewNotFoundError("target namespace not found").WithParams(drb.Namespace)
		return entities.NewCommandResult(false, "target namespace does not exist", toReturn), nil
	}
	log.Debug().Str("namespace", drb.Namespace).Bool("exists", exists).Msg("namespace check")
	exists, err = drb.ExistsEntity(drb.Namespace, "rbac.authorization.k8s.io", "v1", "rolebindings", drb.RoleName)
	if err != nil {
		return entities.NewCommandResult(false, "cannot determine if the role binding exists", err), nil
	}
	log.Debug().Str("roleName", drb.RoleName).Bool("exists", exists).Msg("role binding check")

	if exists {
		err := drb.DeleteEntity(drb.Namespace, "rbac.authorization.k8s.io", "v1", "rolebindings", drb.RoleName)
		if err != nil {
			return entities.NewErrCommand("cannot delete role binding", err), nil
		}
	}
	return entities.NewSuccessCommand([]byte("Role binding deleted")), nil
}

// String returns a string representation
func (drb *DeleteRoleBinding) String() string {
	return fmt.Sprintf("SYNC DeleteRoleBinding %s:%s", drb.Namespace, drb.RoleName)
}

// PrettyPrint returns a simple space indexed string.
func (drb *DeleteRoleBinding) PrettyPrint(indentation int) string {
	return strings.Repeat(" ", indentation) + drb.String()
}

// UserString returns a simple string representation of the command for the user.
func (drb *DeleteRoleBinding) UserString() string {
	return fmt.Sprintf("Deleting role binding %s from %s", drb.RoleName, drb.Namespace)
}
