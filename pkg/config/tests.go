/*
Copyright 2019 Cloudical Deutschland GmbH
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

// Test
type Test struct {
	Type        string      `yaml:"type"`
	RunOptions  RunOptions  `yaml:"runOptions"`
	HostOptions HostOptions `yaml:"hostOptions"`
	Hosts       TestHosts   `yaml:"hosts"`
	IPerf       *IPerf      `yaml:"iperf"`
	IPerf3      *IPerf3     `yaml:"iperf3"`
	Siege       *Siege      `yaml:"siege"`
}

// RunOptions
type RunOptions struct {
	Rounds   int    `yaml:"rounds"`
	Interval string `yaml:"interval"`
}

// HostOptions
type HostOptions struct {
	All bool `yaml:"all"`
}

// TestHosts
type TestHosts struct {
	Sources      []Hosts `yaml:"sources"`
	Destinations []Hosts `yaml:"destinations"`
}

// IPerf
type IPerf struct {
	WindowSizeCalculation IPerfWindowSizeCalculation `yaml:"windowSizeCalculation"`
	AdditionalFlags       IPerfAdditionalFlags       `yaml:"additionalFlags"`
}

// IPerfWindowSizeCalculation
type IPerfWindowSizeCalculation struct {
	Auto bool `yaml:"auto"`
}

// IPerfAdditionalFlags
type IPerfAdditionalFlags struct {
	Client []string `yaml:"client"`
	Server []string `yaml:"server"`
}

// IPerf3
type IPerf3 struct {
}

// Siege
type Siege struct {
}
