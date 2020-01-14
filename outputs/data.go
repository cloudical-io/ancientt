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
	"time"

	"github.com/cloudical-io/ancientt/pkg/config"
)

// Data structured parsed data
type Data struct {
	TestStartTime  time.Time
	TestTime       time.Time
	Tester         string
	ServerHost     string
	ClientHost     string
	AdditionalInfo string
	Data           DataFormat
}

// DataFormat DataFormat interface that must be implemented by data formats, e.g., Table.
type DataFormat interface {
	// Transform run transformations on the `Data`.
	Transform(ts []*config.Transformation) error
}

// Table Data format for data in Table form
type Table struct {
	DataFormat
	Headers []Row
	Rows    [][]Row
}

// Row Row of the Table data format
type Row struct {
	Value interface{}
}

// Transform transformation of table data
func (d Table) Transform(ts []*config.Transformation) error {
	// Iterate over each transformation
	for _, t := range ts {
		i, err := d.GetHeaderIndexByName(t.Key)
		if err != nil {
			return err
		}
		if i == -1 && t.Action != config.TransformationActionAdd {
			return nil
		}

		switch t.Action {
		case config.TransformationActionAdd:
			d.Headers = append(d.Headers, Row{
				Value: t.Key,
			})
			i = len(d.Headers) - 1
		case config.TransformationActionDelete:
			d.Headers[i].Value = nil
		case config.TransformationActionReplace:
			d.Headers[i].Value = t.To
		}

		for _, r := range d.Rows {
			if len(r) >= i {
				continue
			}

			if t.Modifier != nil {
				// TODO use modifier here to generate the new value if set
				// r[i].Value
			}

			switch t.Action {
			case config.TransformationActionAdd:
				r = append(r, r[i])
			case config.TransformationActionDelete:
				r[i].Value = nil
			case config.TransformationActionReplace:
				r[i].Value = r[i]
			}
		}
	}

	return nil
}

// CheckIfHeaderExists check if a header exists by name in the Table
func (d *Table) CheckIfHeaderExists(name interface{}) (int, bool) {
	for k, c := range d.Headers {
		if c.Value == name {
			return k, true
		}
	}

	return 0, false
}

// GetHeaderIndexByName return the header index for a given key (name) string
func (d *Table) GetHeaderIndexByName(name string) (int, error) {
	for i, h := range d.Headers {
		val, ok := h.Value.(string)
		if !ok {
			return -1, fmt.Errorf("failed to cast result header into string, header: %+v", h.Value)
		}
		if val == name {
			return i, nil
		}
	}
	return -1, nil
}
