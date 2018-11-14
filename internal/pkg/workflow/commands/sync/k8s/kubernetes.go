/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package k8s

import (
	"github.com/nalej/derrors"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
	"github.com/rs/zerolog/log"
	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
