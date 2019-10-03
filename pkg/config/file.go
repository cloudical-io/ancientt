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

package config

import (
	"io/ioutil"
	"os"

	"github.com/creasty/defaults"
	"gopkg.in/go-playground/validator.v9"
	"gopkg.in/yaml.v2"
)

var validate *validator.Validate

// Load load the given config file
func Load(cfgFile string) (*Config, error) {
	file, err := os.Open(cfgFile)
	if err != nil {
		return nil, err
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	cfg := New()

	if err := yaml.Unmarshal(content, cfg); err != nil {
		return nil, err
	}

	// Set defaults in the config struct
	if err := defaults.Set(cfg); err != nil {
		return nil, err
	}

	// Validate config struct
	validate = validator.New()
	if err := validate.Struct(cfg); err != nil {
		//validationErrors := err.(validator.ValidationErrors)
		return nil, err
	}

	return cfg, nil
}
