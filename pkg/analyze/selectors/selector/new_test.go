package selector

import (
	"github.com/scylladb/scylla-operator/pkg/analyze/selectors/function"
	"reflect"
	"testing"
)

func newPredicate(label string, f any) *function.Predicate {
	predicate, err := function.NewPredicate(label, f)
	if err != nil {
		panic("TODO")
	}

	return predicate
}

func newRelation(lhs, rhs string, f any) *function.Relation {
	relation, err := function.NewRelation(lhs, rhs, f)
	if err != nil {
		panic("TODO")
	}

	return relation
}

type NewTest struct {
	name       string
	labels     map[string]reflect.Type
	predicates []*function.Predicate
	relations  []*function.Relation
	expected   error
}

func TestNew(t *testing.T) {
	tests := []NewTest{
		NewTest{
			name: "Basic",
			labels: map[string]reflect.Type{
				"alfa": reflect.TypeFor[int](),
			},
			predicates: []*function.Predicate{
				newPredicate("bravo", func() (bool, error) {
					return true, nil
				}),
			},
			relations: []*function.Relation{},
			expected:  nil,
		},
	}

	for _, test := range tests {
		result, err := New(test.labels, test.predicates, test.relations)
		if test.expected == nil && err != nil {
			t.Errorf("Fail: err != nil")
		}

		if test.expected == nil && result == nil {
			t.Errorf("Fail: result == nil")
		}

		if test.expected != nil && err == nil {
			t.Errorf("Fail: err == nil")
		}

		if test.expected != nil && result != nil {
			t.Errorf("Fail: result != nil")
		}
	}
}
