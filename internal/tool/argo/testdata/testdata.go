package testdata

import (
	"path/filepath"
	"runtime"
	"testing"

	argov1alpha1 "github.com/kkb0318/argo-cd-api/api/v1alpha1"
	apiextensionsv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"github.com/kkb0318/kom/internal/utils"
	corev1 "k8s.io/api/core/v1"
	k8sruntime "k8s.io/apimachinery/pkg/runtime"
)

type mockSecretBuilder struct {
	name       string
	namespace  string
	stringData map[string]string
}

func NewMockSecretBuilder() *mockSecretBuilder {
	return &mockSecretBuilder{
    stringData: map[string]string{},
  }
}

func (f *mockSecretBuilder) WithName(val string) *mockSecretBuilder {
	f.name = val
	return f
}

func (f *mockSecretBuilder) WithNamespace(val string) *mockSecretBuilder {
	f.namespace = val
	return f
}

func (f *mockSecretBuilder) WithChartName(val string) *mockSecretBuilder {
	f.stringData["name"] = val
	return f
}
func (f *mockSecretBuilder) WithUrl(val string) *mockSecretBuilder {
	f.stringData["url"] = val
	return f
}

func (f *mockSecretBuilder) Build(t *testing.T, testdataFileName string) *corev1.Secret {
	baseFilePath := filepath.Join(currentDir(t), testdataFileName)
	Secret := &corev1.Secret{}
	utils.LoadYaml(Secret, baseFilePath)
	if f.name != "" {
		Secret.ObjectMeta.SetName(f.name)
	}
	if f.namespace != "" {
		Secret.ObjectMeta.SetNamespace(f.namespace)
	}
  for k, v := range f.stringData {
    Secret.StringData[k] = v
  }
	return Secret
}

type mockApplicationBuilder struct{
	name           string
	namespace      string
	destNamespace  string
	chartName      string
	version        string
	values    *apiextensionsv1.JSON
}

func NewMockApplicationBuilder() *mockApplicationBuilder {
	return &mockApplicationBuilder{}
}

func (f *mockApplicationBuilder) WithName(val string) *mockApplicationBuilder {
	f.name = val
	return f
}

func (f *mockApplicationBuilder) WithNamespace(val string) *mockApplicationBuilder {
	f.namespace = val
	return f
}

func (f *mockApplicationBuilder) WithDestNamespace(val string) *mockApplicationBuilder {
	f.destNamespace = val
	return f
}

func (f *mockApplicationBuilder) WithChartName(val string) *mockApplicationBuilder {
	f.chartName = val
	return f
}

func (f *mockApplicationBuilder) WithVersion(val string) *mockApplicationBuilder {
	f.version = val
	return f
}

func (f *mockApplicationBuilder) WithValues(values *apiextensionsv1.JSON) *mockApplicationBuilder {
	f.values = values
	return f
}

func (f *mockApplicationBuilder) Build(t *testing.T, testdataFileName string) *argov1alpha1.Application {
	baseFilePath := filepath.Join(currentDir(t), testdataFileName)
	application := &argov1alpha1.Application{}
	utils.LoadYaml(application, baseFilePath)
	if f.name != "" {
		application.ObjectMeta.SetName(f.name)
	}
	if f.namespace != "" {
		application.ObjectMeta.SetNamespace(f.namespace)
	}
	if f.destNamespace != "" {
		application.Spec.Destination.Namespace = f.destNamespace
	}
	if f.chartName != "" {
		application.Spec.Source.Chart = f.chartName
	}
	if f.version != "" {
		application.Spec.Source.TargetRevision = f.version
	}
	if f.values != nil {
	  if application.Spec.Source.Helm == nil {
	  	application.Spec.Source.Helm = &argov1alpha1.ApplicationSourceHelm{}
	  }
	  application.Spec.Source.Helm.ValuesObject = &k8sruntime.RawExtension{
	  	Raw: f.values.Raw,
	  }
	}
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
