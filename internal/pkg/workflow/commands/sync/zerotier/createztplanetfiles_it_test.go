/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

 package zerotier

import (
	"github.com/nalej/installer/internal/pkg/utils"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os"
	"path/filepath"
)

/*
RUN_INTEGRATION_TEST=true
IT_K8S_KUBECONFIG=/Users/gaizka/.kube/config
IT_ZT_BINARY=/Users/gaizka/development/zerotier-idtool
*/

func createTempFilePath(name string) string{
	dir, err := ioutil.TempDir("", "ztplanets")
	gomega.Expect(err).Should(gomega.Succeed())
	return filepath.Join(dir, name)
}

var _ = ginkgo.Describe("A Create ZT Planet Files command", func(){

	if ! utils.RunIntegrationTests() {
		log.Warn().Msg("Integration tests are skipped")
		return
	}

	ztBinaryPath := os.Getenv("IT_ZT_BINARY")

	if itKubeConfigFile == "" || ztBinaryPath == "" {
		ginkgo.Fail("missing environment variables")
	}

	ginkgo.It("should be able to update the config map", func(){
		cmd := NewCreateZTPlanetFiles(itKubeConfigFile, ztBinaryPath,
			"managementPublicHost",
			createTempFilePath("identitySecret"),
			createTempFilePath("identityPublic"),
			createTempFilePath("planetJson"),
			createTempFilePath("planet"))
		result, err := cmd.Run("createZtPlanetFiles")
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(result.Success).Should(gomega.BeTrue())
	})

})