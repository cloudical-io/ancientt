/*
Copyright 2018 The Rook Authors. All rights reserved.
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
	corev1 "k8s.io/api/core/v1"
)

// INFO This code has been taken from the Rook project due to Golang dependency and modified to suite acntt's needs.
// Original file https://github.com/rook/rook/blob/master/pkg/operator/k8sutil/node.go

// NodeIsTolerable returns true if the node's taints are all tolerated by the given tolerations.
// There is the option to ignore well known taints defined in WellKnownTaints. See WellKnownTaints
// for more information.
func NodeIsTolerable(node corev1.Node, tolerations []corev1.Toleration) bool {
	for _, taint := range node.Spec.Taints {
		isTolerated := false
		for _, toleration := range tolerations {
			if toleration.ToleratesTaint(&taint) {
				isTolerated = true
				break
			}
		}
		if !isTolerated {
			return false
		}
	}
	return true
}
