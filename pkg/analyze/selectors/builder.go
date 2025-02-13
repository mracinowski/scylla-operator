package selectors

import (
	"github.com/scylladb/scylla-operator/pkg/analyze/sources"
	scyllav1 "github.com/scylladb/scylla-operator/pkg/api/scylla/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
	"reflect"
)

func Type[T any]() reflect.Type {
	return reflect.TypeFor[T]()
}

type builder struct {
	fields      map[string]reflect.Type
	relations   []Relation
	constraints []Constraint
}

func Builder() *builder {
	return &builder{
		fields:      make(map[string]reflect.Type),
		relations:   make([]Relation, 0),
		constraints: make([]Constraint, 0),
	}
}

func (b *builder) New(name string, resource reflect.Type) *builder {
	// TODO: Define field of name `name` of resource type `resource`

	// TODO: Check if exists
	b.fields[name] = resource

	return b
}

func (b *builder) Join(relation Relation) *builder {
	left := reflect.PointerTo(b.fields[relation.Left()])
	right := reflect.PointerTo(b.fields[relation.Right()])
	lhs, rhs := relation.Types()
	// TODO: Check if `relation.Match` for fields named `relation.Lhs()` and
	//       `relation.Rhs` returns true

	if lhs != left || rhs != right {
		// TODO: Write better error message (expected vs lambda.Type())
		panic("Wrong closure type for Join")
	}

	b.relations = append(b.relations, relation)

	return b
}

func (b *builder) Where(constraint Constraint) *builder {
	// TODO: Assert that for all results `lambda` given field `name` returns true
	return b
}

func (b *builder) Any() func(*sources.DataSource) (bool, error) {
	// TODO: Build a lambda that will take k8s state and evaluate query
	// TODO: Maybe return an instance of an interface that is evaluatable

	return func(ds *sources.DataSource) (bool, error) {
		for _, relation := range b.relations {
			klog.Info(relation.Left(), " -- ", relation.Right())
			relation.Match((*scyllav1.ScyllaCluster)(nil), (*v1.Pod)(nil))
		}

		return true, nil
	}
}

func (b *builder) None() func(*sources.DataSource) (bool, error) {
	return func(ds *sources.DataSource) (bool, error) {
		return false, nil
	}
}
