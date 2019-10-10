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

	"github.com/cloudical-io/ancientt/pkg/util"

	"github.com/cloudical-io/ancientt/pkg/ansible"
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
		c.IgnoreSchedulingDisabled = util.BoolPointer(true)
	}
}

// SetDefaults set defaults on config part
func (c *KubernetesTimeouts) SetDefaults() {
	if c.DeleteTimeout == 0 {
		c.DeleteTimeout = 20
	}
	if c.RunningTimeout == 0 {
		c.RunningTimeout = 35
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
