package predicate

import (
	"fmt"
	"reflect"
	"testing"
)

type NewPredicateTest struct {
	name      string
	parameter string
	lambda    any
	expected  reflect.Type
}

func TestNewPredicate(t *testing.T) {
	tests := []NewPredicateTest{
		{
			name:      "example",
			parameter: "x",
			lambda: func(int) (bool, error) {
				return false, nil
			},
			expected: reflect.TypeFor[int](),
		}, {
			name:      "lambda is nil",
			parameter: "x",
			lambda:    nil,
			expected:  nil,
		}, {
			name:      "lambda is not a func",
			parameter: "x",
			lambda:    "not a func",
			expected:  nil,
		}, {
			name:      "too few return arguments",
			parameter: "x",
			lambda: func(int) bool {
				return false
			},
			expected: nil,
		}, {
			name:      "first return value is not a bool",
			parameter: "x",
			lambda: func(int) (string, error) {
				return "not a bool", nil
			},
			expected: nil,
		}, {
			name:      "second return value is not an error",
			parameter: "x",
			lambda: func(int) (bool, string) {
				return false, "not an error"
			},
			expected: nil,
		}, {
			name:      "no parameters",
			parameter: "x",
			lambda: func() (bool, error) {
				return false, nil
			},
			expected: nil,
		}, {
			name:      "too many parameters",
			parameter: "x",
			lambda: func(int, string) (bool, error) {
				return false, nil
			},
			expected: nil,
		},
	}

	for _, test := range tests {
		p, err := New(test.parameter, test.lambda)
		if test.expected != nil && (p == nil || err != nil) {
			t.Errorf("Unexpected error: p=%p error=%s", p, err)
		}

		if test.expected != nil && p != nil {
			name, typ := p.Parameter()

			if name != test.parameter {
				t.Error("Wrong parameter name")
			}

			if typ != test.expected {
				t.Error("Wrong paramater type")
			}
		}

		if test.expected == nil && (p != nil || err == nil) {
			t.Errorf("Expected error: p=%p error=%s", p, err)
		}
	}
}

type PredicateCheckTest struct {
	name          string
	parameter     string
	lambda        any
	argument      any
	expectedValue bool
	expectedError bool
}

func TestPredicateCheck(t *testing.T) {
	tests := []PredicateCheckTest{
		{
			name:      "predicate is true",
			parameter: "x",
			lambda: func(x int) (bool, error) {
				return x == 42, nil
			},
			argument:      42,
			expectedValue: true,
			expectedError: false,
		}, {
			name:      "predicate is false",
			parameter: "x",
			lambda: func(x int) (bool, error) {
				return x == 42, nil
			},
			argument:      24,
			expectedValue: false,
			expectedError: false,
		}, {
			name:      "predicate returns error",
			parameter: "x",
			lambda: func(int) (bool, error) {
				return false, fmt.Errorf("error")
			},
			argument:      42,
			expectedValue: false,
			expectedError: true,
		},
	}

	for _, test := range tests {
		p, err := New(test.parameter, test.lambda)
		if p == nil || err != nil {
			t.Errorf("%s: Unexpected error: p=%p error=%s", test.name, p, err)
			continue
		}

		val, err := p.Test(test.argument)
		if test.expectedValue != val {
			t.Errorf("%s: Expected: %t, but got: %t",
				test.name, test.expectedValue, val)
		}

		if test.expectedError != (err != nil) {
			if test.expectedError {
				t.Errorf("%s: Expected error, but got none", test.name)
			} else {
				t.Errorf("%s: Unexpected error: %s", test.name, err)
			}
		}
	}
}
