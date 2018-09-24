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

// Group command tests
//

package commands



/*
func TestBasicSequence(t *testing.T) {
	utils.EnableDebug()
	cmd1 := sync.NewLogger("cmd1")
	cmd2 := sync.NewSleep("1")
	cmd3 := sync.NewLogger("cmd2")
	g := NewGroup("basicSequence", []entities.Command{cmd1, cmd2, cmd3})
	wID := "TestBasicSequence"
	result, err := g.Run(wID)
	assert.Nil(t, err, "group command should be executed")
	assert.True(t, result.Success, "command should execute correctly")
}

func TestBasicAsync(t *testing.T) {
	utils.EnableDebug()
	cmd1 := sync.NewLogger("cmd1")
	cmd2 := async.NewSleep("1")
	cmd3 := sync.NewLogger("cmd2")
	g := NewGroup("basicSequence", []entities.Command{cmd1, cmd2, cmd3})
	wID := "TestBasicSequence"
	result, err := g.Run(wID)
	assert.Nil(t, err, "group command should be executed")
	assert.True(t, result.Success, "command should execute correctly")
}

func TestBasicSequenceFail(t *testing.T) {
	utils.EnableDebug()
	cmd1 := sync.NewLogger("cmd1")
	cmd2 := sync.NewFail()
	cmd3 := sync.NewSleep("1")
	cmd4 := sync.NewLogger("should not appear")
	g := NewGroup("basicSequence", []entities.Command{cmd1, cmd2, cmd3, cmd4})
	wID := "TestBasicSequenceFail"
	result, err := g.Run(wID)
	assert.Nil(t, err, "group command should be executed")
	assert.False(t, result.Success, "command should execute correctly")
}
*/