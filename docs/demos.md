# Demos

This page contains demos of `ancientt` in action.

Each demo has an [asciinema recording](https://asciinema.org/), the used `testdefinition.yaml` and depending on the used output format a short snippet or the whole output available.

## Kubernetes: `iperf3` One Random Node to All Kubernetes Nodes

[![asciicast](https://asciinema.org/a/kCpLvkjVRAMcyYraBz2ZIp5h6.svg)](https://asciinema.org/a/kCpLvkjVRAMcyYraBz2ZIp5h6)

`testdefinition.yaml`:
```yaml
version: '0'
runner:
  name: kubernetes
  kubernetes:
    kubeconfig: .kube/config
    image: quay.io/galexrt/container-toolbox
    hosts:
      ignoreSchedulingDisabled: true
tests:
- name: iperf3-one-to-all
  type: iperf3
  outputs:
  - name: csv
    csv:
      filePath: ./results
      namePattern: 'ancientt-{{ .TestStartTime }}-{{ .Data.Tester }}.csv'
  runOptions:
    continueOnError: true
    rounds: 1
  hosts:
    clients:
    - name: all-hosts
      all: true
    servers:
    - name: one-host
      count: 1
      random: true
  iperf3:
    additionalFlags:
      clients: ["--interval=1"]
      server: ["--interval=1"]
```

Output `csv (partial snippet):
```csv
test_time,round,tester,server_host,client_host,socket,start,end,seconds,bytes,bits_per_second,retransmits,snd_cwnd,rtt,rttvar,pmtu,omitted,iperf3_version,system_info,additional_info
2019-09-24T20:43:13+0200,0,iperf3,mycoolk8scluster-worker-01,mycoolk8scluster-worker-01,5,0.000000,1.000295,1.000295,1982632968,15856387318.277708,0,3191392,308,384,1500,false,iperf 3.6,Linux ancientt-client-iperf3-b411d4627abf6b743515539a060f8a84292e2267 4.15.0-54-generic #58-Ubuntu SMP Mon Jun 24 10:55:24 UTC 2019 x86_64,ancientt-client-iperf3-b411d4627abf6b743515539a060f8a84292e2267
2019-09-24T20:43:13+0200,0,iperf3,mycoolk8scluster-worker-01,mycoolk8scluster-worker-01,5,1.000295,2.000191,0.999896,2487746560,19904041515.077232,0,3191392,246,310,1500,false,iperf 3.6,Linux ancientt-client-iperf3-b411d4627abf6b743515539a060f8a84292e2267 4.15.0-54-generic #58-Ubuntu SMP Mon Jun 24 10:55:24 UTC 2019 x86_64,ancientt-client-iperf3-b411d4627abf6b743515539a060f8a84292e2267
2019-09-24T20:43:13+0200,0,iperf3,mycoolk8scluster-worker-01,mycoolk8scluster-worker-01,5,2.000191,3.000123,0.999932,2199388160,17596300936.243999,0,3191392,761,622,1500,false,iperf 3.6,Linux ancientt-client-iperf3-b411d4627abf6b743515539a060f8a84292e2267 4.15.0-54-generic #58-Ubuntu SMP Mon Jun 24 10:55:24 UTC 2019 x86_64,ancientt-client-iperf3-b411d4627abf6b743515539a060f8a84292e2267
2019-09-24T20:43:13+0200,0,iperf3,mycoolk8scluster-worker-01,mycoolk8scluster-worker-01,5,3.000123,4.000487,1.000364,2596536320,20764735793.626110,0,3191392,262,322,1500,false,iperf 3.6,Linux ancientt-client-iperf3-b411d4627abf6b743515539a060f8a84292e2267 4.15.0-54-generic #58-Ubuntu SMP Mon Jun 24 10:55:24 UTC 2019 x86_64,ancientt-client-iperf3-b411d4627abf6b743515539a060f8a84292e2267
[...]
2019-09-24T20:43:29+0200,0,iperf3,mycoolk8scluster-worker-01,mycoolk8scluster-worker-02,5,0.000000,1.000507,1.000507,176868460,1414230500.636069,819,1107216,5421,1645,1450,false,iperf 3.6,Linux ancientt-client-iperf3-6ca1e56296767da8d29adf6c8bba8ba6fc32fce4 4.15.0-54-generic #58-Ubuntu SMP Mon Jun 24 10:55:24 UTC 2019 x86_64,ancientt-client-iperf3-6ca1e56296767da8d29adf6c8bba8ba6fc32fce4
[...]
2019-09-24T20:43:43+0200,0,iperf3,mycoolk8scluster-worker-01,mycoolk8scluster-worker-03,5,9.000335,10.000137,0.999802,192675840,1541711805.372557,289,1077858,3027,325,1450,false,iperf 3.6,Linux ancientt-client-iperf3-af06030e7b7611f917582390fb1db0e38e799c08 4.15.0-54-generic #58-Ubuntu SMP Mon Jun 24 10:55:24 UTC 2019 x86_64,ancientt-client-iperf3-af06030e7b7611f917582390fb1db0e38e799c08
```

## More demos to come soon

There will soon be more demos, about the different runners, testers and outputs available in `ancientt`.