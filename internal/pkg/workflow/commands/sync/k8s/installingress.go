/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

// References
// https://github.com/kubernetes/minikube/blob/master/deploy/addons/ingress/ingress-dp.yaml

package k8s

import (
	"encoding/json"
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/installer/internal/pkg/errors"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
	"github.com/rs/zerolog/log"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/resource"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"strings"
)

/*
kind: Service
apiVersion: v1
metadata:
  name: ingress-nginx
  namespace: ingress-nginx
  labels:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
spec:
  externalTrafficPolicy: Local
  type: LoadBalancer
  selector:
    app.kubernetes.io/name: ingress-nginx
    app.kubernetes.io/part-of: ingress-nginx
  ports:
    - name: http
      port: 80
      targetPort: http
    - name: https
      port: 443
      targetPort: https
*/
var HttpPort = v1.ServicePort{
	Name:       "http",
	Protocol:   v1.ProtocolTCP,
	Port:       80,
	TargetPort: intstr.IntOrString{StrVal: "http"},
}
var MinikubeHttpPort = v1.ServicePort{
	Name:       "http",
	Protocol:   v1.ProtocolTCP,
	Port:       8080,
	TargetPort: intstr.IntOrString{IntVal: 80},
	NodePort:   80,
}
var HttpsPort = v1.ServicePort{
	Name:       "https",
	Protocol:   v1.ProtocolTCP,
	Port:       443,
	TargetPort: intstr.IntOrString{StrVal: "https"},
}
var MinikubeHttpsPort = v1.ServicePort{
	Name:       "https",
	Protocol:   v1.ProtocolTCP,
	Port:       9443,
	TargetPort: intstr.IntOrString{IntVal: 443},
	NodePort: 443,
}

// CloudGenericService ingress config based on https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/provider/cloud-generic.yaml
var CloudGenericService = v1.Service{
	TypeMeta: metaV1.TypeMeta{
		Kind:       "Service",
		APIVersion: "v1",
	},
	ObjectMeta: metaV1.ObjectMeta{
		Name:      "default-http-backend",
		Namespace: "kube-system",
		Labels: map[string]string{
			"cluster":                   "management",
			"app.kubernetes.io/name":    "default-http-backend",
			"app.kubernetes.io/part-of": "kube-system",
			"addonmanager.kubernetes.io/mode":"Reconcile",
		},
	},
	Spec: v1.ServiceSpec{
		Ports: []v1.ServicePort{HttpPort, HttpsPort},
		Selector: map[string]string{
			"app.kubernetes.io/name":    "default-http-backend",
		},
		Type: v1.ServiceTypeLoadBalancer,
		ExternalTrafficPolicy: v1.ServiceExternalTrafficPolicyTypeLocal,
	},
}

var MinikubeService = v1.Service{
	TypeMeta: metaV1.TypeMeta{
		Kind:       "Service",
		APIVersion: "v1",
	},
	ObjectMeta: metaV1.ObjectMeta{
		Name:      "nginx-ingress-controller",
		Namespace: "kube-system",
		Labels: map[string]string{
			"cluster":                   "management",
			"app.kubernetes.io/name":    "nginx-ingress-controller",
			"app.kubernetes.io/part-of": "kube-system",
			"addonmanager.kubernetes.io/mode":"Reconcile",
			"kubernetes.io/minikube-addons":"ingress",
			"kubernetes.io/minikube-addons-endpoint":"ingress",
		},
	},
	Spec: v1.ServiceSpec{
		Ports: []v1.ServicePort{MinikubeHttpPort, MinikubeHttpsPort},
		Selector: map[string]string{
			"app.kubernetes.io/name":    "nginx-ingress-controller",
		},
		Type: v1.ServiceTypeNodePort,
		ExternalTrafficPolicy: v1.ServiceExternalTrafficPolicyTypeCluster,
	},
}

var MinikubeServiceDefaultBackend = v1.Service{
	TypeMeta: metaV1.TypeMeta{
		Kind:       "Service",
		APIVersion: "v1",
	},
	ObjectMeta: metaV1.ObjectMeta{
		Name:      "default-http-backend",
		Namespace: "kube-system",
		Labels: map[string]string{
			"cluster":                   "management",
			"app.kubernetes.io/name":    "default-http-backend",
			"app.kubernetes.io/part-of": "kube-system",
			"addonmanager.kubernetes.io/mode":"Reconcile",
			"kubernetes.io/minikube-addons":"ingress",
			"kubernetes.io/minikube-addons-endpoint":"ingress",
		},
	},
	Spec: v1.ServiceSpec{
		Ports: []v1.ServicePort{HttpPort, HttpsPort},
		Selector: map[string]string{
			"app.kubernetes.io/name":    "default-http-backend",
		},
	},
}


// IngressRulesPaths contains the rules for the ingress redirection.
var IngressRulesPaths = &v1beta1.HTTPIngressRuleValue{
	Paths: []v1beta1.HTTPIngressPath{
		v1beta1.HTTPIngressPath{
			Path: "/",
			Backend: v1beta1.IngressBackend{
				ServiceName: "web",
				ServicePort: intstr.IntOrString{IntVal: 80},
			},
		},
		v1beta1.HTTPIngressPath{
			Path: "/v1/login",
			Backend: v1beta1.IngressBackend{
				ServiceName: "login-api",
				ServicePort: intstr.IntOrString{IntVal: 8443},
			},
		},
		v1beta1.HTTPIngressPath{
			Path: "/v1",
			Backend: v1beta1.IngressBackend{
				ServiceName: "public-api",
				ServicePort: intstr.IntOrString{IntVal: 8082},
			},
		},
	},
}

var IngressRules = v1beta1.Ingress{
	TypeMeta: metaV1.TypeMeta{
		Kind:       "Ingress",
		APIVersion: "extensions/v1beta1",
	},
	ObjectMeta: metaV1.ObjectMeta{
		Name:      "ingress-nginx",
		Namespace: "nalej",
		Labels: map[string]string{
			"cluster":   "management",
			"component": "ingress-nginx",
		},
		Annotations: map[string]string{
			"kubernetes.io/ingress.class": "nginx",
			// "nginx.ingress.kubernetes.io/rewrite-target": "/", Not required as we do not need to rewrite the paths.
		},
	},
	Spec: v1beta1.IngressSpec{
		TLS: []v1beta1.IngressTLS{
			v1beta1.IngressTLS{
				Hosts:      []string{"MANAGEMENT_HOST"},
				SecretName: "ingress-tls",
			},
		},
		Rules: []v1beta1.IngressRule{
			v1beta1.IngressRule{
				Host: "MANAGEMENT_HOST",
				IngressRuleValue: v1beta1.IngressRuleValue{
					HTTP: IngressRulesPaths,
				},
			},
		},
	},
}

var SignupAPIIngressRules = v1beta1.Ingress{
	TypeMeta: metaV1.TypeMeta{
		Kind:       "Ingress",
		APIVersion: "extensions/v1beta1",
	},
	ObjectMeta: metaV1.ObjectMeta{
		Name:      "signup-api-ingress",
		Namespace: "nalej",
		Labels: map[string]string{
			"cluster":   "management",
			"component": "ingress-nginx",
		},
		Annotations: map[string]string{
			"kubernetes.io/ingress.class": "nginx",
			//"nginx.ingress.kubernetes.io/ssl-passthrough": "true",
			"nginx.ingress.kubernetes.io/ssl-redirect": "true",
			"nginx.ingress.kubernetes.io/backend-protocol": "GRPC",
		},
	},
	Spec: v1beta1.IngressSpec{
		TLS: []v1beta1.IngressTLS{
			v1beta1.IngressTLS{
				Hosts:      []string{"signup.MANAGEMENT_HOST"},
				SecretName: "signup-server-tls",
			},
		},
		Rules: []v1beta1.IngressRule{
			v1beta1.IngressRule{
				Host: "signup.MANAGEMENT_HOST",
				IngressRuleValue: v1beta1.IngressRuleValue{
					HTTP: &v1beta1.HTTPIngressRuleValue{
						Paths: []v1beta1.HTTPIngressPath{
							v1beta1.HTTPIngressPath{
								Backend: v1beta1.IngressBackend{
									ServiceName: "signup",
									ServicePort: intstr.IntOrString{IntVal: 8180},
								},
							},
						},
					},
				},
			},
		},
	},
}

var LoginAPIIngressRules = v1beta1.Ingress{
	TypeMeta: metaV1.TypeMeta{
		Kind:       "Ingress",
		APIVersion: "extensions/v1beta1",
	},
	ObjectMeta: metaV1.ObjectMeta{
		Name:      "login-api-ingress",
		Namespace: "nalej",
		Labels: map[string]string{
			"cluster":   "management",
			"component": "ingress-nginx",
		},
		Annotations: map[string]string{
			"kubernetes.io/ingress.class": "nginx",
			"nginx.ingress.kubernetes.io/ssl-redirect": "true",
			"nginx.ingress.kubernetes.io/backend-protocol": "GRPC",
		},
	},
	Spec: v1beta1.IngressSpec{
		TLS: []v1beta1.IngressTLS{
			v1beta1.IngressTLS{
				Hosts:      []string{"login.MANAGEMENT_HOST"},
				SecretName: "signup-server-tls",
			},
		},
		Rules: []v1beta1.IngressRule{
			v1beta1.IngressRule{
				Host: "login.MANAGEMENT_HOST",
				IngressRuleValue: v1beta1.IngressRuleValue{
					HTTP: &v1beta1.HTTPIngressRuleValue{
						Paths: []v1beta1.HTTPIngressPath{
							v1beta1.HTTPIngressPath{
								Backend: v1beta1.IngressBackend{
									ServiceName: "login-api",
									ServicePort: intstr.IntOrString{IntVal: 8444},
								},
							},
						},
					},
				},
			},
		},
	},
}

var PublicAPIIngressRules = v1beta1.Ingress{
	TypeMeta: metaV1.TypeMeta{
		Kind:       "Ingress",
		APIVersion: "extensions/v1beta1",
	},
	ObjectMeta: metaV1.ObjectMeta{
		Name:      "public-api-ingress",
		Namespace: "nalej",
		Labels: map[string]string{
			"cluster":   "management",
			"component": "ingress-nginx",
		},
		Annotations: map[string]string{
			"kubernetes.io/ingress.class": "nginx",
			"nginx.ingress.kubernetes.io/ssl-redirect": "true",
			"nginx.ingress.kubernetes.io/backend-protocol": "GRPC",
		},
	},
	Spec: v1beta1.IngressSpec{
		TLS: []v1beta1.IngressTLS{
			v1beta1.IngressTLS{
				Hosts:      []string{"api.MANAGEMENT_HOST"},
				SecretName: "signup-server-tls",
			},
		},
		Rules: []v1beta1.IngressRule{
			v1beta1.IngressRule{
				Host: "api.MANAGEMENT_HOST",
				IngressRuleValue: v1beta1.IngressRuleValue{
					HTTP: &v1beta1.HTTPIngressRuleValue{
						Paths: []v1beta1.HTTPIngressPath{
							v1beta1.HTTPIngressPath{
								Backend: v1beta1.IngressBackend{
									ServiceName: "public-api",
									ServicePort: intstr.IntOrString{IntVal: 8081},
								},
							},
						},
					},
				},
			},
		},
	},
}


// Adapt num replicas to num nodes.
var IngressNumReplicas int32 = 1

// Ingress deployment based on https://github.com/kubernetes/minikube/blob/master/deploy/addons/ingress/ingress-dp.yaml
var IngressDeployment = appsv1.Deployment{
	TypeMeta: metaV1.TypeMeta{
		Kind:       "Deployment",
		APIVersion: "apps/v1",
	},
	ObjectMeta: metaV1.ObjectMeta{
		Name:      "nginx-ingress-controller",
		Namespace: "kube-system",
		Labels: map[string]string{
			"cluster":                         "management",
			"app.kubernetes.io/name":          "nginx-ingress-controller",
			"app.kubernetes.io/part-of":       "kube-system",
			"addonmanager.kubernetes.io/mode": "Reconcile",
		},
	},
	Spec: appsv1.DeploymentSpec{
		Replicas: &IngressNumReplicas,
		Selector: &metaV1.LabelSelector{
			MatchLabels: map[string]string{
				"app.kubernetes.io/name":          "nginx-ingress-controller",
				"app.kubernetes.io/part-of":       "kube-system",
				"addonmanager.kubernetes.io/mode": "Reconcile",
			},
		},
		Template: v1.PodTemplateSpec{
			ObjectMeta: metaV1.ObjectMeta{
				Labels: map[string]string{
					"app.kubernetes.io/name":          "nginx-ingress-controller",
					"app.kubernetes.io/part-of":       "kube-system",
					"addonmanager.kubernetes.io/mode": "Reconcile",
				},
				Annotations: map[string]string{
					"prometheus.io/port":   "10254",
					"prometheus.io/scrape": "true",
				},
			},
			Spec: v1.PodSpec{
				ServiceAccountName: "nginx-ingress",
				Containers: []v1.Container{
					v1.Container{
						Name:  "nginx-ingress-controller",
						Image: "quay.io/kubernetes-ingress-controller/nginx-ingress-controller:0.21.0",
						Args: []string{
							"/nginx-ingress-controller",
							"--default-backend-service=kube-system/default-http-backend",
							"--configmap=kube-system/nginx-load-balancer-conf",
							"--tcp-services-configmap=kube-system/tcp-services",
							"--udp-services-configmap=kube-system/udp-services",
							"--annotations-prefix=nginx.ingress.kubernetes.io",
							//" --enable-ssl-passthrough",
							"--v=4",
						},
						Ports: []v1.ContainerPort{
							v1.ContainerPort{
								Name:          "port80",
								HostPort:      8080,
								ContainerPort: 80,
							},
							v1.ContainerPort{
								Name:          "port443",
								HostPort:      9443,
								ContainerPort: 443,
							},
						},
						Env: []v1.EnvVar{
							v1.EnvVar{
								Name: "POD_NAME",
								ValueFrom: &v1.EnvVarSource{
									FieldRef: &v1.ObjectFieldSelector{
										FieldPath: "metadata.name",
									},
								},
							},
							v1.EnvVar{
								Name: "POD_NAMESPACE",
								ValueFrom: &v1.EnvVarSource{
									FieldRef: &v1.ObjectFieldSelector{
										FieldPath: "metadata.namespace",
									},
								},
							},
						},
						LivenessProbe: &v1.Probe{
							Handler: v1.Handler{
								HTTPGet: &v1.HTTPGetAction{
									Path:   "/healthz",
									Port:   intstr.IntOrString{IntVal: 10254},
									Scheme: "HTTP",
								},
							},
							InitialDelaySeconds: 10,
							TimeoutSeconds:      1,
						},
						ReadinessProbe: &v1.Probe{
							Handler: v1.Handler{
								HTTPGet: &v1.HTTPGetAction{
									Path:   "/healthz",
									Port:   intstr.IntOrString{IntVal: 10254},
									Scheme: "HTTP",
								},
							},
						},
						ImagePullPolicy: "IfNotPresent",
						SecurityContext: &v1.SecurityContext{
							Capabilities: &v1.Capabilities{
								Add:  []v1.Capability{"NET_BIND_SERVICE"},
								Drop: []v1.Capability{"ALL"},
							},
							RunAsUser: &[]int64{33}[0],
						},
					},
				},
				TerminationGracePeriodSeconds: &[]int64{60}[0], // Golang sucks creating literals...
			},
		},
	},
}

// Ingress deployment based on https://github.com/kubernetes/minikube/blob/master/deploy/addons/ingress/ingress-dp.yaml
var IngressDefaultBackend = appsv1.Deployment{
	TypeMeta: metaV1.TypeMeta{
		Kind:       "Deployment",
		APIVersion: "apps/v1",
	},
	ObjectMeta: metaV1.ObjectMeta{
		Name:      "default-http-backend",
		Namespace: "kube-system",
		Labels: map[string]string{
			"cluster":                         "management",
			"app.kubernetes.io/name":          "default-http-backend",
			"app.kubernetes.io/part-of":       "kube-system",
			"addonmanager.kubernetes.io/mode": "Reconcile",
		},
	},
	Spec: appsv1.DeploymentSpec{
		Replicas: &[]int32{1}[0],
		Selector: &metaV1.LabelSelector{
			MatchLabels: map[string]string{
				"app.kubernetes.io/name":          "default-http-backend",
				"addonmanager.kubernetes.io/mode": "Reconcile",
			},
		},
		Template: v1.PodTemplateSpec{
			ObjectMeta: metaV1.ObjectMeta{
				Labels: map[string]string{
					"app.kubernetes.io/name":          "default-http-backend",
					"addonmanager.kubernetes.io/mode": "Reconcile",
				},
			},
			Spec: v1.PodSpec{
				Containers: []v1.Container{
					v1.Container{
						Name:  "default-http-backend",
						Image: "gcr.io/google_containers/defaultbackend:1.4",
						Ports: []v1.ContainerPort{
							v1.ContainerPort{
								Name:          "port8080",
								ContainerPort: 8080,
							},
						},
						LivenessProbe: &v1.Probe{
							Handler: v1.Handler{
								HTTPGet: &v1.HTTPGetAction{
									Path:   "/healthz",
									Port:   intstr.IntOrString{IntVal: 8080},
									Scheme: "HTTP",
								},
							},
							InitialDelaySeconds: 30,
							TimeoutSeconds:      5,
						},
						ImagePullPolicy: "IfNotPresent",
						Resources: v1.ResourceRequirements{
							Limits: map[v1.ResourceName]resource.Quantity{
								"cpu":    resource.MustParse("20m"),
								"memory": resource.MustParse("30Mi"),
							},
							Requests: map[v1.ResourceName]resource.Quantity{
								"cpu":    resource.MustParse("20m"),
								"memory": resource.MustParse("30Mi"),
							},
						},
					},
				},
				TerminationGracePeriodSeconds: &[]int64{60}[0], // Golang sucks creating literals...
			},
		},
	},
}

var IngressServiceAccount = v1.ServiceAccount{
	TypeMeta:                     metaV1.TypeMeta{
		Kind:       "ServiceAccount",
		APIVersion: "v1",
	},
	ObjectMeta:                   metaV1.ObjectMeta{
		Name:                       "nginx-ingress",
		Namespace:                  "kube-system",
		Labels:          map[string]string{
			"cluster":                         "management",
			"addonmanager.kubernetes.io/mode": "Reconcile",
		},
	},
}

var IngressClusterRole = rbacv1.ClusterRole{
	TypeMeta:        metaV1.TypeMeta{
		Kind:       "ClusterRole",
		APIVersion: "v1",
	},
	ObjectMeta:      metaV1.ObjectMeta{
		Name:                       "system:nginx-ingress",
		Namespace:                  "kube-system",
		Labels:          map[string]string{
			"cluster":                         "management",
			"kubernetes.io/bootstrapping":                         "rbac-defaults",
			"addonmanager.kubernetes.io/mode": "Reconcile",
		},
	},
	Rules:           []rbacv1.PolicyRule{
		rbacv1.PolicyRule{
			Verbs:           []string{"list", "watch"},
			APIGroups:       []string{""},
			Resources:       []string{"configmaps", "endpoints", "nodes", "pods", "secrets"},
		},
		rbacv1.PolicyRule{
			Verbs:           []string{"get"},
			APIGroups:       []string{""},
			Resources:       []string{"nodes"},
		},
		rbacv1.PolicyRule{
			Verbs:           []string{"get", "list", "watch"},
			APIGroups:       []string{""},
			Resources:       []string{"services"},
		},
		rbacv1.PolicyRule{
			Verbs:           []string{"get", "list", "watch"},
			APIGroups:       []string{"extensions"},
			Resources:       []string{"ingresses"},
		},
		rbacv1.PolicyRule{
			Verbs:           []string{"create", "patch"},
			APIGroups:       []string{""},
			Resources:       []string{"events"},
		},
		rbacv1.PolicyRule{
			Verbs:           []string{"update"},
			APIGroups:       []string{"extensions"},
			Resources:       []string{"ingresses/status"},
		},
	},
}

var IngressRole = rbacv1.Role{
	TypeMeta:   metaV1.TypeMeta{
		Kind:       "Role",
		APIVersion: "rbac.authorization.k8s.io/v1beta1",
	},
	ObjectMeta: metaV1.ObjectMeta{
		Name:                       "system::nginx-ingress-role",
		Namespace:                  "kube-system",
		Labels:          map[string]string{
			"cluster":                         "management",
			"kubernetes.io/bootstrapping":                         "rbac-defaults",
			"addonmanager.kubernetes.io/mode": "Reconcile",
		},
	},
	Rules:      []rbacv1.PolicyRule{
		rbacv1.PolicyRule{
			Verbs:           []string{"get"},
			APIGroups:       []string{""},
			Resources:       []string{"configmaps", "pods", "secrets", "namespaces"},
		},
		rbacv1.PolicyRule{
			Verbs:           []string{"get", "update"},
			APIGroups:       []string{""},
			Resources:       []string{"ingress-controller-leader-nginx"},
		},
		rbacv1.PolicyRule{
			Verbs:           []string{"create", "get", "update"},
			APIGroups:       []string{""},
			Resources:       []string{"configmaps"},
		},
		rbacv1.PolicyRule{
			Verbs:           []string{"get"},
			APIGroups:       []string{""},
			Resources:       []string{"endpoints"},
		},
	},
}

var IngressRoleBinding = rbacv1.RoleBinding{
	TypeMeta:   metaV1.TypeMeta{
		Kind:       "RoleBinding",
		APIVersion: "rbac.authorization.k8s.io/v1beta1",
	},
	ObjectMeta: metaV1.ObjectMeta{
		Name:                       "system::nginx-ingress-role-binding",
		Namespace:                  "kube-system",
		Labels:          map[string]string{
			"cluster":                         "management",
			"kubernetes.io/bootstrapping":                         "rbac-defaults",
			"addonmanager.kubernetes.io/mode": "EnsureExists",
		},
	},
	Subjects:   []rbacv1.Subject{
		rbacv1.Subject{
			Kind:      "ServiceAccount",
			Name:      "nginx-ingress",
			Namespace: "kube-system",
		},
	},
	RoleRef:    rbacv1.RoleRef{
		APIGroup: "rbac.authorization.k8s.io",
		Kind:     "Role",
		Name:     "system::nginx-ingress-role",
	},
}

var IngressClusterRoleBinding = rbacv1.ClusterRoleBinding{
	TypeMeta:   metaV1.TypeMeta{
		Kind:       "ClusterRoleBinding",
		APIVersion: "rbac.authorization.k8s.io/v1beta1",
	},
	ObjectMeta: metaV1.ObjectMeta{
		Name:                       "system:nginx-ingress",
		Namespace:                  "kube-system",
		Labels:         map[string]string{
			"cluster":                         "management",
			"kubernetes.io/bootstrapping":                         "rbac-defaults",
			"addonmanager.kubernetes.io/mode": "EnsureExists",
		},
	},
	Subjects:   []rbacv1.Subject{
		rbacv1.Subject{
			Kind:      "ServiceAccount",
			Name:      "nginx-ingress",
			Namespace: "kube-system",
		},
	},
	RoleRef:    rbacv1.RoleRef{
		APIGroup: "rbac.authorization.k8s.io",
		Kind:     "ClusterRole",
		Name:     "system:nginx-ingress",
	},
}

var IngressLoadBalancerConfigMap = v1.ConfigMap{
	TypeMeta:   metaV1.TypeMeta{
		Kind:       "ConfigMap",
		APIVersion: "v1",
	},
	ObjectMeta: metaV1.ObjectMeta{
		Name:                       "nginx-load-balancer-conf",
		Namespace:                  "kube-system",
		Labels: map[string]string{
			"cluster":"management",
			"addonmanager.kubernetes.io/mode":"EnsureExists",
		},
	},
	Data: map[string]string{
		"http2":"True",
	},
}

var IngressTCPServiceConfigMap = v1.ConfigMap{
	TypeMeta:   metaV1.TypeMeta{
		Kind:       "ConfigMap",
		APIVersion: "v1",
	},
	ObjectMeta: metaV1.ObjectMeta{
		Name:                       "tcp-services",
		Namespace:                  "kube-system",
		Labels: map[string]string{
			"cluster":"management",
			"addonmanager.kubernetes.io/mode":"EnsureExists",
		},
	},
}

var IngressUDPServiceConfigMap = v1.ConfigMap{
	TypeMeta:   metaV1.TypeMeta{
		Kind:       "ConfigMap",
		APIVersion: "v1",
	},
	ObjectMeta: metaV1.ObjectMeta{
		Name:                       "udp-services",
		Namespace:                  "kube-system",
		Labels: map[string]string{
			"cluster":"management",
			"addonmanager.kubernetes.io/mode":"EnsureExists",
		},
	},
}


type InstallTargetType int32

const (
	MinikubeCluster InstallTargetType = iota + 1
	AzureCluster
	Unknown
)

type InstallIngress struct {
	Kubernetes
	ManagementPublicHost string `json:"management_public_host"`
}

type Ingresses struct {
	HTTPIngress *v1beta1.Ingress
	LoginGRPC *v1beta1.Ingress
	SignupGRPC *v1beta1.Ingress
	PublicAPIGRPC *v1beta1.Ingress
}

func NewInstallIngress(kubeConfigPath string, managementPublicHost string) *InstallIngress {
	return &InstallIngress{
		Kubernetes: Kubernetes{
			GenericSyncCommand: *entities.NewSyncCommand(entities.InstallIngress),
			KubeConfigPath:     kubeConfigPath,
		},
		ManagementPublicHost: managementPublicHost,
	}
}

func NewInstallIngressFromJSON(raw []byte) (*entities.Command, derrors.Error) {
	ccc := &InstallIngress{}
	if err := json.Unmarshal(raw, &ccc); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	ccc.CommandID = entities.GenerateCommandID(ccc.Name())
	var r entities.Command = ccc
	return &r, nil
}

func (ii *InstallIngress) getIngressRules() []*v1beta1.Ingress {
	ingress := IngressRules
	ingress.Spec.TLS[0].Hosts[0] = ii.ManagementPublicHost
	ingress.Spec.Rules[0].Host = ii.ManagementPublicHost

	login := LoginAPIIngressRules
	login.Spec.Rules[0].Host = fmt.Sprintf("login.%s", ii.ManagementPublicHost)

	signup := SignupAPIIngressRules
	signup.Spec.TLS[0].Hosts[0] = fmt.Sprintf("signup.%s", ii.ManagementPublicHost)
	signup.Spec.Rules[0].Host = fmt.Sprintf("signup.%s", ii.ManagementPublicHost)

	api := PublicAPIIngressRules
	api.Spec.TLS[0].Hosts[0] = fmt.Sprintf("api.%s", ii.ManagementPublicHost)
	api.Spec.Rules[0].Host = fmt.Sprintf("api.%s", ii.ManagementPublicHost)

	return []*v1beta1.Ingress{
		&ingress, &login, &signup, &api,
	}

}

func (ii * InstallIngress) getService(installType InstallTargetType) (*v1.Service, *v1.Service) {
	if installType == MinikubeCluster {
		return &MinikubeService, &MinikubeServiceDefaultBackend
	}
	return &CloudGenericService, &MinikubeServiceDefaultBackend
}

// GetExistingIngressOnNamespace checks if an ingress exists on a given namespace.
func (ii *InstallIngress) GetExistingIngressOnNamespace(namespace string) (*v1beta1.Ingress, derrors.Error) {
	client := ii.Client.ExtensionsV1beta1().Ingresses(namespace)
	opts := metaV1.ListOptions{}
	ingresses, err := client.List(opts)
	if err != nil {
		return nil, derrors.NewInternalError("cannot retrieve ingresses", err)
	}
	if len(ingresses.Items) > 0 {
		return &ingresses.Items[0], nil
	}
	return nil, nil
}

// GetExistingIngress retrieves an ingress if it exists on the system.
func (ii *InstallIngress) GetExistingIngress() (*v1beta1.Ingress, derrors.Error) {
	opts := metaV1.ListOptions{}
	namespaces, err := ii.Client.CoreV1().Namespaces().List(opts)
	if err != nil {
		return nil, derrors.NewInternalError("cannot retrieve namespaces", err)
	}
	for _, ns := range namespaces.Items {
		found, err := ii.GetExistingIngressOnNamespace(ns.Name)
		if err != nil {
			return nil, err
		}
		if found != nil {
			return found, nil
		}
	}
	return nil, nil
}

func (ii *InstallIngress) triggerInstall(installType InstallTargetType) derrors.Error {

	err := ii.createNamespacesIfNotExist("nalej")
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating nalej namespace")
		return err
	}

	log.Debug().Msg("Installing ingress service account")
	err = ii.createServiceAccount(&IngressServiceAccount)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating ingress service account")
		return err
	}

	log.Debug().Msg("Installing ingress cluster role")
	err = ii.createClusterRole(&IngressClusterRole)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating ingress cluster role")
		return err
	}

	log.Debug().Msg("Installing ingress role")
	err = ii.createRole(&IngressRole)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating ingress role")
		return err
	}

	log.Debug().Msg("Installing ingress role binding")
	err = ii.createRoleBinding(&IngressRoleBinding)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating ingress role binding")
		return err
	}

	log.Debug().Msg("Installing ingress cluster role binding")
	err = ii.createClusterRoleBinding(&IngressClusterRoleBinding)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating ingress cluster role binding")
		return err
	}

	log.Debug().Msg("Installing ingress load balancer configmap")
	err = ii.createConfigMap(&IngressLoadBalancerConfigMap)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating ingress load balancer configmap")
		return err
	}

	log.Debug().Msg("Installing ingress TCP configmap")
	err = ii.createConfigMap(&IngressTCPServiceConfigMap)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating ingress TCP configmap")
		return err
	}

	log.Debug().Msg("Installing ingress UDP configmap")
	err = ii.createConfigMap(&IngressUDPServiceConfigMap)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating ingress UDP configmap")
		return err
	}

	log.Debug().Msg("Installing ingress service")
	ingressBackend, defaultBackend := ii.getService(installType)

	err = ii.createService(ingressBackend)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating ingress service")
		return err
	}
	err = ii.createService(defaultBackend)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating ingress default service")
		return err
	}
	log.Debug().Msg("Installing ingress rules")
	for _, ingressToInstall := range ii.getIngressRules(){
		err = ii.createIngress(ingressToInstall)
		if err != nil {
			log.Error().Str("trace", err.DebugReport()).Str("name", ingressToInstall.Name).Msg("error creating ingress rules")
			return err
		}
	}

	var ingressDeployment = IngressDeployment

	if installType == MinikubeCluster {
		log.Debug().Msg("Adding extra arguments and ports for Minikube dev install")
		// args - --report-node-internal-ip-address
		args := ingressDeployment.Spec.Template.Spec.Containers[0].Args
		args = append(args, "--report-node-internal-ip-address")
		ingressDeployment.Spec.Template.Spec.Containers[0].Args = args
		statusPort := v1.ContainerPort{
			Name:          "stats", // on /nginx-status
			HostPort:      18080,
			ContainerPort: 18080,
		}
		ports := ingressDeployment.Spec.Template.Spec.Containers[0].Ports
		ports = append(ports, statusPort)
		ingressDeployment.Spec.Template.Spec.Containers[0].Ports = ports
	}

	log.Debug().Msg("installing ingress deployment")
	err = ii.createDeployment(&ingressDeployment)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating ingress deployment")
		return err
	}
	log.Debug().Msg("installing default ingress backend")
	err = ii.createDeployment(&IngressDefaultBackend)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating ingress backend")
		return err
	}

	return nil
}

func (ii *InstallIngress) DetectInstallType(nodes *v1.NodeList) InstallTargetType {
	// Check images for minikube
	for _, n := range nodes.Items {
		log.Debug().Interface("node", n).Msg("Analyzing node to detect install")
		for k, _ := range n.Labels {
			if strings.Contains(k, "kubernetes.azure.com") {
				return AzureCluster
			}
		}
		for _, img := range n.Status.Images {
			if strings.Contains(img.Names[0], "k8s-minikube") {
				return MinikubeCluster
			}
		}
	}
	return Unknown
}

func (ii *InstallIngress) InstallIngress() derrors.Error {
	// Detect the type of target install
	opts := metaV1.ListOptions{}
	client := ii.Client.CoreV1().Nodes()
	nodes, err := client.List(opts)
	if err != nil {
		return derrors.AsError(err, "cannot obtain server nodes")
	}
	detected := ii.DetectInstallType(nodes)
	if detected == MinikubeCluster {
		log.Debug().Msg("Installing ingress in a minikube cluster")
	}
	if detected == AzureCluster {
		log.Debug().Msg("Installing ingress in an Azure cluster")
	}
	if detected == Unknown {
		log.Warn().Msg("Cannot determine cluster type, assuming Minikube")
		detected = MinikubeCluster
	}
	if detected != Unknown {
		return ii.triggerInstall(detected)
	}
	return derrors.NewNotFoundError("cannot determine type of cluster for the ingress service")
}

func (ii *InstallIngress) Run(workflowID string) (*entities.CommandResult, derrors.Error) {
	connectErr := ii.Connect()
	if connectErr != nil {
		return nil, connectErr
	}
	existingIngress, err := ii.GetExistingIngress()
	if err != nil {
		return nil, err
	}
	if existingIngress != nil {
		log.Warn().Interface("ingress", existingIngress).Msg("An ingress has been found")
		return entities.NewSuccessCommand([]byte("[WARN] Ingress has not been installed as it already exists")), nil
	}

	err = ii.InstallIngress()
	if err != nil {
		return entities.NewCommandResult(
			false, "cannot install an ingress", err), nil
	}

	return entities.NewSuccessCommand([]byte("Ingress controller credentials have been created")), nil
}

func (ii *InstallIngress) String() string {
	return fmt.Sprintf("SYNC InstallIngress")
}

func (ii *InstallIngress) PrettyPrint(indentation int) string {
	return strings.Repeat(" ", indentation) + ii.String()
}

func (ii *InstallIngress) UserString() string {
	return fmt.Sprintf("Installing ingress on Kubernetes")
}
