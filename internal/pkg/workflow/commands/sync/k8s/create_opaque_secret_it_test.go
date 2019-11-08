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
	"github.com/nalej/installer/internal/pkg/utils"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"path/filepath"
)

/*
RUN_INTEGRATION_TEST=true
IT_K8S_KUBECONFIG=/Users/gaizka/.kube/config
*/

func createTempFilePath(name string) string {
	dir, err := ioutil.TempDir("", "ztplanets")
	gomega.Expect(err).Should(gomega.Succeed())
	return filepath.Join(dir, name)
}

var _ = ginkgo.Describe("A Create Opaque Secret command", func() {

	if !utils.RunIntegrationTests() {
		log.Warn().Msg("Integration tests are skipped")
		return
	}

	if itKubeConfigFile == "" {
		ginkgo.Fail("missing environment variables")
	}

	testChecker := NewTestChecker(itKubeConfigFile)
	testChecker.Connect()

	ginkgo.It("should be able to create the secret", func() {
		// Create secret in Kubernetes
		cmd := NewCreateOpaqueSecret(itKubeConfigFile, "zt-planet", "planet", "AQAAAH", false, "")
		result, err := cmd.Run("createZtPlanetFiles")
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		// Retrieve secret from kubernetes
		retrieved := testChecker.GetSecret(cmd.SecretName, "nalej")
		gomega.Expect(len(retrieved.Data)).Should(gomega.Equal(1))
		secretContent := retrieved.Data
		gomega.Expect(secretContent[cmd.SecretKey]).Should(gomega.Equal(cmd.SecretValue))
	})
})
