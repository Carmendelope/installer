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

// SCP & SSH Integration tests
//
// Prerequirements:
//   Launch a docker image with an sshd service
//   $ docker run --rm --publish=2222:22 sickp/alpine-sshd:7.5
//
// Copy your PKI credentials
//   $ ssh-copy-id root@localhost -p 2222

package sync

/*
func GetUserPrivateKey(t *testing.T) string {
	usr, err := user.Current()
	assert.Nil(t, err, "current user should be available")
	homeDirectory := usr.HomeDir
	privateKeyFile := path.Join(homeDirectory, ".ssh", "id_rsa")
	privateKey, err := ioutil.ReadFile(privateKeyFile)
	assert.Nil(t, err, "expecting ssh private key to be available")
	return string(privateKey)
}

func SCPToRemote(t *testing.T) {
	testUsername := "root"
	testPassword := "root"
	targetHost := "localhost"
	targetPort := "2222"
	targetPath := "/tmp/"

	content := []byte("this is a testing file")
	tmpfile, err := ioutil.TempFile("", "example")
	assert.Nil(t, err, "temp file should be created")
	defer os.Remove(tmpfile.Name()) // clean up

	size, err := tmpfile.Write(content)
	assert.Nil(t, err, "file should be writable")
	err = tmpfile.Close()
	assert.Nil(t, err, "file should be closed")
	fmt.Println("File: ", tmpfile.Name(), " size: ", size)

	credentials := entities.NewCredentials(testUsername, testPassword)
	cmd := NewSCP(targetHost, targetPort, *credentials, tmpfile.Name(), targetPath)
	result, err := cmd.Run("w1")
	assert.Nil(t, err, "scp should work")
	fmt.Println("result: ", (*result).Success, (*result).Output)
}

func SCPToRemotePKI(t *testing.T) {
	testUsername := "root"
	privateKey := GetUserPrivateKey(t)
	targetHost := "localhost"
	targetPort := "2222"
	targetPath := "/tmp/"

	content := []byte("this is a testing file to be copied with scp over PKI")
	tmpfile, err := ioutil.TempFile("", "examplePKI")
	assert.Nil(t, err, "temp file should be created")
	defer os.Remove(tmpfile.Name()) // clean up

	size, err := tmpfile.Write(content)
	assert.Nil(t, err, "file should be writable")
	err = tmpfile.Close()
	assert.Nil(t, err, "file should be closed")
	fmt.Println("File: ", tmpfile.Name(), " size: ", size)

	credentials := entities.NewPKICredentials(testUsername, string(privateKey))
	cmd := NewSCP(targetHost, targetPort, *credentials, tmpfile.Name(), targetPath)
	result, err := cmd.Run("w1")
	assert.Nil(t, err, "scp should work")
	fmt.Println("result: ", (*result).Success, (*result).Output)
}

func SSHExec(t *testing.T) {
	testUsername := "root"
	testPassword := "root"
	targetHost := "localhost"
	targetPort := "2222"
	credentials := entities.NewCredentials(testUsername, testPassword)
	args := make([]string, 2)
	args[0] = "-lash"
	args[1] = "/var/"
	cmd := NewSSH(targetHost, targetPort, *credentials, "ls", args)
	result, err := cmd.Run("w1")
	assert.Nil(t, err, "ssh should work")
	output := (*result).Output
	assert.True(t, strings.Contains(output, "local"), "local dir should be there")
}

func SSHExecPKI(t *testing.T) {
	testUsername := "root"
	privateKey := GetUserPrivateKey(t)
	targetHost := "localhost"
	targetPort := "2222"
	credentials := entities.NewPKICredentials(testUsername, string(privateKey))
	args := make([]string, 2)
	args[0] = "-lash"
	args[1] = "/var/"
	cmd := NewSSH(targetHost, targetPort, *credentials, "ls", args)
	result, err := cmd.Run("w1")
	assert.Nil(t, err, "ssh should work")
	output := (*result).Output
	assert.True(t, strings.Contains(output, "local"), "local dir should be there")
}

func ProcessCheckTestExist(t *testing.T) {
	testUsername := "root"
	testPassword := "root"
	targetHost := "localhost"
	targetPort := "2222"
	credentials := entities.NewCredentials(testUsername, testPassword)
	cmd := NewProcessCheck(targetHost, targetPort, *credentials, "sshd", true)
	result, err := cmd.Run("w1")
	assert.Nil(t, err, "process check should work")
	assert.True(t, result.Success, "expecting ok")
	output := (*result).Output
	assert.Equal(t, "Process sshd has been found", output, "message should match")
}

func ProcessCheckTestExistPKI(t *testing.T) {
	testUsername := "root"
	privateKey := GetUserPrivateKey(t)
	targetHost := "localhost"
	targetPort := "2222"
	credentials := entities.NewPKICredentials(testUsername, string(privateKey))
	cmd := NewProcessCheck(targetHost, targetPort, *credentials, "sshd", true)
	result, err := cmd.Run("w1")
	assert.Nil(t, err, "process check should work")
	assert.True(t, result.Success, "expecting ok")
	output := (*result).Output
	assert.Equal(t, "Process sshd has been found", output, "message should match")
}

func ProcessCheckTestExistFail(t *testing.T) {
	testUsername := "root"
	testPassword := "root"
	targetHost := "localhost"
	targetPort := "2222"
	credentials := entities.NewCredentials(testUsername, testPassword)
	cmd := NewProcessCheck(targetHost, targetPort, *credentials, "notFound", true)
	result, err := cmd.Run("w1")
	assert.Nil(t, err, "process check should work")
	assert.False(t, result.Success, "expecting fail")
	output := (*result).Output
	assert.Equal(t, "Process notFound has not been found and should exist", output, "message should match")
}

func ProcessCheckTestNotExist(t *testing.T) {
	testUsername := "root"
	testPassword := "root"
	targetHost := "localhost"
	targetPort := "2222"
	credentials := entities.NewCredentials(testUsername, testPassword)
	cmd := NewProcessCheck(targetHost, targetPort, *credentials, "notFound", false)
	result, err := cmd.Run("w1")
	assert.Nil(t, err, "process check should work")
	assert.True(t, result.Success, "expecting ok")
	output := (*result).Output
	assert.Equal(t, "Process notFound has not been found", output, "message should match")
}

func ProcessCheckTestNotExistFail(t *testing.T) {
	testUsername := "root"
	testPassword := "root"
	targetHost := "localhost"
	targetPort := "2222"
	credentials := entities.NewCredentials(testUsername, testPassword)
	cmd := NewProcessCheck(targetHost, targetPort, *credentials, "sshd", false)
	result, err := cmd.Run("w1")
	assert.Nil(t, err, "process check should work")
	assert.False(t, result.Success, "expecting false")
	output := (*result).Output
	assert.Equal(t, "Process sshd has been found and should not exist", output, "message should match")
}

func TestSCP(t *testing.T) {
	fmt.Println("Running SCP integration tests: " + strconv.FormatBool(utils.RunIntegrationTests()))
	if utils.RunIntegrationTests() {
		utils.EnableDebug()
		SCPToRemote(t)
		SSHExec(t)
		ProcessCheckTestExist(t)
		ProcessCheckTestExistFail(t)
		ProcessCheckTestNotExist(t)
		ProcessCheckTestNotExistFail(t)
		// PKI tests
		SSHExecPKI(t)
		SCPToRemotePKI(t)
		ProcessCheckTestExistPKI(t)
	}
}

*/