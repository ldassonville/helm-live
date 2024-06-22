package validation

import (
	"context"
	v12 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	crdv1 "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Resolver struct {
	crdItf crdv1.CustomResourceDefinitionInterface
}

func (r *Resolver) Resolve(ctx context.Context, group, version, kind string) (*v12.CustomResourceDefinition, error) {

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
