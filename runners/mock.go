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
	hosts := &testers.Hosts{
		Clients: map[string]testers.Host{},
		Servers: map[string]testers.Host{},
	}

	mockHosts := generateMockServers()

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
	}

	for _, servers := range test.Hosts.Servers {
		if servers.All {
			for _, mockHost := range mockHosts {
				hosts.Servers[mockHost.Name] = mockHost
			}
		}
		if servers.Random {
			for i := 0; i < servers.Count; i++ {
				mockHost := mockHosts[r.Intn(len(mockHosts))]
				hosts.Servers[mockHost.Name] = mockHost
			}
		}
	}

	return hosts, nil
}

func generateMockServers() []testers.Host {
	hosts := []testers.Host{}
	for i := 0; i < 10; i++ {
		hosts = append(hosts, testers.Host{
			Name:      fmt.Sprintf(mockServerNamePattern, i),
			Addresses: testers.IPAddresses{},
			Labels: map[string]string{
				"i-am-server": fmt.Sprintf(mockServerNamePattern, i),
			},
		})
	}
	return hosts
}

// Execute run the given commands and return the logs of it and / or error
func (k Mock) Execute(cmd, args []string) ([]byte, error) {
	return []byte{}, nil
}
