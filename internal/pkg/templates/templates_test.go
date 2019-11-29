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

package templates

import (
	"github.com/nalej/grpc-installer-go"
	"github.com/nalej/installer/internal/pkg/workflow"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Templates", func() {

	const numNodes = 10
	var parser = workflow.NewParser()

	availablePlatforms := []grpc_installer_go.Platform{
		grpc_installer_go.Platform_MINIKUBE,
		grpc_installer_go.Platform_AZURE,
		grpc_installer_go.Platform_BAREMETAL,
	}

	ginkgo.Context("Install Template", func() {

		for _, platformType := range availablePlatforms {
			ginkgo.Context("installing the management cluster", func() {
				ginkgo.It("should be able to parse the template", func() {
					params := workflow.GetTestInstallParameters(numNodes, false)
					params.InstallRequest.TargetPlatform = platformType
					params.InstallRequest.StaticIpAddresses = &grpc_installer_go.StaticIPAddresses{
						UseStaticIp: false,
						Ingress:     "",
						Dns:         "",
					}
					workflow, err := parser.ParseWorkflow("test", InstallManagementCluster, "InstallManagement", *params)
					gomega.Expect(err).To(gomega.Succeed())
					gomega.Expect(workflow).ShouldNot(gomega.BeNil())
				})
			})
		}

		ginkgo.Context("installing an application cluster", func() {
			ginkgo.It("should be able to parse the template", func() {
				params := workflow.GetTestInstallParameters(numNodes, true)
				params.InstallRequest.StaticIpAddresses = &grpc_installer_go.StaticIPAddresses{
					UseStaticIp: false,
					Ingress:     "",
					Dns:         "",
				}
				workflow, err := parser.ParseWorkflow("test", InstallManagementCluster, "InstallAppCluster", *params)
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(workflow).ShouldNot(gomega.BeNil())
			})
		})

	})

	ginkgo.Context("Uninstall template", func() {
		ginkgo.It("should uninstall a management cluster", func() {
			params := workflow.GetTestUninstallParameters(false)
			workflow, err := parser.ParseWorkflow("test", UninstallCluster, "UninstallManagement", *params)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(workflow).ShouldNot(gomega.BeNil())
		})
		ginkgo.It("should uninstall an application cluster", func() {
			params := workflow.GetTestUninstallParameters(true)
			workflow, err := parser.ParseWorkflow("test", UninstallCluster, "UninstallManagement", *params)
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(workflow).ShouldNot(gomega.BeNil())
		})
	})
})
