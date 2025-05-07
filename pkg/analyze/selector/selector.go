package selector

import (
	"fmt"
	"github.com/scylladb/scylla-operator/pkg/analyze/selector/internal/predicate"
	"github.com/scylladb/scylla-operator/pkg/analyze/selector/internal/relation"
	"github.com/scylladb/scylla-operator/pkg/analyze/selector/internal/spec"
	"github.com/scylladb/scylla-operator/pkg/analyze/snapshot"
	"github.com/scylladb/scylla-operator/pkg/helpers/slices"
	"reflect"
)

type Selector struct {
	spec    *spec.Spec
	filter  map[string]*predicate.Predicate
	nilable map[string]bool
	error
}

func Type[T any]() reflect.Type {
	return reflect.TypeFor[T]()
}

func New() *Selector {
	return &Selector{
		spec:    spec.New(),
		filter:  make(map[string]*predicate.Predicate),
		nilable: make(map[string]bool),
		error:   nil,
	}
}

func (s *Selector) Select(name string, typ reflect.Type, filter any) *Selector {
	if s.error != nil {
		return s
	}

	if !s.spec.Add(name, typ) {
		s.error = fmt.Errorf("Duplicate %s definition", name)
		return s
	}

	if filter != nil {
		p, err := predicate.New(name, filter)
		if err != nil {
			s.error = err
			return s
		}

		s.filter[name] = p
	}

	s.nilable[name] = false

	return s
}

func (s *Selector) SelectWithNil(name string, typ reflect.Type, filter any) *Selector {
	if s.error != nil {
		return s
	}

	s.Select(name, typ, filter)

	s.nilable[name] = true

	return s
}

func (s *Selector) Relate(first, second string, lambda any) *Selector {
	if s.error != nil {
		return s
	}

	relation, err := relation.New(first, second, lambda)
	if err != nil {
		s.error = err
		return s
	}

	if !s.spec.Relate(relation) {
		s.error = fmt.Errorf("Invalid relation between %s and %s", first, second)
		return s
	}

	return s
}

func (s *Selector) Where(name string, lambda any) *Selector {
	if s.error != nil {
		return s
	}

	predicate, err := predicate.New(name, lambda)
	if err != nil {
		s.error = err
		return s
	}

	if !s.spec.Relate(predicate) {
		s.error = fmt.Errorf("Invalid Where condition for %s", name)
		return s
	}

	return s
}

func filterWithError(array []any, filter *predicate.Predicate) ([]any, error) {
	var err error

	res := slices.Filter(array, func(value any) bool {
		r, e := filter.Test(value)
		if e != nil && err == nil {
			err = e
		}

		return r
	})

	return res, err
}

func (s *Selector) IteratorFromSnapshot(snapshot snapshot.Snapshot) (*Iterator, error) {
	if s.error != nil {
		return nil, s.error
	}

	result := make(map[string][]any)

	for name, typ := range s.spec.List() {
		values := snapshot.List(typ)

		if filter := s.filter[name]; filter != nil {
			var err error
			values, err = filterWithError(values, filter)
			if err != nil {
				return nil, err
			}
		}

		if s.nilable[name] {
			values = append(values, nil)
		}

		result[name] = values
	}

	return &Iterator{spec: s.spec, values: result}, nil
}
