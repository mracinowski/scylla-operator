package front

import (
	_ "embed"
	"github.com/scylladb/scylla-operator/pkg/analyze/symptoms"
	"io"
	"text/template"
)

//go:embed front.tmpl
var tmpl_string string
var tmpl *template.Template

func init() {
	tmpl = template.Must(template.New("all").Parse(tmpl_string))
}

func Print(writer io.Writer, issue symptoms.Issue) error {

	err := tmpl.ExecuteTemplate(writer, "issue", issue)

	return err
}
