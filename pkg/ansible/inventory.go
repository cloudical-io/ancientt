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

package ansible

import (
	"encoding/json"
)

/*
// The _meta is ignored
{
    "_meta": {
        "hostvars": {}
    },
    "all": {
        "children": [
            "clients",
            "server",
            "ungrouped"
        ]
    },
    "clients": {
        "hosts": [
            "server1",
            "server2"
        ]
    },
    "server": {
        "hosts": [
            "server4"
        ]
    }
}
*/

// InventoryList Basic `ansible-inventory` JSON output structure
type InventoryList map[string]HostGroup

// HostGroup Host and groups of a group `ansible-inventory` JSON sub-structure
type HostGroup struct {
	Children []string `json:"children"`
	Hosts    []string `json:"hosts"`
}

// Parse raw JSON into
func Parse(in []byte) (*InventoryList, error) {
	inv := &InventoryList{}

	if err := json.Unmarshal(in, inv); err != nil {
		return nil, err
	}

	return inv, nil
}

// GetHostsForGroup return resolved list of hosts for a given group name
func (inv InventoryList) GetHostsForGroup(group string) []string {
	hosts := []string{}

	for k, hg := range inv {
		if k == group {
			if len(hg.Children) > 0 {
				for _, child := range hg.Children {
					hosts = append(hosts, inv.GetHostsForGroup(child)...)
				}
			}
			if len(hg.Hosts) > 0 {
				hosts = append(hosts, hg.Hosts...)
			}
		}
	}

	return hosts
}
