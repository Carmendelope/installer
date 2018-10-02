/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package rke

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/nalej/installer/internal/pkg/errors"
	"github.com/rs/zerolog/log"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"

	"github.com/nalej/derrors"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
	"github.com/nalej/installer/internal/pkg/workflow/handler"

)

// RKEInstall structure defining the fields required to install a cluster using RKE.
type RKEInstall struct {
	entities.GenericSyncCommand
	RkeBinaryPath string `json:"rkeBinaryPath"`
	ClusterConfig
	KubeConfigOutputPath string `json:"kubeConfigOutputPath"`
	installTemplate      string
}

// NewRKEInstall create a new command with all parameters.
func NewRKEInstall(
	rkeBinaryPath string,
	clusterConfig ClusterConfig,
	kubeConfigOutputPath string,
	installTemplate string) *RKEInstall {
	return &RKEInstall{
		*entities.NewSyncCommand(entities.RKEInstall),
		rkeBinaryPath,
		clusterConfig, kubeConfigOutputPath, installTemplate}
}

// NewRKEInstallFromJSON creates a RKE Install command from a JSON object.
func NewRKEInstallFromJSON(raw []byte) (*entities.Command, derrors.Error) {
	f := &RKEInstall{}
	if err := json.Unmarshal(raw, &f); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	f.CommandID = entities.GenerateCommandID(f.Name())
	var r entities.Command = f
	return &r, nil
}

// getTemplate returns the template to be used for the installation process. If empty, the default one will be used.
func (cmd *RKEInstall) getTemplate() string {
	if cmd.installTemplate != "" {
		return cmd.installTemplate
	}
	return ClusterTemplate
}

// CreateClusterConfig generate the RKE cluster.yaml file using the installation parameters.
func (cmd *RKEInstall) CreateClusterConfig() (string, derrors.Error) {
	template := NewRKETemplate(cmd.getTemplate())
	config := cmd.ClusterConfig
	yamlString, err := template.ParseTemplate(&config)
	if err != nil {
		return "", err
	}
	clusterFile, createErr := ioutil.TempFile("", "cluster.yaml")
	if createErr != nil {
		return "", derrors.AsError(createErr, errors.IOError)
	}
	if _, writeErr := clusterFile.Write([]byte(yamlString)); writeErr != nil {
		return "", derrors.AsError(writeErr, errors.IOError)
	}
	clusterFile.Close()
	log.Debug().Str("output file", clusterFile.Name()).Msg("Temporal cluster.yaml stored")
	return clusterFile.Name(), nil
}

// copyToLog copies a reader output to the associated command handler.
func (cmd *RKEInstall) copyToLog(commandHandler handler.CommandHandler, r io.Reader) {
	output := bufio.NewReader(r)
	for {
		line, err := output.ReadString('\n')
		commandHandler.AddLogEntry(cmd.CommandID, strings.TrimSpace(line))
		if err != nil {
			break
		}
	}
}

func (cmd *RKEInstall) copyKubeConfig(clusterConfigFile string) (*entities.CommandResult, derrors.Error) {

	targetName := fmt.Sprintf("kube_config_%s", path.Base(clusterConfigFile))
	kubeFromFile := fmt.Sprintf("%s/%s", path.Dir(clusterConfigFile), targetName)
	log.Debug().Str("kubeConfig", kubeFromFile).Msg("Target kube_config")

	from, err := os.Open(kubeFromFile)
	if err != nil {
		return nil, derrors.AsError(err, errors.IOError)
	}
	defer from.Close()

	kubeToFile := fmt.Sprintf("%s/kube_config_%s_%s.yml", cmd.KubeConfigOutputPath, cmd.ClusterName, cmd.TargetNodes[0])
	to, err := os.OpenFile(kubeToFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return nil, derrors.AsError(err, errors.IOError)
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		return nil, derrors.AsError(err, errors.IOError)
	}
	log.Info().Str("NewKubeConfig", kubeToFile).Msg("KubeConfig available")
	return entities.NewCommandResult(true, "rke finished successfully", nil), nil
}

// Run triggers the execution of the command.
func (cmd *RKEInstall) Run(workflowID string) (*entities.CommandResult, derrors.Error) {
	clusterConfigPath, err := cmd.CreateClusterConfig()
	if err != nil {
		log.Warn().Err(err).Msg("unable to create cluster config")
		return nil, err
	}

	log.Debug().Str("path", cmd.RkeBinaryPath).Msg("RKE binary")
	rke := exec.Command(cmd.RkeBinaryPath, "up", "--config", clusterConfigPath)
	rkeOut, pipeErr := rke.StdoutPipe()
	if pipeErr != nil {
		return nil, derrors.AsError(pipeErr, errors.IOError)
	}

	rkeErr, pipeErr := rke.StderrPipe()
	if pipeErr != nil {
		return nil, derrors.AsError(pipeErr, errors.IOError)
	}

	var wg sync.WaitGroup
	commandHandler := handler.GetCommandHandler()
	log.Debug().Msg("Starting rke binary")
	if err := rke.Start(); err != nil {
		return nil, derrors.AsError(err, errors.OpFail)
	}

	wg.Add(2)
	go func() {
		defer wg.Done()
		cmd.copyToLog(commandHandler, rkeOut)
	}()
	go func() {
		defer wg.Done()
		cmd.copyToLog(commandHandler, rkeErr)
	}()

	// Wait for the stdout and stderr pipes to close.
	wg.Wait()
	// Wait for the command itself to close.
	if err := rke.Wait(); err != nil {
		return entities.NewCommandResult(false, "rke failed", derrors.AsError(err, errors.OpFail)), nil
	}
	return cmd.copyKubeConfig(clusterConfigPath)
}

// Obtain a string representation
func (cmd *RKEInstall) String() string {
	return fmt.Sprintf("SYNC RKE Install on %s", strings.Join(cmd.TargetNodes, ", "))
}

// PrettyPrint returns a simple space indexed string.
func (cmd *RKEInstall) PrettyPrint(indentation int) string {
	outputPath := strings.Repeat("  ", indentation) + fmt.Sprintf("  OutputPath: %s", cmd.KubeConfigOutputPath)
	binaryPath := strings.Repeat("  ", indentation) + fmt.Sprintf("  RKE binary: %s", cmd.RkeBinaryPath)
	return strings.Repeat(" ", indentation) + fmt.Sprintf("SYNC RKE Install on %s\n%s\n%s",
		strings.Join(cmd.TargetNodes, ", "), binaryPath, outputPath)
}

// UserString returns a simple string representation of the command for the user.
func (cmd *RKEInstall) UserString() string {
	return fmt.Sprintf("Installing Kubernetes on %s ", strings.Join(cmd.TargetNodes, ", "))
}
