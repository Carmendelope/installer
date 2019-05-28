/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package k8s

/*
RUN_INTEGRATION_TEST=true
IT_K8S_KUBECONFIG=/Users/daniel/.kube/config
*/

import (
"github.com/nalej/installer/internal/pkg/utils"
"github.com/onsi/ginkgo"
"github.com/onsi/gomega"
"github.com/rs/zerolog/log"
)

var _ = ginkgo.Describe("A create CA certificate command", func() {

	if !utils.RunIntegrationTests() {
		log.Warn().Msg("Integration tests are skipped")
		return
	}

	if itKubeConfigFile == "" {
		ginkgo.Fail("missing environment variables")
	}

	ginkgo.FIt("should be able to create the CA certificate", func() {
		cc := NewCreateCACert(
			itKubeConfigFile, "nalej39.nalej.tech")
		result, err := cc.Run("createCACert")
		if err != nil{
			log.Error().Str("trace", err.DebugReport()).Msg("failed")
		}
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(result.Success).Should(gomega.BeTrue())
	})

})
