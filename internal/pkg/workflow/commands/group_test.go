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

// Group command tests
//

package commands

import (
	"github.com/nalej/installer/internal/pkg/workflow/commands/async"
	"github.com/nalej/installer/internal/pkg/workflow/commands/sync"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Group command", func() {
	ginkgo.It("Must support a basic sequence", func() {
		cmd1 := sync.NewLogger("cmd1")
		cmd2 := sync.NewSleep("1")
		cmd3 := sync.NewLogger("cmd2")
		g := NewGroup("basicSequence", []entities.Command{cmd1, cmd2, cmd3})
		wID := "TestBasicSequence"
		result, err := g.Run(wID)
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(result.Success).To(gomega.BeTrue())
	})

	ginkgo.It("Must support a basic sequence with ASYNC commands", func() {
		cmd1 := sync.NewLogger("cmd1")
		cmd2 := async.NewSleep("1")
		cmd3 := sync.NewLogger("cmd2")
		g := NewGroup("basicSequence", []entities.Command{cmd1, cmd2, cmd3})
		wID := "TestBasicSequence"
		result, err := g.Run(wID)
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(result.Success).To(gomega.BeTrue())
	})

	ginkgo.It("Must stop on fail", func() {
		cmd1 := sync.NewLogger("cmd1")
		cmd2 := sync.NewFail()
		cmd3 := sync.NewSleep("1")
		cmd4 := sync.NewLogger("should not appear")
		g := NewGroup("basicSequence", []entities.Command{cmd1, cmd2, cmd3, cmd4})
		wID := "TestBasicSequenceFail"
		result, err := g.Run(wID)
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(result.Success).To(gomega.BeFalse())
	})
})
