/*
 * Copyright 2020 Nalej
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

package ingress

import (
    "fmt"
    "github.com/nalej/derrors"
    "github.com/rs/zerolog/log"
    "istio.io/api/networking/v1alpha3"
    istioNetworking "istio.io/client-go/pkg/apis/networking/v1alpha3"
    metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/tools/clientcmd"
    versionedclient "istio.io/client-go/pkg/clientset/versioned"
)


func GetIstioGateways(managementPublicHost string) []*istioNetworking.Gateway{
    gw := istioNetworking.Gateway{
        ObjectMeta: metaV1.ObjectMeta{
            Name: "web",
            Namespace: "nalej",
        },
        Spec: v1alpha3.Gateway{
            Selector: map[string]string{"istio":"ingressgateway"},
            Servers: []*v1alpha3.Server{
                {
                    Hosts: []string{
                        fmt.Sprintf("web.%s", managementPublicHost),
                    },
                    Port: &v1alpha3.Port{
                        Name:     "http",
                        Protocol: "HTTP",
                        Number:   80,
                    },
                },
                {
                    Hosts: []string{
                        fmt.Sprintf("web.%s", managementPublicHost),
                    },
                    Port: &v1alpha3.Port{
                        Name:     "https",
                        Protocol: "HTTPS",
                        Number:   443,
                    },
                    Tls: &v1alpha3.Server_TLSOptions{
                        Mode: v1alpha3.Server_TLSOptions_PASSTHROUGH,

                    },
                },
            },
        },
    }
    return []*istioNetworking.Gateway{&gw}
}

// Install a list of Istio gateways using the Istio k8s versioned client.
// params:
//  gateways
//  kubeconfigPath path to find the k8s config file
// return:
//  error if any
func InstallIstioGateways(gateways []*istioNetworking.Gateway, kubeConfigPath string) derrors.Error {
    config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
    if err != nil {
        log.Error().Err(err).Msg("error building configuration from kubeconfig")
        return derrors.AsError(err, "error building configuration from kubeconfig")
    }
    // build versioned client
    ic, err := versionedclient.NewForConfig(config)
    if err != nil {
        log.Error().Err(err).Msg("impossible to build a local Istio client")
        return derrors.NewInternalError("impossible to build a local Istio client", err)
    }

    for _, gw := range gateways {
        _, k8sErr := ic.NetworkingV1alpha3().Gateways("nalej").Create(gw)
        if k8sErr != nil {
            return derrors.NewInternalError("impossible to create gateway", k8sErr)
        }
    }

    return nil
}

/*apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  annotations:
    kubernetes.io/ingress.class: nginx
    nginx.ingress.kubernetes.io/service-upstream: "true"
  creationTimestamp: "2020-01-14T15:33:24Z"
  generation: 1
  labels:
    cluster: management
    component: ingress-nginx
  name: ingress-nginx
  namespace: nalej
  resourceVersion: "7173"
  selfLink: /apis/extensions/v1beta1/namespaces/nalej/ingresses/ingress-nginx
  uid: 33e24d1b-36e3-11ea-8a80-8aa429418bfb
spec:
  rules:
  - host: web.master.jmmaster14.nalej.tech
    http:
      paths:
      - backend:
          serviceName: web
          servicePort: 80
        path: /
      - backend:
          serviceName: login-api
          servicePort: 8443
        path: /v1/login
      - backend:
          serviceName: public-api
          servicePort: 8082
        path: /v1
      - backend:
          serviceName: log-download-manager
          servicePort: 8941
        path: /v1/logs/download
  tls:
  - hosts:
    - web.master.jmmaster14.nalej.tech
    secretName: tls-client-certificate
status:
  loadBalancer:
    ingress:
    - {}
*/

func GetIstioVirtualServices(managementPublicHost string) []*istioNetworking.VirtualService{
    vs := istioNetworking.VirtualService{
        ObjectMeta: metaV1.ObjectMeta{
            Name: "web",
            Namespace: "nalej",
        },
        Spec: v1alpha3.VirtualService {
            Gateways: []string{"web"},
            Hosts: []string{fmt.Sprintf("web.master.%s", managementPublicHost)},
            Http: []*v1alpha3.HTTPRoute{
                {
                    Name: "login-api",
                    Match: []*v1alpha3.HTTPMatchRequest{
                        {
                            Name: "login-api",
                            Port: 80,
                            Uri: &v1alpha3.StringMatch{MatchType: &v1alpha3.StringMatch_Prefix{Prefix:"/v1/login"}},
                        },
                    },
                    Route: []*v1alpha3.HTTPRouteDestination{
                        {
                            Destination: &v1alpha3.Destination{
                                Port: &v1alpha3.PortSelector{Number: 8443},
                                Host: "login-api",
                            },
                        },
                    },
                },
                {
                    Name: "public-api",
                    Match: []*v1alpha3.HTTPMatchRequest{
                        {
                            Name: "public-api",
                            Port: 8082,
                            Uri: &v1alpha3.StringMatch{MatchType: &v1alpha3.StringMatch_Prefix{Prefix:"/v1"}},
                        },
                    },
                    Route: []*v1alpha3.HTTPRouteDestination{
                        {
                            Destination: &v1alpha3.Destination{
                                Port: &v1alpha3.PortSelector{Number: 8082},
                                Host: "public-api",
                            },
                        },
                    },
                },
                {
                    Name: "log-download-manager",
                    Match: []*v1alpha3.HTTPMatchRequest{
                        {
                            Name: "log-download-manager",
                            Port: 8941,
                            Uri: &v1alpha3.StringMatch{MatchType: &v1alpha3.StringMatch_Prefix{Prefix:"/v1/logs/download"}},
                        },
                    },
                    Route: []*v1alpha3.HTTPRouteDestination{
                        {
                            Destination: &v1alpha3.Destination{
                                Port: &v1alpha3.PortSelector{Number: 8941},
                                Host: "log-download-manager",
                            },
                        },
                    },
                },
                {
                    Name: "web",
                    Route: []*v1alpha3.HTTPRouteDestination{
                        {
                            Destination: &v1alpha3.Destination{
                                Port: &v1alpha3.PortSelector{Number: 80},
                                Host: "web",
                            },
                        },
                    },
                },
            },
        },
    }
    return []*istioNetworking.VirtualService{&vs}
}

// Install a list of Istio virtual services using the Istio k8s versioned client.
// params:
//  virtualServices
//  kubeconfigPath path to find the k8s config file
// return:
//  error if any
func InstallIstioVirtualServices(virtualServices []*istioNetworking.VirtualService, kubeConfigPath string) derrors.Error {
    config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
    if err != nil {
        log.Error().Err(err).Msg("error building configuration from kubeconfig")
        return derrors.AsError(err, "error building configuration from kubeconfig")
    }
    // build versioned client
    ic, err := versionedclient.NewForConfig(config)
    if err != nil {
        log.Error().Err(err).Msg("impossible to build a local Istio client")
        return derrors.NewInternalError("impossible to build a local Istio client", err)
    }

    for _, vs := range  virtualServices {
        _, k8sErr := ic.NetworkingV1alpha3().VirtualServices("nalej").Create(vs)
        if k8sErr != nil {
            return derrors.NewInternalError("impossible to create virtual service", k8sErr)
        }
    }

    return nil
}