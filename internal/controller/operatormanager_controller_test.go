package controller

import (
	helmv1 "github.com/fluxcd/helm-controller/api/v2beta2"
	sourcev1 "github.com/fluxcd/source-controller/api/v1beta2"
	komv1alpha1 "github.com/kkb0318/kom/api/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type expected struct {
	source   types.NamespacedName
	fetchers []types.NamespacedName
}

var _ = Describe("OperatorManager controller", func() {
	Context("OperatorManager controller test", func() {
		It("should successfully reconcile a custom resource for kom", func() {
			komName := "test-kom"
			kom := createKom(komName)
			kom.Spec = komv1alpha1.OperatorManagerSpec{
				Prune: true,
				Resource: komv1alpha1.Resource{
					Helm: []komv1alpha1.Helm{
						{
							Name: "name1",
							// Namespace: "",
							Url: "https://test.example.com",
							Charts: []komv1alpha1.Chart{
								{
									Name: "chart1",
									// Namespace: "",
									Version: "x.x.x",
								},
							},
						},
					},
				},
			}
			typeNamespaceName := types.NamespacedName{Name: komName, Namespace: testNamespace}

			expected := expected{
				source: types.NamespacedName{
					Name:      "",
					Namespace: "",
				},
				fetchers: []types.NamespacedName{
					{
						Name:      "",
						Namespace: "",
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
			Eventually(func() error {
				found := &sourcev1.HelmRepository{}
				return k8sClient.Get(ctx, expected.source, found)
			}, timeout).Should(Succeed())
			for _, fetcher := range expected.fetchers {
				Eventually(func() error {
					found := &helmv1.HelmRelease{}
					return k8sClient.Get(ctx, fetcher, found)
				}, timeout).Should(Succeed())
			}
		})
	})
})
