package manifests

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type SecretBuilder struct {
	stringData map[string]string
}

func NewSecretBuilder() *SecretBuilder {
	return &SecretBuilder{}
}

func (b *SecretBuilder) WithGit(url string) *SecretBuilder {
	b.stringData = map[string]string{
		"type":    "git",
		"url":     url,
		"project": "default",
	}
	return b
}

func (b *SecretBuilder) WithHelm(chartName, url string) *SecretBuilder {
	b.stringData = map[string]string{
		"name":    chartName,
		"type":    "helm",
		"url":     url,
		"project": "default",
	}
	return b
}

func (b *SecretBuilder) Build(name, ns string) (*corev1.Secret, error) {
	if b.stringData == nil {
		return nil, fmt.Errorf("argocd Secret.StringData is empty")
	}
	secret := &corev1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: ns,
			Labels: map[string]string{
				"argocd.argoproj.io/secret-type": "repository",
			},
		},
		TypeMeta: v1.TypeMeta{
			APIVersion: corev1.SchemeGroupVersion.String(),
			Kind:       "Secret",
		},
		Type:       corev1.SecretTypeOpaque,
		StringData: b.stringData,
	}
	return secret, nil
}
