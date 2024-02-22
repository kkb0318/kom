package flux

import (
	"strings"
	"time"

	helmv2beta2 "github.com/fluxcd/helm-controller/api/v2beta2"
	sourcev1beta2 "github.com/fluxcd/source-controller/api/v1beta2"
	komv1alpha1 "github.com/kkb0318/kom/api/v1alpha1"
	komtool "github.com/kkb0318/kom/internal/tool"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type FluxHelm struct {
	source *sourcev1beta2.HelmRepository
	helm   []*helmv2beta2.HelmRelease
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
	helmrepo := &sourcev1beta2.HelmRepository{
		ObjectMeta: v1.ObjectMeta{
			Name:      repoName,
			Namespace: namespace,
		},
		TypeMeta: v1.TypeMeta{
			APIVersion: sourcev1beta2.GroupVersion.String(),
			Kind:       sourcev1beta2.HelmRepositoryKind,
		},
		Spec: sourcev1beta2.HelmRepositorySpec{
			Type:     RepositoryType(repoUrl),
			Interval: v1.Duration{Duration: time.Minute},
			URL:      repoUrl,
		},
	}
	hrs := make([]*helmv2beta2.HelmRelease, len(charts))
	for i, chart := range charts {
		chartNs := namespace
		// values := HelmValues{
		// 	FullnameOverride: fmt.Sprintf("%s-controller", chart.Name),
		// }
		// v, err := json.Marshal(values)
		// if err != nil {
		// 	return nil, err
		// }
		hr := &helmv2beta2.HelmRelease{
			ObjectMeta: v1.ObjectMeta{
				Name:      chart.Name,
				Namespace: chartNs,
			},
			TypeMeta: v1.TypeMeta{
				APIVersion: helmv2beta2.GroupVersion.String(),
				Kind:       helmv2beta2.HelmReleaseKind,
			},
			Spec: helmv2beta2.HelmReleaseSpec{
				Chart: helmv2beta2.HelmChartTemplate{
					Spec: helmv2beta2.HelmChartTemplateSpec{
						Chart:   chart.Name,
						Version: chart.Version,
						SourceRef: helmv2beta2.CrossNamespaceObjectReference{
							Kind:      sourcev1beta2.HelmRepositoryKind,
							Name:      repoName,
							Namespace: namespace,
						},
					},
				},
				// Values: &apiextensionsv1.JSON{
				// 	Raw: v,
				// },
			},
		}
		hrs[i] = hr
	}
	f := &FluxHelm{
		source: helmrepo,
		helm:   hrs,
	}
	return f, nil
}
