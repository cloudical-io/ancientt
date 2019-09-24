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

package main

// This file contains the imports for each output, parser, runner and tester.
// Importing them from the, e.g., `outputs` pkg would cause a import cycle.

import (
	// Outputs
	_ "github.com/cloudical-io/ancientt/outputs/csv"
	_ "github.com/cloudical-io/ancientt/outputs/dump"
	_ "github.com/cloudical-io/ancientt/outputs/excelize"
	_ "github.com/cloudical-io/ancientt/outputs/gochart"
	_ "github.com/cloudical-io/ancientt/outputs/mysql"
	_ "github.com/cloudical-io/ancientt/outputs/sqlite"

	// Parsers
	_ "github.com/cloudical-io/ancientt/parsers/iperf3"

	// Runners
	_ "github.com/cloudical-io/ancientt/runners/ansible"
	_ "github.com/cloudical-io/ancientt/runners/kubernetes"
	_ "github.com/cloudical-io/ancientt/runners/mock"

	// Testers
	_ "github.com/cloudical-io/ancientt/testers/iperf3"
)
