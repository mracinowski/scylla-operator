package matcher

import (
	"cmp"
	"fmt"
	"github.com/scylladb/scylla-operator/pkg/analyze/selector/internal/relation"
	"github.com/scylladb/scylla-operator/pkg/analyze/selector/internal/spec"
	"maps"
	"slices"
)

type matcher struct {
	labels []string
	spec   *spec.Spec
	values map[string][]any
}

func ForEach(
	spec *spec.Spec,
	values map[string][]any,
	callback func(map[string]any) (bool, error),
) error {
	labels := make([]string, 0)

	for label := range spec.List() {
		labels = append(labels, label)
		if _, contains := values[label]; !contains {
			return fmt.Errorf("Missing key %s", label)
		}
	}

	slices.SortFunc(labels, func(lhs, rhs string) int {
		return cmp.Compare(len(values[lhs]), len(values[rhs]))
	})

	_, err := (&matcher{
		labels: labels,
		spec:   spec,
		values: values,
	}).forEach(make(map[string]any, len(labels)), callback)

	return err
}

func (it *matcher) forEach(
	prefix map[string]any,
	callback func(map[string]any) (bool, error),
) (bool, error) {
	if len(prefix) >= len(it.labels) {
		return callback(maps.Clone(prefix))
	}

	label := it.labels[len(prefix)]
	for _, value := range it.values[label] {
		prefix[label] = value

		canAppend, err := it.canAppend(prefix, label, value)
		if err != nil {
			return false, err
		}

		if canAppend {
			continu, err := it.forEach(prefix, callback)

			if !continu || err != nil {
				return false, err
			}
		}

		delete(prefix, label)
	}

	return true, nil
}

func (it *matcher) canAppend(
	selection map[string]any,
	newLabel string,
	newValue any,
) (bool, error) {
	for otherLabel, otherValue := range selection {
		relation := it.spec.Relation(otherLabel, newLabel)

		if otherValue != nil && newValue != nil {
			if relation == nil {
				continue
			}

			related, err := relation.Check(
				otherLabel, otherValue,
				newLabel, newValue,
			)
			if !related || err != nil {
				return false, err
			}
		} else if otherValue != nil && newValue == nil {
			result, err := checkRelationWithNil(
				otherLabel, otherValue, newLabel, it.values[newLabel], relation,
			)

			if !result || err != nil {
				return false, err
			}
		} else if otherValue == nil && newValue != nil {
			result, err := checkRelationWithNil(
				newLabel, newValue, otherLabel, it.values[otherLabel], relation,
			)

			if !result || err != nil {
				return false, nil
			}
		} else if relation != nil {
			return false, nil
		}
	}

	return true, nil
}

func checkRelationWithNil(
	presentLabel string,
	presentValue any,
	absentLabel string,
	absentValues []any,
	relation relation.Relation,
) (bool, error) {
	if relation == nil {
		return true, nil
	}

	for _, absentValue := range absentValues {
		if absentValue == nil {
			continue
		}

		related, err := relation.Check(presentLabel, presentValue, absentLabel, absentValue)
		if related || err != nil {
			return false, err
		}
	}

	return true, nil
}
