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

	"github.com/cloudical-io/acntt/pkg/config"
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
	// TOOD Get all Nodes, iterate over them and filter them for each Clients, Servers list
	// Use https://github.com/kubernetes/client-go

	return &testers.Hosts{}, nil
}

// Execute run the given commands and return the logs of it and / or error
func (k Kubernetes) Execute(plan *testers.Plan) (string, error) {
	// TODO Run tasks using Kubernetes Pods
	// TODO Get actual ip addresses of the Pod and use the cmdtemplate.Template() to template it in

	if err := k.prepareKubernetes(); err != nil {
		return "", err
	}

	// Iterate over given plan.Commands to then run each task
	for _, tasks := range plan.Commands {
		// Create the Pods for the server task and client tasks
		if err := k.createPodsForMainTask(tasks); err != nil {
			return "", err
		}
		// Exec into the Pods for the server task and client tasks, run the commands
		if err := k.execTasksInPods(tasks); err != nil {
			return "", err
		}
	}

	return "", nil
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

// createPodsForMainTask
func (k Kubernetes) createPodsForMainTask(tasks []testers.Task) error {
	pod, err := k.k8sclient.CoreV1().Pods(k.config.Namespace).Get("", metav1.GetOptions{})
	if err != nil {
		return err
	}
	_ = pod

	return nil
}

// execTasksInPods
func (k Kubernetes) execTasksInPods(tasks []testers.Task) error {
	// TODO use k8sclient exec to run the commands with retry

	return nil
}
