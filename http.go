package njtransit

import (
	"net/http"
	"net/url"
)

type httpClient interface {
	// Get(url string) (resp *http.Response, err error)
	// Post(url, contentType string, body io.Reader) (resp *Response, err error)
	PostForm(url string, data url.Values) (resp *http.Response, err error)
}
