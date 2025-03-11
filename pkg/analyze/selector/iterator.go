package selector

import (
	"github.com/scylladb/scylla-operator/pkg/analyze/selector/internal/matcher"
	"github.com/scylladb/scylla-operator/pkg/analyze/selector/internal/predicate"
	"github.com/scylladb/scylla-operator/pkg/analyze/selector/internal/spec"
	"github.com/scylladb/scylla-operator/pkg/analyze/snapshot"
)

type Iterator struct {
	spec   *spec.Spec
	values map[string][]any
}

func (s *Selector) FromSnapshot(snapshot snapshot.Snapshot) *Iterator {
	result := make(map[string][]any)

	for name, typ := range s.spec.List() {
		values := snapshot.List(typ)

		if filter := s.filter[name]; filter != nil {
			var err error
			values, err = filterValues(filter, values)
			if err != nil {
				return nil
			}
		}

		if s.nilable[name] {
			values = append(values, nil)
		}

		result[name] = values
	}

	return &Iterator{spec: s.spec, values: result}
}

func filterValues(filter *predicate.Predicate, values []any) ([]any, error) {
	result := make([]any, 0, len(values))

	for _, value := range values {
		res, err := filter.Test(value)
		if err != nil {
			return nil, err
		}

		if res {
			result = append(result, value)
		}
	}

	return result, nil
}

// Constructs a Iterator which calls callback for every match
func (it *Iterator) ForEach(callback func(map[string]any) (bool, error)) error {
	return matcher.ForEach(it.spec, it.values, callback)
}

// Constructs a Iterator which returns a slice of matches
func (it *Iterator) Collect() ([]map[string]any, error) {
	result := make([]map[string]any, 0)
	err := matcher.ForEach(it.spec, it.values, func(values map[string]any) (bool, error) {
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
	err := matcher.ForEach(it.spec, it.values, func(values map[string]any) (bool, error) {
		result = append(result, values)
		return len(result) < n, nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
