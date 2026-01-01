package requests

import (
	"io"
	"net/http"
	"net/url"
	"time"
)

// Request holds request state built from options.
type Request struct {
	method         string
	url            string
	headers        http.Header
	query          url.Values
	body           io.Reader
	timeout        time.Duration
	cookies        []*http.Cookie
	proxy          *url.URL
	redirectMax    *int
	decompressGzip bool
	err            error
}

func newRequest(method, rawURL string, opts ...Option) *Request {
	r := &Request{method: method, url: rawURL}
	for _, opt := range opts {
		if opt != nil {
			opt(r)
		}
	}
	return r
}

func (r *Request) buildURL() (*url.URL, error) {
	u, err := url.Parse(r.url)
	if err != nil {
		return nil, err
	}
	if len(r.query) == 0 {
		return u, nil
	}
	q := u.Query()
	for k, vals := range r.query {
		for _, v := range vals {
			q.Add(k, v)
		}
	}
	u.RawQuery = q.Encode()
	return u, nil
}
