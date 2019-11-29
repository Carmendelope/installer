# installer

This repository contains the installer component in charge of installing new clusters with the Nalej components.

## Getting Started

The installer component is based on the definition of an installation template describing the steps and order in
which a set of commands must be executed. Depending on the request parameters, the template is compiled and the commands
parametrized and executed accordingly.

Notice that the component is expected to receive multiple request and perform those asynchronously. The calling
component is responsible of the overall lifecycle of the request (sending, checking, retrieving result).

### Prerequisites

Detail any component that has to be installed to run this component.

* The components to be deployed are expected to be a set of single-entity Kubernetes YAML files. Those will be
installed as part of the template, and should be available on the components path.

When deploying the component inside Kubernetes, a config file with all the YAMLs is expected to be found in order
to preload the components path. To create such configmap use the following command assuming all YAML files are
 in `./assets/appcluster` 

```
kubectl create configmap installer-configmap --from-file=assets/appcluster/ -nnalej -o yaml --dry-run > ./assets/mngtcluster/installer.configmap.yaml
```
### Build and compile

In order to build and compile this repository use the provided Makefile:

```
make all
```

This operation generates the binaries for this repo, download dependencies,
run existing tests and generate ready-to-deploy Kubernetes files.

### Run tests

Tests are executed using Ginkgo. To run all the available tests:

```
make test
```

### Update dependencies

Dependencies are managed using Godep. For an automatic dependencies download use:

```
make dep
```

In order to have all dependencies up-to-date run:

```
dep ensure -update -v
```

### Integration tests

The following table contains the variables that activate the integration tests. Integration tests are to be
considered unstable as also contain PoC for specific situations. Executing all of them can cause issues/misconfigurations
on the target clusters and may affect all existing namespaces. Execute the tests at your own risk :)

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


## User client interface

A command line interface named `installer-cli` is offered to install the management cluster.

```
$ ./bin/installer-cli install management --consoleLogging
 --binaryPath <binary_path_with_the_rke_executable>
 --componentsPath <components_path_with_the_yaml_files>
 --managementClusterPublicHost=<management_domain> --dnsClusterPublicHost=dns.<management_domain>
 --targetPlatform=AZURE
 --useStaticIPAddresses --ipAddressIngress=<ingress_ip_address> --ipAddressDNS=<dns_ip_address>
 --ipAddressCoreDNS=<app_dns_ip_address> --ipAddressVPNServer=<vpn_server_ip_address>
 --kubeConfigPath=<kubeconfig_file> --targetEnvironment=<environment_type>
```

## Known Issues

* Integration tests will be refactored so they can be properly executed without collateral damage.
* A component YAML file can only contain a single Kubernetes entity
* While partial support for minikube installations is provided, this code path has not been tested in this release.
* The install expects a set of environment variables related to docker registry secrets that are preloaded in
order to install the proper credentials in kubernetes to access private images.
* The installer should make use of the nalej-bus to send operation progress and completion
messages so that it can be easily scaled. The infrastructure-manager should listen to those
events instead of performing active coordination.

## Contributing

Please read [contributing.md](contributing.md) for details on our code of conduct, and the process for submitting pull requests to us.


## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/nalej/installer/tags). 

## Authors

See also the list of [contributors](https://github.com/nalej/installer/contributors) who participated in this project.

## License
This project is licensed under the Apache 2.0 License - see the [LICENSE-2.0.txt](LICENSE-2.0.txt) file for details.