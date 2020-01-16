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

// Deprecated: This set of helping entries are only intended to be used for a total Istio integration ignoring K8S
// ingresses.

package ingress
//
//import (
//    "fmt"
//    "github.com/nalej/derrors"
//    "github.com/nalej/installer/internal/pkg/workflow/entities"
//    "github.com/rs/zerolog/log"
//    "istio.io/api/networking/v1alpha3"
//    istionetworking "istio.io/client-go/pkg/apis/networking/v1alpha3"
//    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
//    "k8s.io/client-go/tools/clientcmd"
//    versionedclient "istio.io/client-go/pkg/clientset/versioned"
//)
//
//func installistioingresses() (*entities.commandresult, error) {
//    // get gateways
//    gws := getistiogateways(ii.managementpublichost)
//    // create gateways
//    gwerr := installistiogateways(gws, ii.kubeconfigpath)
//    if gwerr != nil {
//        return entities.newcommandresult(false, "impossible to create istio gateways", gwerr), gwerr
//    }
//    // get virtualservices
//    vs := getistiovirtualservices(ii.managementpublichost)
//    // create virtual services
//    vserr := installistiovirtualservices(vs, ii.kubeconfigpath)
//    if vserr != nil {
//        return entities.newcommandresult(false, "impossible to create istio virtual services", vserr), vserr
//    }
//}
//
//func getistiogateways(managementpublichost string) []*istionetworking.gateway{
//    gw := istionetworking.gateway{
//        objectmeta: metav1.objectmeta{
//            name: "web",
//            namespace: "nalej",
//        },
//        spec: v1alpha3.gateway{
//            selector: map[string]string{"istio":"ingressgateway"},
//            servers: []*v1alpha3.server{
//                {
//                    hosts: []string{
//                        fmt.sprintf("web.%s", managementpublichost),
//                    },
//                    port: &v1alpha3.port{
//                        name:     "http",
//                        protocol: "http",
//                        number:   80,
//                    },
//                },
//                {
//                    hosts: []string{
//                        fmt.sprintf("web.%s", managementpublichost),
//                    },
//                    port: &v1alpha3.port{
//                        name:     "https",
//                        protocol: "https",
//                        number:   443,
//                    },
//                    tls: &v1alpha3.server_tlsoptions{
//                        mode: v1alpha3.server_tlsoptions_passthrough,
//
//                    },
//                },
//            },
//        },
//    }
//    return []*istionetworking.gateway{&gw}
//}
//
//// install a list of istio gateways using the istio k8s versioned client.
//// params:
////  gateways
////  kubeconfigpath path to find the k8s config file
//// return:
////  error if any
//func installistiogateways(gateways []*istionetworking.gateway, kubeconfigpath string) derrors.error {
//    config, err := clientcmd.buildconfigfromflags("", kubeconfigpath)
//    if err != nil {
//        log.error().err(err).msg("error building configuration from kubeconfig")
//        return derrors.aserror(err, "error building configuration from kubeconfig")
//    }
//    // build versioned client
//    ic, err := versionedclient.newforconfig(config)
//    if err != nil {
//        log.error().err(err).msg("impossible to build a local istio client")
//        return derrors.newinternalerror("impossible to build a local istio client", err)
//    }
//
//    for _, gw := range gateways {
//        _, k8serr := ic.networkingv1alpha3().gateways("nalej").create(gw)
//        if k8serr != nil {
//            return derrors.newinternalerror("impossible to create gateway", k8serr)
//        }
//    }
//
//    return nil
//}
//
///*apiversion: extensions/v1beta1
//kind: ingress
//metadata:
//  annotations:
//    kubernetes.io/ingress.class: nginx
//    nginx.ingress.kubernetes.io/service-upstream: "true"
//  creationtimestamp: "2020-01-14t15:33:24z"
//  generation: 1
//  labels:
//    cluster: management
//    component: ingress-nginx
//  name: ingress-nginx
//  namespace: nalej
//  resourceversion: "7173"
//  selflink: /apis/extensions/v1beta1/namespaces/nalej/ingresses/ingress-nginx
//  uid: 33e24d1b-36e3-11ea-8a80-8aa429418bfb
//spec:
//  rules:
//  - host: web.master.jmmaster14.nalej.tech
//    http:
//      paths:
//      - backend:
//          servicename: web
//          serviceport: 80
//        path: /
//      - backend:
//          servicename: login-api
//          serviceport: 8443
//        path: /v1/login
//      - backend:
//          servicename: public-api
//          serviceport: 8082
//        path: /v1
//      - backend:
//          servicename: log-download-manager
//          serviceport: 8941
//        path: /v1/logs/download
//  tls:
//  - hosts:
//    - web.master.jmmaster14.nalej.tech
//    secretname: tls-client-certificate
//status:
//  loadbalancer:
//    ingress:
//    - {}
//*/
//
//func getistiovirtualservices(managementpublichost string) []*istionetworking.virtualservice{
//    vs := istionetworking.virtualservice{
//        objectmeta: metav1.objectmeta{
//            name: "web",
//            namespace: "nalej",
//        },
//        spec: v1alpha3.virtualservice {
//            gateways: []string{"web"},
//            hosts: []string{fmt.sprintf("web.master.%s", managementpublichost)},
//            http: []*v1alpha3.httproute{
//                {
//                    name: "login-api",
//                    match: []*v1alpha3.httpmatchrequest{
//                        {
//                            name: "login-api",
//                            port: 80,
//                            uri: &v1alpha3.stringmatch{matchtype: &v1alpha3.stringmatch_prefix{prefix:"/v1/login"}},
//                        },
//                    },
//                    route: []*v1alpha3.httproutedestination{
//                        {
//                            destination: &v1alpha3.destination{
//                                port: &v1alpha3.portselector{number: 8443},
//                                host: "login-api",
//                            },
//                        },
//                    },
//                },
//                {
//                    name: "public-api",
//                    match: []*v1alpha3.httpmatchrequest{
//                        {
//                            name: "public-api",
//                            port: 8082,
//                            uri: &v1alpha3.stringmatch{matchtype: &v1alpha3.stringmatch_prefix{prefix:"/v1"}},
//                        },
//                    },
//                    route: []*v1alpha3.httproutedestination{
//                        {
//                            destination: &v1alpha3.destination{
//                                port: &v1alpha3.portselector{number: 8082},
//                                host: "public-api",
//                            },
//                        },
//                    },
//                },
//                {
//                    name: "log-download-manager",
//                    match: []*v1alpha3.httpmatchrequest{
//                        {
//                            name: "log-download-manager",
//                            port: 8941,
//                            uri: &v1alpha3.stringmatch{matchtype: &v1alpha3.stringmatch_prefix{prefix:"/v1/logs/download"}},
//                        },
//                    },
//                    route: []*v1alpha3.httproutedestination{
//                        {
//                            destination: &v1alpha3.destination{
//                                port: &v1alpha3.portselector{number: 8941},
//                                host: "log-download-manager",
//                            },
//                        },
//                    },
//                },
//                {
//                    name: "web",
//                    route: []*v1alpha3.httproutedestination{
//                        {
//                            destination: &v1alpha3.destination{
//                                port: &v1alpha3.portselector{number: 80},
//                                host: "web",
//                            },
//                        },
//                    },
//                },
//            },
//        },
//    }
//    return []*istionetworking.virtualservice{&vs}
//}
//
//// install a list of istio virtual services using the istio k8s versioned client.
//// params:
////  virtualservices
////  kubeconfigpath path to find the k8s config file
//// return:
////  error if any
//func installistiovirtualservices(virtualservices []*istionetworking.virtualservice, kubeconfigpath string) derrors.error {
//    config, err := clientcmd.buildconfigfromflags("", kubeconfigpath)
//    if err != nil {
//        log.error().err(err).msg("error building configuration from kubeconfig")
//        return derrors.aserror(err, "error building configuration from kubeconfig")
//    }
//    // build versioned client
//    ic, err := versionedclient.newforconfig(config)
//    if err != nil {
//        log.error().err(err).msg("impossible to build a local istio client")
//        return derrors.newinternalerror("impossible to build a local istio client", err)
//    }
//
//    for _, vs := range  virtualservices {
//        _, k8serr := ic.networkingv1alpha3().virtualservices("nalej").create(vs)
//        if k8serr != nil {
//            return derrors.newinternalerror("impossible to create virtual service", k8serr)
//        }
//    }
//
//    return nil
//}
