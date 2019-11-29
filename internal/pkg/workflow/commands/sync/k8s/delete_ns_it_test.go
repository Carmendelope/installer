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
	"github.com/nalej/installer/internal/pkg/utils"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"github.com/rs/zerolog/log"
)

var _ = ginkgo.Describe("A delete service account command", func() {

	if !utils.RunIntegrationTests() {
		log.Warn().Msg("Integration tests are skipped")
		return
	}

	if !utils.RunIntegrationTest("delete_ns_it_test") {
		log.Warn().Msg("Integration test is skipped")
		return
	}

	if itKubeConfigFile == "" || itTestTargetNamespace == "" {
		ginkgo.Fail("missing environment variables")
	}

	ginkgo.It("should be able to delete a namespace", func() {
		dsa := NewDeleteNamespace(itKubeConfigFile, itTestTargetNamespace)
		result, err := dsa.Run("deleteNamespace")
		gomega.Expect(err).To(gomega.Succeed())
		if !result.Success {
			log.Debug().Interface("result", result).Msg("failed")
		}
		gomega.Expect(result.Success).Should(gomega.BeTrue())
	})

})
