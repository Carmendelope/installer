/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
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
