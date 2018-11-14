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

// Run command to launch the main component API.

package commands

import (
	"github.com/nalej/installer/internal/pkg/server"
	cfg "github.com/nalej/installer/internal/pkg/server/config"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var config = cfg.Config{}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Launch the server API",
	Long:  `Launch the server API`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		log.Info().Msg("Launching API!")
		server := server.NewService(config)
		server.Run()
	},
}

func init() {

	runCmd.Flags().IntVar(&config.Port, "port", 8900, "Port to launch the Installer")
	runCmd.PersistentFlags().StringVar(&config.ManagementClusterHost, "managementClusterPublicHost", "",
		"Public FQDN where the management cluster is reachable by the application clusters")
	runCmd.MarkPersistentFlagRequired("managementClusterPublicHost")
	runCmd.PersistentFlags().StringVar(&config.ManagementClusterPort, "managementClusterPublicPort", "",
		"Public port where the management cluster is reachable by the application clusters")
	runCmd.MarkPersistentFlagRequired("managementClusterPublicPort")

	runCmd.PersistentFlags().StringVar(&config.ComponentsPath, "componentsPath", "./assets/",
		"Directory with the components to be installed")
	runCmd.PersistentFlags().StringVar(&config.BinaryPath, "binaryPath", "./bin/",
		"Directory with the binary executables")
	runCmd.PersistentFlags().StringVar(&config.TempPath, "tempPath", "./temp/",
		"Directory to store temporal files")

	rootCmd.AddCommand(runCmd)
}