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

package installer

import (
	"context"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-installer-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	Manager Manager
}

func NewHandler(manager Manager) *Handler {
	return &Handler{manager}
}

func (h *Handler) ValidInstallRequest(installRequest *grpc_installer_go.InstallRequest) derrors.Error {
	if installRequest.InstallId == "" {
		return derrors.NewInvalidArgumentError("expecting InstallId")
	}
	if installRequest.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("expecting OrganizationId")
	}
	if installRequest.ClusterId == "" {
		return derrors.NewInvalidArgumentError("expecting ClusterId")
	}
	if installRequest.Hostname == "" {
		return derrors.NewInvalidArgumentError("hostname must be set with the ingress hostname")
	}
	authFound := false

	if installRequest.Username != "" {
		if installRequest.PrivateKey == "" {
			return derrors.NewInvalidArgumentError("expecting PrivateKey with Username")
		}
		if len(installRequest.Nodes) == 0 {
			return derrors.NewInvalidArgumentError("expecting Nodes with Username")
		}
		authFound = true
	}
	if installRequest.KubeConfigRaw != "" {
		if installRequest.Username != "" {
			return derrors.NewInvalidArgumentError("expecting KubeConfigRaw without Username")
		}
		if installRequest.PrivateKey != "" {
			return derrors.NewInvalidArgumentError("expecting KubeConfigRaw without PrivateKey")
		}
		if len(installRequest.Nodes) > 0 {
			return derrors.NewInvalidArgumentError("expecting KubeConfigRaw without Nodes")
		}
		authFound = true
	}
	if !authFound {
		return derrors.NewInvalidArgumentError("expecting KubeConfigRaw or Username, PrivateKey and Nodes")
	}

	return nil
}

func (h *Handler) ValidInstallID(installID *grpc_installer_go.InstallId) derrors.Error {
	if installID.InstallId == "" {
		return derrors.NewInvalidArgumentError("expecting InstallId")
	}
	return nil
}

func (h *Handler) ValidRemoveInstallRequest(removeRequest *grpc_installer_go.RemoveInstallRequest) derrors.Error {
	if removeRequest.InstallId == "" {
		return derrors.NewInvalidArgumentError("expecting InstallId")
	}
	return nil
}

func (h *Handler) InstallCluster(ctx context.Context, installRequest *grpc_installer_go.InstallRequest) (*grpc_installer_go.InstallResponse, error) {
	log.Debug().Str("organizationID", installRequest.OrganizationId).Str("installID", installRequest.InstallId).Msg("install cluster")
	err := h.ValidInstallRequest(installRequest)
	if err != nil {
		log.Warn().Str("trace", err.DebugReport()).Msg(err.Error())
		return nil, conversions.ToGRPCError(err)
	}
	status, err := h.Manager.InstallCluster(*installRequest)
	if err != nil {
		log.Warn().Str("trace", err.DebugReport()).Msg(err.Error())
		return nil, conversions.ToGRPCError(err)
	}
	log.Debug().Str("organizationID", installRequest.OrganizationId).Str("installID", installRequest.InstallId).Msg("install launched")
	return status.ToGRPCInstallResponse(), nil
}

func (h *Handler) CheckProgress(ctx context.Context, installID *grpc_installer_go.InstallId) (*grpc_installer_go.InstallResponse, error) {
	err := h.ValidInstallID(installID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	status, err := h.Manager.GetProgress(installID.InstallId)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return status.ToGRPCInstallResponse(), nil
}

func (h *Handler) RemoveInstall(ctx context.Context, removeRequest *grpc_installer_go.RemoveInstallRequest) (*grpc_common_go.Success, error) {
	err := h.ValidRemoveInstallRequest(removeRequest)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	err = h.Manager.RemoveInstall(removeRequest.InstallId)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return &grpc_common_go.Success{}, nil
}
