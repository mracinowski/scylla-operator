package function

import (
	"testing"
)

type newFunctionErrorsTest struct {
	labels []string
	lambda any
}

var newFunctionErrorsTests = []newFunctionTest{
	newFunctionErrorsTest{labels: []string{"alfa"}, lambda: nil},
}

func TestNewFunctionErrors(t *testing.T) {
	for _, test := range newFunctionTests {
		ret, err := NewFunction[string](test.labels, test.lambda)
		if ret != nil || err == nil {
			t.Errorf("TODO: Mismatch")
			continue
		}
	}

	t.Errorf("TODO: Write the test")
}
