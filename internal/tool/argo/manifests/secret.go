package manifests

import (
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewSecret(name, ns, url, chartName string) *corev1.Secret {
	secret := &corev1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		TypeMeta: v1.TypeMeta{
			APIVersion: corev1.SchemeGroupVersion.String(),
			Kind:       "Secret",
		},
		StringData: map[string]string{
			"name":    chartName,
			"type":    "helm",
			"url":     url,
			"project": "default",
		},
	}
	return secret
}
