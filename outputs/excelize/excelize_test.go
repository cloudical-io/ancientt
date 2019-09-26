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

package excelize

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/cloudical-io/ancientt/outputs/tests"
	"github.com/cloudical-io/ancientt/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDo(t *testing.T) {
	table := tests.GenerateMockTableData(2)

	tempDir := os.TempDir()
	outName := fmt.Sprintf("ancientt-test-%s.xlsx", t.Name())
	tmpOutFile := path.Join(tempDir, outName)
	//defer os.Remove(tmpOutFile)

	outCfg := &config.Output{
		Excelize: &config.Excelize{
			FilePath: config.FilePath{
				FilePath:    tempDir,
				NamePattern: outName,
			},
		},
	}

	e, err := NewExcelizeOutput(nil, outCfg)
	assert.Nil(t, err)
	err = e.Do(table)
	assert.Nil(t, err)

	fInfo, err := os.Stat(tmpOutFile)
	require.Nil(t, err)
	require.NotNil(t, fInfo)

	// For now just check that the file size is bigger than 0
	assert.True(t, fInfo.Size() > 0)
	// TODO Add more sophisticated validation of file content
}
