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
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cloudical-io/acntt/pkg/config"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// NameCSV CSV output name
const NameCSV = "csv"

func init() {
	Factories[NameCSV] = NewCSVOutput
}

// CSV CSV tester structure
type CSV struct {
	Output
	logger *log.Entry
	config *config.CSV
}

// NewCSVOutput return a new CSV tester instance
func NewCSVOutput(cfg *config.Config, outCfg *config.Output) (Output, error) {
	c := CSV{
		logger: log.WithFields(logrus.Fields{"output": NameCSV}),
		config: outCfg.CSV,
	}
	if c.config.NamePattern != "" {
		c.config.NamePattern = "{{ .UnixTime }}-{{ .Data.Tester }}-{{ .Data.ServerHost }}_{{ .Data.ClientHost }}.csv"
	}
	return c, nil
}

// Do make CSV outputs
func (c CSV) Do(data Data) error {
	dataTable, ok := data.Data.(Table)
	if !ok {
		return fmt.Errorf("data not in table for csv output")
	}

	filename, err := getFilenameFromPattern(c.config.NamePattern, data, nil)
	if err != nil {
		return err
	}

	outPath := filepath.Join(c.config.FilePath, filename)
	file, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Iterate over header columns
	for _, column := range dataTable.Headers {
		rowCells := []string{}
		for _, row := range column.Rows {
			rowCells = append(rowCells, fmt.Sprintf("%v", row.Value))
		}
		if len(rowCells) == 0 {
			continue
		}

		if err := writer.Write(rowCells); err != nil {
			return err
		}
	}

	// Iterate over data columns
	for _, column := range dataTable.Columns {
		rowCells := []string{}
		for _, row := range column.Rows {
			rowCells = append(rowCells, fmt.Sprintf("%v", row.Value))
		}
		if len(rowCells) == 0 {
			continue
		}

		if err := writer.Write(rowCells); err != nil {
			return err
		}
	}

	return nil
}
