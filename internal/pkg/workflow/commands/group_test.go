/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
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

var _ = ginkgo.Describe("Group command", func(){
	ginkgo.It("Must support a basic sequence", func(){
		cmd1 := sync.NewLogger("cmd1")
		cmd2 := sync.NewSleep("1")
		cmd3 := sync.NewLogger("cmd2")
		g := NewGroup("basicSequence", []entities.Command{cmd1, cmd2, cmd3})
		wID := "TestBasicSequence"
		result, err := g.Run(wID)
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(result.Success).To(gomega.BeTrue())
	})

	ginkgo.It("Must support a basic sequence with ASYNC commands", func(){
		cmd1 := sync.NewLogger("cmd1")
		cmd2 := async.NewSleep("1")
		cmd3 := sync.NewLogger("cmd2")
		g := NewGroup("basicSequence", []entities.Command{cmd1, cmd2, cmd3})
		wID := "TestBasicSequence"
		result, err := g.Run(wID)
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(result.Success).To(gomega.BeTrue())
	})

	ginkgo.It("Must stop on fail", func(){
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
