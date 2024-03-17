# API Reference

## Packages
- [kom.kkb.jp/v1alpha1](#komkkbjpv1alpha1)


## kom.kkb.jp/v1alpha1

Package v1alpha1 contains API Schema definitions for the  v1alpha1 API group

### Resource Types
- [OperatorManager](#operatormanager)



#### Chart



Chart defines the details of a Helm chart to be managed. Depending on the GitOps tool (Flux or ArgoCD), it corresponds to a HelmRelease CR or an Application CR, respectively.

_Appears in:_
- [Helm](#helm)

| Field | Description |
| --- | --- |
| `name` _string_ | Name is the name of the Helm chart. |
| `version` _string_ | Version specifies the version of the Helm chart to be deployed. |
| `values` _[JSON](#json)_ | Values specifies Helm values to be passed to helm template, defined as a map. |


#### Git



Git defines the configuration for accessing a Git repository.

_Appears in:_
- [Resource](#resource)

| Field | Description |
| --- | --- |
| `name` _string_ | Name is a user-defined identifier for the resource. |
| `namespace` _string_ | Namespace is the Kubernetes namespace where the Helm repository resource is located. The default value depends on the GitOps tool used: "kom-system" for Flux and "argocd" for ArgoCD. |
| `url` _string_ | Url is the URL of the Git repository. |
| `path` _string_ | Path specifies the directory path within the Git repository that contains the desired resources. This allows for selective management of resources located in specific parts of the repository. |
| `reference` _[GitReference](#gitreference)_ | Reference contains the reference information (such as branch, tag, or semver) for the Git repository. This allows for targeting specific versions or configurations of the resources within the repository. |


#### GitReference



GitReference specifies the versioning information for tracking changes in the Git repository.

_Appears in:_
- [Git](#git)

| Field | Description |
| --- | --- |
| `type` _[GitReferenceType](#gitreferencetype)_ | Type indicates the method of versioning used in the repository, applicable only for Flux. Valid options are "branch", "semver", or "tag", allowing for different strategies of version management. |
| `value` _string_ | Value specifies the exact reference to track, such as the name of a branch, a semantic versioning pattern, or a tag. This allows for precise control over which version of the resources is deployed. |


#### GitReferenceType

_Underlying type:_ _string_

GitReferenceType is applicable only for Flux. Valid options are "branch", "semver", or "tag"

_Appears in:_
- [GitReference](#gitreference)



#### Helm



Helm defines the configuration for accessing a Helm repository. Depending on the GitOps tool in use (Flux or ArgoCD), it corresponds to a HelmRepository CR or a Secret, respectively.

_Appears in:_
- [Resource](#resource)

| Field | Description |
| --- | --- |
| `name` _string_ | Name is a user-defined identifier for the resource. |
| `namespace` _string_ | Namespace is the Kubernetes namespace where the Helm repository resource is located. The default value depends on the GitOps tool used: "kom-system" for Flux and "argocd" for ArgoCD. |
| `url` _string_ | Url is the URL of the Helm repository. |
| `charts` _[Chart](#chart) array_ | Charts specifies the Helm charts within the repository to be managed. |


#### OperatorManager



OperatorManager is the Schema for the operatormanagers API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `kom.kkb.jp/v1alpha1`
| `kind` _string_ | `OperatorManager`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.29/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[OperatorManagerSpec](#operatormanagerspec)_ |  |


#### OperatorManagerSpec



OperatorManagerSpec defines the desired state and configuration of an OperatorManager.

_Appears in:_
- [OperatorManager](#operatormanager)

| Field | Description |
| --- | --- |
| `cleanup` _boolean_ | Cleanup, when enabled, allows the OperatorManager to perform garbage collection of resources that are no longer needed or managed. |
| `tool` _[ToolType](#tooltype)_ | Tool specifies the GitOps tool to be used. Users must set this field to either "flux" or "argo". This field is required and determines various default behaviors and configurations. |
| `resource` _[Resource](#resource)_ | Resource specifies the source repository (Helm or Git) for the operators to be managed. |




#### Resource



Resource represents the source repositories for operators, supporting both Helm and Git repositories.

_Appears in:_
- [OperatorManagerSpec](#operatormanagerspec)

| Field | Description |
| --- | --- |
| `helm` _[Helm](#helm) array_ | Helm specifies one or more Helm repositories containing the operators. This field is optional and only needed if operators are to be sourced from Helm repositories. |
| `git` _[Git](#git) array_ | Git specifies one or more Git repositories containing the operators. This field is optional and only needed if operators are to be sourced from Git repositories. |


#### ToolType

_Underlying type:_ _string_

ToolType defines the GitOps tool used for managing resources Valid options are "flux", or "argo".

_Appears in:_
- [OperatorManagerSpec](#operatormanagerspec)



