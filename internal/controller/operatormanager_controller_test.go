package controller

import (
	komv1alpha1 "github.com/kkb0318/kom/api/v1alpha1"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var _ = Describe("Memcached controller", func() {
	Context("Memcached controller test", func() {
		It("should successfully reconcile a custom resource for Memcached", func() {
			komName := "test-kom"
			kom := &komv1alpha1.OperatorManager{}
			kom.Name = komName
			kom.Namespace = testNamespace
			kom.Spec = komv1alpha1.OperatorManagerSpec{
				Prune:    true,
				Resource: komv1alpha1.Resource{},
			}
			typeNamespaceName := types.NamespacedName{Name: komName, Namespace: testNamespace}
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

			_, err = komReconciler.Reconcile(ctx, reconcile.Request{
				NamespacedName: typeNamespaceName,
			})
			Expect(err).To(Not(HaveOccurred()))

			By("Checking if Deployment was successfully created in the reconciliation")
			Eventually(func() error {
				found := &appsv1.Deployment{}
				return k8sClient.Get(ctx, typeNamespaceName, found)
			}, timeout).Should(Succeed())
		})
	})
})
