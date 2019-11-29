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
	"strings"
)

// DeletePodSecurityPolicy structure with the attributes required to delete a given pod security policy.
type DeletePodSecurityPolicy struct {
	// Kubernetes embedded object
	Kubernetes
	// PolicyName with the name of the target policy
	PolicyName string `json:"policy_name"`
	// FailIfNotExists flag determines if the command fails in case the namespace does not exits.
	FailIfNotExists bool `json:"fail_if_not_exists"`
}

// NewDeletePodSecurityPolicy creates a new DeletePodSecurityPolicy command
func NewDeletePodSecurityPolicy(kubeConfigPath string, policyName string) *DeletePodSecurityPolicy {
	return &DeletePodSecurityPolicy{
		Kubernetes: Kubernetes{
			GenericSyncCommand: *entities.NewSyncCommand(entities.DeletePodSecurityPolicy),
			KubeConfigPath:     kubeConfigPath,
		},
		PolicyName: policyName,
	}
}

// NewDeletePodSecurityPolicyFromJSON creates a new DeletePodSecurityPolicy command from a raw JSON representation.
func NewDeletePodSecurityPolicyFromJSON(raw []byte) (*entities.Command, derrors.Error) {
	cmc := &DeletePodSecurityPolicy{}
	if err := json.Unmarshal(raw, &cmc); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	cmc.CommandID = entities.GenerateCommandID(cmc.Name())
	var r entities.Command = cmc
	return &r, nil
}

// Run the current command returning the result or an error.
func (dpsp *DeletePodSecurityPolicy) Run(workflowID string) (*entities.CommandResult, derrors.Error) {
	connectErr := dpsp.Connect()
	if connectErr != nil {
		return nil, connectErr
	}

	exists, err := dpsp.ExistsEntity("", "policy", "v1", "podsecuritypolicies", dpsp.PolicyName)
	if err != nil {
		return entities.NewCommandResult(false, "cannot determine if the pod security policy exists", err), nil
	}
	log.Debug().Str("policyName", dpsp.PolicyName).Bool("exists", exists).Msg("pod security policy check")

	if exists {
		err := dpsp.DeleteEntity("", "policy", "v1", "podsecuritypolicies", dpsp.PolicyName)
		if err != nil {
			return entities.NewErrCommand("cannot delete pod security policy", err), nil
		}
	}
	return entities.NewSuccessCommand([]byte("Pod security policy deleted")), nil
}

// String returns a string representation
func (dpsp *DeletePodSecurityPolicy) String() string {
	return fmt.Sprintf("SYNC DeletePodSecurityPolicy %s", dpsp.PolicyName)
}

// PrettyPrint returns a simple space indexed string.
func (dpsp *DeletePodSecurityPolicy) PrettyPrint(indentation int) string {
	return strings.Repeat(" ", indentation) + dpsp.String()
}

// UserString returns a simple string representation of the command for the user.
func (dpsp *DeletePodSecurityPolicy) UserString() string {
	return fmt.Sprintf("Deleting pod security policy %s", dpsp.PolicyName)
}
