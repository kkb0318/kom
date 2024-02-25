package argo

import (
	argov1alpha1 "github.com/kkb0318/argo-cd-api/api/v1alpha1"
	komv1alpha1 "github.com/kkb0318/kom/api/v1alpha1"
	"github.com/kkb0318/kom/internal/tool/argo/manifests"
	komtool "github.com/kkb0318/kom/internal/tool"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ArgoHelm struct {
	source []*corev1.Secret
	helm   []*argov1alpha1.Application
}

func (f *ArgoHelm) Repositories() []client.Object {
	objs := make([]client.Object, len(f.source))
	for i, source := range f.source {
		objs[i] = source
	}
	return objs
}

func (f *ArgoHelm) Charts() []client.Object {
	objs := make([]client.Object, len(f.helm))
	for i, helm := range f.helm {
		objs[i] = helm
	}
	return objs
}

func NewArgoHelmList(objs []komv1alpha1.Helm) ([]komtool.Resource, error) {
	helmList := make([]komtool.Resource, len(objs))
	for i, obj := range objs {
		helm, err := NewArgoHelm(obj)
		if err != nil {
			return nil, err
		}
		helmList[i] = helm
	}
	return helmList, nil
}

func NewArgoHelm(obj komv1alpha1.Helm) (*ArgoHelm, error) {
	var namespace string
	if obj.Namespace == "" {
		namespace = komv1alpha1.ArgoCDDefaultNamespace
	} else {
		namespace = obj.Namespace
	}
	charts := obj.Charts
	secrets := make([]*corev1.Secret, len(charts))
	apps := make([]*argov1alpha1.Application, len(charts))
	for i, chart := range charts {
		secret := manifests.NewSecret(obj.Name, namespace, obj.Url, chart.Name)
    app := manifests.NewApplication(chart, namespace, obj.Url)
		apps[i] = app
		secrets[i] = secret
	}
	f := &ArgoHelm{
		source: secrets,
		helm:   apps,
	}
	return f, nil
}
