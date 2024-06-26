package httpcache

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
)

type (
	Object interface {
		Key() string
		NewReader() io.Reader
		Length() int64
	}
	HttpResponseObject struct {
		uri           string
		u             *url.URL
		body          io.Reader
		contentLength int64
	}
)

func NewHttpResponseObject(uri string) (*HttpResponseObject, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	return &HttpResponseObject{uri: uri, u: u}, nil
}
func (o *HttpResponseObject) ReadResponse(resp *http.Response) error {
	defer resp.Body.Close()
	var buf bytes.Buffer
	buf.Grow(int(resp.ContentLength))
	if _, err := io.Copy(&buf, resp.Body); err != nil {
		return err
	}
	resp.Body = io.NopCloser(bytes.NewReader(buf.Bytes()))
	o.body = bytes.NewReader(buf.Bytes())
	o.contentLength = resp.ContentLength
	return nil
}
func (o *HttpResponseObject) Key() string {
	ex := filepath.Ext(o.u.Path)
	return md5Hash(o.uri) + ex
}
func (o *HttpResponseObject) NewReader() io.Reader { return o.body }
func (o *HttpResponseObject) Length() int64        { return o.contentLength }

func md5Hash(input string) string {
	hasher := md5.New()

	hasher.Write([]byte(input))

	hashInBytes := hasher.Sum(nil)
	hashString := hex.EncodeToString(hashInBytes)

	return hashString
}
