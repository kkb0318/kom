package manifests

import (
	"errors"
	"strings"
	"time"

	sourcev1 "github.com/fluxcd/source-controller/api/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type HelmRepositoryBuilder struct {
	url      string
	interval v1.Duration
}

func NewHelmRepositoryBuilder() *HelmRepositoryBuilder {
	return &HelmRepositoryBuilder{
		interval: v1.Duration{Duration: time.Minute},
	}
}

func (b *HelmRepositoryBuilder) WithUrl(value string) *HelmRepositoryBuilder {
	b.url = value
	return b
}

func (b *HelmRepositoryBuilder) Build(name, ns string) (*sourcev1.HelmRepository, error) {
	if b.url == "" {
		return nil, errors.New("the 'url' field is empty. Please specify a valid URL")
	}
	helmrepo := &sourcev1.HelmRepository{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		TypeMeta: v1.TypeMeta{
			APIVersion: sourcev1.GroupVersion.String(),
			Kind:       sourcev1.HelmRepositoryKind,
		},
		Spec: sourcev1.HelmRepositorySpec{
			Type:     repositoryType(b.url),
			Interval: b.interval,
			URL:      b.url,
		},
	}
	return helmrepo, nil
}

func repositoryType(url string) string {
	if strings.HasPrefix(url, "oci:") {
		return "oci"
	}
	return "default"
}
