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