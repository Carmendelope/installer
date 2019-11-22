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

package installer_cli

import (
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-installer-go"
	"github.com/nalej/installer/internal/pkg/entities"
	"github.com/nalej/installer/internal/pkg/templates"
	"github.com/nalej/installer/internal/pkg/utils"
	"github.com/nalej/installer/internal/pkg/workflow"
	"github.com/rs/zerolog/log"
	"time"
)

// CLI structure with methods and constructs shared by different CLI commands.
type CLI struct {
	// Params used in the command.
	Params workflow.Parameters
	// Workflow to be executed.
	Workflow *workflow.Workflow
	// kubeConfigContent with the raw contents of the kubeConfig file to be used.
	kubeConfigContent string
}

// NewCLI builds a new CLI command wrapper to interact with the underlying installer logic.
func NewCLI(kubeConfigPath string) (*CLI, derrors.Error) {
	kubeConfigContent, err := utils.GetKubeConfigContent(kubeConfigPath)
	if err != nil {
		return nil, err
	}
	return &CLI{kubeConfigContent: kubeConfigContent}, nil
}

// PrepareInstallCommand prepares the CLI to execute an install command.
func (c *CLI) PrepareInstallCommand(
	installId string,
	installK8s bool,
	username string,
	privateKeyPath string,
	nodes []string,
	targetPlatform string,
	paths workflow.Paths,
	managementClusterHost string,
	dnsClusterHost string,
	dnsClusterPort string,
	useStaticIPAddresses bool,
	ipAddressIngress string,
	ipAddressDNS string,
	ipAddressCoreDNS string,
	ipAddressVPNServer string,
	appClusterInstall bool,
	environment entities.Environment,
) {
	// load the private key content if required.
	privateKeyContent, err := utils.GetPrivateKeyContent(privateKeyPath)
	c.exitOnError(err)

	// set the static IP addresses.
	staticIPAddresses := grpc_installer_go.StaticIPAddresses{
		UseStaticIp: useStaticIPAddresses,
		Ingress:     ipAddressIngress,
		Dns:         ipAddressDNS,
		CorednsExt:  ipAddressCoreDNS,
		VpnServer:   ipAddressVPNServer,
	}
	// Prepare the gRPC request as would have been send to the service.
	request := &grpc_installer_go.InstallRequest{
		InstallId:         installId,
		OrganizationId:    "nalej",
		ClusterId:         "nalej-management-cluster",
		ClusterType:       grpc_infrastructure_go.ClusterType_KUBERNETES,
		InstallBaseSystem: installK8s,
		KubeConfigRaw:     c.kubeConfigContent,
		Hostname:          managementClusterHost,
		Username:          username,
		PrivateKey:        privateKeyContent,
		Nodes:             nodes,
		TargetPlatform:    grpc_installer_go.Platform(grpc_installer_go.Platform_value[targetPlatform]),
		StaticIpAddresses: &staticIPAddresses,
	}
	params := workflow.NewInstallParameters(request, workflow.Assets{},
		paths, managementClusterHost, workflow.DefaultManagementPort,
		dnsClusterHost, dnsClusterPort,
		environment.Target,
		appClusterInstall,
		*workflow.EmptyNetworkConfig, "", "")

	c.Params = *params

}

// PrepareUninstallCommand prepares the CLI to execute an uninstall command.
func (c *CLI) PrepareUninstallCommand(
	requestID string,
	organizationID string,
	clusterID string,
	targetPlatform string,
	appCluster bool,
) {
	// Prepare the gRPC request as would have been send to the service.
	request := &grpc_installer_go.UninstallClusterRequest{
		RequestId:      requestID,
		OrganizationId: organizationID,
		ClusterId:      clusterID,
		ClusterType:    grpc_infrastructure_go.ClusterType_KUBERNETES,
		KubeConfigRaw:  c.kubeConfigContent,
		TargetPlatform: grpc_installer_go.Platform(grpc_installer_go.Platform_value[targetPlatform]),
	}
	params := workflow.NewUninstallParameters(request, appCluster)
	c.Params = *params
}

// Load all the credentials and associated workflow into the installer.
func (c *CLI) loadCredentials() {
	c.exitOnError(c.Params.LoadCredentials())
	c.exitOnError(c.Params.Validate())
	p := workflow.NewParser()
	workflowTemplate := ""
	workflowName := ""
	if c.Params.InstallRequest != nil {
		workflowName = "installCluster"
		workflowTemplate = templates.InstallManagementCluster
	} else if c.Params.UninstallRequest != nil {
		workflowName = "uninstallCluster"
		workflowTemplate = templates.UninstallCluster
	}
	workflow, err := p.ParseWorkflow("cli-install", workflowTemplate, workflowName, c.Params)
	c.exitOnError(err)
	c.Workflow = workflow
}

// exitOnError produces a panic if an error is passed as parameter to finish the execution.
func (c *CLI) exitOnError(err derrors.Error) {
	if err != nil {
		log.Panic().Str("error", err.DebugReport()).Msg("installer-cli exit with error")
	}
}

// logListener receives messages produced by the running workflow.
func (c *CLI) logListener(msg string) {
	log.Info().Msg(msg)
}

// Execute the install/uninstall process.
func (c *CLI) Execute() {
	c.loadCredentials()
	wr := &workflow.WorkflowResult{}
	execHandler := workflow.GetExecutorHandler()
	exec, err := execHandler.Add(c.Workflow, wr.Callback)
	c.exitOnError(err)
	exec.SetLogListener(c.logListener)
	start := time.Now()
	exec, err = execHandler.Execute(c.Workflow.WorkflowID)
	c.exitOnError(err)
	checks := 0
	operation := ""
	if c.Params.InstallRequest != nil {
		if c.Params.AppCluster {
			operation = "Installing application cluster"
		} else {
			operation = "Installing management cluster"
		}
	} else if c.Params.UninstallRequest != nil {
		if c.Params.AppCluster {
			operation = "Uninstalling application cluster"
		} else {
			operation = "Uninstalling management cluster"
		}
	}
	for !wr.Called {
		time.Sleep(time.Second * 15)
		if checks%4 == 0 {
			fmt.Println(operation, string(exec.State), "-", time.Since(start).String())
		}
		checks++
	}
	elapsed := time.Since(start)
	fmt.Println("Operation took ", elapsed)
	if wr.Error != nil {
		fmt.Println("Operation failed due to ", wr.Error.Error())
		log.Fatal().Str("error", wr.Error.DebugReport()).Msg(fmt.Sprintf("%s failed", operation))
	}
}
