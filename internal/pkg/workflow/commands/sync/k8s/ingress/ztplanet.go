/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package ingress

import (
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var AzureVPNServerPort = v1.ServicePort{
	Name:       "vpn-port",
	Protocol:   v1.ProtocolTCP,
	Port:       5555,
	TargetPort: intstr.IntOrString{
		Type:   intstr.Int,
		IntVal: 5555,
	},
}

var MinikubeVPNServerPort = v1.ServicePort{
	Name:       "vpn-port",
	Protocol:   v1.ProtocolTCP,
	Port:       5555,
	TargetPort: intstr.IntOrString{
		Type:   intstr.Int,
		IntVal: 5555,
	},
	NodePort: 5555,
}

var AzureVPNServerService = v1.Service{
	TypeMeta: metaV1.TypeMeta{
		Kind:       "Service",
		APIVersion: "v1",
	},
	ObjectMeta: metaV1.ObjectMeta{
		Name:      "vpn-server",
		Namespace: "nalej",
		Labels: map[string]string{
			"cluster":                   "management",
			"component": "vpn-server",
		},
	},
	Spec: v1.ServiceSpec{
		Ports: []v1.ServicePort{AzureVPNServerPort},
		Selector: map[string]string{
			"cluster":  "management",
			"component":	"vpn-server",
		},
		Type: v1.ServiceTypeLoadBalancer,
	},
}

var MinikubeVPNServerService = v1.Service{
	TypeMeta: metaV1.TypeMeta{
		Kind:       "Service",
		APIVersion: "v1",
	},
	ObjectMeta: metaV1.ObjectMeta{
		Name:      "vpn-servert",
		Namespace: "nalej",
		Labels: map[string]string{
			"cluster":	"management",
			"component": "vpn-server",
		},
	},
	Spec: v1.ServiceSpec{
		Ports: []v1.ServicePort{MinikubeVPNServerPort},
		Selector: map[string]string{
			"cluster":  "management",
			"component":	"vpn.-server",
		},
		Type: v1.ServiceTypeNodePort,
	},
}

