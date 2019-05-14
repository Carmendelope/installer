/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package installer

import (
	"fmt"
	"github.com/nalej/installer/internal/pkg/entities"
	"time"

	"github.com/nalej/derrors"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-installer-go"
	"github.com/nalej/installer/internal/pkg/templates"
	"github.com/nalej/installer/internal/pkg/utils"
	"github.com/nalej/installer/internal/pkg/workflow"
	"github.com/rs/zerolog/log"
)

type Installer struct {
	Params   workflow.Parameters
	Workflow *workflow.Workflow
}

func NewInstallerFromCLI(
	installId string,
	installK8s bool,
	kubeConfigPath string,
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
) (*Installer, derrors.Error) {

	kubeConfigContent, err := utils.GetKubeConfigContent(kubeConfigPath)
	if err != nil {
		return nil, err
	}

	privateKeyContent, err := utils.GetPrivateKeyContent(privateKeyPath)
	if err != nil {
		return nil, err
	}

	staticIPAddresses := grpc_installer_go.StaticIPAddresses{
		UseStaticIp: useStaticIPAddresses,
		Ingress:     ipAddressIngress,
		Dns:         ipAddressDNS,
		CorednsExt: ipAddressCoreDNS,
		VpnServer: ipAddressVPNServer,
	}

	request := grpc_installer_go.InstallRequest{
		InstallId:         installId,
		OrganizationId:    "nalej",
		ClusterId:         "nalej-management-cluster",
		ClusterType:       grpc_infrastructure_go.ClusterType_KUBERNETES,
		InstallBaseSystem: installK8s,
		KubeConfigRaw:     kubeConfigContent,
		Hostname:          managementClusterHost,
		Username:          username,
		PrivateKey:        privateKeyContent,
		Nodes:             nodes,
		TargetPlatform:    grpc_installer_go.Platform(grpc_installer_go.Platform_value[targetPlatform]),
		StaticIpAddresses: &staticIPAddresses,
	}

	registryCredentials := workflow.NewRegistryCredentialsFromEnvironment(environment)

	params := workflow.NewParameters(request, workflow.Assets{},
		paths, managementClusterHost, workflow.DefaultManagementPort,
		dnsClusterHost, dnsClusterPort,
		environment.Target,
		appClusterInstall,
		registryCredentials,
	*workflow.EmptyNetworkConfig, "")
	return NewInstaller(*params), nil
}

func NewInstaller(params workflow.Parameters) *Installer {
	return &Installer{
		Params: params,
	}
}

func (i *Installer) logListener(msg string) {
	log.Info().Msg(msg)
}

// Load all the credentials and associated workflow into the installer.
func (i *Installer) Load() {
	i.exitOnError(i.Params.LoadCredentials())
	i.exitOnError(i.Params.Validate())
	p := workflow.NewParser()
	workflow, err := p.ParseWorkflow("cli-install", templates.InstallManagementCluster, "install-management-cluster", i.Params)
	i.exitOnError(err)
	i.Workflow = workflow
}

func (i *Installer) exitOnError(err derrors.Error) {
	if err != nil {
		log.Panic().Str("error", err.DebugReport()).Msg("installation exit with error")
	}
}

// Run the install process.
func (i *Installer) Run() {
	i.Load()
	wr := &workflow.WorkflowResult{}
	execHandler := workflow.GetExecutorHandler()
	exec, err := execHandler.Add(i.Workflow, wr.Callback)
	i.exitOnError(err)
	exec.SetLogListener(i.logListener)
	start := time.Now()
	exec, err = execHandler.Execute(i.Workflow.WorkflowID)
	i.exitOnError(err)
	checks := 0
	for !wr.Called {
		time.Sleep(time.Second * 15)
		if checks%4 == 0 {
			if i.Params.AppClusterInstall {
				fmt.Println("AppCluster installation", string(exec.State), "-", time.Since(start).String())
			} else {
				fmt.Println("Management cluster installation", string(exec.State), "-", time.Since(start).String())
			}
		}
		checks++
	}
	elapsed := time.Since(start)
	fmt.Println("Installation took ", elapsed)
	if wr.Error != nil {
		fmt.Println("Installation failed due to ", wr.Error.Error())
		log.Fatal().Str("error", wr.Error.DebugReport()).Msg("Installation failed")
	}
}
