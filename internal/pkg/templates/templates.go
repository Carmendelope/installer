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
			// {{if $.InstallRequest.UseCoreDns }}			
			//	{"type":"sync", "name":"updateCoreDNS",
			//		"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
			//		"dns_public_host":"{{$.DNSClusterHost}}",
			//		"dns_public_port":"{{$.DNSClusterPort}}"
			//	},
			// {{end}}
			// {{if $.InstallRequest.UseKubeDns }}			
			//	{"type":"sync", "name":"updateKubeDNS",
			//		"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
			//		"dns_public_host":"{{$.DNSClusterHost}}"
			//	},
			// {{end}}

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
				"platform_type":"{{$.InstallRequest.TargetPlatform}}"
			},
		{{end}}
		{"type":"sync", "name":"installIngress",
				"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
				"management_public_host":"{{$.InstallRequest.Hostname}}",
				"on_management_cluster":{{ not $.AppClusterInstall}}
		},
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

