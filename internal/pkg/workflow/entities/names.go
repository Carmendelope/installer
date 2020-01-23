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

// Enumeration of the supported commands.

package entities

// Exec command to execute local commands in the system.
const Exec = "exec"

// SCP command to copy data to remote hosts.
const SCP = "scp"

// SSH command to execute commands on remote hosts.
const SSH = "ssh"

// Logger command to add logging information to the workflow log.
const Logger = "logger"

// Fail command that aborts the workflow execution.
const Fail = "fail"

// Sleep commands that waits for a given ammount of time.
const Sleep = "sleep"

// GroupCmd command to sequentialize subsets of commands.
const GroupCmd = "group"

// ParallelCmd command to execute several commands in parallel.
const ParallelCmd = "parallel"

// TryCmd command that tries to execute a command, and in case of failure executes an alternative one.
const TryCmd = "try"

// ProcessCheck command to determine if a process is running on a given machine.
const ProcessCheck = "processCheck"

// CheckAsset command to determine if a given asset file exists.
const CheckAsset = "checkAsset"

// RKEInstall command to launch the installation of a new cluster with RKE.
const RKEInstall = "rkeInstall"

// RKERemove command to remove a kubernetes installed with RKE
const RKERemove = "rkeRemove"

// LaunchComponents command to install a set of YAML Kubernetes files
const LaunchComponents = "launchComponents"

// CheckRequirements checks the requirements of the installer against the installed Kubernetes.
const CheckRequirements = "checkRequirements"

// CreateClusterConfig command to create the configmap of the cluster.
const CreateClusterConfig = "createClusterConfig"

// CreateCACert command to create the Nalej CA certificate.
const CreateCACert = "createCACert"

// CreateManagementConfig command to create the configmap with the configuration of the system in the management cluster.
const CreateManagementConfig = "createManagementConfig"

// CreateRegistrySecrets command to create a set of secrets to download images from private registries.
const CreateRegistrySecrets = "createRegistrySecrets"

// UpdateCoreDNS command to update the configuration of CoreDNS.
const UpdateCoreDNS = "updateCoreDNS"

// UpdateKubeDNS command to update the configuration of KubeDNS.
const UpdateKubeDNS = "updateKubeDNS"

// AddClusterUser command to create a user for an application cluster.
const AddClusterUser = "addClusterUser"

// InstallIngress command to create a set of Ingresses for the platform services.
const InstallIngress = "installIngress"

// InstallMngtDNS command to install the Consul DNS load balancer
const InstallMngtDNS = "installMngtDNS"

// InstallExtDNS command to install the external DNS
const InstallExtDNS = "installExtDNS"

// InstallZtPlanetLB command to install a loadbalancer for the Zerotier planet.
const InstallZtPlanetLB = "installZtPlanetLB"

// InstallVpnServerLB command to install a loadbalancer for the edge controller VPN.
const InstallVpnServerLB = "installVpnServerLB"

// CreateZTPlanetFiles command to create the required planet files in a cluster.
const CreateZTPlanetFiles = "createZtPlanetFiles"

// CreateOpaqueSecret command to create an opaque Kubernetes secret.
const CreateOpaqueSecret = "createOpaqueSecret"

// CreateTLSSecret command to create a Kubernetes TLS secret.
const CreateTLSSecret = "createTLSSecret"

// CreateDockerSecret creates a docker secret in Kubernetes.
const CreateDockerSecret = "createDockerSecret"

// DeleteNamespace command to delete a namespace in Kubernetes.
const DeleteNamespace = "deleteNamespace"

// DeleteServiceAccount command to delete a Kubernetes service account entity.
const DeleteServiceAccount = "deleteServiceAccount"

// DeleteNalejNamespace command to delete the contents of the Nalej namespace except the entities created by the provisioner.
const DeleteNalejNamespace = "deleteNalejNamespace"

// DeleteClusterRoleBindind command to delete a Kubernetes cluster role binding.
const DeleteClusterRoleBinding = "deleteClusterRoleBinding"

// DeleteClusterRole command to delete a Kubernetes cluster role.
const DeleteClusterRole = "deleteClusterRole"

// DeleteRole command to delete a Kubernetest role from.
const DeleteRole = "deleteRole"

// DeleteRoleBinding command to delete a Kubernetes role binding.
const DeleteRoleBinding = "deleteRoleBinding"

// DeleteConfigMap command to delete a Kubernetes config map.
const DeleteConfigMap = "deleteConfigMap"

// DeleteService command to delete a Kubernetes service.
const DeleteService = "deleteService"

// DeleteDeployment command to delete a Kubernetes deployment.
const DeleteDeployment = "deleteDeployment"

// DeletePodSecurityPolicy command to delete a Kubernetes pod security policy.
const DeletePodSecurityPolicy = "deletePodSecurityPolicy"

// InstallIstio command to run the istio installation process.
const InstallIstio = "installIstio"
