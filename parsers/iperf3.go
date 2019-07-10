/*
Copyright 2019 Cloudical Deutschland GmbH
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

package parsers

import (
	"bytes"

	"github.com/cloudical-io/acntt/pkg/config"
)

// NameIPerf3 IPerf3 tester name
const NameIPerf3 = "iperf3"

func init() {
	Factories[NameIPerf3] = NewIPerf3Tester
}

// IPerf3 IPerf3 tester structure
type IPerf3 struct {
	Parser
	config *config.IPerf3
}

// NewIPerf3Tester return a new IPerf3 tester instance
func NewIPerf3Tester(cfg *config.Config, test *config.Test) (Parser, error) {
	return IPerf3{
		config: test.IPerf3,
	}, nil
}

// Parse parse IPerf3 JSON responses
func (ip IPerf3) Parse(in bytes.Buffer) ([]byte, error) {
	// TODO parse input

	return []byte{}, nil
}
