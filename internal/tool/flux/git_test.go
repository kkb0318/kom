package flux

import (
	"testing"

	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1"
	sourcev1 "github.com/fluxcd/source-controller/api/v1"
	komv1alpha1 "github.com/kkb0318/kom/api/v1alpha1"
	komtool "github.com/kkb0318/kom/internal/tool"
	"github.com/kkb0318/kom/internal/tool/flux/testdata"
	"github.com/stretchr/testify/assert"
)

func TestFluxGit_New(t *testing.T) {
	tests := []struct {
		name        string
		inputs      []komv1alpha1.Git
		expected    []komtool.Resource
		expectedErr bool
	}{
		{
			name: "continue if not in previous",
			inputs: []komv1alpha1.Git{
				{
					Name:      "repo1",
					Namespace: "repo-ns1",
					Url:       "https://example1.com",
					Path:      "./path1",
					Reference: komv1alpha1.GitReference{
						Type:  komv1alpha1.GitTag,
						Value: "1.0.0",
					},
				},
				{
					Name:      "repo2",
					Namespace: "repo-ns2",
					Url:       "https://example2.com",
					Path:      "./path2",
					Reference: komv1alpha1.GitReference{
						Type:  komv1alpha1.GitBranch,
						Value: "main",
					},
				},
				{
					Name:      "repo3",
					Namespace: "repo-ns3",
					Url:       "https://example3.com",
					Path:      "./path3",
					Reference: komv1alpha1.GitReference{
						Type:  komv1alpha1.GitSemver,
						Value: "x.x.x",
					},
				},
			},
			expected: []komtool.Resource{
				&FluxGit{
					source: testdata.NewMockGitRepositoryBuilder().Build(t, "git_repository.yaml"),
					ks:     testdata.NewMockKustomizationBuilder().Build(t, "kustomization.yaml"),
				},
				&FluxGit{
					source: testdata.NewMockGitRepositoryBuilder().
						WithName("repo2").
						WithNamespace("repo-ns2").
						WithUrl("https://example2.com").
						WithRef(&sourcev1.GitRepositoryRef{Branch: "main"}).
						Build(t, "git_repository.yaml"),
					ks: testdata.NewMockKustomizationBuilder().
						WithName("repo2").
						WithNamespace("repo-ns2").
						WithRef(&kustomizev1.CrossNamespaceSourceReference{Kind: "GitRepository", Name: "repo2", Namespace: "repo-ns2"}).
						WithPath("./path2").
						Build(t, "kustomization.yaml"),
				},
				&FluxGit{
					source: testdata.NewMockGitRepositoryBuilder().
						WithName("repo3").
						WithNamespace("repo-ns3").
						WithUrl("https://example3.com").
						WithRef(&sourcev1.GitRepositoryRef{SemVer: "x.x.x"}).
						Build(t, "git_repository.yaml"),
					ks: testdata.NewMockKustomizationBuilder().
						WithName("repo3").
						WithNamespace("repo-ns3").
						WithRef(&kustomizev1.CrossNamespaceSourceReference{Kind: "GitRepository", Name: "repo3", Namespace: "repo-ns3"}).
						WithPath("./path3").
						Build(t, "kustomization.yaml"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := NewFluxGitList(tt.inputs)
			if tt.expectedErr {
				assert.Error(t, err, "")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, actual)
			}
		})
	}
}
