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
	"io"

	"github.com/cloudical-io/acntt/pkg/config"
)

// Factories contains the list of all available testers.
// The parser can each then be created using the function saved in the map.
var Factories = make(map[string]func(cfg *config.Config, test *config.Test) (Parser, error))

// Parser is the interface a parser has to implement
type Parser interface {
	// Parse
	Parse(stopCh chan struct{}, inCh <-chan Input) error
}

// Input structured parse
type Input struct {
	DataStream     *io.ReadCloser
	Data           []byte
	Tester         string
	ServerHost     string
	ClientHost     string
	AdditionalInfo string
}
