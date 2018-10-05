/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package rke

import (
	"fmt"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
)

func getClusterConfig(numNodes int) *ClusterConfig {
	targetNodes := make([]string, 0)
	for i := 0; i < numNodes; i++ {
		node := fmt.Sprintf("172.1.1.%d", i)
		targetNodes = append(targetNodes, node)
	}
	return NewClusterConfig(
		"testClusterName",
		targetNodes,
		"nodeUsername",
		"privateKeyPath")
}

var _ = ginkgo.Describe("Template", func(){
	ginkgo.It("Should work with a single node", func(){
		config := getClusterConfig(1)
		template := NewRKETemplate(ClusterTemplate)
		yamlString, err := template.ParseTemplate(config)
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(yamlString).ToNot(gomega.BeNil())
		log.Debug().Msg(yamlString)
		err = template.ValidateYAML(yamlString)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.PIt("Should work with 2 nodes", func(){

	})

	ginkgo.PIt("Should work with 3 nodes", func(){

	})

	ginkgo.It("Should work with 10 nodes", func(){
		config := getClusterConfig(10)
		template := NewRKETemplate(ClusterTemplate)
		yamlString, err := template.ParseTemplate(config)
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(yamlString).ToNot(gomega.BeNil())
		log.Debug().Msg(yamlString)
		fmt.Println(yamlString)
		err = template.ValidateYAML(yamlString)
		gomega.Expect(err).To(gomega.BeNil())
	})
})

