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
	"github.com/nalej/derrors"
	"github.com/nalej/installer/internal/pkg/errors"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
	"github.com/rs/zerolog/log"
	v1 "k8s.io/api/apps/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
	"time"
)

const (
	maxRetries = 30
)

// CheckComponents is a command that reads a directory for YAML files and checks the readiness
// of those entities in Kubernetes.
type CheckComponents struct {
	Kubernetes
	Namespaces    []string `json:"namespaces"`
	//ComponentsDir string   `json:"componentsDir"`
	//PlatformType  string   `json:"platform_type"`
}

// NewCheckComponents creates a new CheckComponents command.
func NewCheckComponents(kubeConfigPath string, namespaces []string, componentsDir string, targetPlatform string) *CheckComponents {
	return &CheckComponents{
		Kubernetes: Kubernetes{
			GenericSyncCommand: *entities.NewSyncCommand(entities.CheckComponents),
			KubeConfigPath:     kubeConfigPath,
		},
		Namespaces:    namespaces,
		//ComponentsDir: componentsDir,
		//PlatformType:  targetPlatform,
	}
}

type PlatformResources struct {
	Daemonsets []v1.DaemonSet
	StatefulSets []v1.StatefulSet
	Deployments []v1.Deployment
}

func NewEmptyPlatformResources () *PlatformResources {
	return &PlatformResources{
		Daemonsets:   make([]v1.DaemonSet, 0),
		StatefulSets: make([]v1.StatefulSet, 0),
		Deployments:  make([]v1.Deployment, 0),
	}
}

// NewCheckComponentsFromJSON creates an CheckComponents command from a JSON object.
func NewCheckComponentsFromJSON(raw []byte) (*entities.Command, derrors.Error) {
	cc := &CheckComponents{}
	if err := json.Unmarshal(raw, &cc); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	cc.CommandID = entities.GenerateCommandID(cc.Name())
	var r entities.Command = cc
	return &r, nil
}

// Run the command.
func (cc *CheckComponents) Run(workflowID string) (*entities.CommandResult, derrors.Error) {

	connectErr := cc.Connect()
	if connectErr != nil {
		return nil, connectErr
	}

	for _, target := range cc.Namespaces {
		createErr := cc.CreateNamespaceIfNotExists(target)
		if createErr != nil {
			return nil, createErr
		}
	}
	// Get the preprocessed list of resources to be checked on the target platform.
	resources, err := cc.RetrieveResources()
	if err != nil {
		return nil, err
	}

	// Check daemonsets
	daemonsetsChecked := 0
	for _, daemonset := range resources.Daemonsets {
		log.Info().Str("daemonsetName", daemonset.Name).Msg("checking daemonset")
		err := cc.checkDaemonset(daemonset)
		if err != nil {
			return entities.NewCommandResult(false, "cannot check daemonset", err), nil
		}
		daemonsetsChecked++
	}
	dsMsg := fmt.Sprintf("%d daemonsets have been checked\n", daemonsetsChecked)

	// Check statefulsets
	statefulsetsChecked := 0
	for _, statefulset := range resources.StatefulSets {
		log.Info().Str("statefulsetName", statefulset.Name).Msg("checking statefulset")
		err := cc.checkStatefulset(statefulset)
		if err != nil {
			return entities.NewCommandResult(false, "cannot check statefulset", err), nil
		}
		statefulsetsChecked++
	}
	ssMsg := fmt.Sprintf("%d statefulsets have been checked\n", statefulsetsChecked)

	// Check deployments
	deploymentsChecked := 0
	for _, deployment := range resources.Deployments {
		log.Info().Str("deploymentName", deployment.Name).Msg("checking deployment")
		err := cc.checkDeployment(deployment)
		if err != nil {
			return entities.NewCommandResult(false, "cannot check deployment", err), nil
		}
		statefulsetsChecked++
	}
	dMsg := fmt.Sprintf("%d deployments have been checked\n", deploymentsChecked)

	msg := dsMsg + ssMsg + dMsg

	return entities.NewCommandResult(true, msg, nil), nil
}

func (cc *CheckComponents) RetrieveResources() (*PlatformResources, derrors.Error) {
	namespaces := cc.Namespaces
	emptyOpts := metaV1.ListOptions{}
	resources := NewEmptyPlatformResources()
	for _, ns := range namespaces {
		daemonsets, dsErr := cc.Client.AppsV1().DaemonSets(ns).List(emptyOpts)
		if dsErr != nil {
			return nil, derrors.NewGenericError(dsErr.Error())
		}
		statefulSets, rsErr := cc.Client.AppsV1().StatefulSets(ns).List(emptyOpts)
		if dsErr != nil {
			return nil, derrors.NewGenericError(rsErr.Error())
		}
		deployments, rsErr := cc.Client.AppsV1().Deployments(ns).List(emptyOpts)
		if dsErr != nil {
			return nil, derrors.NewGenericError(rsErr.Error())
		}

		resources.Daemonsets = append(resources.Daemonsets, daemonsets.Items...)
		resources.StatefulSets = append(resources.StatefulSets, statefulSets.Items...)
		resources.Deployments = append(resources.Deployments, deployments.Items...)
	}
	return resources, nil
}

func (cc *CheckComponents) checkDaemonset(ds v1.DaemonSet) derrors.Error {
	for i := 0; i < maxRetries; i++ {
		if ds.Status.NumberAvailable == ds.Status.NumberReady {
			log.Debug().Str("daemonset", ds.Name).Msg("daemonset ready")
			return nil
		} else {
			log.Debug().Str("daemonset", ds.Name).Msg("daemonset not ready, waiting 30s")
			time.Sleep(30 * time.Second)
			i += 1
		}
	}
	return derrors.NewUnavailableError("daemonset unavailable")
}

func (cc *CheckComponents) checkStatefulset(ss v1.StatefulSet) derrors.Error {
	for i := 0; i < maxRetries; i++ {
		if ss.Status.Replicas == ss.Status.ReadyReplicas {
			log.Debug().Str("statefulset", ss.Name).Msg("statefulset ready")
			return nil
		} else {
			log.Debug().Str("statefulset", ss.Name).Msg("statefulset not ready, waiting 30s")
			time.Sleep(30 * time.Second)
			i += 1
		}
	}
	return derrors.NewUnavailableError("statefulset unavailable")
}

func (cc *CheckComponents) checkDeployment(d v1.Deployment) derrors.Error {
	for i := 0; i < maxRetries; i++ {
		if d.Status.Replicas == d.Status.ReadyReplicas {
			log.Debug().Str("deployment", d.Name).Msg("deployment ready")
			return nil
		} else {
			log.Debug().Str("deployment", d.Name).Msg("deployment not ready, waiting 30s")
			time.Sleep(30 * time.Second)
			i += 1
		}
	}
	return derrors.NewUnavailableError("deployment unavailable")
}

func (cc *CheckComponents) String() string {
	return fmt.Sprintf("SYNC CheckComponents from %s", cc.Namespaces)
}

func (cc *CheckComponents) PrettyPrint(indentation int) string {
	simpleIden := strings.Repeat(" ", indentation) + "  "
	entrySep := simpleIden + "  "
	cStr := ""
	resources, err := cc.RetrieveResources()
	if err != nil {
		log.Warn().Err(err).Msg("cannot retrieve resources")
		cStr = cStr + "\n" + entrySep + "<unknown>"
	} else {
		for _, c := range resources.Daemonsets {
			cStr = cStr + "\n" + entrySep + c.Name
		}
		for _, c := range resources.StatefulSets {
			cStr = cStr + "\n" + entrySep + c.Name
		}
		for _, c := range resources.Deployments {
			cStr = cStr + "\n" + entrySep + c.Name
		}
	}
	return strings.Repeat(" ", indentation) + cc.String() + cStr
}

func (cc *CheckComponents) UserString() string {
	return fmt.Sprintf("Checking K8s resources from %s", cc.Namespaces)
}