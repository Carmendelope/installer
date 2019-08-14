package k8s

import (
	"github.com/nalej/installer/internal/pkg/utils"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
)

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
		gomega.Expect(expectedPrivateKeyValue).Should(gomega.Equal(cmd.PrivateKeyValue))
		gomega.Expect(expectedCertValue).Should(gomega.Equal(cmd.CertValue))
	})
})
