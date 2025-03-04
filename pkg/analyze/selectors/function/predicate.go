package function

import (
	"reflect"
)

type Predicate struct {
	Function[bool]
}

func NewPredicate(label string, f any) (*Predicate, error) {
	function, err := NewFunction[bool]([]string{label}, f)

	if function == nil {
		return nil, err
	}

	return &Predicate{Function: *function}, nil
}

func (p *Predicate) Parameter() (string, reflect.Type) {
	return p.labels[0], p.value.Type().In(0)
}

func (p *Predicate) Check(label string, value any) bool {
	return p.Call(map[string]any{label: value})
}
