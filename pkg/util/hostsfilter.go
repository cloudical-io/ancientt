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

package util

import (
	"math/rand"
	"reflect"
	"time"

	"github.com/cloudical-io/acntt/pkg/config"
	"github.com/cloudical-io/acntt/testers"
)

// FilterHostsList filter a given host list
func FilterHostsList(inHosts []*testers.Host, filter config.Hosts) ([]*testers.Host, error) {
	// Create and seed randomness source for the `random` selection of hosts
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	r.Seed(time.Now().UnixNano())

	if filter.All {
		return inHosts, nil
	}

	hosts := []*testers.Host{}

	if len(filter.Hosts) > 0 {
		for _, host := range filter.Hosts {
			hosts = append(hosts, &testers.Host{
				Name: host,
			})
		}
		return hosts, nil
	}

	filteredHosts := filterHostsByLabels(inHosts, filter.HostSelector)

	// Get random server(s)
	if filter.Random {
		for i := 0; i < filter.Count; i++ {
			inHost := filteredHosts[r.Intn(len(filteredHosts))]
			hosts = append(hosts, inHost)
		}
		return hosts, nil
	}

	return hosts, nil
}

// filterHostsByLabels
func filterHostsByLabels(hosts []*testers.Host, labels map[string]string) []*testers.Host {
	if len(labels) == 0 {
		return hosts
	}
	filtered := []*testers.Host{}
	for _, host := range hosts {
		// Compare host and filter labels list
		if reflect.DeepEqual(host.Labels, labels) {
			filtered = append(filtered, host)
			continue
		}
		// TODO implement anti affinity logic based on labels
	}

	return filtered
}
