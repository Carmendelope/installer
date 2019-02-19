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

// Add parameters related to the usage of registries.
func addRegistryOptions(cliCmd *cobra.Command){
	cliCmd.Flags().StringVar(&config.Environment.TargetEnvironment, "targetEnvironment", "PRODUCTION", "Target environment to be installed: PRODUCTION, STAGING, or DEVELOPMENT")
	// Production
	cliCmd.PersistentFlags().StringVar(&config.Environment.ProdRegistryUsername, "prodRegistryUsername", "",
		"Username to download internal images from the production docker registry. Alternatively you may use PROD_REGISTRY_USERNAME")
	cliCmd.PersistentFlags().StringVar(&config.Environment.ProdRegistryPassword, "prodRegistryPassword", "",
		"Password to download internal images from the production docker registry. Alternatively you may use PROD_REGISTRY_PASSWORD")
	cliCmd.PersistentFlags().StringVar(&config.Environment.ProdRegistryURL, "prodRegistryURL", "",
		"URL of the production docker registry. Alternatively you may use PROD_REGISTRY_URL")
	// Staging
	cliCmd.PersistentFlags().StringVar(&config.Environment.StagingRegistryUsername, "stagingRegistryUsername", "",
		"Username to download internal images from the staging docker registry. Alternatively you may use STAGING_REGISTRY_USERNAME")
	cliCmd.PersistentFlags().StringVar(&config.Environment.StagingRegistryPassword, "stagingRegistryPassword", "",
		"Password to download internal images from the staging docker registry. Alternatively you may use STAGING_REGISTRY_PASSWORD")
	cliCmd.PersistentFlags().StringVar(&config.Environment.StagingRegistryURL, "stagingRegistryURL", "",
		"URL of the staging docker registry. Alternatively you may use STAGING_REGISTRY_URL")
	// Development
	cliCmd.PersistentFlags().StringVar(&config.Environment.DevRegistryUsername, "devRegistryUsername", "",
		"Username to download internal images from the development docker registry. Alternatively you may use DEV_REGISTRY_USERNAME")
	cliCmd.PersistentFlags().StringVar(&config.Environment.DevRegistryPassword, "devRegistryPassword", "",
		"Password to download internal images from the development docker registry. Alternatively you may use DEV_REGISTRY_PASSWORD")
	cliCmd.PersistentFlags().StringVar(&config.Environment.DevRegistryURL, "devRegistryURL", "",
		"URL of the development docker registry. Alternatively you may use DEV_REGISTRY_URL")
	// Public
	cliCmd.PersistentFlags().StringVar(&config.Environment.PublicRegistryUsername, "publicRegistryUsername", "",
		"Username to download internal images from the public docker registry. Alternatively you may use PUBLIC_REGISTRY_USERNAME")
	cliCmd.PersistentFlags().StringVar(&config.Environment.PublicRegistryPassword, "publicRegistryPassword", "",
		"Password to download internal images from the public docker registry. Alternatively you may use PUBLIC_REGISTRY_PASSWORD")
	cliCmd.PersistentFlags().StringVar(&config.Environment.PublicRegistryURL, "publicRegistryURL", "",
		"URL of the public docker registry. Alternatively you may use PUBLIC_REGISTRY_URL")
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

	addRegistryOptions(runCmd)

	runCmd.PersistentFlags().StringVar(&config.ZTPlanetSecretPath, "ztPlanetSecretPath", "",
		"Path of the ZeroTier Planet secret file")

	rootCmd.AddCommand(runCmd)
}
