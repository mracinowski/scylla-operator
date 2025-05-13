package selector

import (
	"github.com/scylladb/scylla-operator/pkg/analyze/selector/internal/matcher"
	"github.com/scylladb/scylla-operator/pkg/analyze/selector/internal/spec"
)

type Iterator struct {
	spec   *spec.Spec
	values map[string][]any
}

// Constructs a Iterator which calls callback for every match
func (it *Iterator) ForEach(callback func(map[string]any) (bool, error)) error {
	return matcher.Match(it.spec, it.values, callback)
}

// Constructs a Iterator which returns a slice of matches
func (it *Iterator) Collect() ([]map[string]any, error) {
	result := make([]map[string]any, 0)
	err := matcher.Match(it.spec, it.values, func(values map[string]any) (bool, error) {
		result = append(result, values)
		return true, nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// Constructs a Iterator which returns a slice of at most n matches
func (it *Iterator) Take(n int) ([]map[string]any, error) {
	result := make([]map[string]any, 0)
	err := matcher.Match(it.spec, it.values, func(values map[string]any) (bool, error) {
		result = append(result, values)
		return len(result) < n, nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
