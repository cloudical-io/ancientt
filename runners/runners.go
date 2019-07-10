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

package runners

import (
	"github.com/cloudical-io/acntt/parsers"
	"github.com/cloudical-io/acntt/pkg/config"
	"github.com/cloudical-io/acntt/testers"
)

// Factories contains the list of all available runners.
// The runners can each then be created using the function saved in the map.
var Factories = make(map[string]func(cfg *config.Config) (Runner, error))

// Runner is the interface a runner has to implement.
type Runner interface {
	GetHostsForTest(test config.Test) (*testers.Hosts, error)
	Prepare(runOpts config.RunOptions, plan *testers.Plan) error
	Execute(plan *testers.Plan, parser parsers.Parser) error
	Cleanup(plan *testers.Plan) error
}
