package cmd

import (
	"io"
	"text/tabwriter"
	"text/template"
)

var (
	tmpls *template.Template
)

const (
	storyTemplate = `ID	Name	State{{ range $value := . }}
{{ $value.ID }}	{{ $value.Name }}	{{ $value.CurrentState }}{{ end }}
`

	projectTemplate = `
ID	Name	Velocity{{ range $value := . }}
{{ $value.ID }}	{{ $value.Name }}	{{ $value.InitialVelocity }}{{ end }}
`

	epicTemplate = `
ID	Name	Label{{ range $value := . }}
{{ $value.ID }}	{{ $value.Name }}	{{ $value.Label.Name }}{{ end }}
`

	storyDescribeTemplate = `#{{ .ID }} {{ .Name }} :: {{ .Kind }}

{{ .Description }}
{{ .URL }}

{{ .CreatedAt }} {{ .UpdatedAt }}
`

	epicDescribeTemplate = `#{{ .ID }} {{ .Name }} :: {{ .Label.Name }}

{{ .Description }}
{{ .URL }}

{{ .CreatedAt }} {{ .UpdatedAt }}`
)

func init() {
	tmpls = template.New("get-templates").Funcs(templateFuncMap)
	template.Must(tmpls.New("story").Parse(storyTemplate))
	template.Must(tmpls.New("project").Parse(projectTemplate))
	template.Must(tmpls.New("epic").Parse(epicTemplate))
	template.Must(tmpls.New("describe-epic").Parse(epicDescribeTemplate))
	template.Must(tmpls.New("describe-story").Parse(storyDescribeTemplate))
}

func render(t string, o io.Writer, v interface{}) error {
	w := tabwriter.NewWriter(o, 8, 8, 8, ' ', 0)
	defer w.Flush()
	return tmpls.ExecuteTemplate(w, t, v)
}
