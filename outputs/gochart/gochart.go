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
	"path/filepath"

	"github.com/cloudical-io/ancientt/outputs"
	"github.com/cloudical-io/ancientt/pkg/config"
	"github.com/cloudical-io/ancientt/pkg/util"
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
	files  map[string]struct{}
}

// NewGoChartOutput return a new GoChart tester instance
func NewGoChartOutput(cfg *config.Config, outCfg *config.Output) (outputs.Output, error) {
	goChart := GoChart{
		logger: log.WithFields(logrus.Fields{"output": NameGoChart}),
		config: outCfg.GoChart,
		files:  map[string]struct{}{},
	}
	if goChart.config.NamePattern == "" {
		goChart.config.NamePattern = "ancientt-{{ .TestStartTime }}-{{ .Data.Tester }}-{{ .Data.ServerHost }}_{{ .Data.ClientHost }}-{{ .Extra.Axises }}.png"
	}
	return goChart, nil
}

// Do make GoChart charts
func (gc GoChart) Do(data outputs.Data) error {
	if _, ok := data.Data.(outputs.Table); !ok {
		return fmt.Errorf("data not in data table format for gochart output")
	}

	// Iterate over wanted graph types
	for _, graph := range gc.config.Graphs {
		err := gc.drawAxisChart(graph, data)
		if err != nil {
			return err
		}
	}

	return nil
}

func (gc *GoChart) drawAxisChart(chartOpts *config.GoChartGraph, data outputs.Data) error {
	dataTable, ok := data.Data.(outputs.Table)
	if !ok {
		return fmt.Errorf("data not in table format for gochart output")
	}

	if len(dataTable.Headers) == 0 {
		gc.logger.Warning("no table headers found in data table result, returning")
		return nil
	}

	timeKey := chartOpts.TimeColumn

	var leftYKey string
	if chartOpts.LeftY != "" {
		leftYKey = chartOpts.LeftY
	}
	rightYKey := chartOpts.RightY

	graph := chart.Chart{
		Series: []chart.Series{},
		XAxis: chart.XAxis{
			Name: timeKey,
		},
		YAxis: chart.YAxis{
			Name: rightYKey,
		},
	}

	axises := rightYKey
	if chartOpts.LeftY != "" && leftYKey != "" {
		axises += "_" + leftYKey
	}
	filename, err := outputs.GetFilenameFromPattern(gc.config.FilePath.NamePattern, "", data, map[string]interface{}{
		"Axises": axises,
	})
	if err != nil {
		return err
	}
	outPath := filepath.Join(gc.config.FilePath.FilePath, filename)

	vals := map[string][]float64{}
	for _, search := range []string{chartOpts.TimeColumn, chartOpts.RightY, chartOpts.LeftY} {
		headIndex, err := dataTable.GetHeaderIndexByName(search)
		if err != nil {
			return err
		}

		for _, r := range dataTable.Rows {
			// Skip empty rows
			if len(r) == 0 {
				continue
			}
			if len(r)-1 < headIndex || r[headIndex] == nil {
				return fmt.Errorf("unable to find header with index %d or nil (search: %q)", headIndex, search)
			}

			val, err := util.CastNumberToFloat64(r[headIndex].Value)
			if err != nil {
				return err
			}
			vals[search] = append(vals[search], val)
		}
	}

	series := chart.ContinuousSeries{
		Name: rightYKey,
		Style: chart.Style{
			StrokeColor: chart.GetDefaultColor(1).WithAlpha(64),
			FillColor:   chart.GetDefaultColor(1).WithAlpha(64),
		},
		XValues: vals[timeKey],
		YValues: vals[rightYKey],
	}
	graph.Series = append(graph.Series, series)
	gc.additionalSeries(chartOpts, &graph, &series)

	if chartOpts.LeftY != "" {
		if _, ok := vals[leftYKey]; ok {
			series = chart.ContinuousSeries{
				Name: leftYKey,
				Style: chart.Style{
					StrokeColor: chart.GetDefaultColor(4).WithAlpha(64),
					FillColor:   chart.GetDefaultColor(4).WithAlpha(64),
				},
				YAxis:   chart.YAxisSecondary,
				XValues: vals[timeKey],
				YValues: vals[leftYKey],
			}
			graph.Series = append(graph.Series, series)
			gc.additionalSeries(chartOpts, &graph, &series)
		}
	}

	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	gc.files[outPath] = struct{}{}

	buffer := bytes.NewBuffer([]byte{})
	if err := graph.Render(chart.PNG, buffer); err != nil {
		return fmt.Errorf("failed to render graph to PNG file. %+v", err)
	}

	if err := util.WriteNewTruncFile(outPath, buffer.Bytes()); err != nil {
		return err
	}

	return nil
}

func (gc *GoChart) additionalSeries(chartOpts *config.GoChartGraph, graph *chart.Chart, series *chart.ContinuousSeries) {
	graph.Series = append(graph.Series, chart.LastValueAnnotationSeries(series), chart.LastValueAnnotationSeries(series))

	if chartOpts.WithLinearRegression != nil && *chartOpts.WithLinearRegression {
		linearRegresSeries := &chart.LinearRegressionSeries{
			Name:        fmt.Sprintf("%q - LinearRegress", series.Name),
			InnerSeries: series,
		}
		graph.Series = append(graph.Series, linearRegresSeries)
	}
	if chartOpts.WithSimpleMovingAverage != nil && *chartOpts.WithSimpleMovingAverage {
		smaSeries := &chart.SMASeries{
			Name:        fmt.Sprintf("%q - SimpleMovingAvg", series.Name),
			InnerSeries: series,
		}
		graph.Series = append(graph.Series, smaSeries)
	}
}

// OutputFiles return a list of output files
func (gc GoChart) OutputFiles() []string {
	list := []string{}
	for file := range gc.files {
		list = append(list, file)
	}
	return list
}

// Close NOOP, as graph pictures are written once and closed immediately, no need to do anything here
func (gc GoChart) Close() error {
	return nil
}
