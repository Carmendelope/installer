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

// TODO Define a proper IT. We need to consider that the target K8s is being cleaned up on other tests, so we need to install the required components.

/*
var _ = ginkgo.Describe("A create cluster user command", func(){

	if ! utils.RunIntegrationTests() {
		log.Warn().Msg("Integration tests are skipped")
		return
	}

	var (
		userManagerAddress = os.Getenv("IT_USER_MANAGER_ADDRESS")
	)

	if itKubeConfigFile == "" || userManagerAddress == "" {
		ginkgo.Fail("missing environment variables")
	}

	testChecker := NewTestChecker(itKubeConfigFile)
	testChecker.Connect()

	ginkgo.FIt("should be able to create the cluster user", func(){
		acu := NewAddClusterUser(itKubeConfigFile, "organizationID", "clusterID", userManagerAddress)
		result, err := acu.Run("addClusterUser")
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(result.Success).Should(gomega.BeTrue())
		secret := testChecker.GetSecret(ClusterUserSecretName, "nalej")
		gomega.Expect(len(secret.StringData)).Should(gomega.Equal(2))
	})

})
*/
