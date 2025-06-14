package relation

import (
	"testing"
)

func TestRelationCheck(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		firstParameter  string
		secondParameter string
		lambda          any
		firstArgument   any
		secondArgument  any
		expectedValue   bool
		expectedError   error
	}{
		{
			name:            "simple",
			firstParameter:  "x",
			secondParameter: "y",
			lambda: func(x int, y int) (bool, error) {
				return true, nil
			},
			firstArgument:  1,
			secondArgument: 2,
			expectedValue:  true,
			expectedError:  nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			r, err := New(tc.firstParameter, tc.secondParameter, tc.lambda)
			if r == nil || err != nil {
				t.Fatalf("Unexpected error in New: %s (result: %p)", err, r)
			}

			val, err := r.Check(
				tc.firstParameter, tc.firstArgument,
				tc.secondParameter, tc.secondArgument,
			)

			if tc.expectedError != err {
				if tc.expectedError == nil {
					t.Fatalf("Unexpected error: %s", err)
				} else {
					t.Fatalf("Expected error: %s, but got %s", tc.expectedError, err)
				}
			}

			if tc.expectedValue != val {
				t.Fatalf("Expected: %t, but got: %t", tc.expectedValue, val)
			}
		})
	}
}
