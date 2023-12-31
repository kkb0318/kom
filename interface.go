package internal

import (
	"fmt"

	komv1aplha1 "github.com/kkb0318/kom/api/v1alpha1"
)

type Repository interface {
	Helm() error
	Oci() error
	Git() error
}

func Apply(r Repository, t komv1aplha1.RepositoryType) error {
	switch t {
	case komv1aplha1.HelmRepository:
		return r.Helm()
	case komv1aplha1.OciRepository:
		return r.Oci()
	case komv1aplha1.GitRepository:
		return r.Git()
	}
	return fmt.Errorf("error")
}
