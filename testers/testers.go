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

package testers

import (
	"fmt"
	"time"

	"github.com/cloudical-io/acntt/pkg/config"
)

// Factories contains the list of all available testers.
// The tester can each then be created using the function saved in the map.
var Factories = make(map[string]func(cfg *config.Config, test *config.Test) (Tester, error))

// Tester is the interface a tester has to implement
type Tester interface {
	// Plan return a map of commands and on which servers to run them thanks to the info of the runners.Runner
	Plan(env *Environment, test *config.Test) (*Plan, error)
}

// Environment environment information such as which hosts are doing what (clients, servers)
type Environment struct {
	Hosts *Hosts
}

// Hosts contains a list of clients and servers hosts that will be used in the test environment.
type Hosts struct {
	Clients map[string]*Host `yaml:"clients"`
	Servers map[string]*Host `yaml:"servers"`
}

// Host host information, like labels and addresses (will most of the time be filled by the runners.Runner)
type Host struct {
	Name      string
	Labels    map[string]string
	Addresses *IPAddresses
}

// IPAddresses list of IPv4 and IPv6 addresses a host has
type IPAddresses struct {
	IPv4 []string
	IPv6 []string
}

// Plan contains the information needed to execute the plan
type Plan struct {
	TestStartTime   time.Time         `json:"plannedTime"`
	AffectedServers map[string]*Host  `json:"affectedServers"`
	Commands        [][]*Task         `json:"commands"`
	Tester          string            `json:"tester"`
	RunOptions      config.RunOptions `json:"runOptions"`
}

// PrettyPrint "pretty" prints a plan
func (p Plan) PrettyPrint() {
	fmt.Println("-> BEGIN AffectedServers")
	for _, server := range p.AffectedServers {
		fmt.Println(server.Name)
	}
	fmt.Println("=> END AffectedServers")
	fmt.Println("-> BEGIN Commands")
	for k, commands := range p.Commands {
		round := k + 1
		fmt.Printf("--> BEGIN Round %d\n", round)
		for _, command := range commands {
			if command.Sleep != 0 {
				fmt.Printf("---> Wait for %+v\n", command.Sleep)
				continue
			}
			fmt.Printf("---> BEGIN Server %s\n", command.Host.Name)
			fmt.Printf("----> RUN %s %s (Additional info: %+v; %+v)\n", command.Command, command.Args, command.Ports, command.Sleep)
			for _, task := range command.SubTasks {
				fmt.Printf("-----> BEGIN Client %s\n", task.Host.Name)
				fmt.Printf("------> RUN %s %s (Additional info: %+v)\n", task.Command, task.Args, task.Ports)
				fmt.Printf("=====> END Client %s\n", task.Host.Name)
			}
			fmt.Printf("===> END Server %s\n", command.Host.Name)
		}
		fmt.Printf("==> END Round %d\n", round)
	}
	fmt.Println("=> END Commands")
}

// Task information for the task to execute
type Task struct {
	Host     *Host         `json:"host"`
	Command  string        `json:"command"`
	Args     []string      `json:"args"`
	Sleep    time.Duration `json:"sleep"`
	Ports    Ports         `json:"ports"`
	SubTasks []*Task       `json:"subTasks"`
	Status   *Status       `yaml:"status"`
}

// Ports TCP and UDP ports list
type Ports struct {
	TCP []int32
	UDP []int32
}

// Status status info for a task
type Status struct {
	SuccessfulHosts StatusHosts        `json:"successfulHosts"`
	FailedHosts     StatusHosts        `json:"failedHosts"`
	Errors          map[string][]error `json:"errors"`
}

// StatusHosts status per servers and clients list with counter
type StatusHosts struct {
	Servers map[string]int `json:"servers"`
	Clients map[string]int `json:"clients"`
}

// AddFailedServer add a server host that failed with error to the Status list
func (st *Status) AddFailedServer(host *Host, err error) {
	if _, ok := st.Errors[host.Name]; !ok {
		st.Errors[host.Name] = []error{}
	}
	st.Errors[host.Name] = append(st.Errors[host.Name], err)

	// Increase failed host counter
	if _, ok := st.FailedHosts.Servers[host.Name]; !ok {
		st.FailedHosts.Servers[host.Name] = 1
	} else {
		st.FailedHosts.Servers[host.Name]++
	}
}

// AddFailedClient add a client host that failed with error to the Status list
func (st *Status) AddFailedClient(host *Host, err error) {
	if _, ok := st.Errors[host.Name]; !ok {
		st.Errors[host.Name] = []error{}
	}
	st.Errors[host.Name] = append(st.Errors[host.Name], err)

	// Increase failed host counter
	if _, ok := st.FailedHosts.Clients[host.Name]; !ok {
		st.FailedHosts.Clients[host.Name] = 1
	} else {
		st.FailedHosts.Clients[host.Name]++
	}
}

// AddSuccessfulServer add a successful server host to the list
func (st *Status) AddSuccessfulServer(host *Host) {
	// Increase successful host counter
	if _, ok := st.SuccessfulHosts.Servers[host.Name]; !ok {
		st.SuccessfulHosts.Servers[host.Name] = 1
	} else {
		st.SuccessfulHosts.Servers[host.Name]++
	}
}

// AddSuccessfulClient add a successful client host to the list
func (st *Status) AddSuccessfulClient(host *Host) {
	// Increase successful host counter
	if _, ok := st.SuccessfulHosts.Clients[host.Name]; !ok {
		st.SuccessfulHosts.Clients[host.Name] = 1
	} else {
		st.SuccessfulHosts.Clients[host.Name]++
	}
}
