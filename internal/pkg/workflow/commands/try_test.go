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

package commands

import (
	"testing"

	"github.com/nalej/installer/internal/pkg/workflow/commands/async"
	"github.com/nalej/installer/internal/pkg/workflow/commands/sync"
	"github.com/stretchr/testify/assert"
)

func TestTrySync(t *testing.T) {
	utils.EnableDebug()
	cmd1 := sync.NewLogger("cmd1")
	cmd2 := sync.NewLogger("cmd2")
	try := NewTry("test sync", cmd1, cmd2)
	wID := "testWorkflow"
	result, err := try.Run(wID)
	assert.Nil(t, err, "Try command should be executed")
	assert.True(t, result.Success, "Command should execute correctly")
	assert.Equal(t, "cmd1", result.Output, "Result should match")
}

func TestTryFailSync(t *testing.T) {
	utils.EnableDebug()
	cmd1 := sync.NewFail()
	cmd2 := sync.NewLogger("cmd2")
	try := NewTry("test sync fail", cmd1, cmd2)
	wID := "testWorkflow"
	result, err := try.Run(wID)
	assert.Nil(t, err, "Try command should be executed")
	assert.True(t, result.Success, "Command should execute correctly")
	assert.Equal(t, "cmd2", result.Output, "Result should match")
}

func TestTryAsync(t *testing.T) {
	utils.EnableDebug()
	cmd1 := async.NewSleep("0")
	cmd2 := sync.NewLogger("cmd2")
	try := NewTry("test async", cmd1, cmd2)
	wID := "testWorkflow"
	result, err := try.Run(wID)
	assert.Nil(t, err, "Try command should be executed")
	assert.True(t, result.Success, "Command should execute correctly")
	assert.Equal(t, "Slept for 0", result.Output, "Result should match")
}

func TestTryFailAsync(t *testing.T) {
	utils.EnableDebug()
	cmd1 := async.NewFail()
	cmd2 := async.NewSleep("0")
	try := NewTry("test async", cmd1, cmd2)
	wID := "testWorkflow"
	result, err := try.Run(wID)
	assert.Nil(t, err, "Try command should be executed")
	assert.True(t, result.Success, "Command should execute correctly")
	assert.Equal(t, "Slept for 0", result.Output, "Result should match")
}

func TestNewTryFromJSON(t *testing.T) {
	fromJSON := `
{"type":"sync", "name": "try", "description":"Try",
"cmd": {"type":"sync", "name": "logger", "msg": "This is a logging message"},
"onFail": {"type":"sync", "name": "logger", "msg": "This is a logging message"}}
`
	received, err := NewTryFromJSON([]byte(fromJSON))
	assert.Nil(t, err, "should not fail")
	assert.NotNil(t, (*received).(*Try).TryCommand, "try command must be present")
	assert.NotNil(t, (*received).(*Try).OnFailCommand, "on fail command must be present")
}
