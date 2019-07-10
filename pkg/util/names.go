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

package util

import (
	"crypto/sha1"
	"fmt"
	"os"

	"github.com/cloudical-io/acntt/testers"
)

// GetPNameFromTask get a "persistent" name for a task
// This is done by calculating the checksums of the used names.
func GetPNameFromTask(round int, task testers.Task) string {
	data := fmt.Sprintf("%d-%s-%s", round, task.Host.Name, task.Args)
	return fmt.Sprintf("%s-%s-%x", os.Args[0], task.Command, sha1.Sum([]byte(data)))
}
