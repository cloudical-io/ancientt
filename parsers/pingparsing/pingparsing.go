/*
Copyright 2020 Cloudical Deutschland GmbH. All rights reserved.
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

package pingparsing

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/cloudical-io/ancientt/outputs"
	"github.com/cloudical-io/ancientt/parsers"
	"github.com/cloudical-io/ancientt/pkg/config"
	models "github.com/cloudical-io/ancientt/pkg/models/pingparsing"
	"github.com/cloudical-io/ancientt/pkg/util"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

// NamePingParsing PingParsing tester name
const NamePingParsing = "pingparsing"

func init() {
	parsers.Factories[NamePingParsing] = NewPingParsingTester
}

// PingParsing PingParsing tester structure
type PingParsing struct {
	parsers.Parser
	logger *log.Entry
	config *config.Test
}

// NewPingParsingTester return a new PingParsing tester instance
func NewPingParsingTester(cfg *config.Config, test *config.Test) (parsers.Parser, error) {
	return PingParsing{
		logger: log.WithFields(logrus.Fields{"parers": NamePingParsing}),
		config: test,
	}, nil
}

// Parse parse PingParsing JSON responses
func (p PingParsing) Parse(doneCh chan struct{}, inCh <-chan parsers.Input, dataCh chan<- outputs.Data) error {
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

func (p PingParsing) parse(input parsers.Input, dataCh chan<- outputs.Data) error {
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
	results := models.ClientResults{}
	if err := json.Unmarshal(logs.Bytes(), &results); err != nil {
		return err
	}

	table := &outputs.Table{
		Headers: []*outputs.Row{
			{Value: "test_time"},
			{Value: "round"},
			{Value: "tester"},
			{Value: "server_host"},
			{Value: "client_host"},
			{Value: "target"},
			{Value: "destination"},
			{Value: "packet_transmit"},
			{Value: "packet_receive"},
			{Value: "packet_loss_rate"},
			{Value: "packet_loss_count"},
			{Value: "rtt_min"},
			{Value: "rtt_avg"},
			{Value: "rtt_max"},
			{Value: "rtt_mdev"},
			{Value: "packet_duplicate_rate"},
			{Value: "packet_duplicate_count"},
			{Value: "timestamp"},
			{Value: "icmp_seq"},
			{Value: "ttl"},
			{Value: "time"},
			{Value: "duplicate"},
			{Value: "additional_info"},
		},
		Rows: [][]*outputs.Row{},
	}

	for name, r := range results {
		base := []*outputs.Row{
			{Value: input.TestTime.Format(util.TimeDateFormat)},
			{Value: input.Round},
			{Value: input.Tester},
			{Value: input.ServerHost},
			{Value: input.ClientHost},
			{Value: name},
			{Value: r.Destination},
			{Value: r.PacketTransmit},
			{Value: r.PacketReceive},
			{Value: r.PacketLossRate},
			{Value: r.PacketLossCount},
			{Value: r.RTTMin},
			{Value: r.RTTAvg},
			{Value: r.RTTMax},
			{Value: r.RTTMDev},
			{Value: r.PacketDuplicateRate},
			{Value: r.PacketDuplicateCount},
		}
		for _, e := range r.ICMPReplies {
			table.Rows = append(table.Rows, append(base, []*outputs.Row{
				{Value: e.Timestamp},
				{Value: e.ICMPSeq},
				{Value: e.TTL},
				{Value: e.Time},
				{Value: e.Duplicate},
				{Value: input.AdditionalInfo},
			}...))
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
		Data:           table,
	}

	p.logger.Debug("sending parsed data to dataCh")

	dataCh <- data

	// TODO generate sum and / or end table and send to output

	p.logger.Debug("sent parsed data to dataCh")

	return nil
}
