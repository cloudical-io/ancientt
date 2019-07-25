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

package smokeping

import (
	"github.com/cloudical-io/acntt/testers"
	"github.com/cloudical-io/acntt/pkg/config"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// NameSmokeping Smokeping tester name
const NameSmokeping = "smokeping"

func init() {
	testers.Factories[NameSmokeping] = NewSmokepingTester
}

// Smokeping Smokeping tester structure
type Smokeping struct {
	testers.Tester
	logger *log.Entry
	config *config.Smokeping
}

// NewSmokepingTester return a new Smokeping tester instance
func NewSmokepingTester(cfg *config.Config, test *config.Test) (testers.Tester, error) {
	if test == nil {
		test = &config.Test{
			Smokeping: &config.Smokeping{},
		}
	}

	return Smokeping{
		logger: log.WithFields(logrus.Fields{"tester": NameSmokeping}),
		config: test.Smokeping,
	}, nil
}

// Plan return a plan to run Smokeping from the given config.Test and Environment information (hosts)
func (sp Smokeping) Plan(env *testers.Environment, test *config.Test) (*testers.Plan, error) {
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

			// Set the server that will run the smokeping server in the "main" command
			round.Host = server
			round.Command, round.Args = sp.buildSmokepingServerCommand(server)
			round.Ports = ports

			// Now go over each client and generate their Task
			for _, client := range env.Hosts.Clients {
				// Add client host to AffectedServers list
				if _, ok := plan.AffectedServers[client.Name]; !ok {
					plan.AffectedServers[client.Name] = client
				}

				// Build the Smokeping command
				cmd, args := sp.buildSmokepingClientCommand(server, client)
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

// buildSmokepingServerCommand generate IPer3 server command
func (sp Smokeping) buildSmokepingServerCommand(server *testers.Host) (string, []string) {
	// Base command and args
	cmd := "smokeping"
	args := []string{}

	// TODO Set flags depending on config

	// Append additional server flags to args array
	args = append(args, sp.config.AdditionalFlags.Server...)

	return cmd, args
}

// buildSmokepingClientCommand generate IPer3 client command
func (sp Smokeping) buildSmokepingClientCommand(server *testers.Host, client *testers.Host) (string, []string) {
	// Base command and args
	cmd := "smokeping"
	args := []string{}

	// Append additional client flags to args array
	args = append(args, sp.config.AdditionalFlags.Clients...)

	return cmd, args
}
