/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

/*
RUN_INTEGRATION_TEST=true
IT_K8S_KUBECONFIG=/Users/daniel/.kube/config
*/

package k8s

import (
"github.com/nalej/installer/internal/pkg/utils"
"github.com/onsi/ginkgo"
"github.com/onsi/gomega"
"github.com/rs/zerolog/log"
)

var _ = ginkgo.Describe("A create management config command", func(){

	if ! utils.RunIntegrationTests() {
		log.Warn().Msg("Integration tests are skipped")
		return
	}

	if itKubeConfigFile == "" {
		ginkgo.Fail("missing environment variables")
	}

	ginkgo.It("should be able to create the config maps", func(){
		cmc := NewCreateManagementConfig(
			itKubeConfigFile, "publicHost", "publicPort",
			"MINIKUBE", "PRODUCTION")
		result, err := cmc.Run("createManagementConfig")
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(result.Success).Should(gomega.BeTrue())
	})

})
