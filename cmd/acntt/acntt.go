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
	"time"

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

var (
	rootCmd = &cobra.Command{
		Use:   "acntt",
		Short: "ACNTT is a automated continuous network testing tool which utilizes iperf, siege and others.",
		RunE:  run,
	}
	cfg *config.Config
)

func init() {
	// Set flags, viper binds for flags and viper bind default values
	rootCmd.PersistentFlags().Bool("print-plan", true, "If the plan of the tester should be printed to console.")
	rootCmd.PersistentFlags().Bool("no-cleanup", false, "If runners should run cleanup routines after the tests.")
	rootCmd.PersistentFlags().Bool("yes", false, "Ask for user confirmation for each test before executing it.")
	rootCmd.PersistentFlags().StringP("testdefinition", "c", "", "Path to the testdefinitions to read for the tests.")
	viper.BindPFlag("print-plan", rootCmd.PersistentFlags().Lookup("print-plan"))
	viper.BindPFlag("no-cleanup", rootCmd.PersistentFlags().Lookup("no-cleanup"))
	viper.BindPFlag("yes", rootCmd.PersistentFlags().Lookup("yes"))
	viper.BindPFlag("testdefinition", rootCmd.PersistentFlags().Lookup("testdefinition"))
	viper.SetDefault("print-plan", true)
	viper.SetDefault("no-cleanup", false)
	viper.SetDefault("yes", false)
	viper.SetDefault("testdefinition", "testdefinition.yaml")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// loadConfig load the given config file
func loadConfig() error {
	cfgFile := viper.GetString("testdefinition")
	if cfgFile == "" {
		return fmt.Errorf("empty testdefinition flag given")
	}

	file, err := os.Open(cfgFile)
	if err != nil {
		return err
	}

	cfgContent, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	cfg = config.New()

	if err := yaml.Unmarshal(cfgContent, cfg); err != nil {
		return err
	}

	return nil
}

func run(cmd *cobra.Command, args []string) error {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetReportCaller(false)
	log.SetLevel(log.DebugLevel)

	if err := loadConfig(); err != nil {
		return err
	}

	var runner runners.Runner
	runnerName := strings.ToLower(cfg.Runner.Name)
	runnerNewFunc, ok := runners.Factories[runnerName]
	if !ok {
		return fmt.Errorf("runner with name %s not found", runnerName)
	}
	runner, err := runnerNewFunc(cfg)
	if err != nil {
		return err
	}

	var testErrors []error

	for i, test := range cfg.Tests {
		log.Infof("running task %d of %d", i+1, len(cfg.Tests))

		tester, parser, outputsAssembled, err := prepare(test)
		if err != nil {
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
		plan, err := tester.Plan(env, test)
		if err != nil {
			return err
		}
		// Set TestStartTime for usage in output / results later on
		plan.TestStartTime = time.Now()

		if viper.GetBool("print-plan") {
			// Pretty print the plan of the test to the shell
			fmt.Println("--> BEGIN PLAN")
			plan.PrettyPrint()
			fmt.Println("--> END PLAN")
		}

		if !viper.GetBool("yes") {
			// Ask user if we can continue or not
			if err := askUserForYes(); err != nil {
				return err
			}
		}

		log.WithFields(logrus.Fields{"tester": test.Type}).Info("preparing test")

		// Prepare the runner for the plan
		if err = runner.Prepare(test.RunOptions, plan); err != nil {
			return err
		}

		var wg sync.WaitGroup

		doneCh := make(chan struct{})
		inCh := make(chan parsers.Input)
		dataCh := make(chan outputs.Data)
		errCh := make(chan error)

		go func() {
			select {
			case erro := <-errCh:
				testErrors = append(testErrors, erro)
				log.Error(erro.Error())
			case <-doneCh:
				return
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			// Close dataCh as there won't be anything else coming through
			defer close(dataCh)
			if err := parser.Parse(doneCh, inCh, dataCh); err != nil {
				log.Error(err)
				errCh <- err
				return
			}
		}()

		// Start each output
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := doOutputs(outputsAssembled, test, doneCh, dataCh); err != nil {
				log.Error(err)
				errCh <- err
				return
			}
		}()

		log.WithFields(logrus.Fields{"tester": test.Type}).Info("executing test")

		// Execute the plan
		if err := runner.Execute(plan, inCh); err != nil {
			log.Error(err)
			errCh <- err
		}
		log.WithFields(logrus.Fields{"runner": runnerName}).Debug("runner execute returned, closing inCh and wg.Wait()")

		close(inCh)

		// TODO Do error checking here
		if err := checkForErrors(plan); err != nil {
			if !test.RunOptions.ContinueOnError {
				log.Error(err)
				return err
			}
			log.Warnf("continuing after err. %+v", err)
		}

		wg.Wait()

		close(doneCh)

		// Run runners.Cleanup() func if wanted by the user
		// TODO When wanted, should be run when signal (CTRL+C) received
		if !viper.GetBool("no-cleanup") {
			if err := runnerCleanup(runner, plan); err != nil {
				return err
			}
		}
	}

	log.Info("done with tests")

	return nil
}

func askUserForYes() error {
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

	return nil
}

func prepare(test *config.Test) (testers.Tester, parsers.Parser, map[string]outputs.Output, error) {
	var tester testers.Tester
	var parser parsers.Parser
	outputsAssembled := map[string]outputs.Output{}

	// Get tester for the test
	testerName := strings.ToLower(test.Type)
	testerNewFunc, ok := testers.Factories[testerName]
	if !ok {
		return nil, nil, nil, fmt.Errorf("tester with name %s not found", testerName)
	}
	var err error
	if tester, err = testerNewFunc(cfg, test); err != nil {
		return nil, nil, nil, err
	}

	// Get parser for the tester output
	parserNewFunc, ok := parsers.Factories[testerName]
	if !ok {
		return nil, nil, nil, fmt.Errorf("parser with name %s not found", testerName)
	}
	if parser, err = parserNewFunc(cfg, test); err != nil {
		return nil, nil, nil, err
	}

	for _, outputItem := range test.Outputs {
		outputName := outputItem.Name

		outputNewFunc, ok := outputs.Factories[outputName]
		if !ok {
			return nil, nil, nil, fmt.Errorf("output with name %s not found", outputName)
		}
		var err error
		outputsAssembled[outputName], err = outputNewFunc(cfg, &outputItem)
		if err != nil {
			return nil, nil, nil, err
		}
	}

	return tester, parser, outputsAssembled, err
}

func doOutputs(outputsAssembled map[string]outputs.Output, test *config.Test, doneCh chan struct{}, dataCh chan outputs.Data) error {
	for {
		select {
		case data, ok := <-dataCh:
			if !ok {
				log.Debug("dataCh closed, in doOutputs()")
				return nil
			}
			for _, outputItem := range test.Outputs {
				outputName := outputItem.Name

				if err := outputsAssembled[outputName].Do(data); err != nil {
					return fmt.Errorf("error in output Do() func. %+v", err)
				}
			}
		case <-doneCh:
			return nil
		}
	}
}

func checkForErrors(plan *testers.Plan) error {
	errorOccured := false

	for _, command := range plan.Commands {
		for _, task := range command {
			// FailedHosts
			if len(task.Status.FailedHosts.Servers) > 0 {
				errorOccured = true
				fmt.Println("-> Failed Server Hosts")
				for host, count := range task.Status.FailedHosts.Servers {
					fmt.Printf("%s - %d\n", host, count)
					for _, err := range task.Status.Errors[host] {
						fmt.Print(err)
					}
				}
				fmt.Println("=> Failed Server Hosts")
			}
			if len(task.Status.FailedHosts.Clients) > 0 {
				errorOccured = true
				fmt.Println("-> Failed Client Hosts")
				for host, count := range task.Status.FailedHosts.Clients {
					fmt.Printf("%s - %d\n", host, count)
					for _, err := range task.Status.Errors[host] {
						fmt.Print(err)
					}
				}
				fmt.Println("=> Failed Client Hosts")
			}

			// SuccessfulHosts
			if len(task.Status.SuccessfulHosts.Servers) > 0 {
				errorOccured = true
				fmt.Println("-> Successful Server Hosts")
				for host, count := range task.Status.SuccessfulHosts.Servers {
					fmt.Printf("%s - %d\n", host, count)
					for _, err := range task.Status.Errors[host] {
						fmt.Print(err)
					}
				}
				fmt.Println("=> Successful Server Hosts")
			}
			if len(task.Status.SuccessfulHosts.Clients) > 0 {
				errorOccured = true
				fmt.Println("-> Successful Client Hosts")
				for host, count := range task.Status.SuccessfulHosts.Clients {
					fmt.Printf("%s - %d\n", host, count)
					for _, err := range task.Status.Errors[host] {
						fmt.Print(err)
					}
				}
				fmt.Println("=> Successful Client Hosts")
			}
		}
	}

	if errorOccured {
		return fmt.Errorf("errors occured during task")
	}

	return nil
}

func runnerCleanup(runner runners.Runner, plan *testers.Plan) error {
	log.Info("running runner cleanup func for test")

	if err := runner.Cleanup(plan); err != nil {
		return err
	}

	return nil
}
