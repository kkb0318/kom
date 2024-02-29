package argo

import (
	komv1alpha1 "github.com/kkb0318/kom/api/v1alpha1"
	komtool "github.com/kkb0318/kom/internal/tool"
)

type Argo struct {
	resource komv1alpha1.Resource
}

func NewArgo(obj komv1alpha1.OperatorManager) *Argo {
	return &Argo{obj.Spec.Resource}
}

func (f *Argo) Helm() ([]komtool.Resource, error) {
	helmResources, err := NewArgoHelmList(f.resource.Helm)
	if err != nil {
		return nil, err
	}
	return helmResources, nil
}

func (f *Argo) Git() ([]komtool.Resource, error) {
	gitResources, err := NewArgoGitList(f.resource.Git)
	if err != nil {
		return nil, err
	}
	return gitResources, nil
}
