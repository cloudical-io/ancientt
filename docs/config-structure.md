# Config Structure

This Document documents the types introduced by Ancientt for configuration to be used by users.

> Note this document is generated from code comments. When contributing a change to this document please do so by changing the code comments.

## Table of Contents

* [AdditionalFlags](#additionalflags)
* [AnsibleGroups](#ansiblegroups)
* [CSV](#csv)
* [Config](#config)
* [Dump](#dump)
* [Excelize](#excelize)
* [GoChart](#gochart)
* [Hosts](#hosts)
* [IPerf3](#iperf3)
* [KubernetesHosts](#kuberneteshosts)
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
| Clients | \n List of additional flags for clients | []string | true |
| Server | \n List of additional flags for server | []string | true |

[Back to TOC](#table-of-contents)

## AnsibleGroups

AnsibleGroups server and clients host group names in the used inventory file(s)

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| Server | Server inventory server group name | string | true |
| Clients | Clients inventory clients group name | string | true |

[Back to TOC](#table-of-contents)

## CSV

CSV CSV Output config options

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| FilePath | File base path for output | string | true |
| NamePattern | File name pattern templated from various availables during output generation | string | true |

[Back to TOC](#table-of-contents)

## Config

Config Config object for the config file

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| Version |  | string | true |
| Runner |  | [Runner](#runner) | true |
| Tests |  | []*[Test](#test) | true |

[Back to TOC](#table-of-contents)

## Dump

Dump Dump Output config options

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| FilePath | File base path for output | string | true |
| NamePattern | File name pattern templated from various availables during output generation | string | true |

[Back to TOC](#table-of-contents)

## Excelize

Excelize Excelize Output config options. TODO implement

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| FilePath | File base path for output | string | true |
| NamePattern | File name pattern templated from various availables during output generation | string | true |
| SaveAfterRows | After what amount of rows the Excel file should be saved | int | true |

[Back to TOC](#table-of-contents)

## GoChart

GoChart GoChart Output config options

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| FilePath | File base path for output | string | true |
| NamePattern | File name pattern templated from various availables during output generation | string | true |
| Types | Types of charts to produce from the testers output data | []string | true |

[Back to TOC](#table-of-contents)

## Hosts

Hosts options for hosts selection for a Test

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| Name | Name of this hosts selection. | string | true |
| All | If all hosts available should be used. | bool | true |
| Random | Select `Count` Random hosts from the available hosts list. | bool | true |
| Count | Must be used with `Random`, will cause `Count` times Nodes to be randomly selected from all applicable hosts. | int | true |
| Hosts | Static list of hosts (this list is not checked for accuracy) | []string | true |
| HostSelector | \"Label\" selector for the dynamically generated hosts list, e.g., Kubernetes label selector | map[string]string | true |
| AntiAffinity | AntiAffinity not implemented yet | []string | true |

[Back to TOC](#table-of-contents)

## IPerf3

IPerf3 IPerf3 config structure for testers.Tester config

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| AdditionalFlags | Additional flags for client and server | [AdditionalFlags](#additionalflags) | true |
| UDP | If UDP should be used for the IPerf3 test | *bool | true |

[Back to TOC](#table-of-contents)

## KubernetesHosts

KubernetesHosts hosts selection options for Kubernetes

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| IgnoreSchedulingDisabled | If Nodes that are `SchedulingDisabled` should be ignored | bool | true |
| Tolerations | List of Kubernetes corev1.Toleration to tolerate when selecting Nodes | [][corev1.Toleration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.15/#toleration-v1-core) | true |

[Back to TOC](#table-of-contents)

## KubernetesTimeouts

KubernetesTimeouts timeouts for operations with Kubernetess

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| DeleteTimeout | Timeout for object deletion | int | true |
| RunningTimeout | Timeout for \"Pod running\" check | int | true |
| SucceedTimeout | Timeout for \"Pod succeded\" check (e.g., client Pod exits after Pod) | int | true |

[Back to TOC](#table-of-contents)

## MySQL

MySQL MySQL Output config options

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| DSN | MySQL DSN, format `[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]`, for more information see https://github.com/go-sql-driver/mysql#dsn-data-source-name | string | true |
| TableNamePattern | Pattern used for templating the name of the table used in the MySQL database, the tables are created automatically when MySQL.AutoCreateTables is set to `true` | string | true |
| AutoCreateTables | Automatically create tables in the MySQL database (default `true`) | *bool | true |

[Back to TOC](#table-of-contents)

## Output

Output Output config structure pointing to the other config options for each output

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| Name | Name of this output | string | true |
| CSV | CSV output options | *[CSV](#csv) | true |
| GoChart | GoChart output options | *[GoChart](#gochart) | true |
| Dump | Dump output options | *[Dump](#dump) | true |
| Excelize | Excelize output options | *[Excelize](#excelize) | true |
| SQLite | SQLite output options | *[SQLite](#sqlite) | true |
| MySQL | MySQL output options | *[MySQL](#mysql) | true |

[Back to TOC](#table-of-contents)

## RunOptions

RunOptions options for running the tasks

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| ContinueOnError |  | bool | true |
| Rounds | Amount of test rounds (repetitions) to do for a test plan | int | true |
| Interval | Time interval to sleep / wait between | time.Duration | true |
| Mode | Run mode can be `parallel` or `sequential` (default is `sequential`) | string | true |
| ParallelCount | **NOT IMPLEMENTED YET** amount of test tasks to run when using `parallel` RunOptions.Mode | int | true |

[Back to TOC](#table-of-contents)

## Runner

Runner structure with all available runners config options

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| Name | Name of the runner | string | true |
| Kubernetes | Kubernetes runner options | *[RunnerKubernetes](#runnerkubernetes) | true |
| Ansible | Ansible runner options | *[RunnerAnsible](#runneransible) | true |
| Mock | Mock runner options (userd for testing purposes) | *[RunnerMock](#runnermock) | true |

[Back to TOC](#table-of-contents)

## RunnerAnsible

RunnerAnsible Ansible Runner config options

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| InventoryFilePath | InventoryFilePath Path to inventory file to use | string | true |
| Groups | Groups server and clients group names | *[AnsibleGroups](#ansiblegroups) | true |
| AnsibleCommand | Path to the ansible command (if empty will be searched for in `PATH`) | string | true |
| AnsibleInventoryCommand | Path to the ansible-inventory command (if empty will be searched for in `PATH`) | string | true |
| CommandTimeout | Timeout duration for `ansible` and `ansible-inventory` calls (NOT task command timeouts) | time.Duration | true |
| TaskCommandTimeout | Timeout duration for `ansible` Task command calls | time.Duration | true |

[Back to TOC](#table-of-contents)

## RunnerKubernetes

RunnerKubernetes Kubernetes Runner config options

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| InClusterConfig | If the Kubernetes client should use the in-cluster config for the cluster communication | bool | true |
| Kubeconfig | Path to your kubeconfig file, if not set the `KUBECONFIG` env var will be used and then the default | string | true |
| Image | The image used for the spawned Pods for the tests (default: `quay.io/galexrt/container-toolbox`) | string | true |
| Namespace | Namespace to execute the tests in | string | true |
| HostNetwork | If `hostNetwork` mode should be used for the test Pods | bool | true |
| Timeouts | Timeout settings for operations against the Kubernetes API | *[KubernetesTimeouts](#kubernetestimeouts) | true |
| Annotations | Annotations to put on the test Pods | map[string]string | true |
| Hosts | Host selection specific options | *[KubernetesHosts](#kuberneteshosts) | true |

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
| FilePath | File base path for output | string | true |
| NamePattern | File name pattern templated from various availables during output generation | string | true |
| TableNamePattern | Pattern used for templating the name of the table used in the SQLite database, the tables are created automatically | string | true |

[Back to TOC](#table-of-contents)

## Test

Test Config options for each Test

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| Name | Test name | string | true |
| Type | The tester to use, e.g., for `iperf3` set to `iperf3` and so on | string | true |
| RunOptions | Options for the execution of the test | [RunOptions](#runoptions) | true |
| Outputs | List of Outputs to use for processing data from the testers. | [][Output](#output) | true |
| Hosts | Hosts selection for client and server | [TestHosts](#testhosts) | true |
| IPerf3 | IPerf3 test options | *[IPerf3](#iperf3) | true |

[Back to TOC](#table-of-contents)

## TestHosts

TestHosts list of clients and servers hosts for use in the test(s)

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| Clients | Static list of hosts to use as clients | [][Hosts](#hosts) | true |
| Servers | Static list of hosts to use as server | [][Hosts](#hosts) | true |

[Back to TOC](#table-of-contents)
