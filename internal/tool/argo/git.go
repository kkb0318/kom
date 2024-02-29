package argo

import (
	argov1alpha1 "github.com/kkb0318/argo-cd-api/api/v1alpha1"
	komv1alpha1 "github.com/kkb0318/kom/api/v1alpha1"
	komtool "github.com/kkb0318/kom/internal/tool"
	"github.com/kkb0318/kom/internal/tool/argo/manifests"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ArgoGit struct {
	source *corev1.Secret
	app    *argov1alpha1.Application
}

func (f *ArgoGit) Repositories() []client.Object {
	return []client.Object{f.source}
}

func (f *ArgoGit) Charts() []client.Object {
	return []client.Object{f.app}
}

func NewArgoGitList(objs []komv1alpha1.Git) ([]komtool.Resource, error) {
	gitList := make([]komtool.Resource, len(objs))
	for i, obj := range objs {
		git, err := NewArgoGit(obj)
		if err != nil {
			return nil, err
		}
		gitList[i] = git
	}
	return gitList, nil
}

func NewArgoGit(obj komv1alpha1.Git) (*ArgoGit, error) {
	var namespace string
	if obj.Namespace == "" {
		namespace = komv1alpha1.ArgoCDDefaultNamespace
	} else {
		namespace = obj.Namespace
	}
	secret, err := manifests.NewSecretBuilder().
		WithGit(obj.Url).
		Build(obj.Name, namespace)
	if err != nil {
		return nil, err
	}
	git, err := manifests.NewApplicationBuilder().
		WithGit(obj.Path, obj.Reference.Value, obj.Url).
		Build(obj.Name, namespace)
	if err != nil {
		return nil, err
	}
	f := &ArgoGit{
		source: secret,
		app:    git,
	}
	return f, nil
}
