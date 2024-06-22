package schema

import (
	"errors"
	"io"
	"os"
)

type FileResolver struct {
	GetSchemaPath func(group string, kind string, version string) (string, error)
}

func (r *FileResolver) ResolveSchema(group string, kind string, version string) ([]byte, error) {

	path, err := r.GetSchemaPath(group, kind, version)
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	return io.ReadAll(file)

}
