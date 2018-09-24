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

package handler

/*
type HandlerTestSuite struct {
	suite.Suite
	handler *commandHandler
}

func (suite *HandlerTestSuite) SetupTest() {
	suite.handler = NewCommandHandler().(*commandHandler)
}

func (suite *HandlerTestSuite) TestSimpleWorkflow() {
	suite.Equal(0, len(suite.handler.resultCallbacks), "resultCallbacks map must be empty")
	suite.Equal(0, len(suite.handler.logCallbacks), "logCallbacks map must be empty")
	lines := 0
	finalized := false

	err := suite.handler.AddCommand("id1",
		func(id string, result *entities.CommandResult, error derrors.Error) {
			finalized = true
		},
		func(id string, logEntry string) {
			lines++
		},
	)
	suite.Nil(err, "command must be added")
	suite.Equal(1, len(suite.handler.resultCallbacks), "resultCallbacks has one element")
	suite.Equal(1, len(suite.handler.logCallbacks), "logCallbacks map has one element")
	suite.handler.AddLogEntry("id1", "hello world!")
	suite.handler.FinishCommand("id1", entities.NewSuccessCommand([]byte("OK")), nil)
	time.Sleep(time.Second)
	suite.Equal(1, lines, "must receive one log entry")
	suite.True(finalized, "must finalize")
	suite.Equal(0, len(suite.handler.resultCallbacks), "resultCallbacks map must be empty")
	suite.Equal(0, len(suite.handler.logCallbacks), "logCallbacks map must be empty")
}

func (suite *HandlerTestSuite) TestDuplicatedCommand() {
	err := suite.handler.AddCommand("id1",
		func(id string, result *entities.CommandResult, error derrors.Error) {

		},
		func(id string, logEntry string) {

		},
	)
	suite.Nil(err, "command must be added")
	err = suite.handler.AddCommand("id1",
		func(id string, result *entities.CommandResult, error derrors.Error) {

		},
		func(id string, logEntry string) {

		},
	)
	suite.NotNil(err, "AddCommand must fail")
}
func (suite *HandlerTestSuite) TestTwoCallbacks() {
	lines1 := 0
	finalized1 := false
	err := suite.handler.AddCommand("id1",
		func(id string, result *entities.CommandResult, error derrors.Error) {
			finalized1 = true
		},
		func(id string, logEntry string) {
			lines1++
		},
	)
	suite.Nil(err, "command must be added")
	lines2 := 0
	finalized2 := false
	err = suite.handler.AddCommand("id2",
		func(id string, result *entities.CommandResult, error derrors.Error) {
			finalized2 = true
		},
		func(id string, logEntry string) {
			lines2++
		},
	)
	suite.Nil(err, "command must be added")

	suite.Nil(err, "command must be added")
	suite.Equal(2, len(suite.handler.resultCallbacks), "resultCallbacks has two elements")
	suite.Equal(2, len(suite.handler.logCallbacks), "logCallbacks map has two elements")
	suite.handler.AddLogEntry("id1", "hello world!")
	suite.handler.FinishCommand("id1", entities.NewSuccessCommand([]byte("OK")), nil)
	time.Sleep(time.Second)
	suite.Equal(1, lines1, "must receive one log entry")
	suite.True(finalized1, "must finalize")
	suite.Equal(0, lines2, "must not receive log entries")
	suite.False(finalized2, "must not finalize")
	suite.Equal(1, len(suite.handler.resultCallbacks), "resultCallbacks has one element")
	suite.Equal(1, len(suite.handler.logCallbacks), "logCallbacks has one element")
}
func (suite *HandlerTestSuite) TestFinalizeNotExistingCommand() {
	err := suite.handler.FinishCommand("id1", entities.NewSuccessCommand([]byte("OK")), nil)
	suite.NotNil(err, " FinishCommand must fail")
}

func (suite *HandlerTestSuite) TestAddLogEntryNotExistingCommand() {
	err := suite.handler.AddLogEntry("id1", "hello world!")
	suite.NotNil(err, " FinishCommand must fail")
}

func TestFoundation(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}
*/