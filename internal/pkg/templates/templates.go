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
        {{if eq $.NetworkConfig.NetworkingMode "istio" }}
            {"type":"sync", "name":"installIstio",
                "kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
                "istio_path":"{{$.NetworkConfig.IstioPath}}",
                "cluster_id":"{{$.InstallRequest.ClusterId}}",
                "is_appCluster":{{$.AppCluster}},
                "static_ip_address":"{{$.InstallRequest.StaticIpAddresses.Ingress}}",
                "temp_path":"{{$.Paths.TempPath}}",
                "dns_public_host":"{{$.DNSClusterHost}}"
            },
        {{end}}
		{{if $.AppCluster }}
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
				"on_management_cluster":{{ not $.AppCluster}},
				"use_static_ip":{{$.InstallRequest.StaticIpAddresses.UseStaticIp}},
				"static_ip_address":"{{$.InstallRequest.StaticIpAddresses.Ingress}}",
                "network_mode":"{{$.NetworkConfig.NetworkingMode}}"
		},
		{{if not $.AppCluster }}
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
		[{"type":"sync", "name": "launchComponents",
			"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
			"namespaces":["nalej", "ingress-nginx"],
			"componentsDir":"{{$.Paths.ComponentsPath}}",
			"platform_type":"{{$.InstallRequest.TargetPlatform}}",
			"environment":"{{$.TargetEnvironment}}"
		},
		{"type":"sync", "name": "checkComponents",
			"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
			"namespaces":["nalej", "ingress-nginx"]
		}]
	]
}
`

// UninstallCluster template with the commands required to uninstall the Nalej platform
const UninstallCluster = `
{
	"description": "Uninstall management cluster",
	"commands": [
		{"type":"sync", "name": "logger", "msg": "Checking requirements"},
		{"type":"sync", "name": "checkRequirements",
			"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
			"minVersion":"1.11"
		},
		{"type":"sync", "name": "logger", "msg": "Uninstalling components"},
		{"type":"sync", "name":"deleteServiceAccount",
			"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
			"namespace":"kube-system",
			"service_account":"nginx-ingress",
			"fail_if_not_exists":false
		},
		{"type":"sync", "name":"deleteNamespace",
			"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
			"namespace":"ingress-nginx",
			"fail_if_not_exists":false
		},
		{"type":"sync", "name":"deleteNalejNamespace",
			"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
			"fail_if_not_exists":false
		},
		{"type":"sync", "name":"deleteClusterRoleBinding",
			"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
			"role_binding_name":"system:nginx-ingress",
			"fail_if_not_exists":false
		},

		{{if not $.AppCluster }}
			{"type":"sync", "name":"deleteClusterRoleBinding",
				"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
				"role_binding_name":"deployment-manager",
				"fail_if_not_exists":false
			},
			{"type":"sync", "name":"deleteClusterRoleBinding",
				"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
				"role_binding_name":"kube-state-metrics",
				"fail_if_not_exists":false
			},
			{"type":"sync", "name":"deleteClusterRoleBinding",
				"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
				"role_binding_name":"node-exporter",
				"fail_if_not_exists":false
			},
			{"type":"sync", "name":"deleteClusterRoleBinding",
				"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
				"role_binding_name":"prometheus",
				"fail_if_not_exists":false
			},
			{"type":"sync", "name":"deleteClusterRoleBinding",
				"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
				"role_binding_name":"filebeat",
				"fail_if_not_exists":false
			},
		{{end}}
		{"type":"sync", "name":"deleteClusterRole",
			"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
			"role_name":"system:nginx-ingress",
			"fail_if_not_exists":false
		},
		{{if not $.AppCluster }}
			{"type":"sync", "name":"deleteClusterRole",
				"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
				"role_name":"kube-state-metrics",
				"fail_if_not_exists":false
			},
			{"type":"sync", "name":"deleteClusterRole",
				"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
				"role_name":"node-exporter",
				"fail_if_not_exists":false
			},
			{"type":"sync", "name":"deleteClusterRole",
				"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
				"role_name":"prometheus",
				"fail_if_not_exists":false
			},
			{"type":"sync", "name":"deleteClusterRole",
				"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
				"role_name":"filebeat",
				"fail_if_not_exists":false
			},
		{{end}}
		{"type":"sync", "name":"deleteRole",
			"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
			"namespace":"kube-system",
			"role_name":"system::nginx-ingress-role",
			"fail_if_not_exists":false
		},
		{"type":"sync", "name":"deleteRoleBinding",
			"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
			"namespace":"kube-system",
			"role_name":"system::nginx-ingress-role-binding",
			"fail_if_not_exists":false
		},
		{"type":"sync", "name":"deleteConfigMap",
			"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
			"namespace":"kube-system",
			"config_map_name":"ingress-controller-leader-nginx",
			"fail_if_not_exists":false
		},
		{"type":"sync", "name":"deleteConfigMap",
			"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
			"namespace":"kube-system",
			"config_map_name":"nginx-load-balancer-conf",
			"fail_if_not_exists":false
		},
		{"type":"sync", "name":"deleteConfigMap",
			"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
			"namespace":"kube-system",
			"config_map_name":"tcp-services",
			"fail_if_not_exists":false
		},
		{"type":"sync", "name":"deleteConfigMap",
			"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
			"namespace":"kube-system",
			"config_map_name":"udp-services",
			"fail_if_not_exists":false
		},
		{"type":"sync", "name":"deleteService",
			"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
			"namespace":"kube-system",
			"service_name":"default-http-backend",
			"fail_if_not_exists":false
		},
		{"type":"sync", "name":"deleteService",
			"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
			"namespace":"kube-system",
			"service_name":"nginx-ingress-controller",
			"fail_if_not_exists":false
		},
		{"type":"sync", "name":"deleteDeployment",
			"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
			"namespace":"kube-system",
			"deployment_name":"default-http-backend",
			"fail_if_not_exists":false
		},
		{"type":"sync", "name":"deleteDeployment",
			"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
			"namespace":"kube-system",
			"deployment_name":"nginx-ingress-controller",
			"fail_if_not_exists":false
		},
		{"type":"sync", "name":"deletePodSecurityPolicy",
			"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
			"policy_name":"node-exporter",
			"fail_if_not_exists":false
		}
	]
}
`
