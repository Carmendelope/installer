/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
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

type InstallExtDNS struct {
	k8s.Kubernetes
	PlatformType    string `json:"platform_type"`
	UseStaticIp     bool   `json:"use_static_ip"`
	StaticIpAddress string `json:"static_ip_address"`
}

func NewInstallExtDNS(kubeConfigPath string, platformType string, useStaticIp bool, staticIpAddress string) *InstallExtDNS {
	return &InstallExtDNS{
		Kubernetes: k8s.Kubernetes{
			GenericSyncCommand: *entities.NewSyncCommand(entities.InstallExtDNS),
			KubeConfigPath:     kubeConfigPath,
		},
		PlatformType:    platformType,
		UseStaticIp:     useStaticIp,
		StaticIpAddress: staticIpAddress,
	}
}

func NewInstallExtDNSFromJSON(raw []byte) (*entities.Command, derrors.Error) {
	ccc := &InstallExtDNS{}
	if err := json.Unmarshal(raw, &ccc); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	ccc.CommandID = entities.GenerateCommandID(ccc.Name())
	var r entities.Command = ccc
	return &r, nil
}

func (imd *InstallExtDNS) Run(workflowID string) (*entities.CommandResult, derrors.Error) {
	connectErr := imd.Connect()
	if connectErr != nil {
		return nil, connectErr
	}

	switch imd.PlatformType {
	case grpc_installer_go.Platform_AZURE.String():
		return imd.InstallAzure(workflowID)
	case grpc_installer_go.Platform_BAREMETAL.String():
		// The baremetal type relies on MetalLB so it supports loadbalancers as in Azure.
		return imd.InstallAzure(workflowID)
	case grpc_installer_go.Platform_MINIKUBE.String():
		return imd.InstallMinikube(workflowID)
	}
	log.Warn().Str("platformType", imd.PlatformType).Msg("unsupported platform type")
	return entities.NewCommandResult(
		false, "unsupported platform type", nil), nil
}

func (imd *InstallExtDNS) InstallAzure(workflowID string) (*entities.CommandResult, derrors.Error) {
	azureService := AzureExtDnsService
	if imd.UseStaticIp {
		azureService.Spec.LoadBalancerIP = imd.StaticIpAddress
	}
	err := imd.CreateService(&azureService)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating External DNS service")
		return entities.NewCommandResult(
			false, "cannot install service", err), nil
	}
	msg := fmt.Sprintf("External DNS loadbalancer installed on %s", imd.PlatformType)
	return entities.NewSuccessCommand([]byte(msg)), nil
}

func (imd *InstallExtDNS) InstallMinikube(workflowID string) (*entities.CommandResult, derrors.Error) {
	err := imd.CreateService(&MinikubeExtDnsService)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating External DNS service")
		return entities.NewCommandResult(
			false, "cannot install service", err), nil
	}
	return entities.NewSuccessCommand([]byte("External DNS loadbalancer installed on Minikube")), nil
}

func (imd *InstallExtDNS) String() string {
	return fmt.Sprintf("SYNC InstallExtDNS on %s", imd.PlatformType)
}

func (imd *InstallExtDNS) PrettyPrint(indentation int) string {
	return strings.Repeat(" ", indentation) + imd.String()
}

func (imd *InstallExtDNS) UserString() string {
	return fmt.Sprintf("Installing External DNS loadbalancer")
}
