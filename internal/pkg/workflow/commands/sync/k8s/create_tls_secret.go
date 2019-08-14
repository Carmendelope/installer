package k8s

import (
	"fmt"
	"encoding/json"
	"github.com/nalej/derrors"
	"github.com/nalej/installer/internal/pkg/errors"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/rs/zerolog/log"
	"strings"
)

type CreateTLSSecret struct {
	Kubernetes
	SecretName string `json:"secret_name"`
	PrivateKeyValue string `json:"private_key_value"`
	CertValue string `json:"cert_value"`
}

func NewCreateTLSSecret(
	kubeConfigPath string,
	secretName string,
	privateKeyValue string,
	certValue string) *CreateTLSSecret {
	return &CreateTLSSecret{
		Kubernetes: Kubernetes{
			GenericSyncCommand: *entities.NewSyncCommand(entities.CreateTLSSecret),
			KubeConfigPath: kubeConfigPath,
		},
		SecretName: secretName,
		PrivateKeyValue: privateKeyValue,
		CertValue: certValue,
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
	var certRawValue []byte
	privateKeyRawContent = []byte(cmd.PrivateKeyValue)
	certRawValue = []byte(cmd.CertValue)

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
			"tls.crt": certRawValue,
		},
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
