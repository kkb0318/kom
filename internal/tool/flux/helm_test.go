package flux

import (
	"fmt"
	"testing"

	helmv2beta2 "github.com/fluxcd/helm-controller/api/v2beta2"
	komv1alpha1 "github.com/kkb0318/kom/api/v1alpha1"
	komtool "github.com/kkb0318/kom/internal/tool"
	"github.com/kkb0318/kom/internal/tool/flux/testdata"
	"github.com/stretchr/testify/assert"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"sigs.k8s.io/yaml"
)

func values(hrValues string) *apiextensionsv1.JSON {
	v, err := yaml.YAMLToJSON([]byte(hrValues))
	if err != nil {
		fmt.Println(err)
	}
	return &apiextensionsv1.JSON{Raw: v}
}

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
							Name:    "chart1",
							Version: "x.x.x",
							Values: values(`
                key1: val1
                `),
						},
					},
				},
			},
			expected: []komtool.Resource{
				&FluxHelm{
					source: testdata.NewMockHelmRepositoryBuilder().
						Build(t, "helm_repository.yaml"),
					helm: []*helmv2beta2.HelmRelease{
						testdata.NewMockHelmReleaseBuilder().
							WithValues(values(`
                key1: val1
                `)).
							Build(t, "helm_release.yaml"),
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
