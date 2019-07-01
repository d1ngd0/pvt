package pivotal

import (
	"io"
	"net/http"
	"time"
)

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
	Get(string) (*http.Response, error)
	Post(url, contentType string, body io.Reader) (resp *http.Response, err error)
}

func newHttpClient(token string) httpClient {
	return &http.Client{
		Timeout: time.Second * 3,
		Transport: &transport{
			token: token,
		},
	}
}

type transport struct {
	token string
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	// add the base url and security token in
	req.Header.Add("X-TrackerToken", t.token)
	req.Header.Add("Content-Type", "application/json")

	return http.DefaultTransport.RoundTrip(req)
}
