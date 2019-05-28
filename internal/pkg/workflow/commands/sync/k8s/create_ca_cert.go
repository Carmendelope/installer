/*
 * Copyright (C) 2019 Nalej - All Rights Reserved
 */

package k8s

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/json"
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/installer/internal/pkg/errors"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
	"math/big"
	"strings"
	"time"
)

// CertValidity of 2 years
const CertValidity = time.Hour * 24 * 365 * 2

type CreateCACert struct{
	Kubernetes
	PublicHost     string `json:"public_host"`
	certificate []byte
}

func NewCreateCACert(
	kubeConfigPath string,
	publicHost string) * CreateCACert{
	return &CreateCACert{
		Kubernetes: Kubernetes{
			GenericSyncCommand: *entities.NewSyncCommand(entities.CreateCACert),
			KubeConfigPath:     kubeConfigPath,
		},
		PublicHost: publicHost,
	}
}

func NewCreateCACertFromJSON(raw []byte) (*entities.Command, derrors.Error){
	cmc := &CreateManagementConfig{}
	if err := json.Unmarshal(raw, &cmc); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	cmc.CommandID = entities.GenerateCommandID(cmc.Name())
	var r entities.Command = cmc
	return &r, nil
}

func (cc * CreateCACert) createCACertificate() derrors.Error{

	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	if err != nil{
		return derrors.AsError(err, "cannot create private key for CA cert")
	}

	caCert := x509.Certificate{

		Signature:                   nil,
		SignatureAlgorithm:          0,
		PublicKeyAlgorithm:          0,
		PublicKey:                   nil,
		SerialNumber:                big.NewInt(1),
		Issuer:                      pkix.Name{
			Organization:       []string{"Nalej"},
		},
		Subject:                     pkix.Name{
			Organization: []string{"Nalej"},
		},
		NotBefore:                   time.Now(),
		NotAfter:                    time.Now().Add(CertValidity),
		KeyUsage:                    x509.KeyUsageCertSign | x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:                 []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid:       true,
		IsCA:                  true,
		MaxPathLen:            0,
		MaxPathLenZero:        true,
		DNSNames:              []string{fmt.Sprintf("*.%s", cc.PublicHost)},
	}
	result, err := x509.CreateCertificate(rand.Reader, &caCert, &caCert, privateKey.PublicKey, privateKey)
	if err != nil{
		return derrors.AsError(err, "cannot create CA certificate")
	}
	cc.certificate = result
	return nil
}

func (cc *CreateCACert) Run(workflowID string) (*entities.CommandResult, derrors.Error) {
	connectErr := cc.Connect()
	if connectErr != nil {
		return nil, connectErr
	}

	cErr := cc.CreateNamespacesIfNotExist(TargetNamespace)
	if cErr != nil {
		return entities.NewCommandResult(false, "cannot create namespace", cErr), nil
	}
	return entities.NewSuccessCommand([]byte("CA cert created an installed on cluster")), nil
}

func (cc *CreateCACert) String() string {
	return fmt.Sprintf("SYNC CreateCACert")
}

func (cc *CreateCACert) PrettyPrint(indentation int) string {
	simpleIden := strings.Repeat(" ", indentation) +  "  "
	entrySep := simpleIden +  "  "
	msg := fmt.Sprintf("\n%sCert:\n%sPublicHost: %s",
		simpleIden,
		entrySep, cc.PublicHost,
	)
	return strings.Repeat(" ", indentation) + cc.String() + msg
}

func (cc *CreateCACert) UserString() string {
	return fmt.Sprintf("Creating CA certificate for %s", cc.PublicHost)
}