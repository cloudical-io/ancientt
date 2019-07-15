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
)

// Test Config options for each Test
type Test struct {
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
	Rounds        int           `yaml:"rounds"`
	Interval      time.Duration `yaml:"interval"`
	Mode          string        `yaml:"mode"`
	ParallelCount int           `yaml:"parallelCount"`
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
