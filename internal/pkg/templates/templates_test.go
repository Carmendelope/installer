/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
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
					params := workflow.GetTestParameters(numNodes, false)
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

		ginkgo.Context("installing an application cluster with coredns", func() {
			ginkgo.It("should be able to parse the template", func() {
				params := workflow.GetTestParameters(numNodes, true)
				params.InstallRequest.UseCoreDns = true
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
		ginkgo.Context("installing an application cluster with kubedns", func() {
			ginkgo.It("should be able to parse the template", func() {
				params := workflow.GetTestParameters(numNodes, true)
				params.InstallRequest.UseKubeDns = true
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
})
