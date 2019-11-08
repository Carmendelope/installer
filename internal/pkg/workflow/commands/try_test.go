/*
 * Copyright 2019 Nalej
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
 *
 */

package commands

import (
	"github.com/nalej/installer/internal/pkg/workflow/commands/async"
	"github.com/nalej/installer/internal/pkg/workflow/commands/sync"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Try", func() {
	ginkgo.Context("With SYNC commands", func() {
		ginkgo.It("must support with two sync commands", func() {
			cmd1 := sync.NewLogger("cmd1")
			cmd2 := sync.NewLogger("cmd2")
			try := NewTry("test sync", cmd1, cmd2)
			wID := "testWorkflow"
			result, err := try.Run(wID)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(result.Success).To(gomega.BeTrue())
			gomega.Expect(result.Output).To(gomega.Equal("cmd1"))
		})
		ginkgo.It("on failure must execute the second command", func() {
			cmd1 := sync.NewFail()
			cmd2 := sync.NewLogger("cmd2")
			try := NewTry("test sync fail", cmd1, cmd2)
			wID := "testWorkflow"
			result, err := try.Run(wID)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(result.Success).To(gomega.BeTrue())
			gomega.Expect(result.Output).To(gomega.Equal("cmd2"))
		})
	})

	ginkgo.Context("With ASYNC commands", func() {
		ginkgo.It("must support async commands", func() {
			cmd1 := async.NewSleep("0")
			cmd2 := sync.NewLogger("cmd2")
			try := NewTry("test async", cmd1, cmd2)
			wID := "testWorkflow"
			result, err := try.Run(wID)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(result.Success).To(gomega.BeTrue())
			gomega.Expect(result.Output).To(gomega.Equal("Slept for 0"))
		})
		ginkgo.It("on failure must execute the second command", func() {
			cmd1 := async.NewFail()
			cmd2 := async.NewSleep("0")
			try := NewTry("test async", cmd1, cmd2)
			wID := "testWorkflow"
			result, err := try.Run(wID)
			gomega.Expect(err).To(gomega.BeNil())
			gomega.Expect(result.Success).To(gomega.BeTrue())
			gomega.Expect(result.Output).To(gomega.Equal("Slept for 0"))
		})
	})

	ginkgo.It("Must be buildable from JSON", func() {
		fromJSON := `
{"type":"sync", "name": "try", "description":"Try",
"cmd": {"type":"sync", "name": "logger", "msg": "This is a logging message"},
"onFail": {"type":"sync", "name": "logger", "msg": "This is a logging message"}}
`
		received, err := NewTryFromJSON([]byte(fromJSON))
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect((*received).(*Try).TryCommand).ToNot(gomega.BeNil())
		gomega.Expect((*received).(*Try).OnFailCommand).ToNot(gomega.BeNil())
	})

})
