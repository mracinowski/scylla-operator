package front

import (
	"github.com/scylladb/scylla-operator/pkg/analyze/symptoms"
	"io"
	"text/template"
)

func Print(writer io.Writer, issues []symptoms.Issue) error {

	t := template.Must(template.ParseFiles("pkg/analyze/front/templates.tmpl"))

	t.ExecuteTemplate(writer, "issues", issues)
	//return nil
	//fmt.Println((*issues[0].Symptom).Diagnoses())

	return nil
}
