package flux

import (
	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1"
	sourcev1 "github.com/fluxcd/source-controller/api/v1"
	komv1alpha1 "github.com/kkb0318/kom/api/v1alpha1"
	komtool "github.com/kkb0318/kom/internal/tool"
	"github.com/kkb0318/kom/internal/tool/flux/manifests"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type FluxGit struct {
	source *sourcev1.GitRepository
	ks     *kustomizev1.Kustomization
}

type GitValues struct {
	FullnameOverride string
}

func (f *FluxGit) Repositories() []client.Object {
	return []client.Object{f.source}
}

func (f *FluxGit) Charts() []client.Object {
	return []client.Object{f.ks}
}

func NewFluxGitList(objs []komv1alpha1.Git) ([]komtool.Resource, error) {
	gitList := make([]komtool.Resource, len(objs))
	for i, obj := range objs {
		git, err := NewFluxGit(obj)
		if err != nil {
			return nil, err
		}
		gitList[i] = git
	}
	return gitList, nil
}

func NewFluxGit(obj komv1alpha1.Git) (*FluxGit, error) {
	repoName := obj.Name
	var namespace string
	if obj.Namespace == "" {
		namespace = komv1alpha1.DefaultNamespace
	} else {
		namespace = obj.Namespace
	}
	gitrepo, err := manifests.NewGitRepositoryBuilder().
		WithReference(obj.Reference.Type, obj.Reference.Value).
		WithUrl(obj.Url).
		Build(repoName, namespace)
	if err != nil {
		return nil, err
	}
	ks, err := manifests.NewKustomizationBuilder().
		WithReference(gitrepo.GetName(), gitrepo.GetNamespace()).
		WithPath(obj.Path).
		Build(repoName, namespace)
	if err != nil {
		return nil, err
	}
	f := &FluxGit{
		source: gitrepo,
		ks:     ks,
	}
	return f, nil
}
