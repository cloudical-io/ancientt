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
	"os"
	"path/filepath"

	"github.com/cloudical-io/acntt/pkg/config"
	"github.com/k0kubun/pp"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// NameDump Dump output name
const NameDump = "dump"

func init() {
	Factories[NameDump] = NewDumpOutput
}

// Dump Dump tester structure
type Dump struct {
	Output
	logger *log.Entry
	config *config.Dump
	files  map[string]*os.File
}

// NewDumpOutput return a new Dump tester instance
func NewDumpOutput(cfg *config.Config, outCfg *config.Output) (Output, error) {
	dump := Dump{
		logger: log.WithFields(logrus.Fields{"output": NameDump}),
		config: outCfg.Dump,
		files:  map[string]*os.File{},
	}
	if dump.config.FilePath == "" {
		dump.config.FilePath = "."
	}
	if dump.config.NamePattern != "" {
		dump.config.NamePattern = "acntt-{{ .TestStartTime }}-{{ .Data.Tester }}-{{ .Data.ServerHost }}_{{ .Data.ClientHost }}.txt"
	}
	return dump, nil
}

// Do make Dump outputs
func (d Dump) Do(data Data) error {
	dataTable, ok := data.Data.(Table)
	if !ok {
		return fmt.Errorf("data not in table for dump output")
	}

	filename, err := getFilenameFromPattern(d.config.NamePattern, "", data, nil)
	if err != nil {
		return err
	}

	outPath := filepath.Join(d.config.FilePath, filename)
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
