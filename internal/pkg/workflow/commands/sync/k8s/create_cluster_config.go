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

package k8s

import (
	"encoding/json"
	"fmt"
	"github.com/rs/zerolog/log"
	"strings"

	"github.com/nalej/derrors"
	"github.com/nalej/installer/internal/pkg/errors"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
	"k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CreateClusterConfig struct {
	Kubernetes
	OrganizationID        string `json:"organization_id"`
	ClusterID             string `json:"cluster_id"`
	ManagementPublicHost  string `json:"management_public_host"`
	ManagementPublicPort  string `json:"management_public_port"`
	ClusterPublicHostname string `json:"cluster_public_hostname"`
	DNSPublicHost         string `json:"dns_public_host"`
	DNSPublicPort         string `json:"dns_public_port"`
	PlatformType          string `json:"platform_type"`
}

func NewCreateClusterConfig(
	kubeConfigPath string,
	organizationID string, clusterID string,
	managementPublicHost string, managementPublicPort string,
	clusterPublicHostname string,
	dnsPublicHost string, dnsPublicPort string,
	platformType string,
) *CreateClusterConfig {
	return &CreateClusterConfig{
		Kubernetes: Kubernetes{
			GenericSyncCommand: *entities.NewSyncCommand(entities.CreateClusterConfig),
			KubeConfigPath:     kubeConfigPath,
		},
		OrganizationID:        organizationID,
		ClusterID:             clusterID,
		ManagementPublicHost:  managementPublicHost,
		ManagementPublicPort:  managementPublicPort,
		ClusterPublicHostname: clusterPublicHostname,
		DNSPublicHost:         dnsPublicHost,
		DNSPublicPort:         dnsPublicPort,
		PlatformType:          platformType,
	}
}

func NewCreateClusterConfigFromJSON(raw []byte) (*entities.Command, derrors.Error) {
	ccc := &CreateClusterConfig{}
	if err := json.Unmarshal(raw, &ccc); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	ccc.CommandID = entities.GenerateCommandID(ccc.Name())
	var r entities.Command = ccc
	return &r, nil
}

func (ccc *CreateClusterConfig) Run(workflowID string) (*entities.CommandResult, derrors.Error) {
	connectErr := ccc.Connect()
	if connectErr != nil {
		log.Error().Str("connection error", connectErr.DebugReport()).Str("connection error", connectErr.DebugReport())
		return nil, connectErr
	}

	mgntIPs, rErr := ccc.ResolveIP(ccc.ManagementPublicHost)
	if rErr != nil {
		log.Error().Str("resolve ip error mngt", rErr.DebugReport()).Str("resolve ip error mngt", rErr.DebugReport())
		return nil, rErr
	}

	dnsIPs, rErr := ccc.ResolveIP(ccc.DNSPublicHost)
	if rErr != nil {
		log.Error().Str("resolve ip error dns", rErr.DebugReport()).Str("resolve ip error dns", rErr.DebugReport())
		return nil, rErr
	}

	cErr := ccc.CreateNamespaceIfNotExists("nalej")
	if cErr != nil {
		log.Error().Str("namespace creation error", cErr.DebugReport()).Str("namespace creation error", cErr.DebugReport())
		return entities.NewCommandResult(false, "cannot create namespace", cErr), nil
	}

	log.Debug().Str("creating namespace", "nalej").Msg("creating namespace nalej")
	config := &v1.ConfigMap{
		TypeMeta: v12.TypeMeta{
			Kind:       "ConfigMap",
			APIVersion: "v1",
		},
		ObjectMeta: v12.ObjectMeta{
			Name:      "cluster-config",
			Namespace: "nalej",
			Labels:    map[string]string{"cluster": "application"},
		},
		Data: map[string]string{
			"organization_id":         ccc.OrganizationID,
			"cluster_id":              ccc.ClusterID,
			"management_public_host":  ccc.ManagementPublicHost,
			"management_public_ip":    strings.Join(mgntIPs, ","),
			"management_public_port":  ccc.ManagementPublicPort,
			"cluster_public_hostname": ccc.ClusterPublicHostname,
			"cluster_api_hostname":    fmt.Sprintf("cluster.%s", ccc.ManagementPublicHost),
			"login_api_hostname":      fmt.Sprintf("login.%s", ccc.ManagementPublicHost),
			"dns_public_ips":          strings.Join(dnsIPs, ","),
			"dns_public_port":         ccc.DNSPublicPort,
			"platform_type":           ccc.PlatformType,
		},
	}

	derr := ccc.Create(config)
	if derr != nil {
		return entities.NewCommandResult(false, "cannot create cluster config", derr), nil
	}
	return entities.NewSuccessCommand([]byte("cluster config has been created")), nil
}

func (ccc *CreateClusterConfig) String() string {
	return fmt.Sprintf("SYNC CreateClusterConfig organizationID: %s, clusterID: %s", ccc.OrganizationID, ccc.ClusterID)
}

func (ccc *CreateClusterConfig) PrettyPrint(indentation int) string {
	simpleIden := strings.Repeat(" ", indentation)
	entrySep := simpleIden + "  "
	msg := fmt.Sprintf("\n%sConfig:\n%sManagementPublicHost: %s\n%sClusterPublicHostname: %s",
		entrySep,
		entrySep, ccc.ManagementPublicHost,
		entrySep, ccc.ClusterPublicHostname,
	)
	return strings.Repeat(" ", indentation) + ccc.String() + msg
}

func (ccc *CreateClusterConfig) UserString() string {
	return fmt.Sprintf("Creating cluster config for %s", ccc.ClusterID)
}
