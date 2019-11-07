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

package iperf3

import (
	"fmt"

	"github.com/cloudical-io/ancientt/pkg/config"
	"github.com/cloudical-io/ancientt/testers"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// NameIPerf3 IPerf3 tester name
const NameIPerf3 = "iperf3"

func init() {
	testers.Factories[NameIPerf3] = NewIPerf3Tester
}

// IPerf3 IPerf3 tester structure
type IPerf3 struct {
	testers.Tester
	logger *log.Entry
	config *config.IPerf3
}

// NewIPerf3Tester return a new IPerf3 tester instance
func NewIPerf3Tester(cfg *config.Config, test *config.Test) (testers.Tester, error) {
	if test == nil {
		test = &config.Test{
			IPerf3: &config.IPerf3{},
		}
	}

	return IPerf3{
		logger: log.WithFields(logrus.Fields{"tester": NameIPerf3}),
		config: test.IPerf3,
	}, nil
}

// Plan return a plan to run IPerf3 from the given config.Test and Environment information (hosts)
func (ip IPerf3) Plan(env *testers.Environment, test *config.Test) (*testers.Plan, error) {
	plan := &testers.Plan{
		Tester:          test.Type,
		AffectedServers: map[string]*testers.Host{},
		Commands:        make([][]*testers.Task, test.RunOptions.Rounds),
	}

	var ports testers.Ports
	if ip.config.UDP != nil && *ip.config.UDP {
		ports = testers.Ports{
			UDP: []int32{5601},
		}
	} else {
		ports = testers.Ports{
			TCP: []int32{5601},
		}
	}

	for i := 0; i < test.RunOptions.Rounds; i++ {
		for _, server := range env.Hosts.Servers {
			round := &testers.Task{
				Status: &testers.Status{
					SuccessfulHosts: testers.StatusHosts{
						Servers: map[string]int{},
						Clients: map[string]int{},
					},
					FailedHosts: testers.StatusHosts{
						Servers: map[string]int{},
						Clients: map[string]int{},
					},
					Errors: map[string][]error{},
				},
			}
			// Add server host to AffectedServers list
			if _, ok := plan.AffectedServers[server.Name]; !ok {
				plan.AffectedServers[server.Name] = server
			}

			// Set the server that will run the iperf3 server in the "main" command
			round.Host = server
			round.Command, round.Args = ip.buildIPerf3ServerCommand(server)
			round.Ports = ports

			// Now go over each client and generate their Task
			for _, client := range env.Hosts.Clients {
				// Add client host to AffectedServers list
				if _, ok := plan.AffectedServers[client.Name]; !ok {
					plan.AffectedServers[client.Name] = client
				}

				// Build the IPerf3 command
				cmd, args := ip.buildIPerf3ClientCommand(server, client)
				round.SubTasks = append(round.SubTasks, &testers.Task{
					Host:    client,
					Command: cmd,
					Args:    args,
					Ports:   ports,
				})
			}
			plan.Commands[i] = append(plan.Commands[i], round)

			// Add the given interval after each round except the last one
			if test.RunOptions.Interval != 0 && i != test.RunOptions.Rounds-1 {
				plan.Commands[i] = append(plan.Commands[i], &testers.Task{
					Sleep: test.RunOptions.Interval,
				})
			}
		}
	}

	return plan, nil
}

// buildIPerf3ServerCommand generate IPer3 server command
func (ip IPerf3) buildIPerf3ServerCommand(server *testers.Host) (string, []string) {
	// Base command and args
	cmd := "iperf3"
	args := []string{
		"--json",
		"--port={{ .ServerPort }}",
		"--server",
	}

	// Add --udp flag when UDP should be used
	if ip.config.UDP != nil && *ip.config.UDP {
		args = append(args, "--udp")
	}

	// Append additional server flags to args array
	args = append(args, ip.config.AdditionalFlags.Server...)

	return cmd, args
}

// buildIPerf3ClientCommand generate IPer3 client command
func (ip IPerf3) buildIPerf3ClientCommand(server *testers.Host, client *testers.Host) (string, []string) {
	// Base command and args
	cmd := "iperf3"
	args := []string{
		fmt.Sprintf("--time=%d", ip.config.Duration),
		fmt.Sprintf("--interval=%d", ip.config.Interval),
		"--json",
		"--port={{ .ServerPort }}",
		"--client={{ .ServerAddressV4 }}",
	}

	// Add --udp flag when UDP should be used
	if ip.config.UDP != nil && *ip.config.UDP {
		args = append(args, "--udp")
	}

	// Append additional client flags to args array
	args = append(args, ip.config.AdditionalFlags.Clients...)

	return cmd, args
}
