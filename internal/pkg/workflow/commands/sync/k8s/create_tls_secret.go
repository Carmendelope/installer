package k8s

import (
	"fmt"
	"encoding/json"
	"github.com/nalej/derrors"
	"github.com/nalej/installer/internal/pkg/errors"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
	"io/ioutil"
	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/rs/zerolog/log"
	"strings"
)

type CreateTLSSecret struct {
	Kubernetes
	SecretName string `json:"secret_name"`
	SecretKey    string `json:"secret_key"`
	SecretValue string `json:"secret_value"`
	LoadFromPath bool `json:"load_from_path"`
	SecretValueFromPath string `json:"secret_value_from_path"`
}

func NewCreateTLSSecret(
	kubeConfigPath string,
	secretName string,
	secretKey string,
	secretValue string,
	loadFromPath bool,
	secretValueFromPath string) *CreateTLSSecret {
	return &CreateTLSSecret{
		Kubernetes: Kubernetes{
			GenericSyncCommand: *entities.NewSyncCommand(entities.CreateTLSSecret),
			KubeConfigPath: kubeConfigPath,
		},
		SecretName: secretName,
		SecretValue: secretValue,
		LoadFromPath: loadFromPath,
		SecretValueFromPath: secretValueFromPath,
	}
}

func NewCreateTLSSecretFromJSON (raw []byte) (*entities.Command, derrors.Error) {
	f := &CreateTLSSecret{}
	if err := json.Unmarshal(raw, &f); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	f.CommandID = entities.GenerateCommandID(f.Name())
	var r entities.Command = f
	return &r, nil
}

func (cmd *CreateTLSSecret) createKubernetesSecrets() derrors.Error{

	var secretRawContent []byte
	if cmd.LoadFromPath {
		c, err := ioutil.ReadFile(cmd.SecretValueFromPath)
		if err != nil{
			return derrors.AsError(err, "cannot load secret content")
		}
		secretRawContent = c
	}else{
		secretRawContent = []byte(cmd.SecretValue)
	}

	TLSSecret := &v1.Secret {
		TypeMeta: metaV1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name:         cmd.SecretName,
			GenerateName: "",
			Namespace:    "nalej",
		},
		Data: map[string][]byte{
			cmd.SecretKey: secretRawContent,
		},
		Type: v1.SecretTypeTLS,
	}
	cmd.Connect()
	derr := cmd.Create(TLSSecret)
	if derr != nil {
		return derr
	}

	return nil
}

// Run triggers the execution of the command.
func (cmd *CreateTLSSecret) Run (workflowID string) (*entities.CommandResult, derrors.Error) {
	err := cmd.Connect()
	if err != nil{
		log.Info().Str("kubeConfigPath", cmd.KubeConfigPath).Msg("error connecting to cluster")
		return nil, derrors.NewGenericError("error connecting to cluster", err)
	}

	dErr := cmd.createKubernetesSecrets()
	if dErr != nil{
		return nil, dErr
	}

	return entities.NewSuccessCommand([]byte("Secret successfully created.")), nil
}

func (cmd *CreateTLSSecret) String () string {
	return fmt.Sprintf("SYNC CreateTLSSecret %s", cmd.SecretName)
}

func (cmd *CreateTLSSecret) PrettyPrint (indentation int) string {
	return strings.Repeat(" ", indentation) + cmd.String()
}

func (cmd *CreateTLSSecret) UserString () string {
	return fmt.Sprintf("creating tls secret")
}
