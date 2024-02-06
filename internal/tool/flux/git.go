package flux

import (
	"time"

	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1"
	sourcev1 "github.com/fluxcd/source-controller/api/v1"
	komv1alpha1 "github.com/kkb0318/kom/api/v1alpha1"
	komtool "github.com/kkb0318/kom/internal/tool"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type FluxGit struct {
	source *sourcev1.GitRepository
	ks     *kustomizev1.Kustomization
}

type GitValues struct {
	FullnameOverride string
}

func (f *FluxGit) Repository() client.Object {
	return f.source
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
	repoUrl := obj.Url
	value := obj.Reference.Value
	var namespace string
	if obj.Namespace == "" {
		namespace = komv1alpha1.DefaultNamespace
	} else {
		namespace = obj.Namespace
	}
	ref := &sourcev1.GitRepositoryRef{}
	switch obj.Reference.Type {
	case komv1alpha1.GitBranch:
		ref.Branch = value
	case komv1alpha1.GitTag:
		ref.Tag = value
	case komv1alpha1.GitSemver:
		ref.SemVer = value
	}
	gitrepo := &sourcev1.GitRepository{
		ObjectMeta: v1.ObjectMeta{
			Name:      repoName,
			Namespace: namespace,
		},
		TypeMeta: v1.TypeMeta{
			APIVersion: sourcev1.GroupVersion.String(),
			Kind:       sourcev1.GitRepositoryKind,
		},
		Spec: sourcev1.GitRepositorySpec{
			Interval:  v1.Duration{Duration: time.Minute},
			URL:       repoUrl,
			Reference: ref,
		},
	}
	ks := &kustomizev1.Kustomization{
		ObjectMeta: v1.ObjectMeta{
			Name:      repoName,
			Namespace: namespace,
		},
		TypeMeta: v1.TypeMeta{
			APIVersion: kustomizev1.GroupVersion.String(),
			Kind:       kustomizev1.KustomizationKind,
		},
		Spec: kustomizev1.KustomizationSpec{
			Prune: true,
			Path:  obj.Path,
			SourceRef: kustomizev1.CrossNamespaceSourceReference{
				Kind:      sourcev1.GitRepositoryKind,
				Name:      gitrepo.GetName(),
				Namespace: gitrepo.GetNamespace(),
			},
		},
	}
	f := &FluxGit{
		source: gitrepo,
		ks:     ks,
	}
	return f, nil
}
