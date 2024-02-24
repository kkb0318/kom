package argo


import (
	"testing"

	argov1alpha1 "github.com/kkb0318/argo-cd-api/api/v1alpha1"
	komv1alpha1 "github.com/kkb0318/kom/api/v1alpha1"
	komtool "github.com/kkb0318/kom/internal/tool"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestArgoHelm_New(t *testing.T) {
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
						},
					},
				},
			},
			expected: []komtool.Resource{
				&ArgoHelm{
					source: []*corev1.Secret{
						{
							ObjectMeta: v1.ObjectMeta{
								Name:      "repo1",
								Namespace: "repo-ns1",
							},
							TypeMeta: v1.TypeMeta{
								APIVersion: "source.toolkit.Argocd.io/v1beta2",
								Kind:       "HelmRepository",
							},
						},
					},
					helm: []*argov1alpha1.Application{
						{
							ObjectMeta: v1.ObjectMeta{
								Name:      "chart1",
								Namespace: "repo-ns1",
							},
							TypeMeta: v1.TypeMeta{
								APIVersion: "helm.toolkit.Argocd.io/v2beta2",
								Kind:       "HelmRelease",
							},
						},
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
