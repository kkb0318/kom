package status

import (
	"fmt"

	komv1alpha1 "github.com/kkb0318/kom/api/v1alpha1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

func ToListUnstructed(resources []komv1alpha1.AppliedResource) ([]*unstructured.Unstructured, error) {
	objects := make([]*unstructured.Unstructured, 0)
	for _, r := range resources {
		u, err := ToUnstructured(r)
		if err != nil {
			return nil, err
		}
		objects = append(objects, u)

	}
	return objects, nil
}

// ToUnstructured converts an AppliedResource into an Unstructured object.
// It returns an error if the conversion fails or if the Unstructured object cannot be created.
func ToUnstructured(a komv1alpha1.AppliedResource) (*unstructured.Unstructured, error) {
	gvk := schema.FromAPIVersionAndKind(a.APIVersion, a.Kind)
	// Verify if the GroupVersionKind (GVK) is properly parsed
	if gvk.Group == "" && gvk.Version == "" {
		return nil, fmt.Errorf("failed to parse GroupVersionKind from APIVersion and Kind: %v", gvk)
	}
	// Ensure the resource name is not empty
	if a.Name == "" {
		return nil, fmt.Errorf("resource name is required but was not provided")
	}
	// Ensure the namespace is provided for namespaced resources
	if a.Namespace == "" {
		return nil, fmt.Errorf("namespace is required for namespaced resources but was not provided")
	}
	// Create and populate the Unstructured object
	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(schema.GroupVersionKind{
		Group:   gvk.Group,
		Kind:    gvk.Kind,
		Version: gvk.Version,
	})
	u.SetName(a.Name)
	u.SetNamespace(a.Namespace)
	return u, nil
}

func ToAppliedResource(u unstructured.Unstructured) (*komv1alpha1.AppliedResource, error) {
  a := &komv1alpha1.AppliedResource{}
  a.Name = u.GetName()
  a.Namespace = u.GetNamespace()
  a.Kind = u.GetObjectKind().GroupVersionKind().Kind
  a.APIVersion = u.GetAPIVersion()
  return a, nil
}
