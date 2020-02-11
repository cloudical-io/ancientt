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

	"github.com/cloudical-io/ancientt/pkg/ansible"
	"github.com/cloudical-io/ancientt/pkg/util"
)

// Defaults interface to implement for config parts which allow a "verification" / Setting Defaults
type Defaults interface {
	Defaults()
}

// SetDefaults set defaults on config part
func (c *RunnerKubernetes) SetDefaults() {
	if c.Annotations == nil {
		c.Annotations = map[string]string{}
	}

	if c.Image == "" {
		c.Image = "quay.io/galexrt/container-toolbox:latest"
	}

	if c.Namespace == "" {
		c.Namespace = "ancientt"
	}
}

// SetDefaults set defaults on config part
func (c *KubernetesHosts) SetDefaults() {
	if c.IgnoreSchedulingDisabled == nil {
		c.IgnoreSchedulingDisabled = util.BoolTruePointer()
	}
	if c.Tolerations == nil {
		c.Tolerations = []corev1.Toleration{}
	}
}

// SetDefaults set defaults on config part
func (c *KubernetesTimeouts) SetDefaults() {
	if c.DeleteTimeout == 0 {
		c.DeleteTimeout = 20
	}
	if c.RunningTimeout == 0 {
		c.RunningTimeout = 60
	}
	if c.SucceedTimeout == 0 {
		c.SucceedTimeout = 60
	}
}

// SetDefaults set defaults on config part
func (c *RunnerAnsible) SetDefaults() {
	if c.AnsibleCommand == "" {
		c.AnsibleCommand = ansible.AnsibleCommand
	}
	if c.AnsibleInventoryCommand == "" {
		c.AnsibleInventoryCommand = ansible.AnsibleInventoryCommand
	}

	if c.Timeouts == nil {
		c.Timeouts = &AnsibleTimeouts{}
	}
	if c.Timeouts.CommandTimeout == 0 {
		c.Timeouts.CommandTimeout = 20 * time.Second
	}
	if c.Timeouts.TaskCommandTimeout == 0 {
		c.Timeouts.TaskCommandTimeout = 45 * time.Second
	}

	if c.CommandRetries == nil || *c.CommandRetries == 0 {
		defVal := 10
		c.CommandRetries = &defVal
	}
	if c.ParallelHostFactCalls == nil || *c.ParallelHostFactCalls == 0 {
		defVal := 7
		c.ParallelHostFactCalls = &defVal
	}

	if c.Groups == nil {
		c.Groups = &AnsibleGroups{}
	}
	if c.Groups.Server == "" {
		c.Groups.Server = "server"
	}
	if c.Groups.Clients == "" {
		c.Groups.Clients = "clients"
	}
}

// SetDefaults set defaults on config part
func (c *RunOptions) SetDefaults() {
	if c.ContinueOnError == nil {
		c.ContinueOnError = util.BoolTruePointer()
	}

	if c.Rounds == 0 {
		c.Rounds = 1
	}

	if c.Interval == 0 {
		c.Interval = 10 * time.Second
	}

	if c.Mode == "" {
		c.Mode = RunModeSequential
	}
}

// SetDefaults set defaults on config part
func (c *IPerf3) SetDefaults() {
	if c.Duration == nil {
		defValue := 10
		c.Duration = &defValue
	}

	if c.Interval == nil {
		defValue := 1
		c.Interval = &defValue
	}

	if c.UDP == nil {
		c.UDP = util.BoolFalsePointer()
	}
}

// SetDefaults set defaults on config part
func (c *PingParsing) SetDefaults() {
	if c.Count == nil {
		defValue := 10
		c.Count = &defValue
	}

	if c.Deadline == nil {
		defValue := 15 * time.Second
		c.Deadline = &defValue
	}

	if c.Timeout == nil {
		defValue := 10 * time.Second
		c.Timeout = &defValue
	}
}

// SetDefaults set defaults on config part
func (c *AdditionalFlags) SetDefaults() {
	if c.Server == nil {
		c.Server = []string{}
	}
	if c.Clients == nil {
		c.Clients = []string{}
	}
}

// SetDefaults set defaults on config part
func (c *Excelize) SetDefaults() {
	if c.SaveAfterRows == 0 {
		c.SaveAfterRows = 1
	}
}

// SetDefaults set defaults on config part
func (c *MySQL) SetDefaults() {
	if c.AutoCreateTables == nil {
		c.AutoCreateTables = util.BoolTruePointer()
	}
}

// SetDefaults set defaults on config part
func (c *CSV) SetDefaults() {
	if c.Separator == nil {
		semiColon := ';'
		c.Separator = &semiColon
	}
}

// SetDefaults set defaults on confg part
func (c *GoChartGraph) SetDefaults() {
	if c.WithLinearRegression == nil {
		c.WithLinearRegression = util.BoolFalsePointer()
	}
	if c.WithSimpleMovingAverage == nil {
		c.WithSimpleMovingAverage = util.BoolTruePointer()
	}
}
