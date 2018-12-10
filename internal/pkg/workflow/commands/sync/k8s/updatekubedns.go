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
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/rs/zerolog/log"
	"k8s.io/api/core/v1"
	"strings"
)

const KubeDNSNamespace = "kube-system"
var KubeDNSConfigNames = []string{"kube-dns", "kubedns-kubecfg"}
const KubeDNSSection = "stubDomains"

const KubeDNSUpdateTemplate= `{"service.nalej": [MANAGEMENT_PUBLIC_IPS]}`

type UpdateKubeDNS struct {
	Kubernetes
	ManagementPublicHost string `json:"management_public_host"`
}

func NewUpdateKubeDNS(kubeConfigPath string, managementPublicHost string) * UpdateKubeDNS{
	return &UpdateKubeDNS{
		Kubernetes:    Kubernetes{
			GenericSyncCommand: *entities.NewSyncCommand(entities.UpdateKubeDNS),
			KubeConfigPath:     kubeConfigPath,
		},
		ManagementPublicHost: managementPublicHost,
	}
}

func NewUpdateKubeDNSFromJSON(raw []byte) (*entities.Command, derrors.Error) {
	ccc := &UpdateKubeDNS{}
	if err := json.Unmarshal(raw, &ccc); err != nil {
		return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
	}
	ccc.CommandID = entities.GenerateCommandID(ccc.Name())
	var r entities.Command = ccc
	return &r, nil
}

func (uk * UpdateKubeDNS) Run(workflowID string) (*entities.CommandResult, derrors.Error) {
	connectErr := uk.Connect()
	if connectErr != nil {
		return nil, connectErr
	}
	existing, err := uk.getExistingConfig()
	if err != nil{
		return entities.NewCommandResult(false, "cannot update kube dns", err), nil
	}
	err = uk.updateConfig(existing)
	if err != nil{
		return entities.NewCommandResult(false, "cannot update kube dns", err), nil
	}
	return entities.NewSuccessCommand([]byte("Core DNS config has been updated")), nil
}

func (uk * UpdateKubeDNS) getExistingConfig() (*v1.ConfigMap, derrors.Error){
	client := uk.Client.CoreV1().ConfigMaps(KubeDNSNamespace)
	opts := metaV1.GetOptions{}

	for _, toCheck := range KubeDNSConfigNames{
		cm, err := client.Get(toCheck, opts)
		if err == nil{
			log.Debug().Str("name", cm.Name).Msg("config map found for kubedns")
			return cm, nil
		}
	}
	return nil, derrors.NewNotFoundError("cannot find kubedns config map")
}

func (uk * UpdateKubeDNS) updateConfig(cfg *v1.ConfigMap) derrors.Error {
	mgntIPs, rErr := uk.ResolveIP(uk.ManagementPublicHost)
	if rErr != nil{
		return rErr
	}
	for _, ip := range mgntIPs{
		ip = fmt.Sprintf("\"%s\"", ip)
	}

	toUpdate := strings.Replace(KubeDNSUpdateTemplate, "MANAGEMENT_PUBLIC_IPS", strings.Join(mgntIPs, ", "), 1)
	cfg.Data[KubeDNSSection] = toUpdate
	client := uk.Client.CoreV1().ConfigMaps(KubeDNSNamespace)
	updated, err := client.Update(cfg)
	if err != nil{
		return derrors.NewInternalError("cannot update config map", err)
	}
	log.Debug().Interface("updated", updated).Msg("KubeDNS configmap has been updated")
	return nil
}

func (uk * UpdateKubeDNS) String() string {
	return fmt.Sprintf("SYNC UpdateKubeDNS to %s", uk.ManagementPublicHost)
}

func (uk * UpdateKubeDNS) PrettyPrint(indentation int) string {
	return strings.Repeat(" ", indentation) + uk.String()
}

func (uk * UpdateKubeDNS) UserString() string {
	return fmt.Sprintf("Update cluster KubeDNS config to %s", uk.ManagementPublicHost)
}
