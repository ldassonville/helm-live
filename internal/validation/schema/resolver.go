package schema

import (
	"fmt"
	"github.com/ldassonville/helm-live/internal/evaluation/kubeconform"
	"strings"
)

type Resolver interface {
	ResolveSchema(group string, kind string, APIVersion string) ([]byte, error)
}

type schemaResolver struct {
	Locations []string
	resolvers []Resolver
}

func New(locations []string) Resolver {

	res := &schemaResolver{
		Locations: locations,
	}

	// Init resolvers
	for _, location := range locations {

		var getSchemaPathFnc = func(group string, kind string, version string) (string, error) {
			groupVersion := strings.TrimPrefix(fmt.Sprintf("%s/%s", group, version), "/")
			return kubeconform.GetSchemaPath(location, kind, groupVersion, "master", true)
		}

		if strings.HasPrefix(strings.ToLower(location), "http") {
			res.resolvers = append(res.resolvers, &HttpResolver{
				GetSchemaPath: getSchemaPathFnc,
			})
		} else {
			res.resolvers = append(res.resolvers, &FileResolver{
				GetSchemaPath: getSchemaPathFnc,
			})
		}
	}

	return res
}

func (r *schemaResolver) ResolveSchema(group string, kind string, APIVersion string) ([]byte, error) {

	for _, resolver := range r.resolvers {
		data, err := resolver.ResolveSchema(group, kind, APIVersion)
		if err == nil {
			return data, nil
		}
	}
	return nil, fmt.Errorf("failed to download schema for %s/%s/%s", group, kind, APIVersion)
}
