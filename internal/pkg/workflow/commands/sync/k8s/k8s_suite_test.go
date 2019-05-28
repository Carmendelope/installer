/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package k8s

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"os"
	"testing"
)

const itAuxNamespace = "test-it-launch"
const itNalejNamespace = "nalej"

var itComponentsDir string
var itKubeConfigFile = os.Getenv("IT_K8S_KUBECONFIG")
var itTargetNamespaces = []string {itAuxNamespace, itNalejNamespace}
var itRegistryUsername = os.Getenv("IT_REGISTRY_USERNAME")
var itRegistryPassword = os.Getenv("IT_REGISTRY_PASSWORD")

func TestK8sPackage(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "K8s package suite")
}

var _ = ginkgo.AfterSuite(func(){
	log.Info().Msg("Cleaning test environment")
	if itComponentsDir != "" {
		os.RemoveAll(itComponentsDir)
	}
	if itKubeConfigFile != "" && len(itTargetNamespaces) > 0 {
		// for _, ns := range itTargetNamespaces{
		//	tc := NewTestCleaner(itKubeConfigFile, ns)
		//	gomega.Expect(tc.DeleteAll()).To(gomega.Succeed())
		//}
	}else{
		log.Warn().Msg("TestCleaner skipped")
	}
})