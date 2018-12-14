/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package k8s

import (
	"encoding/json"
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/installer/internal/pkg/errors"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
	"github.com/rs/zerolog/log"
	"github.com/satori/go.uuid"
	"k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

const TargetNamespace = "nalej"

type CreateManagementConfig struct {
	Kubernetes
	PublicHost     string `json:"public_host"`
	PublicPort     string `json:"public_port"`
	DNSHost     string `json:"dns_host"`
	DNSPort     string `json:"dns_port"`
	DockerUsername string `json:"docker_username"`
	DockerPassword string `json:"docker_password"`
	PlatformType string `json:"platform_type"`
}

func NewCreateManagementConfig(
	kubeConfigPath string,
	publicHost string, publicPort string,
	dockerUsername string, dockerPassword string,
	platformType string) *CreateManagementConfig {
	return &CreateManagementConfig{
		Kubernetes: Kubernetes{
			GenericSyncCommand: *entities.NewSyncCommand(entities.CreateManagementConfig),
			KubeConfigPath:     kubeConfigPath,
		},
		PublicHost:     publicHost,
		PublicPort:     publicPort,
		DockerUsername: dockerUsername,
		DockerPassword: dockerPassword,
		PlatformType: platformType,
	}
}

func NewCreateManagementConfigFromJSON(raw []byte) (*entities.Command, derrors.Error) {
	cmc := &CreateManagementConfig{}
	if err := json.Unmarshal(raw, &cmc); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	cmc.CommandID = entities.GenerateCommandID(cmc.Name())
	var r entities.Command = cmc
	return &r, nil
}

func (cmc *CreateManagementConfig) createConfigMap() derrors.Error {
	config := &v1.ConfigMap{
		TypeMeta: v12.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: v12.ObjectMeta{
			Name:      "management-config",
			Namespace: TargetNamespace,
			Labels:    map[string]string{"cluster": "management"},
		},
		Data: map[string]string{
			"public_host": cmc.PublicHost,
			"public_port": cmc.PublicPort,
			"dns_host": cmc.DNSHost,
			"dns_port": cmc.DNSPort,
			"platform_type": cmc.PlatformType,
		},
	}

	client := cmc.Client.CoreV1().ConfigMaps(config.Namespace)
	log.Debug().Interface("configMap", config).Msg("creating management config")
	created, err := client.Create(config)
	if err != nil {
		return derrors.AsError(err, "cannot create configmap")
	}
	log.Debug().Interface("created", created).Msg("new config map has been created")
	return nil
}

func (cmc *CreateManagementConfig) createDockerSecret() derrors.Error {
	docker := &v1.Secret{
		TypeMeta: v12.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: v12.ObjectMeta{
			Name:      "docker-credentials",
			Namespace: TargetNamespace,
			Labels:    map[string]string{"cluster": "management"},
		},
		Data: map[string][]byte{
			"username": []byte(cmc.DockerUsername),
			"password": []byte(cmc.DockerPassword),
		},
		Type: v1.SecretTypeOpaque,
	}
	client := cmc.Client.CoreV1().Secrets(docker.Namespace)
	created, err := client.Create(docker)
	if err != nil {
		return derrors.AsError(err, "cannot create docker secret")
	}
	log.Debug().Interface("created", created).Msg("new secret has been created")
	return nil
}

func (cmc *CreateManagementConfig) createAuthSecret() derrors.Error {
	docker := &v1.Secret{
		TypeMeta: v12.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: v12.ObjectMeta{
			Name:      "authx-secret",
			Namespace: TargetNamespace,
			Labels:    map[string]string{"cluster": "management", "component": "authx"},
		},
		Data: map[string][]byte{
			"secret": []byte(uuid.NewV4().String()),
		},
		Type: v1.SecretTypeOpaque,
	}
	client := cmc.Client.CoreV1().Secrets(docker.Namespace)
	created, err := client.Create(docker)
	if err != nil {
		return derrors.AsError(err, "cannot create authx secret")
	}
	log.Debug().Interface("created", created).Msg("new secret has been created")
	return nil
}

func (cmc *CreateManagementConfig) Run(workflowID string) (*entities.CommandResult, derrors.Error) {
	connectErr := cmc.Connect()
	if connectErr != nil {
		return nil, connectErr
	}

	cErr := cmc.CreateNamespacesIfNotExist(TargetNamespace)
	if cErr != nil {
		return entities.NewCommandResult(false, "cannot create namespace", cErr), nil
	}

	err := cmc.createConfigMap()
	if err != nil {
		return entities.NewCommandResult(
			false, "cannot create management config", err), nil
	}

	err = cmc.createDockerSecret()
	if err != nil {
		return entities.NewCommandResult(
			false, "cannot create management config", err), nil
	}
	err = cmc.createAuthSecret()
	if err != nil {
		return entities.NewCommandResult(
			false, "cannot create management config", err), nil
	}

	return entities.NewSuccessCommand([]byte("management cluster config has been created")), nil
}

func (cmc *CreateManagementConfig) String() string {
	return fmt.Sprintf("SYNC CreateManagementConfig")
}

func (cmc *CreateManagementConfig) PrettyPrint(indentation int) string {
	simpleIden := strings.Repeat(" ", indentation) +  "  "
	entrySep := simpleIden +  "  "
	msg := fmt.Sprintf("\n%sConfig:\n%sPublicHost: %s:%s\n%sDNSHost: %s:%s\n%sDocker credentials: %s:%s\n%sPlatform Type:%s",
		simpleIden,
		entrySep, cmc.PublicHost, cmc.PublicPort,
		entrySep, cmc.DNSHost, cmc.DNSPort,
		entrySep, cmc.DockerUsername, strings.Repeat("*", len(cmc.DockerPassword)),
		entrySep, cmc.PlatformType,
	)
	return strings.Repeat(" ", indentation) + cmc.String() + msg
}

func (cmc *CreateManagementConfig) UserString() string {
	return fmt.Sprintf("Creating management cluster config with public address %s:%s", cmc.PublicHost, cmc.PublicPort)
}
