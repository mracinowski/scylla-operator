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

func NewPredicate(t *testing.T, label string, f any) *predicate.Predicate {
	predicate, err := predicate.New(label, f)
	if err != nil {
		t.Fatal(err)
	}

	return predicate
}

func NewRelation(t *testing.T, lhs, rhs string, f any) relation.Relation {
	relation, err := relation.New(lhs, rhs, f)
	if err != nil {
		t.Fatal(err)
	}

	return relation
}

func MakeRelations(
	t *testing.T,
	types map[string]reflect.Type,
	relations []relation.Relation,
) spec.Spec {
	result := *spec.New()

	for name, ty := range types {
		if !result.Add(name, ty) {
			t.Fatal("Invalid field")
		}
	}

	for _, relation := range relations {
		if !result.Relate(relation) {
			t.Fatal("Invalid relation")
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

func TestForEach(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		relations spec.Spec
		values    map[string][]any
		expected  []map[string]any
	}{
		{
			name: "no relations",
			relations: MakeRelations(t, map[string]reflect.Type{
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
			relations: MakeRelations(t, map[string]reflect.Type{
				"A": reflect.TypeFor[int](),
				"B": reflect.TypeFor[bool](),
			}, []relation.Relation{
				NewRelation(t, "A", "B",
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
			relations: MakeRelations(t, map[string]reflect.Type{
				"A": reflect.TypeFor[int](),
				"B": reflect.TypeFor[bool](),
			}, []relation.Relation{
				NewPredicate(t, "A",
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

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result := make([]map[string]any, 0, len(tc.expected))
			err := ForEach(&tc.relations, tc.values, func(values map[string]any) (bool, error) {
				result = append(result, values)
				return true, nil
			})

			if err != nil {
				t.Fatalf("Unexpected error: %s", err)
			}

			slices.SortFunc(tc.expected, CompareMaps)
			slices.SortFunc(result, CompareMaps)

			if !slices.EqualFunc(tc.expected, result, maps.Equal) {
				for i, match := range result {
					t.Logf("%d: %+v", i, match)
				}

				t.FailNow()
			}
		})
	}
}
