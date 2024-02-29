package factory

import (
	komv1alpha1 "github.com/kkb0318/kom/api/v1alpha1"
	komtool "github.com/kkb0318/kom/internal/tool"
	"github.com/kkb0318/kom/internal/tool/argo"
	"github.com/kkb0318/kom/internal/tool/flux"
)

func NewResourceManager(obj komv1alpha1.OperatorManager) komtool.ResourceManager {
	switch obj.Spec.Tool {
	case komv1alpha1.FluxCDTool:
		return flux.NewFlux(obj)
	case komv1alpha1.ArgoCDTool:
		return argo.NewArgo(obj)
	default:
		return flux.NewFlux(obj)
	}
}
