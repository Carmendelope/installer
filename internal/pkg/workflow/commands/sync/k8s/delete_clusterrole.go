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

// DeleteClusterRole structure with the attributes required to delete a given cluster role.
type DeleteClusterRole struct {
	// Kubernetes embedded object
	Kubernetes
	// RoleName with the name of the target cluster role
	RoleName string `json:"role_name"`
	// FailIfNotExists flag determines if the command fails in case the namespace does not exits.
	FailIfNotExists bool `json:"fail_if_not_exists"`
}

// NewDeleteClusterRole creates a new DeleteClusterRole command
func NewDeleteClusterRole(kubeConfigPath string, roleName string) *DeleteClusterRole {
	return &DeleteClusterRole{
		Kubernetes: Kubernetes{
			GenericSyncCommand: *entities.NewSyncCommand(entities.DeleteClusterRole),
			KubeConfigPath:     kubeConfigPath,
		},
		RoleName: roleName,
	}
}

// NewDeleteClusterRoleFromJSON creates a new DeleteClusterRole command from a raw JSON representation.
func NewDeleteClusterRoleFromJSON(raw []byte) (*entities.Command, derrors.Error) {
	cmc := &DeleteClusterRole{}
	if err := json.Unmarshal(raw, &cmc); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	cmc.CommandID = entities.GenerateCommandID(cmc.Name())
	var r entities.Command = cmc
	return &r, nil
}

// Run the current command returning the result or an error.
func (dcr *DeleteClusterRole) Run(workflowID string) (*entities.CommandResult, derrors.Error) {
	connectErr := dcr.Connect()
	if connectErr != nil {
		return nil, connectErr
	}

	exists, err := dcr.ExistsEntity("", "rbac.authorization.k8s.io", "v1", "clusterroles", dcr.RoleName)
	if err != nil {
		return entities.NewCommandResult(false, "cannot determine if the cluster role exists", err), nil
	}
	log.Debug().Str("roleName", dcr.RoleName).Bool("exists", exists).Msg("cluster role check")

	if exists {
		err := dcr.DeleteEntity("", "rbac.authorization.k8s.io", "v1", "clusterroles", dcr.RoleName)
		if err != nil {
			return entities.NewErrCommand("cannot delete cluster role", err), nil
		}
	}
	return entities.NewSuccessCommand([]byte("Cluster role deleted")), nil
}

// String returns a string representation
func (dcr *DeleteClusterRole) String() string {
	return fmt.Sprintf("SYNC DeleteClusterRole %s", dcr.RoleName)
}

// PrettyPrint returns a simple space indexed string.
func (dcr *DeleteClusterRole) PrettyPrint(indentation int) string {
	return strings.Repeat(" ", indentation) + dcr.String()
}

// UserString returns a simple string representation of the command for the user.
func (dcr *DeleteClusterRole) UserString() string {
	return fmt.Sprintf("Deleting cluster role %s", dcr.RoleName)
}
