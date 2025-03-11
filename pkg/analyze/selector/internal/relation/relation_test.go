package relation

import (
	"testing"
)

type RelationCheckTest struct {
	name            string
	firstParameter  string
	secondParameter string
	lambda          any
	firstArgument   any
	secondArgument  any
	expectedValue   bool
	expectedError   bool
}

func TestRelationCheck(t *testing.T) {
	tests := []RelationCheckTest{
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

	for _, test := range tests {
		r, err := New(test.firstParameter, test.secondParameter, test.lambda)
		if r == nil || err != nil {
			t.Errorf("%s: Unexpected error: p=%p error=%s", test.name, r, err)
			continue
		}

		val, err := r.Check(
			test.firstParameter, test.firstArgument,
			test.secondParameter, test.secondArgument,
		)

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
