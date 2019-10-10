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

package kubernetes

import (
	"fmt"
	"sync"
	"time"

	"github.com/cloudical-io/ancientt/parsers"
	"github.com/cloudical-io/ancientt/pkg/cmdtemplate"
	"github.com/cloudical-io/ancientt/pkg/config"
	"github.com/cloudical-io/ancientt/pkg/hostsfilter"
	"github.com/cloudical-io/ancientt/pkg/k8sutil"
	"github.com/cloudical-io/ancientt/pkg/util"
	"github.com/cloudical-io/ancientt/runners"
	"github.com/cloudical-io/ancientt/testers"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// Name Kubernetes Runner Name
const Name = "kubernetes"

func init() {
	runners.Factories[Name] = NewRunner
}

// Kubernetes Kubernetes runner struct
type Kubernetes struct {
	runners.Runner
	logger     *log.Entry
	config     *config.RunnerKubernetes
	k8sclient  kubernetes.Interface
	runOptions config.RunOptions
}

// NewRunner return a new Kubernetes Runner
func NewRunner(cfg *config.Config) (runners.Runner, error) {
	conf := cfg.Runner.Kubernetes

	clientset, err := k8sutil.NewClient(cfg.Runner.Kubernetes.InClusterConfig, cfg.Runner.Kubernetes.Kubeconfig)
	if err != nil {
		return nil, err
	}

	return &Kubernetes{
		logger:    log.WithFields(logrus.Fields{"runner": Name, "namespace": cfg.Runner.Kubernetes.Namespace}),
		config:    conf,
		k8sclient: clientset,
	}, nil
}

// GetHostsForTest return a mocked list of hots for the given test config
func (k *Kubernetes) GetHostsForTest(test *config.Test) (*testers.Hosts, error) {
	hosts := &testers.Hosts{
		Clients: map[string]*testers.Host{},
		Servers: map[string]*testers.Host{},
	}

	k8sNodes, err := k.k8sNodesToHosts()
	if err != nil {
		return nil, err
	}

	// Go through Hosts Servers list to get the servers hosts
	for _, servers := range test.Hosts.Servers {
		filtered, err := hostsfilter.FilterHostsList(k8sNodes, servers)
		if err != nil {
			return nil, err
		}
		for _, host := range filtered {
			if _, ok := hosts.Servers[host.Name]; !ok {
				hosts.Servers[host.Name] = host
			}
		}
	}

	// Go through Hosts Clients list to get the clients hosts
	for _, clients := range test.Hosts.Clients {
		filtered, err := hostsfilter.FilterHostsList(k8sNodes, clients)
		if err != nil {
			return nil, err
		}
		for _, host := range filtered {
			if _, ok := hosts.Clients[host.Name]; !ok {
				hosts.Clients[host.Name] = host
			}
		}
	}

	k.logger.Debug("returning Kubernetes hosts list")

	return hosts, nil
}

func (k *Kubernetes) k8sNodesToHosts() ([]*testers.Host, error) {
	hosts := []*testers.Host{}
	nodes, err := k.k8sclient.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Quick conversion from a Kubernetes CoreV1 Nodes object to testers.Host
	for _, node := range nodes.Items {
		// Check if node is unschedulable
		if (k.config.Hosts != nil && *k.config.Hosts.IgnoreSchedulingDisabled) && node.Spec.Unschedulable {
			k.logger.WithFields(logrus.Fields{"node": node.ObjectMeta.Name}).Debug("skipping unschedulable node")
			continue
		}

		// Check if the taints on the node match the given tolerations
		if !k8sutil.NodeIsTolerable(node, k.config.Hosts.Tolerations) {
			continue
		}

		hosts = append(hosts, &testers.Host{
			Labels: node.ObjectMeta.Labels,
			Name:   node.ObjectMeta.Name,
		})
	}

	return hosts, nil
}

// Prepare prepare Kubernetes for usage with ancientt, e.g., create Namespace.
func (k *Kubernetes) Prepare(runOpts config.RunOptions, plan *testers.Plan) error {
	k.runOptions = runOpts

	if err := k.prepareKubernetes(); err != nil {
		return err
	}

	return nil
}

// Execute run the given commands and return the logs of it and / or error
func (k *Kubernetes) Execute(plan *testers.Plan, parser chan<- parsers.Input) error {
	// TODO Add option to go through Service IPs instead of Pod IPs

	// Iterate over given plan.Commands to then run each task
	for round, tasks := range plan.Commands {
		k.logger.Infof("running commands round %d of %d", round+1, len(plan.Commands))
		for i, task := range tasks {
			if task.Sleep != 0 {
				k.logger.Infof("waiting %s to pass before continuing next round", task.Sleep.String())
				time.Sleep(task.Sleep)
				continue
			}
			k.logger.Infof("running task round %d of %d", i+1, len(tasks))

			// Create the Pods for the server task and client tasks
			if err := k.createPodsForTasks(round, task, plan.TestStartTime, plan.Tester, util.GetTaskName(plan.Tester, plan.TestStartTime), parser); err != nil {
				if !*plan.RunOptions.ContinueOnError {
					return err
				}
				k.logger.Warnf("continuing after err. %+v", err)
			}
		}
	}

	return nil
}

// prepareKubernetes prepares Kubernetes by creating the namespace if it does not exist
func (k *Kubernetes) prepareKubernetes() error {
	// Check if namespaces exists, if not try create it
	if _, err := k.k8sclient.CoreV1().Namespaces().Get(k.config.Namespace, metav1.GetOptions{}); err != nil {
		// If namespace not found, create it
		if errors.IsNotFound(err) {
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"created-by": "ancientt",
					},
					Name: k.config.Namespace,
				},
			}
			k.logger.Info("trying to create namespace")
			if _, err := k.k8sclient.CoreV1().Namespaces().Create(ns); err != nil {
				return fmt.Errorf("failed to create namespace %s. %+v", k.config.Namespace, err)
			}
			k.logger.Info("created namespace")
		} else {
			return fmt.Errorf("error while getting namespace %s. %+v", k.config.Namespace, err)
		}
	}
	return nil
}

// createPodsForTasks create the Pods that are needed for the task(s)
func (k *Kubernetes) createPodsForTasks(round int, mainTask *testers.Task, plannedTime time.Time, tester string, taskName string, parser chan<- parsers.Input) error {
	logger := k.logger.WithFields(logrus.Fields{"round": round})

	var wg sync.WaitGroup

	// Create server Pod first
	serverPodName := util.GetPNameFromTask(round, mainTask.Host.Name, mainTask.Command, mainTask.Args, util.PNameRoleServer)

	// Create initial cmdtemplate.Variables
	templateVars := cmdtemplate.Variables{
		ServerPort: 5601,
	}

	if err := cmdtemplate.Template(mainTask, templateVars); err != nil {
		erro := fmt.Errorf("failed to template main task command and / or args. %+v", err)
		k.logger.Error(erro)
		mainTask.Status.AddFailedServer(mainTask.Host, erro)
		return nil
	}

	pod := k.getPodSpec(serverPodName, taskName, mainTask)
	k.applyServiceAccountToPod(pod, serverRole)

	logger.WithFields(logrus.Fields{"pod": serverPodName}).Debug("(re)creating server pod")
	if err := k8sutil.PodRecreate(k.k8sclient, pod, k.config.Timeouts.DeleteTimeout); err != nil {
		erro := fmt.Errorf("failed to create server pod %s/%s. %+v", k.config.Namespace, serverPodName, err)
		k.logger.Error(erro)
		mainTask.Status.AddFailedServer(mainTask.Host, erro)
		return nil
	}

	logger.WithFields(logrus.Fields{"pod": serverPodName}).Info("waiting for server pod to run")
	running, err := k8sutil.WaitForPodToRun(k.k8sclient, k.config.Namespace, serverPodName, k.config.Timeouts.RunningTimeout)
	if err != nil {
		erro := fmt.Errorf("failed to wait for server pod %s/%s. %+v", k.config.Namespace, serverPodName, err)
		k.logger.Error(erro)
		mainTask.Status.AddFailedServer(mainTask.Host, erro)
		return nil
	}
	if !running {
		erro := fmt.Errorf("server pod %s/%s not running after runTimeout", k.config.Namespace, serverPodName)
		k.logger.Error(erro)
		mainTask.Status.AddFailedServer(mainTask.Host, erro)
		return nil
	}

	// Get server Pod to have the server IP for each client task
	pod, err = k.k8sclient.CoreV1().Pods(k.config.Namespace).Get(serverPodName, metav1.GetOptions{})
	if err != nil {
		erro := fmt.Errorf("failed to get server pod %s/%s. %+v", k.config.Namespace, serverPodName, err)
		k.logger.Error(erro)
		mainTask.Status.AddFailedServer(mainTask.Host, erro)
		return nil
	}
	if pod.Status.PodIP == "" {
		mainTask.Status.AddFailedServer(mainTask.Host,
			fmt.Errorf("failed to get server pod %s/%s IP, got '%s'", k.config.Namespace, serverPodName, pod.Status.PodIP))
		return nil
	}

	templateVars.ServerAddressV4 = pod.Status.PodIP

	for i, task := range mainTask.SubTasks {
		k.logger.Infof("running sub task %d of %d", i+1, len(mainTask.SubTasks))

		wg.Add(1)
		go func(task *testers.Task) {
			defer wg.Done()

			testTime := time.Now()

			pName := util.GetPNameFromTask(round, task.Host.Name, task.Command, task.Args, util.PNameRoleClient)

			// Template command and args for each task
			if err := cmdtemplate.Template(task, templateVars); err != nil {
				erro := fmt.Errorf("failed to template task command and / or args. %+v", err)
				logger.Errorf("error during createPodsForTasks. %+v", erro)
				mainTask.Status.AddFailedClient(task.Host, erro)
				return
			}

			pod = k.getPodSpec(pName, taskName, task)
			k.applyServiceAccountToPod(pod, clientsRole)

			logger.WithFields(logrus.Fields{"pod": pName}).Debug("(re)creating client pod")
			if err := k8sutil.PodRecreate(k.k8sclient, pod, k.config.Timeouts.DeleteTimeout); err != nil {
				erro := fmt.Errorf("failed to create pod %s/%s. %+v", k.config.Namespace, pName, err)
				logger.Errorf("error during createPodsForTasks. %+v", erro)
				mainTask.Status.AddFailedClient(task.Host, erro)
				return
			}

			logger.WithFields(logrus.Fields{"pod": pName}).Info("waiting for client pod to run or succeed")
			running, err := k8sutil.WaitForPodToRunOrSucceed(k.k8sclient, k.config.Namespace, pName, k.config.Timeouts.RunningTimeout)
			if err != nil {
				erro := fmt.Errorf("failed to wait for pod %s/%s. %+v", k.config.Namespace, pName, err)
				logger.Errorf("error during createPodsForTasks. %+v", erro)
				mainTask.Status.AddFailedClient(task.Host, erro)
				return
			}
			if !running {
				erro := fmt.Errorf("pod %s/%s not running after runTimeout", k.config.Namespace, pName)
				logger.Errorf("error during createPodsForTasks. %+v", erro)
				mainTask.Status.AddFailedClient(task.Host, erro)
				return
			}

			logger.WithFields(logrus.Fields{"pod": pName}).Debug("about to pushLogsToParser")
			if err := k.pushLogsToParser(parser, plannedTime, testTime, round, tester, mainTask.Host.Name, task.Host.Name, pName); err != nil {
				erro := fmt.Errorf("failed to push pod %s/%s logs to parser. %+v", k.config.Namespace, pName, err)
				logger.Errorf("error during createPodsForTasks. %+v", erro)
				mainTask.Status.AddFailedClient(task.Host, erro)
				return
			}

			logger.WithFields(logrus.Fields{"pod": pName}).Info("deleting client pod")
			if err := k8sutil.PodDelete(k.k8sclient, pod, k.config.Timeouts.DeleteTimeout); err != nil {
				erro := fmt.Errorf("failed to delete client pod %s/%s. %+v", k.config.Namespace, pName, err)
				logger.Errorf("error during createPodsForTasks. %+v", erro)
				mainTask.Status.AddFailedClient(task.Host, erro)
				return
			}

			mainTask.Status.AddSuccessfulClient(task.Host)
		}(task)

		if k.runOptions.Mode != config.RunModeParallel {
			wg.Wait()
		}
	}

	// When RunOptions.Mode `parallel` then we wait after all test tasks have been run
	if k.runOptions.Mode == config.RunModeParallel {
		wg.Wait()
	}

	// Delete server pod
	logger.WithFields(logrus.Fields{"pod": serverPodName}).Info("deleting server pod")
	if err := k8sutil.PodDeleteByName(k.k8sclient, k.config.Namespace, serverPodName, k.config.Timeouts.DeleteTimeout); err != nil {
		erro := fmt.Errorf("failed to delete server pod. %+v", err)
		logger.WithFields(logrus.Fields{"pod": serverPodName}).Error(erro)
		mainTask.Status.AddFailedServer(mainTask.Host, erro)
		return nil
	}

	mainTask.Status.AddSuccessfulServer(mainTask.Host)

	logger.Debug("done running tasks for test in kubernetes for plan")

	return nil
}

func (k *Kubernetes) pushLogsToParser(parserInput chan<- parsers.Input, plannedTime time.Time, testTime time.Time, round int, tester string, serverHost string, clientHost string, podName string) error {
	// Wait for the Pod to succeed because that is the "sign" that the test for that Pod is done.
	succeeded, err := k8sutil.WaitForPodToSucceed(k.k8sclient, k.config.Namespace, podName, k.config.Timeouts.SucceedTimeout)
	if err != nil {
		return err
	}

	if succeeded {
		// "Generate" request for logs of Pod
		req := k.k8sclient.CoreV1().Pods(k.config.Namespace).GetLogs(podName, &corev1.PodLogOptions{})

		// Start the log stream
		podLogs, err := req.Stream()
		if err != nil {
			return err
		}
		// Don't close the `podLogs` here, that is the responsibility of the parser!

		// Send the logs to the parser.InputChan
		parserInput <- parsers.Input{
			TestStartTime:  plannedTime,
			TestTime:       testTime,
			Round:          round,
			DataStream:     &podLogs,
			Tester:         tester,
			ServerHost:     serverHost,
			ClientHost:     clientHost,
			AdditionalInfo: podName,
		}
		return nil
	}

	return fmt.Errorf("pod %s/%s has not succeeded", k.config.Namespace, podName)
}

// Cleanup remove all (left behind) Kubernetes resources created for the given Plan.
func (k *Kubernetes) Cleanup(plan *testers.Plan) error {
	var wg sync.WaitGroup

	// Delete all Pods with label XYZ
	if err := k8sutil.PodDeleteByLabels(k.k8sclient, k.config.Namespace, map[string]string{
		k8sutil.TaskIDLabel: util.GetTaskName(plan.Tester, plan.TestStartTime),
	}); err != nil {
		k.logger.Errorf("error during pod delete by labels in cleanup. %+v", err)
		return err
	}
	wg.Wait()

	return nil
}
