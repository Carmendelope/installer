/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package entities

import (
	"github.com/nalej/derrors"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
)

const ProdRegistryName = "nalej-registry"
const StagingRegistryName = "nalej-staging-registry"
const DevRegistryName = "nalej-dev-registry"
const PublicRegistryName = "nalej-public-registry"

const EnvProdRegistryUsername = "PROD_REGISTRY_USERNAME"
const EnvProdRegistryPassword = "PROD_REGISTRY_PASSWORD"
const EnvProdRegistryURL = "PROD_REGISTRY_URL"
const EnvStagingRegistryUsername = "STAGING_REGISTRY_USERNAME"
const EnvStagingRegistryPassword = "STAGING_REGISTRY_PASSWORD"
const EnvStagingRegistryURL = "STAGING_REGISTRY_URL"
const EnvDevRegistryUsername = "DEV_REGISTRY_USERNAME"
const EnvDevRegistryPassword = "DEV_REGISTRY_PASSWORD"
const EnvDevRegistryURL = "DEV_REGISTRY_URL"
const EnvPublicRegistryUsername = "PUBLIC_REGISTRY_USERNAME"
const EnvPublicRegistryPassword = "PUBLIC_REGISTRY_PASSWORD"
const EnvPublicRegistryURL = "PUBLIC_REGISTRY_URL"

// TargetEnvironment defines the type of environment of the install. With this we can manage which type of
// images may be deployed.
type TargetEnvironment int

const (
	// Production clusters only allow production images.
	Production TargetEnvironment = iota + 1
	// Staging clusters allow Production and Staging images.
	Staging
	// Development clusters allow Production, Staging and Development images.
	Development
)

var TargetEnvironmentFromString = map[string]TargetEnvironment {
	"production":Production,
	"PRODUCTION":Production,
	"staging":Staging,
	"STAGING":Staging,
	"development":Development,
	"DEVELOPMENT":Development,
}

var TargetEnvironmentToString = map[TargetEnvironment]string{
	Production:"PRODUCTION",
	Staging:"STAGING",
	Development:"DEVELOPMENT",
}

type Environment struct{
	Target TargetEnvironment
	TargetEnvironment string `json:"target_environment"`
	ProdRegistryUsername string `json:"prod_registry_username"`
	ProdRegistryPassword string `json:"prod_registry_password"`
	ProdRegistryURL string `json:"prod_registry_url"`
	StagingRegistryUsername string `json:"staging_registry_username"`
	StagingRegistryPassword string `json:"staging_registry_password"`
	StagingRegistryURL string `json:"staging_registry_url"`
	DevRegistryUsername string `json:"dev_registry_username"`
	DevRegistryPassword string `json:"dev_registry_password"`
	DevRegistryURL string `json:"dev_registry_url"`
	PublicRegistryUsername string `json:"public_registry_username"`
	PublicRegistryPassword string `json:"public_registry_password"`
	PublicRegistryURL string `json:"public_registry_url"`
}

func NewEnvironment() *Environment{
	return &Environment{}
}

func (e *Environment) ValidateProduction() derrors.Error {
	if e.ProdRegistryUsername == "" || e.ProdRegistryPassword == "" || e.ProdRegistryURL == "" {
		return derrors.NewInvalidArgumentError("production username, password and url must be set")
	}
	return nil
}

func (e *Environment) ValidateStaging() derrors.Error {
	if e.StagingRegistryUsername == "" || e.StagingRegistryPassword == "" || e.StagingRegistryURL == "" {
		return derrors.NewInvalidArgumentError("staging username, password and url must be set")
	}
	return nil
}

func (e *Environment) ValidateDevelopment() derrors.Error {
	if e.DevRegistryUsername == "" || e.DevRegistryPassword == "" || e.DevRegistryURL == "" {
		return derrors.NewInvalidArgumentError("development username, password and url must be set")
	}
	return nil
}

func (e *Environment) ValidatePublic() derrors.Error {
	if e.PublicRegistryUsername == "" || e.PublicRegistryPassword == "" || e.PublicRegistryURL == "" {
		return derrors.NewInvalidArgumentError("public username, password and url must be set")
	}
	return nil
}

func (e*Environment) envOrElse(envName string, paramValue string) string{
	if paramValue != "" {
		return paramValue
	}
	fromEnv := os.Getenv(envName)
	if fromEnv != "" {
		return fromEnv
	}
	return ""
}

// Resolve applies the environment variables
func (e*Environment) Resolve() {
	// Production
	e.ProdRegistryUsername = e.envOrElse(EnvProdRegistryUsername, e.ProdRegistryUsername)
	e.ProdRegistryPassword = e.envOrElse(EnvProdRegistryPassword, e.ProdRegistryPassword)
	e.ProdRegistryURL = e.envOrElse(EnvProdRegistryURL, e.ProdRegistryURL)
	// Staging
	e.StagingRegistryUsername = e.envOrElse(EnvStagingRegistryUsername, e.StagingRegistryUsername)
	e.StagingRegistryPassword = e.envOrElse(EnvStagingRegistryPassword, e.StagingRegistryPassword)
	e.StagingRegistryURL = e.envOrElse(EnvStagingRegistryURL, e.StagingRegistryURL)
	// Development
	e.DevRegistryUsername = e.envOrElse(EnvDevRegistryUsername, e.DevRegistryUsername)
	e.DevRegistryPassword = e.envOrElse(EnvDevRegistryPassword, e.DevRegistryPassword)
	e.DevRegistryURL = e.envOrElse(EnvDevRegistryURL, e.DevRegistryURL)
	// Public
	e.PublicRegistryUsername = e.envOrElse(EnvPublicRegistryUsername, e.PublicRegistryUsername)
	e.PublicRegistryPassword = e.envOrElse(EnvPublicRegistryPassword, e.PublicRegistryPassword)
	e.PublicRegistryURL = e.envOrElse(EnvPublicRegistryURL, e.PublicRegistryURL)
}

// Validate the environment
func (e *Environment) Validate() derrors.Error {
	env, found := TargetEnvironmentFromString[strings.ToLower(e.TargetEnvironment)]
	if !found {
		return derrors.NewNotFoundError("invalid target environment").WithParams(e.TargetEnvironment)
	}
	e.Target = env

	vErr := e.ValidateProduction()
	if vErr != nil{
		return vErr
	}
	vErr = e.ValidatePublic()
	if vErr != nil{
		return vErr
	}
	if e.Target == Staging || e.Target == Development {
		vErr = e.ValidateStaging()
		if vErr != nil{
			return vErr
		}
	}
	if e.Target == Development {
		vErr = e.ValidateDevelopment()
		if vErr != nil{
			return vErr
		}
	}
	return nil
}

func (e*Environment) Print() {
	log.Info().
		Str("username", e.ProdRegistryUsername).
		Str("password", strings.Repeat("*", len(e.ProdRegistryPassword))).
		Str("URL", e.ProdRegistryURL).Msg("Production registry")
	if e.Target == Staging || e.Target == Development{
		log.Info().
			Str("username", e.StagingRegistryUsername).
			Str("password", strings.Repeat("*", len(e.StagingRegistryPassword))).
			Str("URL", e.StagingRegistryURL).Msg("Staging registry")
	}
	if e.Target == Development{
		log.Info().
			Str("username", e.DevRegistryUsername).
			Str("password", strings.Repeat("*", len(e.DevRegistryPassword))).
			Str("URL", e.DevRegistryURL).Msg("Development registry")
	}
	log.Info().
		Str("username", e.PublicRegistryUsername).
		Str("password", strings.Repeat("*", len(e.PublicRegistryPassword))).
		Str("URL", e.PublicRegistryURL).Msg("Public registry")
}