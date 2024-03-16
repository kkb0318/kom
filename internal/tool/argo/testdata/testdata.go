package testdata

import (
	"path/filepath"
	"runtime"
	"testing"

	argov1alpha1 "github.com/kkb0318/argo-cd-api/api/v1alpha1"
	"github.com/kkb0318/kom/internal/utils"
	corev1 "k8s.io/api/core/v1"
)

type mockSecretBuilder struct {
	name       string
	namespace  string
	stringData map[string]string
}

func NewMockSecretBuilder() *mockSecretBuilder {
	return &mockSecretBuilder{}
}

func (f *mockSecretBuilder) WithName(val string) *mockSecretBuilder {
	f.name = val
	return f
}

func (f *mockSecretBuilder) WithNamespace(val string) *mockSecretBuilder {
	f.namespace = val
	return f
}

func (f *mockSecretBuilder) WithType(val string) *mockSecretBuilder {
	f.stringData["type"] = val
	return f
}

func (f *mockSecretBuilder) WithUrl(val string) *mockSecretBuilder {
	f.stringData["url"] = val
	return f
}

func (f *mockSecretBuilder) Build(t *testing.T) *corev1.Secret {
	baseFilePath := filepath.Join(currentDir(t), "git_secret.yaml")
	Secret := &corev1.Secret{}
	utils.LoadYaml(Secret, baseFilePath)
	if f.name != "" {
		Secret.ObjectMeta.SetName(f.name)
	}
	if f.namespace != "" {
		Secret.ObjectMeta.SetNamespace(f.namespace)
	}
	if f.stringData != nil {
		Secret.StringData = f.stringData
	}
	return Secret
}

type mockApplicationBuilder struct{}

func NewMockApplicationBuilder() *mockApplicationBuilder {
	return &mockApplicationBuilder{}
}

func (f *mockApplicationBuilder) Build(t *testing.T) *argov1alpha1.Application {
	baseFilePath := filepath.Join(currentDir(t), "git_application.yaml")
	application := &argov1alpha1.Application{}
	utils.LoadYaml(application, baseFilePath)
	return application
}

func currentDir(t *testing.T) string {
	t.Helper()
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller() failed to get current file path")
	}
	return filepath.Dir(currentFile)
}
