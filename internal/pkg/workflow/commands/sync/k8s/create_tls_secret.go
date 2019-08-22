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
	SecretName     string `json:"secret_name"`
	PrivateKeyPath string `json:"private_key_path"`
	CertPath       string `json:"cert_path"`
}

func NewCreateTLSSecret(
	kubeConfigPath string,
	secretName string,
	privateKeyPath string,
	certPath string) *CreateTLSSecret {
	return &CreateTLSSecret{
		Kubernetes: Kubernetes{
			GenericSyncCommand: *entities.NewSyncCommand(entities.CreateTLSSecret),
			KubeConfigPath: kubeConfigPath,
		},
		SecretName:     secretName,
		PrivateKeyPath: privateKeyPath,
		CertPath:       certPath,
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

	var privateKeyRawContent []byte
	var certRawContent []byte

	pkc, err := ioutil.ReadFile(cmd.PrivateKeyPath)
	if err != nil{
		return derrors.AsError(err, "cannot load private key content")
	}
	privateKeyRawContent = pkc

	cc, err := ioutil.ReadFile(cmd.CertPath)
	if err != nil{
		return derrors.AsError(err, "cannot load cert content")
	}
	privateKeyRawContent = cc

	TLSSecret := &v1.Secret {
		TypeMeta: metaV1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		Type: v1.SecretTypeTLS,
		ObjectMeta: metaV1.ObjectMeta{
			Name:         cmd.SecretName,
			GenerateName: "",
			Namespace:    "nalej",
		},
		Data: map[string][]byte{
			"tls.key": privateKeyRawContent,
			"tls.crt": certRawContent,
		},
	}
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
