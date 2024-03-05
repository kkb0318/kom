package controller

import (
	"fmt"

	komv1alpha1 "github.com/kkb0318/kom/api/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ = Describe("OperatorManager controller", func() {
	Context("OperatorManager controller test", func() {
		It("no garbage collect test. cleanup: false", func() {
			nameSuffix := "-no-cleanup"
			komName := fmt.Sprintf("test-kom%s", nameSuffix)
			kom := createKom(komName)
			kom.Spec = komv1alpha1.OperatorManagerSpec{
				Cleanup: false,
				Tool:    komv1alpha1.FluxCDTool,
				Resource: komv1alpha1.Resource{
					Helm: []komv1alpha1.Helm{
						{
							Name: fmt.Sprintf("helmrepo1%s", nameSuffix),
							Url:  "https://helm.github.io/examples",
							Charts: []komv1alpha1.Chart{
								{
									Name:    "hello-world",
									Version: "x.x.x",
								},
							},
						},
						{
							Name: fmt.Sprintf("helmrepo2%s", nameSuffix),
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
							Name: fmt.Sprintf("gitrepo1%s", nameSuffix),
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
						Name:      fmt.Sprintf("helmrepo1%s", nameSuffix),
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
						Name:      fmt.Sprintf("helmrepo2%s", nameSuffix),
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
						Name:      fmt.Sprintf("gitrepo1%s", nameSuffix),
						Namespace: "kom-system",
					},
					charts: []types.NamespacedName{
						{
							Name:      fmt.Sprintf("gitrepo1%s", nameSuffix),
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

			By("Checking no garbage collect")
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

			By("removing the custom resource for the Kind")
			Eventually(func() error {
				return k8sClient.Delete(ctx, kom)
			}, timeout).Should(Succeed())
			_, err = komReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespaceName,
			})
			Expect(err).To(Not(HaveOccurred()))
			checkNoExist(typeNamespaceName, operatorManager)

			By("Checking no garbage collect")
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
		})
	})
})
