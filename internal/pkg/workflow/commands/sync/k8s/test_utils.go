/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package k8s

import (
	"github.com/nalej/derrors"
	"github.com/rs/zerolog/log"
	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

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

	numDeleted := 0
	for _, ns := range tc.Namespaces {
		deploymentClient := tc.Client.AppsV1().Deployments(ns)
		opts := metaV1.ListOptions{}
		deploymentList, err := deploymentClient.List(opts)
		if err != nil {
			return derrors.AsError(err, "cannot list deployments")
		}
		dOpts := metaV1.DeleteOptions{}
		for _, d := range deploymentList.Items {
			err := deploymentClient.Delete(d.Name, &dOpts)
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
		log.Debug().Interface("n", n).Msg("A namespace")
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