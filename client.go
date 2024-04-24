package httpcache

import (
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	*http.Client
	Cache   Cache
	Handler *Handler
}

var DefaultClient = &Client{
	Client:  http.DefaultClient,
	Cache:   DefaultStorageCache,
	Handler: NewDefaultHandler(),
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	o, err := NewHttpResponseObject(req.RequestURI)
	if err != nil {
		return nil, err
	}
	r, err := c.Handler.Pre(c.Cache, o)
	if err == nil && r != nil {
		var body io.ReadCloser = io.NopCloser(r)
		resp := new(http.Response)
		resp.Body = body
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
