/*
 * Copyright 2020 Nalej
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

package k8s

import (
	"fmt"
	"github.com/nalej/grpc-installer-go"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// CreateTempYAML creates a directory with a set of yaml files.
func CreateTempYAML(numYAML int, numPlatformYAML int, platforms ...string) string {
	dir, err := ioutil.TempDir("", "launch")
	gomega.Expect(err).Should(gomega.Succeed())
	testData := []byte("test")
	for i := 0; i < numYAML; i++ {
		fileName := fmt.Sprintf("%d.yaml", i)
		err := ioutil.WriteFile(filepath.Join(dir, fileName), testData, 0777)
		gomega.Expect(err).Should(gomega.Succeed())
	}

	for _, targetPlatform := range platforms {
		for i := 0; i < numPlatformYAML; i++ {
			fileName := fmt.Sprintf("%d.yaml.%s", i, strings.ToLower(targetPlatform))
			err := ioutil.WriteFile(filepath.Join(dir, fileName), testData, 0777)
			gomega.Expect(err).Should(gomega.Succeed())
		}
	}

	return dir
}

var _ = ginkgo.Describe("A Launch command", func() {
	ginkgo.It("should list components if no platform-dependent are present", func() {
		numYAML := 10
		componentsDir := CreateTempYAML(numYAML, 0)
		launchCmd := NewLaunchComponents("kubeConfigPath", []string{}, componentsDir, grpc_installer_go.Platform_AZURE.String())
		toInstall, err := launchCmd.ListComponents()
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(toInstall).ShouldNot(gomega.BeNil())
		gomega.Expect(len(toInstall)).Should(gomega.Equal(numYAML))
		gomega.Expect(os.RemoveAll(componentsDir)).To(gomega.Succeed())
	})

	ginkgo.It("should list components considering a single conflicting platform", func() {
		numYAML := 10
		componentsDir := CreateTempYAML(numYAML, numYAML, grpc_installer_go.Platform_AZURE.String())
		launchCmd := NewLaunchComponents("kubeConfigPath", []string{}, componentsDir, grpc_installer_go.Platform_AZURE.String())
		toInstall, err := launchCmd.ListComponents()
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(toInstall).ShouldNot(gomega.BeNil())
		gomega.Expect(len(toInstall)).Should(gomega.Equal(numYAML))
		gomega.Expect(os.RemoveAll(componentsDir)).To(gomega.Succeed())
	})

	ginkgo.It("should list components considering a multiple conflicting platforms", func() {
		numYAML := 10
		numPlatformYAML := numYAML / 2
		componentsDir := CreateTempYAML(numYAML, numPlatformYAML, grpc_installer_go.Platform_AZURE.String(), grpc_installer_go.Platform_BAREMETAL.String())
		launchCmd := NewLaunchComponents("kubeConfigPath", []string{}, componentsDir, grpc_installer_go.Platform_AZURE.String())
		toInstall, err := launchCmd.ListComponents()
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(toInstall).ShouldNot(gomega.BeNil())
		gomega.Expect(len(toInstall)).Should(gomega.Equal(numYAML))
		gomega.Expect(os.RemoveAll(componentsDir)).To(gomega.Succeed())
		for i := 0; i < numYAML; i++ {
			expectedName := fmt.Sprintf("%d.yaml.%s", i, strings.ToLower(launchCmd.PlatformType))
			if i >= numPlatformYAML {
				expectedName = fmt.Sprintf("%d.yaml", i)
			}
			gomega.Expect(toInstall[i]).Should(gomega.Equal(expectedName))
		}
	})
})
