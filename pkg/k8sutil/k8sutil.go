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

package k8sutil

import (
	"time"

	"github.com/cloudical-io/acntt/testers"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// PodRecreate delete Pod if it exists and create it again. If the Pod does not exist, create it.
func PodRecreate(k8sclient *kubernetes.Clientset, pod *corev1.Pod) error {
	if err := k8sclient.CoreV1().Pods(pod.ObjectMeta.Namespace).Delete(pod.ObjectMeta.Name, &metav1.DeleteOptions{}); err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	}

	if _, err := k8sclient.CoreV1().Pods(pod.ObjectMeta.Namespace).Create(pod); err != nil {
		if errors.IsAlreadyExists(err) {
			return err
		}
	}
	return nil
}

// WaitForPodToRun wait for a Pod to be in phase Running. In case of phase Running, return true and no error
func WaitForPodToRun(k8sclient *kubernetes.Clientset, namespace string, podName string) (bool, error) {
	// 10 tries with 3 second sleep so 30 seconds in total
	for i := 0; i < 10; i++ {
		pod, err := k8sclient.CoreV1().Pods(namespace).Get(podName, metav1.GetOptions{})
		if err != nil {
			if !errors.IsAlreadyExists(err) {
				return false, err
			}
		}
		if pod.Status.Phase == corev1.PodRunning {
			return true, nil
		}
		time.Sleep(3 * time.Second)
	}
	return false, nil
}

// WaitForPodToSucceed wait for a Pod to be in phase Succeeded. In case of phase Succeeded, return true and no error
func WaitForPodToSucceed(k8sclient *kubernetes.Clientset, namespace string, podName string) (bool, error) {
	// 15 tries with 3 second sleep so 45 seconds in total
	for i := 0; i < 15; i++ {
		pod, err := k8sclient.CoreV1().Pods(namespace).Get(podName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		if pod.Status.Phase == corev1.PodSucceeded {
			return true, nil
		}
		time.Sleep(3 * time.Second)
	}
	return false, nil
}

// PortsListToPorts PortList testers.Port to Kubernetes []corev1.ContainerPort conversion (for TCP and UDP)
func PortsListToPorts(list testers.Ports) []corev1.ContainerPort {
	ports := []corev1.ContainerPort{}
	for _, p := range list.TCP {
		ports = append(ports, corev1.ContainerPort{
			ContainerPort: p,
			Protocol:      corev1.ProtocolTCP,
		})
	}
	for _, p := range list.UDP {
		ports = append(ports, corev1.ContainerPort{
			ContainerPort: p,
			Protocol:      corev1.ProtocolUDP,
		})
	}
	return ports
}
