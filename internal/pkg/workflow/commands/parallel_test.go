/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

// Parallel command tests
//

package commands

import (
	"github.com/nalej/installer/internal/pkg/workflow/commands/sync"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

var _ = ginkgo.Describe("Parallel command", func(){
	ginkgo.It("With 3 sync commands", func(){
		cmd1 := sync.NewLogger("cmd1")
		cmd2 := sync.NewSleep("3")
		cmd3 := sync.NewLogger("cmd2")

		p := NewParallel("test synchronous commands", 3, []entities.Command{cmd1, cmd2, cmd3})
		wID := "testWorkflow"
		result, err := p.Run(wID)
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(result.Success).To(gomega.BeTrue())
	})

	ginkgo.It("Must be buildable from JSON", func(){
		fromJSON := `
{"type":"sync", "name": "parallel", "maxParallelism":3, "commands": [
{"type":"sync", "name": "logger", "msg": "This is a logging message"},
{"type":"sync", "name": "logger", "msg": "This is a logging message"}]}
`
		received, err := NewParallelFromJSON([]byte(fromJSON))
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect((*received).(*Parallel).MaxParallelism).To(gomega.Equal(3))
		fromJSONWithoutParallelism := `
{"type":"sync", "name": "parallel", "commands": [
{"type":"sync", "name": "logger", "msg": "This is a logging message"},
{"type":"sync", "name": "logger", "msg": "This is a logging message"}]}
`
		received, err = NewParallelFromJSON([]byte(fromJSONWithoutParallelism))
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect((*received).(*Parallel).MaxParallelism).To(gomega.Equal(0))
	})

	ginkgo.It("Should stop on failure", func(){
		cmd1 := sync.NewFail()
		cmd2 := sync.NewSleep("10")
		cmd3 := sync.NewSleep("10")

		p := NewParallel("test synchronous commands", 3, []entities.Command{cmd1, cmd2, cmd3})
		wID := "testWorkflow"
		result, err := p.Run(wID)
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(result.Success).To(gomega.BeFalse())
	})

	ginkgo.It("Must support a max level", func(){
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
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(result.Success).To(gomega.BeTrue())
	})

	ginkgo.It("must support higher max levels - DP-1164", func(){
		cmd1 := sync.NewLogger("cmd1")
		cmd2 := sync.NewLogger("cmd2")

		p := NewParallel("test synchronous commands", 3,
			[]entities.Command{cmd1, cmd2})

		wID := "testWorkflow"
		result, err := p.Run(wID)
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(result.Success).To(gomega.BeTrue())
	})

})
