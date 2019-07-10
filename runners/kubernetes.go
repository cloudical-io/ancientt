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
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/cloudical-io/acntt/parsers"
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
	config    *config.RunnerKubernetes
	k8sclient *kubernetes.Clientset
}

// NewKubernetesRunner return a new Kubernetes Runner
func NewKubernetesRunner(cfg *config.Config) (Runner, error) {
	if cfg.Runner.Kubernetes == nil {
		return nil, fmt.Errorf("no kubernetes runner config")
	}
	// use the current context in kubeconfig
	k8sconfig, err := clientcmd.BuildConfigFromFlags("", cfg.Runner.Kubernetes.Kubeconfig)
	if err != nil {
		return nil, err
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(k8sconfig)
	if err != nil {
		return nil, err
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
	for _, node := range nodes.Items {
		hosts = append(hosts, &testers.Host{
			Labels: node.ObjectMeta.Labels,
			Name:   node.ObjectMeta.Name,
		})
	}

	return hosts, nil
}

// Prepare prepare Kubernetes for usage with acntt, e.g., create Namespace.
func (k Kubernetes) Prepare(plan *testers.Plan) error {
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
			if err := k.createPodsForTasks(round, task, taskName); err != nil {
				return err
			}

			// TODO Get log streams of Pods and "push" them to the parser
			if err := k.pushLogsToParser(round, task, taskName, parser); err != nil {
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
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

// createPodsForTasks create the Pods that are needed for the task(s)
func (k Kubernetes) createPodsForTasks(round int, mainTask testers.Task, taskName string) error {
	var wg sync.WaitGroup
	errs := make(chan error)

	// Create server Pod first
	pName := util.GetPNameFromTask(round, mainTask)

	pod := k.getPodSpec(pName, taskName, mainTask)

	if err := k8sutil.PodRecreate(k.k8sclient, pod); err != nil {
		return err
	}

	running, err := k8sutil.WaitForPodToRun(k.k8sclient, pod.ObjectMeta.Namespace, pName)
	if err != nil {
		return err
	}
	if !running {
		return fmt.Errorf("pod %s/%s not running after 10 tries", pod.ObjectMeta.Namespace, pName)
	}

	tasks := []testers.Task{}
	tasks = append(tasks, mainTask.SubTasks...)
	for _, task := range tasks {
		wg.Add(1)
		go func(task testers.Task) {
			defer wg.Done()
			pName := util.GetPNameFromTask(round, task)

			pod := k.getPodSpec(pName, taskName, task)

			if err := k8sutil.PodRecreate(k.k8sclient, pod); err != nil {
				errs <- err
				return
			}

			running, err := k8sutil.WaitForPodToRun(k.k8sclient, pod.ObjectMeta.Namespace, pName)
			if err != nil {
				errs <- err
				return
			}
			if !running {
				errs <- fmt.Errorf("pod %s/%s not running after 10 tries", pod.ObjectMeta.Namespace, pName)
				return
			}
		}(task)
	}
	wg.Wait()

	select {
	case err := <-errs:
		return err
	default:
		return nil
	}
}

func (k Kubernetes) pushLogsToParser(round int, task testers.Task, taskName string, parser parsers.Parser) error {
	//k8sutil.WaitForPodToSucceed()
	return nil
}

// Cleanup remove all (left behind) Kubernetes resources created for the given Plan.
func (k Kubernetes) Cleanup(plan *testers.Plan) error {
	// TODO
	return nil
}
