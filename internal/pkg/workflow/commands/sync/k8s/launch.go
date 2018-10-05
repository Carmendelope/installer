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
	"k8s.io/client-go/kubernetes/scheme"
	"path"

	"io/ioutil"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"
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
	case *appsv1.Deployment:
		return lc.launchDeployment(obj.(*appsv1.Deployment))
	case *appsv1.DaemonSet:
		return derrors.NewUnimplementedError("DaemonSet launch")
	case *v1.Service:
		return lc.launchService(obj.(*v1.Service))
	default:
		return derrors.NewUnimplementedError("object not supported").WithParams(o)
	}

	return derrors.NewInternalError("no case was executed")
}

func (lc * LaunchComponents) launchDeployment(deployment *appsv1.Deployment) derrors.Error {
	deploymentClient := lc.Client.AppsV1().Deployments(deployment.Namespace)
	log.Debug().Interface("deployment", deployment).Msg("unmarshalled")
	created, err := deploymentClient.Create(deployment)
	if err != nil {
		return derrors.AsError(err, "cannot create deployment")
	}
	log.Debug().Interface("created", created).Msg("new component has been created")
	return nil
}

func (lc * LaunchComponents) launchService(service *v1.Service) derrors.Error {
	serviceClient := lc.Client.CoreV1().Services(service.Namespace)
	log.Debug().Interface("service", service).Msg("unmarshalled")
	created, err := serviceClient.Create(service)
	if err != nil {
		return derrors.AsError(err, "cannot create service")
	}
	log.Debug().Interface("created", created).Msg("new component has been created")
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

func (lc * LaunchComponents) String() string {
	return fmt.Sprintf("SYNC LaunchComponents from %s", lc.ComponentsDir)
}

func (lc * LaunchComponents) PrettyPrint(indentation int) string {
	return strings.Repeat(" ", indentation) + lc.String()
}

func (lc * LaunchComponents) UserString() string {
	return fmt.Sprintf("Launching K8s components from %s", lc.ComponentsDir)
}

