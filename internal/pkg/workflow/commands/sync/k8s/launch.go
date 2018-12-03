/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package k8s

import (
	"encoding/json"
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/installer/internal/pkg/errors"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
	"github.com/rs/zerolog/log"
	"k8s.io/api/extensions/v1beta1"
	"k8s.io/client-go/kubernetes/scheme"
	"path"
	"reflect"

	"io/ioutil"
	appsv1 "k8s.io/api/apps/v1"
	batchV1 "k8s.io/api/batch/v1"
	"k8s.io/api/core/v1"
	policyv1beta1 "k8s.io/api/policy/v1beta1"
	rbacv1 "k8s.io/api/rbac/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type LaunchComponents struct {
	Kubernetes
	Namespaces []string `json:"namespaces"`
	ComponentsDir string `json:"componentsDir"`
}

func NewLaunchComponents(kubeConfigPath string, namespaces []string, componentsDir string) * LaunchComponents {
	return &LaunchComponents{
		Kubernetes:    Kubernetes{
			GenericSyncCommand: *entities.NewSyncCommand(entities.LaunchComponents),
			KubeConfigPath:     kubeConfigPath,
		},
		Namespaces: namespaces,
		ComponentsDir: componentsDir,
	}
}

// NewLaunchComponentsFromJSON creates an LaunchComponents command from a JSON object.
func NewLaunchComponentsFromJSON(raw []byte) (*entities.Command, derrors.Error) {
	lc := &LaunchComponents{}
	if err := json.Unmarshal(raw, &lc); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	lc.CommandID = entities.GenerateCommandID(lc.Name())
	var r entities.Command = lc
	return &r, nil
}

func (lc * LaunchComponents) Run(workflowID string) (*entities.CommandResult, derrors.Error) {

		connectErr := lc.Connect()
		if connectErr != nil {
		    return nil, connectErr
		}
		for _, target := range lc.Namespaces{
			createErr := lc.createNamespace(target)
			if createErr != nil{
				return nil, createErr
			}
		}

		fileInfo, err := ioutil.ReadDir(lc.ComponentsDir)
		if err != nil {
			return nil, derrors.AsError(err, "cannot read components dir")
		}
		numLaunched := 0
		for _, file := range fileInfo {
			if strings.HasSuffix(file.Name(), ".yaml"){
				log.Debug().Str("file", file.Name()).Msg("processing component")
				err := lc.launchComponent(path.Join(lc.ComponentsDir, file.Name()))
				if err != nil {
					return entities.NewCommandResult(false, "cannot launch component", err), nil
				}
				numLaunched++
			}
		}
		msg := fmt.Sprintf("%d components have been launched", numLaunched)
		return entities.NewCommandResult(true, msg, nil), nil
}

func (lc * LaunchComponents) launchComponent(componentPath string) derrors.Error {
	log.Debug().Str("path", componentPath).Msg("launch component")

	raw, err := ioutil.ReadFile(componentPath)
	if err != nil {
		return derrors.AsError(err, "cannot read component file")
	}
	log.Debug().Msg("parsing component")

	decode := scheme.Codecs.UniversalDeserializer().Decode

	obj, _, err := decode([]byte(raw), nil, nil)
	if err != nil {
		fmt.Printf("%#v", err)
	}

	switch o := obj.(type) {
	case *batchV1.Job:
		return lc.createJob(obj.(*batchV1.Job))
	case *appsv1.Deployment:
		return lc.createDeployment(obj.(*appsv1.Deployment))
	case *appsv1.DaemonSet:
		return lc.launchDaemonSet(obj.(*appsv1.DaemonSet))
	case *v1.Service:
		return lc.createService(obj.(*v1.Service))
	case *v1.Secret:
		return lc.launchSecret(obj.(*v1.Secret))
	case *v1.ServiceAccount:
		return lc.launchServiceAccount(obj.(*v1.ServiceAccount))
	case *v1.ConfigMap:
		return lc.createConfigMap(obj.(*v1.ConfigMap))
	case *rbacv1.RoleBinding:
		return lc.launchRoleBinding(obj.(*rbacv1.RoleBinding))
	case *rbacv1.ClusterRole:
		return lc.launchClusterRole(obj.(*rbacv1.ClusterRole))
	case *rbacv1.ClusterRoleBinding:
		return lc.launchClusterRoleBinding(obj.(*rbacv1.ClusterRoleBinding))
	case *policyv1beta1.PodSecurityPolicy:
		return lc.launchPodSecurityPolicy(obj.(*policyv1beta1.PodSecurityPolicy))
	case *v1.PersistentVolume:
		return lc.launchPersistentVolume(obj.(*v1.PersistentVolume))
	case *v1.PersistentVolumeClaim:
		return lc.launchPersistentVolumeClaim(obj.(*v1.PersistentVolumeClaim))
	case *policyv1beta1.PodDisruptionBudget:
		return lc.launchPodDisruptionBudget(obj.(*policyv1beta1.PodDisruptionBudget))
	case *appsv1.StatefulSet:
		return lc.launchStatefulSet(obj.(*appsv1.StatefulSet))
	case *v1beta1.Ingress:
		return lc.launchIngress(obj.(*v1beta1.Ingress))
	default:
		log.Warn().Str("type", reflect.TypeOf(o).String()).Msg("Unknown entity")
		return derrors.NewUnimplementedError("object not supported").WithParams(o)
	}

	return derrors.NewInternalError("no case was executed")
}



func (lc * LaunchComponents) launchDaemonSet(daemonSet *appsv1.DaemonSet) derrors.Error {
	client := lc.Client.AppsV1().DaemonSets(daemonSet.Namespace)
	log.Debug().Interface("daemonSet", daemonSet).Msg("unmarshalled")
	created, err := client.Create(daemonSet)
	if err != nil {
		return derrors.AsError(err, "cannot create daemon set")
	}
	log.Debug().Interface("created", created).Msg("new daemon set has been created")
	return nil
}



func (lc * LaunchComponents) launchServiceAccount(serviceAccount *v1.ServiceAccount) derrors.Error {
	client := lc.Client.CoreV1().ServiceAccounts(serviceAccount.Namespace)
	log.Debug().Interface("serviceAccount", serviceAccount).Msg("unmarshalled")
	created, err := client.Create(serviceAccount)
	if err != nil {
		return derrors.AsError(err, "cannot create service account")
	}
	log.Debug().Interface("created", created).Msg("new service account has been created")
	return nil
}

func (lc * LaunchComponents) launchClusterRole(clusterRole *rbacv1.ClusterRole) derrors.Error {
	client := lc.Client.RbacV1().ClusterRoles()
	log.Debug().Interface("clusterRole", clusterRole).Msg("unmarshalled")
	created, err := client.Create(clusterRole)
	if err != nil {
		return derrors.AsError(err, "cannot create cluster role")
	}
	log.Debug().Interface("created", created).Msg("new cluster role has been created")
	return nil
}

func (lc * LaunchComponents) launchClusterRoleBinding(clusterRoleBinding *rbacv1.ClusterRoleBinding) derrors.Error {
	client := lc.Client.RbacV1().ClusterRoleBindings()
	log.Debug().Interface("clusterRoleBinding", clusterRoleBinding).Msg("unmarshalled")
	created, err := client.Create(clusterRoleBinding)
	if err != nil {
		return derrors.AsError(err, "cannot create cluster role binding")
	}
	log.Debug().Interface("created", created).Msg("new cluster role binding has been created")
	return nil
}

func (lc * LaunchComponents) launchRoleBinding(roleBinding *rbacv1.RoleBinding) derrors.Error {
	client := lc.Client.RbacV1().RoleBindings(roleBinding.Namespace)
	log.Debug().Interface("roleBinding", roleBinding).Msg("unmarshalled")
	created, err := client.Create(roleBinding)
	if err != nil {
		return derrors.AsError(err, "cannot create role binding")
	}
	log.Debug().Interface("created", created).Msg("new role binding has been created")
	return nil
}

func (lc * LaunchComponents) launchPodSecurityPolicy(policy *policyv1beta1.PodSecurityPolicy) derrors.Error {
	client := lc.Client.PolicyV1beta1().PodSecurityPolicies()
	log.Debug().Interface("policy", policy).Msg("unmarshalled")
	created, err := client.Create(policy)
	if err != nil {
		return derrors.AsError(err, "cannot create pod security policy")
	}
	log.Debug().Interface("created", created).Msg("new pod security policy has been created")
	return nil
}

func (lc * LaunchComponents) launchSecret(secret *v1.Secret) derrors.Error {
	client := lc.Client.CoreV1().Secrets(secret.Namespace)
	log.Debug().Interface("secret", secret).Msg("unmarshalled")
	created, err := client.Create(secret)
	if err != nil {
		return derrors.AsError(err, "cannot create secret")
	}
	log.Debug().Interface("created", created).Msg("new secret has been created")
	return nil
}

func (lc * LaunchComponents) createNamespace(name string) derrors.Error {
	namespaceClient := lc.Client.CoreV1().Namespaces()
	opts := metaV1.ListOptions{}
	list, err := namespaceClient.List(opts)
	if err != nil{
		return derrors.AsError(err, "cannot obtain the namespace list")
	}
	found := false
	for _, n := range list.Items {
		log.Debug().Interface("n", n).Msg("A namespace")
		if n.Name == name {
			found = true
			break
		}
	}

	if !found {
		toCreate := v1.Namespace{
			ObjectMeta: metaV1.ObjectMeta{
				Name:                       name,
			},
		}
		created, err := namespaceClient.Create(&toCreate)
		if err != nil {
			return derrors.AsError(err,"cannot create namespace")
		}
		log.Debug().Interface("created", created).Msg("namespaces has been created")
	}else{
		log.Debug().Str("namespace", name).Msg("namespace already exists")
	}
	return nil
}

func (lc * LaunchComponents) launchPersistentVolume(pv *v1.PersistentVolume) derrors.Error {
	client := lc.Client.CoreV1().PersistentVolumes()
	log.Debug().Interface("pv", pv).Msg("unmarshalled")
	created, err := client.Create(pv)
	if err != nil {
		return derrors.AsError(err, "cannot create persistent volume")
	}
	log.Debug().Interface("created", created).Msg("new persistent volume has been created")
	return nil
}

func (lc * LaunchComponents) launchPersistentVolumeClaim(pvc *v1.PersistentVolumeClaim) derrors.Error {
	client := lc.Client.CoreV1().PersistentVolumeClaims(pvc.Namespace)
	log.Debug().Interface("pvc", pvc).Msg("unmarshalled")
	created, err := client.Create(pvc)
	if err != nil {
		return derrors.AsError(err, "cannot create persistent volume claim")
	}
	log.Debug().Interface("created", created).Msg("new persistent volume claim has been created")
	return nil
}

func (lc * LaunchComponents) launchPodDisruptionBudget(pdb *policyv1beta1.PodDisruptionBudget) derrors.Error {
	client := lc.Client.PolicyV1beta1().PodDisruptionBudgets(pdb.Namespace)
	log.Debug().Interface("pdb", pdb).Msg("unmarshalled")
	created, err := client.Create(pdb)
	if err != nil {
		return derrors.AsError(err, "cannot create pod disruption budget")
	}
	log.Debug().Interface("created", created).Msg("new pod disruption budget")
	return nil
}

func (lc * LaunchComponents) launchStatefulSet(ss *appsv1.StatefulSet) derrors.Error {
	client := lc.Client.AppsV1().StatefulSets(ss.Namespace)
	log.Debug().Interface("pdb", ss).Msg("unmarshalled")
	created, err := client.Create(ss)
	if err != nil {
		return derrors.AsError(err, "cannot create stateful set")
	}
	log.Debug().Interface("created", created).Msg("new stateful set")
	return nil
}

func (lc * LaunchComponents) launchIngress(ingress *v1beta1.Ingress) derrors.Error {
	client := lc.Client.ExtensionsV1beta1().Ingresses(ingress.Namespace)
	log.Debug().Interface("ingress", ingress).Msg("unmarshalled")
	created, err := client.Create(ingress)
	if err != nil {
		return derrors.AsError(err, "cannot create ingress")
	}
	log.Debug().Interface("created", created).Msg("new ingress set")
	return nil
}

func (lc * LaunchComponents) String() string {
	return fmt.Sprintf("SYNC LaunchComponents from %s", lc.ComponentsDir)
}

func (lc * LaunchComponents) PrettyPrint(indentation int) string {
	return strings.Repeat(" ", indentation) + lc.String()
}

func (lc * LaunchComponents) UserString() string {
	return fmt.Sprintf("Launching K8s components from %s", lc.ComponentsDir)
}

