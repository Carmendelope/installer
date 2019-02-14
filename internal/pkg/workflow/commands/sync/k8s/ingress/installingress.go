/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

// References
// https://github.com/kubernetes/minikube/blob/master/deploy/addons/ingress/ingress-dp.yaml

package ingress

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nalej/derrors"
	"github.com/nalej/installer/internal/pkg/errors"
	"github.com/nalej/installer/internal/pkg/workflow/commands/sync/k8s"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
	"github.com/rs/zerolog/log"
	"k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type InstallTargetType int32

const (
	MinikubeCluster InstallTargetType = iota + 1
	AzureCluster
	Unknown
)

// TODO Refactor using the new parameter to define the target platform instead of detecting it.
type InstallIngress struct {
	k8s.Kubernetes
	ManagementPublicHost string `json:"management_public_host"`
	OnManagementCluster  bool   `json:"on_management_cluster"`
	UseStaticIP          bool   `json:"use_static_ip"`
	StaticIPAddress      string `json:"static_ip_address"`
}

func NewInstallIngress(kubeConfigPath string, managementPublicHost string, useStaticIP bool, staticIPAddress string) *InstallIngress {
	return &InstallIngress{
		Kubernetes: k8s.Kubernetes{
			GenericSyncCommand: *entities.NewSyncCommand(entities.InstallIngress),
			KubeConfigPath:     kubeConfigPath,
		},
		ManagementPublicHost: managementPublicHost,
		UseStaticIP:          useStaticIP,
		StaticIPAddress:      staticIPAddress,
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

func (ii *InstallIngress) getAppClusterIngressRules() []*v1beta1.Ingress {
	appClusterAPI := AppClusterAPIIngressRules
	appClusterAPI.Spec.TLS[0].Hosts[0] = fmt.Sprintf("appcluster.%s", ii.ManagementPublicHost)
	appClusterAPI.Spec.Rules[0].Host = fmt.Sprintf("appcluster.%s", ii.ManagementPublicHost)
	return []*v1beta1.Ingress{&appClusterAPI}
}

func (ii *InstallIngress) getIngressRules() []*v1beta1.Ingress {
	ingress := IngressRules
	ingress.Spec.TLS[0].Hosts[0] = fmt.Sprintf("web.%s", ii.ManagementPublicHost)
	ingress.Spec.Rules[0].Host = fmt.Sprintf("web.%s", ii.ManagementPublicHost)

	login := LoginAPIIngressRules
	login.Spec.TLS[0].Hosts[0] = fmt.Sprintf("login.%s", ii.ManagementPublicHost)
	login.Spec.Rules[0].Host = fmt.Sprintf("login.%s", ii.ManagementPublicHost)

	signup := SignupAPIIngressRules
	signup.Spec.TLS[0].Hosts[0] = fmt.Sprintf("signup.%s", ii.ManagementPublicHost)
	signup.Spec.Rules[0].Host = fmt.Sprintf("signup.%s", ii.ManagementPublicHost)

	api := PublicAPIIngressRules
	api.Spec.TLS[0].Hosts[0] = fmt.Sprintf("api.%s", ii.ManagementPublicHost)
	api.Spec.Rules[0].Host = fmt.Sprintf("api.%s", ii.ManagementPublicHost)

	cluster := ClusterAPIIngressRules
	cluster.Spec.TLS[0].Hosts[0] = fmt.Sprintf("cluster.%s", ii.ManagementPublicHost)
	cluster.Spec.Rules[0].Host = fmt.Sprintf("cluster.%s", ii.ManagementPublicHost)

	device := DeviceAPIIngressRules
	device.Spec.TLS[0].Hosts[0] = fmt.Sprintf("device.%s", ii.ManagementPublicHost)
	device.Spec.Rules[0].Host = fmt.Sprintf("device.%s", ii.ManagementPublicHost)

	deviceLogin := DeviceLoginAPIIngressRules
	deviceLogin.Spec.TLS[0].Hosts[0] = fmt.Sprintf("device-login.%s", ii.ManagementPublicHost)
	deviceLogin.Spec.Rules[0].Host = fmt.Sprintf("device-login.%s", ii.ManagementPublicHost)

	return []*v1beta1.Ingress{
		&ingress, &login, &signup, &api, &cluster, &device, &deviceLogin,
	}

}

func (ii *InstallIngress) getService(installType InstallTargetType) (*v1.Service, *v1.Service) {
	if installType == MinikubeCluster {
		return &MinikubeService, &MinikubeServiceDefaultBackend
	}

	genericService := CloudGenericService
	if ii.UseStaticIP {
		genericService.Spec.LoadBalancerIP = ii.StaticIPAddress
	}

	return &genericService, &CloudGenericServiceDefaultBackend
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
	if ii.OnManagementCluster {
		return ii.triggerManagementInstall(installType)
	}
	return ii.triggerAppClusterInstall(installType)
}

// Trigger the installation of the ingress infrastructure for the application clusters.
func (ii *InstallIngress) triggerAppClusterInstall(installType InstallTargetType) derrors.Error {
	err := ii.CreateNamespacesIfNotExist("nalej")
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating nalej namespace")
		return err
	}

	log.Debug().Msg("Installing ingress service account")
	err = ii.CreateServiceAccount(&IngressServiceAccount)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating ingress service account")
		return err
	}

	log.Debug().Msg("Installing ingress cluster role")
	err = ii.CreateClusterRole(&IngressClusterRole)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating ingress cluster role")
		return err
	}

	log.Debug().Msg("Installing ingress role")
	err = ii.CreateRole(&IngressRole)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating ingress role")
		return err
	}

	log.Debug().Msg("Installing ingress role binding")
	err = ii.CreateRoleBinding(&IngressRoleBinding)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating ingress role binding")
		return err
	}

	log.Debug().Msg("Installing ingress cluster role binding")
	err = ii.CreateClusterRoleBinding(&IngressClusterRoleBinding)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating ingress cluster role binding")
		return err
	}

	log.Debug().Msg("Installing ingress load balancer configmap")
	err = ii.CreateConfigMap(&IngressLoadBalancerConfigMap)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating ingress load balancer configmap")
		return err
	}

	log.Debug().Msg("Installing ingress TCP configmap")
	err = ii.CreateConfigMap(&IngressTCPServiceConfigMap)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating ingress TCP configmap")
		return err
	}

	log.Debug().Msg("Installing ingress UDP configmap")
	err = ii.CreateConfigMap(&IngressUDPServiceConfigMap)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating ingress UDP configmap")
		return err
	}

	log.Debug().Msg("Installing ingress service")
	ingressBackend, defaultBackend := ii.getService(installType)

	err = ii.CreateService(ingressBackend)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating ingress service")
		return err
	}
	err = ii.CreateService(defaultBackend)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating ingress default service")
		return err
	}

	log.Debug().Msg("Installing app cluster ingress rules")
	for _, ingressToInstall := range ii.getAppClusterIngressRules() {
		err = ii.CreateIngress(ingressToInstall)
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
	err = ii.CreateDeployment(&ingressDeployment)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating ingress deployment")
		return err
	}
	log.Debug().Msg("installing default ingress backend")
	err = ii.CreateDeployment(&IngressDefaultBackend)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating ingress backend")
		return err
	}

	return nil
}

// Trigger the installation of the ingress infrastructure for the management cluster.
func (ii *InstallIngress) triggerManagementInstall(installType InstallTargetType) derrors.Error {

	err := ii.CreateNamespacesIfNotExist("nalej")
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating nalej namespace")
		return err
	}

	log.Debug().Msg("Installing ingress service account")
	err = ii.CreateServiceAccount(&IngressServiceAccount)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating ingress service account")
		return err
	}

	log.Debug().Msg("Installing ingress cluster role")
	err = ii.CreateClusterRole(&IngressClusterRole)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating ingress cluster role")
		return err
	}

	log.Debug().Msg("Installing ingress role")
	err = ii.CreateRole(&IngressRole)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating ingress role")
		return err
	}

	log.Debug().Msg("Installing ingress role binding")
	err = ii.CreateRoleBinding(&IngressRoleBinding)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating ingress role binding")
		return err
	}

	log.Debug().Msg("Installing ingress cluster role binding")
	err = ii.CreateClusterRoleBinding(&IngressClusterRoleBinding)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating ingress cluster role binding")
		return err
	}

	log.Debug().Msg("Installing ingress load balancer configmap")
	err = ii.CreateConfigMap(&IngressLoadBalancerConfigMap)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating ingress load balancer configmap")
		return err
	}

	log.Debug().Msg("Installing ingress TCP configmap")
	err = ii.CreateConfigMap(&IngressTCPServiceConfigMap)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating ingress TCP configmap")
		return err
	}

	log.Debug().Msg("Installing ingress UDP configmap")
	err = ii.CreateConfigMap(&IngressUDPServiceConfigMap)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating ingress UDP configmap")
		return err
	}

	log.Debug().Msg("Installing ingress service")
	ingressBackend, defaultBackend := ii.getService(installType)

	err = ii.CreateService(ingressBackend)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating ingress service")
		return err
	}
	err = ii.CreateService(defaultBackend)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating ingress default service")
		return err
	}
	log.Debug().Msg("Installing ingress rules")
	for _, ingressToInstall := range ii.getIngressRules() {
		err = ii.CreateIngress(ingressToInstall)
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
	err = ii.CreateDeployment(&ingressDeployment)
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("error creating ingress deployment")
		return err
	}
	log.Debug().Msg("installing default ingress backend")
	err = ii.CreateDeployment(&IngressDefaultBackend)
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
	return fmt.Sprintf("SYNC InstallIngress on Management: %t", ii.OnManagementCluster)
}

func (ii *InstallIngress) PrettyPrint(indentation int) string {
	msg := strings.Repeat(" ", indentation) + "  Ingresses:"
	var ingresses []*v1beta1.Ingress
	if ii.OnManagementCluster {
		ingresses = ii.getIngressRules()
	} else {
		ingresses = ii.getAppClusterIngressRules()
	}
	for _, ing := range ingresses {
		msg = msg + fmt.Sprintf("\n%s    %s: %s", strings.Repeat(" ", indentation), ing.Name, ing.Spec.Rules[0].Host)
	}
	return strings.Repeat(" ", indentation) + ii.String() + "\n" + msg
}

func (ii *InstallIngress) UserString() string {
	return fmt.Sprintf("Installing ingress on Kubernetes")
}
