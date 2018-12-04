/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package k8s

import (
	"github.com/nalej/derrors"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
	"github.com/rs/zerolog/log"
	"k8s.io/api/core/v1"
	"k8s.io/api/extensions/v1beta1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1 "k8s.io/api/apps/v1"
	batchV1 "k8s.io/api/batch/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type Kubernetes struct {
	entities.GenericSyncCommand
	KubeConfigPath string                `json:"kubeConfigPath"`
	Client         *kubernetes.Clientset `json:"-"`
}

func (k *Kubernetes) Connect() derrors.Error {

	config, err := clientcmd.BuildConfigFromFlags("", k.KubeConfigPath)
	if err != nil {
		log.Error().Err(err).Msg("error building configuration from kubeconfig")
		return derrors.AsError(err, "error building configuration from kubeconfig")
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error().Err(err).Msg("error using configuration to build k8s clientset")
		return derrors.AsError(err, "error using configuration to build k8s clientset")
	}

	k.Client = clientset
	return nil
}

func (k *Kubernetes) existNamespace(name string) (bool, derrors.Error) {
	namespaceClient := k.Client.CoreV1().Namespaces()
	opts := metaV1.ListOptions{}
	list, err := namespaceClient.List(opts)
	if err != nil {
		return false, derrors.AsError(err, "cannot obtain the namespace list")
	}
	found := false
	for _, n := range list.Items {
		log.Debug().Interface("n", n).Msg("A namespace")
		if n.Name == name {
			found = true
			break
		}
	}
	return found, nil
}

func (k *Kubernetes) createNamespace(name string) derrors.Error {
	namespaceClient := k.Client.CoreV1().Namespaces()

	toCreate := v1.Namespace{
		ObjectMeta: metaV1.ObjectMeta{
			Name: name,
		},
	}
	created, err := namespaceClient.Create(&toCreate)
	if err != nil {
		return derrors.AsError(err, "cannot create namespace")
	}
	log.Debug().Interface("created", created).Msg("namespaces has been created")

	return nil
}

func (k *Kubernetes) createNamespacesIfNotExist(name string) derrors.Error {
	found, fErr := k.existNamespace(name)
	if fErr != nil {
		return fErr
	}

	if !found {
		err := k.createNamespace(name)
		if err != nil {
			return err
		}
	} else {
		log.Debug().Str("namespace", name).Msg("namespace already exists")
	}
	return nil
}

func (k *Kubernetes) createService(service *v1.Service) derrors.Error {
	serviceClient := k.Client.CoreV1().Services(service.Namespace)
	log.Debug().Interface("service", service).Msg("unmarshalled")
	created, err := serviceClient.Create(service)
	if err != nil {
		return derrors.AsError(err, "cannot create service")
	}
	log.Debug().Interface("created", created).Msg("new service has been created")
	return nil
}

func (k *Kubernetes) createConfigMap(configMap *v1.ConfigMap) derrors.Error {
	cfClient := k.Client.CoreV1().ConfigMaps(configMap.Namespace)
	log.Debug().Interface("configMap", configMap).Msg("unmarshalled")
	created, err := cfClient.Create(configMap)
	if err != nil {
		return derrors.AsError(err, "cannot create config map")
	}
	log.Debug().Interface("created", created).Msg("new config map has been created")
	return nil
}

func (k *Kubernetes) createIngress(ingress *v1beta1.Ingress) derrors.Error {
	client := k.Client.ExtensionsV1beta1().Ingresses(ingress.Namespace)
	log.Debug().Interface("ingress", ingress).Msg("unmarshalled")
	created, err := client.Create(ingress)
	if err != nil {
		return derrors.AsError(err, "cannot create ingress")
	}
	log.Debug().Interface("created", created).Msg("new ingress set")
	return nil
}

func (k *Kubernetes) createDeployment(deployment *appsv1.Deployment) derrors.Error {
	deploymentClient := k.Client.AppsV1().Deployments(deployment.Namespace)
	log.Debug().Interface("deployment", deployment).Msg("unmarshalled")
	created, err := deploymentClient.Create(deployment)
	if err != nil {
		return derrors.AsError(err, "cannot create deployment")
	}
	log.Debug().Interface("created", created).Msg("new deployment has been created")
	return nil
}

func (k *Kubernetes) createJob(job *batchV1.Job) derrors.Error {
	client := k.Client.BatchV1().Jobs(job.Namespace)
	log.Debug().Interface("job", job).Msg("unmarshalled")
	created, err := client.Create(job)
	if err != nil {
		return derrors.AsError(err, "cannot create job")
	}
	log.Debug().Interface("created", created).Msg("new job has been created")
	return nil
}

func (k *Kubernetes) createServiceAccount(serviceAccount *v1.ServiceAccount) derrors.Error {
	client := k.Client.CoreV1().ServiceAccounts(serviceAccount.Namespace)
	log.Debug().Interface("serviceAccount", serviceAccount).Msg("unmarshalled")
	created, err := client.Create(serviceAccount)
	if err != nil {
		return derrors.AsError(err, "cannot create service account")
	}
	log.Debug().Interface("created", created).Msg("new service account has been created")
	return nil
}

func (k *Kubernetes) createRole(role *rbacv1.Role) derrors.Error {
	client := k.Client.RbacV1().Roles(role.Namespace)
	log.Debug().Interface("role", role).Msg("unmarshalled")
	created, err := client.Create(role)
	if err != nil {
		return derrors.AsError(err, "cannot create role")
	}
	log.Debug().Interface("created", created).Msg("new role has been created")
	return nil
}

func (k *Kubernetes) createClusterRole(clusterRole *rbacv1.ClusterRole) derrors.Error {
	client := k.Client.RbacV1().ClusterRoles()
	log.Debug().Interface("clusterRole", clusterRole).Msg("unmarshalled")
	created, err := client.Create(clusterRole)
	if err != nil {
		return derrors.AsError(err, "cannot create cluster role")
	}
	log.Debug().Interface("created", created).Msg("new cluster role has been created")
	return nil
}

func (k *Kubernetes) createClusterRoleBinding(clusterRoleBinding *rbacv1.ClusterRoleBinding) derrors.Error {
	client := k.Client.RbacV1().ClusterRoleBindings()
	log.Debug().Interface("clusterRoleBinding", clusterRoleBinding).Msg("unmarshalled")
	created, err := client.Create(clusterRoleBinding)
	if err != nil {
		return derrors.AsError(err, "cannot create cluster role binding")
	}
	log.Debug().Interface("created", created).Msg("new cluster role binding has been created")
	return nil
}

func (k *Kubernetes) createRoleBinding(roleBinding *rbacv1.RoleBinding) derrors.Error {
	client := k.Client.RbacV1().RoleBindings(roleBinding.Namespace)
	log.Debug().Interface("roleBinding", roleBinding).Msg("unmarshalled")
	created, err := client.Create(roleBinding)
	if err != nil {
		return derrors.AsError(err, "cannot create role binding")
	}
	log.Debug().Interface("created", created).Msg("new role binding has been created")
	return nil
}

