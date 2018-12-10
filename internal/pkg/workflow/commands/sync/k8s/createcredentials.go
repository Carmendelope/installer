/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package k8s

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/installer/internal/pkg/errors"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
	"github.com/rs/zerolog/log"
	"k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

const DockerRegistryServer = "https://nalejregistry.azurecr.io"

type CreateCredentials struct {
	Kubernetes
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewCreateCredentials(kubeConfigPath string, username string, password string) *CreateCredentials {
	return &CreateCredentials{
		Kubernetes: Kubernetes{
			GenericSyncCommand: *entities.NewSyncCommand(entities.UpdateCoreDNS),
			KubeConfigPath:     kubeConfigPath,
		},
		Username: username,
		Password: password,
	}
}

func NewCreateCredentialsJSON(raw []byte) (*entities.Command, derrors.Error) {
	ccc := &CreateCredentials{}
	if err := json.Unmarshal(raw, &ccc); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	ccc.CommandID = entities.GenerateCommandID(ccc.Name())
	var r entities.Command = ccc
	return &r, nil
}

func (cc *CreateCredentials) getAuth() string {
	toEncode := fmt.Sprintf("%s:%s", cc.Username, cc.Password)
	encoded := base64.StdEncoding.EncodeToString([]byte(toEncode))
	return encoded
}

func (cc *CreateCredentials) getDockerConfigJSON() string {
	template := "{\"auths\":{\"%s\":{\"username\":\"%s\",\"password\":\"%s\",\"email\":\"devops@daisho.group\",\"auth\":\"%s\"}}}"
	toEncode := fmt.Sprintf(template, DockerRegistryServer, cc.Username, cc.Password, cc.getAuth())
	return toEncode
}

func (cc *CreateCredentials) Run(workflowID string) (*entities.CommandResult, derrors.Error) {
	connectErr := cc.Connect()
	if connectErr != nil {
		return nil, connectErr
	}
	cErr := cc.CreateNamespacesIfNotExist("nalej")
	if cErr != nil {
		return entities.NewCommandResult(false, "cannot create namespace", cErr), nil
	}

	secret := &v1.Secret{
		TypeMeta: v12.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: v12.ObjectMeta{
			Name:         "nalej-registry",
			GenerateName: "",
			Namespace:    "nalej",
		},
		Data: map[string][]byte{
			".dockerconfigjson": []byte(cc.getDockerConfigJSON()),
		},
		Type: v1.SecretTypeDockerConfigJson,
	}

	client := cc.Client.CoreV1().Secrets(secret.Namespace)
	created, err := client.Create(secret)
	if err != nil {
		return entities.NewCommandResult(
			false, "cannot create registry credentials", derrors.AsError(err, "cannot create registry credentials")), nil
	}
	log.Debug().Interface("created", created).Msg("new secret has been created")
	return entities.NewSuccessCommand([]byte("registry credentials have been created")), nil
}

func (cc *CreateCredentials) String() string {
	return fmt.Sprintf("SYNC CreateCredentials to %s", DockerRegistryServer)
}

func (cc *CreateCredentials) PrettyPrint(indentation int) string {
	return strings.Repeat(" ", indentation) + cc.String()
}

func (cc *CreateCredentials) UserString() string {
	return fmt.Sprintf("Creating credentials to access %s", DockerRegistryServer)
}
