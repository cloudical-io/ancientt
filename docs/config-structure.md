# Config Structure

This Document documents the types introduced by ACNTT for configuration to be used by users.

> Note this document is generated from code comments. When contributing a change to this document please do so by changing the code comments.

## Table of Contents
* [AdditionalFlags](#additionalflags)
* [CSV](#csv)
* [Config](#config)
* [Dump](#dump)
* [Excelize](#excelize)
* [GoChart](#gochart)
* [Hosts](#hosts)
* [IPerf3](#iperf3)
* [KeyValuePair](#keyvaluepair)
* [KubernetesHosts](#kuberneteshosts)
* [KubernetesTimeouts](#kubernetestimeouts)
* [MySQL](#mysql)
* [Output](#output)
* [RunOptions](#runoptions)
* [Runner](#runner)
* [RunnerKubernetes](#runnerkubernetes)
* [RunnerMock](#runnermock)
* [SQLite](#sqlite)
* [Siege](#siege)
* [Smokeping](#smokeping)
* [Test](#test)
* [TestHosts](#testhosts)

## AdditionalFlags

AdditionalFlags additional flags structure for Server and Clients

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| Clients |  | []string | true |
| Server |  | []string | true |

[Back to TOC](#table-of-contents)

## CSV

CSV CSV Output config options

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| FilePath |  | string | true |
| NamePattern |  | string | true |

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
| FilePath |  | string | true |
| NamePattern |  | string | true |

[Back to TOC](#table-of-contents)

## Excelize

Excelize Excelize Output config options. TODO implement

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| FilePath |  | string | true |
| NamePattern |  | string | true |

[Back to TOC](#table-of-contents)

## GoChart

GoChart GoChart Output config options

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| FilePath |  | string | true |
| NamePattern |  | string | true |
| Types |  | []string | true |

[Back to TOC](#table-of-contents)

## Hosts

Hosts options for hosts selection for a Test

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| Name | Name of this hosts selection. | string | true |
| All | If all hosts available should be used. | bool | true |
| Random | Select `Count` Random hosts from the available hosts list. | bool | true |
| Count | Used with Random to randomly select the Count of hosts. | int | true |
| Hosts |  | []string | true |
| HostSelector |  | map[string]string | true |
| AntiAffinity | AntiAffinity not implemented yet | [][KeyValuePair](#keyvaluepair) | true |

[Back to TOC](#table-of-contents)

## IPerf3

IPerf3 IPerf3 config structure for testers.Tester config

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| AdditionalFlags |  | [AdditionalFlags](#additionalflags) | true |
| UDP |  | *bool | true |

[Back to TOC](#table-of-contents)

## KeyValuePair

KeyValuePair key value string pair

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| Key |  | string | true |
| Value |  | string | true |

[Back to TOC](#table-of-contents)

## KubernetesHosts

KubernetesHosts hosts selection options for Kubernetes

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| IgnoreSchedulingDisabled |  | bool | true |
| Tolerations |  | [][corev1.Toleration](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.15/#toleration-v1-core) | true |

[Back to TOC](#table-of-contents)

## KubernetesTimeouts

KubernetesTimeouts timeouts for operations with Kubernetess

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| DeleteTimeout |  | int | true |
| RunningTimeout |  | int | true |
| SucceedTimeout |  | int | true |

[Back to TOC](#table-of-contents)

## MySQL

MySQL MySQL Output config options

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| DSN |  | string | true |
| TableNamePattern |  | string | true |

[Back to TOC](#table-of-contents)

## Output

Output Output config structure pointing to the other config options for each output

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| Name |  | string | true |
| CSV |  | *[CSV](#csv) | true |
| GoChart |  | *[GoChart](#gochart) | true |
| Dump |  | *[Dump](#dump) | true |
| Excelize |  | *[Excelize](#excelize) | true |
| SQLite |  | *[SQLite](#sqlite) | true |
| MySQL |  | *[MySQL](#mysql) | true |

[Back to TOC](#table-of-contents)

## RunOptions

RunOptions options for running the tasks

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| ContinueOnError |  | bool | true |
| Rounds |  | int | true |
| Interval |  | time.Duration | true |
| Mode |  | string | true |
| ParallelCount |  | int | true |

[Back to TOC](#table-of-contents)

## Runner

Runner structure with all available runners config options

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| Name |  | string | true |
| Kubernetes |  | *[RunnerKubernetes](#runnerkubernetes) | true |
| Mock |  | *[RunnerMock](#runnermock) | true |

[Back to TOC](#table-of-contents)

## RunnerKubernetes

RunnerKubernetes Kubernetes Runner config options

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| Kubeconfig |  | string | true |
| Image |  | string | true |
| Namespace |  | string | true |
| HostNetwork |  | bool | true |
| Timeouts |  | *[KubernetesTimeouts](#kubernetestimeouts) | true |
| Annotations |  | map[string]string | true |
| Hosts |  | *[KubernetesHosts](#kuberneteshosts) | true |

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
| FilePath |  | string | true |
| NamePattern |  | string | true |
| TableNamePattern |  | string | true |

[Back to TOC](#table-of-contents)

## Siege

Siege Siege config structure TODO not implemented yet

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| AdditionalFlags |  | [AdditionalFlags](#additionalflags) | true |
| Benchmark |  | bool | true |
| Headers |  | map[string]string | true |
| URLs |  | []string | true |
| UserAgent |  | string | true |

[Back to TOC](#table-of-contents)

## Smokeping

Smokeping Smokeping config structure TODO not implemented yet

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| AdditionalFlags |  | [AdditionalFlags](#additionalflags) | true |

[Back to TOC](#table-of-contents)

## Test

Test Config options for each Test

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| Name |  | string | true |
| Type |  | string | true |
| RunOptions |  | [RunOptions](#runoptions) | true |
| Outputs |  | [][Output](#output) | true |
| Hosts |  | [TestHosts](#testhosts) | true |
| IPerf3 |  | *[IPerf3](#iperf3) | true |
| Siege |  | *[Siege](#siege) | true |
| Smokeping |  | *[Smokeping](#smokeping) | true |

[Back to TOC](#table-of-contents)

## TestHosts

TestHosts list of clients and servers hosts for use in the test(s)

| Field | Description | Scheme | Required |
| ----- | ----------- | ------ | -------- |
| Clients |  | [][Hosts](#hosts) | true |
| Servers |  | [][Hosts](#hosts) | true |

[Back to TOC](#table-of-contents)
