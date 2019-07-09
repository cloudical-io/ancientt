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

package testers

// Factories contains the list of all available testers.
var Factories = make(map[string]func() (Tester, error))

// Tester is the interface a tester has to implement
type Tester interface {
	Run(params map[string]interface{}, dryRun bool) error
}
