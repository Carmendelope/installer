/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package ingress

import (
	"k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/util/intstr"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1 "k8s.io/api/apps/v1"
	rbacv1 "k8s.io/api/rbac/v1"
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

var CloudGenericHttpPort = v1.ServicePort{
	Name:       "http",
	Protocol:   v1.ProtocolTCP,
	Port:       80,
	TargetPort: intstr.IntOrString{IntVal: 80},
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

var CloudGenericHttpsPort = v1.ServicePort{
	Name:       "https",
	Protocol:   v1.ProtocolTCP,
	Port:       443,
	TargetPort: intstr.IntOrString{IntVal: 443},
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
		Name:      "nginx-ingress-controller",
		Namespace: "kube-system",
		Labels: map[string]string{
			"cluster":                   "management",
			"app.kubernetes.io/name":    "nginx-ingress-controller",
			"app.kubernetes.io/part-of": "kube-system",
			"addonmanager.kubernetes.io/mode":"Reconcile",
		},
	},
	Spec: v1.ServiceSpec{
		Ports: []v1.ServicePort{CloudGenericHttpPort, CloudGenericHttpsPort},
		Selector: map[string]string{
			"app.kubernetes.io/name":    "nginx-ingress-controller",
		},
		Type: v1.ServiceTypeLoadBalancer,
		ExternalTrafficPolicy: v1.ServiceExternalTrafficPolicyTypeCluster,
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

var CloudGenericServiceDefaultBackend = v1.Service{
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
			"certmanager.k8s.io/acme-challenge-type": "http01",
			"certmanager.k8s.io/issuer": "letsencrypt-staging",
			// "nginx.ingress.kubernetes.io/rewrite-target": "/", Not required as we do not need to rewrite the paths.
		},
	},
	Spec: v1beta1.IngressSpec{
		TLS: []v1beta1.IngressTLS{
			v1beta1.IngressTLS{
				Hosts:      []string{"web.MANAGEMENT_HOST"},
				SecretName: "ingress-tls",
			},
		},
		Rules: []v1beta1.IngressRule{
			v1beta1.IngressRule{
				Host: "web.MANAGEMENT_HOST",
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
			"certmanager.k8s.io/acme-challenge-type": "http01",
			"certmanager.k8s.io/issuer": "letsencrypt-staging",
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
			"certmanager.k8s.io/acme-challenge-type": "http01",
			"certmanager.k8s.io/issuer": "letsencrypt-staging",
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
			"certmanager.k8s.io/acme-challenge-type": "http01",
			"certmanager.k8s.io/issuer": "letsencrypt-staging",
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

var ClusterAPIIngressRules = v1beta1.Ingress{
	TypeMeta: metaV1.TypeMeta{
		Kind:       "Ingress",
		APIVersion: "extensions/v1beta1",
	},
	ObjectMeta: metaV1.ObjectMeta{
		Name:      "cluster-api-ingress",
		Namespace: "nalej",
		Labels: map[string]string{
			"cluster":   "management",
			"component": "ingress-nginx",
		},
		Annotations: map[string]string{
			"kubernetes.io/ingress.class": "nginx",
			"certmanager.k8s.io/acme-challenge-type": "http01",
			"certmanager.k8s.io/issuer": "letsencrypt-staging",
			"nginx.ingress.kubernetes.io/ssl-redirect": "true",
			"nginx.ingress.kubernetes.io/backend-protocol": "GRPC",
		},
	},
	Spec: v1beta1.IngressSpec{
		TLS: []v1beta1.IngressTLS{
			v1beta1.IngressTLS{
				Hosts:      []string{"cluster.MANAGEMENT_HOST"},
				SecretName: "signup-server-tls",
			},
		},
		Rules: []v1beta1.IngressRule{
			v1beta1.IngressRule{
				Host: "cluster.MANAGEMENT_HOST",
				IngressRuleValue: v1beta1.IngressRuleValue{
					HTTP: &v1beta1.HTTPIngressRuleValue{
						Paths: []v1beta1.HTTPIngressPath{
							v1beta1.HTTPIngressPath{
								Backend: v1beta1.IngressBackend{
									ServiceName: "cluster-api",
									ServicePort: intstr.IntOrString{IntVal: 8280},
								},
							},
						},
					},
				},
			},
		},
	},
}

var AppClusterAPIIngressRules = v1beta1.Ingress{
	TypeMeta: metaV1.TypeMeta{
		Kind:       "Ingress",
		APIVersion: "extensions/v1beta1",
	},
	ObjectMeta: metaV1.ObjectMeta{
		Name:      "app-cluster-api-ingress",
		Namespace: "nalej",
		Labels: map[string]string{
			"cluster":   "management",
			"component": "ingress-nginx",
		},
		Annotations: map[string]string{
			"kubernetes.io/ingress.class": "nginx",
			"certmanager.k8s.io/acme-challenge-type": "http01",
			"certmanager.k8s.io/issuer": "letsencrypt-staging",
			"nginx.ingress.kubernetes.io/ssl-redirect": "true",
			"nginx.ingress.kubernetes.io/backend-protocol": "GRPC",
		},
	},
	Spec: v1beta1.IngressSpec{
		TLS: []v1beta1.IngressTLS{
			v1beta1.IngressTLS{
				Hosts:      []string{"appcluster.MANAGEMENT_HOST"},
				SecretName: "app-cluster-api-tls",
			},
		},
		Rules: []v1beta1.IngressRule{
			v1beta1.IngressRule{
				Host: "appcluster.MANAGEMENT_HOST",
				IngressRuleValue: v1beta1.IngressRuleValue{
					HTTP: &v1beta1.HTTPIngressRuleValue{
						Paths: []v1beta1.HTTPIngressPath{
							v1beta1.HTTPIngressPath{
								Backend: v1beta1.IngressBackend{
									ServiceName: "app-cluster-api",
									ServicePort: intstr.IntOrString{IntVal: 8281},
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

// The Ingress definition structure contains all the elements required to setup an ingress.
// TODO Refactor into this
type IngressDefinition struct {
	Services []*v1.Service
	Ingresses []*v1beta1.Ingress
	Deployments []*appsv1.Deployment
	ServiceAccounts [] *v1.ServiceAccount
	ClusterRoles [] *rbacv1.ClusterRole
	Roles [] *rbacv1.Role
	RoleBindings []*rbacv1.RoleBinding
	ClusterRoleBindings []*rbacv1.ClusterRoleBinding
	ConfigMaps []*v1.ConfigMap
}
