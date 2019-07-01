package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/d1ngd0/pvt/multiyaml"
	"github.com/d1ngd0/pvt/pivotal"
	"github.com/ghodss/yaml"
	"github.com/spf13/viper"
)

type resource struct {
	Type    string          `json:"type"`
	Project int             `json:"project"`
	Spec    json.RawMessage `json:"spec"`
}

func loadResources(r io.Reader) ([]resource, error) {
	rs := make([]resource, 0)

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return rs, err
	}

	err = multiyaml.Unmarshal(b, &rs)
	return rs, err
}

func emptyResourceYaml(t string, args []string) ([]byte, error) {
	t = normalizeType(t)
	et, err := toType(t)
	if err != nil {
		return nil, err
	}

	applyArguments(et, args)

	b, err := json.Marshal(et)
	if err != nil {
		return nil, err
	}

	res := resource{
		Type:    t,
		Project: viper.GetInt("project"),
		Spec:    b,
	}

	return resourceYaml(res)
}

func resourceYaml(res resource) ([]byte, error) {
	b, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	return yaml.JSONToYAML(b)
}

func toType(t string) (interface{}, error) {
	switch normalizeType(t) {
	case "story":
		story := &pivotal.Story{StoryType: "feature"}
		return story, nil
	case "label-blocker":
		lb := &pivotal.LabelBlocker{}
		return lb, nil
	case "epic":
		e := &pivotal.Epic{}
		return e, nil
	}

	return nil, fmt.Errorf("invalid type %s\n", t)
}

func applyArguments(t interface{}, args []string) {
	switch t.(type) {
	case *pivotal.Story:
		story := t.(*pivotal.Story)

		if args != nil {
			if len(args) > 0 {
				story.Name = args[0]
			}

			if len(args) > 1 {
				story.Description = args[1]
			}

			if len(args) > 2 {
				labels := args[2:]
				for x, l := 0, len(labels); x < l; x++ {
					story.Labels = append(story.Labels, pivotal.Label{Name: labels[x]})
				}
			}
		}
	case *pivotal.LabelBlocker:
		lb := t.(*pivotal.LabelBlocker)

		if len(args) == 2 {
			lb.Matches = args[0]
			lb.Blocks = args[1]
		}
	case *pivotal.Epic:
		e := t.(*pivotal.Epic)

		if len(args) > 0 {
			e.Name = args[0]
		}

		if len(args) > 1 {
			e.Description = args[1]
		}

		if len(args) > 2 {
			e.Label = &pivotal.Label{Name: args[2]}
		}
	}
}

func loadType(b []byte, t string) (interface{}, error) {
	v, err := toType(t)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, v)
	return v, err
}

func updateResource(res resource) error {
	r, err := loadType(res.Spec, res.Type)
	if err != nil {
		return err
	}

	switch r.(type) {
	case *pivotal.Story:
		return c.UpdateStory(res.Project, *r.(*pivotal.Story))
	case *pivotal.Epic:
		return c.UpdateEpic(res.Project, *r.(*pivotal.Epic))
	}

	return errors.New("could not find creator")

}

func createResource(res resource) error {
	r, err := loadType(res.Spec, res.Type)
	if err != nil {
		return err
	}

	switch r.(type) {
	case *pivotal.Story:
		return c.CreateStory(res.Project, *r.(*pivotal.Story))
	case *pivotal.LabelBlocker:
		return c.CreateLabelBlocker(res.Project, *r.(*pivotal.LabelBlocker))
	case *pivotal.Epic:
		return c.CreateEpic(res.Project, *r.(*pivotal.Epic))
	}

	return errors.New("could not find creator")
}

func deleteResource(project int, t string, id int) error {
	t = normalizeType(t)

	switch t {
	case "story":
		return c.DeleteStory(project, id)
	case "epic":
		return c.DeleteEpic(project, id)
	}

	return fmt.Errorf("invalid type %s", t)
}

func getResource(project int, t string, id int) (resource, error) {
	var r resource
	r.Type = normalizeType(t)
	r.Project = project

	var rspec interface{}
	var err error

	switch r.Type {
	case "story":
		rspec, err = c.GetStory(project, id)
	case "epic":
		rspec, err = c.GetEpic(project, id)
	default:
		err = fmt.Errorf("invalid type %s", r.Type)
	}

	if err != nil {
		return r, err
	}

	r.Spec, err = json.Marshal(rspec)
	return r, err
}

func normalizeType(t string) string {
	switch strings.ToLower(t) {
	case "lb", "label-blocker":
		return "label-blocker"
	case "story", "s":
		return "story"
	case "epic", "e":
		return "epic"
	}

	return ""
}
