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

package ingress

import (
	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

var AzureDNSPort = v1.ServicePort{
	Name:     "dns-udp",
	Protocol: v1.ProtocolUDP,
	Port:     53,
	TargetPort: intstr.IntOrString{
		Type:   intstr.String,
		StrVal: "dns-udp",
	},
}

var MinikubeDNSUDPPort = v1.ServicePort{
	Name:     "dns-udp",
	Protocol: v1.ProtocolUDP,
	Port:     53,
	TargetPort: intstr.IntOrString{
		Type:   intstr.String,
		StrVal: "dns-udp",
	},
	NodePort: 53,
}

var MinikubeDNSTCPPort = v1.ServicePort{
	Name:     "dns-tcp",
	Protocol: v1.ProtocolTCP,
	Port:     53,
	TargetPort: intstr.IntOrString{
		StrVal: "dns-tcp",
		Type:   intstr.String,
	},
	NodePort: 53,
}

var MinikubeDNSUIPort = v1.ServicePort{
	Name:     "consul-gui",
	Protocol: v1.ProtocolTCP,
	Port:     8500,
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
			"cluster":   "management",
			"component": "dns-server",
			"release":   "dns-server",
			"app":       "consul",
		},
	},
	Spec: v1.ServiceSpec{
		Ports: []v1.ServicePort{AzureDNSPort},
		Selector: map[string]string{
			"app":     "consul",
			"hasDNS":  "true",
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
			"cluster":   "management",
			"component": "dns-server",
			"release":   "dns-server",
			"app":       "consul",
		},
	},
	Spec: v1.ServiceSpec{
		Ports: []v1.ServicePort{MinikubeDNSUDPPort, MinikubeDNSTCPPort, MinikubeDNSUIPort},
		Selector: map[string]string{
			"app":     "consul",
			"hasDNS":  "true",
			"release": "dns-server",
		},
		Type: v1.ServiceTypeNodePort,
	},
}

var AzureExtDnsService = v1.Service{
	TypeMeta: metaV1.TypeMeta{
		Kind:       "Service",
		APIVersion: "v1",
	},
	ObjectMeta: metaV1.ObjectMeta{
		Name:      "coredns",
		Namespace: "nalej",
		Labels: map[string]string{
			"cluster":   "management",
			"component": "external-dns",
		},
	},
	Spec: v1.ServiceSpec{
		Ports: []v1.ServicePort{AzureDNSPort},
		Selector: map[string]string{
			"cluster":   "management",
			"component": "external-dns",
		},
		Type: v1.ServiceTypeLoadBalancer,
	},
}

var MinikubeExtDnsService = v1.Service{
	TypeMeta: metaV1.TypeMeta{
		Kind:       "Service",
		APIVersion: "v1",
	},
	ObjectMeta: metaV1.ObjectMeta{
		Name:      "coredns",
		Namespace: "nalej",
		Labels: map[string]string{
			"cluster":   "management",
			"component": "external-dns",
		},
	},
	Spec: v1.ServiceSpec{
		Ports: []v1.ServicePort{MinikubeDNSUDPPort, MinikubeDNSTCPPort, MinikubeDNSUIPort},
		Selector: map[string]string{
			"app":     "consul",
			"hasDNS":  "true",
			"release": "dns-server",
		},
		Type: v1.ServiceTypeNodePort,
	},
}
