# ancientt

A tool to automate network testing tools, like iperf3, in dynamic environments such as Kubernetes and more to come dynamic environments.

## Features

**TL;DR** A network test tool, like `iperf3` can be run in, e.g., Kubernetes, cluster from all-to-all Nodes.

* Run network tests with the following projects:
  * `iperf3`
  * Soon other tools will be available as well, like `smokeping`.
* Tests can be run through the following "runners":
  * Ansible (an inventory file is needed)
  * Kubernetes (a kubeconfig connected to a cluster)
* Results of the network tests can be output in different formats:
  * CSV
  * Dump (uses `pp.Sprint()` ([GitHub k0kubun/pp](https://github.com/k0kubun/pp), dump pretty print library))
  * Excel files (Excelize)
  * go-chart Charts (WIP)
  * MySQL
  * SQLite

## Usage

Either [build (`go get`)](#building) or download the Ancientt executable.

A config file containing test definitions must be given by flag `--testdefinition` (or short flag `-c`) or named `testdefinition.yaml` in the current directory.

Below command will try loading `your-testdefinitions.yaml` as the test definitions config:

```shell
# You can also use the short flag `-c
ancientt --testdefinition your-testdefinitions.yaml
```

## Demos

See [Demos](docs/demos.md).

## Goals of this Project

* A bit like Prometheus blackbox exporter which contains "definitions" for probes. The "tests" would be pluggable through a Golang interface.
* "Runner" interface, e.g., for Kubernetes, Ansible, etc. The "runner" abstracts the "how it is run", e.g., for Kubernetes creates a Job, Ansible (download and) trigger a playbook to run the test.
* Store result data in different formats, e.g., CSV, excel, MySQL
  * Up for discussion: graph database ([Dgraph](https://dgraph.io/)) and / or TSDB support
* "Visualization" for humans, e.g., possibility to automatically draw "shiny" graphs from the results.

## Development

**Golang version**: `v1.12` or higher (tested with `v1.12.9` on `linux/amd64`)

### Dependencies

`go mod` is used to manage the depeendencies.

### Building

Quickest way to just get ancientt built is to run the following command:

```bash
go get -u github.com/cloudical-io/ancientt/cmd/ancientt
```

## Licensing

Ancientt is licensed under the Apache 2.0 License.
