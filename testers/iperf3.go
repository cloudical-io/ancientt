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
	"github.com/cloudical-io/acntt/pkg/config"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// NameIPerf3 IPerf3 tester name
const NameIPerf3 = "iperf3"

func init() {
	Factories[NameIPerf3] = NewIPerf3Tester
}

// IPerf3 IPerf3 tester structure
type IPerf3 struct {
	Tester
	logger *log.Entry
	config *config.IPerf3
}

// NewIPerf3Tester return a new IPerf3 tester instance
func NewIPerf3Tester(cfg *config.Config, test *config.Test) (Tester, error) {
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
func (ip IPerf3) Plan(env *Environment, test *config.Test) (*Plan, error) {
	plan := &Plan{
		Tester:          test.Type,
		AffectedServers: map[string]*Host{},
		Commands:        make([][]*Task, test.RunOptions.Rounds),
	}

	var ports Ports
	if ip.config.UDP != nil && *ip.config.UDP {
		ports = Ports{
			UDP: []int32{5601},
		}
	} else {
		ports = Ports{
			TCP: []int32{5601},
		}
	}

	for i := 0; i < test.RunOptions.Rounds; i++ {
		for _, server := range env.Hosts.Servers {
			round := &Task{}
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
				round.SubTasks = append(round.SubTasks, &Task{
					Host:    client,
					Command: cmd,
					Args:    args,
					Ports:   ports,
					Status: Status{
						Errors:      map[string][]error{},
						FailedHosts: []string{},
					},
				})
			}
			plan.Commands[i] = append(plan.Commands[i], round)

			// Add the given interval after each round except the last one
			if test.RunOptions.Interval != 0 && i != test.RunOptions.Rounds-1 {
				plan.Commands[i] = append(plan.Commands[i], &Task{
					Sleep: test.RunOptions.Interval,
				})
			}
		}
	}

	return plan, nil
}

// buildIPerf3ServerCommand generate IPer3 server command
func (ip IPerf3) buildIPerf3ServerCommand(server *Host) (string, []string) {
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
func (ip IPerf3) buildIPerf3ClientCommand(server *Host, client *Host) (string, []string) {
	// Base command and args
	cmd := "iperf3"
	args := []string{
		"--json",
		"--port={{ .ServerPort }}",
		"--client={{ .ServerAddress }}",
	}

	// Add --udp flag when UDP should be used
	if ip.config.UDP != nil && *ip.config.UDP {
		args = append(args, "--udp")
	}

	// Append additional client flags to args array
	args = append(args, ip.config.AdditionalFlags.Clients...)

	return cmd, args
}
