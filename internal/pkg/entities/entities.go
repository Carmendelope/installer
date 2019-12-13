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

package entities

import (
	"github.com/nalej/derrors"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
)

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


var TargetEnvironmentFromString = map[string]TargetEnvironment{
	"production":  Production,
	"PRODUCTION":  Production,
	"staging":     Staging,
	"STAGING":     Staging,
	"development": Development,
	"DEVELOPMENT": Development,
}

var TargetEnvironmentToString = map[TargetEnvironment]string{
	Production:  "PRODUCTION",
	Staging:     "STAGING",
	Development: "DEVELOPMENT",
}

// NetworkingMode indicates the kind of networking solution to be used in the platform
type NetworkingMode string

const (
	NetworkingModeZt = "zt"
	NetworkingModeIstio = "istio"
	// It indicates a non valid mode
	NetworkingModeInvalid = ""
)

var NetworkingModeFromString = map[string] NetworkingMode {
	"zt": NetworkingModeZt,
	"istio": NetworkingModeIstio,
}

var NetworkingModeToString = map[NetworkingMode] string {
	NetworkingModeZt: "zt",
	NetworkingModeIstio: "istio",
}


type Environment struct {
	Target            TargetEnvironment
	TargetEnvironment string `json:"target_environment"`
}

func NewEnvironment() *Environment {
	return &Environment{}
}

func (e *Environment) envOrElse(envName string, paramValue string) string {
	if paramValue != "" {
		return paramValue
	}
	fromEnv := os.Getenv(envName)
	if fromEnv != "" {
		return fromEnv
	}
	return ""
}

// Validate the environment
func (e *Environment) Validate() derrors.Error {
	env, found := TargetEnvironmentFromString[strings.ToLower(e.TargetEnvironment)]
	if !found {
		return derrors.NewNotFoundError("invalid target environment").WithParams(e.TargetEnvironment)
	}
	e.Target = env
	return nil
}

func (e *Environment) Print() {
	log.Info().Str("Environment", TargetEnvironmentToString[e.Target])
}
