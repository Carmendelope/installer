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

import (
	"github.com/nalej/derrors"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"time"
)

var _ = ginkgo.Describe("Handler", func(){

	ginkgo.Context("with a basic workflow", func(){
		handler := NewCommandHandler().(*commandHandler)
		lines := 0
		finalized := false
		ginkgo.It("must support adding a command", func(){
			gomega.Expect(len(handler.resultCallbacks), gomega.Equal(0))
			gomega.Expect(len(handler.logCallbacks), gomega.Equal(0))
			err := handler.AddCommand("id1",
				func(id string, result *entities.CommandResult, error derrors.Error) {
					finalized = true
				},
				func(id string, logEntry string) {
					lines++
				},
			)

			gomega.Expect(err, gomega.BeNil())
			gomega.Expect(len(handler.resultCallbacks), gomega.Equal(1))
			gomega.Expect(len(handler.logCallbacks), gomega.Equal(1))
		})
		handler.AddLogEntry("id1", "hello world!")
		handler.FinishCommand("id1", entities.NewSuccessCommand([]byte("OK")), nil)
		time.Sleep(time.Second)
		ginkgo.It("must receive the callbacks", func(){
			gomega.Expect(lines, gomega.Equal(1))
			gomega.Expect(finalized, gomega.BeTrue())
			gomega.Expect(len(handler.resultCallbacks), gomega.Equal(0))
			gomega.Expect(len(handler.logCallbacks), gomega.Equal(0))
		})
	})

	ginkgo.Context("when adding a duplicated command", func(){
		handler := NewCommandHandler().(*commandHandler)
		err1 := handler.AddCommand("id1",
			func(id string, result *entities.CommandResult, error derrors.Error) {

			},
			func(id string, logEntry string) {

			},
		)
		err2 := handler.AddCommand("id1",
			func(id string, result *entities.CommandResult, error derrors.Error) {

			},
			func(id string, logEntry string) {

			},
		)
		ginkgo.It("must fail on the second command", func(){
			gomega.Expect(err1, gomega.BeNil())
			gomega.Expect(err2, gomega.Not(gomega.BeNil()))
		})
	})

	ginkgo.Context("with two commands", func(){
		handler := NewCommandHandler().(*commandHandler)
		lines1 := 0
		finalized1 := false
		err1 := handler.AddCommand("id1",
			func(id string, result *entities.CommandResult, error derrors.Error) {
				finalized1 = true
			},
			func(id string, logEntry string) {
				lines1++
			},
		)
		lines2 := 0
		finalized2 := false
		err2 := handler.AddCommand("id2",
			func(id string, result *entities.CommandResult, error derrors.Error) {
				finalized2 = true
			},
			func(id string, logEntry string) {
				lines2++
			},
		)

		ginkgo.It("must support adding two commands", func(){
			gomega.Expect(err1, gomega.BeNil())
			gomega.Expect(err2, gomega.BeNil())
			gomega.Expect(len(handler.resultCallbacks), gomega.Equal(2))
			gomega.Expect(len(handler.logCallbacks), gomega.Equal(2))
		})
		handler.AddLogEntry("id1", "hello world!")
		handler.FinishCommand("id1", entities.NewSuccessCommand([]byte("OK")), nil)
		time.Sleep(time.Second)
		ginkgo.Specify("cmd 1 must receive the callbacks", func(){
			gomega.Expect(lines1, gomega.Equal(1))
			gomega.Expect(finalized1, gomega.BeTrue())
			gomega.Expect(lines2, gomega.Equal(0))
			gomega.Expect(finalized2, gomega.BeFalse())
			gomega.Expect(len(handler.resultCallbacks), gomega.Equal(1))
			gomega.Expect(len(handler.logCallbacks), gomega.Equal(1))
		})
	})

	ginkgo.Context("receiving a finish callback on a non registered cmd", func(){
		handler := NewCommandHandler().(*commandHandler)
		err := handler.FinishCommand("id1", entities.NewSuccessCommand([]byte("OK")), nil)
		ginkgo.It("must fail", func(){
			gomega.Expect(err, gomega.Not(gomega.BeNil()))
		})
	})

	ginkgo.Context("receiving an add log callback on a non registered cmd", func(){
		handler := NewCommandHandler().(*commandHandler)
		err := handler.AddLogEntry("id1", "hello world!")
		ginkgo.It("must fail", func(){
			gomega.Expect(err, gomega.Not(gomega.BeNil()))
		})
	})


})
