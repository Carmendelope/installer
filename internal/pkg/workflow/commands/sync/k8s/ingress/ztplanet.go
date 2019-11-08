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

var AzureZTPort = v1.ServicePort{
	Name:     "zt-udp",
	Protocol: v1.ProtocolUDP,
	Port:     9993,
	TargetPort: intstr.IntOrString{
		Type:   intstr.Int,
		IntVal: 9993,
	},
}

var MinikubeZTPort = v1.ServicePort{
	Name:     "zt-udp",
	Protocol: v1.ProtocolUDP,
	Port:     9993,
	TargetPort: intstr.IntOrString{
		Type:   intstr.Int,
		IntVal: 9993,
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
			"cluster":   "management",
			"component": "network-manager",
		},
	},
	Spec: v1.ServiceSpec{
		Ports: []v1.ServicePort{AzureZTPort},
		Selector: map[string]string{
			"cluster":   "management",
			"component": "network-manager",
		},
		Type: v1.ServiceTypeLoadBalancer,
	},
}

var MinikubeZTPlanetService = v1.Service{
	TypeMeta: metaV1.TypeMeta{
		Kind:       "Service",
		APIVersion: "v1",
	},
	ObjectMeta: metaV1.ObjectMeta{
		Name:      "zt-planet",
		Namespace: "nalej",
		Labels: map[string]string{
			"cluster":   "management",
			"component": "network-manager",
		},
	},
	Spec: v1.ServiceSpec{
		Ports: []v1.ServicePort{MinikubeZTPort},
		Selector: map[string]string{
			"cluster":   "management",
			"component": "network-manager",
		},
		Type: v1.ServiceTypeNodePort,
	},
}
