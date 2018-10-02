/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package k8s

import (
	"github.com/nalej/derrors"
	"github.com/rs/zerolog/log"
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

func (tc * TestCleaner) DeleteNamespaces() derrors.Error {
	dOpts := metaV1.DeleteOptions{}
	namespaceClient := tc.Client.CoreV1().Namespaces()
	for _, ns := range tc.Namespaces {
		err := namespaceClient.Delete(ns, &dOpts)
		if err != nil{
			return derrors.AsError(err, "cannot delete namespace")
		}
	}
	return nil
}
