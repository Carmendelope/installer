/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
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
	runCmd.PersistentFlags().StringVar(&config.DNSClusterHost, "dnsClusterPublicHost", "",
		"Public FQDN where the management cluster is reachable for DNS requests by the application clusters")
	runCmd.MarkPersistentFlagRequired("dnsClusterPublicHost")
	runCmd.PersistentFlags().StringVar(&config.DNSClusterPort, "dnsClusterPublicPort", "",
		"Public port where the management cluster is reachable for DNS request by the application clusters")
	runCmd.MarkPersistentFlagRequired("dnsClusterPublicPort")

	runCmd.PersistentFlags().StringVar(&config.ComponentsPath, "componentsPath", "./assets/",
		"Directory with the components to be installed")
	runCmd.PersistentFlags().StringVar(&config.BinaryPath, "binaryPath", "./bin/",
		"Directory with the binary executables")
	runCmd.PersistentFlags().StringVar(&config.TempPath, "tempPath", "./temp/",
		"Directory to store temporal files")

	runCmd.PersistentFlags().StringVar(&config.DockerRegistryUsername, "dockerUsername", "",
		"Username to download internal images from the docker repository")
	runCmd.PersistentFlags().StringVar(&config.DockerRegistryPassword, "dockerPassword", "",
		"Password to download internal images from the docker repository")

	runCmd.PersistentFlags().StringVar(&config.ZTPlanetSecretValue, "planetSecret", "",
		"Secret for the ZeroTier Planet file")

	rootCmd.AddCommand(runCmd)
}
