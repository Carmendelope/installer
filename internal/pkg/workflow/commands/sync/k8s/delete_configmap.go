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

// DeleteConfigMap structure with the attributes required to delete a given configmap from a namespace.
type DeleteConfigMap struct {
	// Kubernetes embedded object
	Kubernetes
	// Namespace with the name of the target namespace
	Namespace string `json:"namespace"`
	// ConfigMapName with the name of the target role.
	ConfigMapName string `json:"config_map_name"`
	// FailIfNotExists flag determines if the command fails in case the namespace does not exits.
	FailIfNotExists bool `json:"fail_if_not_exists"`
}

// NewDeleteConfigMap creates a new DeleteConfigMap command
func NewDeleteConfigMap(kubeConfigPath string, namespace string, configMapName string) *DeleteConfigMap {
	return &DeleteConfigMap{
		Kubernetes: Kubernetes{
			GenericSyncCommand: *entities.NewSyncCommand(entities.DeleteConfigMap),
			KubeConfigPath:     kubeConfigPath,
		},
		Namespace:     namespace,
		ConfigMapName: configMapName,
	}
}

// NewDeleteConfigMapFromJSON creates a new DeleteConfigMap command from a raw JSON representation.
func NewDeleteConfigMapFromJSON(raw []byte) (*entities.Command, derrors.Error) {
	cmc := &DeleteConfigMap{}
	if err := json.Unmarshal(raw, &cmc); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	cmc.CommandID = entities.GenerateCommandID(cmc.Name())
	var r entities.Command = cmc
	return &r, nil
}

// Run the current command returning the result or an error.
func (dcm *DeleteConfigMap) Run(workflowID string) (*entities.CommandResult, derrors.Error) {
	connectErr := dcm.Connect()
	if connectErr != nil {
		return nil, connectErr
	}
	exists, err := dcm.ExistsNamespace(dcm.Namespace)
	if err != nil {
		return entities.NewCommandResult(false, "cannot determine if the namespace exists", err), nil
	}
	if !exists && dcm.FailIfNotExists {
		toReturn := derrors.NewNotFoundError("target namespace not found").WithParams(dcm.Namespace)
		return entities.NewCommandResult(false, "target namespace does not exist", toReturn), nil
	}
	log.Debug().Str("namespace", dcm.Namespace).Bool("exists", exists).Msg("namespace check")
	exists, err = dcm.ExistsEntity(dcm.Namespace, "", "v1", "configmaps", dcm.ConfigMapName)
	if err != nil {
		return entities.NewCommandResult(false, "cannot determine if the config map exists", err), nil
	}
	log.Debug().Str("configMapName", dcm.ConfigMapName).Bool("exists", exists).Msg("config map check")

	if exists {
		err := dcm.DeleteEntity(dcm.Namespace, "", "v1", "configmaps", dcm.ConfigMapName)
		if err != nil {
			return entities.NewErrCommand("cannot delete config map", err), nil
		}
	}
	return entities.NewSuccessCommand([]byte("Config map deleted")), nil
}

// String returns a string representation
func (dcm *DeleteConfigMap) String() string {
	return fmt.Sprintf("SYNC DeleteConfigMap %s:%s", dcm.Namespace, dcm.ConfigMapName)
}

// PrettyPrint returns a simple space indexed string.
func (dcm *DeleteConfigMap) PrettyPrint(indentation int) string {
	return strings.Repeat(" ", indentation) + dcm.String()
}

// UserString returns a simple string representation of the command for the user.
func (dcm *DeleteConfigMap) UserString() string {
	return fmt.Sprintf("Deleting config map %s from %s", dcm.ConfigMapName, dcm.Namespace)
}
