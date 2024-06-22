package schema

import (
	"errors"
	"io"
	"net/http"
)

type HttpResolver struct {
	GetSchemaPath func(group string, kind string, version string) (string, error)
}

func (r *HttpResolver) ResolveSchema(group string, kind string, version string) ([]byte, error) {

	path, err := r.GetSchemaPath(group, kind, version)
	if err != nil {
		return nil, err
	}

	// Try to download the schema
	resp, err := http.Get(path)
	if err != nil {
		return nil, err
	}

	// Check if the response is OK
	if resp.StatusCode == 200 {

		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return b, nil
	}

	if resp.StatusCode == 404 {
		return nil, errors.New("schema not found")
	}

	return nil, errors.New("failed to download schema")
}
