package manifests

import (
	"errors"
	"time"

	sourcev1 "github.com/fluxcd/source-controller/api/v1"
	komv1alpha1 "github.com/kkb0318/kom/api/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type GitRepositoryBuilder struct {
	ref *sourcev1.GitRepositoryRef
	url string
}

func NewGitRepositoryBuilder() *GitRepositoryBuilder {
	return &GitRepositoryBuilder{}
}

// WithReference set GitRepositoryRef
func (b *GitRepositoryBuilder) WithReference(refType komv1alpha1.GitReferenceType, value string) *GitRepositoryBuilder {
	ref := &sourcev1.GitRepositoryRef{}
	switch refType {
	case komv1alpha1.GitBranch:
		ref.Branch = value
	case komv1alpha1.GitTag:
		ref.Tag = value
	case komv1alpha1.GitSemver:
		ref.SemVer = value
	default:
		return b
	}
	b.ref = ref
	return b
}

func (b *GitRepositoryBuilder) WithUrl(value string) *GitRepositoryBuilder {
	b.url = value
	return b
}

func (b *GitRepositoryBuilder) Build(name, ns string) (*sourcev1.GitRepository, error) {
	if b.ref == nil {
		return nil, errors.New("the 'ref' field is nil and must be provided with a valid reference")
	}
	if b.url == "" {
		return nil, errors.New("the 'url' field is empty. Please specify a valid URL")
	}
	gitrepo := &sourcev1.GitRepository{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		TypeMeta: v1.TypeMeta{
			APIVersion: sourcev1.GroupVersion.String(),
			Kind:       sourcev1.GitRepositoryKind,
		},
		Spec: sourcev1.GitRepositorySpec{
			Interval:  v1.Duration{Duration: time.Minute},
			URL:       b.url,
			Reference: b.ref,
		},
	}
	return gitrepo, nil
}
