/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package config

import (
	"os"
	"strings"

	"github.com/nalej/derrors"
	"github.com/nalej/installer/internal/pkg/utils"
	"github.com/nalej/installer/version"
	"github.com/rs/zerolog/log"
)

type Config struct {
	// Address where the API service will listen requests.
	Port                   int
	ComponentsPath         string
	BinaryPath             string
	TempPath               string
	ManagementClusterHost  string
	ManagementClusterPort  string
	DNSClusterHost         string
	DNSClusterPort         string
	DockerRegistryUsername string
	DockerRegistryPassword string
	UseStaticIP            bool
	IPAddressIngress       string
	IPAddressDNS           string
}

func NewConfiguration(
	port int,
	componentsPath string,
	binaryPath string,
	tempPath string,
	clusterPublicHostname string,
	managementClusterHost string,
	managementClusterPort string,
	dnsClusterHost string,
	dnsClusterPort string,
	useStaticIP bool,
	ipAddressIngress string,
	ipAddressDNS string,
) *Config {
	return &Config{
		Port:                  port,
		ComponentsPath:        componentsPath,
		BinaryPath:            binaryPath,
		TempPath:              tempPath,
		ManagementClusterHost: managementClusterHost,
		ManagementClusterPort: managementClusterPort,
		DNSClusterHost:        dnsClusterHost,
		DNSClusterPort:        dnsClusterPort,
		UseStaticIP:           useStaticIP,
		IPAddressIngress:      ipAddressIngress,
		IPAddressDNS:          ipAddressDNS,
	}
}

func (conf *Config) CheckPath(path string) derrors.Error {
	if path == "" {
		return derrors.NewInvalidArgumentError("path cannot be empty")
	}
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return derrors.NewNotFoundError("components path must exist")
	}
	return nil
}

func (conf *Config) Validate() derrors.Error {
	conf.ComponentsPath = utils.GetPath(conf.ComponentsPath)
	conf.BinaryPath = utils.GetPath(conf.BinaryPath)
	conf.TempPath = utils.GetPath(conf.TempPath)

	if err := conf.CheckPath(conf.ComponentsPath); err != nil {
		return derrors.NewInvalidArgumentError("componentsPath").CausedBy(err)
	}
	if err := conf.CheckPath(conf.BinaryPath); err != nil {
		return derrors.NewInvalidArgumentError("binaryPath").CausedBy(err)
	}
	if err := conf.CheckPath(conf.TempPath); err != nil {
		return derrors.NewInvalidArgumentError("tempPath").CausedBy(err)
	}
	if conf.Port == 0 {
		return derrors.NewInvalidArgumentError("port must be set")
	}
	if conf.ManagementClusterHost == "" {
		return derrors.NewInvalidArgumentError("managementClusterHost must be set")
	}
	if conf.ManagementClusterPort == "" {
		return derrors.NewInvalidArgumentError("managementClusterPort must be set")
	}
	if conf.DNSClusterHost == "" {
		return derrors.NewInvalidArgumentError("dnsClusterHost must be set")
	}
	if conf.DNSClusterPort == "" {
		return derrors.NewInvalidArgumentError("dnsClusterPort must be set")
	}
	if conf.DockerRegistryUsername == "" || conf.DockerRegistryPassword == "" {
		return derrors.NewInvalidArgumentError("docker credentials must be set")
	}
	if conf.UseStaticIP && (conf.IPAddressIngress == "" && conf.IPAddressDNS == "") {
		return derrors.NewInvalidArgumentError("if useStaticIp is enabled, ip addresses must be set")
	}

	return nil
}

func (conf *Config) Print() {
	log.Info().Str("app", version.AppVersion).Str("commit", version.Commit).Msg("Version")
	log.Info().Int("port", conf.Port).Msg("gRPC Service")
	log.Info().Str("path", conf.ComponentsPath).Msg("Components")
	log.Info().Str("path", conf.BinaryPath).Msg("Binaries")
	log.Info().Str("path", conf.TempPath).Msg("Temporal files")
	log.Info().Str("host", conf.ManagementClusterHost).
		Str("port", conf.ManagementClusterPort).Msg("Management cluster")
	log.Info().Str("host", conf.DNSClusterHost).
		Str("port", conf.DNSClusterPort).Msg("DNS")
	log.Info().Str("username", conf.DockerRegistryUsername).
		Str("password", strings.Repeat("*", len(conf.DockerRegistryPassword))).Msg("Docker registry")
	log.Info().Bool("ip", conf.UseStaticIP).Msg("Use Static IP addresses")
	log.Info().Str("ip", conf.IPAddressIngress).Msg("Public ingress static IP address")
	log.Info().Str("ip", conf.IPAddressDNS).Msg("Public DNS static IP address")
}
