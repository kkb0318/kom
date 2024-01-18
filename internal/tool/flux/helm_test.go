package flux

import (
	"testing"

	komv1alpha1 "github.com/kkb0318/kom/api/v1alpha1"
	komtool "github.com/kkb0318/kom/internal/tool"
	"gotest.tools/assert"
)

func TestFluxHelm_New(t *testing.T) {
	tests := []struct {
		name        string
		inputs      []komv1alpha1.Helm
		expected    []komtool.Resource
		expectedErr bool
	}{
		{
			name:   "continue if not in previous",
			inputs: []komv1alpha1.Helm{},
      expected: []komtool.Resource{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := NewFluxHelmList(tt.inputs)
			if tt.expectedErr {
				assert.Error(t, err, "")
			} else {
				assert.NilError(t, err)
				assert.DeepEqual(t, tt.expected, actual)
			}
		})
	}
}
