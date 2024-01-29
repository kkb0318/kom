package flux

import (
	"strings"
	"time"

	helmv1 "github.com/fluxcd/helm-controller/api/v2beta2"
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta2"
	komv1alpha1 "github.com/kkb0318/kom/api/v1alpha1"
	komtool "github.com/kkb0318/kom/internal/tool"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type FluxHelm struct {
	source *sourcev1.HelmRepository
	helm   []*helmv1.HelmRelease
}

type HelmValues struct {
	FullnameOverride string
}

func (f *FluxHelm) Repository() client.Object {
	return f.source
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
	helmrepo := &sourcev1.HelmRepository{
		ObjectMeta: v1.ObjectMeta{
			Name:      repoName,
			Namespace: namespace,
		},
		TypeMeta: v1.TypeMeta{
			APIVersion: sourcev1.GroupVersion.String(),
			Kind:       sourcev1.HelmRepositoryKind,
		},
		Spec: sourcev1.HelmRepositorySpec{
			Type:     RepositoryType(repoUrl),
			Interval: v1.Duration{Duration: time.Minute},
			URL:      repoUrl,
		},
	}
	hrs := make([]*helmv1.HelmRelease, len(charts))
	for i, chart := range charts {
		var chartNs string
		if chart.Namespace == "" {
			chartNs = namespace
		} else {
			chartNs = chart.Namespace
		}
		// values := HelmValues{
		// 	FullnameOverride: fmt.Sprintf("%s-controller", chart.Name),
		// }
		// v, err := json.Marshal(values)
		// if err != nil {
		// 	return nil, err
		// }
		hr := &helmv1.HelmRelease{
			ObjectMeta: v1.ObjectMeta{
				Name:      chart.Name,
				Namespace: chartNs,
			},
			TypeMeta: v1.TypeMeta{
				APIVersion: helmv1.GroupVersion.String(),
				Kind:       helmv1.HelmReleaseKind,
			},
			Spec: helmv1.HelmReleaseSpec{
				Chart: helmv1.HelmChartTemplate{
					Spec: helmv1.HelmChartTemplateSpec{
						Chart:   chart.Name,
						Version: chart.Version,
						SourceRef: helmv1.CrossNamespaceObjectReference{
							Kind:      sourcev1.HelmRepositoryKind,
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
