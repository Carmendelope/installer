/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package k8s

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/nalej/derrors"
	"github.com/nalej/grpc-installer-go"
	entities2 "github.com/nalej/installer/internal/pkg/entities"
	"github.com/nalej/installer/internal/pkg/errors"
	"github.com/nalej/installer/internal/pkg/workflow/entities"

	"github.com/rs/zerolog/log"

	"k8s.io/api/core/v1"

	"k8s.io/client-go/kubernetes/scheme"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

const AzureStorageClass = "managed-premium"

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
	// Get the preprocessed list of components to be installed on the target Kubernetes.
	components, err := lc.ListComponents()
	if err != nil {
		return nil, err
	}

	numLaunched := 0
	for _, fileName := range components {
		log.Info().Str("fileName", fileName).Msg("processing component")
		err := lc.launchComponent(path.Join(lc.ComponentsDir, fileName), targetEnvironment)
		if err != nil {
			return entities.NewCommandResult(false, "cannot launch component", err), nil
		}
		numLaunched++
	}
	msg := fmt.Sprintf("%d components have been launched", numLaunched)
	return entities.NewCommandResult(true, msg, nil), nil
}

// ListComponents obtains a list of the files that need to be installed. Platform dependent YAML files overwrite the
// use of the common YAML. For example, if the install is for an Azure cluster, and there are a component.yaml and
// component.yaml.azure files, the later will be used.
func (lc *LaunchComponents) ListComponents() ([]string, derrors.Error) {
	fileInfo, err := ioutil.ReadDir(lc.ComponentsDir)
	if err != nil {
		log.Warn().Err(err).Str("componentsDir", lc.ComponentsDir).Msg("cannot read components dir")
		return nil, derrors.NewInternalError("cannot read component directory", err)
	}
	filesToCreate := make(map[string]bool, 0)
	platformName := strings.ToLower(lc.PlatformType)
	platformSuffix := fmt.Sprintf(".yaml.%s", platformName)
	for _, file := range fileInfo {
		log.Info().Str("fileName", file.Name()).Str("platformSuffix", platformSuffix).Msg("Checking file")
		if strings.HasSuffix(file.Name(), platformSuffix) {
			log.Info().Msg("file has platform suffix, addint to list")
			// A platform specific file is found, delete the common one if exists
			platformIndependentName := strings.TrimSuffix(file.Name(), fmt.Sprintf(".%s", platformName))
			delete(filesToCreate, platformIndependentName)
			// Add the platform specific file to the list.
			filesToCreate[file.Name()] = true
		} else if strings.HasSuffix(file.Name(), ".yaml") {
			log.Info().Msg("file is platform independent")
			// Check if the platform specific equivalent is found
			_, exists := filesToCreate[fmt.Sprintf("%s.%s", file.Name(), platformName)]
			if !exists {
				log.Info().Msg("adding file to list")
				filesToCreate[file.Name()] = true
			}
		}
	}

	result := make([]string, 0)
	for toAdd, _ := range filesToCreate {
		result = append(result, toAdd)
	}
	// Make sure to main the same order as in the listing of the original files.
	sort.Strings(result)
	return result, nil
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
	case *v1.PersistentVolume:
		obj = runtime.Object(lc.patchPersistentVolume(o))
	case *v1.PersistentVolumeClaim:
		obj = runtime.Object(lc.patchPersistentVolumeClaim(o))
	}

	return lc.Create(obj)
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
	components, err := lc.ListComponents()
	if err != nil {
		log.Warn().Err(err).Msg("cannot list components")
		cStr = cStr + "\n" + entrySep + "<unknown>"
	} else {
		for _, c := range components {
			cStr = cStr + "\n" + entrySep + c
		}
	}
	return strings.Repeat(" ", indentation) + lc.String() + cStr
}

func (lc *LaunchComponents) UserString() string {
	return fmt.Sprintf("Launching K8s components from %s for %s", lc.ComponentsDir, lc.Environment)
}
