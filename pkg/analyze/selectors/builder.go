package selectors

import (
	"github.com/scylladb/scylla-operator/pkg/analyze/sources"
	"math"
	"reflect"
)

type Builder struct {
	resources   map[string]reflect.Type
	constraints map[string][]*constraint
	assertion   map[string]*predicate
	relations   []*relation
}

func Type[T any]() reflect.Type {
	return reflect.TypeFor[T]()
}

func Select(label string, typ reflect.Type) *Builder {
	return (&Builder{
		resources:   make(map[string]reflect.Type),
		constraints: make(map[string][]*constraint),
		assertion:   make(map[string]*predicate),
		relations:   make([]*relation, 0),
	}).Select(label, typ)
}

func (b *Builder) Select(label string, typ reflect.Type) *Builder {
	if _, exists := b.resources[label]; exists {
		panic("TODO: Handle duplicate labels")
	}

	b.resources[label] = typ

	return b
}

func (b *Builder) Filter(label string, f any) *Builder {
	typ, defined := b.resources[label]
	if !defined {
		panic("TODO: Handle undefined labels in Filter")
	}

	constraint := newConstraint(label, f)
	if constraint.Labels()[label] != typ {
		panic("TODO: Handle mismatched type in Filter")
	}

	b.constraints[label] = append(b.constraints[label], constraint)

	return b
}

func (b *Builder) Assert(label string, f any) *Builder {
	typ, defined := b.resources[label]
	if !defined {
		panic("TODO: Handle undefined labels in Filter")
	}

	assertion := newPredicate(label, f)
	if assertion.Labels()[label] != typ {
		panic("TODO: Handle mismatched type in Assert")
	}

	b.assertion[label] = assertion

	return b
}

func (b *Builder) Relate(lhs, rhs string, f any) *Builder {
	// TODO: Check input

	relation := newRelation(lhs, rhs, f)

	b.relations = append(b.relations, relation)

	return b
}

func (b *Builder) CollectAll() func(*sources.DataSource2) []map[string]any {
	return b.Collect(math.MaxInt)
}

func (b *Builder) Collect(limit int) func(*sources.DataSource2) []map[string]any {
	executor := newExecutor(
		b.resources,
		b.constraints,
		b.assertion,
		b.relations,
	)

	return func(ds *sources.DataSource2) []map[string]any {
		result := make([]map[string]any, 0)
		count := 0

		executor.execute(ds, func(resources map[string]any) bool {
			if count < limit {
				result = append(result, resources)
				count += 1
				return true
			}
			return false
		})

		return result
	}
}

func (b *Builder) ForEach(labels []string, function any) func(*sources.DataSource2) {
	for _, label := range labels {
		if _, contains := b.resources[label]; !contains {
			panic("TODO: Handle undefined label")
		}
	}

	callback := newFunction[bool](labels, function)
	executor := newExecutor(
		b.resources,
		b.constraints,
		b.assertion,
		b.relations,
	)

	return func(ds *sources.DataSource2) {
		executor.execute(ds, func(resources map[string]any) bool {
			labels := callback.Labels()
			args := make(map[string]any, len(labels))

			for label, resource := range resources {
				if _, exists := labels[label]; !exists {
					continue
				}

				args[label] = resource
			}

			return callback.Call(args)
		})
	}
}

func (b *Builder) Any() func(*sources.DataSource2) bool {
	executor := newExecutor(
		b.resources,
		b.constraints,
		b.assertion,
		b.relations,
	)

	return func(ds *sources.DataSource2) bool {
		result := false

		executor.execute(ds, func(_ map[string]any) bool {
			result = true
			return false
		})

		return result
	}
}
