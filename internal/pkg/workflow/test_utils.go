/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package workflow

import (
	"fmt"
	"github.com/nalej/installer/internal/pkg/entities"

	"github.com/nalej/derrors"
	"github.com/nalej/grpc-infrastructure-go"
	"github.com/nalej/grpc-installer-go"
)

// WorkflowResult structure for testing the callback of a workflow execution.
type WorkflowResult struct {
	Called bool
	Error  derrors.Error
	State  WorkflowState
}

// NewWorkflowResult creates a WorkflowResult.
func NewWorkflowResult() *WorkflowResult {
	return &WorkflowResult{false, nil, InitState}
}

// Finished returns true if the workflow result received the callback.
func (wr *WorkflowResult) Finished() bool {
	return wr.Called
}

// Callback function.
func (wr *WorkflowResult) Callback(workflowID string, error derrors.Error,
	state WorkflowState) {
	wr.Error = error
	wr.Called = true
	wr.State = state
}

func GetTestParameters(numNodes int, appClusterInstall bool) *Parameters {
	nodes := make([]string, 0)
	for i := 0; i < numNodes; i++ {
		toAdd := fmt.Sprintf("10.1.1.%d", i)
		nodes = append(nodes, toAdd)
	}
	request := grpc_installer_go.InstallRequest{
		InstallId:         "TestInstall",
		OrganizationId:    "orgID",
		ClusterId:         "TestCluster",
		ClusterType:       grpc_infrastructure_go.ClusterType_KUBERNETES,
		InstallBaseSystem: false,
		KubeConfigRaw:     "KubeConfigContent",
		Nodes:             nodes,
	}

	assets := NewAssets(make([]string, 0), make([]string, 0))
	paths := NewPaths("assestPath", "binPath", "confPath")
	registryCredentials := NewRegistryCredentials("Production", "dockerUsername", "dockerPassword", "registryURL")
	return NewParameters(
		request,
		*assets,
		*paths,
		"mngtcluster_host", "80",
		"dns_host", "53",
		entities.Production,
		appClusterInstall,
		[]RegistryCredentials{*registryCredentials})
}
