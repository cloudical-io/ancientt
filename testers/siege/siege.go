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

package siege

import (
	"github.com/cloudical-io/ancientt/pkg/config"
	"github.com/cloudical-io/ancientt/testers"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// NameSiege Siege tester name
const NameSiege = "siege"

func init() {
	testers.Factories[NameSiege] = NewSiegeTester
}

// Siege Siege tester structure
type Siege struct {
	testers.Tester
	logger *log.Entry
	config *config.Siege
}

// NewSiegeTester return a new Siege tester instance
func NewSiegeTester(cfg *config.Config, test *config.Test) (testers.Tester, error) {
	if test == nil {
		test = &config.Test{
			Siege: &config.Siege{},
		}
	}

	return Siege{
		logger: log.WithFields(logrus.Fields{"tester": NameSiege}),
		config: test.Siege,
	}, nil
}

// Plan return a plan to run Siege from the given config.Test and Environment information (hosts)
func (sg Siege) Plan(env *testers.Environment, test *config.Test) (*testers.Plan, error) {
	plan := &testers.Plan{
		Tester:          test.Type,
		AffectedServers: map[string]*testers.Host{},
		Commands:        make([][]*testers.Task, test.RunOptions.Rounds),
	}

	var ports testers.Ports
	ports = testers.Ports{
		TCP: []int32{80},
	}

	for i := 0; i < test.RunOptions.Rounds; i++ {
		for _, server := range env.Hosts.Servers {
			round := &testers.Task{}
			// Add server host to AffectedServers list
			if _, ok := plan.AffectedServers[server.Name]; !ok {
				plan.AffectedServers[server.Name] = server
			}

			// Set the server that will run the siege server in the "main" command
			round.Host = server
			round.Command, round.Args = sg.buildSiegeServerCommand(server)
			round.Ports = ports

			// Now go over each client and generate their Task
			for _, client := range env.Hosts.Clients {
				// Add client host to AffectedServers list
				if _, ok := plan.AffectedServers[client.Name]; !ok {
					plan.AffectedServers[client.Name] = client
				}

				// Build the Siege command
				cmd, args := sg.buildSiegeClientCommand(server, client)
				round.SubTasks = append(round.SubTasks, &testers.Task{
					Host:    client,
					Command: cmd,
					Args:    args,
					Ports:   ports,
				})
			}
			plan.Commands[i] = append(plan.Commands[i], round)
		}
	}

	return plan, nil
}

// buildSiegeServerCommand generate IPer3 server command
func (sg Siege) buildSiegeServerCommand(server *testers.Host) (string, []string) {
	// Base command and args
	cmd := "nginx"
	args := []string{
		"-g",
		"daemon off;",
	}

	// Append additional server flags to args array
	args = append(args, sg.config.AdditionalFlags.Server...)

	return cmd, args
}

// buildSiegeClientCommand generate IPer3 client command
func (sg Siege) buildSiegeClientCommand(server *testers.Host, client *testers.Host) (string, []string) {
	// Base command and args
	cmd := "siege"
	args := []string{}

	// Append additional client flags to args array
	args = append(args, sg.config.AdditionalFlags.Clients...)

	return cmd, args
}
