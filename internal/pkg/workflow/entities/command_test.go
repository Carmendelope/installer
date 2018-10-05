/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
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
		gomega.Expect(err).To(gomega.BeNil())
		retrieved := &CommandResultFromJSON{}
		err = json.Unmarshal(result, retrieved)
		gomega.Expect(err).To(gomega.BeNil())
		toCR := retrieved.ToCommandResult()
		gomega.Expect(toCR).To(gomega.Equal(cr))
	})

	ginkgo.It("must be build from a message", func(){
		toReceiveNoError := `
    {"success":true, "output":"output"}
`
		retrieved := &CommandResultFromJSON{}
		err := json.Unmarshal([]byte(toReceiveNoError), retrieved)

		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(retrieved.Success).To(gomega.BeTrue())
		gomega.Expect(retrieved.Output).To(gomega.Equal("output"))
		gomega.Expect(retrieved.Error).To(gomega.BeNil())
		toCR := retrieved.ToCommandResult()
		gomega.Expect(toCR.Success).To(gomega.BeTrue())
		gomega.Expect(toCR.Output).To(gomega.Equal("output"))
		gomega.Expect(toCR.Error).To(gomega.BeNil())
	})
})

