package kubernetes

import (
	"context"

	komv1alpha1 "github.com/kkb0318/kom/api/v1alpha1"
	komstatus "github.com/kkb0318/kom/internal/status"
	komtool "github.com/kkb0318/kom/internal/tool"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

type Handler struct {
	client client.Client
	owner  Owner
}

// NewHelper returns an initialized Helper.
func NewHandler(obj client.Object, c client.Client, owner Owner) (*Handler, error) {
	return &Handler{
		client: c,
		owner:  owner,
	}, nil
}

func (h Handler) ApplyAll(ctx context.Context, r komtool.ResourceManager) ([]komv1alpha1.AppliedResource, error) {
	var appliedResources []komv1alpha1.AppliedResource
	resources, err := r.Helm()
	if err != nil {
		return nil, err
	}
	for _, resource := range resources {
		repo := resource.Repository()
		applied, err := h.Apply(ctx, repo)
		if err != nil {
			return nil, err
		}
		appliedResources = append(appliedResources, *applied)
		charts := resource.Charts()
		for _, chart := range charts {
			applied, err := h.Apply(ctx, chart)
			if err != nil {
				return nil, err
			}
			appliedResources = append(appliedResources, *applied)
		}
	}
	return appliedResources, nil
}

func (h Handler) Apply(ctx context.Context, obj client.Object) (*komv1alpha1.AppliedResource, error) {
	opts := []client.PatchOption{
		client.ForceOwnership,
		client.FieldOwner(h.owner.Field),
	}
	gvk, err := apiutil.GVKForObject(obj, h.client.Scheme())
	if err != nil {
		return nil, err
	}

	u := &unstructured.Unstructured{}
	unstructured, err := runtime.DefaultUnstructuredConverter.ToUnstructured(obj)
	if err != nil {
		return nil, err
	}
	u.Object = unstructured
	u.SetGroupVersionKind(gvk)
	u.SetManagedFields(nil)
	err = h.client.Patch(ctx, u, client.Apply, opts...)
	if err != nil {
		return nil, err
	}
	return komstatus.ToAppliedResource(*u)
}
