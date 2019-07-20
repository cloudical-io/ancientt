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
	"time"
)

func generateMockTableData(length int) Data {
	table := Table{
		Headers: []Column{
			Column{
				Rows: []Row{
					Row{Value: "isthisfloat64"},
					Row{Value: "isthisinteger64"},
					Row{Value: "isittrue"},
					Row{Value: "data"},
				},
			},
		},
		Columns: []Column{},
	}
	column := Column{
		Rows: []Row{
			Row{Value: float64(123.456789)},
			Row{Value: int64(35671233)},
			Row{Value: true},
			Row{Value: "data"},
		},
	}
	for i := 0; i < length; i++ {
		table.Columns = append(table.Columns, column)
	}

	return Data{
		AdditionalInfo: "mock data generated",
		ClientHost:     "host1",
		ServerHost:     "host2",
		TestStartTime:  time.Now(),
		TestTime:       time.Now(),
		Tester:         "foobar",
		Data:           table,
	}
}
