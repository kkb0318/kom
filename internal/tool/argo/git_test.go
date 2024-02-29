package argo

import (
	"testing"

	argoapi "github.com/kkb0318/argo-cd-api/api"
	argov1alpha1 "github.com/kkb0318/argo-cd-api/api/v1alpha1"
	komv1alpha1 "github.com/kkb0318/kom/api/v1alpha1"
	komtool "github.com/kkb0318/kom/internal/tool"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
					source: &corev1.Secret{
						ObjectMeta: v1.ObjectMeta{
							Name:      "repo1",
							Namespace: "repo-ns1",
						},
						TypeMeta: v1.TypeMeta{
							APIVersion: corev1.SchemeGroupVersion.String(),
							Kind:       "Secret",
						},
						StringData: map[string]string{
							"type":    "git",
							"url":     "https://example.com",
							"project": "default",
						},
					},
					app: &argov1alpha1.Application{
						ObjectMeta: v1.ObjectMeta{
							Name:      "repo1",
							Namespace: "repo-ns1",
						},
						TypeMeta: v1.TypeMeta{
							APIVersion: argov1alpha1.SchemeGroupVersion.String(),
							Kind:       argoapi.ApplicationKind,
						},
						Spec: argov1alpha1.ApplicationSpec{
							Source: &argov1alpha1.ApplicationSource{
								Path:           "./path1",
								TargetRevision: "1.0.0",
								RepoURL:        "https://example.com",
							},
							Destination: argov1alpha1.ApplicationDestination{
								Namespace: "repo-ns1",
								Server:    "https://kubernetes.default.svc",
							},
							Project: "default",
							SyncPolicy: &argov1alpha1.SyncPolicy{
								Automated: &argov1alpha1.SyncPolicyAutomated{
									Prune:    true,
									SelfHeal: true,
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
