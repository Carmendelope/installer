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
	"io/ioutil"
	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type CreateOpaqueSecret struct {
	Kubernetes
	SecretName          string `json:"secret_name"`
	SecretKey           string `json:"secret_key"`
	SecretValue         string `json:"secret_value"`
	LoadFromPath        bool   `json:"load_from_path"`
	SecretValueFromPath string `json:"secret_value_from_path"`
}

func NewCreateOpaqueSecret(
	kubeConfigPath string,
	secretName string,
	secretKey string,
	secretValue string,
	loadFromPath bool,
	secretValueFromPath string) *CreateOpaqueSecret {
	return &CreateOpaqueSecret{
		Kubernetes: Kubernetes{
			GenericSyncCommand: *entities.NewSyncCommand(entities.CreateOpaqueSecret),
			KubeConfigPath:     kubeConfigPath,
		},
		SecretName:          secretName,
		SecretKey:           secretKey,
		SecretValue:         secretValue,
		LoadFromPath:        loadFromPath,
		SecretValueFromPath: secretValueFromPath,
	}
}

func NewCreateOpaqueSecretFromJSON(raw []byte) (*entities.Command, derrors.Error) {
	f := &CreateOpaqueSecret{}
	if err := json.Unmarshal(raw, &f); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	f.CommandID = entities.GenerateCommandID(f.Name())
	var r entities.Command = f
	return &r, nil
}

func (cmd *CreateOpaqueSecret) createKubernetesSecrets() derrors.Error {

	var secretRawContent []byte
	if cmd.LoadFromPath {
		c, err := ioutil.ReadFile(cmd.SecretValueFromPath)
		if err != nil {
			return derrors.AsError(err, "cannot load secret content")
		}
		secretRawContent = c
	} else {
		secretRawContent = []byte(cmd.SecretValue)
	}

	OpaqueSecret := &v1.Secret{
		TypeMeta: metaV1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name:         cmd.SecretName,
			GenerateName: "",
			Namespace:    "nalej",
		},
		Data: map[string][]byte{
			cmd.SecretKey: secretRawContent,
		},
		Type: v1.SecretTypeOpaque,
	}
	derr := cmd.Create(OpaqueSecret)
	if derr != nil {
		return derr
	}

	return nil
}

// Run triggers the execution of the command.
func (cmd *CreateOpaqueSecret) Run(workflowID string) (*entities.CommandResult, derrors.Error) {
	err := cmd.Connect()
	if err != nil {
		log.Info().Str("kubeConfigPath", cmd.KubeConfigPath).Msg("error connecting to cluster")
		return nil, derrors.NewGenericError("error connecting to cluster", err)
	}

	dErr := cmd.createKubernetesSecrets()
	if dErr != nil {
		return nil, dErr
	}

	return entities.NewSuccessCommand([]byte("Secret successfully created.")), nil
}

func (cmd *CreateOpaqueSecret) String() string {
	return fmt.Sprintf("SYNC CreateOpaqueSecret %s", cmd.SecretName)
}

func (cmd *CreateOpaqueSecret) PrettyPrint(indentation int) string {
	return strings.Repeat(" ", indentation) + cmd.String()
}

func (cmd *CreateOpaqueSecret) UserString() string {
	return fmt.Sprintf("creating opaque secret")
}
