/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package ingress

import (
	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var AzureDNSPort = v1.ServicePort{
	Name:       "dns-udp",
	Protocol:   v1.ProtocolUDP,
	Port:       53,
	TargetPort: intstr.IntOrString{
		Type:   intstr.String,
		StrVal: "dns-udp",
	},
}

var MinikubeDNSUDPPort = v1.ServicePort{
	Name:       "dns-udp",
	Protocol:   v1.ProtocolUDP,
	Port:       53,
	TargetPort: intstr.IntOrString{
		Type:   intstr.String,
		StrVal: "dns-udp",
	},
	NodePort: 53,
}

var MinikubeDNSTCPPort = v1.ServicePort{
	Name:       "dns-tcp",
	Protocol:   v1.ProtocolTCP,
	Port:       53,
	TargetPort: intstr.IntOrString{
		StrVal: "dns-tcp",
		Type:   intstr.String,
	},
	NodePort: 53,
}

var MinikubeDNSUIPort = v1.ServicePort{
	Name:       "consul-gui",
	Protocol:   v1.ProtocolTCP,
	Port:       8500,
	TargetPort: intstr.IntOrString{
		Type:   intstr.String,
		StrVal: "consul-gui",
	},
	NodePort: 30500,
}

var AzureConsulService = v1.Service{
	TypeMeta: metaV1.TypeMeta{
		Kind:       "Service",
		APIVersion: "v1",
	},
	ObjectMeta: metaV1.ObjectMeta{
		Name:      "dns-server-consul-dns",
		Namespace: "nalej",
		Labels: map[string]string{
			"cluster":                   "management",
			"component": "dns-server",
			"release" : "dns-server",
			"app": "consul",
		},
	},
	Spec: v1.ServiceSpec{
		Ports: []v1.ServicePort{AzureDNSPort},
		Selector: map[string]string{
			"app":    "consul",
			"hasDNS":"true",
			"release": "dns-server",
		},
		Type: v1.ServiceTypeLoadBalancer,
	},
}

var MinikubeConsulService = v1.Service{
	TypeMeta: metaV1.TypeMeta{
		Kind:       "Service",
		APIVersion: "v1",
	},
	ObjectMeta: metaV1.ObjectMeta{
		Name:      "dns-server-consul-dns",
		Namespace: "nalej",
		Labels: map[string]string{
			"cluster":                   "management",
			"component": "dns-server",
			"release" : "dns-server",
			"app": "consul",
		},
	},
	Spec: v1.ServiceSpec{
		Ports: []v1.ServicePort{MinikubeDNSUDPPort, MinikubeDNSTCPPort, MinikubeDNSUIPort},
		Selector: map[string]string{
			"app":    "consul",
			"hasDNS":"true",
			"release": "dns-server",
		},
		Type: v1.ServiceTypeNodePort,
	},
}
