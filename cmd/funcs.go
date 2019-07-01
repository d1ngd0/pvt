package cmd

import (
	"encoding/json"
	"strings"
	"text/template"

	"github.com/davecgh/go-spew/spew"
	"github.com/ghodss/yaml"
)

var templateFuncMap = template.FuncMap{
	"lower": strings.ToLower,
	"upper": strings.ToUpper,
	"dump":  spew.Sdump,
	"json": func(v interface{}) string {
		a, _ := json.MarshalIndent(v, "", "  ")
		return string(a)
	},
	"yaml": func(v interface{}) string {
		a, _ := json.MarshalIndent(v, "", "  ")
		b, _ := yaml.JSONToYAML(a)
		return string(b)
	},
}
