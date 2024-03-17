package manifests

import (
	"errors"

	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1"
	sourcev1 "github.com/fluxcd/source-controller/api/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type KustomizationBuilder struct {
	ref  *kustomizev1.CrossNamespaceSourceReference
	path string
}

func NewKustomizationBuilder() *KustomizationBuilder {
	return &KustomizationBuilder{
		path: "./",
	}
}

// WithReference set KustomizationRef
func (b *KustomizationBuilder) WithReference(name, namespace string) *KustomizationBuilder {
	ref := &kustomizev1.CrossNamespaceSourceReference{
		Kind:      sourcev1.GitRepositoryKind,
		Name:      name,
		Namespace: namespace,
	}
	b.ref = ref
	return b
}

func (b *KustomizationBuilder) WithPath(value string) *KustomizationBuilder {
	b.path = value
	return b
}

func (b *KustomizationBuilder) Build(name, ns string) (*kustomizev1.Kustomization, error) {
	if b.ref == nil {
		return nil, errors.New("the 'ref' field is nil and must be provided with a valid reference")
	}
	gitrepo := &kustomizev1.Kustomization{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		TypeMeta: v1.TypeMeta{
			APIVersion: kustomizev1.GroupVersion.String(),
			Kind:       kustomizev1.KustomizationKind,
		},
		Spec: kustomizev1.KustomizationSpec{
			Prune: true,
			Path:  b.path,
			SourceRef: *b.ref,
		},
	}
	return gitrepo, nil
}
