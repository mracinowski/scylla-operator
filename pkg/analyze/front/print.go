package front

import (
	"github.com/scylladb/scylla-operator/pkg/analyze/symptoms"
	"io"
	"text/template"
)

func Print(writer io.Writer, issues []symptoms.Issue) error {

	tmpl, err := template.New("all").Parse(`
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
No resourcesrelated to this issue.
{{- end -}}
{{- end}}


{{define "issues" -}}
{{- if . -}}
{{- range . -}}
{{- template "issue" . -}}
---
{{end -}}
{{else -}}
No problems found
{{end}}
{{- end}}
`)
	if err != nil {
		return err
	}

	tmpl.ExecuteTemplate(writer, "issues", issues)
	//return nil
	//fmt.Println((*issues[0].Symptom).Diagnoses())

	return nil
}
