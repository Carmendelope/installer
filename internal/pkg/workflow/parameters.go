/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package workflow

import (
	"encoding/json"
	"github.com/nalej/installer/internal/pkg/entities"
	"io/ioutil"

	"github.com/nalej/installer/internal/pkg/errors"

	"github.com/nalej/derrors"
	"github.com/nalej/grpc-installer-go"
)

const DefaultManagementPort = "443"

// Parameters required to transform a template into a workflow.
type Parameters struct {
	// InstallRequest with the details of the installation to be performed.
	InstallRequest grpc_installer_go.InstallRequest
	// Credentials required for the installation of the cluster.
	Credentials InstallCredentials `json:"credentials"`
	// Assets to be installed
	Assets Assets `json:"assets"`
	// Paths defines a set of paths for assets, binaries, and configuration.
	Paths Paths `json:"paths"`
	//ManagementClusterHost is the host where the management cluster is accepting callback requests.
	ManagementClusterHost string `json:"management_cluster_host"`
	//InframgrPort is the port where the management cluster is accepting requests.
	ManagementClusterPort string `json:"management_cluster_port"`
	// DNSClusterHost is the host where the dns service of the management cluster is accepting DNS requests.
	DNSClusterHost string `json:"dns_cluster_host"`
	// DNSClusterPort is the port where the dns service of the management cluster is accepting DNS requests.
	DNSClusterPort string `json:"dns_cluster_port"`
	// TargetEnvironment defines the type of environment being installed: PRODUCTION, STAGING, DEVELOPMENT
	TargetEnvironment string `json:"target_environment"`
	//AppClusterInstall indicates if an application cluster is being installed.
	AppClusterInstall bool `json:"app_cluster_install"`
	//Registries contains the credentials to access the available docker registries to deploy internal images.
	Registries []RegistryCredentials `json:"registries"`
	// NetworkConfig contains the configuration of the networking of the cluster.
	NetworkConfig NetworkConfig `json:"network_config"`
	// AuthSecret contains the secret required to validate JWT tokens.
	AuthSecret string `json:"auth_secret"`
	// CACertPath contains the path to the certificate of a TLS secret
	CACertPath string `json:"ca_cert_path"`
}

var EmptyNetworkConfig = &NetworkConfig{}

// Deprecated: This will be removed as ZT will be removed
type NetworkConfig struct {
	// ZT Planet Secret
	ZTPlanetSecretPath string `json:"zt_planet_secret_path"`
}

// Deprecated: This will be removed as ZT will be removed.
func NewNetworkConfig (ztPlanetSecretPath string) *NetworkConfig {
	return &NetworkConfig{
		ZTPlanetSecretPath: ztPlanetSecretPath,
	}
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

type RegistryCredentials struct {
	// Name of the registry
	Name string `json:"name"`
	//Username to access the docker registry
	Username string `json:"username"`
	//Password to access the docker registry
	Password string `json:"password"`
	// URL of the registry
	URL string `json:"url"`
}

func NewRegistryCredentials(name string, username string, password string, url string) *RegistryCredentials {
	return &RegistryCredentials{
		Name: name,
		Username: username,
		Password: password,
		URL: url,
	}
}

func NewRegistryCredentialsFromEnvironment(env entities.Environment) []RegistryCredentials{
	result := make([]RegistryCredentials, 0)
	prod := NewRegistryCredentials(entities.ProdRegistryName,
		env.ProdRegistryUsername, env.ProdRegistryPassword, env.ProdRegistryURL)
	result = append(result, *prod)
	if env.Target == entities.Staging || env.Target == entities.Development{
		staging := NewRegistryCredentials(entities.StagingRegistryName,
			env.StagingRegistryUsername, env.StagingRegistryPassword, env.StagingRegistryURL)
		result = append(result, *staging)
	}
	if env.Target == entities.Development{
		development := NewRegistryCredentials(entities.DevRegistryName,
			env.DevRegistryUsername, env.DevRegistryPassword, env.DevRegistryURL)
		result = append(result, *development)
	}

	public := NewRegistryCredentials(entities.PublicRegistryName,
		env.PublicRegistryUsername, env.PublicRegistryPassword, env.PublicRegistryURL)
	result = append(result, *public)
	return result
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
	managementClusterHost string,
	managementClusterPort string,
	dnsClusterHost string,
	dnsClusterPort string,
	targetEnvironment entities.TargetEnvironment,
	appClusterInstall bool,
	registryCredentials []RegistryCredentials,
	networkConfig NetworkConfig,
	authxSecret string,
	caCertPath string,
	) *Parameters {
	return &Parameters{
		request,
		InstallCredentials{},
		assets,
		paths,
		managementClusterHost, managementClusterPort,
		dnsClusterHost, dnsClusterPort,
		entities.TargetEnvironmentToString[targetEnvironment],
		appClusterInstall,
		registryCredentials,
		networkConfig,
		authxSecret,
		caCertPath,
	}
}

// NewParametersFromFile extract a parameters object from a file.
func NewParametersFromFile(filePath string) (*Parameters, derrors.Error) {
	content, err := ioutil.ReadFile(filePath)
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
	if p.Credentials.Username == "" && p.Credentials.PrivateKeyPath == "" && p.Credentials.KubeConfigPath == "" {
		return derrors.NewInternalError("credentials have not been loaded. Call LoadCredentials() before Validate()")
	}

	if p.Credentials.KubeConfigPath == "" && len(p.InstallRequest.Nodes) == 0 {
		return derrors.NewInternalError(errors.InvalidNumMaster)
	}

	if len(p.Registries) == 0{
		return derrors.NewInvalidArgumentError("at least one registry must be specified")
	}

	return nil
}

// writeTempFile writes a content to a temporal file
func (p *Parameters) writeTempFile(content string, prefix string) (*string, derrors.Error) {
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

func (p *Parameters) LoadCredentials() derrors.Error {

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
