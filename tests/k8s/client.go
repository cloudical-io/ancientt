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

package k8s

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

// NewClient return a new fake client
func NewClient(nodes int) (*fake.Clientset, error) {
	clientset := fake.NewSimpleClientset()
	for i := 0; i < nodes; i++ {
		n := &v1.Node{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("node-%d", i),
			},
			Status: v1.NodeStatus{
				Conditions: []v1.NodeCondition{
					v1.NodeCondition{
						Type:   v1.NodeReady,
						Status: v1.ConditionTrue,
					},
				},
				Addresses: []v1.NodeAddress{
					{
						Type:    v1.NodeExternalIP,
						Address: fmt.Sprintf("%d.%d.%d.%d", i, i, i, i),
					},
				},
			},
		}
		ctx := context.TODO()
		_, err := clientset.CoreV1().Nodes().Create(ctx, n, metav1.CreateOptions{})
		if err != nil {
			// Something is definitely wrong in the fake client
			return nil, err
		}
	}
	return clientset, nil
}
