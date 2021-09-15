# Config Examples

This page contains some example configurations.

Be sure to also checkout the [Demos](demos.md) page for examples with example `ancientt` output, config and a snippet or whole file of the output results.

## Kubernetes + IPerf3 = CSV Output: IPerf3 test between all to all Nodes

```yaml
version: '0'
runner:
  #name: mock
  name: kubernetes
  kubernetes:
    # Assuming you are in your home directory
    kubeconfig: .kube/config
    image: quay.io/galexrt/container-toolbox:v20210915-101121-713
    namespace: ancientt
    timeouts:
      deleteTimeout: 20
      runningTimeout: 60
      succeedTimeout: 60
    hosts:
      ignoreSchedulingDisabled: true
      tolerations: []
tests:
- name: iperf3-one-rand-to-one-rand
  type: iperf3
  transformations:
  - source: "bits_per_second"
    destination: "gigabits_per_second"
    action: "add"
    modifier: 100000000
    modifierAction: "division"
  outputs:
  - name: csv
    csv:
      filePath: .
      namePattern: 'ancientt-{{ .TestStartTime }}-{{ .Data.Tester }}.csv'
      # If you want one CSV per server and client host test run, you can use the following:
      #namePattern: 'ancientt-{{ .TestStartTime }}-{{ .Data.Tester }}-{{ .Data.ServerHost }}_{{ .Data.ClientHost }}.csv'
  runOptions:
    continueOnError: true
    # If you wanna do the test(s) more than once in one go, set to higher than 1
    rounds: 1
    # Wait 10 seconds between each round
    interval: 10s
    mode: "sequential"
    parallelcount: 1
  # This hosts section would cause iperf3 to be run from all hosts to the hosts selected in the `destinations` section
  # Each entry will be merged into one list
  hosts:
    clients:
    - name: all-hosts
      all: true
    servers:
    - name: all-hosts
      all: true
  iperf3:
    udp: false
    duration: 10
    interval: 1
    additionalFlags:
      clients: []
      server: []
```
