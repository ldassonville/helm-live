package validation

import (
	"context"
	"github.com/ldassonville/helm-live/internal/validation/crd"
	"github.com/rs/zerolog/log"
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

func (r *Registry) ImportSchema(ctx context.Context, group, kind string) error {

	err := r.importer.Import(ctx, r.SchemaPath, group, kind)
	if err != nil {
		log.Err(err).Msgf("Failed to import schema for %s/%s", group, kind)
		return err
	}
	return nil

}
