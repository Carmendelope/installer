/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package zerotier

import (
	"encoding/json"
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/installer/internal/pkg/errors"
	"github.com/nalej/installer/internal/pkg/workflow/commands/sync/k8s"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
	"github.com/rs/zerolog/log"
	"io/ioutil"
	"os/exec"
	"strings"
)

type CreateZTPlanetFiles struct {
	k8s.Kubernetes
	ZtIdToolBinaryPath string `json:"ztIdToolBinaryPath"`
	MgmtClusterFQDN    string `json:"mgmtClusterFqdn"`
	IdentitySecretPath string `json:"identitySecretPath"`
	IdentityPublicPath string `json:"identityPublicPath"`
	PlanetJsonPath string `json:"planetJsonPath"`
}

type ZTPlanetJson struct {
	Id string `json:"id"`
	ObjType string `json:"objtype"`
	Roots [] struct{
		Identity string `json:"identity"`
		StableEndpoints []string `json:"stableEndpoints"`
	} `json:"roots"`
	SigningKey string `json:"signinKey"`
	SigningKeySecret string `json:"signingKey_SECRET"`
	UpdatesMustBeSignedBy string `json:"updatesMustBeSignedBy"`
	WorldType string `json:"worldType"`
}

func NewCreateZTPlanetFiles (
	kubeConfigPath string,
	ztIdToolBinaryPath string,
	mgmtClusterFqdn string,
	identitySecretPath string,
	identityPublicPath string,
	planetJsonPath string) *CreateZTPlanetFiles {
	return &CreateZTPlanetFiles{
		Kubernetes: k8s.Kubernetes{
			GenericSyncCommand: *entities.NewSyncCommand(entities.CreateZTPlanetFiles),
			KubeConfigPath: kubeConfigPath,
		},
		ZtIdToolBinaryPath: ztIdToolBinaryPath,
		MgmtClusterFQDN: mgmtClusterFqdn,
		IdentitySecretPath: identitySecretPath,
		IdentityPublicPath: identityPublicPath,
		PlanetJsonPath: planetJsonPath,
	}
}

func NewCreateZTPlanetFilesFromJSON (raw []byte) (*entities.Command, derrors.Error) {
	f := &CreateZTPlanetFiles{}
	if err := json.Unmarshal(raw, &f); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	f.CommandID = entities.GenerateCommandID(f.Name())
	var r entities.Command = f
	return &r, nil
}

// Run triggers the execution of the command.
func (cmd *CreateZTPlanetFiles) Run (workflowID string) (*entities.CommandResult, derrors.Error) {
	log.Debug().Str("path", cmd.ZtIdToolBinaryPath).Msg("ZT ID Tool binary")

	// Generate ZT Planet IDs
	generateIds := exec.Command(cmd.ZtIdToolBinaryPath, "generate", cmd.IdentitySecretPath, cmd.IdentityPublicPath)
	_, pipeErr := generateIds.StderrPipe()
	if pipeErr != nil {
		log.Error().Msg("Error while executing generate command")
		return nil, derrors.AsError(pipeErr, errors.IOError)
	}
	err := generateIds.Run()
	if err != nil {
		log.Error().Msg("Error while executing generate command")
		derrors.NewGenericError("Error generating ZT Planet ID files")
	}

	// Init Moon
	initMoon := exec.Command(cmd.ZtIdToolBinaryPath, "initmoon", "$("+cmd.IdentityPublicPath+")")
	_, pipeErr = initMoon.StderrPipe()
	if pipeErr != nil {
		log.Error().Msg("Error while executing initmoon command")
		return nil, derrors.AsError(pipeErr, errors.IOError)
	}
	initMoonOut, err := initMoon.StdoutPipe()
	if err != nil {
		log.Error().Msg("Error initializing ZT Planet")
		return nil, derrors.NewGenericError("Error initializing ZT Planet")
	}
	if err := initMoon.Start(); err != nil {
		log.Error().Msg("Error starting ZT Planet initialization command")
		return nil, derrors.NewGenericError("Error starting ZT Planet initialization command")
	}

	var planet ZTPlanetJson
	if err := json.NewDecoder(initMoonOut).Decode(planet); err != nil {
		log.Error().Msg("Error parsing Planet JSON")
		return nil, derrors.NewGenericError("Error parsing Planet JSON")
	}

	if err := initMoon.Wait(); err != nil {
		log.Error().Msg("Error waiting for Planet JSON parse")
		return nil, derrors.NewGenericError("Error waiting for Planet JSON parse")
	}

	log.Debug().Interface("zeroTierPlanet", planet).Msg("Empty ZT planet")

	if len(planet.Roots) != 1 {
		log.Error().Msg("Unexpected roots found in zerotier planet file")
		return nil, derrors.NewGenericError("Unexpected roots found in zerotier planet file")
	}

	planet.Roots[0].StableEndpoints = []string{cmd.MgmtClusterFQDN}
	planet.WorldType = "planet"

	log.Debug().Interface("zeroTierPlanet", planet).Msg("Final ZT planet")

	planetJson, err := json.Marshal(planet)
	if err != nil {
		log.Error().Msg("Error marshalling ZT Planet JSON")
		return nil, derrors.NewGenericError("Error marshalling ZT Planet JSON")
	}
	err = ioutil.WriteFile(cmd.PlanetJsonPath, planetJson, 0644)
	if err != nil {
		log.Error().Msg("Error saving ZT Planet JSON file")
		return nil, derrors.NewGenericError("Error saving ZT Planet JSON file")
	}


	// Generate Planet file
	generateMoon := exec.Command(cmd.ZtIdToolBinaryPath, "genmoon", cmd.PlanetJsonPath)
	_, pipeErr = generateMoon.StderrPipe()
	if pipeErr != nil {
		log.Error().Msg("Error while executing genmoon command")
		return nil, derrors.AsError(pipeErr, errors.IOError)
	}
	generateMoonOut, err := generateMoon.StdoutPipe()
	if err != nil {
		log.Error().Msg("Error generating ZT Planet file")
		return nil, derrors.NewGenericError("Error generating ZT Planet file")
	}
	if err := generateMoon.Start(); err != nil {
		log.Error().Msg("Error starting ZT Planet file creation command")
		return nil, derrors.NewGenericError("Error starting ZT Planet file creation command")
	}

	var planetFile string

	return entities.NewSuccessCommand([]byte("ZT Planet files and secrets successfully created.")), nil
}

func (cmd *CreateZTPlanetFiles) String () string {
	return fmt.Sprintf("SYNC CreateZTPlanetFiles on %s", cmd.KubeConfigPath)
}

func (cmd *CreateZTPlanetFiles) PrettyPrint (indentation int) string {
	return strings.Repeat(" ", indentation) + cmd.String()
}

func (cmd *CreateZTPlanetFiles) UserString () string {
	return fmt.Sprintf("Creating ZT Planet files")
}