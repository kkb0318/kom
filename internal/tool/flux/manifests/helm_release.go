package manifests

import (
	"errors"

	helmv2 "github.com/fluxcd/helm-controller/api/v2"
	sourcev1 "github.com/fluxcd/source-controller/api/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type HelmReleaseBuilder struct {
	ref     *helmv2.CrossNamespaceObjectReference
	chart   string
	version string
	values  *apiextensionsv1.JSON
}

func NewHelmReleaseBuilder() *HelmReleaseBuilder {
	return &HelmReleaseBuilder{
		version: "x",
	}
}

func (b *HelmReleaseBuilder) WithChart(value string) *HelmReleaseBuilder {
	b.chart = value
	return b
}

func (b *HelmReleaseBuilder) WithVersion(value string) *HelmReleaseBuilder {
	b.version = value
	return b
}

func (b *HelmReleaseBuilder) WithReference(name, namespace string) *HelmReleaseBuilder {
	ref := &helmv2.CrossNamespaceObjectReference{
		Kind:      sourcev1.HelmRepositoryKind,
		Name:      name,
		Namespace: namespace,
	}
	b.ref = ref
	return b
}

func (b *HelmReleaseBuilder) WithValues(values *apiextensionsv1.JSON) *HelmReleaseBuilder {
	b.values = values
	return b
}

func (b *HelmReleaseBuilder) Build(name, ns string) (*helmv2.HelmRelease, error) {
	if b.ref == nil {
		return nil, errors.New("the 'ref' field is nil and must be provided with a valid reference")
	}
	if b.chart == "" {
		return nil, errors.New("the 'chart' field is empty. Please specify a valid URL")
	}
	helmrelease := &helmv2.HelmRelease{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		TypeMeta: v1.TypeMeta{
			APIVersion: helmv2.GroupVersion.String(),
			Kind:       helmv2.HelmReleaseKind,
		},
		Spec: helmv2.HelmReleaseSpec{
			Values: b.values,
			Chart: &helmv2.HelmChartTemplate{
				Spec: helmv2.HelmChartTemplateSpec{
					Chart:     b.chart,
					Version:   b.version,
					SourceRef: *b.ref,
				},
			},
		},
	}
	return helmrelease, nil
}
