package httpcache

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	Client  httpClient
	Cache   Cache
	Handler *Handler
}

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

var DefaultClient = &Client{
	Client:  http.DefaultClient,
	Cache:   DefaultStorageCache,
	Handler: NewDefaultHandler(),
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	o := NewHttpResponseObject(req.URL)
	r, err := c.Handler.Pre(c.Cache, o)
	if err == nil && r != nil {
		resp := new(http.Response)
		resp.Body = io.NopCloser(r)
		resp.Status = http.StatusText(http.StatusOK)
		resp.StatusCode = http.StatusOK
		return resp, nil
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}
	if err := o.ReadResponse(resp); err != nil {
		return resp, err
	}
	if err := c.Handler.Post(c.Cache, o); err != nil {
		return resp, err
	}

	return resp, nil
}

func (c *Client) Get(url string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *Client) Post(url, contentType string, body io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", contentType)
	return c.Do(req)
}

func (c *Client) PostForm(url string, data url.Values) (resp *http.Response, err error) {
	return c.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}
