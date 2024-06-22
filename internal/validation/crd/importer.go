package crd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ldassonville/helm-live/internal/evaluation/kubeconform"
	"os"
	"path/filepath"
	"strings"
)

type Importer struct {
	resolver Resolver
}

func NewImporter(resolver Resolver) *Importer {

	i := &Importer{
		resolver: resolver,
	}
	return i
}

func (i *Importer) Import(ctx context.Context, location string, group, kind string) error {

	crd, err := i.resolver.Resolve(ctx, group, kind)
	if err != nil {
		return err
	}

	if crd != nil {
		for _, version := range crd.Spec.Versions {
			groupVersion := strings.TrimPrefix(fmt.Sprintf("%s/%s", group, version.Name), "/")

			msgLoc, err := kubeconform.GetSchemaPath(location, kind, groupVersion, "master", true)
			if err != nil {
				return err
			}

			dir := filepath.Dir(msgLoc)
			if err = os.MkdirAll(dir, os.ModePerm); err != nil {
				return err
			}

			data, err := json.MarshalIndent(version.Schema.OpenAPIV3Schema, "", "  ")
			if err != nil {
				return err
			}

			err = i.writeFile(data, msgLoc)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (i *Importer) writeFile(data []byte, filepath string) error {

	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return err
	}
	return nil

}
