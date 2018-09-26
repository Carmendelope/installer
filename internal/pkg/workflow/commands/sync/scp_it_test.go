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

// SCP Integration tests
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
	"io/ioutil"
	"os"
	"os/user"
	"path"
)

func getUserPrivateKey() string {
	var privateKey []byte
	usr, err := user.Current()

	gomega.Expect(err).To(gomega.BeNil())
	homeDirectory := usr.HomeDir
	privateKeyFile := path.Join(homeDirectory, ".ssh", "id_rsa")
	privateKey, err = ioutil.ReadFile(privateKeyFile)
	gomega.Expect(err).To(gomega.BeNil())

	return string(privateKey)
}


var _ = ginkgo.Describe("An SCP command", func(){
	if ! utils.RunIntegrationTests() {
		log.Warn().Msg("Integration tests are skipped")
		return
	}
	var (
		testUsername = "root"
		testPassword = "root"
		targetHost = os.Getenv("IT_SSH_HOST")
		targetPort = os.Getenv("IT_SSH_PORT")
		targetPath = "/tmp/"
	)

	if targetHost == "" || targetPort == "" {
		ginkgo.Fail("missing environment variables")
	}

	ginkgo.It("must be able to copy a file", func(){
		content := []byte("this is a testing file")
		tmpfile, err := ioutil.TempFile("", "example")
		gomega.Expect(err).To(gomega.BeNil())
		defer os.Remove(tmpfile.Name()) // clean up

		size, err := tmpfile.Write(content)
		gomega.Expect(err).To(gomega.BeNil())
		err = tmpfile.Close()
		gomega.Expect(err).To(gomega.BeNil())
		log.Debug().Str("file", tmpfile.Name()).Int("size", size).Msg("file is written")


		credentials := entities.NewCredentials(testUsername, testPassword)
		cmd := NewSCP(targetHost, targetPort, *credentials, tmpfile.Name(), targetPath)
		result, err := cmd.Run("w1")
		gomega.Expect(err).To(gomega.BeNil())
		log.Debug().Bool("result", (*result).Success).Str("output", (*result).Output).Msg("scp has been executed")
	})

	ginkgo.It("must be able to copy using PKI", func(){
		privateKey := getUserPrivateKey()
		content := []byte("this is a testing file to be copied with scp over PKI")
		tmpfile, err := ioutil.TempFile("", "example")
		gomega.Expect(err).To(gomega.BeNil())
		defer os.Remove(tmpfile.Name()) // clean up

		size, err := tmpfile.Write(content)
		gomega.Expect(err).To(gomega.BeNil())
		err = tmpfile.Close()
		gomega.Expect(err).To(gomega.BeNil())
		log.Debug().Str("file", tmpfile.Name()).Int("size", size).Msg("file is written")


		credentials := entities.NewPKICredentials(testUsername, string(privateKey))
		cmd := NewSCP(targetHost, targetPort, *credentials, tmpfile.Name(), targetPath)
		result, err := cmd.Run("w1")
		gomega.Expect(err).To(gomega.BeNil())
		log.Debug().Bool("result", (*result).Success).Str("output", (*result).Output).Msg("scp has been executed")
	})


})

