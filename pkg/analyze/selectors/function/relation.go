package function

import (
	"reflect"
)

type Relation struct {
	Function[bool]
}

func NewRelation(lhs, rhs string, f any) (*Relation, error) {
	function, err := NewFunction[bool]([]string{lhs, rhs}, f)

	if function == nil {
		return nil, err
	}

	return &Relation{Function: *function}, nil
}

/* TODO: Remove
func (r *Relation) Labels() (string, string) {
	labels := make([]string, 0, 2)
	for label := range r.Function.Labels() {
		labels = append(labels, label)
	}
	return labels[0], labels[1]
}
*/

func (r *Relation) FirstParameter() (string, reflect.Type) {
	return r.labels[0], r.value.Type().In(0)
}

func (r *Relation) SecondParameter() (string, reflect.Type) {
	return r.labels[1], r.value.Type().In(1)
}

func (r *Relation) Check(
	lhsLabel string, lhsValue any,
	rhsLabel string, rhsValue any,
) bool {
	if lhsValue == nil {
		panic("lhs is nil")
	}

	if rhsValue == nil {
		panic("rhs is nil")
	}

	return r.Call(map[string]any{lhsLabel: lhsValue, rhsLabel: rhsValue})
}
