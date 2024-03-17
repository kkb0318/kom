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
	"slices"

	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Constants for OperatorManager configurations.
const (
	// OperatorManagerKind represents the kind attribute of an OperatorManager resource.
	OperatorManagerKind = "OperatorManager"
	// DefaultNamespace is the default namespace where the OperatorManager operates.
	DefaultNamespace = "kom-system"
	// ArgoCDDefaultNamespace is the default namespace for ArgoCD resources.
	ArgoCDDefaultNamespace = "argocd"
)

// OperatorManagerSpec defines the desired state and configuration of an OperatorManager.
type OperatorManagerSpec struct {
	// Cleanup, when enabled, allows the OperatorManager to perform garbage collection
	// of resources that are no longer needed or managed.
	// +required
	Cleanup bool `json:"cleanup"`

	// Tool specifies the GitOps tool to be used. Users must set this field to either "flux" or "argo".
	// This field is required and determines various default behaviors and configurations.
	// +required
	Tool ToolType `json:"tool"`

	// Resource specifies the source repository (Helm or Git) for the operators to be managed.
	// +required
	Resource Resource `json:"resource"`
}

// Resource represents the source repositories for operators, supporting both Helm and Git repositories.
type Resource struct {
	// Helm specifies one or more Helm repositories containing the operators.
	// This field is optional and only needed if operators are to be sourced from Helm repositories.
	Helm []Helm `json:"helm,omitempty"`

	// Git specifies one or more Git repositories containing the operators.
	// This field is optional and only needed if operators are to be sourced from Git repositories.
	Git []Git `json:"git,omitempty"`
}

// Helm defines the configuration for accessing a Helm repository.
// Depending on the GitOps tool in use (Flux or ArgoCD), it corresponds to a HelmRepository CR or a Secret, respectively.
type Helm struct {
	// Name is a user-defined identifier for the resource.
	Name string `json:"name,omitempty"`

	// Namespace is the Kubernetes namespace where the Helm repository resource is located.
	// The default value depends on the GitOps tool used: "kom-system" for Flux and "argocd" for ArgoCD.
	Namespace string `json:"namespace,omitempty"`

	// Url is the URL of the Helm repository.
	Url string `json:"url,omitempty"`

	// Charts specifies the Helm charts within the repository to be managed.
	Charts []Chart `json:"charts,omitempty"`
}

// Chart defines the details of a Helm chart to be managed.
// Depending on the GitOps tool (Flux or ArgoCD), it corresponds to a HelmRelease CR or an Application CR, respectively.
type Chart struct {
	// Name is the name of the Helm chart.
	Name string `json:"name,omitempty"`

	// Version specifies the version of the Helm chart to be deployed.
	Version string `json:"version,omitempty"`

	// Values specifies Helm values to be passed to helm template, defined as a map.
	// +optional
	Values *apiextensionsv1.JSON `json:"values,omitempty"`
}

// Git defines the configuration for accessing a Git repository.
type Git struct {
	// Name is a user-defined identifier for the resource.
	Name string `json:"name,omitempty"`

	// Namespace is the Kubernetes namespace where the Helm repository resource is located.
	// The default value depends on the GitOps tool used: "kom-system" for Flux and "argocd" for ArgoCD.
	Namespace string `json:"namespace,omitempty"`

	// Url is the URL of the Git repository.
	Url string `json:"url,omitempty"`

	// Path specifies the directory path within the Git repository that contains the desired resources.
	// This allows for selective management of resources located in specific parts of the repository.
	Path string `json:"path,omitempty"`

	// Reference contains the reference information (such as branch, tag, or semver) for the Git repository.
	// This allows for targeting specific versions or configurations of the resources within the repository.
	Reference GitReference `json:"reference,omitempty"`
}

// GitReference specifies the versioning information for tracking changes in the Git repository.
type GitReference struct {
	// Type indicates the method of versioning used in the repository, applicable only for Flux.
	// Valid options are "branch", "semver", or "tag", allowing for different strategies of version management.
	Type GitReferenceType `json:"type,omitempty"`

	// Value specifies the exact reference to track, such as the name of a branch, a semantic versioning pattern, or a tag.
	// This allows for precise control over which version of the resources is deployed.
	Value string `json:"value,omitempty"`
}

// GitReferenceType is applicable only for Flux. Valid options are "branch", "semver", or "tag"
type GitReferenceType string

const (
	GitBranch GitReferenceType = "branch"
	GitSemver GitReferenceType = "semver"
	GitTag    GitReferenceType = "tag"
)

// ToolType defines the GitOps tool used for managing resources Valid options are "flux", or "argo".
type ToolType string

const (
	FluxCDTool ToolType = "flux"
	ArgoCDTool ToolType = "argo"
)

// OperatorManagerStatus defines the observed state of OperatorManager
type OperatorManagerStatus struct {
	// Inventory of applied resources
	AppliedResources AppliedResourceList `json:"appliedResources,omitempty"`
}
type AppliedResourceList []AppliedResource

// Unique identifier for the resource,  "namespace-name-kind-group-apiversion"
type AppliedResource struct {
	// Kind of the Kubernetes resource, e.g., Deployment, Service, etc.
	Kind string `json:"kind"`
	// APIVersion of the resource, e.g., "apps/v1"
	APIVersion string `json:"apiVersion"`
	// Name of the resource
	Name string `json:"name"`
	// Namespace of the resource, if applicable
	Namespace string `json:"namespace,omitempty"`
}

func (a AppliedResource) Equal(b AppliedResource) bool {
	return a.Name == b.Name &&
		a.Namespace == b.Namespace &&
		a.Kind == b.Kind &&
		a.APIVersion == b.APIVersion
}

// Diff returns the resourceList that exist in listA, but not in listB (A - B).
func (listA AppliedResourceList) Diff(listB AppliedResourceList) AppliedResourceList {
	diff := append(AppliedResourceList{}, listA...)
	diff = slices.DeleteFunc(diff, func(a AppliedResource) bool {
		return slices.ContainsFunc(listB, func(b AppliedResource) bool {
			return b.Equal(a)
		})
	})
	return diff
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
