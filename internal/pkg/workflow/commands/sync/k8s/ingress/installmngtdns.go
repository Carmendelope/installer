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

package ingress

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nalej/derrors"
	"github.com/nalej/grpc-installer-go"
	"github.com/nalej/installer/internal/pkg/errors"
	"github.com/nalej/installer/internal/pkg/workflow/commands/sync/k8s"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
	"github.com/rs/zerolog/log"
)

type InstallMngtDNS struct {
	k8s.Kubernetes
	PlatformType    string `json:"platform_type"`
	UseStaticIp     bool   `json:"use_static_ip"`
	StaticIpAddress string `json:"static_ip_address"`
}

func NewInstallMngtDNS(kubeConfigPath string, platformType string, useStaticIp bool, staticIpAddress string) *InstallMngtDNS {
	return &InstallMngtDNS{
		Kubernetes: k8s.Kubernetes{
			GenericSyncCommand: *entities.NewSyncCommand(entities.InstallMngtDNS),
			KubeConfigPath:     kubeConfigPath,
		},
		PlatformType:    platformType,
		UseStaticIp:     useStaticIp,
		StaticIpAddress: staticIpAddress,
	}
}

func NewInstallMngtDNSFromJSON(raw []byte) (*entities.Command, derrors.Error) {
	ccc := &InstallMngtDNS{}
	if err := json.Unmarshal(raw, &ccc); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	ccc.CommandID = entities.GenerateCommandID(ccc.Name())
	var r entities.Command = ccc
	return &r, nil
}

func (imd *InstallMngtDNS) Run(workflowID string) (*entities.CommandResult, derrors.Error) {
	connectErr := imd.Connect()
	if connectErr != nil {
		return nil, connectErr
	}

	switch imd.PlatformType {
	case grpc_installer_go.Platform_AZURE.String():
		return imd.InstallLoadBalancer(workflowID)
	case grpc_installer_go.Platform_BAREMETAL.String():
		// The baremetal type relies on MetalLB so it supports loadbalancers as in Azure.
		return imd.InstallLoadBalancer(workflowID)
	case grpc_installer_go.Platform_MINIKUBE.String():
		return imd.InstallMinikube(workflowID)
	}
	log.Warn().Str("platformType", imd.PlatformType).Msg("unsupported platform type")
	return entities.NewCommandResult(
		false, "unsupported platform type", nil), nil
}

func (imd *InstallMngtDNS) InstallLoadBalancer(workflowID string) (*entities.CommandResult, derrors.Error) {
	azureService := AzureConsulService
	if imd.UseStaticIp {
		azureService.Spec.LoadBalancerIP = imd.StaticIpAddress
	}
	err := imd.Create(&azureService)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating DNS service")
		return entities.NewCommandResult(
			false, "cannot install service", err), nil
	}
	msg := fmt.Sprintf("DNS loadbalancer installed on %s", imd.PlatformType)
	return entities.NewSuccessCommand([]byte(msg)), nil
}

func (imd *InstallMngtDNS) InstallMinikube(workflowID string) (*entities.CommandResult, derrors.Error) {
	err := imd.Create(&MinikubeConsulService)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating DNS service")
		return entities.NewCommandResult(
			false, "cannot install service", err), nil
	}
	return entities.NewSuccessCommand([]byte("DNS loadbalancer installed on Minikube")), nil
}

func (imd *InstallMngtDNS) String() string {
	return fmt.Sprintf("SYNC InstallMngtDNS on %s", imd.PlatformType)
}

func (imd *InstallMngtDNS) PrettyPrint(indentation int) string {
	return strings.Repeat(" ", indentation) + imd.String()
}

func (imd *InstallMngtDNS) UserString() string {
	return fmt.Sprintf("Installing DNS loadbalancer")
}
