# installer

This repository contains the installer component in charge of installing new clusters with the Nalej components.

## Installing the Nalej platform

```
$ ./bin/installer-cli install management --debug --consoleLogging
```

## Creating the config map with the YAMLs for the application cluster

Assuming all YAML files are in `./assets/appcluster` use:

```
kubectl create configmap installer-configmap --from-file=assets/appcluster/ -nnalej -o yaml --dry-run > ./assets/mngtcluster/installer.configmap.yaml
```

## Installing the platform (local development environment)

Prerequirements:
1. Minikube
2. Signup component
3. Rke downloaded
4. Kubectl
5. Public api cli


Steps:

1. Compile all repositories manually.
2. Copy the yaml files to the assets directory

```
$ mkdir -p installer/assets/mngtcluster
$ mkdir -p installer/assets/appcluster
$ cp */bin/yaml/mngtcluster/* installer/assets/mngtcluster/.
$ kubectl create configmap installer-config --from-file=installer/assets/appcluster/ -nnalej -o yaml --dry-run > installer/bin/yaml/mngtcluster/installer.configmap.yaml
```

3. Load the images in the local minikube environment

```
./scripts/loadImagesInMinikube.sh
```

4. Launch the installer

```
$ cd installer
$ ./bin/installer-cli install management --debug --consoleLogging --binaryPath ~/development/rke/ --managementClusterPublicHost=192.168.99.100
```

5. Create a test organization

```
$ ../signup/bin/signup-cli signup --debug --signupAddress=192.168.99.100:32180 --orgName=nalej --ownerEmail=admin1@nalej.com --ownerName=Admin --ownerPassword=password
```

6. Setup the options (Optional step)

Notice: You may need to open nodeports for the login and public api components.

```
$ cd public-api
$ ./bin/public-api-cli options set --key=organizationID --value=<your_organization>
$ ./bin/public-api-cli options set --key=nalejAddress --value=192.168.99.100
$ ./bin/public-api-cli options set --key=port --value=31405
```

7. Login

```
$ ./bin/public-api-cli login --debug --consoleLogging --nalejAddress=192.168.99.100 --loginPort=30211 --email=admin1@nalej.com --password=password
```

8. Test

```
./bin/public-api-cli org info
```

## Installing an application cluster

# Integration tests

The following table contains the variables that activate the integration tests

| Variable  | Example Value | Description |
| ------------- | ------------- |------------- |
| RUN_INTEGRATION_TEST  | true | Run integration tests |
| IT_SSH_HOST | localhost | Host where a docker sshd image is running for SCP/SSH commands. |
| IT_SSH_PORT | 2222 | Port of the sshd server. |
| IT_RKE_PRIVATE_KEY| /private/tmp/it_test/.vagrant/machines/default/virtualbox/private_key | Private Key of the target vagrant machine |
| IT_RKE_BINARY | /Users/<yourUser>/Downloads/rke_darwin-amd64 | Path of the RKE binary |
| IT_RKE_TARGET_NODES | 172.28.128.3 | List of nodes to be installed |
| IT_K8S_KUBECONFIG | /Users/daniel/.kube/config| KubeConfig for the minikube credentials |
| IT_REGISTRY_USERNAME | <k8s_service_account_login.user_id> | Username to access the nalej repository. Use terraform output to obtain the value |
| IT_REGISTRY_PASSWORD | <k8s_service_account_login.password> | Password to access the nalej repository. Use terraform output to obtain the value |