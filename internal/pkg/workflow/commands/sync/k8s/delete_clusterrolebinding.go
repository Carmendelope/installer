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
	"strings"
)

// DeleteClusterRoleBinding structure with the attributes required to delete a given cluster role binding.
type DeleteClusterRoleBinding struct {
	// Kubernetes embedded object
	Kubernetes
	// RoleBindingName with the name of the target cluster role binding.
	RoleBindingName string `json:"role_binding_name"`
	// FailIfNotExists flag determines if the command fails in case the entity does not exits.
	FailIfNotExists bool `json:"fail_if_not_exists"`
}

// NewDeleteClusterRoleBinding creates a new DeleteClusterRoleBinding command
func NewDeleteClusterRoleBinding(kubeConfigPath string) *DeleteClusterRoleBinding {
	return &DeleteClusterRoleBinding{
		Kubernetes: Kubernetes{
			GenericSyncCommand: *entities.NewSyncCommand(entities.DeleteClusterRoleBinding),
			KubeConfigPath:     kubeConfigPath,
		},
	}
}

// NewDeleteClusterRoleBindingFromJSON creates a new DeleteClusterRoleBinding command from a raw JSON representation.
func NewDeleteClusterRoleBindingFromJSON(raw []byte) (*entities.Command, derrors.Error) {
	cmc := &DeleteClusterRoleBinding{}
	if err := json.Unmarshal(raw, &cmc); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	cmc.CommandID = entities.GenerateCommandID(cmc.Name())
	var r entities.Command = cmc
	return &r, nil
}

// Run the current command returning the result or an error.
func (dcrb *DeleteClusterRoleBinding) Run(workflowID string) (*entities.CommandResult, derrors.Error) {
	connectErr := dcrb.Connect()
	if connectErr != nil {
		return nil, connectErr
	}

	exists, err := dcrb.ExistClusterRoleBinding(dcrb.RoleBindingName)
	if err != nil {
		return entities.NewCommandResult(false, "cannot determine if the cluster role binding exists", err), nil
	}
	if !exists && dcrb.FailIfNotExists {
		toReturn := derrors.NewNotFoundError("cluster role binding not found").WithParams(dcrb.RoleBindingName)
		return entities.NewCommandResult(false, "cluster role binding does not exist", toReturn), nil
	}

	if exists {
		err := dcrb.DeleteEntity("", "rbac.authorization.k8s.io", "v1", "clusterrolebindings", dcrb.RoleBindingName)
		if err != nil {
			return entities.NewErrCommand("cannot delete cluster role binding", err), nil
		}
	}
	return entities.NewSuccessCommand([]byte("Cluster role binding deleted")), nil
}

// String returns a string representation
func (dcrb *DeleteClusterRoleBinding) String() string {
	return fmt.Sprintf("SYNC DeleteClusterRoleBinding %s", dcrb.RoleBindingName)
}

// PrettyPrint returns a simple space indexed string.
func (dcrb *DeleteClusterRoleBinding) PrettyPrint(indentation int) string {
	return strings.Repeat(" ", indentation) + dcrb.String()
}

// UserString returns a simple string representation of the command for the user.
func (dcrb *DeleteClusterRoleBinding) UserString() string {
	return fmt.Sprintf("Deleting cluster role binding %s", dcrb.RoleBindingName)
}
