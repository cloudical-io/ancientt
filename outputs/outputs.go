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
	"bytes"
	"html/template"
	"time"

	"github.com/cloudical-io/acntt/pkg/config"
)

// Factories contains the list of all available outputs.
// The outputs can each then be created using the function saved in the map.
var Factories = make(map[string]func(cfg *config.Config, outCfg *config.Output) (Output, error))

// Output is the interface a output has to implement.
type Output interface {
	Do(data Data) error
}

// Data structured parsed data
type Data struct {
	Tester         string
	ServerHost     string
	ClientHost     string
	AdditionalInfo string
	Data           DataFormat
}

// DataFormat
type DataFormat interface {
}

// Table
type Table struct {
	DataFormat
	Headers []Column
	Columns []Column
}

// Column
type Column struct {
	Rows []Row
}

// Row
type Row struct {
	Value interface{}
}

// getFilenameFromPattern
func getFilenameFromPattern(pattern string, data Data, additionalData map[string]interface{}) (string, error) {
	t, err := template.New("main").Parse(pattern)
	if err != nil {
		return "", err
	}

	variables := struct {
		Data     Data
		UnixTime int64
		Extra    map[string]interface{}
	}{
		Data:     data,
		UnixTime: time.Now().Unix(),
		Extra:    additionalData,
	}

	var out bytes.Buffer
	if err = t.ExecuteTemplate(&out, "main", variables); err != nil {
		return "", err
	}
	return out.String(), nil
}
