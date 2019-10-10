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

package mock

import (
	"fmt"

	"github.com/cloudical-io/ancientt/parsers"
	"github.com/cloudical-io/ancientt/pkg/config"
	"github.com/cloudical-io/ancientt/pkg/hostsfilter"
	"github.com/cloudical-io/ancientt/runners"
	"github.com/cloudical-io/ancientt/testers"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

const (
	// Name Mock Runner Name
	Name                  = "mock"
	mockServerNamePattern = "servers-%d"
)

func init() {
	runners.Factories[Name] = NewRunner
}

// Mock Mock Runner struct
type Mock struct {
	runners.Runner
	logger *log.Entry
}

// NewRunner returns a new Mock Runner
func NewRunner(cfg *config.Config) (runners.Runner, error) {
	return Mock{
		logger: log.WithFields(logrus.Fields{"runner": Name}),
	}, nil
}

// GetHostsForTest return a mocked list of hots for the given test config
func (m Mock) GetHostsForTest(test *config.Test) (*testers.Hosts, error) {
	// Pre create the structure to return
	hosts := &testers.Hosts{
		Clients: map[string]*testers.Host{},
		Servers: map[string]*testers.Host{},
	}

	mockHosts := generateMockServers()

	// Go through Hosts Servers list to get the servers hosts
	for _, servers := range test.Hosts.Servers {
		filtered, err := hostsfilter.FilterHostsList(mockHosts, servers)
		if err != nil {
			return nil, err
		}
		for _, host := range filtered {
			if _, ok := hosts.Servers[host.Name]; !ok {
				hosts.Servers[host.Name] = host
			}
		}
	}

	// Go through Hosts Clients list to get the clients hosts
	for _, clients := range test.Hosts.Clients {
		filtered, err := hostsfilter.FilterHostsList(mockHosts, clients)
		if err != nil {
			return nil, err
		}
		for _, host := range filtered {
			if _, ok := hosts.Clients[host.Name]; !ok {
				hosts.Clients[host.Name] = host
			}
		}
	}

	return hosts, nil
}

// generateMockServers generate a list of mcoekd servers for testing purposes
func generateMockServers() []*testers.Host {
	hosts := []*testers.Host{}
	for i := 0; i < 10; i++ {
		hosts = append(hosts, &testers.Host{
			Name:      fmt.Sprintf(mockServerNamePattern, i),
			Addresses: &testers.IPAddresses{},
			Labels: map[string]string{
				"i-am-server": fmt.Sprintf(mockServerNamePattern, i),
			},
		})
	}
	return hosts
}

// Prepare NOOP because there is nothing to prepare because this is Mock.
func (m Mock) Prepare(runOpts config.RunOptions, plan *testers.Plan) error {
	m.logger.Info("Mock.Prepare() called")
	return nil
}

// Execute run the given testers.Plan and return the logs of each step and / or error
func (m Mock) Execute(plan *testers.Plan, parser chan<- parsers.Input) error {
	m.logger.Info("Mock.Execute() called")
	// Return nothing because we don't do anything in the Mock
	return nil
}

// Cleanup NOOP because Mock doesn't create any resource nor connection or so to any hosts.
func (m Mock) Cleanup(plan *testers.Plan) error {
	m.logger.Info("Mock.Cleanup() called")
	// Return nothing because we don't do anything in the Mock
	return nil
}
