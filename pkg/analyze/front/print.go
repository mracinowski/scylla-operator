package front

import (
	"github.com/scylladb/scylla-operator/pkg/analyze/symptoms"
	"io"
	"text/template"
)

var tmpl *template.Template

func init() {
	tmpl = template.Must(template.New("all").Parse(`
{{ define "resource" -}}
	{{- if .GetObjectKind.GroupVersionKind.Empty -}}
		{{- "\t" -}} No GVK set {{- "\n" -}}
	{{- else -}}
		{{- "\t" -}}{{- .GetObjectKind.GroupVersionKind.String -}} {{- "\n" -}}
	{{- end -}}
{{end}}

{{ define "symptom" -}}
	{{.Name}} {{- "\n" -}}
	{{- if .Diagnoses -}}
		Diagnoses: {{- "\n" -}}
		{{- range .Diagnoses -}}
			{{ "\t" }}{{.}} {{- "\n" -}}
		{{- end -}}
	{{- else -}}
		No Diagnoses {{- "\n" -}}
	{{- end -}}
	{{- if .Suggestions -}}
		Suggestions: {{- "\n" -}}
		{{- range .Suggestions -}}
			{{ "\t" }}{{.}} {{- "\n" -}}
	{{- end -}}
	{{- else -}}
		No suggestions {{- "\n" -}}
	{{- end -}}
{{end}}


{{ define "issue" -}}
	{{if .Symptom}}
		{{- template "symptom" .Symptom -}}
	{{else -}}
		No symptom {{- "\n" -}}
	{{- end}}
	{{- if .Resources -}}
		Resources GVK: {{- "\n" -}}
		{{- range .Resources -}}
			{{- template "resource" . -}}
		{{- end -}}
	{{- else -}}
		No resources related to this issue. {{- "\n" -}}
	{{- end -}}
{{- end}}
`))
}

func Print(writer io.Writer, issue symptoms.Issue) error {

	err := tmpl.ExecuteTemplate(writer, "issue", issue)

	return err
}
