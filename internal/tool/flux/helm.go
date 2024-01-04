package flux

import (
	"fmt"
	"time"

	helmv1 "github.com/fluxcd/helm-controller/api/v2beta2"
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta2"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

type FluxHelm struct {
	source sourcev1.HelmRepository
	helm   []helmv1.HelmRelease
}

func NewFluxHelm() *FluxHelm {
	pkgs := []string{"s3"}
	helmrepo := sourcev1.HelmRepository{
		ObjectMeta: v1.ObjectMeta{
			Name:      "aws-controller-k8s",
			Namespace: "flux-system",
		},
		TypeMeta: v1.TypeMeta{
			APIVersion: "source.toolkit.fluxcd.io/v1beta2",
			Kind:       "HelmRepository",
		},
		Spec: sourcev1.HelmRepositorySpec{
			Type:     "oci",
			Interval: v1.Duration{Duration: time.Minute},
			URL:      "",
		},
	}
	var hrs []helmv1.HelmRelease
	for _, pkg := range pkgs {
		hr := &helmv1.HelmRelease{
			ObjectMeta: v1.ObjectMeta{
				Name:      pkg,
				Namespace: "flux-system",
			},
			TypeMeta: v1.TypeMeta{
				APIVersion: "source.toolkit.fluxcd.io/v2beta2",
				Kind:       "HelmRelease",
			},
      Spec: helmv1.HelmReleaseSpec{
        Chart: helmv1.HelmChartTemplate{
          Spec: helmv1.HelmChartTemplateSpec{
            Chart: fmt.Sprintf("%s-chart", pkg),
            Version: "",
            SourceRef: helmv1.CrossNamespaceObjectReference{
              Kind: "HelmRepository",
              Name: "aws-controller-k8s",
            },
          },
        },
        Values: &apiextensionsv1.JSON{
        },
      },
		}
		hrs = append(hrs, *hr)
	}
	f := &FluxHelm{
		source: helmrepo,
		helm:   hrs,
	}
	return f
}
