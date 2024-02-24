package argo

import (
	argoapi "github.com/kkb0318/argo-cd-api/api"
	argov1alpha1 "github.com/kkb0318/argo-cd-api/api/v1alpha1"
	komv1alpha1 "github.com/kkb0318/kom/api/v1alpha1"
	komtool "github.com/kkb0318/kom/internal/tool"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ArgoHelm struct {
	source []*corev1.Secret
	helm   []*argov1alpha1.Application
}

func (f *ArgoHelm) Repositories() []client.Object {
	objs := make([]client.Object, len(f.source))
	for i, source := range f.source {
		objs[i] = source
	}
	return objs
}

func (f *ArgoHelm) Charts() []client.Object {
	objs := make([]client.Object, len(f.helm))
	for i, helm := range f.helm {
		objs[i] = helm
	}
	return objs
}

func NewArgoHelmList(objs []komv1alpha1.Helm) ([]komtool.Resource, error) {
	helmList := make([]komtool.Resource, len(objs))
	for i, obj := range objs {
		helm, err := NewArgoHelm(obj)
		if err != nil {
			return nil, err
		}
		helmList[i] = helm
	}
	return helmList, nil
}

func NewArgoHelm(obj komv1alpha1.Helm) (*ArgoHelm, error) {
	repoName := obj.Name
	var namespace string
	if obj.Namespace == "" {
		namespace = komv1alpha1.ArgoCDDefaultNamespace
	} else {
		namespace = obj.Namespace
	}
	repoUrl := obj.Url
	charts := obj.Charts
	secrets := make([]*corev1.Secret, len(charts))
	apps := make([]*argov1alpha1.Application, len(charts))
	for i, chart := range charts {
		secret := &corev1.Secret{
			ObjectMeta: v1.ObjectMeta{
				Name:      repoName,
				Namespace: namespace,
			},
			TypeMeta: v1.TypeMeta{
				APIVersion: corev1.SchemeGroupVersion.String(),
				Kind:       "Secret",
			},
			StringData: map[string]string{
				"name":    chart.Name,
				"type":    "helm",
				"url":     repoUrl,
				"project": "default",
			},
		}
		chartNs := namespace
		app := &argov1alpha1.Application{
			ObjectMeta: v1.ObjectMeta{
				Name:      chart.Name,
				Namespace: chartNs,
			},
			TypeMeta: v1.TypeMeta{
				APIVersion: argov1alpha1.SchemeGroupVersion.String(),
				Kind:       argoapi.ApplicationKind,
			},
			Spec: argov1alpha1.ApplicationSpec{
				Source: &argov1alpha1.ApplicationSource{
					Chart:          chart.Name,
					TargetRevision: chart.Version,
					RepoURL:        repoUrl,
				},
				Destination: argov1alpha1.ApplicationDestination{
					Namespace: chartNs,
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
		apps[i] = app
		secrets[i] = secret
	}
	f := &ArgoHelm{
		source: secrets,
		helm:   apps,
	}
	return f, nil
}
