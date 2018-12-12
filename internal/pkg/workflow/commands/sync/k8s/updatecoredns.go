/*
 * Copyright (C) 2018 Nalej - All Rights Reserved
 */

package k8s

import (
	"encoding/json"
	"fmt"
	"github.com/nalej/derrors"
	"github.com/nalej/installer/internal/pkg/errors"
	"github.com/nalej/installer/internal/pkg/workflow/entities"
	"github.com/rs/zerolog/log"
	"k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

const CoreDNSNamespace = "kube-system"
const CoreDNSConfigName = "coredns"
const CoreDNSSection = "Corefile"
const CoreDNSUpdateTemplate = `.:53 {
    errors
    health
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        upstream
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    proxy . /etc/resolv.conf
    cache 30
    reload
}

service.nalej {
    log stdout
    proxy . DNS_PUBLIC_IPS
    cache 30
}
`

type UpdateCoreDNS struct {
	Kubernetes
	DNSPublicHost string `json:"dns_public_host"`
	DNSPublicPort string `json:"dns_public_port"`
}

func NewUpdateCoreDNS(kubeConfigPath string, dnsPublicHost string) * UpdateCoreDNS {
	return &UpdateCoreDNS{
		Kubernetes:    Kubernetes{
			GenericSyncCommand: *entities.NewSyncCommand(entities.UpdateCoreDNS),
			KubeConfigPath:     kubeConfigPath,
		},
		DNSPublicHost: dnsPublicHost,
	}
}

func NewUpdateCoreDNSFromJSON(raw []byte) (*entities.Command, derrors.Error) {
	ccc := &UpdateCoreDNS{}
	if err := json.Unmarshal(raw, &ccc); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	ccc.CommandID = entities.GenerateCommandID(ccc.Name())
	var r entities.Command = ccc
	return &r, nil
}

func (uc * UpdateCoreDNS) Run(workflowID string) (*entities.CommandResult, derrors.Error) {
	connectErr := uc.Connect()
	if connectErr != nil {
		return nil, connectErr
	}
	existing, err := uc.getExistingConfig()
	if err != nil{
		return entities.NewCommandResult(false, "cannot update core dns", err), nil
	}
	err = uc.updateConfig(existing)
	if err != nil{
		return entities.NewCommandResult(false, "cannot update core dns", err), nil
	}
	return entities.NewSuccessCommand([]byte("Core DNS config has been updated")), nil
}

func (uc * UpdateCoreDNS) getExistingConfig() (*v1.ConfigMap, derrors.Error){
	client := uc.Client.CoreV1().ConfigMaps(CoreDNSNamespace)
	opts := metaV1.GetOptions{}
	cm, err := client.Get(CoreDNSConfigName, opts)
	if err != nil{
		return nil, derrors.NewNotFoundError("cannot obtain coredns config map", err)
	}
	return cm, nil
}

func (uc * UpdateCoreDNS) updateConfig(cfg *v1.ConfigMap) derrors.Error {
	log.Debug().Interface("data", cfg.Data[CoreDNSSection]).Msg("current data")
	mgntIPs, rErr := uc.ResolveIP(uc.DNSPublicHost)
	if rErr != nil{
		return rErr
	}
	for _, ip := range mgntIPs{
		ip = fmt.Sprintf("%s:%s", ip, uc.DNSPublicPort)
	}

	toUpdate := strings.Replace(CoreDNSUpdateTemplate, "DNS_PUBLIC_IPS", strings.Join(mgntIPs, " "), 1)
	cfg.Data[CoreDNSSection] = toUpdate
	client := uc.Client.CoreV1().ConfigMaps(CoreDNSNamespace)
	updated, err := client.Update(cfg)
	if err != nil{
		return derrors.NewInternalError("cannot update config map", err)
	}
	log.Debug().Interface("updated", updated).Msg("CoreDNS configmap has been updated")
	return nil
}

func (uc * UpdateCoreDNS) String() string {
	return fmt.Sprintf("SYNC UpdateCoreDNS to %s", uc.DNSPublicHost)
}

func (uc * UpdateCoreDNS) PrettyPrint(indentation int) string {
	return strings.Repeat(" ", indentation) + uc.String()
}

func (uc * UpdateCoreDNS) UserString() string {
	return fmt.Sprintf("Update cluster CoreDNS config to %s", uc.DNSPublicHost)
}