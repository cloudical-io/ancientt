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

	"github.com/cloudical-io/acntt/pkg/config"
)

const NameIPerf3 = "iperf3"

func init() {
	Factories[NameIPerf3] = NewIPerf3Tester
}

// IPerf3
type IPerf3 struct {
	Tester
}

func NewIPerf3Tester() (Tester, error) {
	return IPerf3{}, nil
}

// Plan
func (ip IPerf3) Plan(env *Environment, test *config.Test) (*Plan, error) {
	plan := &Plan{
		Tester:          test.Type,
		AffectedServers: map[string]Host{},
		Commands:        []map[string][]Task{},
	}

	for i := 0; i < test.RunOptions.Rounds; i++ {
		round := map[string][]Task{}
		for _, server := range env.Hosts.Servers {
			// Add server host to AffectedServers list
			if _, ok := plan.AffectedServers[server.Name]; !ok {
				plan.AffectedServers[server.Name] = server
			}
			round[server.Name] = []Task{}
			for _, client := range env.Hosts.Clients {
				// Add client host to AffectedServers list
				if _, ok := plan.AffectedServers[client.Name]; !ok {
					plan.AffectedServers[client.Name] = client
				}

				command, err := buildIPerf3Command(server, client)
				if err != nil {
					return plan, err
				}
				round[server.Name] = append(round[server.Name], Task{
					Command: command,
				})
			}
		}
		plan.Commands = append(plan.Commands, round)
	}

	return plan, nil
}

func buildIPerf3Command(server Host, client Host) (string, error) {
	return fmt.Sprintf("iperf3 -c %s -p %d %s", client.Name, 5601, ""), nil
}
