package manifests

import (
	argoapi "github.com/kkb0318/argo-cd-api/api"
	argov1alpha1 "github.com/kkb0318/argo-cd-api/api/v1alpha1"
	komv1alpha1 "github.com/kkb0318/kom/api/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewApplication(chart komv1alpha1.Chart, ns, url string) *argov1alpha1.Application {
	app := &argov1alpha1.Application{
		ObjectMeta: v1.ObjectMeta{
			Name:      chart.Name,
			Namespace: ns,
		},
		TypeMeta: v1.TypeMeta{
			APIVersion: argov1alpha1.SchemeGroupVersion.String(),
			Kind:       argoapi.ApplicationKind,
		},
		Spec: argov1alpha1.ApplicationSpec{
			Source: &argov1alpha1.ApplicationSource{
				Chart:          chart.Name,
				TargetRevision: chart.Version,
				RepoURL:        url,
			},
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
	return app
}
