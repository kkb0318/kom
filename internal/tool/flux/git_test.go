package flux

import (
	"testing"
	"time"

	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1"
	sourcev1 "github.com/fluxcd/source-controller/api/v1"
	komv1alpha1 "github.com/kkb0318/kom/api/v1alpha1"
	komtool "github.com/kkb0318/kom/internal/tool"
	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
					source: &sourcev1.GitRepository{
						ObjectMeta: v1.ObjectMeta{
							Name:      "repo1",
							Namespace: "repo-ns1",
						},
						TypeMeta: v1.TypeMeta{APIVersion: "source.toolkit.fluxcd.io/v1", Kind: "GitRepository"},
						Spec: sourcev1.GitRepositorySpec{
							Interval:  v1.Duration{Duration: time.Minute},
							URL:       "https://example1.com",
							Reference: &sourcev1.GitRepositoryRef{Tag: "1.0.0"},
						},
					},
					ks: &kustomizev1.Kustomization{
						ObjectMeta: v1.ObjectMeta{
							Name:      "repo1",
							Namespace: "repo-ns1",
						},
						TypeMeta: v1.TypeMeta{APIVersion: "kustomize.toolkit.fluxcd.io/v1", Kind: "Kustomization"},
						Spec: kustomizev1.KustomizationSpec{
							Prune: true,
							Path:  "./path1",
							SourceRef: kustomizev1.CrossNamespaceSourceReference{
								Kind:      "GitRepository",
								Name:      "repo1",
								Namespace: "repo-ns1",
							},
						},
					},
				},
				&FluxGit{
					source: &sourcev1.GitRepository{
						ObjectMeta: v1.ObjectMeta{
							Name:      "repo2",
							Namespace: "repo-ns2",
						},
						TypeMeta: v1.TypeMeta{APIVersion: "source.toolkit.fluxcd.io/v1", Kind: "GitRepository"},
						Spec: sourcev1.GitRepositorySpec{
							Interval:  v1.Duration{Duration: time.Minute},
							URL:       "https://example2.com",
							Reference: &sourcev1.GitRepositoryRef{Branch: "main"},
						},
					},
					ks: &kustomizev1.Kustomization{
						ObjectMeta: v1.ObjectMeta{
							Name:      "repo2",
							Namespace: "repo-ns2",
						},
						TypeMeta: v1.TypeMeta{APIVersion: "kustomize.toolkit.fluxcd.io/v1", Kind: "Kustomization"},
						Spec: kustomizev1.KustomizationSpec{
							Prune: true,
							Path:  "./path2",
							SourceRef: kustomizev1.CrossNamespaceSourceReference{
								Kind:      "GitRepository",
								Name:      "repo2",
								Namespace: "repo-ns2",
							},
						},
					},
				},
				&FluxGit{
					source: &sourcev1.GitRepository{
						ObjectMeta: v1.ObjectMeta{
							Name:      "repo3",
							Namespace: "repo-ns3",
						},
						TypeMeta: v1.TypeMeta{APIVersion: "source.toolkit.fluxcd.io/v1", Kind: "GitRepository"},
						Spec: sourcev1.GitRepositorySpec{
							Interval:  v1.Duration{Duration: time.Minute},
							URL:       "https://example3.com",
							Reference: &sourcev1.GitRepositoryRef{SemVer: "x.x.x"},
						},
					},
					ks: &kustomizev1.Kustomization{
						ObjectMeta: v1.ObjectMeta{
							Name:      "repo3",
							Namespace: "repo-ns3",
						},
						TypeMeta: v1.TypeMeta{APIVersion: "kustomize.toolkit.fluxcd.io/v1", Kind: "Kustomization"},
						Spec: kustomizev1.KustomizationSpec{
							Prune: true,
							Path:  "./path3",
							SourceRef: kustomizev1.CrossNamespaceSourceReference{
								Kind:      "GitRepository",
								Name:      "repo3",
								Namespace: "repo-ns3",
							},
						},
					},
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
