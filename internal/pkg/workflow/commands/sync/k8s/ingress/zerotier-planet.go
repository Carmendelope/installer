/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package ingress

import (
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

//apiVersion: v1
//kind: Service
//metadata:
//  labels:
//    cluster: management
//    component: network-manager
//  name: zerotier-planet
//  namespace: nalej
//spec:
//  externalTrafficPolicy: Cluster
//  ports:
//  - name: zt-udp
//    port: 9993
//    protocol: UDP
//    targetPort: 9993
//  selector:
//    cluster: management
//    component: network-manager
//  type: LoadBalancer

var AzureZTPort = v1.ServicePort{
	Name:       "zt-udp",
	Protocol:   v1.ProtocolUDP,
	Port:       9993,
	TargetPort: intstr.IntOrString{
		Type:   intstr.String,
		StrVal: "zt-udp",
	},
}

var MinikubeZTPort = v1.ServicePort{
	Name:       "zt-udp",
	Protocol:   v1.ProtocolUDP,
	Port:       9993,
	TargetPort: intstr.IntOrString{
		Type:   intstr.String,
		StrVal: "zt-udp",
	},
	NodePort: 9993,
}

var AzureZTPlanetService = v1.Service{
	TypeMeta: metaV1.TypeMeta{
		Kind:       "Service",
		APIVersion: "v1",
	},
	ObjectMeta: metaV1.ObjectMeta{
		Name:      "zt-planet",
		Namespace: "nalej",
		Labels: map[string]string{
			"cluster":                   "management",
			"component": "network-manager",
		},
	},
	Spec: v1.ServiceSpec{
		Ports: []v1.ServicePort{AzureZTPort},
		Selector: map[string]string{
			"cluster":  "management",
			"component":	"network-manager",
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
		Name:      "zt-planet",
		Namespace: "nalej",
		Labels: map[string]string{
			"cluster":	"management",
			"component": "network-manager",
		},
	},
	Spec: v1.ServiceSpec{
		Ports: []v1.ServicePort{MinikubeZTPort},
		Selector: map[string]string{
			"cluster":  "management",
			"component":	"network-manager",
		},
		Type: v1.ServiceTypeNodePort,
	},
}

