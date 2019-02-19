/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package k8s

import (
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

// CreateRegistrySecrets creates the secrets related to the docker registries available for internal components. Two
// types of secrets may be created. First, the docker credentials are created so that nalej images can be downloaded.
// Additionally, a secret is generated on the management cluster with the values required to create secrets on the
// application cluster.
type CreateRegistrySecrets struct {
	Kubernetes
	OnManagementCluster  bool   `json:"on_management_cluster"`
	CredentialsName string `json:"credentials_name"`
	Username string `json:"username"`
	Password string `json:"password"`
	URL string `json:"url"`
}

func NewCreateRegistrySecrets(
	kubeConfigPath string,
	onManagementCluster bool,
	credentialsName string, username string, password string, url string) *CreateRegistrySecrets {
	return &CreateRegistrySecrets{
		Kubernetes: Kubernetes{
			GenericSyncCommand: *entities.NewSyncCommand(entities.CreateManagementConfig),
			KubeConfigPath:     kubeConfigPath,
		},
		OnManagementCluster: onManagementCluster,
		CredentialsName:                credentialsName,
		Username:            username,
		Password:            password,
		URL:                 url,
	}
}

func NewCreateRegistrySecretsFromJSON(raw []byte) (*entities.Command, derrors.Error) {
	cmd := &CreateRegistrySecrets{}
	if err := json.Unmarshal(raw, &cmd); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	cmd.CommandID = entities.GenerateCommandID(cmd.Name())
	var r entities.Command = cmd
	return &r, nil
}

// createEnvironmentSecret creates the secret that will be mounted by the installer to be able to trigger
// the install of application clusters.
func (cmd* CreateRegistrySecrets) createEnvironmentSecret() derrors.Error{
	envSecret := &v1.Secret{
		TypeMeta: v12.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: v12.ObjectMeta{
			Name:      fmt.Sprintf("credentials-%s", cmd.CredentialsName),
			Namespace: TargetNamespace,
			Labels:    map[string]string{"cluster": "management"},
		},
		Data: map[string][]byte{
			"credentials_name":[]byte(cmd.CredentialsName),
			"username": []byte(cmd.Username),
			"password": []byte(cmd.Password),
			"url": []byte(cmd.URL),
		},
		Type: v1.SecretTypeOpaque,
	}
	client := cmd.Client.CoreV1().Secrets(envSecret.Namespace)
	created, err := client.Create(envSecret)
	if err != nil {
		return derrors.AsError(err, "cannot create environment registry-credentials secret")
	}
	log.Debug().Interface("created", created).Msg("new secret has been created")
	return nil
}

func (cmd * CreateRegistrySecrets) createDockerSecrets(workflowID string) derrors.Error{
	// Reuse the existing create docker secret commands

	// Create the production secret
	secret := NewCreateDockerSecret(cmd.KubeConfigPath, cmd.CredentialsName,
		cmd.Username, cmd.Password, cmd.URL)
	result, err := secret.Run(workflowID)
	if err != nil{
		return err
	}
	if !result.Success{
		return result.Error
	}

	return nil
}

func (cmd *CreateRegistrySecrets) Run(workflowID string) (*entities.CommandResult, derrors.Error) {
	connectErr := cmd.Connect()
	if connectErr != nil {
		return nil, connectErr
	}

	cErr := cmd.CreateNamespacesIfNotExist("nalej")
	if cErr != nil {
		return entities.NewCommandResult(false, "cannot create namespace", cErr), nil
	}
	if cmd.OnManagementCluster{
		sErr := cmd.createEnvironmentSecret()
		if sErr != nil{
			return entities.NewCommandResult(false, "cannot create environment secret", sErr), nil
		}
	}
	sErr := cmd.createDockerSecrets(workflowID)
	if sErr != nil{
		return entities.NewCommandResult(false, "cannot create docker registry secret", sErr), nil
	}
	// Create Docker secrets
	log.Debug().Msg("management registry secret has been created")
	return entities.NewSuccessCommand([]byte("management registry secret has been created")), nil
}

func (cmd *CreateRegistrySecrets) String() string {
	return fmt.Sprintf("SYNC CreateRegistrySecrets for a %s environment", cmd.CredentialsName)
}

func (cmd *CreateRegistrySecrets) PrettyPrint(indentation int) string {
	return strings.Repeat(" ", indentation) + cmd.String()
}

func (cmd *CreateRegistrySecrets) UserString() string {
	return fmt.Sprintf("Creating managment registry secrets for a %s environment", cmd.CredentialsName)
}
