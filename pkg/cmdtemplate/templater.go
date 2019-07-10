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

package cmdtemplate

import (
	"bytes"
	"text/template"

	"github.com/cloudical-io/acntt/testers"
)

// Variables variables used for templating
type Variables struct {
	ServerAddress string
	ServerPort    int32
}

// Template template a given cmd and args with the given host information struct
func Template(task *testers.Task, variables Variables) error {
	templatedArgs := []string{}

	// TODO take care of what variables to give because,
	// e.g., hosts with or without IPv6 must only get IPv4 addresses
	var err error
	task.Command, err = templateString(task.Command, variables)
	if err != nil {
		return err
	}

	for _, arg := range task.Args {
		arg, err = templateString(arg, variables)
		if err != nil {
			return err
		}
		templatedArgs = append(templatedArgs, arg)
	}
	task.Args = templatedArgs
	return nil
}

func templateString(in string, variable interface{}) (string, error) {
	t, err := template.New("main").Parse(in)
	if err != nil {
		return "", err
	}

	var out bytes.Buffer
	if err = t.ExecuteTemplate(&out, "main", variable); err != nil {
		return "", err
	}
	return out.String(), err
}
