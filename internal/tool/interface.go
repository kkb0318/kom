package internal

import (
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ResourceManager interface {
	Helm() ([]Resource, error)
	Git() ([]Resource, error)
}

type Resource interface {
	Repository() client.Object
	Charts() []client.Object
}
