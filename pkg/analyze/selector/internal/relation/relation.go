package relation

import (
	"fmt"
	"reflect"
)

// A Relation is a function from two parameters of given types to (bool, error)
type Relation interface {
	FirstParameter() (string, reflect.Type)
	SecondParameter() (string, reflect.Type)
	Check(string, any, string, any) (bool, error)
}

type relation struct {
	first  string
	second string
	value  reflect.Value
}

func New(first string, second string, lambda any) (*relation, error) {
	typ := reflect.TypeOf(lambda)

	if typ == nil || typ.Kind() != reflect.Func {
		return nil, fmt.Errorf("Not a Func")
	}

	if typ.NumOut() != 2 ||
		typ.Out(0) != reflect.TypeFor[bool]() ||
		typ.Out(1) != reflect.TypeFor[error]() {
		return nil, fmt.Errorf("Return type must be (bool, error)")
	}

	if typ.NumIn() != 2 {
		return nil, fmt.Errorf("There must be exactly two parameters")
	}

	return &relation{
		first:  first,
		second: second,
		value:  reflect.ValueOf(lambda),
	}, nil
}

func (p *relation) FirstParameter() (string, reflect.Type) {
	return p.first, p.value.Type().In(0)
}

func (p *relation) SecondParameter() (string, reflect.Type) {
	return p.second, p.value.Type().In(1)
}

func (p *relation) Check(
	firstParameter string, firstArgument any,
	secondParameter string, secondArgument any,
) (bool, error) {
	if p.first == secondParameter && p.second == firstParameter {
		firstParameter, secondParameter = secondParameter, firstParameter
		firstArgument, secondArgument = secondArgument, firstArgument
	}

	if p.first != firstParameter || p.second != secondParameter {
		return false, fmt.Errorf("Wrong parameters")
	}

	if !reflect.TypeOf(firstArgument).AssignableTo(p.value.Type().In(0)) {
		return false, fmt.Errorf("Argument %s of type %s not assignable to %s",
			firstParameter, reflect.TypeOf(firstArgument), p.value.Type().In(0))
	}

	if !reflect.TypeOf(secondArgument).AssignableTo(p.value.Type().In(1)) {
		return false, fmt.Errorf("Argument %s of type %s not assignable to %s",
			secondParameter, reflect.TypeOf(secondArgument), p.value.Type().In(1))
	}

	result := p.value.Call([]reflect.Value{
		reflect.ValueOf(firstArgument),
		reflect.ValueOf(secondArgument),
	})

	if result[1].IsNil() {
		return result[0].Interface().(bool), nil
	}

	return result[0].Interface().(bool), result[1].Interface().(error)
}
