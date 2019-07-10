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
	"github.com/cloudical-io/acntt/pkg/k8sutil"
	"github.com/cloudical-io/acntt/testers"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k Kubernetes) getPodSpec(pName string, taskName string, task testers.Task) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Labels: k8sutil.GetPodLabels(pName, taskName),
			Name:   pName,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "acntt",
					Image:   k.config.Image,
					Command: []string{task.Command},
					Args:    task.Args,
					Ports:   k8sutil.PortsListToPorts(task.Ports),
				},
			},
			HostNetwork:   k.config.HostNetwork,
			RestartPolicy: corev1.RestartPolicyOnFailure,
		},
	}
}
