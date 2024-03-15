package argo

import (
	"path/filepath"
	"runtime"
	"testing"

	argov1alpha1 "github.com/kkb0318/argo-cd-api/api/v1alpha1"
	komv1alpha1 "github.com/kkb0318/kom/api/v1alpha1"
	komtool "github.com/kkb0318/kom/internal/tool"
	"github.com/kkb0318/kom/internal/utils"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
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
					source: newMockSecretBuilder().Build(filepath.Join(currentDir(t), "testdata", "git_secret.yaml")),
					app: newMockApplicationBuilder().Build(filepath.Join(currentDir(t), "testdata", "git_application.yaml")),
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

type mockSecretBuilder struct{}

func newMockSecretBuilder() *mockSecretBuilder {
	return &mockSecretBuilder{}
}

func (f *mockSecretBuilder) Build(baseFilePath string) *corev1.Secret {
	Secret := &corev1.Secret{}
	utils.LoadYaml(Secret, baseFilePath)
	return Secret
}

type mockApplicationBuilder struct{}

func newMockApplicationBuilder() *mockApplicationBuilder {
	return &mockApplicationBuilder{}
}

func (f *mockApplicationBuilder) Build(baseFilePath string) *argov1alpha1.Application {
	application := &argov1alpha1.Application{}
	utils.LoadYaml(application, baseFilePath)
	return application
}

func currentDir(t *testing.T) string {
	t.Helper()
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller() failed to get current file path")
	}
	return filepath.Dir(currentFile)
}
