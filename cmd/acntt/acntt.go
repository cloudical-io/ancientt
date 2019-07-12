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

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/cloudical-io/acntt/outputs"
	"github.com/cloudical-io/acntt/parsers"
	"github.com/cloudical-io/acntt/pkg/config"
	"github.com/cloudical-io/acntt/runners"
	"github.com/cloudical-io/acntt/testers"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
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
	rootCmd.PersistentFlags().Bool("no-cleanup", false, "If runners should run cleanup routines after the tests.")
	rootCmd.PersistentFlags().Bool("yes", false, "Ask for user confirmation for each test before executing it.")
	rootCmd.PersistentFlags().StringP("testdefinition", "c", "", "Path to the testdefinitions to read for the tests.")
	viper.BindPFlag("no-cleanup", rootCmd.PersistentFlags().Lookup("no-cleanup"))
	viper.BindPFlag("yes", rootCmd.PersistentFlags().Lookup("yes"))
	viper.BindPFlag("testdefinition", rootCmd.PersistentFlags().Lookup("testdefinition"))
	viper.SetDefault("yes", false)
	viper.SetDefault("no-cleanup", false)
	viper.SetDefault("testdefinition", "testdefinition.yaml")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetReportCaller(false)
	log.SetLevel(log.DebugLevel)

	cfgFile := viper.GetString("testdefinition")

	file, err := os.Open(cfgFile)
	if err != nil {
		return err
	}

	cfgContent, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	cfg := config.New()

	if err := yaml.Unmarshal(cfgContent, cfg); err != nil {
		return err
	}

	var runner runners.Runner
	runnerName := strings.ToLower(cfg.Runner.Name)
	runnerNewFunc, ok := runners.Factories[runnerName]
	if !ok {
		return fmt.Errorf("runner with name %s not found", runnerName)
	}
	if runner, err = runnerNewFunc(cfg); err != nil {
		return err
	}

	for _, test := range cfg.Tests {
		// Get tester for the test
		var tester testers.Tester
		testerName := strings.ToLower(test.Type)
		testerNewFunc, ok := testers.Factories[testerName]
		if !ok {
			return fmt.Errorf("tester with name %s not found", testerName)
		}
		if tester, err = testerNewFunc(cfg, &test); err != nil {
			return err
		}

		// Get parser for the tester output
		var parser parsers.Parser
		parserNewFunc, ok := parsers.Factories[testerName]
		if !ok {
			return fmt.Errorf("parser with name %s not found", testerName)
		}
		if parser, err = parserNewFunc(cfg, &test); err != nil {
			return err
		}

		// Get hosts for the test
		hosts, err := runner.GetHostsForTest(test)
		if err != nil {
			return err
		}
		// Create testers.Environment with the hosts
		env := &testers.Environment{
			Hosts: hosts,
		}
		// Get plan from testers.Plan()
		plan, err := tester.Plan(env, &test)
		if err != nil {
			return err
		}
		// Print the plan
		fmt.Println("--> BEGIN PLAN <--")
		plan.PrettyPrint()
		fmt.Println("-->  END PLAN  <--")

		if !viper.GetBool("yes") {
			// Ask user if we can continue or not
			fmt.Println("===================")
			for {
				reader := bufio.NewReader(os.Stdin)
				fmt.Print("Are you sure you want to continue with above test plan? ('Yes' or 'No'): ")
				userInput, err := reader.ReadString('\n')
				if err != nil {
					return err
				}
				userInput = strings.TrimSpace(strings.ToLower(userInput))
				if userInput == "yes" || userInput == "y" {
					break
				} else if userInput == "no" || userInput == "n" {
					return fmt.Errorf("aborted by user")
				}
			}
			fmt.Println("===================")
		}

		log.WithFields(logrus.Fields{"testers": test.Type}).Info("preparing test")

		// Prepare the runner for the plan
		if err = runner.Prepare(test.RunOptions, plan); err != nil {
			return err
		}

		var wg sync.WaitGroup

		doneCh := make(chan struct{})
		inCh := make(chan parsers.Input)
		dataCh := make(chan outputs.Data)

		errs := make(chan error)

		go func() {
			select {
			case erro := <-errs:
				log.Error(erro.Error())
			case <-doneCh:
				return
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			err := parser.Parse(doneCh, inCh, dataCh)
			// Close dataCh as there won't be anything else coming through
			close(dataCh)
			if err != nil {
				errs <- err
				return
			}
		}()

		// Start each output
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case data, ok := <-dataCh:
					if !ok {
						return
					}
					for _, outputItem := range test.Outputs {
						outputName := outputItem.Name

						var output outputs.Output
						outputNewFunc, ok := outputs.Factories[outputName]
						if !ok {
							errs <- fmt.Errorf("output with name %s not found", outputName)
							return
						}
						output, err = outputNewFunc(cfg, &outputItem)
						if err != nil {
							errs <- err
							return
						}

						if err := output.Do(data); err != nil {
							errs <- err
							return
						}
					}
				case <-doneCh:
					return
				}
			}
		}()

		log.WithFields(logrus.Fields{"testers": test.Type}).Info("executing test")

		// Execute the plan
		if err := runner.Execute(plan, inCh); err != nil {
			return err
		}
		log.WithFields(logrus.Fields{"testers": test.Type}).Debug("runner execute returned, closing doneCh")

		close(inCh)
		wg.Wait()

		close(doneCh)

		// Run runners.Cleanup() func if wanted by the user
		// TODO When wanted, should be run when signal (CTRL+C) received
		if !viper.GetBool("no-cleanup") {
			log.WithFields(logrus.Fields{"testers": test.Type}).Info("running cleanup for test")
			if err := runner.Cleanup(plan); err != nil {
				return err
			}
		}
	}

	log.Info("done with tests")

	return nil
}
