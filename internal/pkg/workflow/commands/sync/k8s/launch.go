/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package k8s

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/nalej/derrors"
	"github.com/nalej/grpc-installer-go"
	entities2 "github.com/nalej/installer/internal/pkg/entities"
	"github.com/nalej/installer/internal/pkg/errors"
	"github.com/nalej/installer/internal/pkg/workflow/entities"

	"github.com/rs/zerolog/log"

	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/api/core/v1"

	"k8s.io/client-go/kubernetes/scheme"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

const AzureStorageClass = "managed-premium"

var ProductionImagePullSecret = &v1.LocalObjectReference{
	Name: entities2.ProdRegistryName,
}

var StagingImagePullSecret = &v1.LocalObjectReference{
	Name: entities2.StagingRegistryName,
}

var DevImagePullSecret = &v1.LocalObjectReference{
	Name: entities2.DevRegistryName,
}

var ProductionImagePullSecrets = []v1.LocalObjectReference{*ProductionImagePullSecret}
var StagingImagePullSecrets = []v1.LocalObjectReference{*ProductionImagePullSecret, *StagingImagePullSecret}
var DevImagePullSecrets = []v1.LocalObjectReference{*ProductionImagePullSecret, *StagingImagePullSecret, *DevImagePullSecret}

// LaunchComponents is a command that reads a directory for YAML files and triggers the creation
// of those entities in Kubernetes.
type LaunchComponents struct {
	Kubernetes
	Namespaces    []string `json:"namespaces"`
	ComponentsDir string   `json:"componentsDir"`
	PlatformType  string   `json:"platform_type"`
	Environment   string   `json:"environment"`
}

// NewLaunchComponents creates a new LaunchComponents command.
func NewLaunchComponents(kubeConfigPath string, namespaces []string, componentsDir string, targetPlatform string) *LaunchComponents {
	return &LaunchComponents{
		Kubernetes: Kubernetes{
			GenericSyncCommand: *entities.NewSyncCommand(entities.LaunchComponents),
			KubeConfigPath:     kubeConfigPath,
		},
		Namespaces:    namespaces,
		ComponentsDir: componentsDir,
		PlatformType:  targetPlatform,
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

// Run the command.
func (lc *LaunchComponents) Run(workflowID string) (*entities.CommandResult, derrors.Error) {

	connectErr := lc.Connect()
	if connectErr != nil {
		return nil, connectErr
	}

	targetEnvironment, found := entities2.TargetEnvironmentFromString[lc.Environment]
	if !found {
		return nil, derrors.NewInvalidArgumentError("cannot determine target environment").WithParams(lc.Environment)
	}

	for _, target := range lc.Namespaces {
		createErr := lc.CreateNamespaceIfNotExists(target)
		if createErr != nil {
			return nil, createErr
		}
	}

	fileInfo, err := ioutil.ReadDir(lc.ComponentsDir)
	if err != nil {
		return nil, derrors.AsError(err, "cannot read components dir")
	}
	numLaunched := 0
	for _, file := range fileInfo {
		if strings.HasSuffix(file.Name(), ".yaml") {
			log.Info().Str("file", file.Name()).Msg("processing component")
			err := lc.launchComponent(path.Join(lc.ComponentsDir, file.Name()), targetEnvironment)
			if err != nil {
				return entities.NewCommandResult(false, "cannot launch component", err), nil
			}
			numLaunched++
		}
	}
	msg := fmt.Sprintf("%d components have been launched", numLaunched)
	return entities.NewCommandResult(true, msg, nil), nil
}

// ListComponents obtains a list of the files that need to be installed.
// TODO Overwrite files if a *.yaml.minikube file is found on the same entity with a MINIKUBE environment.
func (lc *LaunchComponents) ListComponents() []string {
	fileInfo, err := ioutil.ReadDir(lc.ComponentsDir)
	if err != nil {
		log.Fatal().Err(err).Str("componentsDir", lc.ComponentsDir).Msg("cannot read components dir")
	}
	result := make([]string, 0)
	for _, file := range fileInfo {
		if strings.HasSuffix(file.Name(), ".yaml") {
			result = append(result, file.Name())
		}
	}
	return result
}

// launchComponent triggers the creation of a given component from a YAML file
func (lc *LaunchComponents) launchComponent(componentPath string, targetEnvironment entities2.TargetEnvironment) derrors.Error {
	log.Debug().
		Str("path", componentPath).
		Str("targetEnvironment", entities2.TargetEnvironmentToString[targetEnvironment]).
		Msg("launch component")

	f, err := os.Open(componentPath)
	if err != nil {
		return derrors.NewPermissionDeniedError("cannot read component file", err)
	}
	defer f.Close()
	log.Debug().Str("path", componentPath).Msg("parsing component")

	// We use a YAML decoder to decode the resource straight into an
	// unstructured object. This way, we can deal with resources that are
	// not known to this client - like CustomResourceDefinitions
	obj := runtime.Object(&unstructured.Unstructured{})

	yamlDecoder := yaml.NewYAMLOrJSONDecoder(f, 1024)
	err = yamlDecoder.Decode(obj)
	if err != nil {
		return derrors.NewInvalidArgumentError("cannot parse component file", err)
	}
	gvk := obj.GetObjectKind().GroupVersionKind()
	log.Debug().Str("resource", gvk.String()).Msg("decoded resource")

	// Now let's see if it's a resource we know and can type, so we can
	// decide if we need to do some modifications. We ignore the error
	// because that just means we don't have the specific implementation of
	// the resource type and that's ok
	clientScheme := scheme.Scheme
	typed, _ := scheme.Scheme.New(gvk)
	if typed != nil {
		// Ah, we can convert this to something specific to deal with!
		err := clientScheme.Convert(obj, typed, nil)
		if err != nil {
			return derrors.NewInternalError("cannot convert resource to specific type", err)
		}
	}

	// Implement specific resource modifications for known types here. We
	// make sure to cast it to a generic object again so we can assign it
	// to the same variable as we had for the unstructured object.
	// obj -> typed -> o -> obj
	// We can do this switch even if typed might be nil.
	switch o := typed.(type) {
	case *appsv1.Deployment:
		obj = runtime.Object(lc.patchDeployment(o, targetEnvironment))
	case *v1.PersistentVolume:
		obj = runtime.Object(lc.patchPersistentVolume(o))
	case *v1.PersistentVolumeClaim:
		obj = runtime.Object(lc.patchPersistentVolumeClaim(o))
	}

	return lc.Create(obj)
}

// patchDeployment modifies the deployment to include image pull secrets depending on the type of environment.
func (lc *LaunchComponents) patchDeployment(deployment *appsv1.Deployment, targetEnvironment entities2.TargetEnvironment) *appsv1.Deployment {
	aux := deployment
	switch targetEnvironment {
	case entities2.Production:
		aux.Spec.Template.Spec.ImagePullSecrets = ProductionImagePullSecrets
	case entities2.Staging:
		aux.Spec.Template.Spec.ImagePullSecrets = StagingImagePullSecrets
	case entities2.Development:
		aux.Spec.Template.Spec.ImagePullSecrets = DevImagePullSecrets
	}
	return aux
}

// patchPersistenceVolume modifies the storage class
func (lc *LaunchComponents) patchPersistentVolume(pv *v1.PersistentVolume) *v1.PersistentVolume {
	if lc.PlatformType == grpc_installer_go.Platform_AZURE.String() {
		log.Debug().Msg("Modifying storageClass")
		patched := pv.DeepCopy()
		sc := AzureStorageClass
		patched.Spec.StorageClassName = sc
		pv = patched
	}
	return pv
}

// patchPersistenceVolumeClaim modifies the storage class of a pvc
func (lc *LaunchComponents) patchPersistentVolumeClaim(pvc *v1.PersistentVolumeClaim) *v1.PersistentVolumeClaim {
	if lc.PlatformType == grpc_installer_go.Platform_AZURE.String() {
		log.Debug().Msg("Modifying storageClass")
		patched := pvc.DeepCopy()
		sc := AzureStorageClass
		patched.Spec.StorageClassName = &sc
		pvc = patched
	}

	return pvc
}

func (lc *LaunchComponents) String() string {
	return fmt.Sprintf("SYNC LaunchComponents from %s", lc.ComponentsDir)
}

func (lc *LaunchComponents) PrettyPrint(indentation int) string {
	simpleIden := strings.Repeat(" ", indentation) + "  "
	entrySep := simpleIden + "  "
	cStr := ""
	for _, c := range lc.ListComponents() {
		cStr = cStr + "\n" + entrySep + c
	}
	return strings.Repeat(" ", indentation) + lc.String() + cStr
}

func (lc *LaunchComponents) UserString() string {
	return fmt.Sprintf("Launching K8s components from %s for %s", lc.ComponentsDir, lc.Environment)
}
