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
		expectedError   bool
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
			expectedError:  false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			r, err := New(tc.firstParameter, tc.secondParameter, tc.lambda)
			if r == nil || err != nil {
				t.Fatalf("%s: Unexpected error: p=%p error=%s", tc.name, r, err)
			}

			val, err := r.Check(
				tc.firstParameter, tc.firstArgument,
				tc.secondParameter, tc.secondArgument,
			)

			if tc.expectedValue != val {
				t.Fatalf("Expected: %t, but got: %t", tc.expectedValue, val)
			}

			if tc.expectedError != (err != nil) {
				if tc.expectedError {
					t.Fatal("Expected error, but got none")
				} else {
					t.Fatalf("Unexpected error: %s", err)
				}
			}
		})
	}
}
