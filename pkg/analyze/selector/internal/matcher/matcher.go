package matcher

import (
	"cmp"
	"fmt"
	"github.com/scylladb/scylla-operator/pkg/analyze/selector/internal/relation"
	"github.com/scylladb/scylla-operator/pkg/analyze/selector/internal/spec"
	"maps"
	"slices"
)

// Match calls callback for every selection of values from resources that
// satisfy all constraints in spec.
// Given callback should return whether calls on further results should be
// performed.
func Match(
	spec *spec.Spec,
	resources map[string][]any,
	callback func(map[string]any) (bool, error),
) error {
	kinds := make([]string, 0)

	for kind := range spec.List() {
		kinds = append(kinds, kind)
		if _, contains := resources[kind]; !contains {
			return fmt.Errorf("Missing key %s", kind)
		}
	}

	slices.SortFunc(kinds, func(lhs, rhs string) int {
		return cmp.Compare(len(resources[lhs]), len(resources[rhs]))
	})

	_, err := match(spec, kinds, resources, make(map[string]any, len(kinds)), callback)

	return err
}

func match(
	spec *spec.Spec,
	orderedKinds []string,
	resources map[string][]any,
	prefix map[string]any,
	callback func(map[string]any) (bool, error),
) (bool, error) {
	if len(prefix) >= len(orderedKinds) {
		return callback(maps.Clone(prefix))
	}

	kind := orderedKinds[len(prefix)]
	for _, resource := range resources[kind] {
		prefix[kind] = resource

		appendable, err := doesMatchSpec(spec, resources, prefix, kind, resource)
		if err != nil {
			return false, err
		}

		if appendable {
			continu, err := match(spec, orderedKinds, resources, prefix, callback)

			if !continu || err != nil {
				return false, err
			}
		}

		delete(prefix, kind)
	}

	return true, nil
}

func doesMatchSpec(
	spec *spec.Spec,
	resources map[string][]any,
	selection map[string]any,
	newKind string,
	newResource any,
) (bool, error) {
	for otherKind, otherResource := range selection {
		relation := spec.Relation(otherKind, newKind)

		if otherResource != nil && newResource != nil {
			if relation == nil {
				continue
			}

			related, err := relation.Check(
				otherKind, otherResource,
				newKind, newResource,
			)
			if !related || err != nil {
				return false, err
			}
		} else if otherResource != nil && newResource == nil {
			result, err := checkRelationWithNil(
				otherKind, otherResource, newKind, resources[newKind], relation,
			)

			if !result || err != nil {
				return false, err
			}
		} else if otherResource == nil && newResource != nil {
			result, err := checkRelationWithNil(
				newKind, newResource, otherKind, resources[otherKind], relation,
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
	presentKind string,
	presentResource any,
	absentKind string,
	absentResources []any,
	relation relation.Relation,
) (bool, error) {
	if relation == nil {
		return true, nil
	}

	for _, absentResource := range absentResources {
		if absentResource == nil {
			continue
		}

		related, err := relation.Check(presentKind, presentResource, absentKind, absentResource)
		if related || err != nil {
			return false, err
		}
	}

	return true, nil
}
