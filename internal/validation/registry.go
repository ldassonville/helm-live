package validation

import (
	"context"
	"github.com/ldassonville/helm-live/internal/validation/crd"
	v1 "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
)

func NewRegistry(schemaPath string, crdi v1.CustomResourceDefinitionInterface) *Registry {

	resolver := crd.NewResolver(crdi)

	return &Registry{
		SchemaPath: schemaPath,
		importer:   crd.NewImporter(resolver),
	}
}

type Registry struct {
	SchemaPath string
	importer   *crd.Importer
}

func (r *Registry) ImportSchema(ctx context.Context, group, version, kind string) error {

	err := r.importer.Import(ctx, r.SchemaPath, group, version, kind)
	if err != nil {
		return err
	}

	return nil

}
