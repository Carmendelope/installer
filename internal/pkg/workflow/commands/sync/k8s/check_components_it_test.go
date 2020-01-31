/*
 * Copyright 2020 Nalej
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

// Launch a simple test to deploy some components in Kubernetes
// Requirements
// 1.- Launch minikube

/*
RUN_INTEGRATION_TEST=true
IT_K8S_KUBECONFIG=/Users/gaizka/.kube/config
*/

package k8s

import (
	"github.com/nalej/installer/internal/pkg/utils"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"io/ioutil"
)

var _ = ginkgo.Describe("A check components command", func() {

	const numDeployments = 2

	if !utils.RunIntegrationTests() {
		log.Warn().Msg("Integration tests are skipped")
		return
	}

	if !utils.RunIntegrationTest("launch_it_test") {
		log.Warn().Msg("Integration test is skipped")
		return
	}

	if itKubeConfigFile == "" {
		ginkgo.Fail("missing environment variables")
	}

	var componentsDir string

	ginkgo.BeforeSuite(func() {
		cd, err := ioutil.TempDir("", "checkComponentsIT")
		gomega.Expect(err).To(gomega.Succeed())
		componentsDir = cd

		for i := 0; i < numDeployments; i++ {
			createDeployment(componentsDir, itAuxNamespace, i)
		}
	})

	ginkgo.It("should create the deployments on kubernetes", func() {
		cc := NewCheckComponents(itKubeConfigFile, []string{itAuxNamespace})
		result, err := cc.Run("testCheckComponents")
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(result).ShouldNot(gomega.BeNil())
		gomega.Expect(result.Success).Should(gomega.BeTrue())
	})

})