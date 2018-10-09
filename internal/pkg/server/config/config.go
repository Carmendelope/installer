/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package config

import (
	"github.com/nalej/derrors"
	"github.com/nalej/installer/internal/pkg/utils"
	"github.com/rs/zerolog/log"
	"os"
)

type Config struct {
	// Address where the API service will listen requests.
	Port int
	ComponentsPath string
	BinaryPath string
	TempPath string
}

func NewConfiguration(port int, componentsPath string, binaryPath string, tempPath string) * Config {
	return &Config{
		Port: port,
		ComponentsPath: componentsPath,
		BinaryPath:     binaryPath,
		TempPath:       tempPath,
	}
}

func (conf * Config) CheckPath(path string) derrors.Error {
	if path == "" {
		return derrors.NewInvalidArgumentError("path cannot be empty")
	}
	_, err := os.Stat(path);
	if os.IsNotExist(err) {
		return derrors.NewNotFoundError("components path must exist")
	}
	return nil
}

func (conf * Config) Validate() derrors.Error {
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

	return nil
}

func (conf *Config) Print() {
	log.Info().Int("port", conf.Port).Msg("gRPC Service")
	log.Info().Str("path", conf.ComponentsPath).Msg("Components")
	log.Info().Str("path", conf.BinaryPath).Msg("Binaries")
	log.Info().Str("path", conf.TempPath).Msg("Temporal files")
}