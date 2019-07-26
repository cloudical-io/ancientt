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

package tests

import (
	"github.com/cloudical-io/acntt/outputs"
	"time"
)

func GenerateMockTableData(length int) outputs.Data {
	table := outputs.Table{
		Headers: []outputs.Column{
			outputs.Column{
				Rows: []outputs.Row{
					outputs.Row{Value: "isthisfloat64"},
					outputs.Row{Value: "isthisinteger64"},
					outputs.Row{Value: "isittrue"},
					outputs.Row{Value: "data"},
				},
			},
		},
		Columns: []outputs.Column{},
	}
	column := outputs.Column{
		Rows: []outputs.Row{
			outputs.Row{Value: float64(123.456789)},
			outputs.Row{Value: int64(35671233)},
			outputs.Row{Value: true},
			outputs.Row{Value: "data"},
		},
	}
	for i := 0; i < length; i++ {
		table.Columns = append(table.Columns, column)
	}

	return outputs.Data{
		AdditionalInfo: "mock data generated",
		ClientHost:     "host1",
		ServerHost:     "host2",
		TestStartTime:  time.Now(),
		TestTime:       time.Now(),
		Tester:         "foobar",
		Data:           table,
	}
}
