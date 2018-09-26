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

 // RKE Install test
 // Prerequirements:
 // 1. Download a target VM
 //    $ vagrant init debian/stretch64
 //	1.1. Modify the Vagrantfile adding
 // config.vm.network "private_network", type: "dhcp"
 // 2. Launch the machine and obtain the private key
 //    $ vagrant up
 //    $ vagrant ssh-config
 // 3. Setup required software
 //    $ sudo apt-get install curl
 //    $ curl https://releases.rancher.com/install-docker/17.03.sh | sh
 //    $ sudo usermod -aG docker vagrant

/*
RUN_INTEGRATION_TEST=true
IT_RKE_TARGET_NODES=172.28.128.3
IT_RKE_PRIVATE_KEY=/private/tmp/it_test/.vagrant/machines/default/virtualbox/private_key
IT_RKE_BINARY=/Users/daniel/Downloads/rke_darwin-amd64
*/

package rke

import (
	"github.com/nalej/installer/internal/pkg/utils"
	"github.com/nalej/installer/internal/pkg/workflow/handler"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"os"
	"strings"

	"github.com/nalej/derrors"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
)

const TestTemplate string = ClusterTemplate

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

var _ = ginkgo.Describe("RKE Install", func(){
	if !utils.RunIntegrationTests() {
		return
	}

	var (
		privateKeyPath = os.Getenv("IT_RKE_PRIVATE_KEY")
		rkeBinaryPath = os.Getenv("IT_RKE_BINARY")
		targetNodes = os.Getenv("IT_RKE_TARGET_NODES")
	)

	if privateKeyPath == "" || rkeBinaryPath == "" || targetNodes == "" {
		ginkgo.Fail("missing environment variables")
	}

	ginkgo.It("should be able to install a Kubernetes cluster", func(){
	    nodes := strings.Split(targetNodes, ",")
		cmd := NewRKEInstall(
			rkeBinaryPath,
			*NewClusterConfig(
				"testClusterIT",
				nodes,
				"vagrant",
				privateKeyPath), "/tmp/", TestTemplate)
		commandHandler := handler.GetCommandHandler()
		gomega.Expect(commandHandler, gomega.Not(gomega.BeNil()))
		helper := NewHandlerHelper()
		err := commandHandler.AddCommand(cmd.ID(), helper.resultCallback, helper.logCallback)
		gomega.Expect(err).To(gomega.BeNil())
		result, err := cmd.Run("workflowID")
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(result).ToNot(gomega.BeNil())
		log.Debug().Msg(result.String())
		gomega.Expect(result.Success).To(gomega.BeTrue())
	})

})

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