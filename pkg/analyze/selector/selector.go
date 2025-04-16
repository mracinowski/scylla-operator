package selector

import (
	"fmt"
	"github.com/scylladb/scylla-operator/pkg/analyze/selector/internal/predicate"
	"github.com/scylladb/scylla-operator/pkg/analyze/selector/internal/relation"
	"github.com/scylladb/scylla-operator/pkg/analyze/selector/internal/spec"
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

func Select(name string, typ reflect.Type, filter any) *Selector {
	return New().Select(name, typ, filter)
}

func SelectWithNil(name string, typ reflect.Type, filter any) *Selector {
	return New().SelectWithNil(name, typ, filter)
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
