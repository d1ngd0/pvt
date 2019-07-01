package pivotal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
)

const (
	pivotalBase = "https://www.pivotaltracker.com/services/v5"
)

// Client is a struct
type Client struct {
	http httpClient
	base string
}

// New creates a new http Client
func New(token string) *Client {
	return &Client{
		http: newHttpClient(token),
		base: pivotalBase,
	}
}

func (c *Client) delete(url string) error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/%s", c.base, url), nil)

	res, err := c.http.Do(req)
	if err != nil {
		return err
	}

	return isError(res)
}

func (c *Client) put(url string, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, fmt.Sprintf("%s/%s", c.base, url), bytes.NewReader(b))

	res, err := c.http.Do(req)
	if err != nil {
		return err
	}

	return isError(res)
}

func (c *Client) get(url string, v interface{}) error {
	res, err := c.http.Get(fmt.Sprintf("%s/%s", c.base, url))
	if err != nil {
		return fmt.Errorf("failed fetching %s: %s", url, err)
	}

	if err = isError(res); err != nil {
		return nil
	}

	return unmarshalResponse(res, v)
}

func (c *Client) post(url string, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	res, err := c.http.Post(fmt.Sprintf("%s/%s", c.base, url), "application/json", bytes.NewReader(b))

	if err != nil {
		return err
	}

	return isError(res)
}

func unmarshalResponse(res *http.Response, v interface{}) error {
	b, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.New("could not read response body for request")
	}
	res.Body.Close()

	return json.Unmarshal(b, v)
}

func getFields(v interface{}) string {
	var s []string
	st := reflect.TypeOf(v)

	for x := 0; x < st.NumField(); x++ {
		f := st.Field(x)
		if tag := f.Tag.Get("json"); tag != "" {
			chunks := strings.Split(tag, ",")
			s = append(s, chunks[0])
		}
	}

	return strings.Join(s, ",")
}

func isError(res *http.Response) error {
	if res.StatusCode >= http.StatusMultipleChoices {
		b, _ := ioutil.ReadAll(res.Body)
		res.Body.Close()
		return fmt.Errorf("%s: failed: %s", res.Status, b)
	}

	return nil
}
