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

// DeleteNamespace structure with the attributes required to delete a namespace and all its content from Kubernetes.
type DeleteNamespace struct {
	// Kubernetes embedded object
	Kubernetes
	// Namespace with the name of the target namespace
	Namespace string `json:"namespace"`
	// FailIfNotExists flag determines if the command fails in case the namespace does not exits.
	FailIfNotExists bool `json:"fail_if_not_exists"`
}

// NewDeleteNamespace creates a new DeleteNamespace command
func NewDeleteNamespace(kubeConfigPath string, namespace string) *DeleteNamespace {
	return &DeleteNamespace{
		Kubernetes: Kubernetes{
			GenericSyncCommand: *entities.NewSyncCommand(entities.DeleteNamespace),
			KubeConfigPath:     kubeConfigPath,
		},
		Namespace: namespace,
	}
}

// NewDeleteNamespaceFromJSON creates a new DeleteServiceAccount command from a raw JSON representation.
func NewDeleteNamespaceFromJSON(raw []byte) (*entities.Command, derrors.Error) {
	cmc := &DeleteNamespace{}
	if err := json.Unmarshal(raw, &cmc); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	cmc.CommandID = entities.GenerateCommandID(cmc.Name())
	var r entities.Command = cmc
	return &r, nil
}

// Run the current command returning the result or an error.
func (dn *DeleteNamespace) Run(workflowID string) (*entities.CommandResult, derrors.Error) {
	connectErr := dn.Connect()
	if connectErr != nil {
		return nil, connectErr
	}
	exists, err := dn.ExistsNamespace(dn.Namespace)
	if err != nil {
		return entities.NewCommandResult(false, "cannot determine if the namespace exists", err), nil
	}
	if !exists && dn.FailIfNotExists {
		toReturn := derrors.NewNotFoundError("target namespace not found").WithParams(dn.Namespace)
		return entities.NewCommandResult(false, "target namespace does not exist", toReturn), nil
	}
	if exists {
		err := dn.DeleteEntity("", "", "v1", "namespaces", dn.Namespace)
		if err != nil {
			return entities.NewErrCommand("cannot delete namespace", err), nil
		}
	}
	return entities.NewSuccessCommand([]byte("Namespace deleted")), nil
}

// String returns a string representation
func (dn *DeleteNamespace) String() string {
	return fmt.Sprintf("SYNC DeleteNamespace %s", dn.Namespace)
}

// PrettyPrint returns a simple space indexed string.
func (dn *DeleteNamespace) PrettyPrint(indentation int) string {
	return strings.Repeat(" ", indentation) + dn.String()
}

// UserString returns a simple string representation of the command for the user.
func (dn *DeleteNamespace) UserString() string {
	return fmt.Sprintf("Deleting namespace %s", dn.Namespace)
}
