package flux

import (
	"strings"

	helmv2 "github.com/fluxcd/helm-controller/api/v2"
	sourcev1 "github.com/fluxcd/source-controller/api/v1"
	komv1alpha1 "github.com/kkb0318/kom/api/v1alpha1"
	komtool "github.com/kkb0318/kom/internal/tool"
	"github.com/kkb0318/kom/internal/tool/flux/manifests"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type FluxHelm struct {
	source *sourcev1.HelmRepository
	helm   []*helmv2.HelmRelease
}

type HelmValues struct {
	FullnameOverride string
}

func (f *FluxHelm) Repositories() []client.Object {
	return []client.Object{f.source}
}

func (f *FluxHelm) Charts() []client.Object {
	objs := make([]client.Object, len(f.helm))
	for i, helm := range f.helm {
		objs[i] = helm
	}
	return objs
}

func NewFluxHelmList(objs []komv1alpha1.Helm) ([]komtool.Resource, error) {
	helmList := make([]komtool.Resource, len(objs))
	for i, obj := range objs {
		helm, err := NewFluxHelm(obj)
		if err != nil {
			return nil, err
		}
		helmList[i] = helm
	}
	return helmList, nil
}

func RepositoryType(url string) string {
	if strings.HasPrefix(url, "oci:") {
		return "oci"
	}
	return "default"
}

func NewFluxHelm(obj komv1alpha1.Helm) (*FluxHelm, error) {
	repoName := obj.Name
	var namespace string
	if obj.Namespace == "" {
		namespace = komv1alpha1.DefaultNamespace
	} else {
		namespace = obj.Namespace
	}
	repoUrl := obj.Url
	charts := obj.Charts
	helmrepo, err := manifests.NewHelmRepositoryBuilder().
		WithUrl(repoUrl).
		Build(repoName, namespace)
	if err != nil {
		return nil, err
	}
	hrs := make([]*helmv2.HelmRelease, len(charts))
	for i, chart := range charts {
		hr, err := manifests.NewHelmReleaseBuilder().
			WithReference(repoName, namespace).
			WithChart(chart.Name).
			WithVersion(chart.Version).
			WithValues(chart.Values).
			Build(chart.Name, namespace)
		if err != nil {
			return nil, err
		}
		hrs[i] = hr
	}
	f := &FluxHelm{
		source: helmrepo,
		helm:   hrs,
	}
	return f, nil
}
