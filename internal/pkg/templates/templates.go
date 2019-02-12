/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package templates

const InstallManagementCluster = `
{
	"description": "Install management cluster",
	"commands": [
		// Prerequirements
		{"type":"sync", "name":"checkAsset", "path":"{{$.Paths.BinaryPath}}/rke"},
		{"type":"sync", "name":"checkAsset", "path":"{{$.Paths.BinaryPath}}/zerotier-idtool"},
		{"type":"sync", "name":"checkAsset", "path":"{{$.Paths.BinaryPath}}/zerotier-cli"},
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
		{{else}}
			{"type":"sync", "name":"createManagementConfig",
				"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
				"public_host":"{{$.ManagementClusterHost}}",
				"public_port":"{{$.ManagementClusterPort}}",
				"dns_host":"{{$.DNSClusterHost}}",
				"dns_port":"{{$.DNSClusterPort}}",
				"docker_username":"{{$.Registry.Username}}",
				"docker_password":"{{$.Registry.Password}}",
				"platform_type":"{{$.InstallRequest.TargetPlatform}}"
			},
			{"type":"sync", "name":"installMngtDNS",
				"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
				"platform_type":"{{$.InstallRequest.TargetPlatform}}",
				"use_static_ip":{{$.InstallRequest.StaticIpAddresses.UseStaticIp}},
				"static_ip_address":"{{$.InstallRequest.StaticIpAddresses.Dns}}"
			},
		{{end}}
		{"type":"sync", "name":"installIngress",
				"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
				"management_public_host":"{{$.InstallRequest.Hostname}}",
				"on_management_cluster":{{ not $.AppClusterInstall}},
				"use_static_ip":{{$.InstallRequest.StaticIpAddresses.UseStaticIp}},
				"static_ip_address":"{{$.InstallRequest.StaticIpAddresses.Ingress}}"
		},
		{{if not $.AppClusterInstall }}
			{"type":"sync", "name":"installZtPlanetLB",
				"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
				"platform_type":"{{$.InstallRequest.TargetPlatform}}"
			},
			{"type":"sync", "name":"createZtPlanetConfig",
				"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
				"ztIdToolBinaryPath":"{{$.Paths.BinaryPath}}/zerotier-idtool",
				"management_public_host":"{{$.InstallRequest.Hostname}}",
				"identitySecretPath":"{{$.Paths.TempPath}}/identity.secret",
				"identityPublicPath":"{{$.Paths.TempPath}}/identity.public",
				"planetJsonPath":"{{$.Paths.TempPath}}/planet.json",
			},
		{{end}}
		{"type":"sync", "name":"createCredentials",
				"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
				"username":"{{$.Registry.Username}}",
				"password":"{{$.Registry.Password}}"
		},
		{"type":"sync", "name": "launchComponents",
			"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
			"namespaces":["nalej", "ingress-nginx"],
			"componentsDir":"{{$.Paths.ComponentsPath}}",
			"platform_type":"{{$.InstallRequest.TargetPlatform}}"
		}
	]
}
`
