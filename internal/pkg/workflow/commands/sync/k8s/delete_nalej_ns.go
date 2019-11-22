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

// NalejNamespace constant with the name of the nalej namespace
const NalejNamespace = "nalej"

// ExcludedSecrets contains the name of the secrets that will not be deleted.
var ExcludedSecrets = []string{"tls-client-certificate"}

// DeleteNalejNamespace structure with the attributes required to delete the contents of the Nalej namespace.
type DeleteNalejNamespace struct {
	// Kubernetes embedded object
	Kubernetes
	// FailIfNotExists flag determines if the command fails in case the namespace does not exits.
	FailIfNotExists bool `json:"fail_if_not_exists"`
}

// NewDeleteNalejNamespace creates a new DeleteNalejNamespace command
func NewDeleteNalejNamespace(kubeConfigPath string) *DeleteNalejNamespace {
	return &DeleteNalejNamespace{
		Kubernetes: Kubernetes{
			GenericSyncCommand: *entities.NewSyncCommand(entities.DeleteNalejNamespace),
			KubeConfigPath:     kubeConfigPath,
		},
	}
}

// NewDeleteNalejNamespaceFromJSON creates a new DeleteNalejNamespace command from a raw JSON representation.
func NewDeleteNalejNamespaceFromJSON(raw []byte) (*entities.Command, derrors.Error) {
	cmc := &DeleteNalejNamespace{}
	if err := json.Unmarshal(raw, &cmc); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	cmc.CommandID = entities.GenerateCommandID(cmc.Name())
	var r entities.Command = cmc
	return &r, nil
}

// Run the current command returning the result or an error.
func (dnn *DeleteNalejNamespace) Run(workflowID string) (*entities.CommandResult, derrors.Error) {
	connectErr := dnn.Connect()
	if connectErr != nil {
		return nil, connectErr
	}
	exists, err := dnn.ExistsNamespace(NalejNamespace)
	if err != nil {
		return entities.NewCommandResult(false, "cannot determine if the namespace exists", err), nil
	}
	if !exists && dnn.FailIfNotExists {
		toReturn := derrors.NewNotFoundError("target namespace not found").WithParams(NalejNamespace)
		return entities.NewCommandResult(false, "target namespace does not exist", toReturn), nil
	}
	if exists {
		// Delete deployments
		if err = dnn.DeleteAllEntities(NalejNamespace, "apps", "v1", "deployments"); err != nil {
			return entities.NewErrCommand("cannot delete Nalej deployments", err), nil
		}
		// Delete services
		if err = dnn.DeleteAllEntities(NalejNamespace, "", "v1", "services"); err != nil {
			return entities.NewErrCommand("cannot delete Nalej services", err), nil
		}
		// Delete configmaps
		if err = dnn.DeleteAllEntities(NalejNamespace, "", "v1", "configmaps"); err != nil {
			return entities.NewErrCommand("cannot delete Nalej configmaps", err), nil
		}
		// Delete secrets
		if err = dnn.DeleteAllEntities(NalejNamespace, "", "v1", "serviceaccounts", "default"); err != nil {
			return entities.NewErrCommand("cannot delete Nalej service accounts", err), nil
		}
		// Delete secrets
		if err = dnn.DeleteAllEntities(NalejNamespace, "", "v1", "secrets", ExcludedSecrets...); err != nil {
			return entities.NewErrCommand("cannot delete Nalej secrets", err), nil
		}
		// Delete daemon set
		if err = dnn.DeleteAllEntities(NalejNamespace, "apps", "v1", "daemonsets"); err != nil {
			return entities.NewErrCommand("cannot delete Nalej daemon sets", err), nil
		}
		// Stateful set
		if err = dnn.DeleteAllEntities(NalejNamespace, "apps", "v1", "statefulsets"); err != nil {
			return entities.NewErrCommand("cannot delete Nalej stateful sets", err), nil
		}
	}
	return entities.NewSuccessCommand([]byte("Nalej namespace contents deleted")), nil
}

// String returns a string representation
func (dnn *DeleteNalejNamespace) String() string {
	return fmt.Sprintf("SYNC DeleteNalejNamespace")
}

// PrettyPrint returns a simple space indexed string.
func (dnn *DeleteNalejNamespace) PrettyPrint(indentation int) string {
	return strings.Repeat(" ", indentation) + dnn.String()
}

// UserString returns a simple string representation of the command for the user.
func (dnn *DeleteNalejNamespace) UserString() string {
	return fmt.Sprintf("Deleting contents of Nalej namespace")
}
