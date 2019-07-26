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

package gochart

import (
	"bytes"
	"fmt"

	"github.com/cloudical-io/acntt/outputs"
	"github.com/cloudical-io/acntt/pkg/config"
	"github.com/cloudical-io/acntt/pkg/util"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	chart "github.com/wcharczuk/go-chart"
)

// NameGoChart GoChart output name
const NameGoChart = "gochart"

func init() {
	outputs.Factories[NameGoChart] = NewGoChartOutput
}

// GoChart GoChart tester structure
type GoChart struct {
	outputs.Output
	logger *log.Entry
	config *config.GoChart
}

// NewGoChartOutput return a new GoChart tester instance
func NewGoChartOutput(cfg *config.Config, outCfg *config.Output) (outputs.Output, error) {
	goChart := GoChart{
		logger: log.WithFields(logrus.Fields{"output": NameGoChart}),
		config: outCfg.GoChart,
	}
	if goChart.config.FilePath == "" {
		goChart.config.FilePath = "."
	}
	if goChart.config.NamePattern != "" {
		goChart.config.NamePattern = "acntt-{{ .TestStartTime }}-{{ .Data.Tester }}-{{ .Data.ServerHost }}_{{ .Data.ClientHost }}-{{ .Extra.Header }}-{{ .Extra.Type }}.png"
	}
	return goChart, nil
}

// Do make GoChart charts
func (gc GoChart) Do(data outputs.Data) error {
	dataTable, ok := data.Data.(outputs.Table)
	if !ok {
		return fmt.Errorf("data not in table for csv output")
	}

	// Iterate over wanted graph types
	// TODO Allow certain header columns to be selected per graphType
	for _, graphType := range gc.config.Types {
		for _, column := range dataTable.Headers {
			for _, row := range column.Rows {
				filename, err := outputs.GetFilenameFromPattern(gc.config.NamePattern, "", data, map[string]interface{}{
					"Type":   graphType,
					"Header": row.Value,
				})
				if err != nil {
					return err
				}

				graph := chart.Chart{Series: []chart.Series{chart.ContinuousSeries{}}}

				// TODO create graphs per column

				// Iterate over header columns

				// Iterate over data columns
				for _, column := range dataTable.Columns {
					rowCells := []string{}
					for _, row := range column.Rows {
						rowCells = append(rowCells, fmt.Sprintf("%v", row.Value))
					}
					if len(rowCells) == 0 {
						continue
					}

				}

				//XValues: []float64{1.0, 2.0, 3.0, 4.0},
				//YValues: []float64{1.0, 2.0, 3.0, 4.0},

				buffer := bytes.NewBuffer([]byte{})
				if err := graph.Render(chart.PNG, buffer); err != nil {
					return err
				}

				if err := util.WriteNewTruncFile(filename, buffer.Bytes()); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// Close NOOP, as graph pictures are written once (= closed immediately)
func (gc GoChart) Close() error {
	return nil
}
