package function

import (
	"fmt"
	"reflect"
)

type Function[R any] struct {
	labels []string
	types  map[string]reflect.Type
	value  reflect.Value
}

func NewFunction[R any](labels []string, lambda any) (*Function[R], error) {
	typ := reflect.TypeOf(lambda)

	if typ == nil || typ.Kind() != reflect.Func {
		return nil, fmt.Errorf("Not a Func")
	}

	if typ.NumOut() != 1 {
		return nil, fmt.Errorf("Number of return values is not one")
	}

	if typ.Out(0) != reflect.TypeFor[R]() {
		return nil, fmt.Errorf("Wrong return type")
	}

	if typ.NumIn() != len(labels) {
		return nil, fmt.Errorf("Mismatched number of lables")
	}

	types := make(map[string]reflect.Type, typ.NumIn())
	for idx := range typ.NumIn() {
		types[labels[idx]] = typ.In(idx)
	}

	result := &Function[R]{
		labels: labels,
		value:  reflect.ValueOf(lambda),
	}

	// TODO: Check if all labels are unique

	return result, nil
}

func (f *Function[R]) Parameters() map[string]reflect.Type {
	result := make(map[string]reflect.Type, len(f.labels))

	typ := f.value.Type()
	for i := range typ.NumIn() {
		result[f.labels[i]] = typ.In(i)
	}

	return result
}

func (f *Function[R]) Call(args map[string]any) R {
	in := make([]reflect.Value, 0, len(f.labels))

	for i, label := range f.labels {
		arg, exists := args[label]
		if !exists {
			fmt.Printf("expected labels: %v got args: %v", f.labels, args)
			panic("TODO: Missing argument")
		}

		if arg == nil {
			in = append(in, reflect.Zero(f.value.Type().In(i)))
		} else {
			in = append(in, reflect.ValueOf(arg))
		}
	}

	result := f.value.Call(in)

	return result[0].Interface().(R)
}
