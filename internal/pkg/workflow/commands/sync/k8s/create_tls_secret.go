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

type CreateTLSSecret struct {
	Kubernetes
	SecretName     string `json:"secret_name"`
	PrivateKeyPath string `json:"private_key_path"`
	CertPath       string `json:"cert_path"`
}

func NewCreateTLSSecret(
	kubeConfigPath string,
	secretName string,
	privateKeyPath string,
	certPath string) *CreateTLSSecret {
	return &CreateTLSSecret{
		Kubernetes: Kubernetes{
			GenericSyncCommand: *entities.NewSyncCommand(entities.CreateTLSSecret),
			KubeConfigPath:     kubeConfigPath,
		},
		SecretName:     secretName,
		PrivateKeyPath: privateKeyPath,
		CertPath:       certPath,
	}
}

func NewCreateTLSSecretFromJSON(raw []byte) (*entities.Command, derrors.Error) {
	f := &CreateTLSSecret{}
	if err := json.Unmarshal(raw, &f); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	f.CommandID = entities.GenerateCommandID(f.Name())
	var r entities.Command = f
	return &r, nil
}

func (cmd *CreateTLSSecret) createKubernetesSecrets() derrors.Error {

	var privateKeyRawContent []byte
	var certRawContent []byte

	if cmd.PrivateKeyPath != "" {
		pkc, err := ioutil.ReadFile(cmd.PrivateKeyPath)
		if err != nil {
			return derrors.AsError(err, "cannot load private key content")
		}
		privateKeyRawContent = pkc
	} else {
		privateKeyRawContent = make([]byte, 0)
	}

	cc, err := ioutil.ReadFile(cmd.CertPath)
	if err != nil {
		return derrors.AsError(err, "cannot load cert content")
	}
	certRawContent = cc

	TLSSecret := &v1.Secret{
		TypeMeta: metaV1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		Type: v1.SecretTypeTLS,
		ObjectMeta: metaV1.ObjectMeta{
			Name:         cmd.SecretName,
			GenerateName: "",
			Namespace:    "nalej",
		},
		Data: map[string][]byte{
			"tls.key": privateKeyRawContent,
			"tls.crt": certRawContent,
		},
	}
	derr := cmd.Create(TLSSecret)
	if derr != nil {
		return derr
	}

	return nil
}

// Run triggers the execution of the command.
func (cmd *CreateTLSSecret) Run(workflowID string) (*entities.CommandResult, derrors.Error) {
	err := cmd.Connect()
	if err != nil {
		log.Info().Str("kubeConfigPath", cmd.KubeConfigPath).Msg("error connecting to cluster")
		return nil, derrors.NewGenericError("error connecting to cluster", err)
	}

	dErr := cmd.createKubernetesSecrets()
	if dErr != nil {
		log.Error().Str("trace", dErr.DebugReport()).Msg("cannot create kubernetes secrets")
		return nil, dErr
	}

	return entities.NewSuccessCommand([]byte("Secret successfully created.")), nil
}

func (cmd *CreateTLSSecret) String() string {
	return fmt.Sprintf("SYNC CreateTLSSecret %s", cmd.SecretName)
}

func (cmd *CreateTLSSecret) PrettyPrint(indentation int) string {
	return strings.Repeat(" ", indentation) + cmd.String()
}

func (cmd *CreateTLSSecret) UserString() string {
	return fmt.Sprintf("creating tls secret")
}
