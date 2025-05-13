package front

import (
	"bytes"
	"github.com/google/go-cmp/cmp"
	"github.com/scylladb/scylla-operator/pkg/analyze/symptoms"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestPrint(t *testing.T) {
	tt := []struct {
		name     string
		symptom      symptoms.Symptom
		resources      map[string]any
		expected string
	}{
		{
			name: "Simple issue",
			symptom : symptoms.NewSymptom("name", []string{"diag1", "diag2"}, []string{"sugg1", "sugg2"}, nil),
			resources : map[string]any{
				"pod": &v1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name: "pod1",
					},
				},
				"serviceAccount": &v1.ServiceAccount{
					ObjectMeta: metav1.ObjectMeta{
						Name: "serviceAccount1",
					},
				},
			},
			expected: "",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			issue := symptoms.NewIssue(&(tc.symptom), tc.resources)
			var buf bytes.Buffer
			Print(&buf, issue)
			got := buf.String()
			if diff := cmp.Diff(tc.expected, got); diff != "" {
				t.Errorf("Expected and actual output differ:\n%s", diff)
			}
		})
	}

}
