package controller

import (
	komv1alpha1 "github.com/kkb0318/kom/api/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type expected struct {
	source types.NamespacedName
	charts []types.NamespacedName
}

var _ = Describe("OperatorManager controller", func() {
	Context("OperatorManager controller test", func() {
		It("should successfully reconcile a custom resource for kom", func() {
			komName := "test-kom-argo"
			kom := createKom(komName)
			kom.Spec = komv1alpha1.OperatorManagerSpec{
				Cleanup: true,
				Resource: komv1alpha1.Resource{
					Helm: []komv1alpha1.Helm{
						{
							Name: "helmrepo1",
							Url:  "https://helm.github.io/examples",
							Charts: []komv1alpha1.Chart{
								{
									Name:    "hello-world",
									Version: "x.x.x",
								},
							},
						},
						{
							Name: "helmrepo2",
							Url:  "https://stefanprodan.github.io/podinfo",
							Charts: []komv1alpha1.Chart{
								{
									Name:    "podinfo",
									Version: "x.x.x",
								},
							},
						},
					},
					Git: []komv1alpha1.Git{
						{
							Name: "gitrepo1",
							Url:  "https://github.com/operator-framework/operator-sdk",
							Path: "testdata/helm/memcached-operator/config/default",
							Reference: komv1alpha1.GitReference{
								Type:  komv1alpha1.GitTag,
								Value: "v1.33.0",
							},
						},
					},
				},
			}
			typeNamespaceName := types.NamespacedName{Name: komName, Namespace: testNamespace}

			expectedHelmResources := []expected{
				{
					source: types.NamespacedName{
						Name:      "helmrepo1",
						Namespace: "kom-system",
					},
					charts: []types.NamespacedName{
						{
							Name:      "hello-world",
							Namespace: "kom-system",
						},
					},
				},
				{
					source: types.NamespacedName{
						Name:      "helmrepo2",
						Namespace: "kom-system",
					},
					charts: []types.NamespacedName{
						{
							Name:      "podinfo",
							Namespace: "kom-system",
						},
					},
				},
			}

			expectedGitResources := []expected{
				{
					source: types.NamespacedName{
						Name:      "gitrepo1",
						Namespace: "kom-system",
					},
					charts: []types.NamespacedName{
						{
							Name:      "gitrepo1",
							Namespace: "kom-system",
						},
					},
				},
			}

			err := k8sClient.Create(ctx, kom)
			Expect(err).To(Not(HaveOccurred()))

			By("Checking if the custom resource was successfully created")
			Eventually(func() error {
				found := &komv1alpha1.OperatorManager{}
				return k8sClient.Get(ctx, typeNamespaceName, found)
			}, timeout).Should(Succeed())

			By("Reconciling the custom resource created")
			komReconciler := &OperatorManagerReconciler{
				Client: k8sClient,
				Scheme: k8sClient.Scheme(),
			}
			// add finalizer
			_, err = komReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespaceName,
			})
			Expect(err).To(Not(HaveOccurred()))
			_, err = komReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespaceName,
			})
			Expect(err).To(Not(HaveOccurred()))

			By("Checking if Resources were successfully created in the reconciliation")
			for _, expected := range expectedHelmResources {
				checkExist(expected.source, helmRepo)
				for _, fetcher := range expected.charts {
					checkExist(fetcher, helmRelease)
				}
			}
			for _, expected := range expectedGitResources {
				checkExist(expected.source, gitRepo)
				for _, fetcher := range expected.charts {
					checkExist(fetcher, kustomization)
				}
			}

			By("Checking garbage collect of partially deletion")
			k8sClient.Get(ctx, typeNamespaceName, kom)
			Eventually(func() error {
				// delete resource[1]
				kom.Spec.Resource.Helm = []komv1alpha1.Helm{kom.Spec.Resource.Helm[0]}
				return k8sClient.Update(ctx, kom)
			}, timeout).Should(Succeed())
			_, err = komReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespaceName,
			})
			Expect(err).To(Not(HaveOccurred()))
			// resource[0] exists
			checkExist(expectedHelmResources[0].source, helmRepo)
			for _, fetcher := range expectedHelmResources[0].charts {
				checkExist(fetcher, helmRelease)
			}
			// resource[1] does not exist
			checkNoExist(expectedHelmResources[1].source, helmRepo)
			for _, fetcher := range expectedHelmResources[1].charts {
				checkNoExist(fetcher, helmRelease)
			}

			By("removing the custom resource for the Kind")
			Eventually(func() error {
				return k8sClient.Delete(ctx, kom)
			}, timeout).Should(Succeed())
			_, err = komReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespaceName,
			})
			Expect(err).To(Not(HaveOccurred()))
			checkNoExist(typeNamespaceName, operatorManager)

			By("Checking if Resources were successfully deleted in the reconciliation")
			for _, expected := range expectedHelmResources {
				checkNoExist(expected.source, helmRepo)
				for _, fetcher := range expected.charts {
					checkNoExist(fetcher, helmRelease)
				}
			}
			for _, expected := range expectedGitResources {
				checkNoExist(expected.source, gitRepo)
				for _, fetcher := range expected.charts {
					checkNoExist(fetcher, kustomization)
				}
			}
		})
	})
})
