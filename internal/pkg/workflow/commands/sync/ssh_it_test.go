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

// SSH Integration tests
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
	"strings"
)

var _ = ginkgo.Describe("An SCP command", func() {
	if !utils.RunIntegrationTests() {
		log.Warn().Msg("Integration tests are skipped")
		return
	}
	var (
		testUsername = "root"
		testPassword = "root"
		targetHost   = os.Getenv("IT_SSH_HOST")
		targetPort   = os.Getenv("IT_SSH_PORT")
	)

	if targetHost == "" || targetPort == "" {
		ginkgo.Fail("missing environment variables")
	}

	ginkgo.It("must be able to execute a simple command", func() {
		credentials := entities.NewCredentials(testUsername, testPassword)
		args := make([]string, 2)
		args[0] = "-lash"
		args[1] = "/var/"
		cmd := NewSSH(targetHost, targetPort, *credentials, "ls", args)
		result, err := cmd.Run("w1")
		gomega.Expect(err).To(gomega.BeNil())
		output := (*result).Output
		gomega.Expect(strings.Contains(output, "local")).To(gomega.BeTrue())
	})

	ginkgo.It("must be able to execute a command using PKI", func() {
		privateKey := getUserPrivateKey()
		credentials := entities.NewPKICredentials(testUsername, string(privateKey))
		args := make([]string, 2)
		args[0] = "-lash"
		args[1] = "/var/"
		cmd := NewSSH(targetHost, targetPort, *credentials, "ls", args)
		result, err := cmd.Run("w1")
		gomega.Expect(err).To(gomega.BeNil())
		output := (*result).Output
		gomega.Expect(strings.Contains(output, "local")).To(gomega.BeTrue())
	})

})
