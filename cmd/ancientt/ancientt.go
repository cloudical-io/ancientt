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
	"os"
	"strings"
	"sync"
	"time"

	"github.com/cloudical-io/ancientt/outputs"
	"github.com/cloudical-io/ancientt/parsers"
	"github.com/cloudical-io/ancientt/pkg/config"
	"github.com/cloudical-io/ancientt/runners"
	"github.com/cloudical-io/ancientt/testers"
	au "github.com/logrusorgru/aurora"
	"github.com/mattn/go-isatty"
	"github.com/prometheus/common/version"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	outputSeparator = aurora.Red("===================")
	aurora          = au.NewAurora(isatty.IsTerminal(os.Stdout.Fd()))

	rootCmd = &cobra.Command{
		Use:   "ancientt",
		Short: "Ancientt is a tool to automate network testing tools, like iperf3, in dynamic environments such as Kubernetes and more to come dynamic environments.",
		RunE:  run,
	}
	cfg      *config.Config
	logLevel string
)

func init() {
	// Set flags, viper binds for flags and viper bind default values
	rootCmd.PersistentFlags().Bool("version", false, "Print version info and exit.")
	rootCmd.PersistentFlags().BoolP("only-print-plan", "p", false, "Only print plan for the testdefinitions to console and exit.")
	rootCmd.PersistentFlags().Bool("no-cleanup", false, "If runners should run cleanup routines after the tests.")
	rootCmd.PersistentFlags().BoolP("yes", "y", false, "Ask for user confirmation for each test before executing it.")
	rootCmd.PersistentFlags().StringP("testdefinition", "c", "testdefinition.yaml", "Path to the testdefinitions to read for the tests.")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "INFO", "Log level (DEBUG, INFO, WARN, ERROR, default: INFO).")
	viper.BindPFlag("version", rootCmd.PersistentFlags().Lookup("version"))
	viper.BindPFlag("only-print-plan", rootCmd.PersistentFlags().Lookup("only-print-plan"))
	viper.BindPFlag("no-cleanup", rootCmd.PersistentFlags().Lookup("no-cleanup"))
	viper.BindPFlag("yes", rootCmd.PersistentFlags().Lookup("yes"))
	viper.BindPFlag("testdefinition", rootCmd.PersistentFlags().Lookup("testdefinition"))
	viper.SetDefault("version", false)
	viper.SetDefault("only-print-plan", false)
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

	var err error
	cfg, err = config.Load(cfgFile)

	return err
}

func run(cmd *cobra.Command, args []string) error {
	if viper.GetBool("version") {
		fmt.Print(version.Print(os.Args[0]))
		return nil
	}

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetReportCaller(false)

	log.WithFields(logrus.Fields{
		"program":   os.Args[0],
		"version":   version.Version,
		"branch":    version.Branch,
		"revision":  version.Revision,
		"go":        version.GoVersion,
		"buildUser": version.BuildUser,
		"buildDate": version.BuildDate,
	}).Info("starting ancientt")

	level, err := log.ParseLevel(logLevel)
	if err != nil {
		return fmt.Errorf("failed to parse given log level. %+v", err)
	}
	log.SetLevel(level)

	if err := loadConfig(); err != nil {
		return err
	}

	var runner runners.Runner
	runnerName := strings.ToLower(cfg.Runner.Name)
	runnerNewFunc, ok := runners.Factories[runnerName]
	if !ok {
		return fmt.Errorf("runner with name %s not found", runnerName)
	}
	runner, err = runnerNewFunc(cfg)
	if err != nil {
		return err
	}

	for i, test := range cfg.Tests {
		log.WithFields(logrus.Fields{"runner": runnerName}).Infof("doing test '%s', %d of %d", test.Name, i+1, len(cfg.Tests))

		logger, tester, parser, outputsAssembled, err := prepare(test, runnerName)
		if err != nil {
			logger.Errorf("error preparing test run. %+v", err)
			if !*test.RunOptions.ContinueOnError {
				return err
			}
			logger.Warnf("skippinmg test %d of %d due to error in initial prepare step", i+1, len(cfg.Tests))
			continue
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

		fmt.Println(outputSeparator)
		// Pretty print the plan of the test to the shell
		fmt.Println("--> BEGIN PLAN")
		plan.PrettyPrint()
		fmt.Println("--> END PLAN")
		if viper.GetBool("only-print-plan") {
			return nil
		}

		if !viper.GetBool("yes") {
			// Ask user if we can continue or not
			if err := askUserForYes(); err != nil {
				return err
			}
		}

		logger.Info("preparing test")

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
				logger.Error(erro.Error())
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
				errCh <- err
				return
			}
		}()

		// Start each output
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := doOutputs(outputsAssembled, test, doneCh, dataCh); err != nil {
				errCh <- err
				return
			}
		}()

		logger.Info("executing test")

		// Execute the plan
		if err := runner.Execute(plan, inCh); err != nil {
			errCh <- err
		}
		logger.Debug("runner execute returned, closing inCh and wg.Wait()")

		close(inCh)

		if err := checkForErrors(plan); err != nil {
			logger.Error(err)
			if !*test.RunOptions.ContinueOnError {
				return err
			}
			logger.Warnf("continue on error run option given for test, continuing")
		}

		wg.Wait()

		close(doneCh)

		fmt.Println(outputSeparator)
		fmt.Println(aurora.Magenta("Following files have been created / used:"))
		for outName, output := range outputsAssembled {
			for _, file := range output.OutputFiles() {
				fmt.Printf("%s (output: %s)\n", file, outName)
			}
		}
		fmt.Println(outputSeparator)

		// Run runners.Cleanup() func if wanted by the user
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
	fmt.Println(outputSeparator)

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print(aurora.Underline("Are you sure you want to continue with above test plan? ('Yes' or 'No'):"), " ")
		userInput, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		userInput = strings.TrimSpace(strings.ToLower(userInput))
		fmt.Println(outputSeparator)
		if userInput == "yes" || userInput == "y" {
			break
		} else if userInput == "no" || userInput == "n" {
			return fmt.Errorf("aborted by user")
		}
	}

	return nil
}

func prepare(test *config.Test, runnerName string) (*log.Entry, testers.Tester, parsers.Parser, map[string]outputs.Output, error) {
	var tester testers.Tester
	var parser parsers.Parser
	outputsAssembled := map[string]outputs.Output{}

	// Get tester for the test
	testerName := strings.ToLower(test.Type)
	logger := log.WithFields(logrus.Fields{"tester": testerName, "parser": testerName, "runner": runnerName})

	testerNewFunc, ok := testers.Factories[testerName]
	if !ok {
		return logger, nil, nil, nil, fmt.Errorf("tester with name %s not found", testerName)
	}
	var err error
	if tester, err = testerNewFunc(cfg, test); err != nil {
		return logger, nil, nil, nil, err
	}

	// Get parser for the tester output
	parserNewFunc, ok := parsers.Factories[testerName]
	if !ok {
		return logger, nil, nil, nil, fmt.Errorf("parser with name %s not found", testerName)
	}
	if parser, err = parserNewFunc(cfg, test); err != nil {
		return logger, nil, nil, nil, err
	}

	for _, outputItem := range test.Outputs {
		outputName := outputItem.Name

		outputNewFunc, ok := outputs.Factories[outputName]
		if !ok {
			return logger, nil, nil, nil, fmt.Errorf("output with name %s not found", outputName)
		}
		var err error
		outputsAssembled[outputName], err = outputNewFunc(cfg, &outputItem)
		if err != nil {
			return logger, nil, nil, nil, err
		}
	}

	return logger, tester, parser, outputsAssembled, err
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
					// TODO Run all ouputs and concat errors
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

	fmt.Println(outputSeparator)
	for _, command := range plan.Commands {
		for _, task := range command {
			if task.Status == nil || task.Sleep != 0 {
				log.Debug("task status is empty or is sleep task, continuing")
				continue
			}
			// FailedHosts
			if len(task.Status.FailedHosts.Servers) > 0 {
				errorOccured = true
				fmt.Println(aurora.Yellow("-> Failed Server Hosts"))
				for host, count := range task.Status.FailedHosts.Servers {
					fmt.Printf("%s - %d\n", host, count)
					for _, err := range task.Status.Errors[host] {
						fmt.Println(err)
					}
				}
			}
			if len(task.Status.FailedHosts.Clients) > 0 {
				errorOccured = true
				fmt.Println(aurora.Yellow("-> Failed Client Hosts"))
				for host, count := range task.Status.FailedHosts.Clients {
					fmt.Printf("%s - %d\n", host, count)
					for _, err := range task.Status.Errors[host] {
						fmt.Println(err)
					}
				}
			}

			// SuccessfulHosts
			if len(task.Status.SuccessfulHosts.Servers) > 0 {
				fmt.Println(aurora.Green("-> Successful Server Hosts"))
				for host, count := range task.Status.SuccessfulHosts.Servers {
					fmt.Printf("%s - %d\n", host, count)
				}
			}
			if len(task.Status.SuccessfulHosts.Clients) > 0 {
				fmt.Println(aurora.Green("-> Successful Client Hosts"))
				for host, count := range task.Status.SuccessfulHosts.Clients {
					fmt.Printf("%s - %d\n", host, count)
				}
			}
		}
	}

	fmt.Println(outputSeparator)

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
