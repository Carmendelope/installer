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

package entities

/*
func TestCommandResult_FromToJSON(t *testing.T) {
	de := derrors.NewGenericError("some error")
	cr := NewCommandResult(true, "output", de)
	result, err := json.Marshal(cr)
	fmt.Println(string(result))
	assert.Nil(t, err, "marshal should not fail")
	retrieved := &CommandResultFromJSON{}
	err = json.Unmarshal(result, retrieved)
	assert.Nil(t, err, "unmarshall should work")
	toCR := retrieved.ToCommandResult()
	assert.EqualValues(t, cr, toCR, "commands should match")
}

func TestCommandResult_JSONString(t *testing.T) {

	toReceiveNoError := `
    {"success":true, "output":"output"}
`
	retrieved := &CommandResultFromJSON{}
	err := json.Unmarshal([]byte(toReceiveNoError), retrieved)
	assert.Nil(t, err, "unmarshall should work")
	assert.True(t, retrieved.Success)
	assert.Equal(t, "output", retrieved.Output)
	assert.Nil(t, retrieved.Error)
	toCR := retrieved.ToCommandResult()
	assert.True(t, toCR.Success, "should report success")
	assert.Equal(t, "output", toCR.Output, "output should match")
	assert.Nil(t, toCR.Error, "error should be nil")
}

*/