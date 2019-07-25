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

package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

	"github.com/cloudical-io/acntt/pkg/config"
	"github.com/cloudical-io/acntt/pkg/util"
	"github.com/cloudical-io/acntt/outputs"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// NameCSV CSV output name
const NameCSV = "csv"

func init() {
	outputs.Factories[NameCSV] = NewCSVOutput
}

// CSV CSV tester structure
type CSV struct {
	outputs.Output
	logger  *log.Entry
	config  *config.CSV
	files   map[string]*os.File
	writers map[string]*csv.Writer
}

// NewCSVOutput return a new CSV tester instance
func NewCSVOutput(cfg *config.Config, outCfg *config.Output) (outputs.Output, error) {
	c := CSV{
		logger:  log.WithFields(logrus.Fields{"output": NameCSV}),
		config:  outCfg.CSV,
		files:   map[string]*os.File{},
		writers: map[string]*csv.Writer{},
	}
	if c.config.FilePath == "" {
		c.config.FilePath = "."
	}
	if c.config.NamePattern == "" {
		c.config.NamePattern = "acntt-{{ .TestStartTime }}-{{ .Data.Tester }}-{{ .Data.ServerHost }}_{{ .Data.ClientHost }}.csv"
	}
	return c, nil
}

// Do make CSV outputs
func (c CSV) Do(data outputs.Data) error {
	dataTable, ok := data.Data.(outputs.Table)
	if !ok {
		return fmt.Errorf("data not in table for csv output")
	}

	filename, err := outputs.GetFilenameFromPattern(c.config.NamePattern, "", data, nil)
	if err != nil {
		return err
	}

	var writeHeaders bool

	outPath := filepath.Join(c.config.FilePath, filename)
	writer, ok := c.writers[outPath]
	if !ok {
		file, ok := c.files[outPath]
		if !ok {
			file, err = os.Create(outPath)
			if err != nil {
				return err
			}
			c.files[outPath] = file
		}

		writer = csv.NewWriter(file)
		c.writers[outPath] = writer
		writeHeaders = true
	}

	defer writer.Flush()

	if writeHeaders {
		// Iterate over header columns
		for _, column := range dataTable.Headers {
			rowCells := []string{}
			for _, row := range column.Rows {
				rowCells = append(rowCells, util.CastToString(row.Value))
			}
			if len(rowCells) == 0 {
				continue
			}

			if err := writer.Write(rowCells); err != nil {
				return err
			}
		}
	}

	// Iterate over data columns
	for _, column := range dataTable.Columns {
		rowCells := []string{}
		for _, row := range column.Rows {
			rowCells = append(rowCells, util.CastToString(row.Value))
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

// Close close all file descriptors here
func (c CSV) Close() error {
	for name, writer := range c.writers {
		c.logger.WithFields(logrus.Fields{"filepath": name}).Debug("closing file")
		writer.Flush()
		if err := writer.Error(); err != nil {
			c.logger.WithFields(logrus.Fields{"filepath": name}).Errorf("error flushing file. %+v", err)
		}
	}

	for name, file := range c.files {
		c.logger.WithFields(logrus.Fields{"filepath": name}).Debug("closing file")
		if err := file.Close(); err != nil {
			c.logger.WithFields(logrus.Fields{"filepath": name}).Errorf("error closing file. %+v", err)
		}
	}

	return nil

}
