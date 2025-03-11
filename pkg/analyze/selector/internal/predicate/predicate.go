package predicate

import (
	"fmt"
	"reflect"
)

// A Predicate holds a function from any given type to (bool, error)
// Additionally can be used as a relation which holds iff both arguments are
// the same and the predicate for both of them is true.
type Predicate struct {
	parameter string
	value     reflect.Value
}

func New(parameter string, lambda any) (*Predicate, error) {
	typ := reflect.TypeOf(lambda)

	if typ == nil || typ.Kind() != reflect.Func {
		return nil, fmt.Errorf("Not a Func")
	}

	if typ.NumOut() != 2 ||
		typ.Out(0) != reflect.TypeFor[bool]() ||
		typ.Out(1) != reflect.TypeFor[error]() {
		return nil, fmt.Errorf("Return type must be (bool, error)")
	}

	if typ.NumIn() != 1 {
		return nil, fmt.Errorf("There must be exactly one parameter")
	}

	return &Predicate{
		parameter: parameter,
		value:     reflect.ValueOf(lambda),
	}, nil
}

func (p *Predicate) Parameter() (string, reflect.Type) {
	return p.parameter, p.value.Type().In(0)
}

func (p *Predicate) FirstParameter() (string, reflect.Type) {
	return p.parameter, p.value.Type().In(0)
}

func (p *Predicate) SecondParameter() (string, reflect.Type) {
	return p.FirstParameter()
}

func (p *Predicate) Test(argument any) (bool, error) {
	if !reflect.TypeOf(argument).AssignableTo(p.value.Type().In(0)) {
		return false, fmt.Errorf("Argument not assignable")
	}

	result := p.value.Call([]reflect.Value{reflect.ValueOf(argument)})

	if result[1].IsNil() {
		return result[0].Interface().(bool), nil
	}

	return result[0].Interface().(bool), result[1].Interface().(error)
}

func (p *Predicate) Check(
	firstParameter string, firstArgument any,
	secondParameter string, secondArgument any,
) (bool, error) {
	if firstParameter != p.parameter || secondParameter != p.parameter {
		return false, fmt.Errorf("Wrong parameters %s and %s",
			firstParameter, secondParameter)
	}

	if !reflect.TypeOf(firstArgument).AssignableTo(p.value.Type().In(0)) {
		return false, fmt.Errorf("Argument %s of type %s not assignable to %s",
			firstParameter, reflect.TypeOf(firstArgument), p.value.Type().In(0))
	}

	if !reflect.TypeOf(secondArgument).AssignableTo(p.value.Type().In(0)) {
		return false, fmt.Errorf("Argument %s of type %s not assignable to %s",
			secondParameter, reflect.TypeOf(secondArgument), p.value.Type().In(0))
	}

	return p.Test(firstArgument)
}
