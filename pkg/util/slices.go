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

// UniqueStringSlice return deduplicated string slice from one or more ins.
func UniqueStringSlice(ins ...[]string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, in := range ins {
		for _, k := range in {
			if _, value := keys[k]; !value {
				keys[k] = true
				list = append(list, k)
			}
		}
	}
	return list
}
