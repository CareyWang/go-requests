package requests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Option mutates a Request.
type Option func(*Request)

// WithHeader sets a single header value (overwrites existing).
func WithHeader(key, value string) Option {
	return func(r *Request) {
		if r.headers == nil {
			r.headers = make(http.Header)
		}
		r.headers.Set(key, value)
	}
}

// WithHeaders sets multiple headers (overwrites existing keys).
func WithHeaders(h map[string]string) Option {
	return func(r *Request) {
		if r.headers == nil {
			r.headers = make(http.Header)
		}
		for k, v := range h {
			r.headers.Set(k, v)
		}
	}
}

// WithQuery appends query parameters.
func WithQuery(q map[string]string) Option {
	return func(r *Request) {
		if r.query == nil {
			r.query = make(url.Values)
		}
		for k, v := range q {
			r.query.Add(k, v)
		}
	}
}

// WithTimeout sets per-request timeout.
func WithTimeout(d time.Duration) Option {
	return func(r *Request) {
		r.timeout = d
	}
}

// WithDecompressGzip enables gzip auto-decompression for response bodies.
func WithDecompressGzip() Option {
	return func(r *Request) {
		r.decompressGzip = true
	}
}

// WithJSON encodes v as JSON and sets Content-Type if missing.
func WithJSON(v any) Option {
	return func(r *Request) {
		if r.err != nil {
			return
		}
		b, err := json.Marshal(v)
		if err != nil {
			r.err = err
			return
		}
		r.body = bytes.NewReader(b)
		if r.headers == nil {
			r.headers = make(http.Header)
		}
		if r.headers.Get("Content-Type") == "" {
			r.headers.Set("Content-Type", "application/json")
		}
	}
}

// WithForm encodes form values and sets Content-Type.
func WithForm(values map[string]string) Option {
	return func(r *Request) {
		if r.err != nil {
			return
		}
		form := make(url.Values)
		for k, v := range values {
			form.Set(k, v)
		}
		r.body = strings.NewReader(form.Encode())
		if r.headers == nil {
			r.headers = make(http.Header)
		}
		r.headers.Set("Content-Type", "application/x-www-form-urlencoded")
	}
}

// WithBody sets a raw body reader.
// Note: the reader is used as-is; do not reuse it across concurrent requests.
func WithBody(body io.Reader) Option {
	return func(r *Request) {
		r.body = body
	}
}

// WithCookies adds cookies to the request.
func WithCookies(cookies ...*http.Cookie) Option {
	return func(r *Request) {
		r.cookies = append(r.cookies, cookies...)
	}
}

// WithProxy sets a proxy URL for the request.
func WithProxy(rawURL string) Option {
	return func(r *Request) {
		if r.err != nil {
			return
		}
		u, err := url.Parse(rawURL)
		if err != nil {
			r.err = err
			return
		}
		r.proxy = u
	}
}

// WithRedirect sets max redirects. max=0 disables redirects.
func WithRedirect(max int) Option {
	return func(r *Request) {
		r.redirectMax = &max
	}
}
