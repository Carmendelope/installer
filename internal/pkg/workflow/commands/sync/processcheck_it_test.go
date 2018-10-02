/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

// ProcessChecks Integration tests
//
// Prerequirements:
//   Launch a docker image with an sshd service
//   $ docker run --rm --publish=2222:22 sickp/alpine-sshd:7.5
//
// Copy your PKI credentials
//   $ ssh-copy-id root@localhost -p 2222

/*
RUN_INTEGRATION_TEST=true
IT_SSH_HOST=localhost
IT_SSH_PORT=2222
 */

package sync

import (
	"github.com/nalej/installer/internal/pkg/utils"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"os"
)

var _ = ginkgo.Describe("An SCP command", func() {
	if ! utils.RunIntegrationTests() {
		log.Warn().Msg("Integration tests are skipped")
		return
	}
	var (
		testUsername= "root"
		testPassword= "root"
		targetHost= os.Getenv("IT_SSH_HOST")
		targetPort= os.Getenv("IT_SSH_PORT")
	)

	if targetHost == "" || targetPort == "" {
		ginkgo.Fail("missing environment variables")
	}

	ginkgo.It("should be able to check that a process exists when expecting it to exist", func(){
		credentials := entities.NewCredentials(testUsername, testPassword)
		cmd := NewProcessCheck(targetHost, targetPort, *credentials, "sshd", true)
		result, err := cmd.Run("w1")
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(result.Success).To(gomega.BeTrue())
		output := (*result).Output
		gomega.Expect(output).To(gomega.Equal("Process sshd has been found"))
	})

	ginkgo.It("should be able to check that a process exists when expecting it to exist using PKI", func(){
		privateKey := getUserPrivateKey()
		credentials := entities.NewPKICredentials(testUsername, string(privateKey))
		cmd := NewProcessCheck(targetHost, targetPort, *credentials, "sshd", true)
		result, err := cmd.Run("w1")
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(result.Success).To(gomega.BeTrue())
		output := (*result).Output
		gomega.Expect(output).To(gomega.Equal("Process sshd has been found"))
	})

	ginkgo.It("should fail when a command does not exists and the process expects it to", func(){
		credentials := entities.NewCredentials(testUsername, testPassword)
		cmd := NewProcessCheck(targetHost, targetPort, *credentials, "notFound", true)
		result, err := cmd.Run("w1")
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(result.Success).To(gomega.BeFalse())
		output := (*result).Output
		gomega.Expect(output).To(gomega.Equal("Process notFound has not been found and should exist"))
	})

	ginkgo.It("should work when a command does not exists and the process does not expect it to", func(){
		credentials := entities.NewCredentials(testUsername, testPassword)
		cmd := NewProcessCheck(targetHost, targetPort, *credentials, "notFound", false)
		result, err := cmd.Run("w1")
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(result.Success).To(gomega.BeTrue())
		output := (*result).Output
		gomega.Expect(output).To(gomega.Equal("Process notFound has not been found"))
	})

	ginkgo.It("should fail when a process exists and the process does not expects it to", func(){
		credentials := entities.NewCredentials(testUsername, testPassword)
		cmd := NewProcessCheck(targetHost, targetPort, *credentials, "sshd", false)
		result, err := cmd.Run("w1")
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(result.Success).To(gomega.BeFalse())
		output := (*result).Output
		gomega.Expect(output).To(gomega.Equal("Process sshd has been found and should not exist"))
	})

})
