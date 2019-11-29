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
	installer_cli "github.com/nalej/installer/internal/app/installer-cli"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"strings"
)

var appCluster bool

var uninstallLongHelp = `
Uninstall the Nalej components deployed by the installer

This command will remove the contents of the Nalej namespace, and all the Kubernetes
entities created during installation. Notice that the cluster certificate will be
removed on the decomission process attending to the certificate manager used to
created it.
`

var uninstallExample = `

# Uninstall a management cluster
installer-cli uninstall nalej/mngtCluster.yaml

# Uninstall an application cluster
installer-cli uninstall nalej/appCluster.yaml --appCluster

# Show the uninstall plan
installer-cli uninstall nalej/mngtCluster.yaml --explainPlan
`

var uninstallClusterCmd = &cobra.Command{
	Use:     "uninstall <kubeConfigPath>",
	Short:   "Uninstall a Nalej cluster",
	Long:    uninstallLongHelp,
	Example: uninstallExample,
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		LaunchUninstall(args[0])
	},
}

func init() {
	uninstallClusterCmd.Flags().BoolVar(&explainPlan, "explainPlan", false,
		"Show install plan instead of performing the uninstall")
	uninstallClusterCmd.Flags().BoolVar(&appCluster, "appCluster", false,
		"Set to true if the target cluster is an application cluster.")
	rootCmd.AddCommand(uninstallClusterCmd)
}

// LaunchUninstall triggers the uninstall process of a given cluster.
func LaunchUninstall(kubeConfig string) {
	inst, err := installer_cli.NewCLI(kubeConfig)
	if err != nil {
		log.Panic().Str("error", err.DebugReport()).Msg("cannot create CLI installer")
	}
	inst.PrepareUninstallCommand(
		"cli-uninstall",
		"nalej",
		"cli-cluster-request",
		strings.ToUpper(targetPlatform),
		appCluster)

	if explainPlan {
		inst.LoadCredentials()
		fmt.Println(inst.Workflow.PrettyPrint())
	} else {
		inst.Execute()
	}
}
