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

package iperf3

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/cloudical-io/ancientt/outputs"
	"github.com/cloudical-io/ancientt/parsers"
	"github.com/cloudical-io/ancientt/pkg/config"
	models "github.com/cloudical-io/ancientt/pkg/models/iperf3"
	"github.com/cloudical-io/ancientt/pkg/util"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// NameIPerf3 IPerf3 tester name
const NameIPerf3 = "iperf3"

func init() {
	parsers.Factories[NameIPerf3] = NewIPerf3Tester
}

// IPerf3 IPerf3 tester structure
type IPerf3 struct {
	parsers.Parser
	logger *log.Entry
	config *config.Test
}

// NewIPerf3Tester return a new IPerf3 tester instance
func NewIPerf3Tester(cfg *config.Config, test *config.Test) (parsers.Parser, error) {
	return IPerf3{
		logger: log.WithFields(logrus.Fields{"parers": NameIPerf3}),
		config: test,
	}, nil
}

// Parse parse IPerf3 JSON responses
func (p IPerf3) Parse(doneCh chan struct{}, inCh <-chan parsers.Input, dataCh chan<- outputs.Data) error {
	for {
		select {
		case <-doneCh:
			return nil
		case input, ok := <-inCh:
			if !ok {
				return nil
			}
			if input.ClientHost == "" && input.ServerHost == "" && input.Tester == "" {
				log.Warn("received input.Data with empty input.Tester and others are empty, 'signal' channel closed")
				close(dataCh)
				return nil
			}
			if err := p.parse(input, dataCh); err != nil {
				return err
			}
		}
	}
}

func (p IPerf3) parse(input parsers.Input, dataCh chan<- outputs.Data) error {
	var logs *bytes.Buffer
	if input.DataStream != nil {
		logs = new(bytes.Buffer)
		if _, err := io.Copy(logs, *input.DataStream); err != nil {
			return fmt.Errorf("error in copy information from logs to buffer")
		}
		if err := (*input.DataStream).Close(); err != nil {
			return fmt.Errorf("error during closing input.DataStream. %+v", err)
		}
	} else if len(input.Data) > 0 {
		// Directly pump the data in the logs var
		p.logger.Warn("received input.Data instead of input.DataStream, who wrote that runners without stream support")
		logs = bytes.NewBuffer(input.Data)
	} else {
		return fmt.Errorf("no data stream nor data from Input channel")
	}

	// Parse JSON response
	result := &models.ClientResult{}
	if err := json.Unmarshal(logs.Bytes(), result); err != nil {
		return err
	}

	intervalTable := &outputs.Table{
		Headers: []*outputs.Row{
			{Value: "test_time"},
			{Value: "round"},
			{Value: "tester"},
			{Value: "server_host"},
			{Value: "client_host"},
			{Value: "socket"},
			{Value: "start"},
			{Value: "end"},
			{Value: "seconds"},
			{Value: "bytes"},
			{Value: "bits_per_second"},
			{Value: "retransmits"},
			{Value: "snd_cwnd"},
			{Value: "rtt"},
			{Value: "rttvar"},
			{Value: "pmtu"},
			{Value: "omitted"},
			{Value: "iperf3_version"},
			{Value: "system_info"},
			{Value: "additional_info"},
		},
		Rows: [][]*outputs.Row{},
	}

	for _, interval := range result.Intervals {
		for _, stream := range interval.Streams {
			intervalTable.Rows = append(intervalTable.Rows, []*outputs.Row{
				{Value: input.TestTime.Format(util.TimeDateFormat)},
				{Value: input.Round},
				{Value: input.Tester},
				{Value: input.ServerHost},
				{Value: input.ClientHost},
				{Value: stream.Socket},
				{Value: stream.Start},
				{Value: stream.End},
				{Value: stream.Seconds},
				{Value: stream.Bytes},
				{Value: stream.BitsPerSecond},
				{Value: stream.Retransmits},
				{Value: stream.SndCwnd},
				{Value: stream.RTT},
				{Value: stream.RTTVar},
				{Value: stream.PMTU},
				{Value: stream.Omitted},
				{Value: result.Start.Version},
				{Value: result.Start.SystemInfo},
				{Value: input.AdditionalInfo},
			})
		}
	}

	p.logger.Debug("parsed data input")

	// Transform Input into outputs.Data struct
	data := outputs.Data{
		TestStartTime:  input.TestStartTime,
		TestTime:       input.TestTime,
		AdditionalInfo: input.AdditionalInfo,
		ServerHost:     input.ServerHost,
		ClientHost:     input.ClientHost,
		Tester:         input.Tester,
		Data:           intervalTable,
	}

	p.logger.Debug("sending parsed data to dataCh")

	dataCh <- data

	// TODO generate sum and / or end table and send to output

	p.logger.Debug("sent parsed data to dataCh")

	return nil
}
