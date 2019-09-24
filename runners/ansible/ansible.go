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

package ansible

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/cloudical-io/ancientt/parsers"
	"github.com/cloudical-io/ancientt/pkg/ansible"
	"github.com/cloudical-io/ancientt/pkg/cmdtemplate"
	"github.com/cloudical-io/ancientt/pkg/config"
	"github.com/cloudical-io/ancientt/pkg/executor"
	"github.com/cloudical-io/ancientt/pkg/util"
	"github.com/cloudical-io/ancientt/runners"
	"github.com/cloudical-io/ancientt/testers"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

const (
	// Name Ansible Runner Name
	Name = "ansible"
	// AnsibleCommand `ansible` command
	AnsibleCommand = "ansible"
	// AnsibleInventoryCommand `ansible-inventory` command
	AnsibleInventoryCommand = "ansible-inventory"
)

var (
	jsonHeadCleanRegex = regexp.MustCompile(`(?sm)^.*(=> \{| >>$\n\{)`)
	jsonTailCleanRegex = regexp.MustCompile(`(?sm)(^\}.*)`)
)

func init() {
	runners.Factories[Name] = NewRunner
}

// Ansible Ansible runner struct
type Ansible struct {
	runners.Runner
	logger         *log.Entry
	config         *config.RunnerAnsible
	runOptions     config.RunOptions
	executor       executor.Executor
	additionalInfo string
}

// NewRunner return a new Ansible Runner
func NewRunner(cfg *config.Config) (runners.Runner, error) {
	if cfg.Runner.Ansible == nil {
		return nil, fmt.Errorf("no ansible runner config")
	}

	if cfg.Runner.Ansible.AnsibleCommand == "" {
		var err error
		cfg.Runner.Ansible.AnsibleCommand, err = exec.LookPath("ansible")
		if err != nil {
			return nil, err
		}
	}

	if cfg.Runner.Ansible.AnsibleInventoryCommand == "" {
		var err error
		cfg.Runner.Ansible.AnsibleInventoryCommand, err = exec.LookPath("ansible-inventory")
		if err != nil {
			return nil, err
		}
	}

	if cfg.Runner.Ansible.CommandTimeout == 0 {
		cfg.Runner.Ansible.CommandTimeout = 20 * time.Second
	}

	if cfg.Runner.Ansible.TaskCommandTimeout == 0 {
		cfg.Runner.Ansible.TaskCommandTimeout = 45 * time.Second
	}

	if cfg.Runner.Ansible.InventoryFilePath == "" {
		return nil, fmt.Errorf("no inventory file path given")
	}
	if cfg.Runner.Ansible.Groups == nil {
		cfg.Runner.Ansible.Groups = &config.AnsibleGroups{}
	}

	if cfg.Runner.Ansible.Groups.Clients == "" {
		cfg.Runner.Ansible.Groups.Clients = "clients"
	} else if cfg.Runner.Ansible.Groups.Clients == "_meta" {
		return nil, fmt.Errorf("ansible clients group can't be named `_meta`")
	}
	if cfg.Runner.Ansible.Groups.Server == "" {
		cfg.Runner.Ansible.Groups.Server = "server"
	} else if cfg.Runner.Ansible.Groups.Server == "_meta" {
		return nil, fmt.Errorf("ansible server group can't be named `_meta`")
	}

	return &Ansible{
		logger:   log.WithFields(logrus.Fields{"runner": Name, "inventoryfile": cfg.Runner.Ansible.InventoryFilePath}),
		config:   cfg.Runner.Ansible,
		executor: executor.NewCommandExecutor("runner:ansible"),
	}, nil
}

// GetHostsForTest return a mocked list of hots for the given test config
func (a *Ansible) GetHostsForTest(test *config.Test) (*testers.Hosts, error) {
	ctx, cancel := context.WithTimeout(context.Background(), a.config.CommandTimeout)
	defer cancel()

	out, err := a.executor.ExecuteCommandWithOutputByte(ctx, "runner:ansible: list hosts from inventory", a.config.AnsibleInventoryCommand, []string{
		fmt.Sprintf("--inventory=%s", a.config.InventoryFilePath),
		"--list",
	}...)
	if err != nil {
		return nil, err
	}

	inv, err := ansible.Parse(out)
	if err != nil {
		return nil, err
	}

	clients := map[string]*testers.Host{}
	for _, client := range inv.GetHostsForGroup(a.config.Groups.Clients) {
		addresses, err := a.getHostNetworkAddress(client)
		if err != nil {
			return nil, err
		}

		clients[client] = &testers.Host{
			Name:      client,
			Labels:    map[string]string{},
			Addresses: addresses,
		}
	}

	servers := map[string]*testers.Host{}
	for _, server := range inv.GetHostsForGroup(a.config.Groups.Server) {
		addresses, err := a.getHostNetworkAddress(server)
		if err != nil {
			return nil, err
		}

		servers[server] = &testers.Host{
			Name:      server,
			Labels:    map[string]string{},
			Addresses: addresses,
		}
	}

	hosts := &testers.Hosts{
		Clients: clients,
		Servers: servers,
	}

	return hosts, nil
}

/*
ansible_default_ipv4.interface and ansible_default_ipv6.interface
```
{
    "ansible_facts": {
    [...]
    "ansible_default_ipv4": {
        "address": "172.16.5.100",
        [...]
    },
    "ansible_default_ipv6": {
        "address": "2a02:8071:22c8:c486:ea5f:3fc7:5039:8b74",
        [...]
    },
[...]
}
```
*/
type facts struct {
	AnsibleFacts networkInterface `json:"ansible_facts"`
}

type networkInterface struct {
	AnsibleDefaultIPv4 networkInterfaceAddress `json:"ansible_default_ipv4"`
	AnsibleDefaultIPv6 networkInterfaceAddress `json:"ansible_default_ipv6"`
}

type networkInterfaceAddress struct {
	Address string `json:"address"`
}

func (a *Ansible) getHostNetworkAddress(host string) (*testers.IPAddresses, error) {
	ctx, cancel := context.WithTimeout(context.Background(), a.config.CommandTimeout)
	defer cancel()

	out, err := a.executor.ExecuteCommandWithOutputByte(ctx, "runner:ansible: list hosts from inventory", a.config.AnsibleCommand, []string{
		fmt.Sprintf("--inventory=%s", a.config.InventoryFilePath),
		host,
		"--module-name=setup",
		"--args=gather_subset=!all,!any,network",
	}...)
	if err != nil {
		return nil, err
	}

	out = cleanAnsibleOutput(out)

	facts := &facts{}
	if err := json.Unmarshal(out, facts); err != nil {
		return nil, err
	}

	addresses := &testers.IPAddresses{}

	if facts.AnsibleFacts.AnsibleDefaultIPv4.Address != "" {
		addresses.IPv4 = []string{
			facts.AnsibleFacts.AnsibleDefaultIPv4.Address,
		}
	}
	if facts.AnsibleFacts.AnsibleDefaultIPv6.Address != "" {
		addresses.IPv6 = []string{
			facts.AnsibleFacts.AnsibleDefaultIPv6.Address,
		}
	}

	if facts.AnsibleFacts.AnsibleDefaultIPv4.Address == "" && facts.AnsibleFacts.AnsibleDefaultIPv6.Address == "" {
		return nil, fmt.Errorf("no default IP addresses for ansible host %s", host)
	}

	return addresses, nil
}

// Prepare prepare Ansible runner for usage, though right now there isn't really anything in need of preparations
func (a *Ansible) Prepare(runOpts config.RunOptions, plan *testers.Plan) error {
	a.runOptions = runOpts

	ctx, cancel := context.WithTimeout(context.Background(), a.config.CommandTimeout)
	defer cancel()

	out, err := a.executor.ExecuteCommandWithOutput(ctx, "runner:ansible: get ansible version", a.config.AnsibleCommand, "--version")
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(strings.NewReader(out))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "ansible ") {
			a.additionalInfo = line
			break
		}
	}

	return nil
}

// Execute run the given commands and return the logs of it and / or error
func (a *Ansible) Execute(plan *testers.Plan, parser chan<- parsers.Input) error {
	for round, tasks := range plan.Commands {
		a.logger.Infof("running commands round %d of %d", round+1, len(plan.Commands))
		for i, task := range tasks {
			if task.Sleep != 0 {
				a.logger.Infof("waiting %s to pass before continuing next round", task.Sleep.String())
				time.Sleep(task.Sleep)
				continue
			}
			a.logger.Infof("running task round %d of %d", i+1, len(tasks))

			if err := a.runTasks(round, task, plan.TestStartTime, plan.Tester, util.GetTaskName(plan), parser); err != nil {
				if !plan.RunOptions.ContinueOnError {
					return err
				}
				a.logger.Warnf("continuing after err. %+v", err)
			}
		}
	}

	return nil
}

func (a *Ansible) runTasks(round int, mainTask *testers.Task, plannedTime time.Time, tester string, taskName string, parser chan<- parsers.Input) error {
	logger := a.logger.WithFields(logrus.Fields{"round": round})

	// Create initial cmdtemplate.Variables
	templateVars := cmdtemplate.Variables{
		ServerPort: 5601,
	}
	if len(mainTask.Host.Addresses.IPv4) > 0 {
		templateVars.ServerAddressV4 = mainTask.Host.Addresses.IPv4[0]
	}
	if len(mainTask.Host.Addresses.IPv6) > 0 {
		templateVars.ServerAddressV6 = mainTask.Host.Addresses.IPv6[0]
	}

	if err := cmdtemplate.Template(mainTask, templateVars); err != nil {
		erro := fmt.Errorf("failed to template main task command and / or args. %+v", err)
		logger.Error(erro)
		mainTask.Status.AddFailedServer(mainTask.Host, erro)
		return erro
	}

	var mainWG sync.WaitGroup
	var wg sync.WaitGroup

	mainTaskStopped := false
	mainCtx, mainCancel := context.WithCancel(context.Background())
	defer mainCancel()

	mainWG.Add(1)
	go func() {
		defer mainWG.Done()
		err := a.executor.ExecuteCommand(mainCtx, "runner:ansible: run main task command", a.config.AnsibleCommand, []string{
			fmt.Sprintf("--inventory=%s", a.config.InventoryFilePath),
			mainTask.Host.Name,
			"--module-name=shell",
			fmt.Sprintf("--args=%s %s", mainTask.Command, strings.Join(mainTask.Args, " ")),
		}...)
		if err != nil {
			if exiterr, ok := err.(*exec.ExitError); ok {
				fmt.Printf("EXITERR: %+v - %+v - %+v\n", exiterr, exiterr.Pid(), exiterr.ProcessState)

				if err := syscall.Kill(-exiterr.Pid(), syscall.SIGKILL); err != nil {
					log.Println("failed to kill: ", err)
				}
			}
			// Ignore any error after the main task is stopped
			if mainTaskStopped {
				logger.Debug(err)
				return
			}

			logger.Error(err)
			mainTask.Status.AddFailedServer(mainTask.Host, err)
			return
		}
	}()

	time.Sleep(250 * time.Millisecond)

	ready := false
	checkCtx, checkCancel := context.WithTimeout(context.Background(), a.config.TaskCommandTimeout)
	defer checkCancel()

	tries := 5
	for i := 0; i <= tries; i++ {
		err := a.executor.ExecuteCommand(checkCtx, fmt.Sprintf("runner:ansible: check if main task is running (try: %d/%d)", i, tries), a.config.AnsibleCommand, []string{
			fmt.Sprintf("--inventory=%s", a.config.InventoryFilePath),
			mainTask.Host.Name,
			"--module-name=shell",
			fmt.Sprintf("--args=pgrep %s", mainTask.Command),
		}...)
		if err == nil {
			ready = true
			break
		}
		logger.Error(err)

		logger.Infof("main task not running yet, sleeping 3 seconds (try: %d/%d) ...", i, tries)
		time.Sleep(3 * time.Second)
	}

	if ready {
		for i, task := range mainTask.SubTasks {
			logger.Infof("running sub task %d of %d", i+1, len(mainTask.SubTasks))

			wg.Add(1)
			go func(task *testers.Task) {
				ctx, cancel := context.WithTimeout(context.Background(), a.config.TaskCommandTimeout)
				defer cancel()

				defer wg.Done()

				// Template command and args for each task
				if err := cmdtemplate.Template(task, templateVars); err != nil {
					erro := fmt.Errorf("failed to template task command and / or args. %+v", err)
					logger.Errorf("error during createPodsForTasks. %+v", erro)
					mainTask.Status.AddFailedClient(task.Host, erro)
					return
				}

				testTime := time.Now()

				out, err := a.executor.ExecuteCommandWithOutputByte(ctx, "runner:ansible: run sub task command", a.config.AnsibleCommand, []string{
					fmt.Sprintf("--inventory=%s", a.config.InventoryFilePath),
					task.Host.Name,
					"--module-name=shell",
					fmt.Sprintf("--args=%s %s", task.Command, strings.Join(task.Args, " ")),
				}...)
				if err != nil {
					logger.Error(err)
					mainTask.Status.AddFailedClient(task.Host, err)
					return
				}

				mainTask.Status.AddSuccessfulClient(task.Host)

				// Clean, "transform" to io.Reader compatible interface and send logs to parsers
				out = cleanAnsibleOutput(out)
				r := ioutil.NopCloser(bytes.NewReader(out))

				parser <- parsers.Input{
					TestStartTime:  plannedTime,
					TestTime:       testTime,
					Round:          round,
					DataStream:     &r,
					Tester:         tester,
					ServerHost:     mainTask.Host.Name,
					ClientHost:     task.Host.Name,
					AdditionalInfo: a.additionalInfo,
				}
			}(task)

			if a.runOptions.Mode != config.RunModeParallel {
				wg.Wait()
			}
		}

		// When RunOptions.Mode `parallel` then we wait after all test tasks have been run
		if a.runOptions.Mode == config.RunModeParallel {
			wg.Wait()
		}

		mainTask.Status.AddSuccessfulServer(mainTask.Host)
		mainTaskStopped = true
	} else {
		err := fmt.Errorf("ansible main test task is not running")
		mainTask.Status.AddFailedServer(mainTask.Host, err)
		return err
	}

	logger.Info("stopping main task")
	mainCancel()
	mainWG.Wait()

	logger.Debug("done running tasks for test in ansible for plan")

	return nil
}

// Cleanup remove all (left behind) Ansible resources created for the given Plan.
func (a *Ansible) Cleanup(plan *testers.Plan) error {
	// Nothing to do here for Ansible (yet)
	return nil
}

func cleanAnsibleOutput(in []byte) []byte {
	return jsonTailCleanRegex.ReplaceAll(
		jsonHeadCleanRegex.ReplaceAll(in, []byte("{")),
		[]byte("}"))
}
