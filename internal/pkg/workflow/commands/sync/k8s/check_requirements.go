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
	"strconv"
	"strings"
)

type CheckRequirements struct {
	Kubernetes
	MinVersion string `json:"minVersion"`
}

func NewCheckRequirements(minVersion string, kubeConfigPath string) *CheckRequirements {
	return &CheckRequirements{
		Kubernetes: Kubernetes{
			GenericSyncCommand: *entities.NewSyncCommand(entities.CheckRequirements),
			KubeConfigPath:     kubeConfigPath,
		},
		MinVersion: minVersion,
	}
}

// NewCheckRequirementsFromJSON creates an CheckRequirements command from a JSON object.
func NewCheckRequirementsFromJSON(raw []byte) (*entities.Command, derrors.Error) {
	cr := &CheckRequirements{}
	if err := json.Unmarshal(raw, &cr); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	cr.CommandID = entities.GenerateCommandID(cr.Name())
	var r entities.Command = cr
	return &r, nil
}

func (cr *CheckRequirements) CheckVersion(major string, minor string) bool {
	if cr.MinVersion == "" {
		return false
	}

	vSplit := strings.Split(cr.MinVersion, ".")
	if len(vSplit) < 2 {
		return false
	}

	majorServer, emas := strconv.Atoi(major)
	minorServer, emis := strconv.Atoi(minor)

	majorRequired, emar := strconv.Atoi(vSplit[0])
	minorRequired, emir := strconv.Atoi(vSplit[1])

	if emas != nil || emis != nil || emar != nil || emir != nil {
		log.Warn().Str("required", cr.MinVersion).Str("major", major).Str("minor", minor).
			Msg("Cannot parse version")
		return false
	}

	return (majorServer >= majorRequired) && (minorServer >= minorRequired)
}

func (cr *CheckRequirements) Run(workflowID string) (*entities.CommandResult, derrors.Error) {

	connectErr := cr.Connect()
	if connectErr != nil {
		return nil, connectErr
	}
	// Check the server version.
	sv, err := cr.Client.Discovery().ServerVersion()

	if err != nil {
		return nil, derrors.NewInternalError("cannot connect to K8s", err)
	}

	log.Debug().Str("version", sv.String()).Msg("Server")
	if !cr.CheckVersion(sv.Major, sv.Minor) {
		msg := fmt.Sprintf("expecting %s, found %s.%s", cr.MinVersion, sv.Major, sv.Minor)
		return entities.NewCommandResult(false, msg, nil), nil
	}
	msg := fmt.Sprintf("Version OK, found %s.%s", sv.Major, sv.Minor)
	return entities.NewSuccessCommand([]byte(msg)), nil
}

func (cr *CheckRequirements) String() string {
	return "SYNC CheckRequirements minVersion: " + cr.MinVersion
}

func (cr *CheckRequirements) PrettyPrint(indentation int) string {
	return strings.Repeat(" ", indentation) + cr.String()
}

func (cr *CheckRequirements) UserString() string {
	return fmt.Sprintf("Checking Kubernetes requirements")
}
