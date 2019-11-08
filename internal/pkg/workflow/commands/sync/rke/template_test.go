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

var _ = ginkgo.Describe("Template", func() {
	ginkgo.It("Should work with a single node", func() {
		config := getClusterConfig(1)
		template := NewRKETemplate(ClusterTemplate)
		yamlString, err := template.ParseTemplate(config)
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(yamlString).ToNot(gomega.BeNil())
		log.Debug().Msg(yamlString)
		err = template.ValidateYAML(yamlString)
		gomega.Expect(err).To(gomega.BeNil())
	})

	ginkgo.PIt("Should work with 2 nodes", func() {

	})

	ginkgo.PIt("Should work with 3 nodes", func() {

	})

	ginkgo.It("Should work with 10 nodes", func() {
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
