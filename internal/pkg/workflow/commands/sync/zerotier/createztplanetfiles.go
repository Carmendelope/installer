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
	"k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
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

	generateMoonOut, err := generateMoon.Output()
	if err != nil {
		log.Error().Msg("Error while executing genmoon command")
		return nil, derrors.NewGenericError("Error while executing genmoon command")
	}
	generateMoonOutStr := strings.Fields(string (generateMoonOut))
	moonName := generateMoonOutStr [1]
	fmt.Print(moonName)

	// Move moon file to planet file
	err = os.Rename("/bin/"+moonName,cmd.PlanetJsonPath)
	if err != nil {
		log.Error().Msg("Error while renaming moon file")
		return nil, derrors.NewGenericError("Error while renaming moon file")
	}

	// Planet Secret
	planetData, err := ioutil.ReadFile(cmd.PlanetJsonPath)
	if err != nil {
		log.Error().Msg("cannot read planet file")
		return nil, derrors.NewGenericError("cannot read planet file")
	}
	ztPlanetSecret := &v1.Secret{
		TypeMeta: v12.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: v12.ObjectMeta{
			Name:         "zt-planet",
			GenerateName: "",
			Namespace:    "nalej",
		},
		Data: map[string][]byte{
		 	"planet": planetData,
		 },
		Type: v1.SecretTypeDockerConfigJson,
	}
	client := cmd.Client.CoreV1().Secrets(ztPlanetSecret.Namespace)
	created, err := client.Create(ztPlanetSecret)
	if err != nil {
		log.Error().Msg("Error creating zt-planet secret")
		return nil, derrors.NewGenericError("Error creating zt-planet secret")
	}
	log.Debug().Interface("created", created).Msg("zt-planet secret has been created")

	// Identity Secret Secret
	identitySecretData, err := ioutil.ReadFile(cmd.IdentitySecretPath)
	if err != nil {
		log.Error().Msg("cannot read identity.secret file")
		return nil, derrors.NewGenericError("cannot read identity.secret file")
	}
	ztIdentitySecretSecret := &v1.Secret{
		TypeMeta: v12.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: v12.ObjectMeta{
			Name:         "zt-identity-secret",
			GenerateName: "",
			Namespace:    "nalej",
		},
		Data: map[string][]byte{
			"identity.secret": identitySecretData,
		},
		Type: v1.SecretTypeDockerConfigJson,
	}
	client = cmd.Client.CoreV1().Secrets(ztIdentitySecretSecret.Namespace)
	created, err = client.Create(ztIdentitySecretSecret)
	if err != nil {
		log.Error().Msg("Error creating zt-identity-secret secret")
		return nil, derrors.NewGenericError("Error creating zt-identity-secret secret")
	}
	log.Debug().Interface("created", created).Msg("zt-identity-secret secret has been created")

	// Identity Public Secret
	identityPublicData, err := ioutil.ReadFile(cmd.IdentityPublicPath)
	if err != nil {
		log.Error().Msg("cannot read identity.public file")
		return nil, derrors.NewGenericError("cannot read identity.public file")
	}
	ztIdentityPublicSecret := &v1.Secret{
		TypeMeta: v12.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: v12.ObjectMeta{
			Name:         "zt-identity-public",
			GenerateName: "",
			Namespace:    "nalej",
		},
		Data: map[string][]byte{
			"planet": identityPublicData,
		},
		Type: v1.SecretTypeDockerConfigJson,
	}
	client = cmd.Client.CoreV1().Secrets(ztIdentityPublicSecret.Namespace)
	created, err = client.Create(ztIdentityPublicSecret)
	if err != nil {
		log.Error().Msg("Error creating zt-planet secret")
		return nil, derrors.NewGenericError("Error creating zt-planet secret")
	}
	log.Debug().Interface("created", created).Msg("zt-planet secret has been created")

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