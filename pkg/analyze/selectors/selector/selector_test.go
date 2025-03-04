package selector

import (
	"github.com/scylladb/scylla-operator/pkg/analyze/selectors/function"
	"reflect"
	"testing"
)

type ObjectA struct {
	Name string
}

func (a ObjectA) String() string {
	return a.Name
}

type SelectTest struct {
	labels     map[string]reflect.Type
	predicates []*function.Predicate
	relations  []*function.Relation
	values     map[string][]any
	expected   []map[string]any
}

var SelectTests = []SelectTest{
	{
		labels: map[string]reflect.Type{
			"alfa":  reflect.TypeFor[int](),
			"bravo": reflect.TypeFor[string](),
		},
		predicates: []*function.Predicate{
			newPredicate("alfa", func() (bool, error) {
				return false, nil
			}),
		},
		relations: []*function.Relation{
			newRelation("alfa", "bravo", func() (bool, error) {
				return false, nil
			}),
		},
		values: map[string][]any{
			"alpha": []any{},
		},
		expected: []map[string]any{},
	},
}

func Assert(t *testing.T, expected map[string]any, selector Selector) {
	err := selector.Select(
		nil,
		func(values map[string]any) (bool, error) {
			t.Log(values)
			return true, nil
		},
	)

	if err != nil {
		t.Error(err)
	}
}

type SelectTest struct {
	name       string
	labels     map[string]reflect.Type
	values     map[string][]any
	predicates []*function.Predicate
	relations  []*function.Relation
	expected   []map[string]any
}

func TestSelect(t *testing.T) {
	tests := []SelectTest{
		{
			name:   "",
			labels: map[string]reflect.Type{},
		},
	}
	for _, test := range SelectTests {

		selector, _ := New(test.labels, test.predicates, test.relations)
		err := selector.Select(test.values, func(values map[string]any) {

		})
		if !reflect.DeepEqual(test.expected, result) {
			t.Errorf("Fail")
		}
	}
}
