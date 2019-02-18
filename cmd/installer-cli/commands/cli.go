/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package commands

import (
	"github.com/nalej/installer/internal/pkg/entities"
	"os"

	"github.com/nalej/derrors"
	"github.com/nalej/installer/internal/pkg/utils"
	"github.com/nalej/installer/internal/pkg/workflow"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var explainPlan bool

var installKubernetes bool
var kubeConfigPath string
var username string
var privateKeyPath string
var nodes string
var targetPlatform string

var managementPublicHost string

var useStaticIPAddresses bool
var ipAddressIngress string
var ipAddressDNS string

var dnsClusterHost string
var dnsClusterPort int

var componentsPath string
var binaryPath string
var confPath string
var tempPath string

var environment entities.Environment

var cliCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the Nalej platform",
	Long:  `Install the Nalej platform`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		cmd.Help()
	},
}

func init() {
	cliCmd.PersistentFlags().BoolVar(&explainPlan, "explainPlan", false,
		"Show install plan instead of performing the install")
	cliCmd.PersistentFlags().BoolVar(&installKubernetes, "installK8s", false,
		"Whether kubernetes should be installed")
	cliCmd.PersistentFlags().StringVar(&kubeConfigPath, "kubeConfigPath", "~/.kube/config",
		"Specify the Kubernetes config path")
	cliCmd.PersistentFlags().StringVar(&username, "username", "",
		"Specify the username to connect to the remote machines (Only if installK8s is selected)")
	cliCmd.PersistentFlags().StringVar(&privateKeyPath, "privateKeyPath", "~/.ssh/id_rsa",
		"Specify the private key path to connect to the remote machine (Only if installK8s is selected)")
	cliCmd.PersistentFlags().StringVar(&nodes, "nodes", "",
		"List of IPs of the nodes to be installed separated by comma (Only if installK8s is selected)")
	cliCmd.PersistentFlags().StringVar(&targetPlatform, "targetPlatform", "MINIKUBE", "Target platform: MINIKUBE or AZURE")
	cliCmd.PersistentFlags().StringVar(&managementPublicHost, "managementClusterPublicHost", "",
		"Public FQDN where the management cluster is reachable by the application clusters")
	cliCmd.MarkPersistentFlagRequired("managementClusterPublicHost")

	cliCmd.PersistentFlags().BoolVar(&useStaticIPAddresses, "useStaticIPAddresses", false,
		"Use statically assigned IP Addresses for the public facing services")
	cliCmd.PersistentFlags().StringVar(&ipAddressIngress, "ipAddressIngress", "",
		"Public IP Address assigned to the public ingress service")
	cliCmd.PersistentFlags().StringVar(&ipAddressDNS, "ipAddressDNS", "",
		"Public IP Address assigned to the DNS server service")

	cliCmd.PersistentFlags().StringVar(&dnsClusterHost, "dnsClusterPublicHost", "",
		"Public FQDN where the management cluster is reachable for DNS requests by the application clusters")
	cliCmd.MarkPersistentFlagRequired("dnsClusterPublicHost")
	cliCmd.PersistentFlags().IntVar(&dnsClusterPort, "dnsClusterPublicPort", 53,
		"Public port where the management cluster is reachable for DNS request by the application clusters")

	cliCmd.PersistentFlags().StringVar(&componentsPath, "componentsPath", "./assets/",
		"Directory with the components to be installed")
	cliCmd.PersistentFlags().StringVar(&binaryPath, "binaryPath", "./bin/",
		"Directory with the binary executables")
	cliCmd.PersistentFlags().StringVar(&confPath, "confPath", "./conf/",
		"Directory with the configuration files")
	cliCmd.PersistentFlags().StringVar(&tempPath, "tempPath", "./temp/",
		"Directory to store temporal files")

	addRegistryOptions(cliCmd)

	rootCmd.AddCommand(cliCmd)
}

// Add parameters related to the usage of registries.
func addRegistryOptions(cliCmd *cobra.Command){
	cliCmd.PersistentFlags().StringVar(&environment.TargetEnvironment, "targetEnvironment", "PRODUCTION", "Target environment to be installed: PRODUCTION, STAGING, or DEVELOPMENT")
	// Production
	cliCmd.PersistentFlags().StringVar(&environment.ProdRegistryUsername, "prodRegistryUsername", "",
		"Username to download internal images from the production docker registry. Alternatively you may use PROD_REGISTRY_USERNAME")
	cliCmd.PersistentFlags().StringVar(&environment.ProdRegistryPassword, "prodRegistryPassword", "",
		"Password to download internal images from the production docker registry. Alternatively you may use PROD_REGISTRY_PASSWORD")
	cliCmd.PersistentFlags().StringVar(&environment.ProdRegistryURL, "prodRegistryURL", "",
		"URL of the production docker registry. Alternatively you may use PROD_REGISTRY_URL")
	// Staging
	cliCmd.PersistentFlags().StringVar(&environment.StagingRegistryUsername, "stagingRegistryUsername", "",
		"Username to download internal images from the staging docker registry. Alternatively you may use STAGING_REGISTRY_USERNAME")
	cliCmd.PersistentFlags().StringVar(&environment.StagingRegistryPassword, "stagingRegistryPassword", "",
		"Password to download internal images from the staging docker registry. Alternatively you may use STAGING_REGISTRY_PASSWORD")
	cliCmd.PersistentFlags().StringVar(&environment.StagingRegistryURL, "stagingRegistryURL", "",
		"URL of the staging docker registry. Alternatively you may use STAGING_REGISTRY_URL")
	// Development
	cliCmd.PersistentFlags().StringVar(&environment.DevRegistryUsername, "devRegistryUsername", "",
		"Username to download internal images from the development docker registry. Alternatively you may use DEV_REGISTRY_USERNAME")
	cliCmd.PersistentFlags().StringVar(&environment.DevRegistryPassword, "devRegistryPassword", "",
		"Password to download internal images from the development docker registry. Alternatively you may use DEV_REGISTRY_PASSWORD")
	cliCmd.PersistentFlags().StringVar(&environment.DevRegistryURL, "devRegistryURL", "",
		"URL of the development docker registry. Alternatively you may use DEV_REGISTRY_URL")
	// Public
	cliCmd.PersistentFlags().StringVar(&environment.PublicRegistryUsername, "publicRegistryUsername", "",
		"Username to download internal images from the public docker registry. Alternatively you may use PUBLIC_REGISTRY_USERNAME")
	cliCmd.PersistentFlags().StringVar(&environment.PublicRegistryPassword, "publicRegistryPassword", "",
		"Password to download internal images from the public docker registry. Alternatively you may use PUBLIC_REGISTRY_PASSWORD")
	cliCmd.PersistentFlags().StringVar(&environment.PublicRegistryURL, "publicRegistryURL", "",
		"URL of the public docker registry. Alternatively you may use PUBLIC_REGISTRY_URL")
}


func CheckExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func GetPaths() (*workflow.Paths, derrors.Error) {

	components := utils.ExtendComponentsPath(utils.GetPath(componentsPath), false)
	binary := utils.GetPath(binaryPath)
	temp := utils.GetPath(tempPath)

	if !CheckExists(components) {
		return nil, derrors.NewNotFoundError("components directory does not exist").WithParams(components)
	}

	if !CheckExists(binary) {
		return nil, derrors.NewNotFoundError("binary directory does not exists").WithParams(binary)
	}

	if !CheckExists(temp) {
		err := os.MkdirAll(temp, os.ModePerm)
		if err != nil {
			return nil, derrors.AsError(err, "cannot create temp directory")
		}
	}

	log.Info().Str("path", components).Msg("Components")
	log.Info().Str("path", binary).Msg("Binaries")
	log.Info().Str("path", temp).Msg("Temporal files")

	return &workflow.Paths{
		ComponentsPath: components,
		BinaryPath:     binary,
		TempPath:       temp,
	}, nil
}

func ValidateInstallParameters() derrors.Error {
	if installKubernetes {
		if username == "" || privateKeyPath == "" {
			return derrors.NewInvalidArgumentError("username and privateKeyPath expected on kubernetes install mode")
		}
		if nodes == "" {
			return derrors.NewInvalidArgumentError("nodes expected on kubernetes install mode")
		}
	} else {
		if kubeConfigPath == "" {
			return derrors.NewInvalidArgumentError("kubeConfig path expected")
		}
	}
	log.Info().Bool("set", installKubernetes).Msg("Install Kubernetes")
	if installKubernetes {
		log.Info().Str("value", username).Msg("Username")
		log.Info().Str("path", privateKeyPath).Msg("Private Key")
	}
	log.Info().Str("path", kubeConfigPath).Msg("KubeConfig")
	return nil
}
