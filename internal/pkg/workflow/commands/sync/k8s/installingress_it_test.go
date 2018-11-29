/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package k8s

import (
	"github.com/nalej/installer/internal/pkg/utils"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
)

var _ = ginkgo.Describe("A install ingress command", func(){

	if ! utils.RunIntegrationTests() {
		log.Warn().Msg("Integration tests are skipped")
		return
	}

	if itKubeConfigFile == "" {
		ginkgo.Fail("missing environment variables")
	}

	ginkgo.FIt("should be able to install an ingress", func(){
		ii := NewInstallIngress(itKubeConfigFile)
		result, err := ii.Run("installIngress")
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(result.Success).Should(gomega.BeTrue())
	})

})