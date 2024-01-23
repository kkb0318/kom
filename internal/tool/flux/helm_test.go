package flux

import (
	"testing"
	"time"

	helmv1 "github.com/fluxcd/helm-controller/api/v2beta2"
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta2"
	komv1alpha1 "github.com/kkb0318/kom/api/v1alpha1"
	komtool "github.com/kkb0318/kom/internal/tool"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestFluxHelm_New(t *testing.T) {
	tests := []struct {
		name        string
		inputs      []komv1alpha1.Helm
		expected    []komtool.Resource
		expectedErr bool
	}{
		{
			name: "continue if not in previous",
			inputs: []komv1alpha1.Helm{
				{
					Name:      "repo1",
					Namespace: "repo-ns1",
					Url:       "https://example.com",
					Charts: []komv1alpha1.Chart{
						{
							Name:      "chart1",
							Namespace: "chart-ns1",
							Version:   "x.x.x",
						},
					},
				},
			},
			expected: []komtool.Resource{
				&FluxHelm{
					source: &sourcev1.HelmRepository{
						ObjectMeta: v1.ObjectMeta{
							Name:      "repo1",
							Namespace: "repo-ns1",
						},
						TypeMeta: v1.TypeMeta{
							APIVersion: "source.toolkit.fluxcd.io/v1beta2",
							Kind:       "HelmRepository",
						},
						Spec: sourcev1.HelmRepositorySpec{
							Type:     "default",
							Interval: v1.Duration{Duration: time.Minute},
							URL:      "https://example.com",
						},
					},
					helm: []*helmv1.HelmRelease{
						{
							ObjectMeta: v1.ObjectMeta{
								Name:      "chart1",
								Namespace: "chart-ns1",
							},
							TypeMeta: v1.TypeMeta{
								APIVersion: "helm.toolkit.fluxcd.io/v2beta2",
								Kind:       "HelmRelease",
							},
							Spec: helmv1.HelmReleaseSpec{
								Chart: helmv1.HelmChartTemplate{
									Spec: helmv1.HelmChartTemplateSpec{
										Chart:   "chart1",
										Version: "x.x.x",
										SourceRef: helmv1.CrossNamespaceObjectReference{
											Kind:      "HelmRepository",
											Name:      "repo1",
											Namespace: "repo-ns1",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := NewFluxHelmList(tt.inputs)
			if tt.expectedErr {
				assert.Error(t, err, "")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, actual)
			}
		})
	}
}
