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

package excelize

import (
	"fmt"
	"path"

	//include excelize library for .xlsx output
	"github.com/360EntSecGroup-Skylar/excelize/v2"

	"github.com/cloudical-io/ancientt/outputs"
	"github.com/cloudical-io/ancientt/pkg/config"
	"github.com/cloudical-io/ancientt/pkg/util"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// NameExcelize Excelize output name
const NameExcelize = "excelize"

func init() {
	outputs.Factories[NameExcelize] = NewExcelizeOutput
}

// Excelize Excelize tester structure
type Excelize struct {
	outputs.Output
	logger *log.Entry
	config *config.Excelize
	files  map[string]*fileState
}

type fileState struct {
	file *excelize.File
	row  int
}

// NewExcelizeOutput return a new Excelize tester instance
func NewExcelizeOutput(cfg *config.Config, outCfg *config.Output) (outputs.Output, error) {
	excelize := Excelize{
		logger: log.WithFields(logrus.Fields{"output": NameExcelize}),
		config: outCfg.Excelize,
		files:  map[string]*fileState{},
	}
	if excelize.config.FilePath.NamePattern == "" {
		excelize.config.FilePath.NamePattern = "ancientt-{{ .TestStartTime }}-{{ .Data.Tester }}-{{ .Data.ServerHost }}_{{ .Data.ClientHost }}.xlsx"
	}
	if excelize.config.SaveAfterRows == 0 {
		excelize.config.SaveAfterRows = 200
	}
	return excelize, nil
}

// Do Inputs the data into the excel sheet, contains all logic necessary to perform this task
func (e Excelize) Do(data outputs.Data) error {
	dataTable, ok := data.Data.(outputs.Table)
	if !ok {
		return fmt.Errorf("data not in Table data type for excel output")
	}

	outputFilename, err := outputs.GetFilenameFromPattern(e.config.FilePath.NamePattern, "", data, nil)
	if err != nil {
		return err
	}
	filePath := path.Join(e.config.FilePath.FilePath, outputFilename)

	// Check if the file had already been opened, reuse it if so
	var fState *fileState
	if _, ok := e.files[filePath]; !ok {
		// Create new excelize file
		excelFile := excelize.NewFile()
		excelFile.Path = filePath

		// Initial state for a new file
		// The fileState of a file will keep the *excelize.Fileand the current row
		// Current row is needed if the file is reused as otherwise it would start
		// at the first row again
		state := &fileState{
			file: excelFile,
			row:  1,
		}
		fState = state
		e.files[filePath] = state
	} else {
		fState = e.files[filePath]
	}

	// Initially save file on each (re-)use
	if err = fState.file.Save(); err != nil {
		return err
	}

	if fState.row == 1 {
		if err := e.inputData(fState.row, [][]outputs.Row{dataTable.Headers}, fState); err != nil {
			return err
		}
	}
	if err := e.inputData(fState.row, dataTable.Rows, fState); err != nil {
		return err
	}

	// NOTE If this isn't enough use the `Close()` func
	if err = fState.file.Save(); err != nil {
		return err
	}

	return nil
}

func (e Excelize) inputData(startRow int, rows [][]outputs.Row, fState *fileState) error {
	// Iterate over data columns to get the first row of data.
	for i, row := range rows {
		fState.row++

		// Set each cell value
		for j, row := range row {
			if err := fState.file.SetCellValue("Sheet1", fmt.Sprintf("%s%d", util.IntToChar(j+1), startRow+i), row.Value); err != nil {
				// TODO Return a final concated error after the whole data has been written
				e.logger.WithFields(logrus.Fields{"filepath": fState.file.Path}).Errorf("unable to set cell value in excelize file. %+v", err)
			}
		}

		if i != 0 && (e.config.SaveAfterRows%i) == 0 {
			if err := fState.file.Save(); err != nil {
				return err
			}
		}
	}

	return nil
}

// OutputFiles return a list of output files
func (e Excelize) OutputFiles() []string {
	list := []string{}
	for file := range e.files {
		list = append(list, file)
	}
	return list
}

// Close Nothing to do here, all files are written during "creation" and "usage" (writing data) to the files
func (e Excelize) Close() error {
	return nil
}

func toChar(i int) rune {
	return rune('A' - 1 + i)
}
