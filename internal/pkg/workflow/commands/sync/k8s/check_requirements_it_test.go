/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

// Launch a simple test to deploy some components in Kubernetes
// Prerequirements
// 1.- Launch minikube

/*
RUN_INTEGRATION_TEST=true
IT_K8S_KUBECONFIG=/Users/daniel/.kube/config
 */

package k8s

import (
	"fmt"
	"github.com/nalej/grpc-utils/pkg/conversions"
	"github.com/nalej/installer/internal/pkg/utils"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"os"
)

var _ = ginkgo.Describe("A check requirements command", func(){

	if ! utils.RunIntegrationTests() {
		log.Warn().Msg("Integration tests are skipped")
		return
	}
	var (
		kubeConfigFile = os.Getenv("IT_K8S_KUBECONFIG")
	)

	if kubeConfigFile == "" {
		ginkgo.Fail("missing environment variables")
	}

	ginkgo.It("should pass the requirements on a common config", func(){
	    cr := NewCheckRequirements("1.9", kubeConfigFile)
		result, err := cr.Run("checkRequirements")
		if err != nil {
			fmt.Println(conversions.ToDerror(err))
		}
		fmt.Println(result.String())
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(result).ShouldNot(gomega.BeNil())
		gomega.Expect(result.Success).Should(gomega.BeTrue())
	})

	ginkgo.It("should fail on a non existing higher version", func(){

	})
})