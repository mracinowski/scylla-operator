package snapshot

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"reflect"

	"k8s.io/apimachinery/pkg/runtime"
)

func NewDataSourceFromFS(fsys fs.FS, decoder runtime.Decoder) (Snapshot, error) {
	var ds Snapshot = &snapshot{
		objects: make(map[reflect.Type][]interface{}),
	}

	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || filepath.Ext(path) != ".yaml" {
			return nil
		}
		content, err := fs.ReadFile(fsys, path)
		if err != nil {
			return fmt.Errorf("can't read file %q: %w", path, err)
		}
		obj, _, err := decoder.Decode(content, nil, nil)

		if err != nil {
			if !runtime.IsNotRegisteredError(err) {
				return fmt.Errorf("can't deserialize file %q: %w", path, err)
			}
			return nil
		}
		ds.Add(obj)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("can't walk the file tree: %w", err)
	}
	return ds, nil
}
