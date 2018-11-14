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
	"k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type CreateManagementConfig struct {
	Kubernetes
	PublicHost string `json:"public_host"`
	PublicPort string `json:"public_port"`
}

func NewCreateManagementConfig(kubeConfigPath string, publicHost string, publicPort string) *CreateManagementConfig {
	return &CreateManagementConfig{
		Kubernetes: Kubernetes{
			GenericSyncCommand: *entities.NewSyncCommand(entities.CreateManagementConfig),
			KubeConfigPath:     kubeConfigPath,
		},
		PublicHost: publicHost,
		PublicPort: publicPort,
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

func (cmc *CreateManagementConfig) Run(workflowID string) (*entities.CommandResult, derrors.Error) {
	connectErr := cmc.Connect()
	if connectErr != nil {
		return nil, connectErr
	}

	cErr := cmc.createNamespacesIfNotExist("nalej")
	if cErr != nil {
		return entities.NewCommandResult(false, "cannot create namespace", cErr), nil
	}

	config := &v1.ConfigMap{
		TypeMeta: v12.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: v12.ObjectMeta{
			Name:      "management-config",
			Namespace: "nalej",
			Labels:    map[string]string{"cluster": "management"},
		},
		Data: map[string]string{
			"public_host": cmc.PublicHost,
			"public_port": cmc.PublicPort},
	}

	client := cmc.Client.CoreV1().ConfigMaps(config.Namespace)
	log.Debug().Interface("configMap", config).Msg("creating management config")
	created, err := client.Create(config)
	if err != nil {
		return entities.NewCommandResult(
			false, "cannot create management config", derrors.AsError(err, "cannot create configmap")), nil
	}
	log.Debug().Interface("created", created).Msg("new config map has been created")
	return entities.NewSuccessCommand([]byte("management cluster config has been created")), nil
}

func (cmc *CreateManagementConfig) String() string {
	return fmt.Sprintf("SYNC CreateManagementConfig publicHost: %s, publicPort: %s", cmc.PublicHost, cmc.PublicPort)
}

func (cmc *CreateManagementConfig) PrettyPrint(indentation int) string {
	return strings.Repeat(" ", indentation) + cmc.String()
}

func (cmc *CreateManagementConfig) UserString() string {
	return fmt.Sprintf("Creating management cluster config with public address %s:%s", cmc.PublicHost, cmc.PublicPort)
}
