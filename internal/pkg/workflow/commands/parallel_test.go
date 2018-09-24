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

// Parallel command tests
//

package commands

/*
func TestParallelSync(t *testing.T) {
	utils.EnableDebug()
	cmd1 := sync.NewLogger("cmd1")
	cmd2 := sync.NewSleep("3")
	cmd3 := sync.NewLogger("cmd2")

	p := NewParallel("test synchronous commands", 3, []entities.Command{cmd1, cmd2, cmd3})
	wID := "testWorkflow"
	result, err := p.Run(wID)
	assert.Nil(t, err, "Parallel command should be executed")
	assert.True(t, result.Success, "Command should execute correctly")
}

func TestFromJSON(t *testing.T) {
	fromJSON := `
{"type":"sync", "name": "parallel", "maxParallelism":3, "commands": [
{"type":"sync", "name": "logger", "msg": "This is a logging message"},
{"type":"sync", "name": "logger", "msg": "This is a logging message"}]}
`
	received, err := NewParallelFromJSON([]byte(fromJSON))
	assert.Nil(t, err, "should not fail")
	assert.Equal(t, 3, (*received).(*Parallel).MaxParallelism, "should match")

	fromJSONWithoutParallelism := `
{"type":"sync", "name": "parallel", "commands": [
{"type":"sync", "name": "logger", "msg": "This is a logging message"},
{"type":"sync", "name": "logger", "msg": "This is a logging message"}]}
`
	received, err = NewParallelFromJSON([]byte(fromJSONWithoutParallelism))
	assert.Nil(t, err, "should not fail")
	assert.Equal(t, 0, (*received).(*Parallel).MaxParallelism, "should match")
}

func TestParallelStopOnFail(t *testing.T) {
	utils.EnableDebug()
	cmd1 := sync.NewFail()
	cmd2 := sync.NewSleep("10")
	cmd3 := sync.NewSleep("10")

	p := NewParallel("test synchronous commands", 3, []entities.Command{cmd1, cmd2, cmd3})
	wID := "testWorkflow"
	result, err := p.Run(wID)
	assert.Nil(t, err, "Parallel command should be executed")
	assert.False(t, result.Success, "Command should execute correctly")
}

func TestParallelWithLevel(t *testing.T) {
	utils.EnableDebug()
	cmd1 := sync.NewLogger("cmd1")
	cmd2 := sync.NewLogger("cmd2")
	cmd3 := sync.NewLogger("cmd3")
	cmd4 := sync.NewLogger("cmd4")
	cmd5 := sync.NewLogger("cmd5")
	cmd6 := sync.NewLogger("cmd6")
	cmd7 := sync.NewLogger("cmd7")
	cmd8 := sync.NewLogger("cmd8")

	p := NewParallel("test synchronous commands", 2,
		[]entities.Command{cmd1, cmd2, cmd3, cmd4, cmd5, cmd6, cmd7, cmd8})

	wID := "testWorkflow"
	result, err := p.Run(wID)
	assert.Nil(t, err, "Parallel command should be executed")
	assert.True(t, result.Success, "Command should execute correctly")
}

func TestParallelDP1164(t *testing.T) {
	utils.EnableDebug()
	cmd1 := sync.NewLogger("cmd1")
	cmd2 := sync.NewLogger("cmd2")

	p := NewParallel("test synchronous commands", 3,
		[]entities.Command{cmd1, cmd2})

	wID := "testWorkflow"
	result, err := p.Run(wID)
	assert.Nil(t, err, "Parallel command should be executed")
	assert.True(t, result.Success, "Command should execute correctly")
}
*/