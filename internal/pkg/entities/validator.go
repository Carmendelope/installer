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

package entities

import (
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-installer-go"
)

// ValidInstallRequest validates that an install request contains the required fields.
func ValidInstallRequest(installRequest *grpc_installer_go.InstallRequest) derrors.Error {
	if installRequest.InstallId == "" {
		return derrors.NewInvalidArgumentError("expecting install_id")
	}
	if installRequest.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("expecting organization_id")
	}
	if installRequest.ClusterId == "" {
		return derrors.NewInvalidArgumentError("expecting cluster_id")
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

// ValidInstallID checks that the request contains the required fields.
func ValidInstallID(installID *grpc_installer_go.InstallId) derrors.Error {
	if installID.InstallId == "" {
		return derrors.NewInvalidArgumentError("expecting install_id")
	}
	return nil
}

// ValidRemoveInstallRequest checks that the request contains the required fields.
func ValidRemoveInstallRequest(removeRequest *grpc_installer_go.RemoveInstallRequest) derrors.Error {
	if removeRequest.InstallId == "" {
		return derrors.NewInvalidArgumentError("expecting install_id")
	}
	return nil
}

func ValidUninstallClusterRequest(request *grpc_installer_go.UninstallClusterRequest) derrors.Error {
	if request.RequestId == "" {
		return derrors.NewInvalidArgumentError("expecting request_id")
	}
	if request.OrganizationId == "" {
		return derrors.NewInvalidArgumentError("expecting organization_id")
	}
	if request.ClusterId == "" {
		return derrors.NewInvalidArgumentError("expecting cluster_id")
	}
	if request.KubeConfigRaw == "" {
		return derrors.NewInvalidArgumentError("expecting kube_config_raw")
	}
	return nil
}
