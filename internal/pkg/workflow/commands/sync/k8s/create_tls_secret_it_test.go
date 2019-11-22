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
)

var _ = ginkgo.Describe("A Create TLS Secret command", func() {

	if !utils.RunIntegrationTests() {
		log.Warn().Msg("Integration tests are skipped")
		return
	}

	if !utils.RunIntegrationTest("create_tls_secret_it") {
		log.Warn().Msg("Integration test is skipped")
		return
	}

	if itKubeConfigFile == "" {
		ginkgo.Fail("missing environment variables")
	}

	testChecker := NewTestChecker(itKubeConfigFile)
	testChecker.Connect()

	ginkgo.It("should be able to create the secret", func() {
		// Create secret in Kubernetes
		cmd := NewCreateTLSSecret(itKubeConfigFile, "tls-secret", "", "AQAAAH")
		result, err := cmd.Run("createTLSSecret")
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		// Retrieve secret from kubernetes
		retrieved := testChecker.GetSecret(cmd.SecretName, "nalej")
		gomega.Expect(len(retrieved.Data)).Should(gomega.Equal(2))
		secretContent := retrieved.Data
		expectedPrivateKeyValue := string(secretContent["tls.key"])
		expectedCertValue := string(secretContent["tls.crt"])
		gomega.Expect(expectedPrivateKeyValue).Should(gomega.Equal(cmd.PrivateKeyPath))
		gomega.Expect(expectedCertValue).Should(gomega.Equal(cmd.CertPath))
	})
})
