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

package dump

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cloudical-io/ancientt/outputs"
	"github.com/cloudical-io/ancientt/pkg/config"
	"github.com/k0kubun/pp"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// NameDump Dump output name
const NameDump = "dump"

func init() {
	outputs.Factories[NameDump] = NewDumpOutput
}

// Dump Dump tester structure
type Dump struct {
	outputs.Output
	logger *log.Entry
	config *config.Dump
	files  map[string]*os.File
}

// NewDumpOutput return a new Dump tester instance
func NewDumpOutput(cfg *config.Config, outCfg *config.Output) (outputs.Output, error) {
	dump := Dump{
		logger: log.WithFields(logrus.Fields{"output": NameDump}),
		config: outCfg.Dump,
		files:  map[string]*os.File{},
	}
	if dump.config.FilePath.NamePattern == "" {
		dump.config.FilePath.NamePattern = "ancientt-{{ .TestStartTime }}-{{ .Data.Tester }}-{{ .Data.ServerHost }}_{{ .Data.ClientHost }}.txt"
	}
	return dump, nil
}

// Do make Dump outputs
func (d Dump) Do(data outputs.Data) error {
	dataTable, ok := data.Data.(outputs.Table)
	if !ok {
		return fmt.Errorf("data not in data table format for dump output")
	}

	filename, err := outputs.GetFilenameFromPattern(d.config.FilePath.NamePattern, "", data, nil)
	if err != nil {
		return err
	}

	outPath := filepath.Join(d.config.FilePath.FilePath, filename)
	file, ok := d.files[outPath]
	if !ok {
		file, err = os.Create(outPath)
		if err != nil {
			return err
		}
		d.files[outPath] = file
	}

	// FIXME should the output be improved?

	if _, err := file.WriteString(pp.Sprint(dataTable)); err != nil {
		return err
	}

	return nil
}

// OutputFiles return a list of output files
func (d Dump) OutputFiles() []string {
	list := []string{}
	for file := range d.files {
		list = append(list, file)
	}
	return list
}

// Close close open files
func (d Dump) Close() error {
	for name, file := range d.files {
		d.logger.WithFields(logrus.Fields{"filepath": name}).Debug("closing file")
		if err := file.Close(); err != nil {
			d.logger.WithFields(logrus.Fields{"filepath": name}).Errorf("error closing file. %+v", err)
		}
	}

	return nil
}
