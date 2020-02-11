/*
Copyright 2020 Cloudical Deutschland GmbH. All rights reserved.
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

package pingparsing

import (
	"fmt"

	"github.com/cloudical-io/ancientt/pkg/config"
	"github.com/cloudical-io/ancientt/testers"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// NamePingParsing PingParsing tester name
const NamePingParsing = "pingparsing"

func init() {
	testers.Factories[NamePingParsing] = NewPingParsingTester
}

// PingParsing PingParsing tester structure
type PingParsing struct {
	testers.Tester
	logger *log.Entry
	config *config.PingParsing
}

// NewPingParsingTester return a new PingParsing tester instance
func NewPingParsingTester(cfg *config.Config, test *config.Test) (testers.Tester, error) {
	if test == nil {
		test = &config.Test{
			PingParsing: &config.PingParsing{},
		}
	}

	return PingParsing{
		logger: log.WithFields(logrus.Fields{"tester": NamePingParsing}),
		config: test.PingParsing,
	}, nil
}

// Plan
func (t PingParsing) Plan(env *testers.Environment, test *config.Test) (*testers.Plan, error) {
	plan := &testers.Plan{
		Tester:          test.Type,
		AffectedServers: map[string]*testers.Host{},
		Commands:        make([][]*testers.Task, test.RunOptions.Rounds),
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

			// Setting a server is not really needed, but we set it to `sleep 99999`
			// for compatibility with Runners such as Kubernetes where the IP is only
			// available when a Server Pod is running
			round.Host = server
			round.Command, round.Args = t.buildPingParsingServerCommand(server)

			// Now go over each client and generate their Task
			for _, client := range env.Hosts.Clients {
				// Add client host to AffectedServers list
				if _, ok := plan.AffectedServers[client.Name]; !ok {
					plan.AffectedServers[client.Name] = client
				}

				// Build the PingParsing command
				cmd, args := t.buildPingParsingClientCommand(server, client)
				round.SubTasks = append(round.SubTasks, &testers.Task{
					Host:    client,
					Command: cmd,
					Args:    args,
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

// buildPingParsingServerCommand
func (t PingParsing) buildPingParsingServerCommand(server *testers.Host) (string, []string) {
	return "sleep", []string{"9999999"}
}

// buildPingParsingClientCommand
func (t PingParsing) buildPingParsingClientCommand(server *testers.Host, client *testers.Host) (string, []string) {
	// Base command and args
	cmd := "pingparsing"
	args := []string{
		"--icmp-reply",
		"--timestamp=datetime",
		fmt.Sprintf("-c=%d", *t.config.Count),
		fmt.Sprintf("-w=%s", *t.config.Deadline),
		fmt.Sprintf("--timeout=%s", *t.config.Timeout),
		fmt.Sprintf("-I=%s", t.config.Interface),
		"{{ .ServerAddressV4 }}",
	}

	return cmd, args
}
