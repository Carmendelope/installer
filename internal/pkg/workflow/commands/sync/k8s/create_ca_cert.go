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

package k8s

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/installer/internal/pkg/errors"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
	"github.com/rs/zerolog/log"
	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"math/big"
	"strings"
	"time"
)

// CertValidity of 2 years
const CertValidity = time.Hour * 24 * 365 * 2

type CreateCACert struct {
	Kubernetes
	PublicHost     string `json:"public_host"`
	certificate    []byte
	certificatePEM string
	privateKeyPEM  string
}

func NewCreateCACert(
	kubeConfigPath string,
	publicHost string) *CreateCACert {
	return &CreateCACert{
		Kubernetes: Kubernetes{
			GenericSyncCommand: *entities.NewSyncCommand(entities.CreateCACert),
			KubeConfigPath:     kubeConfigPath,
		},
		PublicHost: publicHost,
	}
}

func NewCreateCACertFromJSON(raw []byte) (*entities.Command, derrors.Error) {
	cmc := &CreateCACert{}
	if err := json.Unmarshal(raw, &cmc); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	cmc.CommandID = entities.GenerateCommandID(cmc.Name())
	var r entities.Command = cmc
	return &r, nil
}

func (cc *CreateCACert) createCACertificate() derrors.Error {

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return derrors.AsError(err, "cannot create private key for CA cert")
	}

	caCert := x509.Certificate{

		SerialNumber: big.NewInt(1),
		Issuer: pkix.Name{
			Organization: []string{"Nalej"},
		},
		Subject: pkix.Name{
			Organization: []string{"Nalej"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(CertValidity),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            0,
		MaxPathLenZero:        true,
		DNSNames:              []string{fmt.Sprintf("*.%s", cc.PublicHost)},
	}
	publicKey := &privateKey.PublicKey
	result, err := x509.CreateCertificate(rand.Reader, &caCert, &caCert, publicKey, privateKey)
	if err != nil {
		return derrors.AsError(err, "cannot create CA certificate")
	}
	cc.certificate = result

	// Export the content to PEM
	CAOut := &bytes.Buffer{}
	err = pem.Encode(CAOut, &pem.Block{Type: "CERTIFICATE", Bytes: cc.certificate})
	if err != nil {
		return derrors.AsError(err, "cannot transform certificate to PEM")
	}
	cc.certificatePEM = CAOut.String()

	PKOut := &bytes.Buffer{}
	err = pem.Encode(PKOut, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err != nil {
		return derrors.AsError(err, "cannot transform private key to PEM")
	}
	cc.privateKeyPEM = PKOut.String()
	return nil
}

func (cc *CreateCACert) createCertSecret() derrors.Error {
	tlsSecret := &v1.Secret{
		TypeMeta: metaV1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name:         "mngt-ca-cert",
			GenerateName: "",
			Namespace:    "nalej",
		},
		Data: nil,
		StringData: map[string]string{
			"tls.crt": cc.certificatePEM,
			"tls.key": cc.privateKeyPEM,
		},
		Type: v1.SecretTypeTLS,
	}
	cc.Connect()
	derr := cc.Create(tlsSecret)
	if derr != nil {
		return derr
	}
	return nil
}

func (cc *CreateCACert) Run(workflowID string) (*entities.CommandResult, derrors.Error) {
	connectErr := cc.Connect()
	if connectErr != nil {
		return nil, connectErr
	}

	cErr := cc.CreateNamespaceIfNotExists(TargetNamespace)
	if cErr != nil {
		return entities.NewCommandResult(false, "cannot create namespace", cErr), nil
	}

	// Create certificate
	err := cc.createCACertificate()
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot create CA certificate")
		return entities.NewCommandResult(false, "cannot create CA certificate", err), nil
	}

	// Create secret in kubernetes
	err = cc.createCertSecret()
	if err != nil {
		log.Error().Str("trace", err.DebugReport()).Msg("cannot create CA certificate secret")
		return entities.NewCommandResult(false, "cannot create CA certificate secret", err), nil
	}

	return entities.NewSuccessCommand([]byte("CA cert created an installed on cluster")), nil
}

func (cc *CreateCACert) String() string {
	return fmt.Sprintf("SYNC CreateCACert")
}

func (cc *CreateCACert) PrettyPrint(indentation int) string {
	simpleIden := strings.Repeat(" ", indentation) + "  "
	entrySep := simpleIden + "  "
	msg := fmt.Sprintf("\n%sCert:\n%sPublicHost: %s",
		simpleIden,
		entrySep, cc.PublicHost,
	)
	return strings.Repeat(" ", indentation) + cc.String() + msg
}

func (cc *CreateCACert) UserString() string {
	return fmt.Sprintf("Creating CA certificate for %s", cc.PublicHost)
}
