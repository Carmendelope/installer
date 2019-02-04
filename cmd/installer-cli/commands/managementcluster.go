/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package commands

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/nalej/installer/cmd/installer-cli/commands/installer"
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

func getDockerUsername() string {
	if dockerRegistryUsername != "" {
		return dockerRegistryUsername
	}
	fromEnv := os.Getenv("DOCKER_USER")
	if fromEnv != "" {
		return fromEnv
	}
	log.Fatal().Msg("docker username must be set either by parameter or as env variable")
	return ""
}

func getDockerPassword() string {
	if dockerRegistryPassword != "" {
		return dockerRegistryPassword
	}
	fromEnv := os.Getenv("DOCKER_PASSWORD")
	if fromEnv != "" {
		return fromEnv
	}
	log.Fatal().Msg("docker password must be set either by parameter or as env variable")
	return ""
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

	inst, err := installer.NewInstallerFromCLI("cli-install",
		installKubernetes,
		kubeConfigPath,
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
		false,
		getDockerUsername(), getDockerPassword())

	if err != nil {
		log.Panic().Str("error", err.DebugReport()).Msg("cannot generate installer")
	}

	inst.Load()

	if explainPlan {
		fmt.Println(inst.Workflow.PrettyPrint())
	} else {
		inst.Run()
	}

}
