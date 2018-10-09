/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package templates

import (
	"github.com/nalej/installer/internal/pkg/workflow"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Templates", func() {

	const numNodes = 10
	var parser = workflow.NewParser()

	ginkgo.Context("Install Template", func() {

		ginkgo.Context("installing the management cluster", func() {
			ginkgo.It("should be able to parse the template", func(){
				params := workflow.GetTestParameters(numNodes, false)
				workflow, err := parser.ParseWorkflow("test", InstallManagementCluster, "InstallManagement", *params)
				gomega.Expect(err).To(gomega.Succeed())
				gomega.Expect(workflow).ShouldNot(gomega.BeNil())
			})
		})
	})
})
