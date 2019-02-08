package ingress

import (
	"encoding/json"
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-installer-go"
	"github.com/nalej/installer/internal/pkg/errors"
	"github.com/nalej/installer/internal/pkg/workflow/commands/sync/k8s"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
	"github.com/rs/zerolog/log"
	"strings"
)

type InstallZtPlanetLB struct {
	k8s.Kubernetes
	PlatformType    string `json:"platform_type"`
}

func NewInstallZtPlanetLB (kubeConfigPath string, platformType string) *InstallZtPlanetLB {
	return &InstallZtPlanetLB{
		Kubernetes: k8s.Kubernetes{
			GenericSyncCommand: *entities.NewSyncCommand(entities.InstallZtPlanetLB),
			KubeConfigPath:     kubeConfigPath,
		},
		PlatformType:    platformType,
	}
}

func NewInstallZtPlanetLBFromJSON (raw []byte) (*entities.Command, derrors.Error) {
	ccc := &InstallZtPlanetLB{}
	if err := json.Unmarshal(raw, &ccc); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	ccc.CommandID = entities.GenerateCommandID(ccc.Name())
	var r entities.Command = ccc
	return &r, nil
}

func (imd *InstallZtPlanetLB) Run (workflowID string) (*entities.CommandResult, derrors.Error) {
	connectErr := imd.Connect()
	if connectErr != nil {
		return nil, connectErr
	}

	switch imd.PlatformType {
	case grpc_installer_go.Platform_AZURE.String():
		return imd.InstallAzure(workflowID)
	case grpc_installer_go.Platform_MINIKUBE.String():
		return imd.InstallMinikube(workflowID)
	}
	log.Warn().Str("platformType", imd.PlatformType).Msg("unsupported platform type")
	return entities.NewCommandResult(
		false, "unsupported platform type", nil), nil
}

func (imd *InstallZtPlanetLB) InstallAzure (workflowID string) (*entities.CommandResult, derrors.Error) {
	azureService := AzureZTPlanetService
	err := imd.CreateService(&azureService)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating ZT Planet LB service")
		return entities.NewCommandResult(
			false, "cannot install service", err), nil
	}
	return entities.NewSuccessCommand([]byte("ZT planet installed on Azure")), nil
}

func (imd *InstallZtPlanetLB) InstallMinikube (workflowID string) (*entities.CommandResult, derrors.Error) {
	err := imd.CreateService(&MinikubeConsulService)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating  ZT Planet LB service")
		return entities.NewCommandResult(
			false, "cannot install service", err), nil
	}
	return entities.NewSuccessCommand([]byte("ZT planet installed on Minikube")), nil
}

func (imd *InstallZtPlanetLB) String () string {
	return fmt.Sprintf("SYNC InstallZtPlanetLB on %s", imd.PlatformType)
}

func (imd *InstallZtPlanetLB) PrettyPrint (indentation int) string {
	return strings.Repeat(" ", indentation) + imd.String()
}

func (imd *InstallZtPlanetLB) UserString () string {
	return fmt.Sprintf("Installing ZT planet loadbalancer")
}