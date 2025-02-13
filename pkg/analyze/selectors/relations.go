package selectors

import (
	"reflect"
)

type Relation interface {
	Left() string
	Right() string
	Types() (reflect.Type, reflect.Type)
	Match(any, any) (bool, error)
}

type FuncRelation[L, R any] struct {
	Lhs string
	Rhs string
	F   func(L, R) (bool, error)
}

func (r *FuncRelation[L, R]) Left() string {
	return r.Lhs
}

func (r *FuncRelation[L, R]) Right() string {
	return r.Rhs
}

func (r *FuncRelation[L, R]) Types() (reflect.Type, reflect.Type) {
	return reflect.TypeFor[L](), reflect.TypeFor[R]()
}

func (r *FuncRelation[L, R]) Match(lhs any, rhs any) (bool, error) {
	return r.F(lhs.(L), rhs.(R))
}
