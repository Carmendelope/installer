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

		{"type":"sync", "name": "logger", "msg": "Installing Management component"},
		{"type":"sync", "name": "checkRequirements",
			"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
			"minVersion":"1.9"
		},
		{"type":"sync", "name": "launchComponents",
			"kubeConfigPath":"{{$.Credentials.KubeConfigPath}}",
			"namespaces":["nalej", "ingress-nginx"],
			"componentsDir":"{{$.Paths.ComponentsPath}}"
		}
	]
}
`

