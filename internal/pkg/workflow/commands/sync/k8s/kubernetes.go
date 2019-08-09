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
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/restmapper"
	"k8s.io/client-go/tools/clientcmd"
	"net"
)

type Kubernetes struct {
	entities.GenericSyncCommand
	KubeConfigPath string		`json:"kubeConfigPath"`
	Client	 *kubernetes.Clientset `json:"-"`

	// Discovery client for REST mapper to use, so we can figure out
	// the right endpoints for reserves
	discoveryClient *discovery.DiscoveryClient
	// Dynamic client used to create all resources
	dynClient dynamic.Interface
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

	// Create the discovery client
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		return derrors.NewInternalError("failed to create discovery client", err)
	}
	k.discoveryClient = discoveryClient

	// Create the dynamic client that can be used to create any object
	// by specifying what resource we're dealing with by using the REST mapper
	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return derrors.NewInternalError("failed to create dynamic client", err)
	}
	k.dynClient = dynClient

	return nil
}

func (k *Kubernetes) ResolveIP(address string) ([]string, derrors.Error){
	result := make([]string, 0)
	ips, err := net.LookupIP(address)
	if err != nil {
		log.Error().Err(err).Str("address", address).Msg("cannot resolve IP address")
		return nil, derrors.AsError(err, "cannot resolve IP address")
	}
	for _, ip := range ips {
		if len(ip) == net.IPv4len {
			result = append(result, ip.String())
		}
	}
	return result, nil
}

func (k *Kubernetes) ExistsNamespace(name string) (bool, derrors.Error) {
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

func (k *Kubernetes) CreateNamespace(name string) derrors.Error {
	toCreate := v1.Namespace{
		TypeMeta: metaV1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name: name,
		},
	}
	err := k.Create(&toCreate)
	if err != nil {
		return derrors.AsError(err, "cannot create namespace")
	}

	return nil
}

func (k *Kubernetes) CreateNamespaceIfNotExists(name string) derrors.Error {
	found, fErr := k.ExistsNamespace(name)
	if fErr != nil {
		return fErr
	}

	if !found {
		err := k.CreateNamespace(name)
		if err != nil {
			return err
		}
	} else {
		log.Debug().Str("namespace", name).Msg("namespace already exists")
	}
	return nil
}


func (k *Kubernetes) Create(obj runtime.Object) derrors.Error {
	// Create unstructured object
	unstructuredMap, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return derrors.NewInvalidArgumentError("cannot convert object to unstructured", err).WithParams(obj)
	}
	unstructuredObj := &unstructured.Unstructured{
		Object: unstructuredMap,
	}

	gvk, derr := getKind(obj)
	if derr != nil {
		return derr
	}

	// Items in list resources need to be sent to the server one by one
	if unstructuredObj.IsList() {
		log.Debug().Str("resource", gvk.String()).Msg("creating each item in list resource")
		list, err := unstructuredObj.ToList()
		if err != nil {
			return derrors.NewInternalError("cannot create unstructured list", err)
		}
		err = list.EachListItem(func (obj runtime.Object) error { return k.Create(obj).(error) })
		if err != nil {
			return err.(derrors.Error)
		}
		log.Debug().Str("resource", gvk.String()).Msg("created all items in list resource")
		return nil
	}

	// Create the REST mapper through a discovery client
	// We do this every time we create a resource, because if we created
	// a custom resource definition in a previous step, we need to
	// update the list of supported resources.
	resources, err := restmapper.GetAPIGroupResources(k.discoveryClient)
	if err != nil {
		return derrors.NewInternalError("failed to get api group resources", err)
	}
	mapper := restmapper.NewDiscoveryRESTMapper(resources)

	// Get the right REST endpoint through the mapper
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		return derrors.NewInternalError("unable to get REST mapping for object", err).WithParams(unstructuredObj)
	}

	var client dynamic.ResourceInterface
	namespace := unstructuredObj.GetNamespace()
	if namespace != "" {
		client = k.dynClient.Resource(mapping.Resource).Namespace(namespace)
	} else {
		client = k.dynClient.Resource(mapping.Resource)
	}

	log.Debug().Interface("obj", unstructuredObj).Msg("creating resource")

	created, err := client.Create(unstructuredObj, metaV1.CreateOptions{})
	if err != nil {
		return derrors.NewInternalError("unable to create object", err).WithParams(unstructuredObj)
	}

	log.Debug().Str("resource", created.GetSelfLink()).Msg("created")

	return nil
}

func getKind(obj interface{}) (schema.GroupVersionKind, derrors.Error) {
	kinds, _, err := scheme.Scheme.ObjectKinds(obj.(runtime.Object))
	if err != nil {
		return schema.GroupVersionKind{}, derrors.NewInvalidArgumentError("invalid object received")
	}

	// Not sure what to do if an object matches multiple kinds, let's
	// at least warn
	if len(kinds) > 1 {
		kindLog := log.Warn()
		for _, k := range(kinds) {
			kindLog = kindLog.Str("candidate", k.String())
		}
		kindLog.Msg("received ambiguous object, picking first candidate")
	}

	kind := kinds[0]

	return kind, nil
}
