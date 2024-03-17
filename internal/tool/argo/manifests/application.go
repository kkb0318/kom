package manifests

import (
	"fmt"

	argoapi "github.com/kkb0318/argo-cd-api/api"
	argov1alpha1 "github.com/kkb0318/argo-cd-api/api/v1alpha1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type ApplicationBuilder struct {
	source *argov1alpha1.ApplicationSource
}

func NewApplicationBuilder() *ApplicationBuilder {
	return &ApplicationBuilder{}
}

func (b *ApplicationBuilder) WithHelm(name, version, url string) *ApplicationBuilder {
	source := &argov1alpha1.ApplicationSource{
		Chart:          name,
		TargetRevision: version,
		RepoURL:        url,
	}
	b.source = source
	return b
}

func (b *ApplicationBuilder) WithHelmValues(values *apiextensionsv1.JSON) *ApplicationBuilder {
	if values == nil {
		return b
	}
	b.source.Helm.ValuesObject = &runtime.RawExtension{
		Raw: values.Raw,
	}
	return b
}

func (b *ApplicationBuilder) WithGit(path, version, url string) *ApplicationBuilder {
	source := &argov1alpha1.ApplicationSource{
		Path:           path,
		TargetRevision: version,
		RepoURL:        url,
	}
	b.source = source
	return b
}

func (b *ApplicationBuilder) Build(name, ns string) (*argov1alpha1.Application, error) {
	if b.source == nil {
		return nil, fmt.Errorf("argocd ApplicationSource is empty")
	}
	app := &argov1alpha1.Application{
		ObjectMeta: v1.ObjectMeta{
			Name:      name,
			Namespace: ns,
		},
		TypeMeta: v1.TypeMeta{
			APIVersion: argov1alpha1.SchemeGroupVersion.String(),
			Kind:       argoapi.ApplicationKind,
		},
		Spec: argov1alpha1.ApplicationSpec{
			Source: b.source,
			Destination: argov1alpha1.ApplicationDestination{
				Namespace: ns,
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
	}
	return app, nil
}
