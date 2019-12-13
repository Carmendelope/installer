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

package commands

import (
	"fmt"
	"github.com/nalej/installer/internal/app/installer-cli"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var managementClusterCmd = &cobra.Command{
	Use:   "management",
	Short: "Install the Nalej management cluster",
	Long:  `Install the Nalej management cluster`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		LaunchManagementInstall()
	},
}

func init() {
	cliCmd.AddCommand(managementClusterCmd)
}

func LaunchManagementInstall() {
	log.Info().Msg("Installing management cluster")
	err := ValidateInstallParameters()
	if err != nil {
		log.Panic().Str("error", err.DebugReport()).Msg("parameter validation failed")
	}
	paths, err := GetPaths()
	if err != nil {
		log.Panic().Str("error", err.DebugReport()).Msg("cannot obtain paths")
	}

	vErr := environment.Validate()
	if vErr != nil {
		log.Fatal().Str("trace", vErr.DebugReport()).Msg("Invalid environment")
	}
	environment.Print()

	inst, err := installer_cli.NewCLI(kubeConfigPath)
	if err != nil {
		log.Panic().Str("error", err.DebugReport()).Msg("cannot create CLI installer")
	}
	// Prepare the parameters.
	inst.PrepareInstallCommand(
		"cli-install",
		installKubernetes,
		username,
		privateKeyPath,
		strings.Split(nodes, ","),
		strings.ToUpper(targetPlatform),
		*paths,
		managementPublicHost,
		dnsClusterHost,
		strconv.Itoa(dnsClusterPort),
		useStaticIPAddresses,
		ipAddressIngress,
		ipAddressDNS,
		ipAddressCoreDNS,
		ipAddressVPNServer,
		false,
		environment,
		networkingMode,
		istioPath,
		istioCertsPath)

	if explainPlan {
		inst.LoadCredentials()
		fmt.Println(inst.Workflow.PrettyPrint())
	} else {
		inst.Execute()
	}
}
