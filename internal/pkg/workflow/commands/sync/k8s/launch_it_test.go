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
	"github.com/nalej/installer/internal/pkg/utils"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"path"
	"strings"
)

const SampleDevelopment = `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: NAME
  namespace: NAMESPACE
  labels:
    app: nginx
spec:
  replicas: 1
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.7.9
        ports:
        - containerPort: 80
`

func createDeployment(basePath string, namespace string, index int) {
	toWrite := strings.Replace(SampleDevelopment, "NAMESPACE", namespace, 1)
	toWrite = strings.Replace(toWrite, "NAME", fmt.Sprintf("nginx-%d", index), 1)
	outputPath := path.Join(basePath, fmt.Sprintf("component%d.yaml", index))
	err := ioutil.WriteFile(outputPath, []byte(toWrite), 777)
	gomega.Expect(err).To(gomega.Succeed())
	log.Debug().Str("file", outputPath).Msg("deployment has been created")
}

var _ = ginkgo.Describe("A launch command", func(){

	const numDeployments = 2

	if ! utils.RunIntegrationTests() {
		log.Warn().Msg("Integration tests are skipped")
		return
	}

	if itKubeConfigFile == "" {
		ginkgo.Fail("missing environment variables")
	}

	var componentsDir string

	ginkgo.BeforeSuite(func(){
		cd, err := ioutil.TempDir("", "launchIT")
		gomega.Expect(err).To(gomega.Succeed())
		componentsDir = cd

		for i:= 0; i< numDeployments; i++{
			createDeployment(componentsDir, itAuxNamespace, i)
		}
	})

	ginkgo.It("should create the deployments on kubernetes", func(){
	    lc := NewLaunchComponents(itKubeConfigFile, []string{itAuxNamespace}, componentsDir, "MINIKUBE")
	    result, err := lc.Run("testLaunchComponents")
	    gomega.Expect(err).To(gomega.Succeed())
	    gomega.Expect(result).ShouldNot(gomega.BeNil())
	    gomega.Expect(result.Success).Should(gomega.BeTrue())
	})

})