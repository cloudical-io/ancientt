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

package hostsfilter

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/cloudical-io/ancientt/testers"
)

func TestAntiAffinity(t *testing.T) {
	hosts := []*testers.Host{
		{
			Name: "host-a",
			Labels: map[string]string{
				"foo": "bar",
			},
		},
		{
			Name: "host-b",
			Labels: map[string]string{
				"foo": "bar",
			},
		},
		{
			Name: "host-c",
			Labels: map[string]string{
				"foo": "notbar",
			},
		},
		{
			Name: "host-d",
			Labels: map[string]string{
				"bar": "foo",
			},
		},
	}
	antiAffinityFilter := []string{
		"foo",
	}

	filteredHosts := checkAntiAffinity(hosts, antiAffinityFilter)
	assert.Equal(t, 3, len(filteredHosts))

	// host-b will not be in the list because `host-a` has the same anti affinity label as `host-a`
	assert.Equal(t, "host-a", filteredHosts[0].Name)
	assert.Equal(t, "host-c", filteredHosts[1].Name)
	assert.Equal(t, "host-d", filteredHosts[2].Name)
}
