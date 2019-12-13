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

// This command follows the instructions provided by Istio to install a multiple clusters with a shared plane control.
// See https://istio.io/docs/setup/install/multicluster/shared-gateways/ for more details
// TODO Detach this process from the istioctl binary execution

package istio

import (
    "encoding/json"
    "fmt"
    "github.com/nalej/derrors"
    "github.com/nalej/installer/internal/pkg/errors"
    "github.com/nalej/installer/internal/pkg/workflow/commands/sync"
    "github.com/nalej/installer/internal/pkg/workflow/commands/sync/k8s"
    "github.com/nalej/installer/internal/pkg/workflow/entities"
    "github.com/rs/zerolog/log"
    "io/ioutil"
    "istio.io/api/networking/v1alpha3"
    istioNetworking "istio.io/client-go/pkg/apis/networking/v1alpha3"
    istioClient "istio.io/client-go/pkg/clientset/versioned"
    "k8s.io/api/core/v1"
    metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/clientcmd"
    "strings"
    "time"
)

const (
    //IstioNamespace the name of the namespace used by Istio
    IstioNamespace = "istio-system"
    //IstioIngressGateway the name of the gateway service
    IstioIngressGateway = "istio-ingressgateway"
    //IstioSecretName name of the certificates used
    IstioSecretName = "cacerts"
    // Time between checks
    IstioTimeSleep = time.Second * 5
    // Time before timeout
    IstioTimeout = time.Second * 300
)

type InstallIstio struct {
    k8s.Kubernetes
    // Istio client to create specific Istio entities
    Istio *istioClient.Clientset
    // Path where Istio can be found
    IstioPath string        `json:"istio_path"`
    IstioCertsPath string   `json:"istio_certs_path"`
    ClusterID string        `json:"cluster_id"`
    IsAppCluster bool       `json:"is_appCluster"`
}

func NewInstallIstio(kubeConfigPath string, istioPath string, istioCertsPath string, clusterID string, isAppCluster bool) *InstallIstio {

    // use the current context in kubeconfig
    config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
    if err != nil {
        // --> Error
        return nil
    }

    istCli, err := istioClient.NewForConfig(config)
    if err != nil {
        // --> Error
        return nil
    }

    return &InstallIstio{
        Kubernetes: k8s.Kubernetes{
            GenericSyncCommand: *entities.NewSyncCommand(entities.AddClusterUser),
            KubeConfigPath:     kubeConfigPath,
        },
        IstioPath: istioPath,
        IstioCertsPath: istioCertsPath,
        Istio: istCli,
        ClusterID: clusterID,
        IsAppCluster: isAppCluster,
    }
}

// NewAddClusterUserFromJSON creates an InstallIstio command from a JSON object.
func NewInstallIstioFromJSON(raw []byte) (*entities.Command, derrors.Error) {
    lc := &InstallIstio{}
    if err := json.Unmarshal(raw, &lc); err != nil {
        return nil, derrors.NewInvalidArgumentError(errors.UnmarshalError, err)
    }

    // instantiate the Istio client
    // use the current context in kubeconfig
    config, err := clientcmd.BuildConfigFromFlags("", lc.KubeConfigPath)
    if err != nil {
        return nil, derrors.NewInternalError("impossible to get kubeconfig path", err)
    }

    istCli, err := istioClient.NewForConfig(config)
    if err != nil {
        return nil, derrors.NewInternalError("impossible to instantiate istio client")
    }

    lc.Istio = istCli

    lc.CommandID = entities.GenerateCommandID(lc.Name())
    var r entities.Command = lc
    return &r, nil
}


func (i *InstallIstio) Run(workflowID string) (*entities.CommandResult, derrors.Error) {
    // Create namespace
    connectErr := i.Connect()
    if connectErr != nil {
        return nil, connectErr
    }
    err := i.CreateNamespace(IstioNamespace)
    if err != nil {
        return nil, derrors.NewInternalError("impossible to create namespace for istio", err)
    }

    // Create secrets
    err = i.createSecrets()
    if err != nil {
        return nil, derrors.NewInternalError("impossible to create Istio secrets", err)
    }

    // Run Istioctl installer
    if i.IsAppCluster {
        // Install Istio in the application cluster
        err = i.installInSlave(i.ClusterID)
    } else {
        // Install Istio in the master
        err = i.installInMaster()
        // Create gateway
        i.installGateway()
    }

    // Wait for the gateway to have a valid ip
    // i.waitForGatewayIP()

    return entities.NewSuccessCommand([]byte("istio has been installed successfully")), nil
}

// waitForGatewayIP periodically checks the availability of the Istio gateway. The function terminates
// if and only if the gateway is available and it has its own IP address.
func (i *InstallIstio) waitForGatewayIP() derrors.Error {

    log.Info().Msg("wait for Istio ingress gateway service to be available")
    ticker := time.NewTicker(IstioTimeSleep)
    timeout := make(chan bool)
    ip := make(chan string)

    go func() {
        for {
            select {
            case <- ticker.C:
                svc, err := i.Client.CoreV1().Services(IstioNamespace).Get(IstioIngressGateway, metaV1.GetOptions{})
                if err == nil {
                    // check if we have a valid ip
                    if len(svc.Status.LoadBalancer.Ingress) > 0 {
                        svcIP := svc.Status.LoadBalancer.Ingress[0].IP
                        if len(svcIP) != 0 {
                            ip <- svcIP
                            log.Info().Msgf("Istio gateway has the associated IP: %s", svcIP)
                        }
                    }
                }
            case <- ip:
                return
            case <- timeout:
                log.Info().Msg("timeout reached when waiting for gateway service")
                return
            }
        }
    }()

    // wait until the Istio gateway service has an assigned IP
    for {
        select {
        case <- time.After(IstioTimeout):
            timeout <- true
            return derrors.NewDeadlineExceededError("timeout reached when waiting for gateway service")
        case <- ip:
            return nil
        }
    }

    return nil
}

// createSecrets builds and generates the K8s secrets to be used by Istio.
func (i *InstallIstio) createSecrets() derrors.Error {

    var caCert []byte
    var caKey []byte
    var certChain []byte
    var rootCert []byte

    caCert, err := ioutil.ReadFile(fmt.Sprintf("%s/ca-cert.pem", i.IstioCertsPath))
    if err != nil {
        return derrors.NewInternalError("error reading istio cacert",err)
    }

    caKey, err = ioutil.ReadFile(fmt.Sprintf("%s/ca-key.pem", i.IstioCertsPath))
    if err != nil {
        return derrors.NewInternalError("error reading istio ca-key",err)
    }

    certChain, err = ioutil.ReadFile(fmt.Sprintf("%s/cert-chain.pem", i.IstioCertsPath))
    if err != nil {
        return derrors.NewInternalError("error reading istio ca-key",err)
    }

    rootCert, err = ioutil.ReadFile(fmt.Sprintf("%s/root-cert.pem", i.IstioCertsPath))
    if err != nil {
        return derrors.NewInternalError("error reading istio ca-key",err)
    }

    // Generate the certificates
    secret := &v1.Secret{
        TypeMeta: metaV1.TypeMeta{
            Kind:       "Secret",
            APIVersion: "v1",
        },
        ObjectMeta: metaV1.ObjectMeta{
            Name:         IstioSecretName,
            GenerateName: "",
            Namespace:    IstioNamespace,
        },
        Data: map[string][]byte{
            "ca-cert.pem": caCert,
            "ca-key.pem": caKey,
            "cert-chain.pem": certChain,
            "root-cert.pem": rootCert,
        },
    }

    connectErr := i.Connect()
    if connectErr != nil {
        return connectErr
    }

    err = i.Create(secret)
    if err != nil {
        log.Error().Err(err).Msg("error creating istio cacerts secret")
        return derrors.NewInternalError("error creating istio cacerts secret", err)
    }

    return nil
}

func (i *InstallIstio) installInMaster() derrors.Error {
    //istioctl manifest apply --context=$CTX_CLUSTER1 \
    //  -f install/kubernetes/operator/examples/multicluster/values-istio-multicluster-primary.yaml
    log.Info().Msg("call Istioctl to install the master cluster")
    args := []string{
        "manifest",
        "apply",
        fmt.Sprintf("--kubeconfig=%s", i.KubeConfigPath),
        "-f",
        fmt.Sprintf("%s/install/kubernetes/operator/examples/multicluster/values-istio-multicluster-primary.yaml", i.IstioPath),
    }

    rExec := sync.NewExec(fmt.Sprintf("%s/bin/istioctl", i.IstioPath),args)
    _, err := rExec.Run("")
    if err != nil {
        return err
    }

    return nil
}

func (i *InstallIstio) installInSlave(clusterID string) derrors.Error {

    // We create a local kube client to check the istio ingress ip

    config, err := rest.InClusterConfig()
    if err != nil {
        return derrors.NewInternalError("impossible to get master cluster k8s configuration", err)
    }
    localClient, err := kubernetes.NewForConfig(config)
    if err != nil {
        return derrors.NewInternalError("impossible to instantiate k8s client for master cluster", err)
    }

    svc, err := localClient.CoreV1().Services(IstioNamespace).Get(IstioIngressGateway, metaV1.GetOptions{})
    if err != nil {
        log.Error().Err(err).Msg("impossible to find istio gateway service IP")
        return derrors.NewInternalError("impossible to find istio gateway service IP", err)
    }

    if len(svc.Status.LoadBalancer.Ingress) == 0 {
        return derrors.NewInternalError("no available Istio ingress IP for master cluster")
    }

    gatewayIP := svc.Status.LoadBalancer.Ingress[0].IP
    if gatewayIP == "" {
        return derrors.NewInternalError("there is no public IP for istio master gateway")
    }

     args := []string{
         "manifest",
         "apply",
         fmt.Sprintf("--kubeconfig=%s", i.KubeConfigPath),
         "--set", "values.global.mtls.enabled=true",
         "--set", "values.gateways.enabled=true",
         "--set", "values.security.selfSigned=false",
         "--set", "values.global.controlPlaneSecurityEnabled=true",
         "--set", "values.global.createRemoteSvcEndpoints=true",
         "--set", "values.global.remotePilotCreateSvcEndpoint=true",
         "--set", "values.global.remotePilotAddress="+gatewayIP,
         "--set", "values.global.remotePolicyAddress="+gatewayIP,
         "--set", "values.global.remoteTelemetryAddress="+gatewayIP,
         "--set", "values.gateways.istio-ingressgateway.env.ISTIO_META_NETWORK="+clusterID,
         "--set", "values.global.network="+clusterID,
         "--set", "autoInjection.enabled=true",
     }

    rExec := sync.NewExec(fmt.Sprintf("%s/bin/istioctl",i.IstioPath),args)
    _, execErr := rExec.Run("")
    if execErr != nil {
        return execErr
    }

    return nil
}

// installGateway to provide the master with a gateway entry point for master
func (i *InstallIstio) installGateway() derrors.Error {
    gw := istioNetworking.Gateway{
        ObjectMeta: metaV1.ObjectMeta{
            Name: "cluster-aware-gateway",
            Namespace: IstioNamespace,
        },
        Spec: v1alpha3.Gateway{
            Selector: map[string]string{
                "istio": "ingressgateway",
            },
            Servers: []*v1alpha3.Server{
                {
                    Port: &v1alpha3.Port{
                        Name: "tls",
                        Number: 443,
                        Protocol: "TLS",
                    },
                    Hosts: []string{
                        "*.local",
                    },
                    Tls: &v1alpha3.Server_TLSOptions{
                        Mode: v1alpha3.Server_TLSOptions_AUTO_PASSTHROUGH,
                    },
                },
            },
        },
    }

    _, err := i.Istio.NetworkingV1alpha3().Gateways(IstioNamespace).Create(&gw)
    if err != nil {
        return derrors.NewInternalError("error generating error", err)
    }

    return nil
}



func (i *InstallIstio) String() string {
    return fmt.Sprintf("SYNC InstallIstio")
}

func (i *InstallIstio) PrettyPrint(indentation int) string {
    return strings.Repeat(" ", indentation) + i.String()
}

func (i *InstallIstio) UserString() string {
    return fmt.Sprintf("Installing Istio")
}
