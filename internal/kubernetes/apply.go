package kubernetes

import (
	"context"

	komtool "github.com/kkb0318/kom/internal/tool"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Handler struct {
	Client client.Client
	Owner  Owner
}

func (h Handler) ApplyAll(ctx context.Context, r komtool.ResourceManager) error {
	resources, err := r.Helm()
	if err != nil {
		return err
	}
	for _, resource := range resources {
		repo := resource.Repository()
		err = h.Apply(ctx, repo)
		if err != nil {
			return err
		}
		charts := resource.Charts()
		for _, chart := range charts {
			err = h.Apply(ctx, chart)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (h Handler) Apply(ctx context.Context, obj client.Object) error {
	opts := []client.PatchOption{
		client.ForceOwnership,
		client.FieldOwner(h.Owner.Field),
	}

	return h.Client.Patch(ctx, obj, client.Apply, opts...)
}
