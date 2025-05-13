package front

import (
	"bytes"
	"github.com/google/go-cmp/cmp"
	"github.com/scylladb/scylla-operator/pkg/analyze/symptoms"
	//v1 "k8s.io/api/core/v1"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestPrint(t *testing.T) {

	/*
		stdout := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		defer func() {
			os.Stdout = stdout
		}()

		symptom := symptoms.NewSymptom("name", []string{"diag1", "diag2"}, []string{"sugg1", "sugg2"}, nil)
		resources := map[string]any{
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
		}
		issue := symptoms.NewIssue(&symptom, resources)

		const (
			bold       = "\033[1m"
			red        = "\033[31m"
			yellow     = "\033[93m"
			blue       = "\033[34m"
			resetStyle = "\033[0m"
		)
		expected := bold + red + "name" + resetStyle + "\n\n"
		expected += red + "diag1" + resetStyle + "\n" + red + "diag2" + resetStyle + "\n\n"
		expected += blue + "sugg1" + resetStyle + "\n" + blue + "sugg2" + resetStyle + "\n\n"
		expected += yellow + "Related resources:" + resetStyle + "\n"
		expected += yellow + "             pod -> pod1" + resetStyle + "\n"
		expected += yellow + "  serviceAccount -> serviceAccount1" + resetStyle + "\n\n"

		expected += bold + red + "name" + resetStyle + "\n"

		Print(issue, false)
		Print(issue, true)

		var buf bytes.Buffer
		w.Close()
		buf.ReadFrom(r)
		r.Close()
		output := buf.String()

		if output != expected {
			t.Errorf("expected %q got %q", expected, output)
		}
	*/

	tt := []struct {
		name     string
		sym      []symptoms.Symptom
		res      []map[string]any
		expected string
	}{
		{
			name: "Simple issue",
			sym: []symptoms.Symptom{
				symptoms.NewSymptom("name", []string{"diag1", "diag2"}, []string{"sugg1", "sugg2"}, nil),
			},
			res: []map[string]any{
				nil,
			},
			expected: "",
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			issues := make([]symptoms.Issue, 0)
			if len(tc.sym)!=len(tc.res){
				t.Fatalf("Bad test")
			}
			for i := 0; i<len(tc.sym); i++ {
				issue := symptoms.NewIssue(&(tc.sym[i]), tc.res[i])
				issues = append(issues, issue)
			}
			var buf bytes.Buffer
			Print(&buf, issues)
			got := buf.String()
			if diff := cmp.Diff(tc.expected, got); diff != "" {
				t.Errorf("Expected and actual output differ:\n%s", diff)
			}
		})
	}

}
