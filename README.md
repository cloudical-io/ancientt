# acntt

Automated Continous network testing tool using existing projects like iperf3, siege, etc.

## Points to achieve

* A bit like Prometheus blackbox exporter which contains "definitions" for probes. The "tests" would be pluggable through a Golang interface.
* "Runner" interface, e.g., for Kubernetes, Ansible, etc. The "runner" abstracts the "how it is run", e.g., for Kubernetes creates a Job, Ansible (download and) trigger a playbook to run the test.
* Store result data in, e.g., graph database https://dgraph.io/ or a TSDB?
* Visualization for humans.
