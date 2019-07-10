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

// GetLabels return a default set of labels for "any" object acntt is going to create.
func GetLabels() map[string]string {
	name := "acntt"
	return map[string]string{
		"app.kubernetes.io/part-of":    name,
		"app.kubernetes.io/managed-by": name,
		"app.kubernetes.io/version":    "0.0.1",
	}
}

// GetPodLabels default labels combined with additional labels for Pods.
func GetPodLabels(podName string, taskName string) map[string]string {
	labels := GetLabels()
	labels["app.kubernetes.io/instance"] = podName
	labels["acntt/task-id"] = taskName
	return labels
}
