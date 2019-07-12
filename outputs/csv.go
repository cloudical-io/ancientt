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

package outputs

import (
	"fmt"

	"github.com/cloudical-io/acntt/pkg/config"
)

// NameCSV CSV output name
const NameCSV = "csv"

func init() {
	Factories[NameCSV] = NewCSVOutput
}

// CSV CSV tester structure
type CSV struct {
	Output
	config *config.CSV
}

// NewCSVOutput return a new CSV tester instance
func NewCSVOutput(cfg *config.Config, outCfg *config.Output) (Output, error) {
	return CSV{
		config: outCfg.CSV,
	}, nil
}

// Do parse CSV JSON responses
func (ip CSV) Do(data Data) error {
	dataTable, ok := data.Data.(Table)
	if !ok {
		return fmt.Errorf("data not in table for csv output")
	}
	for _, column := range dataTable.Headers {
		for _, row := range column.Rows {
			fmt.Print(row.Value, ",")
		}
		fmt.Print("\n")
	}
	for _, column := range dataTable.Columns {
		for _, row := range column.Rows {
			fmt.Print(row.Value, ",")
		}
		fmt.Print("\n")
	}
	return nil
}
