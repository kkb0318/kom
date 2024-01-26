/*
Copyright 2023.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// OperatorManagerKind is the string representation of a OperatorManager.OperatorManagerKind
	OperatorManagerKind = "OperatorManager"
	DefaultNamespace    = "kom-system"
)

// OperatorManagerSpec defines the desired state of OperatorManager
type OperatorManagerSpec struct {
	// Prune enables garbage collection.
	// +required
	Prune bool `json:"prune"`

	Tool ToolType `json:"tool,omitempty"`
	// +required
	Resource Resource `json:"resource"`
}

type Resource struct {
	Helm []Helm `json:"helm,omitempty"`
}

type Helm struct {
	Name      string  `json:"name,omitempty"`
	Namespace string  `json:"namespace,omitempty"`
	Url       string  `json:"url,omitempty"`
	Charts    []Chart `json:"charts,omitempty"`
}
type Chart struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
	Version   string `json:"version,omitempty"`
}
type ToolType string

const (
	FluxCDTool ToolType = "flux"
	ArgoCDTool ToolType = "argo"
)

// OperatorManagerStatus defines the observed state of OperatorManager
type OperatorManagerStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// OperatorManager is the Schema for the operatormanagers API
type OperatorManager struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OperatorManagerSpec   `json:"spec,omitempty"`
	Status OperatorManagerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// OperatorManagerList contains a list of OperatorManager
type OperatorManagerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OperatorManager `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OperatorManager{}, &OperatorManagerList{})
}
