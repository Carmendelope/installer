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

const CheckAsset = "checkAsset"

// RKEInstall command to launch the installation of a new cluster with RKE.
const RKEInstall = "rkeInstall"

const RKERemove = "rkeRemove"

const LaunchComponents = "launchComponents"

const CheckRequirements = "checkRequirements"

const CreateClusterConfig = "createClusterConfig"

const CreateCACert = "createCACert"

const CreateManagementConfig = "createManagementConfig"

const CreateRegistrySecrets = "createRegistrySecrets"

const CreateDockerSecret = "createDockerSecret"

const UpdateCoreDNS = "updateCoreDNS"

const UpdateKubeDNS = "updateKubeDNS"

const CreateCredentials = "createCredentials"

const CreateAuthxSecret = "createAuthxSecret"

const AddClusterUser = "addClusterUser"

const InstallIngress = "installIngress"

const InstallMngtDNS = "installMngtDNS"

const InstallZtPlanetLB = "installZtPlanetLB"

const InstallVpnServerLB = "installVpnServerLB"

const CreateZTPlanetFiles = "createZtPlanetFiles"

const CreateOpaqueSecret = "createOpaqueSecret"

const InstallExtDNS = "installExtDNS"

const CreateTLSSecret = "createTLSSecret"
