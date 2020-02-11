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
	Headers []*Row
	Rows    [][]*Row
}

// Row Row of the Table data format
type Row struct {
	Value interface{}
}

// Transform transformation of table data
func (d Table) Transform(ts []*config.Transformation) error {
	// Iterate over each transformation
	for _, t := range ts {
		index, err := d.GetHeaderIndexByName(t.Source)
		if err != nil {
			return err
		}
		if index == -1 {
			return nil
		}

		switch t.Action {
		case config.TransformationActionAdd:
			d.Headers = append(d.Headers, &Row{
				Value: t.Destination,
			})
		case config.TransformationActionDelete:
			d.Headers[index] = nil
		case config.TransformationActionReplace:
			toHeader := t.Destination
			if toHeader == "" {
				toHeader = t.Source
			}
			d.Headers[index].Value = toHeader
		}

		for i := range d.Rows {
			if len(d.Rows[i]) < index {
				continue
			}

			switch t.Action {
			case config.TransformationActionAdd:
				d.Rows[i] = append(d.Rows[i], &Row{
					Value: d.modifyValue(d.Rows[i][index].Value, t),
				})
			case config.TransformationActionDelete:
				d.Rows[i][index] = nil
			case config.TransformationActionReplace:
				d.Rows[i][index].Value = d.modifyValue(d.Rows[i][index].Value, t)
			}
		}
	}

	return nil
}

func (d *Table) modifyValue(in interface{}, t *config.Transformation) interface{} {
	value, ok := in.(float64)
	if !ok {
		valInt, ok := in.(int64)
		if !ok {
			return in
		}
		value = float64(valInt)
	}

	switch t.ModifierAction {
	case config.ModifierActionAddition:
		return value + *t.Modifier
	case config.ModifierActionSubstract:
		return value - *t.Modifier
	case config.ModifierActionDivison:
		return value / *t.Modifier
	case config.ModifierActionMultiply:
		return value * *t.Modifier
	}

	return in
}

// CheckIfHeaderExists check if a header exists by name in the Table
func (d *Table) CheckIfHeaderExists(name interface{}) (int, bool) {
	for k, c := range d.Headers {
		if c == nil {
			continue
		}
		if c.Value == name {
			return k, true
		}
	}

	return 0, false
}

// GetHeaderIndexByName return the header index for a given key (name) string
func (d *Table) GetHeaderIndexByName(name string) (int, error) {
	for i := range d.Headers {
		if d.Headers[i] == nil {
			continue
		}
		val, ok := d.Headers[i].Value.(string)
		if !ok {
			return -1, fmt.Errorf("failed to cast result header into string, header: %+v", d.Headers[i].Value)
		}
		if val == name {
			return i, nil
		}
	}
	return -1, nil
}
