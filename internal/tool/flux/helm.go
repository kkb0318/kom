package flux

import (
	"encoding/json"
	"fmt"
	"time"

	helmv1 "github.com/fluxcd/helm-controller/api/v2beta2"
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta2"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	komv1alpha1 "github.com/kkb0318/kom/api/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type FluxHelm struct {
	source sourcev1.HelmRepository
	helm   []helmv1.HelmRelease
}

type HelmValues struct {
	FullnameOverride string
}

func NewFluxHelm() (*FluxHelm, error) {
	repoName := "aws-controller-k8s"
	namespace := "ack-system"
	repoUrl := "oci://public.ecr.aws/aws-controllers-k8s"
	charts := []komv1alpha1.Chart{
    {
      Name: "s3-chart",
      Version: "*.*.*",

    },
  }
	helmrepo := sourcev1.HelmRepository{
		ObjectMeta: v1.ObjectMeta{
			Name:      repoName,
			Namespace: namespace,
		},
		TypeMeta: v1.TypeMeta{
			APIVersion: "source.toolkit.fluxcd.io/v1beta2",
			Kind:       "HelmRepository",
		},
		Spec: sourcev1.HelmRepositorySpec{
			Type:     "oci",
			Interval: v1.Duration{Duration: time.Minute},
			URL:      repoUrl,
		},
	}
	var hrs []helmv1.HelmRelease
	for _, chart := range charts {
		values := HelmValues{
			FullnameOverride: fmt.Sprintf("%s-controller", chart),
		}
		v, err := json.Marshal(values)
		if err != nil {
			return nil, err
		}
		hr := &helmv1.HelmRelease{
			ObjectMeta: v1.ObjectMeta{
				Name:      chart.Name,
				Namespace: namespace,
			},
			TypeMeta: v1.TypeMeta{
				APIVersion: "source.toolkit.fluxcd.io/v2beta2",
				Kind:       "HelmRelease",
			},
			Spec: helmv1.HelmReleaseSpec{
				Chart: helmv1.HelmChartTemplate{
					Spec: helmv1.HelmChartTemplateSpec{
						Chart:   chart.Name,
						Version: chart.Version,
						SourceRef: helmv1.CrossNamespaceObjectReference{
							Kind:      "HelmRepository",
							Name:      repoName,
							Namespace: namespace,
						},
					},
				},
				Values: &apiextensionsv1.JSON{
					Raw: v,
				},
			},
		}
		hrs = append(hrs, *hr)
	}
	f := &FluxHelm{
		source: helmrepo,
		helm:   hrs,
	}
	return f, nil
}
