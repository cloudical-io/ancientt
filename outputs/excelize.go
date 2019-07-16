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
	"os"

	"github.com/cloudical-io/acntt/pkg/config"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// NameExcelize Excelize output name
const NameExcelize = "excelize"

func init() {
	Factories[NameExcelize] = NewExcelizeOutput
}

// Excelize Excelize tester structure
type Excelize struct {
	Output
	logger *log.Entry
	config *config.Excelize
	files  map[string]*os.File
}

// NewExcelizeOutput return a new Excelize tester instance
func NewExcelizeOutput(cfg *config.Config, outCfg *config.Output) (Output, error) {
	excelize := Excelize{
		logger: log.WithFields(logrus.Fields{"output": NameExcelize}),
		config: outCfg.Excelize,
		files:  map[string]*os.File{},
	}
	if excelize.config.FilePath == "" {
		excelize.config.FilePath = "."
	}
	if excelize.config.NamePattern != "" {
		excelize.config.NamePattern = "acntt-{{ .TestStartTime }}-{{ .Data.Tester }}-{{ .Data.ServerHost }}_{{ .Data.ClientHost }}.xlsx"
	}
	return excelize, nil
}

// Do TODO Implement
func (e Excelize) Do(data Data) error {
	return nil
}

// Close TODO Implement
func (e Excelize) Close() error {
	return nil
}
