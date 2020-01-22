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
	// Version right now is just `0`, so we can keep track of config structure versioning.
	Version string `yaml:"version"`
	// Runner Runner configuration to use.
	Runner Runner `yaml:"runner"`
	// Tests List of `Test`s to run.
	Tests []*Test `yaml:"tests" validate:"required,min=1"`
}

// New return a new Config object with the `Version` set by default
func New() *Config {
	return &Config{
		Version: "0",
		Runner:  Runner{},
		Tests:   []*Test{},
	}
}

// Hosts options for hosts selection for a Test
type Hosts struct {
	// Name of this hosts selection.
	Name string `yaml:"name"`
	// If all hosts available should be used (default: `false`).
	All *bool `yaml:"all,omitempty"`
	// Select `Count` Random hosts from the available hosts list (default: `false`).
	Random *bool `yaml:"random,omitempty"`
	// Must be used with `Random`, will cause `Count` times Nodes to be randomly selected from all applicable hosts.
	Count int `yaml:"count"`
	// Static list of hosts (this list is not checked for accuracy)
	Hosts []string `yaml:"hosts"`
	// "Label" selector for the dynamically generated hosts list, e.g., Kubernetes label selector
	HostSelector map[string]string `yaml:"hostSelector"`
	// AntiAffinity **not implemented yet**
	AntiAffinity []string `yaml:"antiAffinity,omitempty"`
}

// Output Output config structure pointing to the other config options for each output
type Output struct {
	// Name of this output
	Name string `yaml:"name" validate:"required,min=3"`
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
	// Transformations transformations to be applied to the output data for the chosen output
	Transformations []*Transformation `yaml:"transformations,omitempty"`
}

// FilePath file path and name pattern for outputs file generation
type FilePath struct {
	// File base path for output
	FilePath string `yaml:"filePath" validate:"required,min=1"`
	// File name pattern templated from various availables during output generation
	NamePattern string `yaml:"namePattern" validate:"required,min=1"`
}

// TransformationAction Transformation action value
type TransformationAction string

const (
	// TransformationActionAdd Transformation add action
	TransformationActionAdd TransformationAction = "add"
	// TransformationActionReplace Transformation replace action
	TransformationActionReplace TransformationAction = "replace"
	// TransformationActionDelete Transformation delete action
	TransformationActionDelete TransformationAction = "delete"
)

// IsValidTransformationAction function to check if a TransformationAction is valid
// ("in range of available transformationactions")
func IsValidTransformationAction(a TransformationAction) bool {
	switch a {
	case TransformationActionAdd:
	case TransformationActionReplace:
	case TransformationActionDelete:
	default:
		return false
	}
	return true
}

// ModifierAction action to use with a modifier (float64)
type ModifierAction string

const (
	// ModifierActionMultiply Modifier multiply action
	ModifierActionMultiply ModifierAction = "multiply"
	// ModifierActionDivison Modifier divison action
	ModifierActionDivison ModifierAction = "division"
	// ModifierActionAddition Modifier addition action
	ModifierActionAddition ModifierAction = "addition"
	// ModifierActionSubstract Modifier substract action
	ModifierActionSubstract ModifierAction = "substract"
)

// Transformation data transformation instructions
type Transformation struct {
	// Source name of the (data) column to use for the transformation
	Source string `yaml:"key" validate:"required"`
	// Action transformation action to use on the Source key
	Action TransformationAction `yaml:"action" validate:"required,min=3"`
	// Destination used for the "replace" TransformationAction for targetting the key to overwrite
	Destination string `yaml:"from,omitempty"`
	// Modifier value to use in combination with the ModifierAction to modify the values (e.g., settin git to `1000` and ModifierAction to `divison` will divise the value by 1000)
	Modifier *float64 `yaml:"modifier,omitempty"`
	// ModifierAction action to run on the values together with the Modifier
	ModifierAction ModifierAction `yaml:"modifierAction"`
}

// CSV CSV Output config options
type CSV struct {
	// FilePath struct fields which are inherited by this struct.
	// The fields of the FilePath struct must be written directly to this struct.
	FilePath `yaml:",inline"`
	// Separator which rune to use as a separator in the CSV file (default: `;`).
	Separator *rune `yaml:"separator"`
}

// GoChart GoChart Output config options
type GoChart struct {
	// FilePath struct fields which are inherited by this struct.
	// The fields of the FilePath struct must be written directly to this struct.
	FilePath `yaml:",inline"`
	// Graphs definitions of graphs to produce from the testers output data
	Graphs []*GoChartGraph `yaml:"graphs" validate:"required,min=1"`
}

// GoChartGraph Type and columns for a one or two Y-axis chart (+ X-axis) to be generated based on this information.
type GoChartGraph struct {
	// TimeColumn column with the time / interval to use for the X-axis
	TimeColumn string `yaml:"timeColumn" validate:"required"`
	// LeftY name of the column / data column to use for the the left Y axis
	LeftY string `yaml:"leftY,omitempty"`
	// RightY name of the column / data column to use for the the right Y axis
	RightY string `yaml:"rightY" validate:"required"`
	// WithLinearRegression if a linear regression series should be added to each data series (default: `false`).
	WithLinearRegression *bool `yaml:"withLinearRegression,omitempty"`
	// WithSimpleMovingAverage if a simple moving average should be added to each data series(default: `false`).
	WithSimpleMovingAverage *bool `yaml:"withSimpleMovingAverage,omitempty"`
}

// Dump Dump Output config options
type Dump struct {
	// FilePath struct fields which are inherited by this struct.
	// The fields of the FilePath struct must be written directly to this struct.
	FilePath `yaml:",inline"`
}

// Excelize Excelize Output config options. TODO implement
type Excelize struct {
	// FilePath struct fields which are inherited by this struct.
	// The fields of the FilePath struct must be written directly to this struct.
	FilePath `yaml:",inline"`
	// After what amount of rows the Excel file should be saved (default: `1`)
	SaveAfterRows int `yaml:"saveAfterRows,omitempty" validate:"required,min=1"`
}

// SQLite SQLite Output config options
type SQLite struct {
	// FilePath struct fields which are inherited by this struct.
	// The fields of the FilePath struct must be written directly to this struct.
	FilePath `yaml:",inline"`
	// Pattern used for templating the name of the table used in the SQLite database, the tables are created automatically
	TableNamePattern string `yaml:"tableNamePattern"`
}

// MySQL MySQL Output config options
type MySQL struct {
	// MySQL DSN, format `[username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]`, for more information see [GitHub go-sql-driver/mysql - DSN (Data Source Name)](https://github.com/go-sql-driver/mysql#dsn-data-source-name)
	DSN string `yaml:"dsn"`
	// Pattern used for templating the name of the table used in the MySQL database, the tables are created automatically when MySQL.AutoCreateTables is set to `true`
	TableNamePattern string `yaml:"tableNamePattern"`
	// Automatically create tables in the MySQL database (default: `true`)
	AutoCreateTables *bool `yaml:"autoCreateTables,omitempty"`
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
	// Path to your kubeconfig file, if not set the following order will be tried out, `KUBECONFIG` and `$HOME/.kube/config`
	Kubeconfig string `yaml:"kubeconfig,omitempty"`
	// The image used for the spawned Pods for the tests (default: `quay.io/galexrt/container-toolbox`)
	Image string `yaml:"image,omitempty"`
	// Namespace to execute the tests in
	Namespace string `yaml:"namespace" validate:"max=63"`
	// If `hostNetwork` mode should be used for the test Pods
	HostNetwork *bool `yaml:"hostNetwork,omitempty"`
	// Timeout settings for operations against the Kubernetes API
	Timeouts *KubernetesTimeouts `yaml:"timeouts,omitempty"`
	// Annotations to put on the test Pods
	Annotations map[string]string `yaml:"annotations,omitempty"`
	// Host selection specific options
	Hosts *KubernetesHosts `yaml:"hosts,omitempty"`
	// ServiceAccounst to use server and client Pods
	ServiceAccounts *KubernetesServiceAccounts `yaml:"serviceaccounts,omitempty"`
}

// KubernetesTimeouts timeouts for operations with the Kubernetess API (in secconds)
type KubernetesTimeouts struct {
	// Timeout for object deletion in seconds (default: `20`)
	DeleteTimeout int `yaml:"deleteTimeout,omitempty"`
	// Timeout for "Pod running" check in seconds (default: `60`)
	RunningTimeout int `yaml:"runningTimeout,omitempty"`
	// Timeout for "Pod succeded" check in seconds (e.g., client Pod exits after Pod; default: `60`)
	SucceedTimeout int `yaml:"succeedTimeout,omitempty"`
}

// KubernetesHosts hosts selection options for Kubernetes
type KubernetesHosts struct {
	// If Nodes that are `SchedulingDisabled` should be ignored (default: `true`)
	IgnoreSchedulingDisabled *bool `yaml:"ignoreSchedulingDisabled,omitempty"`
	// List of Kubernetes corev1.Toleration to tolerate when selecting Nodes
	Tolerations []corev1.Toleration `yaml:"tolerations,omitempty"`
}

// KubernetesServiceAccounts server and client ServiceAccount name to use for the created Pods
type KubernetesServiceAccounts struct {
	// Server ServiceAccount name to use for server Pods
	Server string `yaml:"server,omitempty"`
	// Clients ServiceAccount name to use for client Pods
	Clients string `yaml:"clients,omitempty"`
}

// RunnerAnsible Ansible Runner config options
type RunnerAnsible struct {
	// InventoryFilePath Path to inventory file to use
	InventoryFilePath string `yaml:"inventoryFilePath"`
	// Groups server and clients group names
	Groups *AnsibleGroups `yaml:"groups"`
	// Path to the ansible command (if empty will be searched for in `PATH`; default: `ansble`)
	AnsibleCommand string `yaml:"ansibleCommand,omitempty"`
	// Path to the ansible-inventory command (if empty will be searched for in `PATH`; default: `ansble-inventory`)
	AnsibleInventoryCommand string `yaml:"ansibleInventoryCommand,omitempty"`
	// Timeout settings for ansible command runs
	Timeouts *AnsibleTimeouts `yaml:"timeouts,omitempty"`
	// CommandRetries amount of tries before to fail waiting for the server (main) task to start (default: `10`)
	CommandRetries *int `yaml:"commandRetries,omitempty"`
	// ParallelHostFactCalls the amount of host facts calls to make in parallel (default: `7`)
	ParallelHostFactCalls *int `yaml:"parallelHostFactCalls,omitempty"`
}

// AnsibleGroups server and clients host group names in the used inventory file(s)
type AnsibleGroups struct {
	// Server inventory server group name
	Server string `yaml:"server"`
	// Clients inventory clients group name
	Clients string `yaml:"clients"`
}

// AnsibleTimeouts timeouts for Ansible command runs
type AnsibleTimeouts struct {
	// Timeout duration for `ansible` and `ansible-inventory` calls (NOT task command timeouts; default: `20s`)
	CommandTimeout time.Duration `yaml:"commandTimeout,omitempty"`
	// Timeout duration for `ansible` Task command calls (default: `45s`)
	TaskCommandTimeout time.Duration `yaml:"taskCommandTimeout,omitempty"`
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
	RunOptions RunOptions `yaml:"runOptions,omitempty"`
	// List of Outputs to use for processing data from the testers.
	Outputs []Output `yaml:"outputs" validate:"required,min=1"`
	// Transformations transformations to be applied to Output data
	Transformations []*Transformation `yaml:"transformations,omitempty"`
	// Hosts selection for client and server
	Hosts TestHosts `yaml:"hosts"`
	// IPerf3 test options
	IPerf3 *IPerf3 `yaml:"iperf3"`
}

// RunMode custom run mode const type for
type RunMode string

const (
	// RunModeSequential run tasks in sequential / serial order
	RunModeSequential RunMode = "sequential"
	// RunModeParallel run tasks in parallel (WARNING! Be sure what you cause with this, e.g., 100 iperfs might not be good for a production environment)
	RunModeParallel RunMode = "parallel"
)

// RunOptions options for running the tasks
type RunOptions struct {
	// Continue on error during test runs (recommended to set to `true`) (default: is `true`)
	ContinueOnError *bool `yaml:"continueOnError,omitempty"`
	// Amount of test rounds (repetitions) to do for a test plan (default: `1`)
	Rounds int `yaml:"rounds,omitempty"`
	// Time interval to sleep / wait between (default: `10s`)
	Interval time.Duration `yaml:"interval,omitempty"`
	// Run mode can be `parallel` or `sequential` (see `RunMode`, default: is `sequential`)
	Mode RunMode `yaml:"mode,omitempty"`
	// **NOT IMPLEMENTED YET** amount of test tasks to run when using `RunModeParallel` (value: `parallel`).
	ParallelCount int `yaml:"parallelCount,omitempty"`
}

// TestHosts list of clients and servers hosts for use in the test(s)
type TestHosts struct {
	// Static list of hosts to use as clients
	Clients []Hosts `yaml:"clients" validate:"required,min=1"`
	// Static list of hosts to use as server
	Servers []Hosts `yaml:"servers" validate:"required,min=1"`
}

// AdditionalFlags additional flags structure for Server and Clients
type AdditionalFlags struct {
	// List of additional flags for clients
	Clients []string `yaml:"clients,omitempty"`
	// List of additional flags for server
	Server []string `yaml:"server,omitempty"`
}

// IPerf3 IPerf3 config structure for testers.Tester config
type IPerf3 struct {
	// Additional flags for client and server
	AdditionalFlags AdditionalFlags `yaml:"additionalFlags,omitempty"`
	// Duration Time in seconds the IPerf3 test should transmit / receive (default: `10`).
	// In case of the Ansible Runner, you need to increase the Ansible runners `timeouts.taskCommandTimeout` option when
	// increasing the Duration. The Ansible Runner `timeouts.taskCommandTimeout` option should be set to `Duration + some extra time`
	// (e.g., 10 seconds).
	Duration *int `yaml:"duration,omitempty" validate:"required,min=1"`
	// Interval Interval in seconds which IPerf3 will print / return periodic throughput reports (default: `1`).
	Interval *int `yaml:"interval,omitempty" validate:"required,min=1"`
	// If UDP should be used for the IPerf3 test
	UDP *bool `yaml:"udp,omitempty"`
}
