package spec

import (
	"github.com/scylladb/scylla-operator/pkg/analyze/selector/internal/relation"
	"iter"
	"maps"
	"reflect"
)

// A Spec is a container for relation.Relation instances
// It enforces types of relation.Relation's arguments and a single relation per
// unordered pair of kinds
type Spec struct {
	types     map[string]reflect.Type
	relations map[string]map[string]relation.Relation
}

func New() *Spec {
	return &Spec{
		types:     make(map[string]reflect.Type),
		relations: make(map[string]map[string]relation.Relation),
	}
}

func (s *Spec) Add(name string, typ reflect.Type) bool {
	if _, contains := s.types[name]; contains {
		return false
	}

	s.types[name] = typ
	s.relations[name] = make(map[string]relation.Relation)

	return true
}

func (s *Spec) List() iter.Seq2[string, reflect.Type] {
	return maps.All(s.types)
}

func (s *Spec) Relate(relation relation.Relation) bool {
	if relation == nil {
		return false
	}

	firstName, firstType := relation.FirstParameter()
	secondName, secondType := relation.SecondParameter()

	if typ, exists := s.types[firstName]; !exists || firstType != typ {
		return false
	}

	if typ, exists := s.types[secondName]; !exists || secondType != typ {
		return false
	}

	if firstName > secondName {
		firstName, secondName = secondName, firstName
	}

	if _, exists := s.relations[firstName][secondName]; exists {
		return false
	}

	s.relations[firstName][secondName] = relation

	return true
}

func (s *Spec) Relation(first, second string) relation.Relation {
	if first > second {
		first, second = second, first
	}

	return s.relations[first][second]
}
