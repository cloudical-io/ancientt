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

// Template template a given cmd and args with the given host information struct
func Template(cmd string, args []string, host *testers.Host) (string, []string, error) {
	// Anonymous variables structure for the command and args templating
	variables := struct {
		ServerAddress *testers.IPAddresses
	}{
		ServerAddress: host.Addresses,
	}

	templatedArgs := []string{}

	// TODO take care of what variables to give because,
	// e.g., hosts with or without IPv6 must only get IPv4 addresses
	var err error
	cmd, err = templateString(cmd, variables)
	if err != nil {
		return "", nil, err
	}

	for _, arg := range args {
		arg, err = templateString(arg, variables)
		if err != nil {
			return "", nil, err
		}
		templatedArgs = append(templatedArgs, arg)
	}
	return cmd, templatedArgs, nil
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
