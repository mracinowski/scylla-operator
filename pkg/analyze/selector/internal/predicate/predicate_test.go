package predicate

import (
	"fmt"
	"reflect"
	"testing"
)

func TestNewPredicate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		parameter string
		lambda    any
		expected  reflect.Type
	}{
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

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			p, err := New(tc.parameter, tc.lambda)
			if tc.expected != nil && (p == nil || err != nil) {
				t.Fatalf("Unexpected error: p=%p error=%s", p, err)
			}

			if tc.expected != nil && p != nil {
				name, ty := p.Parameter()

				if name != tc.parameter {
					t.Fatal("Wrong parameter name")
				}

				if ty != tc.expected {
					t.Fatal("Wrong paramater type")
				}
			}

			if tc.expected == nil && (p != nil || err == nil) {
				t.Fatalf("Expected error: p=%p error=%s", p, err)
			}
		})
	}
}

func TestPredicateCheck(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		parameter     string
		lambda        any
		argument      any
		expectedValue bool
		expectedError bool
	}{
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

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			p, err := New(tc.parameter, tc.lambda)
			if p == nil || err != nil {
				t.Fatalf("%s: Unexpected error: p=%p error=%s", tc.name, p, err)
			}

			val, err := p.Test(tc.argument)
			if tc.expectedValue != val {
				t.Fatalf("Expected: %t, but got: %t", tc.expectedValue, val)
			}

			if tc.expectedError != (err != nil) {
				if tc.expectedError {
					t.Fatalf("Expected error, but got none")
				} else {
					t.Fatalf("Unexpected error: %s", err)
				}
			}
		})
	}
}
