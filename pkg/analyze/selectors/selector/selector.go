// Package selector implements... TODO: https://tip.golang.org/doc/comment
package selector

import (
	"cmp"
	"github.com/scylladb/go-set/strset"
	"github.com/scylladb/scylla-operator/pkg/analyze/selectors/function"
	"maps"
	"reflect"
	"slices"
)

type set[T comparable] map[T]struct{}

type Selector struct {
	labels     strset.Set
	predicates map[string]*function.Predicate
	relations  map[string]map[string]*function.Relation
}

// TODO: Document
func New(
	labels map[string]reflect.Type,
	predicates []*function.Predicate,
	relations []*function.Relation,
) (*Selector, error) {
	mappedPredicates := make(map[string]*function.Predicate, len(predicates))
	for _, predicate := range predicates {
		parameter, parameterType := predicate.Parameter()
		expectedType, defined := labels[parameter]
		if !defined {
			panic("TODO")
		}

		if expectedType != parameterType {
			panic("TODO")
		}

		mappedPredicates[parameter] = predicate
	}

	mappedRelations := make(map[string]map[string]*function.Relation)
	for _, relation := range relations {
		firstParameter, firstType := relation.FirstParameter()
		expectedType, defined := labels[firstParameter]
		if !defined {
			panic("TODO")
		}

		if expectedType != firstType {
			panic("TODO")
		}

		secondParameter, secondType := relation.SecondParameter()
		expectedType, defined = labels[secondParameter]
		if !defined {
			panic("TODO")
		}

		if expectedType != secondType {
			panic("TODO")
		}

		if firstParameter > secondParameter {
			firstParameter, secondParameter = secondParameter, firstParameter
		}

		mappedRelations[firstParameter][secondParameter] = relation
	}

	return &Selector{
		labels:     *strset.New(slices.Collect(maps.Keys(labels))...),
		predicates: mappedPredicates,
		relations:  mappedRelations,
	}, nil
}

// TODO: Document
func (s *Selector) Select(
	values map[string][]any,
	callback func(map[string]any) (bool, error),
) error {
	panic("TODO")
	labels := order(values)

	selection := make(map[string]any, len(labels))
	_, err := s.match(labels, values, selection, callback)
	return err
}

func order(values map[string][]any) []string {
	result := make([]string, len(values))

	for label := range values {
		result = append(result, label)
	}

	slices.SortFunc(result, func(lhs, rhs string) int {
		return cmp.Compare(len(values[lhs]), len(values[rhs]))
	})

	return result
}

func (s *Selector) match(
	labels []string,
	values map[string][]any,
	selection map[string]any,
	callback func(map[string]any) (bool, error),
) (bool, error) {
	if len(selection) >= len(labels) {
		return callback(maps.Clone(selection))
	}

	label := labels[len(selection)]
	for _, value := range values {
		if !s.canAppend(values, selection, label, value) {
			continue
		}

		selection[label] = value
		continu, err := s.match(labels, values, selection, callback)
		delete(selection, label)

		if continu == false || err != nil {
			return false, err
		}
	}

	return true, nil
}

func (s *Selector) canAppend(
	resources map[string][]any,
	selection map[string]any,
	newLabel string,
	newValue any,
) bool {
	for otherLabel, otherValue := range selection {
		relation := s.relation(otherLabel, newLabel)

		if otherValue != nil && newValue != nil {
			if relation != nil && relation.Check(
				otherLabel, otherValue, newLabel, newValue,
			) {
				return false
			}

		} else if otherValue != nil && newValue == nil {
			if !checkRelationWithNil(
				otherLabel, otherValue, newLabel, resources[newLabel], relation,
			) {
				return false
			}

		} else if otherValue == nil && newValue != nil {
			if !checkRelationWithNil(
				newLabel, newValue, otherLabel, resources[otherLabel], relation,
			) {
				return false
			}

		} else if relation != nil {
			return false
		}
	}

	return true
}

func areKeysEqual[K comparable, V any](lhs, rhs map[K]V) bool {
	if len(lhs) != len(rhs) {
		return false
	}

	for key := range lhs {
		if _, contains := rhs[key]; !contains {
			return false
		}
	}

	return true
}

func (s *Selector) relation(lhs, rhs string) *function.Relation {
	if lhs > rhs {
		lhs, rhs = rhs, lhs
	}

	return s.relations[lhs][rhs]
}

func checkRelationWithNil(
	presentLabel string,
	presentValue any,
	absentLabel string,
	absentValues []any,
	relation *function.Relation,
) bool {
	if relation == nil {
		return true
	}

	for _, absentValue := range absentValues {
		if absentValue == nil {
			continue
		}

		if relation.Check(presentLabel, presentValue, absentLabel, absentValue) {
			return false
		}
	}

	return true
}
