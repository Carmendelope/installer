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

const basicWorkflow = `
{
 "description": "basicWorkflow",
 "commands": [
  {"type":"sync", "name": "logger", "msg": "Starting basicWorkflow execution"},
  {"type":"sync", "name": "exec", "cmd": "mkdir", "args":["/tmp/basicWorkflow"]},
  {"type":"sync", "name": "exec", "cmd": "touch", "args":["/tmp/basicWorkflow/file"]},
  {"type":"sync", "name": "exec", "cmd": "ls", "args":["-las", "/tmp/basicWorkflow"]},
  {"type":"sync", "name": "exec", "cmd": "rm", "args":["/tmp/basicWorkflow/file"]},
  {"type":"sync", "name": "exec", "cmd": "rm", "args":["-r", "/tmp/basicWorkflow"]}
 ]
}
`

const basicParallelWorkflow = `
{
 "description": "basicParallelWorkflow",
 "commands": [
  {"type":"sync", "name":"parallel", "description":"Execute the following commands in parallel", 
    "commands":[
      {"type":"sync", "name": "logger", "msg": "Starting basicParallelWorkflow execution"},
      {"type":"sync", "name": "logger", "msg": "Ending basicParallelWorkflow execution"}
    ]}
  ]
}
`

const failParallelWorkflow = `
{
 "description": "failParallelWorkflow",
 "commands": [
  {"type":"sync", "name":"parallel", "description":"Execute the following commands in parallel", 
    "commands":[
      {"type":"sync", "name": "logger", "msg": "This commands works"},
      {"type":"async", "name": "fail"}
    ]}
  ]
}
`

const parallelMaxParallelismWorkflow = `
{
 "description": "basicParallelWorkflow", 
 "commands": [
  {"type":"sync", "name":"parallel", "description":"Execute the following commands in parallel", "maxParallelism":2,
    "commands":[
      {"type":"sync", "name": "logger", "msg": "msg1"},
      {"type":"sync", "name": "logger", "msg": "msg2"},
      {"type":"sync", "name": "logger", "msg": "msg3"},
      {"type":"sync", "name": "logger", "msg": "msg4"},
      {"type":"sync", "name": "logger", "msg": "msg5"},
      {"type":"sync", "name": "logger", "msg": "msg6"},
      {"type":"sync", "name": "logger", "msg": "msg7"},
      {"type":"sync", "name": "logger", "msg": "msg8"},
      {"type":"sync", "name": "logger", "msg": "msg9"},
      {"type":"sync", "name": "logger", "msg": "msg10"}
    ]}
  ]
}
`

/*
func getWorkflow(t *testing.T, name string, template string) *Workflow {
	p := NewParser()
	workflow, err := p.ParseWorkflow(template, name, EmptyParameters)
	assert.Nil(t, err, "expecting workflow")
	return workflow
}

func TestBasicWorkflow(t *testing.T) {
	maxWait := 5
	w := getWorkflow(t, "TestBasicWorkflow", basicWorkflow)
	wr := &WorkflowResult{}

	exec := NewWorkflowExecutor(w, wr.Callback)
	exec.Exec()
	// Wait for the workflow to finish
	for i := 0; i < maxWait && !wr.Finished(); i++ {
		time.Sleep(time.Second * 1)
	}
	fmt.Println(strings.Join(exec.Log(), "\n"))
	assert.True(t, wr.Called, "workflow should finish")
	assert.Nil(t, wr.Error, "workflow should not fail")
}

func TestBasicParallel(t *testing.T) {
	maxWait := 5
	w := getWorkflow(t, "TestBasicParallel", basicParallelWorkflow)
	wr := &WorkflowResult{}

	exec := NewWorkflowExecutor(w, wr.Callback)
	exec.Exec()
	// Wait for the workflow to finish
	for i := 0; i < maxWait && !wr.Finished(); i++ {
		time.Sleep(time.Second * 1)
	}
	fmt.Println(strings.Join(exec.Log(), "\n"))

	assert.True(t, wr.Called, "workflow should finish")
	assert.Nil(t, wr.Error, "workflow should not fail")
}

func TestFailParallel(t *testing.T) {
	maxWait := 5
	w := getWorkflow(t, "TestFailParallel", failParallelWorkflow)
	wr := &WorkflowResult{}

	exec := NewWorkflowExecutor(w, wr.Callback)
	exec.Exec()
	// Wait for the workflow to finish
	for i := 0; i < maxWait && !wr.Finished(); i++ {
		time.Sleep(time.Second * 1)
	}
	fmt.Println(strings.Join(exec.Log(), "\n"))

	assert.True(t, wr.Called, "workflow should finish")
	assert.NotNil(t, wr.Error, "workflow should not fail")
}

func TestMaxParallel(t *testing.T) {
	maxWait := 5
	w := getWorkflow(t, "TestMaxParallel", parallelMaxParallelismWorkflow)
	wr := &WorkflowResult{}

	exec := NewWorkflowExecutor(w, wr.Callback)
	exec.Exec()
	// Wait for the workflow to finish
	for i := 0; i < maxWait && !wr.Finished(); i++ {
		time.Sleep(time.Second * 1)
	}
	fmt.Println(strings.Join(exec.Log(), "\n"))

	assert.True(t, wr.Called, "workflow should finish")
	assert.Nil(t, wr.Error, "workflow should not fail")
}
*/