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

// Output Output config structure pointing to the other config options for each output
type Output struct {
	Name     string    `yaml:"name"`
	CSV      *CSV      `yaml:"csv"`
	GoChart  *GoChart  `yaml:"goChart"`
	Dump     *Dump     `yaml:"dump"`
	Excelize *Excelize `yaml:"excelize"`
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
