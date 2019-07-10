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

package runners

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"sync"
	"time"

	"github.com/cloudical-io/acntt/parsers"
	"github.com/cloudical-io/acntt/pkg/cmdtemplate"
	"github.com/cloudical-io/acntt/pkg/config"
	"github.com/cloudical-io/acntt/pkg/k8sutil"
	"github.com/cloudical-io/acntt/pkg/util"
	"github.com/cloudical-io/acntt/testers"
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

	return Kubernetes{
		config:    cfg.Runner.Kubernetes,
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

	// Create and seed randomness source for the `random` selection of hosts
	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)
	r.Seed(time.Now().UnixNano())

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
		if len(servers.Hosts) > 0 {
			for _, host := range servers.Hosts {
				hosts.Servers[host] = &testers.Host{
					Name: host,
				}
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
		if len(clients.Hosts) > 0 {
			for _, host := range clients.Hosts {
				hosts.Clients[host] = &testers.Host{
					Name: host,
				}
			}
		}
	}

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
func (k Kubernetes) Execute(plan *testers.Plan, parser parsers.Parser) error {
	// TODO Get actual ip addresses of the Pod and use the cmdtemplate.Template() to template it in

	// Iterate over given plan.Commands to then run each task
	for round, tasks := range plan.Commands {
		for _, task := range tasks {
			taskName := fmt.Sprintf("acntt-%d", time.Now().Unix())

			// Create the Pods for the server task and client tasks
			if err := k.createPodsForTasks(round, task, taskName, parser); err != nil {
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
			if _, err := k.k8sclient.CoreV1().Namespaces().Create(ns); err != nil {
				return fmt.Errorf("failed to create namespace %s. %+v", k.config.Namespace, err)
			}
		} else {
			return fmt.Errorf("error while getting namespace %s. %+v", k.config.Namespace, err)
		}
	}
	return nil
}

// createPodsForTasks create the Pods that are needed for the task(s)
func (k Kubernetes) createPodsForTasks(round int, mainTask testers.Task, taskName string, parser parsers.Parser) error {
	var wg sync.WaitGroup
	errs := make(chan error)

	// Create server Pod first
	pName := util.GetPNameFromTask(round, mainTask)

	// Create initial cmdtemplate.Variables
	// TODO the port does not need to be mapped in Kubernetes case, but for other runners (e.g., Ansible)
	templateVars := cmdtemplate.Variables{
		ServerPort: 5601,
	}

	if err := cmdtemplate.Template(&mainTask, templateVars); err != nil {
		return fmt.Errorf("failed to template main task command and / or args. %+v", err)
	}

	pod := k.getPodSpec(pName, taskName, mainTask)

	if err := k8sutil.PodRecreate(k.k8sclient, pod); err != nil {
		return fmt.Errorf("failed to create pod %s/%s. %+v", k.config.Namespace, pName, err)
	}

	running, err := k8sutil.WaitForPodToRun(k.k8sclient, k.config.Namespace, pName)
	if err != nil {
		return fmt.Errorf("failed to wait for pod %s/%s. %+v", k.config.Namespace, pName, err)
	}
	if !running {
		return fmt.Errorf("pod %s/%s not running after 10 tries", k.config.Namespace, pName)
	}

	// Get server Pod to have the server IP for each client task
	pod, err = k.k8sclient.CoreV1().Pods(k.config.Namespace).Get(pName, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("failed to get pod %s/%s. %+v", k.config.Namespace, pName, err)
	}
	templateVars.ServerAddress = pod.Status.PodIP

	tasks := []testers.Task{}
	tasks = append(tasks, mainTask.SubTasks...)
	for _, task := range tasks {
		wg.Add(1)
		go func(task testers.Task) {
			defer wg.Done()
			pName := util.GetPNameFromTask(round, task)

			// Template command and args for each task
			if err := cmdtemplate.Template(&task, templateVars); err != nil {
				errs <- fmt.Errorf("failed to template task command and / or args. %+v", err)
				return
			}

			pod = k.getPodSpec(pName, taskName, task)

			if err := k8sutil.PodRecreate(k.k8sclient, pod); err != nil {
				errs <- fmt.Errorf("failed to create pod %s/%s. %+v", k.config.Namespace, pName, err)
				return
			}

			running, err := k8sutil.WaitForPodToRun(k.k8sclient, k.config.Namespace, pName)
			if err != nil {
				errs <- fmt.Errorf("failed to wait for pod %s/%s. %+v", k.config.Namespace, pName, err)
				return
			}
			if !running {
				errs <- fmt.Errorf("pod %s/%s not running after 10 tries", k.config.Namespace, pName)
				return
			}

			// In case of running sequentiell, just send the logs to the parser ASAP
			// TODO Find a better way to do this. Feels a bit broken to do it here instead of, e.g., using a goroutined channel ;-)
			if k.runOptions.Mode != config.RunModeParallel {
				if err := k.pushLogsToParser(pName, parser); err != nil {
					errs <- fmt.Errorf("failed to push pod %s/%s logs to parser. %+v", k.config.Namespace, pName, err)
					return
				}
			}
		}(task)

		if k.runOptions.Mode != config.RunModeParallel {
			wg.Wait()
		}
	}
	// When RunOptions.Mode `parallel` then we wait after all test tasks have been run
	if k.runOptions.Mode == config.RunModeParallel {
		wg.Wait()
	}

	select {
	case err := <-errs:
		return err
	default:
		return nil
	}
}

func (k Kubernetes) pushLogsToParser(podName string, parser parsers.Parser) error {
	succeeded, err := k8sutil.WaitForPodToSucceed(k.k8sclient, k.config.Namespace, podName)
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
		// Close it afterwards
		defer podLogs.Close()

		buf := new(bytes.Buffer)
		_, err = io.Copy(buf, podLogs)
		if err != nil {
			return fmt.Errorf("error in copy information from podLogs to buffer")
		}

		// Send the logs to the parser.Parser() func
		// TODO Find a better way to do this. Feels a bit broken to do it here instead of, e.g., using a goroutined channel ;-)
		parsed, err := parser.Parse(buf)
		_ = parsed
		if err != nil {
			return err
		}
		return nil
	}

	return fmt.Errorf("pod %s/%s has not succeeded", k.config.Namespace, podName)
}

// Cleanup remove all (left behind) Kubernetes resources created for the given Plan.
func (k Kubernetes) Cleanup(plan *testers.Plan) error {
	// TODO
	return nil
}
