/*
 * Copyright 2019 Nalej
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Launch a simple test to deploy some components in Kubernetes
// Prerequirements
// 1.- Launch minikube

/*
RUN_INTEGRATION_TEST=true
IT_K8S_KUBECONFIG=/Users/daniel/.kube/config
IT_RKE_BINARY=/Users/daniel/development/rke/rke
*/

package installer

import (
	"context"
	"fmt"
	grpc_common_go "github.com/nalej/grpc-common-go"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-installer-go"
	"github.com/nalej/grpc-utils/pkg/test"
	cfg "github.com/nalej/installer/internal/pkg/server/config"
	"github.com/nalej/installer/internal/pkg/utils"
	"github.com/nalej/installer/internal/pkg/workflow/commands/sync/k8s"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

const SampleComponent = `
apiVersion: apps/v1
kind: Deployment
metadata:
  cluster: application
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
	toWrite := strings.Replace(SampleComponent, "NAMESPACE", namespace, 1)
	toWrite = strings.Replace(toWrite, "NAME", fmt.Sprintf("nginx-%d", index), 1)
	outputPath := path.Join(basePath, fmt.Sprintf("component%d.yaml", index))
	err := ioutil.WriteFile(outputPath, []byte(toWrite), 777)
	gomega.Expect(err).To(gomega.Succeed())
	log.Debug().Str("file", outputPath).Msg("deployment has been created")
}

var _ = ginkgo.Describe("Installer", func() {

	const numDeployments = 2
	const targetNamespace = "test-it-install"

	if !utils.RunIntegrationTests() {
		log.Warn().Msg("Integration tests are skipped")
		return
	}
	var (
		kubeConfigFile = os.Getenv("IT_K8S_KUBECONFIG")
		rkeBinary      = os.Getenv("IT_RKE_BINARY")
	)

	if kubeConfigFile == "" || rkeBinary == "" {
		ginkgo.Fail("missing environment variables")
	}

	var componentsDir string
	var binaryDir string
	var tempDir string
	var kubeConfigRaw string

	// gRPC server
	var server *grpc.Server
	// grpc test listener
	var listener *bufconn.Listener
	// client
	var client grpc_installer_go.InstallerClient

	ginkgo.BeforeSuite(func() {

		// Load data and ENV variables.
		kubeConfigContent, lErr := utils.GetKubeConfigContent(kubeConfigFile)
		gomega.Expect(lErr).To(gomega.Succeed())
		kubeConfigRaw = kubeConfigContent

		cd, err := ioutil.TempDir("", "installITComponents")
		gomega.Expect(err).To(gomega.Succeed())
		componentsDir = cd

		td, err := ioutil.TempDir("", "installITTemp")
		gomega.Expect(err).To(gomega.Succeed())
		tempDir = td

		binaryDir = filepath.Dir(rkeBinary)

		for i := 0; i < numDeployments; i++ {
			createDeployment(componentsDir, targetNamespace, i)
		}

		tu := k8s.NewTestK8sUtils(kubeConfigFile)
		gomega.Expect(tu.Connect()).To(gomega.Succeed())
		gomega.Expect(tu.CreateNamespace(targetNamespace)).To(gomega.Succeed())

		config := cfg.Config{
			ComponentsPath: componentsDir,
			BinaryPath:     binaryDir,
			TempPath:       tempDir,
		}

		// Launch gRPC server
		listener = test.GetDefaultListener()

		server = grpc.NewServer()

		manager := NewManager(config)
		handler := NewHandler(manager)
		grpc_installer_go.RegisterInstallerServer(server, handler)

		test.LaunchServer(server, listener)

		conn, err := test.GetConn(*listener)
		gomega.Expect(err).To(gomega.Succeed())
		client = grpc_installer_go.NewInstallerClient(conn)

	})

	ginkgo.AfterSuite(func() {
		server.Stop()
		listener.Close()
		os.RemoveAll(componentsDir)
		tc := k8s.NewTestCleaner(kubeConfigFile, targetNamespace)
		gomega.Expect(tc.DeleteAll()).To(gomega.Succeed())
	})

	ginkgo.PContext("On a base system", func() {
		ginkgo.PIt("should be able to install an application cluster from scratch", func() {

		})
	})

	ginkgo.Context("On a kubernetes cluster", func() {
		ginkgo.It("should be able to install an application cluster", func() {
			ginkgo.By("installing the cluster")
			installRequest := &grpc_installer_go.InstallRequest{
				RequestId:         "test-install-id",
				OrganizationId:    "test-org-id",
				ClusterId:         "test-cluster-id",
				ClusterType:       grpc_infrastructure_go.ClusterType_KUBERNETES,
				InstallBaseSystem: false,
				KubeConfigRaw:     kubeConfigRaw,
			}
			response, err := client.InstallCluster(context.Background(), installRequest)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(response).ToNot(gomega.BeNil())
			gomega.Expect(response.RequestId).Should(gomega.Equal(installRequest.RequestId))

			// Wait for it to finish
			maxWait := 1000
			finished := false

			requestID := &grpc_common_go.RequestId{
				RequestId: installRequest.RequestId,
			}
			ginkgo.By("checking the install progress")
			log.Info().Msg("Checking progress")
			for i := 0; i < maxWait && !finished; i++ {
				time.Sleep(time.Second)
				progress, err := client.CheckProgress(context.Background(), requestID)
				gomega.Expect(err).To(gomega.Succeed())
				log.Debug().Interface("progress", progress).Msg("Check progress")
				finished = (progress.Status == grpc_common_go.OpStatus_SUCCESS) ||
					(progress.Status == grpc_common_go.OpStatus_FAILED)
				log.Debug().Bool("finished", finished).Msg("workflow has finished?")
			}
			log.Info().Msg("obtain final progress")
			progress, err := client.CheckProgress(context.Background(), requestID)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(progress.Status).Should(gomega.Equal(grpc_common_go.OpStatus_SUCCESS))
			ginkgo.By("removing the install")
			log.Info().Msg("removing the install")

			client.RemoveInstall(context.Background(), requestID)
			log.Info().Msg("Finished!!!!!!")
		})
	})

})
