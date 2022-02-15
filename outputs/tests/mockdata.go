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
	"math/rand"
	"time"

	"github.com/cloudical-io/ancientt/outputs"
)

// GenerateMockTableData generate some mock DataTable data for testing purposes
func GenerateMockTableData(length int) outputs.Data {
	table := &outputs.Table{
		Headers: []*outputs.Row{
			{Value: "isthisfloat64"},
			{Value: "iamafloat64part2"},
			{Value: "isthisinteger64"},
			{Value: "isittrue"},
			{Value: "data"},
			{Value: "interval"},
		},
		Rows: [][]*outputs.Row{},
	}

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	for i := 0; i < length; i++ {
		f := float64(i)
		r := []*outputs.Row{
			{Value: (r.Float64() * f) + f},
			{Value: (r.Float64() * f) + f},
			{Value: int64(r.Intn(99999))},
			{Value: true},
			{Value: "data"},
			{Value: i},
		}
		table.Rows = append(table.Rows, r)
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
