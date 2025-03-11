package matcher

import (
	"cmp"
	"fmt"
	"github.com/scylladb/scylla-operator/pkg/analyze/selector/internal/predicate"
	"github.com/scylladb/scylla-operator/pkg/analyze/selector/internal/relation"
	"github.com/scylladb/scylla-operator/pkg/analyze/selector/internal/spec"
	"maps"
	"reflect"
	"slices"
	"testing"
)

func NewPredicate(label string, f any) *predicate.Predicate {
	predicate, err := predicate.New(label, f)
	if err != nil {
		panic(err)
	}

	return predicate
}

func NewRelation(lhs, rhs string, f any) relation.Relation {
	relation, err := relation.New(lhs, rhs, f)
	if err != nil {
		panic(err)
	}

	return relation
}

func MakeRelations(
	types map[string]reflect.Type,
	relations []relation.Relation,
) spec.Spec {
	result := *spec.New()

	for name, typ := range types {
		if !result.Add(name, typ) {
			panic("Invalid field")
		}
	}

	for _, relation := range relations {
		if !result.Relate(relation) {
			panic("Invalid relation")
		}
	}

	return result
}

func CompareMaps(x, y map[string]any) int {
	result := cmp.Compare(len(x), len(y))
	if result != 0 {
		return result
	}

	keys := slices.Sorted(maps.Keys(x))
	otherKeys := slices.Sorted(maps.Keys(y))
	result = slices.Compare(keys, otherKeys)

	if result != 0 {
		return result
	}

	for _, key := range keys {
		result = cmp.Compare(
			fmt.Sprintf("%+v", x[key]),
			fmt.Sprintf("%+v", y[key]),
		)
		if result != 0 {
			return result
		}
	}

	return 0
}

type ForEachTest struct {
	name      string
	relations spec.Spec
	values    map[string][]any
	expected  []map[string]any
}

func TestForEach(t *testing.T) {
	tests := []ForEachTest{
		{
			name: "no relations",
			relations: MakeRelations(map[string]reflect.Type{
				"A": reflect.TypeFor[int](),
				"B": reflect.TypeFor[bool](),
			}, []relation.Relation{}),
			values: map[string][]any{
				"A": {1, 2, 3},
				"B": {false, true},
			},
			expected: []map[string]any{
				{"A": 1, "B": false},
				{"A": 1, "B": true},
				{"A": 2, "B": false},
				{"A": 2, "B": true},
				{"A": 3, "B": false},
				{"A": 3, "B": true},
			},
		},
		{
			name: "single relation",
			relations: MakeRelations(map[string]reflect.Type{
				"A": reflect.TypeFor[int](),
				"B": reflect.TypeFor[bool](),
			}, []relation.Relation{
				NewRelation("A", "B",
					func(a int, b bool) (bool, error) {
						return (a%2 == 0) == b, nil
					}),
			}),
			values: map[string][]any{
				"A": {1, 2, 3},
				"B": {false, true},
			},
			expected: []map[string]any{
				{"A": 1, "B": false},
				{"A": 2, "B": true},
				{"A": 3, "B": false},
			},
		},
		{
			name: "single predicate",
			relations: MakeRelations(map[string]reflect.Type{
				"A": reflect.TypeFor[int](),
				"B": reflect.TypeFor[bool](),
			}, []relation.Relation{
				NewPredicate("A",
					func(a int) (bool, error) {
						return a%2 != 0, nil
					}),
			}),
			values: map[string][]any{
				"A": {1, 2, 3},
				"B": {false, true},
			},
			expected: []map[string]any{
				{"A": 1, "B": false},
				{"A": 1, "B": true},
				{"A": 3, "B": false},
				{"A": 3, "B": true},
			},
		},
	}

	for _, test := range tests {
		result := make([]map[string]any, 0, len(test.expected))
		err := ForEach(&test.relations, test.values, func(values map[string]any) (bool, error) {
			result = append(result, values)
			return true, nil
		})

		if err != nil {
			t.Errorf("%s: Unexpected error: %s", test.name, err)
		}

		slices.SortFunc(test.expected, CompareMaps)
		slices.SortFunc(result, CompareMaps)

		if !slices.EqualFunc(test.expected, result, maps.Equal) {
			t.Errorf("%s: Fail", test.name)

			for i, match := range result {
				t.Logf("%d: %+v", i, match)
			}
		}
	}
}
