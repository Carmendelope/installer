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

package workflow

import (
	"github.com/nalej/installer/internal/pkg/workflow/commands/sync"
	"testing"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

const basicDefinitionNoTemplate = `
{
 "description": "basicDefinitionNoTemplate",
 // This is a comment
 "commands": [
  {"type":"sync", "name": "exec", "cmd": "cmd1"},
  {"type":"sync", "name": "exec", "cmd": "cmd2", "args":["arg1"]}
 ]
}
`

const basicDefinitionTwoCommands = `
{
 "description": "basicDefinitionNoTemplate",
 "commands": [
  {"type":"sync", "name": "exec", "cmd": "cmd1"},
  {"type":"sync", "name": "scp", "targetHost": "127.0.0.1", "credentials":{"username": "username", "password":"passwd", "privateKey":""}, "source":"script.sh", "destination":"/opt/scripts/."}
 ]
}
`

const basicDefinitionNodeTemplate = `
{
 "description": "basicDefinitionNodeTemplate",
 "commands": [
  {"type":"sync", "name": "exec", "cmd": "cmd1", "args":["{{.NetworkID}}", "{{.ClusterID}}"]}
  {{- range .Nodes.Slaves }}
  ,{"type":"sync", "name": "exec", "cmd": "cmd2", "args":["{{.PublicIP}}", "{{.Username}}", "{{.Password}}"]}
  {{- end}}
 ]
}
`

const basicDefinitionTestComma = `
{
 "description": "basicDefinitionNodeTemplate",
 "commands": [
  {"type":"sync", "name": "exec", "cmd": "cmd1", "args":["{{.NetworkID}}", "{{.ClusterID}}"]},
  {{ range $index, $node := .Nodes.Slaves }}
  {{if $index}},{{end}} {"type":"sync", "name": "exec", "cmd": "cmd2", "args":["{{$node.PublicIP}}", "{{$node.Username}}", "{{$node.Password}}"]}
  {{ end}}
 ]
}
`

func TestParser(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Parser test suite")
}

var _ = ginkgo.Describe("Parser", func() {
	var parser = NewParser()

	ginkgo.Context("parses a workflow not requiring the template", func(){
		workflow, err := parser.ParseWorkflow(basicDefinitionNoTemplate, "TestParseWorkflow_Basic", EmptyParameters)
		ginkgo.It("must be returned", func(){
			gomega.Expect(err, gomega.BeNil())
			gomega.Expect(workflow, gomega.Not(gomega.BeNil()))
		})
		ginkgo.It("must contain cmd1", func(){
			gomega.Expect(len((*workflow).Commands), gomega.Equal(2))
			cmd1 := (*workflow).Commands[0]
			gomega.Expect(cmd1.(*sync.Exec).Cmd, gomega.Equal("cmd1"))
		})
	})
})

/*
var parser = NewParser()

func getTestParameters(numNodes int) *Parameters {
	slaves := make([]smEntities.Node, 0)
	for i := 0; i < numNodes; i++ {
		toAdd := smEntities.NewNode(TestNetworkID, TestClusterID,
			"name"+strconv.Itoa(i), "desc"+strconv.Itoa(i), []string{},
			"public"+strconv.Itoa(i), "private"+strconv.Itoa(i),
			false, "user"+strconv.Itoa(i), "pass"+strconv.Itoa(i), "")
		slaves = append(slaves, *toAdd)
	}
	installCredentials := entities.NewInstallCredentials("username", "privateKeyPath")
	paths := NewPaths("assestPath", "binPath", "confPath")
	return NewParameters("", TestNetworkID, TestClusterID,
		*NewNodes(nil, make([]smEntities.Node, 0), make([]smEntities.Node, 0), slaves),
		*NewAssets(constants.SamuraiAssets, constants.SamuraiServiceNames), *paths, *installCredentials,
		"systemModelHost", "inframgrHost", false, false,
		*NewEmptyPhoneHomeParams(), *NewEmptyPostInstallParams(), "")
}

func TestParseWorkflow_Basic(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	workflow, err := parser.ParseWorkflow(basicDefinitionNoTemplate, "TestParseWorkflow_Basic", EmptyParameters)
	assert.Nil(t, err, "error should be nil")
	assert.NotNil(t, workflow, "workflow should be returned")
	assert.Equal(t, 2, len((*workflow).Commands))
	cmd1 := (*workflow).Commands[0]
	assert.Equal(t, "cmd1", cmd1.(*sync.Exec).Cmd, "command should match")
	// fmt.Println((*workflow).(* GenericWorkflow).PrettyPrint())
}

func TestParseWorkflow_SimpleTemplate(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	numNodes := 10
	params := getTestParameters(numNodes)
	workflow, err := parser.ParseWorkflow(basicDefinitionNodeTemplate, "TestParseWorkflow_SimpleTemplate", *params)
	assert.Nil(t, err, "error should be nil")
	assert.NotNil(t, workflow, "workflow should be returned")
	assert.Equal(t, numNodes+1, len((*workflow).Commands))
	fmt.Println((*workflow).PrettyPrint())
	cmd1 := (*workflow).Commands[0]
	assert.Equal(t, TestNetworkID, cmd1.(*sync.Exec).Args[0], "command should match")
	assert.Equal(t, TestClusterID, cmd1.(*sync.Exec).Args[1], "command should match")
	assert.Equal(t, "cmd1", cmd1.(*sync.Exec).Cmd, "command should match")
	for i := 0; i < numNodes; i++ {
		toCheck := (*workflow).Commands[i+1]
		assert.Equal(t, "public"+strconv.Itoa(i), toCheck.(*sync.Exec).Args[0], "command should match")
	}
	// fmt.Println((*workflow).(* GenericWorkflow).PrettyPrint())
}

func TestParseWorkflow_CommaConstruct(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	numNodes := 10
	params := getTestParameters(numNodes)
	workflow, err := parser.ParseWorkflow(basicDefinitionTestComma, "TestParseWorkflow_CommaConstruct", *params)
	assert.Nil(t, err, "error should be nil")
	assert.NotNil(t, workflow, "workflow should be returned")
	assert.Equal(t, numNodes+1, len((*workflow).Commands))
	fmt.Println((*workflow).PrettyPrint())
	cmd1 := (*workflow).Commands[0]
	assert.Equal(t, TestNetworkID, cmd1.(*sync.Exec).Args[0], "command should match")
	assert.Equal(t, TestClusterID, cmd1.(*sync.Exec).Args[1], "command should match")
	assert.Equal(t, "cmd1", cmd1.(*sync.Exec).Cmd, "command should match")
	for i := 0; i < numNodes; i++ {
		toCheck := (*workflow).Commands[i+1]
		assert.Equal(t, "public"+strconv.Itoa(i), toCheck.(*sync.Exec).Args[0], "command should match")
	}
	// fmt.Println((*workflow).(* GenericWorkflow).PrettyPrint())
}

func TestParseWorkflow_TwoCommands(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	workflow, err := parser.ParseWorkflow(basicDefinitionTwoCommands, "TestParseWorkflow_TwoCommands", EmptyParameters)
	assert.Nil(t, err, "error should be nil")
	assert.NotNil(t, workflow, "workflow should be returned")
	assert.Equal(t, 2, len((*workflow).Commands))
	fmt.Println((*workflow).PrettyPrint())
	cmd1 := (*workflow).Commands[0]
	assert.Equal(t, entities.Exec, cmd1.Name(), "name should match")
	assert.Equal(t, "cmd1", cmd1.(*sync.Exec).Cmd, "command should match")
	cmd2 := (*workflow).Commands[1]
	assert.Equal(t, entities.SCP, cmd2.Name(), "name should match")
	assert.Equal(t, "127.0.0.1", cmd2.(*sync.SCP).TargetHost, "command should match")
}
*/