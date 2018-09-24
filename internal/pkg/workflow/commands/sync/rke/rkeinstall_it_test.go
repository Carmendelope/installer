/*
 * Copyright 2018 Nalej
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
 */

package rke

import (
	"github.com/rs/zerolog/log"
	"os"

	"github.com/nalej/derrors"
	"github.com/nalej/installer/internal/pkg/workflow/entities"

)

const TestTemplate string = ClusterTemplate

func RunRKEInstallTest() bool {
	var runIntegration = os.Getenv("RUN_INTEGRATION_TEST")
	var privateKeyPath = os.Getenv("RKE_PRIVATE_KEY")
	var rkeBinaryPath = os.Getenv("RKE_BINARY")
	var target = os.Getenv("RKE_TARGET_NODES")
	return runIntegration == "true" && privateKeyPath != "" && rkeBinaryPath != "" && target != ""
}

type HandlerHelper struct {
}

func NewHandlerHelper() *HandlerHelper {
	return &HandlerHelper{}
}

func (hh *HandlerHelper) resultCallback(id string, result *entities.CommandResult, error derrors.Error) {
	log.Debug().Str("id", id).Str("result", result.String()).Msg("resultCallback()")
}

func (hh *HandlerHelper) logCallback(id string, logEntry string) {
	log.Debug().Str("id", id).Str("logEntry", logEntry).Msg("logCallback()")
}

/*
func installVagrant(t *testing.T) {

	var rkeBinaryPath = os.Getenv("RKE_BINARY")
	var privateKeyPath = os.Getenv("RKE_PRIVATE_KEY")
	var target = os.Getenv("RKE_TARGET_NODES")
	targetNodes := []string{target}

	cmd := NewRKEInstall(
		rkeBinaryPath,
		*NewClusterConfig(
			"testClusterIT",
			targetNodes[0],
			targetNodes,
			targetNodes,
			"vagrant",
			privateKeyPath), "/tmp/", TestTemplate)

	commandHandler := handler.GetCommandHandler()
	helper := NewHandlerHelper()
	commandHandler.AddCommand(cmd.ID(), helper.resultCallback, helper.logCallback)

	result, err := cmd.Run("workflowID")
	assert.Nil(t, err, "expecting no error")
	assert.NotNil(t, result, "expecting result")
	log.Debug().Msg(result.String())
	assert.True(t, result.Success, "expecting command to succedded")
}

func TestRKEInstall(t *testing.T) {
	if RunRKEInstallTest() {
		utils.EnableDebug()
		installVagrant(t)
	} else {
		log.Info().Msg("skipping RKE install test")
	}
}
*/