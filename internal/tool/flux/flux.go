package flux

import (
	komv1alpha1 "github.com/kkb0318/kom/api/v1alpha1"
	komtool "github.com/kkb0318/kom/internal/tool"
)

type Flux struct {
	obj komv1alpha1.OperatorManager
}

func NewFlux(obj komv1alpha1.OperatorManager) *Flux {
	return &Flux{obj}
}

func (f *Flux) Helm() ([]komtool.Resource, error) {
	return nil, nil
}

func (f *Flux) Oci() ([]komtool.Resource, error) {
	return nil, nil
}

func (f *Flux) Git() ([]komtool.Resource, error) {
	return nil, nil
}
