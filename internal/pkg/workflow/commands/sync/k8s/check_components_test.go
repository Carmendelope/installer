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
	grpc_installer_go "github.com/nalej/grpc-installer-go"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
	"os"
)

var _ = ginkgo.Describe("A Check Components command", func() {
	ginkgo.FIt("should check components if no platform-dependent are present", func() {
		numYAML := 5
		componentsDir := CreateTempYAML(numYAML, 0)
		ns := []string{itAuxNamespace}
		ccCmd := NewCheckComponents("kubeConfigPath", ns)

		log.Debug().Msg("retrieving resources")
		toCheck, err := ccCmd.RetrieveResources()
		for _, d := range toCheck.Deployments {
			log.Debug().Str("deploy", d.Name).Msg("deploy resources")

		}
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(toCheck.Deployments).ShouldNot(gomega.BeNil())
		gomega.Expect(len(toCheck.Deployments)).Should(gomega.Equal(numYAML))
		gomega.Expect(os.RemoveAll(componentsDir)).To(gomega.Succeed())
	})

	ginkgo.It("should check components considering a single conflicting platform", func() {
		numYAML := 5
		componentsDir := CreateTempYAML(numYAML, numYAML, grpc_installer_go.Platform_AZURE.String())
		launchCmd := NewCheckComponents("kubeConfigPath", []string{})
		toCheck, err := launchCmd.RetrieveResources()
		gomega.Expect(err).To(gomega.Succeed())
		gomega.Expect(toCheck).ShouldNot(gomega.BeNil())
		gomega.Expect(len(toCheck.Deployments)).Should(gomega.Equal(numYAML))
		gomega.Expect(os.RemoveAll(componentsDir)).To(gomega.Succeed())
	})
})
