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

package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TODO Add a simple but useful test to verify config loading and basic functionality

func TestAncienttRunCommand(t *testing.T) {
	err := run(nil, []string{})
	// Not nill because there is no config file named `testdefinition.yaml` where the current working directory is
	assert.NotNil(t, err)

	// TODO Generate temp dir and create basic `testdefinition.yaml` in it, switch cwd to tmp dir, no error should be returned then
	/*
	   	tmpFile, err := ioutil.TempFile(os.TempDir(), "ancienttgotests")
	   	require.Nil(t, err)
	       // defer delete tmpfile

	   	err = run(nil, []string{})
	   	// Not nill because there is no config file named `testdefinition.yaml` where the current working directory is
	   	assert.Nil(t, err)
	*/
}
