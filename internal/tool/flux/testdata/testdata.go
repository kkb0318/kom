package testdata

import (
	"path/filepath"
	"runtime"
	"testing"

	helmv2 "github.com/fluxcd/helm-controller/api/v2"
	kustomizev1 "github.com/fluxcd/kustomize-controller/api/v1"
	sourcev1 "github.com/fluxcd/source-controller/api/v1"
	"github.com/kkb0318/kom/internal/utils"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
)

type mockGitRepositoryBuilder struct {
	name      string
	namespace string
	url       string
	ref       *sourcev1.GitRepositoryRef
}

func NewMockGitRepositoryBuilder() *mockGitRepositoryBuilder {
	return &mockGitRepositoryBuilder{}
}

func (f *mockGitRepositoryBuilder) WithName(val string) *mockGitRepositoryBuilder {
	f.name = val
	return f
}

func (f *mockGitRepositoryBuilder) WithNamespace(val string) *mockGitRepositoryBuilder {
	f.namespace = val
	return f
}

func (f *mockGitRepositoryBuilder) WithUrl(val string) *mockGitRepositoryBuilder {
	f.url = val
	return f
}

func (f *mockGitRepositoryBuilder) WithRef(val *sourcev1.GitRepositoryRef) *mockGitRepositoryBuilder {
	f.ref = val
	return f
}

func (f *mockGitRepositoryBuilder) Build(t *testing.T, testdataFileName string) *sourcev1.GitRepository {
	baseFilePath := filepath.Join(currentDir(t), testdataFileName)
	gitrepo := &sourcev1.GitRepository{}
	utils.LoadYaml(gitrepo, baseFilePath)
	if f.name != "" {
		gitrepo.SetName(f.name)
	}
	if f.namespace != "" {
		gitrepo.SetNamespace(f.namespace)
	}
	if f.ref != nil {
		gitrepo.Spec.Reference = f.ref
	}
	if f.url != "" {
		gitrepo.Spec.URL = f.url
	}
	return gitrepo
}

type mockKustomizationBuilder struct {
	name      string
	namespace string
	path      string
	url       string
	ref       *kustomizev1.CrossNamespaceSourceReference
}

func NewMockKustomizationBuilder() *mockKustomizationBuilder {
	return &mockKustomizationBuilder{}
}

func (f *mockKustomizationBuilder) WithName(val string) *mockKustomizationBuilder {
	f.name = val
	return f
}

func (f *mockKustomizationBuilder) WithNamespace(val string) *mockKustomizationBuilder {
	f.namespace = val
	return f
}

func (f *mockKustomizationBuilder) WithPath(val string) *mockKustomizationBuilder {
	f.path = val
	return f
}

func (f *mockKustomizationBuilder) WithRef(val *kustomizev1.CrossNamespaceSourceReference) *mockKustomizationBuilder {
	f.ref = val
	return f
}

func (f *mockKustomizationBuilder) Build(t *testing.T, testdataFileName string) *kustomizev1.Kustomization {
	baseFilePath := filepath.Join(currentDir(t), testdataFileName)
	ks := &kustomizev1.Kustomization{}
	utils.LoadYaml(ks, baseFilePath)
	if f.name != "" {
		ks.SetName(f.name)
	}
	if f.namespace != "" {
		ks.SetNamespace(f.namespace)
	}
	if f.path != "" {
		ks.Spec.Path = f.path
	}
	if f.ref != nil {
		ks.Spec.SourceRef = *f.ref
	}
	return ks
}

type mockHelmRepositoryBuilder struct {
	name      string
	namespace string
	url       string
}

func NewMockHelmRepositoryBuilder() *mockHelmRepositoryBuilder {
	return &mockHelmRepositoryBuilder{}
}

func (f *mockHelmRepositoryBuilder) WithName(val string) *mockHelmRepositoryBuilder {
	f.name = val
	return f
}

func (f *mockHelmRepositoryBuilder) WithNamespace(val string) *mockHelmRepositoryBuilder {
	f.namespace = val
	return f
}

func (f *mockHelmRepositoryBuilder) WithUrl(val string) *mockHelmRepositoryBuilder {
	f.url = val
	return f
}

func (f *mockHelmRepositoryBuilder) Build(t *testing.T, testdataFileName string) *sourcev1.HelmRepository {
	baseFilePath := filepath.Join(currentDir(t), testdataFileName)
	helmrepo := &sourcev1.HelmRepository{}
	utils.LoadYaml(helmrepo, baseFilePath)
	if f.name != "" {
		helmrepo.SetName(f.name)
	}
	if f.namespace != "" {
		helmrepo.SetNamespace(f.namespace)
	}
	if f.url != "" {
		helmrepo.Spec.URL = f.url
	}
	return helmrepo
}

type mockHelmReleaseBuilder struct {
	name      string
	namespace string
	values    *apiextensionsv1.JSON
}

func NewMockHelmReleaseBuilder() *mockHelmReleaseBuilder {
	return &mockHelmReleaseBuilder{}
}

func (f *mockHelmReleaseBuilder) WithName(val string) *mockHelmReleaseBuilder {
	f.name = val
	return f
}

func (f *mockHelmReleaseBuilder) WithNamespace(val string) *mockHelmReleaseBuilder {
	f.namespace = val
	return f
}

func (f *mockHelmReleaseBuilder) WithValues(values *apiextensionsv1.JSON) *mockHelmReleaseBuilder {
	f.values = values
	return f
}

func (f *mockHelmReleaseBuilder) Build(t *testing.T, testdataFileName string) *helmv2.HelmRelease {
	baseFilePath := filepath.Join(currentDir(t), testdataFileName)
	hr := &helmv2.HelmRelease{}
	utils.LoadYaml(hr, baseFilePath)
	if f.name != "" {
		hr.SetName(f.name)
	}
	if f.namespace != "" {
		hr.SetNamespace(f.namespace)
	}
	if f.values != nil {
		hr.Spec.Values = f.values
	}

	return hr
}

func currentDir(t *testing.T) string {
	t.Helper()
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller() failed to get current file path")
	}
	return filepath.Dir(currentFile)
}
