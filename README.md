# ancientt

A tool to automate network testing tools, like iperf3, in dynamic environments such as Kubernetes and more to come dynamic environments.

Container Image available from:

* [GHCR.io](https://github.com/users/cloudical-io/packages/container/package/ancientt)

Container Image Tags:

* `main` - Latest build of the `main` branch.
* `vx.y.z` - Tagged build of the application.

## Features

**TL;DR** A network test tool, like `iperf3` can be run in, e.g., Kubernetes, cluster from all-to-all Nodes.

* Run network tests with the following projects:
  * [`iperf3`](https://iperf.fr/)
  * [PingParsing](https://github.com/thombashi/pingparsing)
  * Soon more tools will be available as well, see [GitHub Issues with "testers" Label](https://github.com/cloudical-io/ancientt/issues?utf8=%E2%9C%93&q=is%3Aissue+is%3Aopen+label%3Atesters+).
* Tests can be run through the following "runners":
  * Ansible (an inventory file is needed)
  * Kubernetes (a kubeconfig connected to a cluster)
* Results of the network tests can be output in different formats:
  * CSV
  * Dump (uses `pp.Sprint()` ([GitHub k0kubun/pp](https://github.com/k0kubun/pp), pretty print library))
  * Excel files (Excelize)
  * go-chart Charts (WIP)
  * MySQL
  * SQLite

## Usage

Either [build (`go get`)](#building), download the Ancientt executable from the GitHub release page or use the Container image.

A config file containing test definitions must be given by flag `--testdefinition` (or short flag `-c`) or named `testdefinition.yaml` in the current directory.

Below command will try loading `your-testdefinitions.yaml` as the test definitions config:

```shell
$ ancientt --testdefinition your-testdefinitions.yaml
# You can also use the short flag `-c` instead of `--testdefinition`
# and also with `-y` run the tests immediately
$ ancientt -c your-testdefinitions.yaml -y
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

**Golang version**: `v1.17` or higher (tested with `v1.17.6` on `linux/amd64`)

### Creating Release

1. Add new entry for release to [`CHANGELOG.md`](CHANGELOG.md).
2. Update [`VERSION`](VERSION) with new version number.
3. `git commit` and `git push` both changes (e.g., `version: update to VERSION_HERE`).
4. Now create the git tag and push the tag `git tag VERSION_HERE` followed by `git push --tags`.

### Dependencies

`go mod` is used to manage the dependencies.

### Building

Quickest way to just get ancientt built is to run the following command:

```bash
go get -u github.com/cloudical-io/ancientt/cmd/ancientt
```

## Licensing

Ancientt is licensed under the Apache 2.0 License.
