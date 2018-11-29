/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package k8s

import (
	"encoding/json"
	"fmt"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"github.com/nalej/derrors"
	"github.com/nalej/installer/internal/pkg/errors"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
	"github.com/rs/zerolog/log"
	"k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	TargetPort: intstr.IntOrString{StrVal:"http"},
}
var HttpsPort = v1.ServicePort{
	Name:       "https",
	Protocol:   v1.ProtocolTCP,
	Port:       443,
	TargetPort: intstr.IntOrString{StrVal:"https"},
}

// CloudGenericService ingress config based on https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/provider/cloud-generic.yaml
var CloudGenericService = v1.Service{
	TypeMeta:   metaV1.TypeMeta{
		Kind: "Service",
		APIVersion: "v1",
	},
	ObjectMeta: metaV1.ObjectMeta{
		Name:                       "ingress-nginx",
		Namespace:                  "nalej",
		Labels: map[string]string{
			"app.kubernetes.io/name":"ingress-nginx",
			"app.kubernetes.io/part-of":"ingress-nginx",
		},
	},
	Spec:       v1.ServiceSpec{
		Ports:                    []v1.ServicePort{HttpPort, HttpsPort},
		Selector: map[string]string{
			"app.kubernetes.io/name":"ingress-nginx",
			"app.kubernetes.io/part-of":"ingress-nginx",
		},
		Type:                     v1.ServiceTypeLoadBalancer,
		ExternalTrafficPolicy:    v1.ServiceExternalTrafficPolicyTypeLocal,
	},
}

// IngressRulesPaths contains the rules for the ingress redirection.
var IngressRulesPaths = &extensions.HTTPIngressRuleValue{
	Paths: []extensions.HTTPIngressPath{
		extensions.HTTPIngressPath{
			Path:    "/",
			Backend: extensions.IngressBackend{
				ServiceName: "web",
				ServicePort: intstr.IntOrString{IntVal: 80},
			},
		},
		extensions.HTTPIngressPath{
			Path:    "/v1/login",
			Backend: extensions.IngressBackend{
				ServiceName: "login-api",
				ServicePort: intstr.IntOrString{IntVal: 8443},
			},
		},
		extensions.HTTPIngressPath{
			Path:    "/v1",
			Backend: extensions.IngressBackend{
				ServiceName: "public-api",
				ServicePort: intstr.IntOrString{IntVal: 8082},
			},
		},
	},
}

var IngressRules = extensions.Ingress{
	TypeMeta:   metaV1.TypeMeta{
		Kind: "Ingress",
		APIVersion: "extensions/v1beta1",
	},
	ObjectMeta: metaV1.ObjectMeta{
		Name:                       "ingress-nginx",
		Namespace:                  "nalej",
		Labels: map[string]string{
			"cluster":"management",
			"component":"ingress-nginx",
		},
		Annotations: map[string]string{
			"nginx.ingress.kubernetes.io/rewrite-target":"/",
		},
	},
	Spec:       extensions.IngressSpec{
		TLS:     []extensions.IngressTLS{
			extensions.IngressTLS{
				Hosts:      []string{"MANAGEMENT_HOST"},
				SecretName: "ingress-tls",
			},
		},
		Rules:   []extensions.IngressRule{
			extensions.IngressRule{
				Host:             "MANAGEMENT_HOST",
				IngressRuleValue: extensions.IngressRuleValue{
					HTTP: IngressRulesPaths,
				},
			},
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

func NewInstallIngress(kubeConfigPath string) *InstallIngress {
	return &InstallIngress{
		Kubernetes: Kubernetes{
			GenericSyncCommand: *entities.NewSyncCommand(entities.InstallIngress),
			KubeConfigPath:     kubeConfigPath,
		},
	}
}

func NewInstallIngressJSON(raw []byte) (*entities.Command, derrors.Error) {
	ccc := &InstallIngress{}
	if err := json.Unmarshal(raw, &ccc); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	ccc.CommandID = entities.GenerateCommandID(ccc.Name())
	var r entities.Command = ccc
	return &r, nil
}

func (ii * InstallIngress) getIngressRules() *extensions.Ingress {
	toReturn := IngressRules
	toReturn.Spec.TLS[0].Hosts[0] = ii.ManagementPublicHost
	toReturn.Spec.Rules[0].Host = ii.ManagementPublicHost
	return &toReturn
}

// GetExistingIngressOnNamespace checks if an ingress exists on a given namespace.
func (ii * InstallIngress) GetExistingIngressOnNamespace(namespace string) (*v1beta1.Ingress, derrors.Error){
	client := ii.Client.ExtensionsV1beta1().Ingresses(namespace)
	opts := metaV1.ListOptions{}
	ingresses, err := client.List(opts)
	if err != nil{
		return nil, derrors.NewInternalError("cannot retrieve ingresses", err)
	}
	if len(ingresses.Items) > 0 {
		return &ingresses.Items[0], nil
	}
	return nil, nil
}

// GetExistingIngress retrieves an ingress if it exists on the system.
func (ii *InstallIngress) GetExistingIngress() (*v1beta1.Ingress, derrors.Error){
	opts := metaV1.ListOptions{}
	namespaces, err := ii.Client.CoreV1().Namespaces().List(opts)
	if err != nil{
		return nil, derrors.NewInternalError("cannot retrieve namespaces", err)
	}
	for _, ns := range namespaces.Items{
		found, err := ii.GetExistingIngressOnNamespace(ns.Name)
		if err != nil{
			return nil, err
		}
		if found != nil{
			return found, nil
		}
	}
	return nil, nil
}

func (ii * InstallIngress) triggerInstall(_ InstallTargetType) derrors.Error{
	err := ii.createService(&CloudGenericService)
	if err != nil{
		return err
	}

	return nil
}

func (ii * InstallIngress) DetectInstallType(nodes *v1.NodeList) InstallTargetType {
	// Check images for minikube
	for _, n := range nodes.Items{
		log.Debug().Interface("node", n).Msg("Analyzing node to detect install")
		for k, _ := range n.Labels {
			if strings.Contains(k, "kubernetes.azure.com"){
				return AzureCluster
			}
		}
		for _, img := range n.Status.Images{
			if strings.Contains(img.Names[0], "k8s-minikube"){
				return MinikubeCluster
			}
		}
	}
	return Unknown
}

func (ii * InstallIngress) InstallIngress() derrors.Error{
	// Detect the type of target install
	opts := metaV1.ListOptions{}
	client := ii.Client.CoreV1().Nodes()
	nodes, err := client.List(opts)
	if err != nil{
		return derrors.AsError(err, "cannot obtain server nodes")
	}
	detected := ii.DetectInstallType(nodes)
	if detected == MinikubeCluster {
		log.Debug().Msg("Installing ingress in a minikube cluster")
	}
	if detected == AzureCluster {
		log.Debug().Msg("Installing ingress in an Azure cluster")
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
	if err != nil{
		return nil, err
	}
	if existingIngress != nil{
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