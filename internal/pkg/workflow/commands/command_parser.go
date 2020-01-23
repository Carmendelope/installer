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

// This file contains the command parsing facilities to avoid import cycles.

package commands

import (
	"encoding/json"
	"github.com/nalej/derrors"
	"github.com/nalej/installer/internal/pkg/errors"
	"github.com/nalej/installer/internal/pkg/workflow/commands/async"
	"github.com/nalej/installer/internal/pkg/workflow/commands/sync"
	"github.com/nalej/installer/internal/pkg/workflow/commands/sync/istio"
	"github.com/nalej/installer/internal/pkg/workflow/commands/sync/k8s"
	"github.com/nalej/installer/internal/pkg/workflow/commands/sync/k8s/ingress"
	"github.com/nalej/installer/internal/pkg/workflow/commands/sync/rke"
	"github.com/nalej/installer/internal/pkg/workflow/commands/sync/zerotier"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
)

// CmdParser structure for the command parsing.
type CmdParser struct {
}

// NewCmdParser creates a new command parser.
func NewCmdParser() *CmdParser {
	return &CmdParser{}
}

// ParseCommand extracts a command from a raw JSON message.
func (cp *CmdParser) ParseCommand(raw []byte) (*entities.Command, derrors.Error) {
	var gc entities.GenericCommand
	if err := json.Unmarshal(raw, &gc); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	cmd, err := cp.parseCommand(gc, raw)
	if err != nil {
		return nil, err
	}
	return cmd, nil
}

func (cp *CmdParser) parseCommand(generic entities.GenericCommand, raw []byte) (*entities.Command, derrors.Error) {
	switch generic.CommandType {
	case entities.SyncCommandType:
		return cp.parseSyncCommand(generic, raw)
	case entities.AsyncCommandType:
		return cp.parseAsyncCommand(generic, raw)
	default:
		return nil, derrors.NewInvalidArgumentError(errors.UnsupportedCommandType).WithParams(generic)
	}
}

func (cp *CmdParser) parseSyncCommand(generic entities.GenericCommand, raw []byte) (*entities.Command, derrors.Error) {
	switch generic.CommandName {
	case entities.Exec:
		return sync.NewExecFromJSON(raw)
	case entities.SCP:
		return sync.NewSCPFromJSON(raw)
	case entities.SSH:
		return sync.NewSSHFromJSON(raw)
	case entities.Logger:
		return sync.NewLoggerFromJSON(raw)
	case entities.Sleep:
		return sync.NewSleepFromJSON(raw)
	case entities.Fail:
		return sync.NewFailFromJSON(raw)
	case entities.ParallelCmd:
		return NewParallelFromJSON(raw)
	case entities.GroupCmd:
		return NewGroupFromJSON(raw)
	case entities.TryCmd:
		return NewTryFromJSON(raw)
	case entities.ProcessCheck:
		return sync.NewProcessCheckFromJSON(raw)
	case entities.RKEInstall:
		return rke.NewRKEInstallFromJSON(raw)
	case entities.RKERemove:
		return rke.NewRKERemoveFromJSON(raw)
	case entities.CheckAsset:
		return sync.NewCheckAssetFromJSON(raw)
	case entities.LaunchComponents:
		return k8s.NewLaunchComponentsFromJSON(raw)
	case entities.CheckRequirements:
		return k8s.NewCheckRequirementsFromJSON(raw)
	case entities.CreateClusterConfig:
		return k8s.NewCreateClusterConfigFromJSON(raw)
	case entities.CreateManagementConfig:
		return k8s.NewCreateManagementConfigFromJSON(raw)
	case entities.UpdateCoreDNS:
		return k8s.NewUpdateCoreDNSFromJSON(raw)
	case entities.UpdateKubeDNS:
		return k8s.NewUpdateKubeDNSFromJSON(raw)
	case entities.CreateRegistrySecrets:
		return k8s.NewCreateRegistrySecretsFromJSON(raw)
	case entities.AddClusterUser:
		return k8s.NewAddClusterUserFromJSON(raw)
	case entities.InstallIngress:
		return ingress.NewInstallIngressFromJSON(raw)
	case entities.InstallMngtDNS:
		return ingress.NewInstallMngtDNSFromJSON(raw)
	case entities.InstallZtPlanetLB:
		return ingress.NewInstallZtPlanetLBFromJSON(raw)
	case entities.InstallVpnServerLB:
		return ingress.NewInstallVpnServerLBFromJSON(raw)
	case entities.CreateZTPlanetFiles:
		return zerotier.NewCreateZTPlanetFilesFromJSON(raw)
	case entities.CreateOpaqueSecret:
		return k8s.NewCreateOpaqueSecretFromJSON(raw)
	case entities.InstallExtDNS:
		return ingress.NewInstallExtDNSFromJSON(raw)
	case entities.CreateCACert:
		return k8s.NewCreateCACertFromJSON(raw)
	case entities.CreateTLSSecret:
		return k8s.NewCreateTLSSecretFromJSON(raw)
	case entities.DeleteNamespace:
		return k8s.NewDeleteNamespaceFromJSON(raw)
	case entities.DeleteNalejNamespace:
		return k8s.NewDeleteNalejNamespaceFromJSON(raw)
	case entities.DeleteServiceAccount:
		return k8s.NewDeleteServiceAccountFromJSON(raw)
	case entities.DeleteClusterRoleBinding:
		return k8s.NewDeleteClusterRoleBindingFromJSON(raw)
	case entities.DeleteClusterRole:
		return k8s.NewDeleteClusterRoleFromJSON(raw)
	case entities.DeleteRole:
		return k8s.NewDeleteRoleFromJSON(raw)
	case entities.DeleteRoleBinding:
		return k8s.NewDeleteRoleBindingFromJSON(raw)
	case entities.DeleteConfigMap:
		return k8s.NewDeleteConfigMapFromJSON(raw)
	case entities.DeleteService:
		return k8s.NewDeleteServiceFromJSON(raw)
	case entities.DeleteDeployment:
		return k8s.NewDeleteDeploymentFromJSON(raw)
	case entities.DeletePodSecurityPolicy:
		return k8s.NewDeletePodSecurityPolicyFromJSON(raw)
	case entities.InstallIstio:
		return istio.NewInstallIstioFromJSON(raw)
	default:
		return nil, derrors.NewInvalidArgumentError(errors.UnsupportedCommand).WithParams(generic)
	}
}

func (cp *CmdParser) parseAsyncCommand(generic entities.GenericCommand, raw []byte) (*entities.Command, derrors.Error) {
	switch generic.CommandName {
	case entities.Fail:
		return async.NewFailFromJSON(raw)
	case entities.Sleep:
		return async.NewSleepFromJSON(raw)
	default:
		return nil, derrors.NewInvalidArgumentError(errors.UnsupportedCommand).WithParams(generic)
	}
}
