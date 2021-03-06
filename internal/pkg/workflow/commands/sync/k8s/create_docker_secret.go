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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/installer/internal/pkg/errors"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
	"k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

// CreateDockerSecret is a generic command to create docker secrets in kubernetes.
type CreateDockerSecret struct {
	Kubernetes
	SecretName string `json:"secret_name"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	URL        string `json:"url"`
}

func NewCreateDockerSecret(
	kubeConfigPath string,
	name string, username string, password string, url string) *CreateDockerSecret {
	return &CreateDockerSecret{
		Kubernetes: Kubernetes{
			GenericSyncCommand: *entities.NewSyncCommand(entities.CreateDockerSecret),
			KubeConfigPath:     kubeConfigPath,
		},
		SecretName: name,
		Username:   username,
		Password:   password,
		URL:        url,
	}
}

func NewCreateDockerSecretFromJSON(raw []byte) (*entities.Command, derrors.Error) {
	cmd := &CreateDockerSecret{}
	if err := json.Unmarshal(raw, &cmd); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	cmd.CommandID = entities.GenerateCommandID(cmd.Name())
	var r entities.Command = cmd
	return &r, nil
}

func (cmd *CreateDockerSecret) getAuth() string {
	toEncode := fmt.Sprintf("%s:%s", cmd.Username, cmd.Password)
	encoded := base64.StdEncoding.EncodeToString([]byte(toEncode))
	return encoded
}

func (cmd *CreateDockerSecret) getDockerConfigJSON() string {
	template := "{\"auths\":{\"%s\":{\"username\":\"%s\",\"password\":\"%s\",\"email\":\"devops@daisho.group\",\"auth\":\"%s\"}}}"
	toEncode := fmt.Sprintf(template, cmd.URL, cmd.Username, cmd.Password, cmd.getAuth())
	return toEncode
}

func (cmd *CreateDockerSecret) Run(workflowID string) (*entities.CommandResult, derrors.Error) {
	connectErr := cmd.Connect()
	if connectErr != nil {
		return nil, connectErr
	}
	cErr := cmd.CreateNamespaceIfNotExists("nalej")
	if cErr != nil {
		return entities.NewCommandResult(false, "cannot create namespace", cErr), nil
	}

	secret := &v1.Secret{
		TypeMeta: v12.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: v12.ObjectMeta{
			Name:         cmd.SecretName,
			GenerateName: "",
			Namespace:    "nalej",
		},
		Data: map[string][]byte{
			".dockerconfigjson": []byte(cmd.getDockerConfigJSON()),
		},
		Type: v1.SecretTypeDockerConfigJson,
	}

	derr := cmd.Create(secret)
	if derr != nil {
		return entities.NewCommandResult(
			false, "cannot create docker registry credentials", derrors.AsError(derr, "cannot create registry credentials")), nil
	}
	return entities.NewSuccessCommand([]byte("docker registry credentials have been created")), nil
}

func (cmd *CreateDockerSecret) String() string {
	return fmt.Sprintf("SYNC CreateDockerSecret for %s", cmd.URL)
}

func (cmd *CreateDockerSecret) PrettyPrint(indentation int) string {
	return strings.Repeat(" ", indentation) + cmd.String()
}

func (cmd *CreateDockerSecret) UserString() string {
	return fmt.Sprintf("Creating docker secrets for %s", cmd.URL)
}
