package multiyaml

import (
	"bytes"
	"encoding/json"
	"errors"
	"reflect"

	"github.com/ghodss/yaml"
)

// unmarshal will each document being a different thing
func Unmarshal(b []byte, v interface{}) error {
	bb := bytes.Split(b, []byte("---"))
	c := reflect.ValueOf(v)
	t := c.Type()

	if t.Kind() == reflect.Ptr {
		c = reflect.ValueOf(v).Elem()
		t = c.Type()
	}

	if t.Kind() != reflect.Slice {
		return errors.New("interface given must be of type slice")
	}

	tmp := reflect.MakeSlice(t, len(bb), len(bb))

	for x := 0; x < len(bb); x++ {
		e := reflect.New(t.Elem())
		jb, err := yaml.YAMLToJSON(bb[x])
		if err != nil {
			return err
		}

		i := e.Interface()
		if err := json.Unmarshal(jb, i); err != nil {
			return err
		}

		tmp.Index(x).Set(reflect.ValueOf(i).Elem())
	}
	c.Set(tmp)
	return nil
}

func Marshal(v interface{}) ([]byte, error) {
	bb := make([][]byte, 0)
	c := reflect.ValueOf(v)
	t := c.Type()

	if t.Kind() == reflect.Ptr {
		c = reflect.ValueOf(v).Elem()
		t = c.Type()
	}

	if t.Kind() != reflect.Slice {
		return make([]byte, 0), errors.New("interface given must be of type slice")
	}

	for x := 0; x < c.Len(); x++ {
		i := c.Index(x)
		jb, err := json.Marshal(i.Interface())
		if err != nil {
			return nil, err
		}

		if b, err := yaml.JSONToYAML(jb); err != nil {
			return make([]byte, 0), err
		} else {
			bb = append(bb, b)
		}
	}

	return bytes.Join(bb, []byte("\n---\n\n")), nil
}
