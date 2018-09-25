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

import (
	"encoding/json"
	"github.com/nalej/derrors"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Command structure", func(){
	ginkgo.It("must be parsed from JSON", func(){
		de := derrors.NewGenericError("some error")
		cr := NewCommandResult(true, "output", de)
		result, err := json.Marshal(cr)
		gomega.Expect(err, gomega.BeNil())
		retrieved := &CommandResultFromJSON{}
		err = json.Unmarshal(result, retrieved)
		gomega.Expect(err, gomega.BeNil())
		toCR := retrieved.ToCommandResult()
		gomega.Expect(toCR, gomega.Equal(cr))
	})

	ginkgo.It("must be build from a message", func(){
		toReceiveNoError := `
    {"success":true, "output":"output"}
`
		retrieved := &CommandResultFromJSON{}
		err := json.Unmarshal([]byte(toReceiveNoError), retrieved)

		gomega.Expect(err, gomega.BeNil())
		gomega.Expect(retrieved.Success, gomega.BeTrue())
		gomega.Expect(retrieved.Output, gomega.Equal("output"))
		gomega.Expect(retrieved.Error, gomega.BeNil())
		toCR := retrieved.ToCommandResult()
		gomega.Expect(toCR.Success, gomega.BeTrue())
		gomega.Expect(toCR.Output, gomega.Equal("output"))
		gomega.Expect(toCR.Error, gomega.BeNil())
	})
})

