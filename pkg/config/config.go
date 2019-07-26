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

import (
	"time"

	corev1 "k8s.io/api/core/v1"
)

// Config Config object for the config file
type Config struct {
	Version string  `yaml:"version"`
	Runner  Runner  `yaml:"runner"`
	Tests   []*Test `yaml:"tests"`
}

// KeyValuePair key value string pair
type KeyValuePair struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

// New return a new Config object with the `Version` set by default
func New() *Config {
	return &Config{
		Runner:  Runner{},
		Tests:   []*Test{},
		Version: "0",
	}
}

// Hosts options for hosts selection for a Test
type Hosts struct {
	// Name of this hosts selection.
	Name string `yaml:"name"`
	// If all hosts available should be used.
	All bool `yaml:"all"`
	// Select `Count` Random hosts from the available hosts list.
	Random bool `yaml:"random"`
	// Used with Random to randomly select the Count of hosts.
	Count        int               `yaml:"count"`
	Hosts        []string          `yaml:"hosts"`
	HostSelector map[string]string `yaml:"hostSelector"`
	// AntiAffinity not implemented yet
	AntiAffinity []KeyValuePair `yaml:"antiAffinity"`
}

// Output Output config structure pointing to the other config options for each output
type Output struct {
	Name     string    `yaml:"name"`
	CSV      *CSV      `yaml:"csv"`
	GoChart  *GoChart  `yaml:"goChart"`
	Dump     *Dump     `yaml:"dump"`
	Excelize *Excelize `yaml:"excelize"`
	SQLite   *SQLite   `yaml:"sqlite"`
	MySQL    *MySQL    `yaml:"mysql"`
}

// CSV CSV Output config options
type CSV struct {
	FilePath    string `yaml:"filePath"`
	NamePattern string `yaml:"namePattern"`
}

// GoChart GoChart Output config options
type GoChart struct {
	FilePath    string   `yaml:"filePath"`
	NamePattern string   `yaml:"namePattern"`
	Types       []string `yaml:"types"`
}

// Dump Dump Output config options
type Dump struct {
	FilePath    string `yaml:"filePath"`
	NamePattern string `yaml:"namePattern"`
}

// Excelize Excelize Output config options. TODO implement
type Excelize struct {
	FilePath    string `yaml:"filePath"`
	NamePattern string `yaml:"namePattern"`
}

// SQLite SQLite Output config options
type SQLite struct {
	FilePath         string `yaml:"filePath"`
	NamePattern      string `yaml:"namePattern"`
	TableNamePattern string `yaml:"tableNamePattern"`
}

// MySQL MySQL Output config options
type MySQL struct {
	DSN              string `yaml:"dsn"`
	TableNamePattern string `yaml:"tableNamePattern"`
}

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

// KubernetesHosts hosts selection options for Kubernetes
type KubernetesHosts struct {
	IgnoreSchedulingDisabled bool                `yaml:"ignoreSchedulingDisabled"`
	Tolerations              []corev1.Toleration `yaml:"tolerations"`
}

// RunnerMock Mock Runner config options (here for good measure)
type RunnerMock struct {
}

// Test Config options for each Test
type Test struct {
	Name       string     `yaml:"name"`
	Type       string     `yaml:"type"`
	RunOptions RunOptions `yaml:"runOptions"`
	Outputs    []Output   `yaml:"outputs"`
	Hosts      TestHosts  `yaml:"hosts"`
	IPerf3     *IPerf3    `yaml:"iperf3"`
	Siege      *Siege     `yaml:"siege"`
	Smokeping  *Smokeping `yaml:"smokeping"`
}

const (
	// RunModeSequential run tasks in sequential / serial order
	RunModeSequential = "sequential"
	// RunModeParallel run tasks in parallel (WARNING! Be sure what you cause with this, e.g., 100 iperfs might not be good for a production environment)
	RunModeParallel = "parallel"
)

// RunOptions options for running the tasks
type RunOptions struct {
	ContinueOnError bool          `yaml:"continueOnError"`
	Rounds          int           `yaml:"rounds"`
	Interval        time.Duration `yaml:"interval"`
	Mode            string        `yaml:"mode"`
	ParallelCount   int           `yaml:"parallelCount"`
}

// TestHosts list of clients and servers hosts for use in the test(s)
type TestHosts struct {
	Clients []Hosts `yaml:"clients"`
	Servers []Hosts `yaml:"servers"`
}

// AdditionalFlags additional flags structure for Server and Clients
type AdditionalFlags struct {
	Clients []string `yaml:"clients"`
	Server  []string `yaml:"server"`
}

// IPerf3 IPerf3 config structure for testers.Tester config
type IPerf3 struct {
	AdditionalFlags AdditionalFlags `yaml:"additionalFlags"`
	UDP             *bool           `yaml:"udp"`
}

// Siege Siege config structure TODO not implemented yet
type Siege struct {
	AdditionalFlags AdditionalFlags   `yaml:"additionalFlags"`
	Benchmark       bool              `yaml:"benchmark"`
	Headers         map[string]string `yaml:"headers"`
	URLs            []string          `yaml:"urls"`
	UserAgent       string            `yaml:"userAgent"`
	// TODO Add more options from SIEGERC config file
}

// Smokeping Smokeping config structure TODO not implemented yet
type Smokeping struct {
	AdditionalFlags AdditionalFlags `yaml:"additionalFlags"`
}
