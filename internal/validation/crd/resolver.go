package crd

import (
	"context"
	v12 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	crdv1 "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Resolver interface {
	Resolve(ctx context.Context, group, kind string) (*v12.CustomResourceDefinition, error)
}

func NewResolver(crdItf crdv1.CustomResourceDefinitionInterface) Resolver {
	return &resolver{
		crdItf: crdItf,
	}
}

type resolver struct {
	crdItf crdv1.CustomResourceDefinitionInterface
}

// Resolve returns the CustomResourceDefinition for the given group and kind
// return nil if not found
func (r *resolver) Resolve(ctx context.Context, group, kind string) (*v12.CustomResourceDefinition, error) {

	listOptions := v1.ListOptions{}
	list, err := r.crdItf.List(ctx, listOptions)
	if err != nil {
		return nil, err
	}

	for _, crd := range list.Items {
		if crd.Spec.Group == group && crd.Spec.Names.Kind == kind {
			return &crd, nil
		}
	}
	return nil, nil
}
