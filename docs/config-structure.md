# Config Structure

This Document documents the types introduced by Ancientt for configuration to be used by users.

> Note this document is generated from code comments. When contributing a change to this document please do so by changing the code comments.

## Table of Contents

* [AdditionalFlags](#additionalflags)
* [AnsibleGroups](#ansiblegroups)
* [AnsibleTimeouts](#ansibletimeouts)
* [CSV](#csv)
* [Config](#config)
* [Dump](#dump)
* [Excelize](#excelize)
* [FilePath](#filepath)
* [GoChart](#gochart)
* [Hosts](#hosts)
* [IPerf3](#iperf3)
* [KubernetesHosts](#kuberneteshosts)
* [KubernetesServiceAccounts](#kubernetesserviceaccounts)
* [KubernetesTimeouts](#kubernetestimeouts)
* [MySQL](#mysql)
* [Output](#output)
* [RunOptions](#runoptions)
* [Runner](#runner)
* [RunnerAnsible](#runneransible)
* [RunnerKubernetes](#runnerkubernetes)
* [RunnerMock](#runnermock)
* [SQLite](#sqlite)
* [Test](#test)
* [TestHosts](#testhosts)

## AdditionalFlags

AdditionalFlags additional flags structure for Server and Clients

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| clients | \n List of additional flags for clients | []string | true |
| server | \n List of additional flags for server | []string | true |

[Back to TOC](#table-of-contents)

## AnsibleGroups

AnsibleGroups server and clients host group names in the used inventory file(s)

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| server | Server inventory server group name | string | true |
| clients | Clients inventory clients group name | string | true |

[Back to TOC](#table-of-contents)

## AnsibleTimeouts

AnsibleTimeouts timeouts for Ansible command runs

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| commandTimeout | Timeout duration for `ansible` and `ansible-inventory` calls (NOT task command timeouts) | time.Duration | true |
| taskCommandTimeout | Timeout duration for `ansible` Task command calls | time.Duration | true |

[Back to TOC](#table-of-contents)

## CSV

CSV CSV Output config options

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| FilePath |  | [FilePath](#filepath) | false |

[Back to TOC](#table-of-contents)

## Config

Config Config object for the config file

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| version |  | string | true |
| runner |  | [Runner](#runner) | true |
| tests |  | []*[Test](#test) | true |

[Back to TOC](#table-of-contents)

## Dump

Dump Dump Output config options

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| FilePath |  | [FilePath](#filepath) | false |

[Back to TOC](#table-of-contents)

## Excelize

Excelize Excelize Output config options. TODO implement

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| FilePath |  | [FilePath](#filepath) | false |
| saveAfterRows | After what amount of rows the Excel file should be saved | int | true |

[Back to TOC](#table-of-contents)

## FilePath

FilePath file path and name pattern for outputs file generation

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| filePath | File base path for output | string | true |
| namePattern | File name pattern templated from various availables during output generation | string | true |

[Back to TOC](#table-of-contents)

## GoChart

GoChart GoChart Output config options

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| FilePath |  | [FilePath](#filepath) | false |
| types | Types of charts to produce from the testers output data | []string | true |

[Back to TOC](#table-of-contents)

## Hosts

Hosts options for hosts selection for a Test

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| name | Name of this hosts selection. | string | true |
| all | If all hosts available should be used. | bool | true |
| random | Select `Count` Random hosts from the available hosts list. | bool | true |
| count | Must be used with `Random`, will cause `Count` times Nodes to be randomly selected from all applicable hosts. | int | true |
| hosts | Static list of hosts (this list is not checked for accuracy) | []string | true |
| hostSelector | \"Label\" selector for the dynamically generated hosts list, e.g., Kubernetes label selector | map[string]string | true |
| antiAffinity | AntiAffinity not implemented yet | []string | true |

[Back to TOC](#table-of-contents)

## IPerf3

IPerf3 IPerf3 config structure for testers.Tester config

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| additionalFlags | Additional flags for client and server | [AdditionalFlags](#additionalflags) | true |
| udp | If UDP should be used for the IPerf3 test | *bool | true |

[Back to TOC](#table-of-contents)

## KubernetesHosts

KubernetesHosts hosts selection options for Kubernetes

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| ignoreSchedulingDisabled | If Nodes that are `SchedulingDisabled` should be ignored | *bool | true |
| tolerations | List of Kubernetes corev1.Toleration to tolerate when selecting Nodes | [][corev1.Toleration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.15/#toleration-v1-core) | true |

[Back to TOC](#table-of-contents)

## KubernetesServiceAccounts

KubernetesServiceAccounts server and client ServiceAccount name to use for the created Pods

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| server | Server ServiceAccount name to use for server Pods | string | false |
| clients | Clients ServiceAccount name to use for client Pods | string | false |

[Back to TOC](#table-of-contents)

## KubernetesTimeouts

KubernetesTimeouts timeouts for operations with Kubernetess

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| deleteTimeout | Timeout for object deletion | int | true |
| runningTimeout | Timeout for \"Pod running\" check | int | true |
| succeedTimeout | Timeout for \"Pod succeded\" check (e.g., client Pod exits after Pod) | int | true |

[Back to TOC](#table-of-contents)

## MySQL

MySQL MySQL Output config options

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| dsn | MySQL DSN, format `[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]`, for more information see https://github.com/go-sql-driver/mysql#dsn-data-source-name | string | true |
| tableNamePattern | Pattern used for templating the name of the table used in the MySQL database, the tables are created automatically when MySQL.AutoCreateTables is set to `true` | string | true |
| autoCreateTables | Automatically create tables in the MySQL database (default `true`) | *bool | true |

[Back to TOC](#table-of-contents)

## Output

Output Output config structure pointing to the other config options for each output

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| name | Name of this output | string | true |
| csv | CSV output options | *[CSV](#csv) | true |
| goChart | GoChart output options | *[GoChart](#gochart) | true |
| dump | Dump output options | *[Dump](#dump) | true |
| excelize | Excelize output options | *[Excelize](#excelize) | true |
| sqlite | SQLite output options | *[SQLite](#sqlite) | true |
| mysql | MySQL output options | *[MySQL](#mysql) | true |

[Back to TOC](#table-of-contents)

## RunOptions

RunOptions options for running the tasks

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| continueOnError | Continue on error during test runs (recommended to set to `true`) (default is `true`) | *bool | false |
| rounds | Amount of test rounds (repetitions) to do for a test plan | int | true |
| interval | Time interval to sleep / wait between (default `10s`) | time.Duration | false |
| mode | Run mode can be `parallel` or `sequential` (see `RunMode`, default is `sequential`) | RunMode | true |
| parallelCount | **NOT IMPLEMENTED YET** amount of test tasks to run when using `parallel` RunOptions.Mode | int | true |

[Back to TOC](#table-of-contents)

## Runner

Runner structure with all available runners config options

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| name | Name of the runner | string | true |
| kubernetes | Kubernetes runner options | *[RunnerKubernetes](#runnerkubernetes) | true |
| ansible | Ansible runner options | *[RunnerAnsible](#runneransible) | true |
| mock | Mock runner options (userd for testing purposes) | *[RunnerMock](#runnermock) | true |

[Back to TOC](#table-of-contents)

## RunnerAnsible

RunnerAnsible Ansible Runner config options

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| inventoryFilePath | InventoryFilePath Path to inventory file to use | string | true |
| groups | Groups server and clients group names | *[AnsibleGroups](#ansiblegroups) | true |
| ansibleCommand | Path to the ansible command (if empty will be searched for in `PATH`) | string | true |
| ansibleInventoryCommand | Path to the ansible-inventory command (if empty will be searched for in `PATH`) | string | true |
| timeouts | Timeout settings for ansible command runs | *[AnsibleTimeouts](#ansibletimeouts) | true |

[Back to TOC](#table-of-contents)

## RunnerKubernetes

RunnerKubernetes Kubernetes Runner config options

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| inClusterConfig | If the Kubernetes client should use the in-cluster config for the cluster communication | bool | true |
| kubeconfig | Path to your kubeconfig file, if not set the following order will be tried out, `KUBECONFIG` and `$HOME/.kube/config` | string | false |
| image | The image used for the spawned Pods for the tests (default `quay.io/galexrt/container-toolbox`) | string | true |
| namespace | Namespace to execute the tests in | string | true |
| hostNetwork | If `hostNetwork` mode should be used for the test Pods | *bool | false |
| timeouts | Timeout settings for operations against the Kubernetes API | *[KubernetesTimeouts](#kubernetestimeouts) | false |
| annotations | Annotations to put on the test Pods | map[string]string | false |
| hosts | Host selection specific options | *[KubernetesHosts](#kuberneteshosts) | true |
| serviceaccounts | ServiceAccounst to use server and client Pods | *[KubernetesServiceAccounts](#kubernetesserviceaccounts) | false |

[Back to TOC](#table-of-contents)

## RunnerMock

RunnerMock Mock Runner config options (here for good measure)

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |

[Back to TOC](#table-of-contents)

## SQLite

SQLite SQLite Output config options

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| FilePath |  | [FilePath](#filepath) | false |
| tableNamePattern | Pattern used for templating the name of the table used in the SQLite database, the tables are created automatically | string | true |

[Back to TOC](#table-of-contents)

## Test

Test Config options for each Test

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| name | Test name | string | true |
| type | The tester to use, e.g., for `iperf3` set to `iperf3` and so on | string | true |
| runOptions | Options for the execution of the test | [RunOptions](#runoptions) | true |
| outputs | List of Outputs to use for processing data from the testers. | [][Output](#output) | true |
| hosts | Hosts selection for client and server | [TestHosts](#testhosts) | true |
| iperf3 | IPerf3 test options | *[IPerf3](#iperf3) | true |

[Back to TOC](#table-of-contents)

## TestHosts

TestHosts list of clients and servers hosts for use in the test(s)

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| clients | Static list of hosts to use as clients | [][Hosts](#hosts) | true |
| servers | Static list of hosts to use as server | [][Hosts](#hosts) | true |

[Back to TOC](#table-of-contents)
