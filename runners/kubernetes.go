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
func (k Kubernetes) Execute(cmd, args []string) ([]byte, error) {
	// TODO Run tasks using Kubernetes Job objects

	return []byte{}, nil
}
