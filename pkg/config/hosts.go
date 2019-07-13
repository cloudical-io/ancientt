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

// Hosts options for hosts selection for a Test
type Hosts struct {
	Name         string            `yaml:"name"`
	All          bool              `yaml:"all"`
	Random       bool              `yaml:"random"`
	Count        int               `yaml:"count"`
	Hosts        []string          `yaml:"hosts"`
	HostSelector map[string]string `yaml:"hostSelector"`
	// AntiAffinity not implemented yet
	AntiAffinity []KeyValuePair `yaml:"antiAffinity"`
}
