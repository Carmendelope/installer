package k8s

import (
	"github.com/nalej/installer/internal/pkg/utils"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
)

/*
RUN_INTEGRATION_TEST=true
IT_K8S_KUBECONFIG=/Users/gaizka/.kube/config
*/

var _ = ginkgo.Describe("A Create TLS Secret command", func(){

	if ! utils.RunIntegrationTests() {
		log.Warn().Msg("Integration tests are skipped")
		return
	}

	if itKubeConfigFile == "" {
		ginkgo.Fail("missing environment variables")
	}

	testChecker := NewTestChecker(itKubeConfigFile)
	testChecker.Connect()

	ginkgo.It("should be able to create the secret", func(){
		// Create secret in Kubernetes
		cmd := NewCreateTLSSecret(itKubeConfigFile, "tls-secret", "cert", "AQAAAH", false, "")
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
