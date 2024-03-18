package argo

import (
	"fmt"
	"testing"

	argov1alpha1 "github.com/kkb0318/argo-cd-api/api/v1alpha1"
	komv1alpha1 "github.com/kkb0318/kom/api/v1alpha1"
	komtool "github.com/kkb0318/kom/internal/tool"
	"github.com/kkb0318/kom/internal/tool/argo/testdata"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/yaml"
)

func TestArgoHelm_New(t *testing.T) {
	tests := []struct {
		name        string
		inputs      []komv1alpha1.Helm
		expected    []komtool.Resource
		expectedErr bool
	}{
		{
			name: "create manifest for installing helm with argo",
			inputs: []komv1alpha1.Helm{
				{
					Name:      "repo1",
					Namespace: "repo-ns1",
					Url:       "https://example.com",
					Charts: []komv1alpha1.Chart{
						{
							Name:    "chart1",
							Version: "1.0.0",
							Values: values(`
                key1: val1
                `),
						},
					},
				},
				{
					Name:      "repo2",
					Namespace: "repo-ns2",
					Url:       "https://example.com",
					Charts: []komv1alpha1.Chart{
						{
							Name:    "chart2",
							Version: "x.x.x",
						},
					},
				},
			},
			expected: []komtool.Resource{
				&ArgoHelm{
					source: []*corev1.Secret{
						testdata.NewMockSecretBuilder().Build(t, "helm_secret.yaml"),
					},
					helm: []*argov1alpha1.Application{
						testdata.NewMockApplicationBuilder().
							WithValues(values(`
                key1: val1
                `),
							).
							Build(t, "helm_application.yaml"),
					},
				},
				&ArgoHelm{
					source: []*corev1.Secret{
						testdata.NewMockSecretBuilder().
							WithName("repo2").
							WithNamespace("repo-ns2").
							WithChartName("chart2").
							Build(t, "helm_secret.yaml"),
					},
					helm: []*argov1alpha1.Application{
						testdata.NewMockApplicationBuilder().
							WithName("chart2").
							WithNamespace("repo-ns2").
							WithDestNamespace("repo-ns2").
							WithChartName("chart2").
							WithVersion("x.x.x").
							Build(t, "helm_application.yaml"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := NewArgoHelmList(tt.inputs)
			if tt.expectedErr {
				assert.Error(t, err, "")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, actual)
			}
		})
	}
}

func values(hrValues string) *apiextensionsv1.JSON {
	v, err := yaml.YAMLToJSON([]byte(hrValues))
	if err != nil {
		fmt.Println(err)
	}
	return &apiextensionsv1.JSON{Raw: v}
}
