/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package k8s

import (
	"github.com/nalej/derrors"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

const KubeSystemNamespace = "kube-system"

type TestCleaner struct {
	Namespaces []string `json:"namespace"`
	KubeConfigPath string `json:"kubeConfig"`
	Client * kubernetes.Clientset `json:"-"`
}

func NewTestCleaner(kubeConfigPath string, namespaces ...string) * TestCleaner {
	return &TestCleaner{
		Namespaces: namespaces,
		KubeConfigPath: kubeConfigPath,
	}
}

func (tc * TestCleaner) Connect() derrors.Error {

	config, err := clientcmd.BuildConfigFromFlags("", tc.KubeConfigPath)
	if err != nil {
		log.Error().Err(err).Msg("error building configuration from kubeconfig")
		return derrors.AsError(err, "error building configuration from kubeconfig")
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error().Err(err).Msg("error using configuration to build k8s clientset")
		return derrors.AsError(err,"error using configuration to build k8s clientset")
	}

	tc.Client = clientset
	return nil
}

func (tc * TestCleaner) DeleteAll() derrors.Error {
	err := tc.DeleteDeployments()
	if err != nil {
		return err
	}
	err = tc.DeleteServices()
	if err != nil {
		return err
	}
	err = tc.DeleteServiceAccounts()
	if err != nil {
		return err
	}
	err = tc.DeleteClusterRoles()
	if err != nil{
		return err
	}
	err = tc.DeleteRoles()
	if err != nil{
		return err
	}
	err = tc.DeleteRoleBindings()
	if err != nil{
		return err
	}
	err = tc.DeleteClusterRoleBindings()
	if err != nil{
		return err
	}
	err = tc.DeleteConfigMaps()
	if err != nil{
		return err
	}
	err = tc.DeleteNamespaces()
	if err != nil {
		return err
	}
	return nil
}

func (tc * TestCleaner) DeleteDeployments() derrors.Error {
	if tc.Client == nil {
		err := tc.Connect()
		if err != nil{
			return err
		}
	}

	opts := metaV1.ListOptions{}
	dOpts := metaV1.DeleteOptions{}
	numDeleted := 0
	for _, ns := range tc.Namespaces {
		deploymentClient := tc.Client.AppsV1().Deployments(ns)
		deploymentList, err := deploymentClient.List(opts)
		if err != nil {
			return derrors.AsError(err, "cannot list deployments")
		}
		for _, d := range deploymentList.Items {
			log.Debug().Str("name", d.Name).Msg("deleting deployment")
			err := deploymentClient.Delete(d.Name, &dOpts)
			if err != nil {
				return derrors.AsError(err, "cannot delete deployment")
			}
			numDeleted++
		}
	}

	// Kube-System
	client := tc.Client.AppsV1().Deployments(KubeSystemNamespace)
	ds, err := client.List(opts)
	if err != nil {
		return derrors.AsError(err, "cannot list config maps")
	}
	for _, d := range ds.Items {
		_, exists := d.Labels["cluster"]
		if exists {
			log.Debug().Str("name", d.Name).Msg("deleting deployment")
			err := client.Delete(d.Name, &dOpts)
			if err != nil {
				return derrors.AsError(err, "cannot delete deployment")
			}
			numDeleted++
		}
	}

	log.Debug().Int("deleted", numDeleted).Msg("deployments deleted")
	return nil
}

func (tc * TestCleaner) existNamespace(name string) (bool, derrors.Error) {
	namespaceClient := tc.Client.CoreV1().Namespaces()
	opts := metaV1.ListOptions{}
	list, err := namespaceClient.List(opts)
	if err != nil{
		return false, derrors.AsError(err, "cannot obtain the namespace list")
	}
	found := false
	for _, n := range list.Items {
		if n.Name == name {
			found = true
			break
		}
	}
	return found, nil
}

func (tc * TestCleaner) DeleteNamespaces() derrors.Error {
	dOpts := metaV1.DeleteOptions{}
	namespaceClient := tc.Client.CoreV1().Namespaces()
	for _, ns := range tc.Namespaces {

		found, fErr := tc.existNamespace(ns)
		if fErr != nil{
			return fErr
		}
		if found {
			err := namespaceClient.Delete(ns, &dOpts)
			if err != nil{
				return derrors.AsError(err, "cannot delete namespace")
			}
			log.Debug().Str("namespace", ns).Msg("deleted")
		}
	}
	log.Debug().Int("deleted", len(tc.Namespaces)).Msg("namespaces deleted")
	return nil
}


func (tc * TestCleaner) DeleteServices() derrors.Error {
	numDeleted := 0
	for _, ns := range tc.Namespaces {
		serviceClient := tc.Client.CoreV1().Services(ns)
		opts := metaV1.ListOptions{}
		serviceList, err := serviceClient.List(opts)
		if err != nil {
			return derrors.AsError(err, "cannot list services")
		}
		dOpts := metaV1.DeleteOptions{}
		for _, s := range serviceList.Items {
			log.Debug().Str("name", s.Name).Msg("deleting service")
			err := serviceClient.Delete(s.Name, &dOpts)
			if err != nil {
				return derrors.AsError(err, "cannot delete service")
			}
			numDeleted++
		}
	}
	// Kube-System
	serviceClient := tc.Client.CoreV1().Services(KubeSystemNamespace)
	opts := metaV1.ListOptions{}
	serviceList, err := serviceClient.List(opts)
	if err != nil {
		return derrors.AsError(err, "cannot list services")
	}
	dOpts := metaV1.DeleteOptions{}
	for _, s := range serviceList.Items {
		_, exists := s.Labels["cluster"]
		if exists {
			log.Debug().Str("name", s.Name).Msg("deleting service")
			err := serviceClient.Delete(s.Name, &dOpts)
			if err != nil {
				return derrors.AsError(err, "cannot delete service")
			}
			numDeleted++
		}
	}
	log.Debug().Int("deleted", numDeleted).Msg("services deleted")
	return nil
}

func (tc * TestCleaner) DeleteServiceAccounts() derrors.Error {
	numDeleted := 0
	for _, ns := range tc.Namespaces {
		client := tc.Client.CoreV1().ServiceAccounts(ns)
		opts := metaV1.ListOptions{}
		accounts, err := client.List(opts)
		if err != nil {
			return derrors.AsError(err, "cannot list services")
		}
		dOpts := metaV1.DeleteOptions{}
		for _, acc := range accounts.Items{
			log.Debug().Str("name", acc.Name).Msg("deleting service account")
			err := client.Delete(acc.Name, &dOpts)
			if err != nil {
				return derrors.AsError(err, "cannot delete service account")
			}
			numDeleted++
		}
	}

	// Kube-System
	client := tc.Client.CoreV1().ServiceAccounts(KubeSystemNamespace)
	opts := metaV1.ListOptions{}
	serviceList, err := client.List(opts)
	if err != nil {
		return derrors.AsError(err, "cannot list services")
	}
	dOpts := metaV1.DeleteOptions{}
	for _, s := range serviceList.Items {
		_, exists := s.Labels["cluster"]
		if exists {
			log.Debug().Str("name", s.Name).Msg("deleting service account")
			err := client.Delete(s.Name, &dOpts)
			if err != nil {
				return derrors.AsError(err, "cannot delete service")
			}
			numDeleted++
		}
	}
	log.Debug().Int("deleted", numDeleted).Msg("service accounts deleted")
	return nil
}

func (tc * TestCleaner) DeleteClusterRoles() derrors.Error {
	numDeleted := 0

	client := tc.Client.RbacV1().ClusterRoles()
	opts := metaV1.ListOptions{}
	roles, err := client.List(opts)
	if err != nil {
		return derrors.AsError(err, "cannot list services")
	}
	dOpts := metaV1.DeleteOptions{}

		for _, cr := range roles.Items{
			_, exists := cr.Labels["cluster"]
			if exists {
				log.Debug().Str("name", cr.Name).Msg("deleting cluster role")
				err := client.Delete(cr.Name, &dOpts)
				if err != nil {
					return derrors.AsError(err, "cannot delete cluster role")
				}
				numDeleted++
			}
		}



	log.Debug().Int("deleted", numDeleted).Msg("service cluster roles deleted")
	return nil
}

func (tc * TestCleaner) DeleteRoles() derrors.Error {
	numDeleted := 0
	for _, ns := range tc.Namespaces {
		client := tc.Client.RbacV1().Roles(ns)
		opts := metaV1.ListOptions{}
		roles, err := client.List(opts)
		if err != nil {
			return derrors.AsError(err, "cannot list roles")
		}
		dOpts := metaV1.DeleteOptions{}
		for _, rol := range roles.Items{
			log.Debug().Str("name", rol.Name).Msg("deleting role")
			err := client.Delete(rol.Name, &dOpts)
			if err != nil {
				return derrors.AsError(err, "cannot delete role")
			}
			numDeleted++
		}
	}

	// Kube-System
	client := tc.Client.RbacV1().Roles(KubeSystemNamespace)
	opts := metaV1.ListOptions{}
	roleList, err := client.List(opts)
	if err != nil {
		return derrors.AsError(err, "cannot list roles")
	}
	dOpts := metaV1.DeleteOptions{}
	for _, r := range roleList.Items {
		_, exists := r.Labels["cluster"]
		if exists {
			log.Debug().Str("name", r.Name).Msg("deleting role")
			err := client.Delete(r.Name, &dOpts)
			if err != nil {
				return derrors.AsError(err, "cannot delete role")
			}
			numDeleted++
		}
	}
	log.Debug().Int("deleted", numDeleted).Msg("roles deleted")
	return nil
}

func (tc * TestCleaner) DeleteClusterRoleBindings() derrors.Error {
	numDeleted := 0

	client := tc.Client.RbacV1().ClusterRoleBindings()
	opts := metaV1.ListOptions{}
	roles, err := client.List(opts)
	if err != nil {
		return derrors.AsError(err, "cannot list cluster role bindings")
	}
	dOpts := metaV1.DeleteOptions{}

	for _, cr := range roles.Items{
		_, exists := cr.Labels["cluster"]
		if exists {
			log.Debug().Str("name", cr.Name).Msg("deleting cluster role binding")
			err := client.Delete(cr.Name, &dOpts)
			if err != nil {
				return derrors.AsError(err, "cannot delete cluster role binding")
			}
			numDeleted++
		}
	}

	log.Debug().Int("deleted", numDeleted).Msg("service cluster role bindings deleted")
	return nil
}

func (tc * TestCleaner) DeleteRoleBindings() derrors.Error {
	numDeleted := 0
	for _, ns := range tc.Namespaces {
		client := tc.Client.RbacV1().RoleBindings(ns)
		opts := metaV1.ListOptions{}
		roles, err := client.List(opts)
		if err != nil {
			return derrors.AsError(err, "cannot list roles")
		}
		dOpts := metaV1.DeleteOptions{}
		for _, rol := range roles.Items{
			log.Debug().Str("name", rol.Name).Msg("deleting role")
			err := client.Delete(rol.Name, &dOpts)
			if err != nil {
				return derrors.AsError(err, "cannot delete role")
			}
			numDeleted++
		}
	}

	// Kube-System
	client := tc.Client.RbacV1().RoleBindings(KubeSystemNamespace)
	opts := metaV1.ListOptions{}
	roleList, err := client.List(opts)
	if err != nil {
		return derrors.AsError(err, "cannot list roles")
	}
	dOpts := metaV1.DeleteOptions{}
	for _, r := range roleList.Items {
		_, exists := r.Labels["cluster"]
		if exists {
			log.Debug().Str("name", r.Name).Msg("deleting role")
			err := client.Delete(r.Name, &dOpts)
			if err != nil {
				return derrors.AsError(err, "cannot delete role")
			}
			numDeleted++
		}
	}
	log.Debug().Int("deleted", numDeleted).Msg("roles deleted")
	return nil
}

func (tc * TestCleaner) DeleteConfigMaps() derrors.Error {
	numDeleted := 0
	opts := metaV1.ListOptions{}
	dOpts := metaV1.DeleteOptions{}
	for _, ns := range tc.Namespaces {
		client := tc.Client.CoreV1().ConfigMaps(ns)
		cms, err := client.List(opts)
		if err != nil {
			return derrors.AsError(err, "cannot list config maps")
		}
		for _, cm := range cms.Items{
			log.Debug().Str("name", cm.Name).Msg("deleting config map")
			err := client.Delete(cm.Name, &dOpts)
			if err != nil {
				return derrors.AsError(err, "cannot delete config map")
			}
			numDeleted++
		}
	}

	// Kube-System
	client := tc.Client.CoreV1().ConfigMaps(KubeSystemNamespace)
	cms, err := client.List(opts)
	if err != nil {
		return derrors.AsError(err, "cannot list config maps")
	}
	for _, cm := range cms.Items {
		_, exists := cm.Labels["cluster"]
		if exists {
			log.Debug().Str("name", cm.Name).Msg("deleting config map")
			err := client.Delete(cm.Name, &dOpts)
			if err != nil {
				return derrors.AsError(err, "cannot delete config map")
			}
			numDeleted++
		}
	}
	log.Debug().Int("deleted", numDeleted).Msg("config maps deleted")
	return nil
}


type TestK8sUtils struct {
	KubeConfigPath string
	Client * kubernetes.Clientset
}

func NewTestK8sUtils(kubeConfigPath string) * TestK8sUtils {
	return &TestK8sUtils{
		KubeConfigPath: kubeConfigPath,
	}
}

func (tu * TestK8sUtils) Connect() derrors.Error {

	config, err := clientcmd.BuildConfigFromFlags("", tu.KubeConfigPath)
	if err != nil {
		log.Error().Err(err).Msg("error building configuration from kubeconfig")
		return derrors.AsError(err, "error building configuration from kubeconfig")
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error().Err(err).Msg("error using configuration to build k8s clientset")
		return derrors.AsError(err,"error using configuration to build k8s clientset")
	}

	tu.Client = clientset
	return nil
}


func (tu * TestK8sUtils) CreateNamespace(name string) derrors.Error {
	namespaceClient := tu.Client.CoreV1().Namespaces()
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

type TestChecker struct {
	KubeConfigPath string `json:"kubeConfig"`
	Client * kubernetes.Clientset `json:"-"`
}

func NewTestChecker(kubeConfigPath string) * TestChecker {
	return &TestChecker{
		KubeConfigPath: kubeConfigPath,
	}
}

func (tc * TestChecker) Connect() derrors.Error {

	config, err := clientcmd.BuildConfigFromFlags("", tc.KubeConfigPath)
	if err != nil {
		log.Error().Err(err).Msg("error building configuration from kubeconfig")
		return derrors.AsError(err, "error building configuration from kubeconfig")
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error().Err(err).Msg("error using configuration to build k8s clientset")
		return derrors.AsError(err,"error using configuration to build k8s clientset")
	}

	tc.Client = clientset
	return nil
}

func (tc * TestChecker) GetSecret(secretName string, namespace string) *v1.Secret {
	sc := tc.Client.CoreV1().Secrets(namespace)
	opts := metaV1.GetOptions{}
	found, err := sc.Get(secretName, opts)
	gomega.Expect(err).To(gomega.Succeed())
	return found
}

