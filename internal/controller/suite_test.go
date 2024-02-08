/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	helmv1 "github.com/fluxcd/helm-controller/api/v2beta2"
	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1"
	sourcev1 "github.com/fluxcd/source-controller/api/v1"
	sourcev1beta2 "github.com/fluxcd/source-controller/api/v1beta2"
	komv1alpha1 "github.com/kkb0318/kom/api/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var (
	timeout       = time.Second * 10
	cfg           *rest.Config
	k8sClient     client.Client
	testEnv       *envtest.Environment
	testNamespace = "kom-system"
	ctx           = ctrl.SetupSignalHandler()
)

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Controller Suite")
}

var _ = BeforeSuite(func() {
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{
			filepath.Join("..", "..", "config", "crd", "bases"),
			filepath.Join(".", "testdata", "crds"),
		},
		ErrorIfCRDPathMissing: true,
	}

	var err error
	// cfg is defined in this file globally.
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = komv1alpha1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())
	// for flux
	err = sourcev1beta2.SchemeBuilder.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())
	err = helmv1.SchemeBuilder.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())
	err = sourcev1.SchemeBuilder.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())
	err = kustomizev1.SchemeBuilder.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())
	k8sClient.Create(ctx, createNamespace(testNamespace))
})

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	authUser, _ := testEnv.AddUser(envtest.User{Name: "test", Groups: []string{"system:masters"}}, &rest.Config{})
	kubectl, _ := authUser.Kubectl()
	if os.Getenv("DEBUG") == "true" {
		fmt.Printf(`
      You can use the following command to investigate the failure:
      kubectl %s
      
      When you have finished investigation, clean up with the following commands:
      pkill kube-apiserver
      pkill etcd
      rm -rf %s
      `, strings.Join(kubectl.Opts, " "), testEnv.ControlPlane.APIServer.CertDir)
		return
	}
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})

func createKom(name string) *komv1alpha1.OperatorManager {
	return &komv1alpha1.OperatorManager{
		TypeMeta: metav1.TypeMeta{
			Kind:       komv1alpha1.OperatorManagerKind,
			APIVersion: komv1alpha1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: testNamespace,
		},
	}
}

func createNamespace(ns string) *corev1.Namespace {
	return &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{Name: ns},
	}
}

func createHelmRepository(name string) *sourcev1beta2.HelmRepository {
	return &sourcev1beta2.HelmRepository{
		TypeMeta: metav1.TypeMeta{
			Kind:       sourcev1beta2.HelmRepositoryKind,
			APIVersion: sourcev1beta2.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: testNamespace,
		},
	}
}

func createHelmRelease(name string) *helmv1.HelmRelease {
	return &helmv1.HelmRelease{
		TypeMeta: metav1.TypeMeta{
			Kind:       helmv1.HelmReleaseKind,
			APIVersion: helmv1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: testNamespace,
		},
	}
}

func checkExist(expected types.NamespacedName, newFunc func() client.Object) {
	Eventually(func() error {
		found := newFunc()
		return k8sClient.Get(ctx, expected, found)
	}, timeout).Should(Succeed())
}

func checkNoExist(expected types.NamespacedName, newFunc func() client.Object) {
	Eventually(func() error {
		found := newFunc()
		return k8sClient.Get(ctx, expected, found)
	}, timeout).Should(Not(Succeed()))
}

// -----used for assertion-----
func helmRepo() client.Object {
	return &sourcev1beta2.HelmRepository{}
}

func helmRelease() client.Object {
	return &helmv1.HelmRelease{}
}

func gitRepo() client.Object {
	return &sourcev1.GitRepository{}
}

func kustomization() client.Object {
	return &kustomizev1.Kustomization{}
}

func operatorManager() client.Object {
	return &komv1alpha1.OperatorManager{}
}
