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

package ingress

import (
	"github.com/nalej/installer/internal/pkg/utils"
	"github.com/onsi/ginkgo"
	"github.com/rs/zerolog/log"
)

var _ = ginkgo.Describe("A install ingress command", func() {

	if !utils.RunIntegrationTests() {
		log.Warn().Msg("Integration tests are skipped")
		return
	}
	/*
		if itKubeConfigFile == "" {
			ginkgo.Fail("missing environment variables")
		}

		ginkgo.FIt("should be able to install an ingress", func(){
			ii := NewInstallIngress(itKubeConfigFile, "nalej.cluster.local")
			result, err := ii.Run("installIngress")
			gomega.Expect(err).To(gomega.Succeed())
			gomega.Expect(result.Success).Should(gomega.BeTrue())
		})
	*/
})
