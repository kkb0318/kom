package status

import (
	"testing"

	komv1alpha1 "github.com/kkb0318/kom/api/v1alpha1"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestStatus_Diff(t *testing.T) {
	tests := []struct {
		name         string
		oldList      komv1alpha1.AppliedResourceList
		newList      komv1alpha1.AppliedResourceList
		expectedDiff []*unstructured.Unstructured
		expectErr  bool
	}{
		{
			name: "Diff with one removed resource",
			oldList: komv1alpha1.AppliedResourceList{
				{Name: "name1", Namespace: "ns1", Kind: "Kind1", APIVersion: "v1"},
				{Name: "name2", Namespace: "ns2", Kind: "Kind2", APIVersion: "v2"},
			},
			newList: komv1alpha1.AppliedResourceList{
				{Name: "name1", Namespace: "ns1", Kind: "Kind1", APIVersion: "v1"},
			},
			expectedDiff: []*unstructured.Unstructured{
				{
					Object: map[string]interface{}{
						"metadata":   map[string]interface{}{"name": "name2", "namespace": "ns2"},
						"apiVersion": "v2", "kind": "Kind2",
					},
				},
			},
		},
		{
			name: "Diff with multiple changes",
			oldList: komv1alpha1.AppliedResourceList{
				{Name: "name1", Namespace: "ns1", Kind: "Kind1", APIVersion: "v1"},
				{Name: "name2", Namespace: "ns2", Kind: "Kind2", APIVersion: "v2"},
				{Name: "name3", Namespace: "ns3", Kind: "Kind3", APIVersion: "v3"},
			},
			newList: komv1alpha1.AppliedResourceList{
				{Name: "name2", Namespace: "ns2", Kind: "Kind2", APIVersion: "v2"},
				{Name: "xxxxx", Namespace: "xxxxx", Kind: "xxxxx", APIVersion: "xxxxx"},
			},
			expectedDiff: []*unstructured.Unstructured{
				{
					Object: map[string]interface{}{
						"metadata":   map[string]interface{}{"name": "name1", "namespace": "ns1"},
						"apiVersion": "v1", "kind": "Kind1",
					},
				},
				{
					Object: map[string]interface{}{
						"metadata":   map[string]interface{}{"name": "name3", "namespace": "ns3"},
						"apiVersion": "v3", "kind": "Kind3",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := Diff(tt.oldList, tt.newList)
			if tt.expectErr {
				assert.Error(t, err, "")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedDiff, actual)
			}
		})
	}
}

func TestStatus_ToListUnstructured(t *testing.T) {
	testCases := []struct {
		name         string
		resourceList komv1alpha1.AppliedResourceList
		expected     []*unstructured.Unstructured
		expectErr    bool
	}{
		{
			name: "Convert valid resource list to unstructured",
			resourceList: komv1alpha1.AppliedResourceList{
				{Name: "name1", Namespace: "ns1", Kind: "Kind1", APIVersion: "v1"},
				{Name: "name2", Namespace: "ns2", Kind: "Kind2", APIVersion: "v2"},
			},
			expected: []*unstructured.Unstructured{
				{
					Object: map[string]interface{}{
						"metadata":   map[string]interface{}{"name": "name1", "namespace": "ns1"},
						"apiVersion": "v1", "kind": "Kind1",
					},
				},
				{
					Object: map[string]interface{}{
						"metadata":   map[string]interface{}{"name": "name2", "namespace": "ns2"},
						"apiVersion": "v2", "kind": "Kind2",
					},
				},
			},
			expectErr: false,
		},
		{
			name: "Fail conversion when resource lacks kind",
			resourceList: komv1alpha1.AppliedResourceList{
				{Name: "name1", Namespace: "ns1", Kind: "Kind1", APIVersion: "v1"},
				// Omitting Kind to simulate a malformed resource
				{Name: "name2", Namespace: "ns2", APIVersion: "v2"},
			},
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := ToListUnstructured(tc.resourceList)
			if tc.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, result)
			}
		})
	}
}


func TestStatus_ToUnstructured(t *testing.T) {
    testCases := []struct {
        name string 
        resource    komv1alpha1.AppliedResource
        expected    *unstructured.Unstructured
        expectErr bool 
    }{
        {
            name: "Convert a valid resource to unstructured",
            resource:    komv1alpha1.AppliedResource{Name: "name1", Namespace: "ns1", Kind: "Kind1", APIVersion: "v1"},
            expected: &unstructured.Unstructured{
                Object: map[string]interface{}{
                    "metadata":   map[string]interface{}{"name": "name1", "namespace": "ns1"},
                    "apiVersion": "v1", "kind": "Kind1",
                },
            },
        },
        {
            name: "Fail to convert resource with missing name",
            resource:    komv1alpha1.AppliedResource{Name: "", Namespace: "ns1", Kind: "Kind1", APIVersion: "v1"},
            expectErr: true,
        },
        {
            name: "Fail to convert resource with missing namespace",
            resource:    komv1alpha1.AppliedResource{Name: "name1", Namespace: "", Kind: "Kind1", APIVersion: "v1"},
            expectErr: true,
        },
        {
            name: "Fail to convert resource with missing kind",
            resource:    komv1alpha1.AppliedResource{Name: "name1", Namespace: "ns1", Kind: "", APIVersion: "v1"},
            expectErr: true,
        },
        {
            name: "Fail to convert resource with missing API version",
            resource:    komv1alpha1.AppliedResource{Name: "name1", Namespace: "ns1", Kind: "Kind1", APIVersion: ""},
            expectErr: true,
        },
    }
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result, err := ToUnstructured(tc.resource) 
            if tc.expectErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tc.expected, result)
            }
        })
    }
}

func TestStatus_ToAppliedResource(t *testing.T) {
    testCases := []struct {
        name  string
        input        unstructured.Unstructured
        expected     *komv1alpha1.AppliedResource
        expectErr  bool
    }{
        {
            name: "Convert a fully populated unstructured object",
            input: unstructured.Unstructured{
                Object: map[string]interface{}{
                    "apiVersion": "v1",
                    "kind":       "Kind1",
                    "metadata": map[string]interface{}{
                        "name":      "name1",
                        "namespace": "ns1",
                    },
                },
            },
            expected: &komv1alpha1.AppliedResource{
                Name:       "name1",
                Namespace:  "ns1",
                Kind:       "Kind1",
                APIVersion: "v1",
            },
            expectErr: false,
        },
        {
            name: "Fail to convert due to missing name",
            input: unstructured.Unstructured{
                Object: map[string]interface{}{
                    "apiVersion": "v1",
                    "kind":       "Kind1",
                    "metadata":   map[string]interface{}{"namespace": "ns1"},
                },
            },
            expectErr: true,
        },
        {
            name: "Fail to convert due to missing namespace",
            input: unstructured.Unstructured{
                Object: map[string]interface{}{
                    "apiVersion": "v1",
                    "kind":       "Kind1",
                    "metadata":   map[string]interface{}{"name": "name1"},
                },
            },
            expectErr: true,
        },
        {
            name: "Fail to convert due to missing kind",
            input: unstructured.Unstructured{
                Object: map[string]interface{}{
                    "apiVersion": "v1",
                    "metadata": map[string]interface{}{
                        "name":      "name1",
                        "namespace": "ns1",
                    },
                },
            },
            expectErr: true,
        },
        {
            name: "Fail to convert due to missing apiVersion",
            input: unstructured.Unstructured{
                Object: map[string]interface{}{
                    "kind": "Kind1",
                    "metadata": map[string]interface{}{
                        "name":      "name1",
                        "namespace": "ns1",
                    },
                },
            },
            expectErr: true,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result, err := ToAppliedResource(tc.input)
            if tc.expectErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tc.expected, result)
            }
        })
    }
}
