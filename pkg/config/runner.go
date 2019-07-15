/*
Copyright 2019 Cloudical Deutschland GmbH. All rights reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

// Runner structure with all available runners config options
type Runner struct {
	Name       string            `yaml:"name"`
	Kubernetes *RunnerKubernetes `yaml:"kubernetes"`
	Mock       *RunnerMock       `yaml:"mock"`
}

// RunnerKubernetes Kubernetes Runner config options
type RunnerKubernetes struct {
	Kubeconfig  string              `yaml:"kubeconfig"`
	Image       string              `yaml:"image"`
	Namespace   string              `yaml:"namespace"`
	HostNetwork bool                `yaml:"hostNetwork"`
	Timeouts    *KubernetesTimeouts `yaml:"timeouts"`
	Annotations map[string]string   `yaml:"annotations"`
	Hosts       *KubernetesHosts    `yaml:"hosts"`
}

// KubernetesTimeouts timeouts for operations with Kubernetess
type KubernetesTimeouts struct {
	DeleteTimeout  int `yaml:"deleteTimeout"`
	RunningTimeout int `yaml:"runningTimeout"`
	SucceedTimeout int `yaml:"succeedTimeout"`
}

// KubernetesHosts
type KubernetesHosts struct {
	IgnoreSchedulingDisabled bool `yaml:"ignoreSchedulingDisabled"`
}

// RunnerMock Mock Runner config options (here for good measure)
type RunnerMock struct {
}
