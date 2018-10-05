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
	Credentials InstallCredentials `json:"credentials"`
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

// TODO Remove assets if not used anymore
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
	// TempPath contains the path of the temporal files used for the installs.
	TempPath string `json:"tempPath"`
}

func NewPaths(componentsPath string, binaryPath string, tempPath string) *Paths {
	return &Paths{componentsPath, binaryPath, tempPath}
}

type InstallCredentials struct {
	// Username for the SSH credentials.
	Username string `json:"username"`
	// PrivateKeyPath with the path of the private key.
	PrivateKeyPath string `json:"privateKeyPath"`
	// KubeConfigPath with the path of the kubeconfig file
	KubeConfigPath string `json:"kubeConfigPath"`
	// RemoveCredentials indicates that the credentials files must be removed after the installation.
	RemoveCredentials bool `json:"removeCredentials"`
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
		InstallCredentials{},
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

	if p.Credentials.Username == "" && p.Credentials.PrivateKeyPath == "" && p.Credentials.KubeConfigPath == "" {
		return derrors.NewInternalError("credentials have not been loaded. Call LoadCredentials() before Validate()")
	}

	return nil
}

// writeTempFile writes a content to a temporal file
func (p * Parameters) writeTempFile(content string, prefix string) (*string, derrors.Error) {
	tmpfile, err := ioutil.TempFile(p.Paths.TempPath, prefix)
	if err != nil {
		return nil, derrors.AsError(err, "cannot create temporal file")
	}
	_, err = tmpfile.Write([]byte(content))
	if err != nil {
		return nil, derrors.AsError(err, "cannot write temporal file")
	}
	err = tmpfile.Close()
	if err != nil {
		return nil, derrors.AsError(err, "cannot close temporal file")
	}
	tmpName := tmpfile.Name()
	return &tmpName, nil
}

func (p * Parameters) LoadCredentials() derrors.Error {

	p.Credentials.Username = p.InstallRequest.Username

	if p.InstallRequest.PrivateKey != "" {
		f, err := p.writeTempFile(p.InstallRequest.PrivateKey, "pk")
		if err != nil {
			return err
		}
		p.Credentials.PrivateKeyPath = *f
	}

	if p.InstallRequest.KubeConfigRaw != "" {
		f, err := p.writeTempFile(p.InstallRequest.KubeConfigRaw, "kc")
		if err != nil {
			return err
		}
		p.Credentials.KubeConfigPath = *f
	}

	return nil
}


