# Config Examples

## Kubernetes + IPerf3 between all hosts

```yaml
version: '0'
runner:
  #name: mock
  name: kubernetes
  kubernetes:
    # Assuming you are in your home directory
    kubeconfig: .kube/config
    image: quay.io/galexrt/container-toolbox
    namespace: acntt
    timeouts:
      deleteTimeout: 20
      runningTimeout: 35
      succeedTimeout: 60
    hosts:
      ignoreSchedulingDisabled: true
      tolerations: []
tests:
- name: iperf3-one-rand-to-one-rand
  type: iperf3
  outputs:
  - name: csv
    csv:
      filePath: .
      namePattern: 'acntt-{{ .TestStartTime }}-{{ .Data.Tester }}-{{ .Data.ServerHost }}_{{ .Data.ClientHost }}.csv'
  runOptions:
    continueOnError: true
    rounds: 2
    # wait 10 seconds between each round
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
    additionalFlags:
      clients: []
      server: []
```