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

package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/cloudical-io/acntt/pkg/config"
	_ "github.com/cloudical-io/acntt/runners"
	_ "github.com/cloudical-io/acntt/testers"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

var rootCmd = &cobra.Command{
	Use:   "acntt",
	Short: "ACNTT is a automated continuous network testing tool which utilizes iperf, siege and others.",
	RunE:  run,
}

func init() {
	rootCmd.PersistentFlags().StringP("testdefinition", "c", "", "Path to the testdefinitions to read for the tests.")
	viper.BindPFlag("testdefinition", rootCmd.PersistentFlags().Lookup("testdefinition"))
	viper.SetDefault("testdefinition", "testdefinition.yaml")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	cfgFile := viper.GetString("testdefinition")

	file, err := os.Open(cfgFile)
	if err != nil {
		return err
	}

	cfgContent, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	cfg := &config.Config{}

	if err := yaml.Unmarshal(cfgContent, cfg); err != nil {
		return err
	}

	fmt.Printf("TEST: %+v\n", cfg)

	return nil
}
