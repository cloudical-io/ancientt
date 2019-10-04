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

package k8sutil

import (
	"fmt"
	"os"
	"path/filepath"

	// Load all k8s client auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"github.com/mitchellh/go-homedir"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

// NewClient create a new Kubernetes clientset
func NewClient(inClusterConfig bool, kubeconfig string) (kubernetes.Interface, error) {
	var k8sconfig *rest.Config
	if inClusterConfig {
		var err error
		k8sconfig, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("kubeconfig in-cluster configuration error. %+v", err)
		}
	} else {
		var kubeconfig string
		// Try to fallback to the `KUBECONFIG` env var
		if kubeconfig == "" {
			kubeconfig = os.Getenv("KUBECONFIG")
		}
		// If the `KUBECONFIG` is empty, default to home dir default kube config path
		if kubeconfig == "" {
			home, err := homedir.Dir()
			if err != nil {
				return nil, fmt.Errorf("kubeconfig unable to get home dir. %+v", err)
			}
			kubeconfig = filepath.Join(home, ".kube", "config")
		}
		var err error
		// This will simply use the current context in the kubeconfig
		k8sconfig, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			return nil, fmt.Errorf("kubeconfig out-of-cluster configuration (%s) error. %+v", kubeconfig, err)
		}
	}

	// Create the clientset
	clientset, err := kubernetes.NewForConfig(k8sconfig)
	if err != nil {
		return nil, fmt.Errorf("kubernetes new client error. %+v", err)
	}
	return clientset, nil
}
