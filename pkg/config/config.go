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
	Count int `yaml:"count"`
	// Static list of hosts (this list is not checked for accuracy)
	Hosts []string `yaml:"hosts"`
	// "Label" selector for the dynamically generated hosts list, e.g., Kubernetes label selector
	HostSelector map[string]string `yaml:"hostSelector"`
	// AntiAffinity not implemented yet
	AntiAffinity []string `yaml:"antiAffinity"`
}

// Output Output config structure pointing to the other config options for each output
type Output struct {
	// Name of this output
	Name string `yaml:"name"`
	// CSV output options
	CSV *CSV `yaml:"csv"`
	// GoChart output options
	GoChart *GoChart `yaml:"goChart"`
	// Dump output options
	Dump *Dump `yaml:"dump"`
	// Excelize output options
	Excelize *Excelize `yaml:"excelize"`
	// SQLite output options
	SQLite *SQLite `yaml:"sqlite"`
	// MySQL output options
	MySQL *MySQL `yaml:"mysql"`
}

// CSV CSV Output config options
type CSV struct {
	// File base path for output
	FilePath string `yaml:"filePath"`
	// File name pattern templated from various availables during output generation
	NamePattern string `yaml:"namePattern"`
}

// GoChart GoChart Output config options
type GoChart struct {
	// File base path for output
	FilePath string `yaml:"filePath"`
	// File name pattern templated from various availables during output generation
	NamePattern string `yaml:"namePattern"`
	// Types of charts to produce from the testers output data
	Types []string `yaml:"types"`
}

// Dump Dump Output config options
type Dump struct {
	// File base path for output
	FilePath string `yaml:"filePath"`
	// File name pattern templated from various availables during output generation
	NamePattern string `yaml:"namePattern"`
}

// Excelize Excelize Output config options. TODO implement
type Excelize struct {
	// File base path for output
	FilePath string `yaml:"filePath"`
	// File name pattern templated from various availables during output generation
	NamePattern string `yaml:"namePattern"`
	// After what amount of rows the Excel file should be saved
	SaveAfterRows int `yaml:"saveAfterRows"`
}

// SQLite SQLite Output config options
type SQLite struct {
	// File base path for output
	FilePath string `yaml:"filePath"`
	// File name pattern templated from various availables during output generation
	NamePattern string `yaml:"namePattern"`
	// Pattern used for templating the name of the table used in the SQLite database, the tables are created automatically
	TableNamePattern string `yaml:"tableNamePattern"`
}

// MySQL MySQL Output config options
type MySQL struct {
	// MySQL DSN, format `[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]`, for more information see https://github.com/go-sql-driver/mysql#dsn-data-source-name
	DSN string `yaml:"dsn"`
	// Pattern used for templating the name of the table used in the MySQL database, the tables are created automatically when MySQL.AutoCreateTables is set to `true`
	TableNamePattern string `yaml:"tableNamePattern"`
	// Automatically create tables in the MySQL database (default `true`)
	AutoCreateTables *bool `yaml:"autoCreateTables"`
}

// Runner structure with all available runners config options
type Runner struct {
	// Name of the runner
	Name string `yaml:"name"`
	// Kubernetes runner options
	Kubernetes *RunnerKubernetes `yaml:"kubernetes"`
	// Ansible runner options
	Ansible *RunnerAnsible `yaml:"ansible"`
	// Mock runner options (userd for testing purposes)
	Mock *RunnerMock `yaml:"mock"`
}

// RunnerKubernetes Kubernetes Runner config options
type RunnerKubernetes struct {
	// If the Kubernetes client should use the in-cluster config for the cluster communication
	InClusterConfig bool `yaml:"inClusterConfig"`
	// Path to your kubeconfig file, if not set the `KUBECONFIG` env var will be used and then the default
	Kubeconfig string `yaml:"kubeconfig"`
	// The image used for the spawned Pods for the tests (default: `quay.io/galexrt/container-toolbox`)
	Image string `yaml:"image"`
	// Namespace to execute the tests in
	Namespace string `yaml:"namespace"`
	// If `hostNetwork` mode should be used for the test Pods
	HostNetwork bool `yaml:"hostNetwork"`
	// Timeout settings for operations against the Kubernetes API
	Timeouts *KubernetesTimeouts `yaml:"timeouts"`
	// Annotations to put on the test Pods
	Annotations map[string]string `yaml:"annotations"`
	// Host selection specific options
	Hosts *KubernetesHosts `yaml:"hosts"`
}

// KubernetesTimeouts timeouts for operations with Kubernetess
type KubernetesTimeouts struct {
	// Timeout for object deletion
	DeleteTimeout int `yaml:"deleteTimeout"`
	// Timeout for "Pod running" check
	RunningTimeout int `yaml:"runningTimeout"`
	// Timeout for "Pod succeded" check (e.g., client Pod exits after Pod)
	SucceedTimeout int `yaml:"succeedTimeout"`
}

// KubernetesHosts hosts selection options for Kubernetes
type KubernetesHosts struct {
	// If Nodes that are `SchedulingDisabled` should be ignored
	IgnoreSchedulingDisabled bool `yaml:"ignoreSchedulingDisabled"`
	// List of Kubernetes corev1.Toleration to tolerate when selecting Nodes
	Tolerations []corev1.Toleration `yaml:"tolerations"`
}

// RunnerAnsible Ansible Runner config options
type RunnerAnsible struct {
	// InventoryFilePath Path to inventory file to use
	InventoryFilePath string `yaml:"inventoryFilePath"`
	// Groups server and clients group names
	Groups *AnsibleGroups `yaml:"groups"`
	// Path to the ansible command (if empty will be searched for in `PATH`)
	AnsibleCommand string `yaml:"ansibleCommand"`
	// Path to the ansible-inventory command (if empty will be searched for in `PATH`)
	AnsibleInventoryCommand string `yaml:"ansibleInventoryCommand"`
	// Timeout duration for `ansible` and `ansible-inventory` calls (NOT task command timeouts)
	CommandTimeout time.Duration `yaml:"commandTimeout"`
	// Timeout duration for `ansible` Task command calls
	TaskCommandTimeout time.Duration `yaml:"taskCommandTimeout"`
}

// AnsibleGroups server and clients host group names in the used inventory file(s)
type AnsibleGroups struct {
	// Server inventory server group name
	Server string `yaml:"server"`
	// Clients inventory clients group name
	Clients string `yaml:"clients"`
}

// RunnerMock Mock Runner config options (here for good measure)
type RunnerMock struct {
}

// Test Config options for each Test
type Test struct {
	// Test name
	Name string `yaml:"name"`
	// The tester to use, e.g., for `iperf3` set to `iperf3` and so on
	Type string `yaml:"type"`
	// Options for the execution of the test
	RunOptions RunOptions `yaml:"runOptions"`
	// List of Outputs to use for processing data from the testers.
	Outputs []Output `yaml:"outputs"`
	// Hosts selection for client and server
	Hosts TestHosts `yaml:"hosts"`
	// IPerf3 test options
	IPerf3 *IPerf3 `yaml:"iperf3"`
	// **NOT IMPLEMENTED** Siege test options
	Siege *Siege `yaml:"siege"`
}

const (
	// RunModeSequential run tasks in sequential / serial order
	RunModeSequential = "sequential"
	// RunModeParallel run tasks in parallel (WARNING! Be sure what you cause with this, e.g., 100 iperfs might not be good for a production environment)
	RunModeParallel = "parallel"
)

// RunOptions options for running the tasks
type RunOptions struct {
	//
	ContinueOnError bool `yaml:"continueOnError"`
	// Amount of test rounds (repetitions) to do for a test plan
	Rounds int `yaml:"rounds"`
	// Time interval to sleep / wait between
	Interval time.Duration `yaml:"interval"`
	// Run mode can be `parallel` or `sequential` (default is `sequential`)
	Mode string `yaml:"mode"`
	// **NOT IMPLEMENTED YET** amount of test tasks to run when using `parallel` RunOptions.Mode
	ParallelCount int `yaml:"parallelCount"`
}

// TestHosts list of clients and servers hosts for use in the test(s)
type TestHosts struct {
	// Static list of hosts to use as clients
	Clients []Hosts `yaml:"clients"`
	// Static list of hosts to use as server
	Servers []Hosts `yaml:"servers"`
}

// AdditionalFlags additional flags structure for Server and Clients
type AdditionalFlags struct {
	//  List of additional flags for clients
	Clients []string `yaml:"clients"`
	//  List of additional flags for server
	Server []string `yaml:"server"`
}

// IPerf3 IPerf3 config structure for testers.Tester config
type IPerf3 struct {
	// Additional flags for client and server
	AdditionalFlags AdditionalFlags `yaml:"additionalFlags"`
	// If UDP should be used for the IPerf3 test
	UDP *bool `yaml:"udp"`
}

// Siege Siege config structure TODO not implemented yet
type Siege struct {
	// Additional flags for client and server
	AdditionalFlags AdditionalFlags   `yaml:"additionalFlags"`
	Benchmark       bool              `yaml:"benchmark"`
	Headers         map[string]string `yaml:"headers"`
	URLs            []string          `yaml:"urls"`
	UserAgent       string            `yaml:"userAgent"`
	// TODO Add more options from SIEGERC config file
}
