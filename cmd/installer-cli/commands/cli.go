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

package commands

import (
	"github.com/nalej/derrors"
	"github.com/nalej/installer/cmd/installer-cli/commands/installer"
	"github.com/nalej/installer/internal/pkg/workflow"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"
)

var explainPlan bool

var installKubernetes bool
var kubeConfigPath string
var username string
var privateKeyPath string
var nodes string

var componentsPath string
var binaryPath string
var confPath string
var tempPath string

var cliCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the Nalej platform",
	Long:  `Install the Nalej platform`,
	Run: func(cmd *cobra.Command, args []string) {
		SetupLogging()
		log.Info().Msg("Launching API!")
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

	cliCmd.PersistentFlags().StringVar(&componentsPath, "componentsPath", "./assets/",
		"Directory with the components to be installed")
	cliCmd.PersistentFlags().StringVar(&binaryPath, "binaryPath", "./bin/",
		"Directory with the binary executables")
	cliCmd.PersistentFlags().StringVar(&confPath, "confPath", "./conf/",
		"Directory with the configuration files")
	cliCmd.PersistentFlags().StringVar(&tempPath, "tempPath", "./temp/",
		"Directory to store temporal files")

	rootCmd.AddCommand(cliCmd)
}

func CheckExists(path string) bool {
	_, err := os.Stat(path);
	return !os.IsNotExist(err)
}

func GetPaths() (* workflow.Paths, derrors.Error) {

	components := installer.GetPath(componentsPath)
	binary := installer.GetPath(binaryPath)
	temp := installer.GetPath(tempPath)

	if !CheckExists(components) {
		return nil, derrors.NewNotFoundError("components directory does not exist").WithParams(components)
	}

	if !CheckExists(binary) {
		return nil, derrors.NewNotFoundError("binary directory does not exists").WithParams(binary)
	}

	if ! CheckExists(temp) {
		err := os.MkdirAll(temp, os.ModePerm)
		if err != nil{
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
	}else {
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

