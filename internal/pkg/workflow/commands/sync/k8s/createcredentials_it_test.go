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

var _ = ginkgo.Describe("A create credentials command", func(){

	if ! utils.RunIntegrationTests() {
		log.Warn().Msg("Integration tests are skipped")
		return
	}

	if itKubeConfigFile == "" || itRegistryUsername == "" || itRegistryPassword == "" {
		ginkgo.Fail("missing environment variables")
	}

	ginkgo.It("should be able to create the config", func(){
		uc := NewCreateCredentials(itKubeConfigFile, itRegistryUsername, itRegistryPassword)
		result, err := uc.Run("createCredentials")
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(result.Success).Should(gomega.BeTrue())
	})

})
