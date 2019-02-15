/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

// This file contains the command parsing facilities to avoid import cycles.

package commands

import (
	"encoding/json"
	"github.com/nalej/installer/internal/pkg/errors"
	"github.com/nalej/installer/internal/pkg/workflow/commands/sync/k8s"
	"github.com/nalej/installer/internal/pkg/workflow/commands/sync/k8s/ingress"

	"github.com/nalej/derrors"
	"github.com/nalej/installer/internal/pkg/workflow/commands/async"
	"github.com/nalej/installer/internal/pkg/workflow/commands/sync"
	"github.com/nalej/installer/internal/pkg/workflow/commands/sync/rke"
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
