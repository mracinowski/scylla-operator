package front

import (
	"github.com/scylladb/scylla-operator/pkg/analyze/symptoms"
	"io"
	"text/template"
)

var tmpl *template.Template

func init(){
	tmpl = template.Must(template.New("all").Parse(`
{{ define "resource" -}}
{{- "\t" -}}    Group: {{ .GetObjectKind.GroupVersionKind.Group }}
{{- "\n\t" -}}  Version: {{ .GetObjectKind.GroupVersionKind.Version }}
{{- "\n\t" -}}  Kind: {{ .GetObjectKind.GroupVersionKind.Kind }}
{{- "\n\t" -}}  ---
{{end}}


{{ define "issue" -}}
{{.Symptom.Name}}
Diagnoses:
{{- range .Symptom.Diagnoses}}
{{ "\t" }}{{.}}
{{- end -}}
{{- if .Symptom.Suggestions}}
Suggestions:
{{- range .Symptom.Suggestions}}
{{ "\t" }}{{.}}
{{- end -}}
{{else}}
No suggestions
{{- end -}}
{{if .Resources}}
Resources:
{{range .Resources -}}
{{- template "resource" . -}}
{{- end -}}
{{else}}
No resources related to this issue.
{{- end -}}
{{- end}}
`))
}

func Print(writer io.Writer, issue symptoms.Issue) error {

	err := tmpl.ExecuteTemplate(writer, "issue", issue)

	return err
}
