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
	"github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-installer-go"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/installer/internal/pkg/entities"
	"github.com/rs/zerolog/log"
)

type Handler struct {
	Manager Manager
}

func NewHandler(manager Manager) *Handler {
	return &Handler{manager}
}

// InstallCluster triggers the installation of a new application cluster.
func (h *Handler) InstallCluster(ctx context.Context, installRequest *grpc_installer_go.InstallRequest) (*grpc_common_go.OpResponse, error) {
	log.Debug().Str("organizationID", installRequest.OrganizationId).Str("installID", installRequest.RequestId).Msg("install cluster")
	err := entities.ValidInstallRequest(installRequest)
	if err != nil {
		log.Warn().Str("trace", err.DebugReport()).Msg(err.Error())
		return nil, conversions.ToGRPCError(err)
	}
	status, err := h.Manager.InstallCluster(*installRequest)
	if err != nil {
		log.Warn().Str("trace", err.DebugReport()).Msg(err.Error())
		return nil, conversions.ToGRPCError(err)
	}
	log.Debug().Str("organizationID", installRequest.OrganizationId).Str("requestID", installRequest.RequestId).Msg("install launched")
	return status.ToGRPCOpResponse(), nil
}

// UninstallCluster proceeds to remove all Nalej created elements in that cluster.
func (h *Handler) UninstallCluster(ctx context.Context, request *grpc_installer_go.UninstallClusterRequest) (*grpc_common_go.OpResponse, error) {
	log.Debug().Str("organizationID", request.OrganizationId).Str("requestID", request.RequestId).Msg("uninstall cluster")
	err := entities.ValidUninstallClusterRequest(request)
	if err != nil {
		log.Warn().Str("trace", err.DebugReport()).Msg(err.Error())
		return nil, conversions.ToGRPCError(err)
	}
	response, err := h.Manager.UninstallCluster(*request)
	if err != nil {
		log.Warn().Str("trace", err.DebugReport()).Msg(err.Error())
		return nil, conversions.ToGRPCError(err)
	}
	log.Debug().Str("organizationID", request.OrganizationId).Str("requestID", request.RequestId).Msg("uninstall launched")
	return response.ToGRPCOpResponse(), nil
}

// CheckProgress gets an updated state of an install request.
func (h *Handler) CheckProgress(ctx context.Context, requestID *grpc_common_go.RequestId) (*grpc_common_go.OpResponse, error) {
	err := entities.ValidRequestID(requestID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	status, err := h.Manager.GetProgress(requestID.RequestId)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return status.ToGRPCOpResponse(), nil
}

// RemoveInstall cancels and ongoing install or removes the information of an already processed install.
func (h *Handler) RemoveInstall(ctx context.Context, requestID *grpc_common_go.RequestId) (*grpc_common_go.Success, error) {
	err := entities.ValidRequestID(requestID)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	err = h.Manager.RemoveInstall(requestID.RequestId)
	if err != nil {
		return nil, conversions.ToGRPCError(err)
	}
	return &grpc_common_go.Success{}, nil
}
