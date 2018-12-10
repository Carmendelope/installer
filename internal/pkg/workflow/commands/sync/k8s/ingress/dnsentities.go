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
	Name:       "dns",
	Protocol:   v1.ProtocolUDP,
	Port:       53,
	TargetPort: intstr.IntOrString{IntVal: 53},
}

var MinikubeDNSPort = v1.ServicePort{
	Name:       "dns",
	Protocol:   v1.ProtocolUDP,
	Port:       53,
	TargetPort: intstr.IntOrString{IntVal: 53},
}

var AzureConsulService = v1.Service{
	TypeMeta: metaV1.TypeMeta{
		Kind:       "Service",
		APIVersion: "v1",
	},
	ObjectMeta: metaV1.ObjectMeta{
		Name:      "dns-consul-server",
		Namespace: "kube-system",
		Labels: map[string]string{
			"cluster":                   "management",
			"app.kubernetes.io/name":    "dns-consul-server",
			"app.kubernetes.io/part-of": "kube-system",
			"addonmanager.kubernetes.io/mode":"Reconcile",
		},
	},
	Spec: v1.ServiceSpec{
		Ports: []v1.ServicePort{AzureDNSPort},
		Selector: map[string]string{
			"app":    "consul",
			"component":"dns-server",
			"release": "dns-server",
		},
		Type: v1.ServiceTypeLoadBalancer,
		ExternalTrafficPolicy: v1.ServiceExternalTrafficPolicyTypeCluster,
	},
}

var MinikubeConsulService = v1.Service{
	TypeMeta: metaV1.TypeMeta{
		Kind:       "Service",
		APIVersion: "v1",
	},
	ObjectMeta: metaV1.ObjectMeta{
		Name:      "nginx-ingress-controller",
		Namespace: "kube-system",
		Labels: map[string]string{
			"cluster":                   "management",
			"app.kubernetes.io/name":    "dns-consul-server",
			"app.kubernetes.io/part-of": "kube-system",
			"addonmanager.kubernetes.io/mode":"Reconcile",
			"kubernetes.io/minikube-addons":"ingress",
			"kubernetes.io/minikube-addons-endpoint":"ingress",
		},
	},
	Spec: v1.ServiceSpec{
		Ports: []v1.ServicePort{MinikubeHttpPort, MinikubeHttpsPort},
		Selector: map[string]string{
			"app.kubernetes.io/name":    "dns-consul-server",
		},
		Type: v1.ServiceTypeNodePort,
		ExternalTrafficPolicy: v1.ServiceExternalTrafficPolicyTypeCluster,
	},
}
