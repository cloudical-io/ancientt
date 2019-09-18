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
	"time"

	"github.com/cloudical-io/acntt/pkg/config"
	"github.com/cloudical-io/acntt/testers"
)

// FilterHostsList filter a given host list
func FilterHostsList(inHosts []*testers.Host, filter config.Hosts) ([]*testers.Host, error) {
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

	filteredHosts = checkAntiAffinity(filteredHosts, filter.AntiAffinity)

	if len(filteredHosts) == 0 {
		return filteredHosts, nil
	}

	if filter.All {
		return filteredHosts, nil
	}

	// Create and seed randomness source for the `random` selection of hosts
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	r.Seed(time.Now().UnixNano())

	// Get random server(s)
	if filter.Random {
		for i := 0; i < filter.Count; i++ {
			inHost := filteredHosts[r.Intn(len(filteredHosts))]
			hosts = append(hosts, inHost)
		}
		return hosts, nil
	}

	return filteredHosts, nil
}

// filterHostsByLabels all labels must match
func filterHostsByLabels(hosts []*testers.Host, labels map[string]string) []*testers.Host {
	if len(labels) == 0 {
		return hosts
	}

	filtered := []*testers.Host{}
	for _, host := range hosts {
		// Compare host and filter labels list, all labels list must match
		match := true
		for k, v := range labels {
			if labelValue, ok := host.Labels[k]; ok {
				if labelValue != v {
					match = false
					break
				}
			} else {
				match = false
			}
		}
		if match {
			filtered = append(filtered, host)
		}
	}

	return filtered
}

func checkAntiAffinity(hosts []*testers.Host, labels []string) []*testers.Host {
	if len(labels) == 0 {
		return hosts
	}

	filtered := []*testers.Host{}
	usedLabels := map[string][]string{}

	for _, host := range hosts {
		// Compare host and filter labels list, all labels list must match
		match := true
		for _, label := range labels {
			hostLabelVal, ok := host.Labels[label]
			if !ok {
				continue
			}
			// Check if label and value is in usedLabels list and if it is the host should not be added
			if usedLabelValues, ok := usedLabels[label]; ok {
				for _, usedLabelVal := range usedLabelValues {
					if usedLabelVal == hostLabelVal {
						match = false
						break
					}
				}
				usedLabels[label] = append(usedLabels[label], hostLabelVal)
			} else {
				usedLabels[label] = []string{hostLabelVal}
			}
		}
		if match {
			filtered = append(filtered, host)
		}
	}
	return filtered
}
