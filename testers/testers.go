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

package testers

import (
	"fmt"
	"time"

	"github.com/cloudical-io/acntt/pkg/config"
)

// Factories contains the list of all available testers.
var Factories = make(map[string]func() (Tester, error))

// Tester is the interface a tester has to implement
type Tester interface {
	// Plan return a map of commands and on which servers to run them thanks to the info of the runners.Runner
	Plan(env *Environment, test *config.Test) (*Plan, error)
}

// Environment
type Environment struct {
	Hosts *Hosts
}

// TestHosts
type Hosts struct {
	Clients map[string]Host `yaml:"clients"`
	Servers map[string]Host `yaml:"servers"`
}

// Host
type Host struct {
	Name      string
	Labels    map[string]string
	Addresses IPAddresses
}

// IPAddresses list of IPv4 and IPv6 addresses a host has
type IPAddresses struct {
	IPv4 []string
	IPv6 []string
}

// Plan contains the information needed to execute the plan
type Plan struct {
	AffectedServers map[string]Host     `json:"affectedServers"`
	Commands        []map[string][]Task `json:"commands"`
	Tester          string              `json:"tester"`
}

// PrettyPrint "pretty" prints a plan
func (p Plan) PrettyPrint() {
	fmt.Println("-> BEGIN AffectedServers")
	for _, server := range p.AffectedServers {
		fmt.Println(server.Name)
	}
	fmt.Println("=> END AffectedServers")
	fmt.Println("-> BEGIN Commands")
	for k, command := range p.Commands {
		fmt.Printf("--> BEGIN Round %d\n", k)
		for server, tasks := range command {
			fmt.Printf("---> BEGIN Server %s will run\n", server)
			for _, task := range tasks {
				fmt.Printf("----> %s (Additional info: %+v; %+v)\n", task.Command, task.Ports, task.Sleep)
			}
			fmt.Printf("---> END Server %s will run\n", server)
		}
		fmt.Printf("--> END Round %d\n", k)
	}
	fmt.Println("=> END Commands")
}

// Task information for the task to execute
type Task struct {
	Command string        `json:"command"`
	Sleep   time.Duration `json:"sleep"`
	// TODO Implement CommandBuilder to build commands on the fly in the runner
	// This will be useful when the runner part can handle the "port assignment" / "port mapping".
	// Might be handled by the testers.Tester itself when, e.g., RunOptions.Mode `parallel` is used.
	CommandBuilder func(server Host, client Host) (string, error)
	Ports          Ports `json:"ports"`
}

// Ports TCP and UDP ports list
type Ports struct {
	TCP []int16
	UDP []int16
}
