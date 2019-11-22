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

// DeleteDeployment structure with the attributes required to delete a given deployment from a namespace.
type DeleteDeployment struct {
	// Kubernetes embedded object
	Kubernetes
	// Namespace with the name of the target namespace
	Namespace string `json:"namespace"`
	// DeploymentName with the name of the target deployment.
	DeploymentName string `json:"deployment_name"`
	// FailIfNotExists flag determines if the command fails in case the namespace does not exits.
	FailIfNotExists bool `json:"fail_if_not_exists"`
}

// NewDeleteDeployment creates a new DeleteDeployment command
func NewDeleteDeployment(kubeConfigPath string, namespace string, deploymentName string) *DeleteDeployment {
	return &DeleteDeployment{
		Kubernetes: Kubernetes{
			GenericSyncCommand: *entities.NewSyncCommand(entities.DeleteDeployment),
			KubeConfigPath:     kubeConfigPath,
		},
		Namespace:      namespace,
		DeploymentName: deploymentName,
	}
}

// NewDeleteDeploymentFromJSON creates a new DeleteDeployment command from a raw JSON representation.
func NewDeleteDeploymentFromJSON(raw []byte) (*entities.Command, derrors.Error) {
	cmc := &DeleteDeployment{}
	if err := json.Unmarshal(raw, &cmc); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	cmc.CommandID = entities.GenerateCommandID(cmc.Name())
	var r entities.Command = cmc
	return &r, nil
}

// Run the current command returning the result or an error.
func (dd *DeleteDeployment) Run(workflowID string) (*entities.CommandResult, derrors.Error) {
	connectErr := dd.Connect()
	if connectErr != nil {
		return nil, connectErr
	}
	exists, err := dd.ExistsNamespace(dd.Namespace)
	if err != nil {
		return entities.NewCommandResult(false, "cannot determine if the namespace exists", err), nil
	}
	if !exists && dd.FailIfNotExists {
		toReturn := derrors.NewNotFoundError("target namespace not found").WithParams(dd.Namespace)
		return entities.NewCommandResult(false, "target namespace does not exist", toReturn), nil
	}
	log.Debug().Str("namespace", dd.Namespace).Bool("exists", exists).Msg("namespace check")
	exists, err = dd.ExistsEntity(dd.Namespace, "apps", "v1", "deployments", dd.DeploymentName)
	if err != nil {
		return entities.NewCommandResult(false, "cannot determine if the deployment exists", err), nil
	}
	log.Debug().Str("deploymentName", dd.DeploymentName).Bool("exists", exists).Msg("deployment check")

	if exists {
		err := dd.DeleteEntity(dd.Namespace, "apps", "v1", "deployments", dd.DeploymentName)
		if err != nil {
			return entities.NewErrCommand("cannot delete deployment", err), nil
		}
	}
	return entities.NewSuccessCommand([]byte("Deployment deleted")), nil
}

// String returns a string representation
func (dd *DeleteDeployment) String() string {
	return fmt.Sprintf("SYNC DeleteDeployment %s:%s", dd.Namespace, dd.DeploymentName)
}

// PrettyPrint returns a simple space indexed string.
func (dd *DeleteDeployment) PrettyPrint(indentation int) string {
	return strings.Repeat(" ", indentation) + dd.String()
}

// UserString returns a simple string representation of the command for the user.
func (dd *DeleteDeployment) UserString() string {
	return fmt.Sprintf("Deleting deployment %s from %s", dd.DeploymentName, dd.Namespace)
}
