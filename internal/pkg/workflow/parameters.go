/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package workflow

import (
	"encoding/json"
	"github.com/nalej/installer/internal/pkg/errors"
	"io/ioutil"

	"github.com/nalej/derrors"
	"github.com/nalej/grpc-installer-go"
)


// Parameters required to transform a template into a workflow.
type Parameters struct {
	InstallRequest grpc_installer_go.InstallRequest
	// Assets to be installed
	Assets Assets `json:"assets"`
	// Paths defines a set of paths for assets, binaries, and configuration.
	Paths Paths `json:"paths"`
	//InframgrHost is the host where the inframgr is accepting callback requests.
	InframgrHost string `json:"inframgrHost"`
	//InframgrPort is the port where the inframgr is accepting requests.
	InframgrPort string `json:"inframgrPort"`
	//AppClusterInstall indicates if an application cluster is being installed.
	AppClusterInstall bool `json:"appClusterInstall"`
}

type Assets struct {
	// Names is an array of asset names
	Names []string `json:"assets"`
	// Services is an array of the service associated with the assets
	Services []string `json:"services"`
}

func NewAssets(names []string, services []string) *Assets {
	return &Assets{names, services}
}

type Paths struct {
	// ComponentsPath inside the installer machine with the yamls files.
	ComponentsPath string `json:"componentsPath"`
	// BinaryPath contains the path for the auxiliar binaries to be executed (e.g., rke).
	BinaryPath string `json:"binaryPath"`
	// ConfPath contains the path of the configuration files.
	ConfPath string `json:"confPath"`
}

func NewPaths(assetsPath string, binaryPath string, confPath string) *Paths {
	return &Paths{assetsPath, binaryPath, confPath}
}

// EmptyParameters structure that can be used whenever no parameters are passed to the parser.
var EmptyParameters = Parameters{}

// NewParameters creates a Parameters structure.
func NewParameters(
	request grpc_installer_go.InstallRequest,
	assets Assets,
	paths Paths,
	inframgrHost string, appClusterInstall bool) *Parameters {
	return &Parameters{
		request,
		assets,
		paths,
		inframgrHost, "8860", appClusterInstall}
}

// NewParametersFromFile extract a parameters object from a file.
func NewParametersFromFile(filepath string) (*Parameters, derrors.Error) {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, derrors.NewInternalError(errors.CannotParseParameters, err)
	}
	parameters := new(Parameters)
	err = json.Unmarshal(content, parameters)
	if err != nil {
		return nil, derrors.NewInternalError(errors.CannotParseParameters, err)
	}
	return parameters, nil
}

// Validate checks the parameters to determine if the workflow can be executed.
func (p *Parameters) Validate() derrors.Error {
	if len(p.InstallRequest.Nodes) == 0 {
		return derrors.NewInternalError(errors.InvalidNumMaster)
	}
	return nil
}
