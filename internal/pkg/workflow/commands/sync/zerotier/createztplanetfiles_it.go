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
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type CreateZTPlanetFiles struct {
	k8s.Kubernetes
	ZtIdToolBinaryPath string `json:"ztIdToolBinaryPath"`
	MgmtClusterFQDN    string `json:"management_public_host"`
	IdentitySecretPath string `json:"identitySecretPath"`
	IdentityPublicPath string `json:"identityPublicPath"`
	PlanetJsonPath string `json:"planetJsonPath"`
	PlanetPath string `json:"planetPath"`
}

func NewCreateZTPlanetFiles (
	kubeConfigPath string,
	ztIdToolBinaryPath string,
	mgmtClusterFqdn string,
	identitySecretPath string,
	identityPublicPath string,
	planetJsonPath string,
	planetPath string) *CreateZTPlanetFiles {
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
		PlanetPath: planetPath,
	}
}

type ZTPlanetJson struct {
	Id string `json:"id"`
	ObjType string `json:"objtype"`
	Roots [] struct {
		Identity string `json:"identity"`
		StableEndpoints []string `json:"stableEndpoints"`
	} `json:"roots"`
	SigningKey string `json:"signingKey"`
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

func (cmd * CreateZTPlanetFiles) generateZTIdentityFiles() derrors.Error{
	// Generate ZT Planet IDs
	generateIds := exec.Command(cmd.ZtIdToolBinaryPath, "generate", cmd.IdentitySecretPath, cmd.IdentityPublicPath)
	_, pipeErr := generateIds.StderrPipe()
	if pipeErr != nil {
		log.Error().Msg("Error while executing generate command")
		return derrors.AsError(pipeErr, errors.IOError)
	}
	err := generateIds.Run()
	if err != nil {
		log.Error().Msg("Error while executing generate command")
		return derrors.NewGenericError("Error generating ZT Planet ID files")
	}
	return nil
}

func (cmd * CreateZTPlanetFiles) initMoon() derrors.Error{
	// Init Moon
	log.Info().Msg(cmd.ZtIdToolBinaryPath+" initmoon "+"$(cat "+cmd.IdentityPublicPath+")")
	identityPublicRaw, err := ioutil.ReadFile(cmd.IdentityPublicPath)
	if err != nil{
		return derrors.NewGenericError("cannot read identity public", err)
	}
	initMoon := exec.Command(cmd.ZtIdToolBinaryPath, "initmoon", string(identityPublicRaw))
	initMoonOut, err := initMoon.StdoutPipe()
	if err != nil {
		log.Error().Msg("Error obtaining stdout for initmoon")
		return derrors.NewGenericError("Error initializing ZT Planet", err)
	}
	// redirect pipes
	initMoon.Stderr = initMoon.Stdout

	if err := initMoon.Start(); err != nil {
		log.Error().Msg("Error launching initmoon")
		return derrors.NewGenericError("Error launching initmoon", err)
	}

	planetRaw, err := ioutil.ReadAll(initMoonOut)
	if err != nil{
		return derrors.NewInternalError("cannot read planet from pipe", err)
	}

	if err := initMoon.Wait(); err != nil {
		log.Error().Msg("Error waiting for Planet JSON parse")
		return derrors.NewGenericError("Error waiting for Planet JSON parse", err)
	}

	planet := &ZTPlanetJson {}
	if err := json.Unmarshal(planetRaw, planet); err != nil {
		log.Error().Msg("Error parsing Planet JSON")
		return derrors.NewGenericError("Error parsing Planet JSON", err)
	}

	log.Info().Interface("json", planet).Msg("Planet")

	if len(planet.Roots) != 1 {
		log.Error().Msg("Unexpected number of roots found in zerotier planet file")
		return derrors.NewGenericError("Unexpected roots found in zerotier planet file")
	}

	planet.Roots[0].StableEndpoints = []string{cmd.MgmtClusterFQDN}
	planet.WorldType = "planet"

	log.Debug().Interface("zeroTierPlanet", planet).Msg("Final ZT planet")

	planetJson, err := json.Marshal(planet)
	if err != nil {
		log.Error().Msg("Error marshalling ZT Planet JSON")
		return derrors.NewGenericError("Error marshalling ZT Planet JSON", err)
	}
	err = ioutil.WriteFile(cmd.PlanetJsonPath, planetJson, 0644)
	if err != nil {
		log.Error().Msg("Error saving ZT Planet JSON file")
		return derrors.NewGenericError("Error saving ZT Planet JSON file", err)
	}
	return nil
}

func (cmd * CreateZTPlanetFiles) generatePlanet() derrors.Error{
	// Generate Planet file
	generateMoon := exec.Command(cmd.ZtIdToolBinaryPath, "genmoon", cmd.PlanetJsonPath)
	generateMoon.Dir = filepath.Dir(cmd.PlanetPath)
	_, pipeErr := generateMoon.StderrPipe()
	if pipeErr != nil {
		log.Error().Msg("Error while executing genmoon command")
		return derrors.AsError(pipeErr, errors.IOError)
	}

	generateMoonOut, err := generateMoon.Output()
	if err != nil {
		log.Error().Msg("Error while executing genmoon command")
		return derrors.NewGenericError("Error while executing genmoon command", err)
	}
	generateMoonOutStr := strings.Fields(string (generateMoonOut))
	moonName := generateMoonOutStr [1]
	log.Debug().Str("moonName", moonName).Msg("Moon")

	// Move moon file to planet file
	sourcePath := filepath.Join(filepath.Dir(cmd.PlanetPath), moonName)
	err = os.Rename(sourcePath, cmd.PlanetPath)
	if err != nil {
		log.Error().Msg("Error while renaming moon file")
		return derrors.NewGenericError("Error while renaming moon file", err)
	}
	return nil
}

func (cmd * CreateZTPlanetFiles) createKubernetesSecrets() derrors.Error{
	// Planet Secret
	planetData, err := ioutil.ReadFile(cmd.PlanetPath)
	if err != nil {
		log.Error().Msg("cannot read planet file")
		return derrors.NewGenericError("cannot read planet file", err)
	}
	ztPlanetSecret := &v1.Secret{
		TypeMeta: metaV1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name:         "zt-planet",
			GenerateName: "",
			Namespace:    "nalej",
		},
		Data: map[string][]byte{
			"planet": planetData,
		},
		Type: v1.SecretTypeOpaque,
	}
	cmd.Connect()
	client := cmd.Client.CoreV1().Secrets(ztPlanetSecret.Namespace)
	created, err := client.Create(ztPlanetSecret)
	if err != nil {
		log.Error().Msg("Error creating zt-planet secret")
		return derrors.NewGenericError("Error creating zt-planet secret", err)
	}
	log.Debug().Interface("created", created).Msg("zt-planet secret has been created")

	// Identity Secret Secret
	identitySecretData, err := ioutil.ReadFile(cmd.IdentitySecretPath)
	if err != nil {
		log.Error().Msg("cannot read identity.secret file")
		return derrors.NewGenericError("cannot read identity.secret file", err)
	}
	ztIdentitySecretSecret := &v1.Secret{
		TypeMeta: metaV1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name:         "zt-identity-secret",
			GenerateName: "",
			Namespace:    "nalej",
		},
		Data: map[string][]byte{
			"identity.secret": identitySecretData,
		},
		Type: v1.SecretTypeOpaque,
	}
	created, err = client.Create(ztIdentitySecretSecret)
	if err != nil {
		log.Error().Msg("Error creating zt-identity-secret secret")
		return derrors.NewGenericError("Error creating zt-identity-secret secret", err)
	}
	log.Debug().Interface("created", created).Msg("zt-identity-secret secret has been created")

	// Identity Public Secret
	identityPublicData, err := ioutil.ReadFile(cmd.IdentityPublicPath)
	if err != nil {
		log.Error().Msg("cannot read identity.public file")
		return derrors.NewGenericError("cannot read identity.public file", err)
	}
	ztIdentityPublicSecret := &v1.Secret{
		TypeMeta: metaV1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name:         "zt-identity-public",
			GenerateName: "",
			Namespace:    "nalej",
		},
		Data: map[string][]byte{
			"planet": identityPublicData,
		},
		Type: v1.SecretTypeOpaque,
	}
	created, err = client.Create(ztIdentityPublicSecret)
	if err != nil {
		log.Error().Msg("Error creating zt-identity-public secret")
		return derrors.NewGenericError("Error creating zt-identity-public secret", err)
	}
	log.Debug().Interface("created", created).Msg("zt-identity-public secret has been created")

	return nil
}

// Run triggers the execution of the command.
func (cmd *CreateZTPlanetFiles) Run (workflowID string) (*entities.CommandResult, derrors.Error) {
	log.Debug().Str("path", cmd.ZtIdToolBinaryPath).Msg("ZT ID Tool binary")

	dErr := cmd.generateZTIdentityFiles()
	if dErr != nil{
		return nil, dErr
	}

	dErr = cmd.initMoon()
	if dErr != nil {
		return nil, dErr
	}

	dErr = cmd.generatePlanet()
	if dErr != nil{
		return nil, dErr
	}

	dErr = cmd.createKubernetesSecrets()
	if dErr != nil{
		return nil, dErr
	}

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