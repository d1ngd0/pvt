/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"
	"text/template"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

// templateCmd represents the template command
var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("you must supply at least the template")
		}

		return nil
	},
	Run: runTemplateCmd,
}

var debug bool

func init() {
	templateCmd.Flags().BoolVar(&debug, "debug", false, "output the generated yaml instead of executing it")
	createCmd.AddCommand(templateCmd)
}

func runTemplateCmd(cmd *cobra.Command, args []string) {
	b, err := ioutil.ReadFile(args[0])
	if err != nil {
		log.Fatal(fmt.Errorf("could not open file %s: %s", args[0], err))
	}

	var params []string
	if len(args) > 1 {
		params = args[1:]
	}

	tmpl, err := newTemplate(b, params)
	if err != nil {
		log.Fatal(fmt.Errorf("failed parsing template: %s", err))
	}

	var bs bytes.Buffer
	if err = tmpl.Execute(&bs); err != nil {
		log.Fatal(fmt.Errorf("could not execute template: %s", err))
	}

	if debug {
		fmt.Println(string(bs.Bytes()))
		return
	}

	if err = createResourceReader(&bs); err != nil {
		log.Fatal(err)
	}
}

func newTemplate(b []byte, args []string) (*Template, error) {
	var t Template
	if err := yaml.Unmarshal(b, &t); err != nil {
		return nil, err
	}

	tmpl, err := template.New("").Parse(t.Body)
	if err != nil {
		return nil, err
	}

	t.Template = tmpl.Funcs(templateFuncMap)

	if err = t.parseParameters(args); err != nil {
		return nil, err
	}

	return &t, nil
}

type Template struct {
	Properties templateProperties `yaml:"properties"`
	Body       string             `yaml:"template"`
	Template   *template.Template `yaml:"-"`
}

func (t *Template) parseParameters(params []string) error {
	raw, err := parseRawParameters(params)
	if err != nil {
		return err
	}

	for x := 0; x < len(t.Properties); x++ {
		p := &t.Properties[x]
		if value, ok := raw[p.Name]; ok {
			p.Value = value
		}

		if p.Required && p.Value == nil {
			return fmt.Errorf("Property %s is required", p.Name)
		}

		if p.Value == nil {
			p.Value = p.Default
		}
	}

	return nil
}

func (t *Template) Execute(w io.Writer) error {
	return t.Template.Execute(w, t.Properties)
}

type TemplateProperty struct {
	Name     string   `yaml:"name"`
	Required bool     `yaml:"required"`
	Default  []string `yaml:"default"`
	Value    []string `yaml:"-"`
}

type templateProperties []TemplateProperty

func (t templateProperties) Get(k string) string {
	return t.GetDefault(k, "")
}

func (t templateProperties) GetDefault(k, def string) string {
	v := t.GetMany(k)
	if v == nil {
		return def
	}

	return v[0]
}

func (t templateProperties) GetMany(k string) []string {
	for x, l := 0, len(t); x < l; x++ {
		if t[x].Name == k {
			return t[x].Value
		}
	}

	return nil
}

func parseRawParameters(raw []string) (map[string][]string, error) {
	v := make(map[string][]string)

	for x, l := 0, len(raw); x < l; x++ {
		chunks := strings.Split(raw[x], ":")
		if len(chunks) != 2 {
			return v, fmt.Errorf("%s is not a valid parameter, expecting key:value", raw[x])
		}

		if _, ok := v[chunks[0]]; ok {
			v[chunks[0]] = append(v[chunks[0]], chunks[1])
		} else {
			v[chunks[0]] = []string{chunks[1]}
		}
	}

	return v, nil
}
