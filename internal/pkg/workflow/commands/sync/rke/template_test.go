/*
 * Copyright 2018 Nalej
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
 */

package rke

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getClusterConfig(numNodes int) *ClusterConfig {
	targetNodes := make([]string, 0)
	for i := 0; i < numNodes; i++ {
		node := fmt.Sprintf("172.1.1.%d", i)
		targetNodes = append(targetNodes, node)
	}
	m2Nodes := make([]string, 0)
	for i := 0; i < 2 && i < numNodes; i++ {
		m2Nodes = append(m2Nodes, targetNodes[i])
	}
	return NewClusterConfig(
		"testClusterName",
		"SamuraiIP",
		m2Nodes,
		targetNodes,
		"nodeUsername",
		"privateKeyPath")
}

func TestSingleNode(t *testing.T) {
	config := getClusterConfig(1)
	template := NewRKETemplate(ClusterTemplate)
	yamlString, err := template.ParseTemplate(config)
	assert.Nil(t, err, "expecting yaml file")
	assert.NotEmpty(t, yamlString, "expecting yaml file")
	fmt.Println(yamlString)
	err = template.ValidateYAML(yamlString)
	assert.Nil(t, err, "yaml should be valid")
}

func Test10Nodes(t *testing.T) {
	config := getClusterConfig(10)
	template := NewRKETemplate(ClusterTemplate)
	yamlString, err := template.ParseTemplate(config)
	assert.Nil(t, err, "expecting yaml file")
	assert.NotEmpty(t, yamlString, "expecting yaml file")
	fmt.Println(yamlString)
	err = template.ValidateYAML(yamlString)
	assert.Nil(t, err, "yaml should be valid")
}
