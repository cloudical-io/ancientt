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

package runners

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/cloudical-io/acntt/parsers"
	"github.com/cloudical-io/acntt/pkg/config"
	"github.com/cloudical-io/acntt/testers"
)

const (
	// NameMock Mock Runner Name
	NameMock              = "mock"
	mockServerNamePattern = "servers-%d"
)

func init() {
	Factories[NameMock] = NewMockRunner
}

// Mock Mock Runner struct
type Mock struct {
	Runner
}

// NewMockRunner returns a new Mock Runner
func NewMockRunner(cfg *config.Config) (Runner, error) {
	return Mock{}, nil
}

// GetHostsForTest return a mocked list of hots for the given test config
func (k Mock) GetHostsForTest(test config.Test) (*testers.Hosts, error) {
	// Pre create the structure to return
	hosts := &testers.Hosts{
		Clients: map[string]*testers.Host{},
		Servers: map[string]*testers.Host{},
	}

	mockHosts := generateMockServers()

	// Create and seed randomness source for the `random` selection of hosts
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	r.Seed(time.Now().UnixNano())

	for _, clients := range test.Hosts.Clients {
		if clients.All {
			for _, mockHost := range mockHosts {
				hosts.Clients[mockHost.Name] = mockHost
			}
		}
		if clients.Random {
			for i := 0; i < clients.Count; i++ {
				mockHost := mockHosts[r.Intn(len(mockHosts))]
				hosts.Clients[mockHost.Name] = mockHost
			}
		}
		if len(clients.Hosts) > 0 {
			// Just mock any hosts which are in a list format directly given
			for k, mockHost := range clients.Hosts {
				hosts.Clients[mockHost] = &testers.Host{
					Name: mockHost,
					Addresses: &testers.IPAddresses{
						IPv4: []string{
							fmt.Sprintf("%d.%d.%d.%d", k, k, k, k),
						},
						IPv6: []string{
							fmt.Sprintf("2001:db8:abcd:0012::%d", k),
						},
					},
				}
			}
		}
	}

	for _, servers := range test.Hosts.Servers {
		if servers.All {
			for _, mockHost := range mockHosts {
				hosts.Servers[mockHost.Name] = mockHost
			}
			continue
		}
		if servers.Random {
			for i := 0; i < servers.Count; i++ {
				mockHost := mockHosts[r.Intn(len(mockHosts))]
				hosts.Servers[mockHost.Name] = mockHost
			}
			// TODO Filter on labelSelector basis and antiAffinity
			continue
		}
		if len(servers.Hosts) > 0 {
			// Just mock any hosts which are in a list format directly given
			for _, mockHost := range servers.Hosts {
				hosts.Servers[mockHost] = &testers.Host{
					Name: mockHost,
					Addresses: &testers.IPAddresses{
						IPv4: []string{
							"1.1.1.1",
						},
						IPv6: []string{
							"2001:db8:abcd:0012::1",
						},
					},
				}
			}
			continue
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
func (k Mock) Prepare(runOpts config.RunOptions, plan *testers.Plan) error {
	return nil
}

// Execute run the given testers.Plan and return the logs of each step and / or error
func (k Mock) Execute(plan *testers.Plan, parser parsers.Parser) error {
	// Return nothing because we didn't do anything in the mock
	return nil
}

// Cleanup NOOP because Mock doesn't create any resource nor connection or so to any hosts.
func (k Mock) Cleanup(plan *testers.Plan) error {
	// TODO
	return nil
}
