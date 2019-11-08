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

package templates

const InstallManagementCluster = `
{
	"description": "Install management cluster",
	"commands": [
		// Prerequirements
		{"type":"sync", "name":"checkAsset", "path":"{{$.Paths.BinaryPath}}/rke"},
		// Install K8s
		{{if $.InstallRequest.InstallBaseSystem }}
			{"type":"sync", "name": "logger", "msg": "Installing base system"},
		{{end}}

		{"type":"sync", "name": "logger", "msg": "Checking requirements"},
		{"type":"sync", "name": "checkRequirements",
			"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
			"minVersion":"1.11"
		},
		{"type":"sync", "name": "logger", "msg": "Installing components"},
		{{if $.AppClusterInstall }}
			{"type":"sync", "name":"createClusterConfig",
				"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
				"organization_id":"{{$.InstallRequest.OrganizationId}}",
				"cluster_id":"{{$.InstallRequest.ClusterId}}",
				"management_public_host":"{{$.ManagementClusterHost}}",
				"management_public_port":"{{$.ManagementClusterPort}}",
				"cluster_public_hostname":"{{$.InstallRequest.Hostname}}",
				"dns_public_host":"{{$.DNSClusterHost}}",
				"dns_public_port":"{{$.DNSClusterPort}}",
				"platform_type":"{{$.InstallRequest.TargetPlatform}}"
			},
			{"type":"sync", "name":"addClusterUser",
				"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
				"organization_id":"{{$.InstallRequest.OrganizationId}}",
				"cluster_id":"{{$.InstallRequest.ClusterId}}",
				"user_manager_address":"user-manager.nalej:8920"
			},
			{"type":"sync", "name":"createOpaqueSecret",
				"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
				"secret_name":"authx-secret",
				"secret_key":"secret",
				"load_from_path":false,
				"secret_value":"{{$.AuthSecret}}"
			},
			{"type":"sync", "name":"createOpaqueSecret",
				"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
				"secret_name":"ca-certificate",
				"secret_key":"ca.crt",
				"load_from_path":true,
				"secret_value_from_path":"{{$.CACertPath}}"
			},
		{{else}}
			{"type":"sync", "name":"createManagementConfig",
				"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
				"public_host":"{{$.ManagementClusterHost}}",
				"public_port":"{{$.ManagementClusterPort}}",
				"dns_host":"{{$.DNSClusterHost}}",
				"dns_port":"{{$.DNSClusterPort}}",
				"platform_type":"{{$.InstallRequest.TargetPlatform}}",
				"environment":"{{$.TargetEnvironment}}"
			},
			{"type":"sync", "name":"installMngtDNS",
				"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
				"platform_type":"{{$.InstallRequest.TargetPlatform}}",
				"use_static_ip":{{$.InstallRequest.StaticIpAddresses.UseStaticIp}},
				"static_ip_address":"{{$.InstallRequest.StaticIpAddresses.Dns}}"
			},
			{"type":"sync", "name":"createCACert",
				"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
				"public_host":"{{$.ManagementClusterHost}}"
			},
		{{end}}
		{"type":"sync", "name":"installIngress",
				"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
				"platform_type":"{{$.InstallRequest.TargetPlatform}}",
				"management_public_host":"{{$.InstallRequest.Hostname}}",
				"on_management_cluster":{{ not $.AppClusterInstall}},
				"use_static_ip":{{$.InstallRequest.StaticIpAddresses.UseStaticIp}},
				"static_ip_address":"{{$.InstallRequest.StaticIpAddresses.Ingress}}"
		},
		{{if not $.AppClusterInstall }}
			{"type":"sync", "name":"installExtDNS",
				"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
				"platform_type":"{{$.InstallRequest.TargetPlatform}}",
				"use_static_ip":{{$.InstallRequest.StaticIpAddresses.UseStaticIp}},
				"static_ip_address":"{{$.InstallRequest.StaticIpAddresses.CorednsExt}}"
			},
			{"type":"sync", "name":"installVpnServerLB",
				"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
				"platform_type":"{{$.InstallRequest.TargetPlatform}}",
				"use_static_ip":{{$.InstallRequest.StaticIpAddresses.UseStaticIp}},
				"static_ip_address":"{{$.InstallRequest.StaticIpAddresses.VpnServer}}"
			},
		{{end}}

		{{ range $index, $registry := $.Registries }}
			{"type":"sync", "name":"createRegistrySecrets",
				"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
				"on_management_cluster":{{ not $.AppClusterInstall}},
				"credentials_name":"{{$registry.Name}}",
				"username":"{{$registry.Username}}",
				"password":"{{$registry.Password}}",
				"url":"{{$registry.URL}}"
			},
		{{end}}

		{"type":"sync", "name": "launchComponents",
			"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
			"namespaces":["nalej", "ingress-nginx"],
			"componentsDir":"{{$.Paths.ComponentsPath}}",
			"platform_type":"{{$.InstallRequest.TargetPlatform}}",
			"environment":"{{$.TargetEnvironment}}"
		}
	]
}
`
