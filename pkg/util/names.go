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

package util

import (
	"fmt"
	"time"
)

// PNameRole task role names type
type PNameRole string

const (
	// PNameRoleClient Client role name
	PNameRoleClient PNameRole = "client"
	// PNameRoleServer Server role name
	PNameRoleServer PNameRole = "server"
)

// GetPNameFromTask get a "persistent" name for a task
// This is done by calculating the checksums of the used names.
func GetPNameFromTask(round int, hostname string, command string, role PNameRole, testStartTime time.Time) string {
	return fmt.Sprintf("ancientt-%s-%s-%d", role, command, testStartTime.UnixNano())
}

// GetTaskName get a task name
func GetTaskName(tester string, testStartTime time.Time) string {
	return fmt.Sprintf("ancientt-%s-%d", tester, testStartTime.Unix())
}
