package argo

import (
	"testing"

	komv1alpha1 "github.com/kkb0318/kom/api/v1alpha1"
	komtool "github.com/kkb0318/kom/internal/tool"
	"github.com/kkb0318/kom/internal/tool/argo/testdata"
	"github.com/stretchr/testify/assert"
)

func TestArgoGit_New(t *testing.T) {
	tests := []struct {
		name        string
		inputs      []komv1alpha1.Git
		expected    []komtool.Resource
		expectedErr bool
	}{
		{
			name: "create manifest for installing helm with argo",
			inputs: []komv1alpha1.Git{
				{
					Name:      "repo1",
					Namespace: "repo-ns1",
					Url:       "https://example.com",
					Path:      "./path1",
					Reference: komv1alpha1.GitReference{
						Value: "1.0.0",
					},
				},
			},
			expected: []komtool.Resource{
				&ArgoGit{
					source: testdata.NewMockSecretBuilder().Build(t, "git_secret.yaml"),
					app:    testdata.NewMockApplicationBuilder().Build(t, "git_application.yaml"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := NewArgoGitList(tt.inputs)
			if tt.expectedErr {
				assert.Error(t, err, "")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, actual)
			}
		})
	}
}

