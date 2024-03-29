package flux

import (
	komv1alpha1 "github.com/kkb0318/kom/api/v1alpha1"
	komtool "github.com/kkb0318/kom/internal/tool"
)

type Flux struct {
	resource komv1alpha1.Resource
}

func NewFlux(obj komv1alpha1.OperatorManager) *Flux {
	return &Flux{obj.Spec.Resource}
}

func (f *Flux) Helm() ([]komtool.Resource, error) {
	helmResources, err := NewFluxHelmList(f.resource.Helm)
	if err != nil {
		return nil, err
	}
	return helmResources, nil
}

func (f *Flux) Git() ([]komtool.Resource, error) {
	gitResources, err := NewFluxGitList(f.resource.Git)
	if err != nil {
		return nil, err
	}
	return gitResources, nil
}
