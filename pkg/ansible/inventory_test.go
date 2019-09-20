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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetHostsForTest(t *testing.T) {
	inv, err := Parse([]byte(`{
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
		},
		"test123": {
			"children": [
				"all"
			],
			"hosts": [
				"server8",
				"server9"
			]
		}
	}`))
	require.NotNil(t, inv)
	require.Nil(t, err)

	all := inv.GetHostsForGroup("all")
	assert.Equal(t, 3, len(all))
	assert.Contains(t, all, "server1")
	assert.Contains(t, all, "server2")
	assert.Contains(t, all, "server4")

	clients := inv.GetHostsForGroup("clients")
	assert.Equal(t, 2, len(clients))
	assert.Contains(t, clients, "server1")
	assert.Contains(t, clients, "server2")
	assert.NotContains(t, clients, "server4")

	server := inv.GetHostsForGroup("server")
	assert.Equal(t, 1, len(server))
	assert.NotContains(t, server, "server1")
	assert.NotContains(t, server, "server2")
	assert.Contains(t, server, "server4")

	meta := inv.GetHostsForGroup("_meta")
	assert.Equal(t, 0, len(meta))

	test123 := inv.GetHostsForGroup("test123")
	assert.Equal(t, 5, len(test123))
	assert.Contains(t, test123, "server1")
	assert.Contains(t, test123, "server2")
	assert.Contains(t, test123, "server4")
	assert.Contains(t, test123, "server8")
	assert.Contains(t, test123, "server9")
}
