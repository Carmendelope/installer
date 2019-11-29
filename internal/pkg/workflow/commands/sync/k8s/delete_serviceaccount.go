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

// DeleteServiceAccount structure with the attributes required to delete a given service account from a namespace.
type DeleteServiceAccount struct {
	// Kubernetes embedded object
	Kubernetes
	// Namespace with the name of the target namespace
	Namespace string `json:"namespace"`
	// ServiceAccount with the name of the target service account.
	ServiceAccount string `json:"service_account"`
	// FailIfNotExists flag determines if the command fails in case the namespace does not exits.
	FailIfNotExists bool `json:"fail_if_not_exists"`
}

// NewDeleteServiceAccount creates a new DeleteServiceAccount command
func NewDeleteServiceAccount(kubeConfigPath string, namespace string, serviceAccount string) *DeleteServiceAccount {
	return &DeleteServiceAccount{
		Kubernetes: Kubernetes{
			GenericSyncCommand: *entities.NewSyncCommand(entities.DeleteServiceAccount),
			KubeConfigPath:     kubeConfigPath,
		},
		Namespace:      namespace,
		ServiceAccount: serviceAccount,
	}
}

// NewDeleteServiceAccountFromJSON creates a new DeleteServiceAccount command from a raw JSON representation.
func NewDeleteServiceAccountFromJSON(raw []byte) (*entities.Command, derrors.Error) {
	cmc := &DeleteServiceAccount{}
	if err := json.Unmarshal(raw, &cmc); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	cmc.CommandID = entities.GenerateCommandID(cmc.Name())
	var r entities.Command = cmc
	return &r, nil
}

// Run the current command returning the result or an error.
func (dsa *DeleteServiceAccount) Run(workflowID string) (*entities.CommandResult, derrors.Error) {
	connectErr := dsa.Connect()
	if connectErr != nil {
		return nil, connectErr
	}
	exists, err := dsa.ExistsNamespace(dsa.Namespace)
	if err != nil {
		return entities.NewCommandResult(false, "cannot determine if the namespace exists", err), nil
	}
	if !exists && dsa.FailIfNotExists {
		toReturn := derrors.NewNotFoundError("target namespace not found").WithParams(dsa.Namespace)
		return entities.NewCommandResult(false, "target namespace does not exist", toReturn), nil
	}
	log.Debug().Str("namespace", dsa.Namespace).Bool("exists", exists).Msg("namespace check")
	exists, err = dsa.ExistsEntity(dsa.Namespace, "", "v1", "serviceaccounts", dsa.ServiceAccount)
	if err != nil {
		return entities.NewCommandResult(false, "cannot determine if the service account exists", err), nil
	}
	log.Debug().Str("serviceAccount", dsa.ServiceAccount).Bool("exists", exists).Msg("service account check")

	if exists {
		err := dsa.DeleteEntity(dsa.Namespace, "", "v1", "serviceaccounts", dsa.ServiceAccount)
		if err != nil {
			return entities.NewErrCommand("cannot delete service account", err), nil
		}
	}
	return entities.NewSuccessCommand([]byte("Service account deleted")), nil
}

// String returns a string representation
func (dsa *DeleteServiceAccount) String() string {
	return fmt.Sprintf("SYNC DeleteServiceAccount %s:%s", dsa.Namespace, dsa.ServiceAccount)
}

// PrettyPrint returns a simple space indexed string.
func (dsa *DeleteServiceAccount) PrettyPrint(indentation int) string {
	return strings.Repeat(" ", indentation) + dsa.String()
}

// UserString returns a simple string representation of the command for the user.
func (dsa *DeleteServiceAccount) UserString() string {
	return fmt.Sprintf("Deleting service account %s from %s", dsa.ServiceAccount, dsa.Namespace)
}
