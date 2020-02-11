# Runners: Ansible

```bash
# Check the generated plan and confirm by typing 'yes'
ancientt
# To just print the plan
ancientt --only-print-plan
# Generate and execute the plan without user prompt
ancientt --yes
```

## Depdency Playbooks

* [`depdendency-iperf3.yaml`](dependency-iperf3.yaml) - Installs `iperf3` and `procps` package for IPerf3 testing.
* [`depdendency-pingparsing.yaml`](dependency-pingparsing.yaml) - Installs `pingparsing` (using `pip`) and `procps` package for IPerf3 testing. This also checks that Python 3.5+ is used for Ansible on the target hosts.
