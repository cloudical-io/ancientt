# ancientt

A tool to automate network testing tools, like iperf3, in dynamic environments such as Kubernetes and more to come dynamic environments.

## Goals to achieve

* A bit like Prometheus blackbox exporter which contains "definitions" for probes. The "tests" would be pluggable through a Golang interface.
* "Runner" interface, e.g., for Kubernetes, Ansible, etc. The "runner" abstracts the "how it is run", e.g., for Kubernetes creates a Job, Ansible (download and) trigger a playbook to run the test.
* Store result data in different formats, e.g., CSV, excel, MySQL
  * Up for discussion: graph database ([Dgraph](https://dgraph.io/)) and / or TSDB support
* "Visualization" for humans, e.g., possibility to automatically draw "shiny" graphs from the results.

## Usage

Compile or download the Ancientt binary.

```shell
# You can also use the short flag `-c
ancientt --testdefinition your-testdefinitions.yaml
```

## Building

**Golang version**: `v1.12` or higher (tested with `v1.12.7` on `linux/amd64`)

## Licensing

Ancientt is under the Apache 2.0 License.
