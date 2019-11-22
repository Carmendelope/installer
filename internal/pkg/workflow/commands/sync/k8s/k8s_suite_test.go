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

package k8s

import (
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"os"
	"testing"
)

const itAuxNamespace = "test-it-launch"
const itNalejNamespace = "nalej"

var itComponentsDir string
var itKubeConfigFile = os.Getenv("IT_K8S_KUBECONFIG")
var itTargetNamespaces = []string{itAuxNamespace, itNalejNamespace}
var itRegistryUsername = os.Getenv("IT_REGISTRY_USERNAME")
var itRegistryPassword = os.Getenv("IT_REGISTRY_PASSWORD")
var itTestTargetNamespace = os.Getenv("IT_TARGET_NAMESPACE")

func TestK8sPackage(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "K8s package suite")
}

/*
var _ = ginkgo.AfterSuite(func() {
	log.Info().Msg("Cleaning test environment")
	if itComponentsDir != "" {
		os.RemoveAll(itComponentsDir)
	}
	if itKubeConfigFile != "" && len(itTargetNamespaces) > 0 {
		// for _, ns := range itTargetNamespaces{
		//	tc := NewTestCleaner(itKubeConfigFile, ns)
		//	gomega.Expect(tc.DeleteAll()).To(gomega.Succeed())
		//}
	} else {
		log.Warn().Msg("TestCleaner skipped")
	}
})

*/
