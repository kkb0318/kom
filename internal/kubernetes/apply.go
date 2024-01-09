package apply

import (
	komtool "github.com/kkb0318/kom/internal/tool"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ApplyAll(r komtool.ResourceManager) error {
	resources, err := r.Helm()
	if err != nil {
		return err
	}
	for _, resource := range resources {
		repo := resource.Repository()
		err = Apply(repo)
		if err != nil {
			return err
		}
		charts := resource.Charts()
		for _, chart := range charts {
			err = Apply(chart)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func Apply(obj client.Object) error {
	return nil
}
