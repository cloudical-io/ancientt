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

	"github.com/cloudical-io/ancientt/outputs"
	"github.com/cloudical-io/ancientt/pkg/config"
	"github.com/cloudical-io/ancientt/pkg/util"
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
	if c.config.FilePath.NamePattern == "" {
		c.config.FilePath.NamePattern = "ancientt-{{ .TestStartTime }}-{{ .Data.Tester }}.csv"
	}
	return c, nil
}

// Do make CSV outputs
func (c CSV) Do(data outputs.Data) error {
	dataTable, ok := data.Data.(outputs.Table)
	if !ok {
		return fmt.Errorf("data not in Table interface format for csv output")
	}

	filename, err := outputs.GetFilenameFromPattern(c.config.FilePath.NamePattern, "", data, nil)
	if err != nil {
		return err
	}

	var writeHeaders bool

	outPath := filepath.Join(c.config.FilePath.FilePath, filename)
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
		writer.Comma = *c.config.Separator
		c.writers[outPath] = writer
		writeHeaders = true
	}

	defer writer.Flush()

	if writeHeaders {
		// Iterate over header columns
		headers := []string{}
		for _, r := range dataTable.Headers {
			if r == nil {
				continue
			}
			headers = append(headers, util.CastToString(r.Value))
		}
		if err := writer.Write(headers); err != nil {
			return err
		}
	}

	// Iterate over data columns
	for _, row := range dataTable.Rows {
		cells := []string{}
		for _, r := range row {
			if r == nil {
				continue
			}
			cells = append(cells, util.CastToString(r.Value))
		}
		if len(cells) == 0 {
			continue
		}

		if err := writer.Write(cells); err != nil {
			return err
		}
	}

	return nil
}

// OutputFiles return a list of output files
func (c CSV) OutputFiles() []string {
	list := []string{}
	for file := range c.files {
		list = append(list, file)
	}
	return list
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
