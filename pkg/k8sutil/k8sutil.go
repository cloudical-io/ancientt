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
	"context"
	"fmt"
	"time"

	"github.com/cloudical-io/ancientt/testers"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

// PodRecreate delete Pod if it exists and create it again. If the Pod does not exist, create it.
func PodRecreate(k8sclient kubernetes.Interface, pod *corev1.Pod, delTimeout int) error {
	// Delete Pod if it exists
	if err := PodDelete(k8sclient, pod, delTimeout); err != nil {
		return err
	}

	// Create Pod again
	ctx := context.TODO()
	if _, err := k8sclient.CoreV1().Pods(pod.ObjectMeta.Namespace).Create(ctx, pod, metav1.CreateOptions{}); err != nil {
		if !errors.IsAlreadyExists(err) {
			return err
		}
	}

	return nil
}

// PodDelete delete Pod if it exists, wait for it till it has been for custom amount deleted
func PodDelete(k8sclient kubernetes.Interface, pod *corev1.Pod, timeout int) error {
	namespace := pod.ObjectMeta.Namespace
	podName := pod.ObjectMeta.Name

	// Delete Pod
	ctx := context.TODO()
	if err := k8sclient.CoreV1().Pods(namespace).Delete(ctx, podName, metav1.DeleteOptions{}); err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	for i := 0; i < timeout; i++ {
		// Check if Pod still exists
		ctx := context.TODO()
		if _, err := k8sclient.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{}); err != nil {
			if errors.IsNotFound(err) {
				return nil
			}
			return err
		}

		time.Sleep(1 * time.Second)
	}

	return fmt.Errorf("pod %s/%s not deleted after 30s", namespace, podName)
}

// PodDeleteByName delete Pod by namespace and name if it exists
func PodDeleteByName(k8sclient kubernetes.Interface, namespace string, podName string, timeout int) error {
	return PodDelete(k8sclient, &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      podName,
		},
	}, timeout)
}

// PodDeleteByLabels delete Pods by labels
func PodDeleteByLabels(k8sclient kubernetes.Interface, namespace string, selectorLabels map[string]string) error {
	set := labels.Set(selectorLabels)

	ctx := context.TODO()
	pods, err := k8sclient.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: set.AsSelector().String(),
	})
	if err != nil {
		if errors.IsNotFound(err) {
			return nil
		}
		return err
	}

	if pods.Items == nil || len(pods.Items) == 0 {
		return nil
	}

	for _, pod := range pods.Items {
		// Delete Pods by labels
		ctx := context.TODO()
		if err := k8sclient.CoreV1().Pods(namespace).Delete(ctx, pod.ObjectMeta.Name, metav1.DeleteOptions{}); err != nil {
			if errors.IsNotFound(err) {
				return nil
			}
			return err
		}
	}
	return nil
}

// WaitForPodToRun wait for a Pod to be in phase Running. In case of phase Running, return true and no error
func WaitForPodToRun(k8sclient kubernetes.Interface, namespace string, podName string, timeout int) (bool, error) {
	for i := 0; i < timeout; i++ {
		ctx := context.TODO()
		pod, err := k8sclient.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
		if err != nil {
			if !errors.IsAlreadyExists(err) {
				return false, err
			}
		}
		if pod.Status.Phase == corev1.PodRunning {
			return true, nil
		}

		time.Sleep(1 * time.Second)
	}

	return false, nil
}

// WaitForPodToSucceed wait for a Pod to be in phase Succeeded. In case of phase Succeeded, return true and no error
func WaitForPodToSucceed(k8sclient kubernetes.Interface, namespace string, podName string, timeout int) (bool, error) {
	for i := 0; i < timeout; i++ {
		ctx := context.TODO()
		pod, err := k8sclient.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		if pod.Status.Phase == corev1.PodSucceeded {
			return true, nil
		}

		time.Sleep(1 * time.Second)
	}

	return false, nil
}

// WaitForPodToRunOrSucceed wait for a Pod to be in phase Running or Succeeded. In case of one of the phases, return true and no error
func WaitForPodToRunOrSucceed(k8sclient kubernetes.Interface, namespace string, podName string, timeout int) (bool, error) {
	for i := 0; i < timeout; i++ {
		ctx := context.TODO()
		pod, err := k8sclient.CoreV1().Pods(namespace).Get(ctx, podName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}
		if pod.Status.Phase == corev1.PodRunning || pod.Status.Phase == corev1.PodSucceeded {
			return true, nil
		}

		time.Sleep(1 * time.Second)
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
