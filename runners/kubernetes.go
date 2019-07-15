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

package runners

import (
	"fmt"
	"sync"
	"time"

	"github.com/cloudical-io/acntt/parsers"
	"github.com/cloudical-io/acntt/pkg/cmdtemplate"
	"github.com/cloudical-io/acntt/pkg/config"
	"github.com/cloudical-io/acntt/pkg/k8sutil"
	"github.com/cloudical-io/acntt/pkg/util"
	"github.com/cloudical-io/acntt/testers"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// NameKubernetes Kubernetes Runner Name
const NameKubernetes = "kubernetes"

func init() {
	Factories[NameKubernetes] = NewKubernetesRunner
}

// Kubernetes Kubernetes runner struct
type Kubernetes struct {
	Runner
	logger     *log.Entry
	config     *config.RunnerKubernetes
	k8sclient  *kubernetes.Clientset
	runOptions config.RunOptions
}

// NewKubernetesRunner return a new Kubernetes Runner
func NewKubernetesRunner(cfg *config.Config) (Runner, error) {
	if cfg.Runner.Kubernetes == nil {
		return nil, fmt.Errorf("no kubernetes runner config")
	}
	// Use the current context in kubeconfig
	k8sconfig, err := clientcmd.BuildConfigFromFlags("", cfg.Runner.Kubernetes.Kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("kubeconfig configuration error. %+v", err)
	}

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(k8sconfig)
	if err != nil {
		return nil, fmt.Errorf("kubernetes client configuration error. %+v", err)
	}

	var k8sConfig *config.RunnerKubernetes
	if cfg.Runner.Kubernetes != nil {
		k8sConfig = cfg.Runner.Kubernetes
	} else {
		k8sConfig = &config.RunnerKubernetes{}
	}

	if k8sConfig.Annotations == nil {
		k8sConfig.Annotations = map[string]string{}
	}

	if k8sConfig.Hosts == nil {
		k8sConfig.Hosts = &config.KubernetesHosts{
			IgnoreSchedulingDisabled: true,
		}
	}

	if k8sConfig.Namespace == "" {
		k8sConfig.Namespace = "acntt"
	}
	if k8sConfig.Timeouts == nil {
		k8sConfig.Timeouts = &config.KubernetesTimeouts{
			DeleteTimeout:  20,
			RunningTimeout: 35,
			SucceedTimeout: 60,
		}
	}

	return Kubernetes{
		logger:    log.WithFields(logrus.Fields{"runner": NameKubernetes, "namespace": cfg.Runner.Kubernetes.Namespace}),
		config:    k8sConfig,
		k8sclient: clientset,
	}, nil
}

// GetHostsForTest return a mocked list of hots for the given test config
func (k Kubernetes) GetHostsForTest(test config.Test) (*testers.Hosts, error) {
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
		filtered, err := util.FilterHostsList(k8sNodes, servers)
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
		filtered, err := util.FilterHostsList(k8sNodes, clients)
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

func (k Kubernetes) k8sNodesToHosts() ([]*testers.Host, error) {
	hosts := []*testers.Host{}
	nodes, err := k.k8sclient.CoreV1().Nodes().List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	// Quick conversion from a Kubernetes CoreV1 Nodes object to testers.Host
	for _, node := range nodes.Items {
		// Check if node is unschedulable
		// TODO Add checks for taints (e.g., https://github.com/rook/rook/blob/master/pkg/operator/k8sutil/node.go)
		if k.config.Hosts.IgnoreSchedulingDisabled && node.Spec.Unschedulable {
			k.logger.WithFields(logrus.Fields{"node": node.ObjectMeta.Name}).Debug("skipping unschedulable node")
			continue
		}
		hosts = append(hosts, &testers.Host{
			Labels: node.ObjectMeta.Labels,
			Name:   node.ObjectMeta.Name,
		})
	}

	return hosts, nil
}

// Prepare prepare Kubernetes for usage with acntt, e.g., create Namespace.
func (k Kubernetes) Prepare(runOpts config.RunOptions, plan *testers.Plan) error {
	k.runOptions = runOpts

	if err := k.prepareKubernetes(); err != nil {
		return err
	}

	return nil
}

// Execute run the given commands and return the logs of it and / or error
func (k Kubernetes) Execute(plan *testers.Plan, parser chan<- parsers.Input) error {
	// TODO Add option to go through Service IPs instead of Pod IPs

	// Iterate over given plan.Commands to then run each task
	for round, tasks := range plan.Commands {
		for _, task := range tasks {
			if task.Sleep != 0 {
				k.logger.Infof("waiting %s to pass before continuing next round", task.Sleep.String())
				time.Sleep(task.Sleep)
				continue
			}

			// Create the Pods for the server task and client tasks
			if err := k.createPodsForTasks(round, task, plan.PlannedTime, plan.Tester, util.GetTaskName(plan), parser); err != nil {
				return err
			}
		}
	}

	return nil
}

// prepareKubernetes prepares Kubernetes by creating the namespace if it does not exist
func (k Kubernetes) prepareKubernetes() error {
	// Check if namespaces exists, if not try create it
	if _, err := k.k8sclient.CoreV1().Namespaces().Get(k.config.Namespace, metav1.GetOptions{}); err != nil {
		// If namespace not found, create it
		if errors.IsNotFound(err) {
			ns := &corev1.Namespace{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"created-by": "acntt",
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
func (k Kubernetes) createPodsForTasks(round int, mainTask testers.Task, plannedTime time.Time, tester string, taskName string, parser chan<- parsers.Input) error {
	logger := k.logger.WithFields(logrus.Fields{"round": round})

	var wg sync.WaitGroup
	errs := make(chan error)

	// Create server Pod first
	serverPodName := util.GetPNameFromTask(round, mainTask, util.PNameRoleServer)

	// Create initial cmdtemplate.Variables
	// TODO the port does not need to be mapped in Kubernetes case, but for other runners (e.g., Ansible) need to map the port
	// Find a way to do so, in a good way.
	templateVars := cmdtemplate.Variables{
		ServerPort: 5601,
	}

	if err := cmdtemplate.Template(&mainTask, templateVars); err != nil {
		return fmt.Errorf("failed to template main task command and / or args. %+v", err)
	}

	pod := k.getPodSpec(serverPodName, taskName, mainTask)

	logger.WithFields(logrus.Fields{"pod": serverPodName}).Debug("(re)creating server pod")
	if err := k8sutil.PodRecreate(k.k8sclient, pod, k.config.Timeouts.DeleteTimeout); err != nil {
		return fmt.Errorf("failed to create pod %s/%s. %+v", k.config.Namespace, serverPodName, err)
	}

	logger.WithFields(logrus.Fields{"pod": serverPodName}).Info("waiting for server pod to run")
	running, err := k8sutil.WaitForPodToRun(k.k8sclient, k.config.Namespace, serverPodName, k.config.Timeouts.RunningTimeout)
	if err != nil {
		return fmt.Errorf("failed to wait for pod %s/%s. %+v", k.config.Namespace, serverPodName, err)
	}
	if !running {
		return fmt.Errorf("pod %s/%s not running after runTimeout", k.config.Namespace, serverPodName)
	}

	// Get server Pod to have the server IP for each client task
	pod, err = k.k8sclient.CoreV1().Pods(k.config.Namespace).Get(serverPodName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get pod %s/%s. %+v", k.config.Namespace, serverPodName, err)
	}
	templateVars.ServerAddress = pod.Status.PodIP

	go func() {
		// TODO Fix this code to be.. well "good"(?)
		errsList := []string{}
		for erro := range errs {
			logger.Errorf("error during createPodsForTasks. %+v", erro)
			errsList = append(errsList, erro.Error())
		}
	}()

	tasks := []testers.Task{}
	tasks = append(tasks, mainTask.SubTasks...)
	for _, task := range tasks {
		testTime := time.Now()

		wg.Add(1)
		go func(task testers.Task, testTime time.Time) {
			defer wg.Done()
			pName := util.GetPNameFromTask(round, task, util.PNameRoleClient)

			// Template command and args for each task
			if err := cmdtemplate.Template(&task, templateVars); err != nil {
				errs <- fmt.Errorf("failed to template task command and / or args. %+v", err)
				return
			}

			pod = k.getPodSpec(pName, taskName, task)

			logger.WithFields(logrus.Fields{"pod": pName}).Debug("(re)creating client pod")
			if err := k8sutil.PodRecreate(k.k8sclient, pod, k.config.Timeouts.DeleteTimeout); err != nil {
				errs <- fmt.Errorf("failed to create pod %s/%s. %+v", k.config.Namespace, pName, err)
				return
			}

			logger.WithFields(logrus.Fields{"pod": pName}).Info("waiting for client pod to run or succeed")
			running, err := k8sutil.WaitForPodToRunOrSucceed(k.k8sclient, k.config.Namespace, pName, k.config.Timeouts.RunningTimeout)
			if err != nil {
				errs <- fmt.Errorf("failed to wait for pod %s/%s. %+v", k.config.Namespace, pName, err)
				return
			}
			if !running {
				errs <- fmt.Errorf("pod %s/%s not running after runTimeout", k.config.Namespace, pName)
				return
			}

			logger.WithFields(logrus.Fields{"pod": pName}).Debug("about to pushLogsToParser")
			if err := k.pushLogsToParser(parser, plannedTime, testTime, round, tester, mainTask.Host.Name, task.Host.Name, pName); err != nil {
				errs <- fmt.Errorf("failed to push pod %s/%s logs to parser. %+v", k.config.Namespace, pName, err)
				return
			}

			logger.WithFields(logrus.Fields{"pod": pName}).Info("deleting client pod")
			if err := k8sutil.PodDelete(k.k8sclient, pod, k.config.Timeouts.DeleteTimeout); err != nil {
				errs <- fmt.Errorf("failed to delete client pod %s/%s. %+v", k.config.Namespace, pName, err)
				return
			}
		}(task, testTime)

		if k.runOptions.Mode != config.RunModeParallel {
			wg.Wait()
		}
	}

	// Delete server pod
	logger.WithFields(logrus.Fields{"pod": serverPodName}).Info("deleting server pod")
	if err := k8sutil.PodDeleteByName(k.k8sclient, k.config.Namespace, serverPodName, k.config.Timeouts.DeleteTimeout); err != nil {
		logger.WithFields(logrus.Fields{"pod": serverPodName}).Errorf("failed to delete server pod. %+v", err)
	}

	// When RunOptions.Mode `parallel` then we wait after all test tasks have been run
	if k.runOptions.Mode == config.RunModeParallel {
		wg.Wait()
	}

	logger.Debug("done running tests in kubernetes for plan")

	close(errs)
	return nil
}

func (k Kubernetes) pushLogsToParser(parserInput chan<- parsers.Input, plannedTime time.Time, testTime time.Time, round int, tester string, serverHost string, clientHost string, podName string) error {
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
			PlannedTime:    plannedTime,
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
func (k Kubernetes) Cleanup(plan *testers.Plan) error {
	var wg sync.WaitGroup

	// Delete all Pods with label XYZ
	if err := k8sutil.PodDeleteByLabels(k.k8sclient, k.config.Namespace, map[string]string{
		k8sutil.TaskIDLabel: util.GetTaskName(plan),
	}); err != nil {
		k.logger.Errorf("error during pod delete by labels in cleanup. %+v", err)
		return err
	}
	wg.Wait()

	return nil
}
