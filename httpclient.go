package requests

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
)

// Get sends a GET request.
func Get(url string, opts ...Option) (*Response, error) { return do(http.MethodGet, url, opts...) }

// Post sends a POST request.
func Post(url string, opts ...Option) (*Response, error) { return do(http.MethodPost, url, opts...) }

// Put sends a PUT request.
func Put(url string, opts ...Option) (*Response, error) { return do(http.MethodPut, url, opts...) }

// Patch sends a PATCH request.
func Patch(url string, opts ...Option) (*Response, error) { return do(http.MethodPatch, url, opts...) }

// Delete sends a DELETE request.
func Delete(url string, opts ...Option) (*Response, error) {
	return do(http.MethodDelete, url, opts...)
}

// Head sends a HEAD request.
func Head(url string, opts ...Option) (*Response, error) { return do(http.MethodHead, url, opts...) }

// Options sends an OPTIONS request.
func Options(url string, opts ...Option) (*Response, error) {
	return do(http.MethodOptions, url, opts...)
}

func do(method, rawURL string, opts ...Option) (*Response, error) {
	req := newRequest(method, rawURL, opts...)
	if req.err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRequest, req.err)
	}
	u, err := req.buildURL()
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRequest, err)
	}
	httpReq, err := http.NewRequest(method, u.String(), req.body)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrRequest, err)
	}
	if req.headers != nil {
		httpReq.Header = req.headers
	}
	for _, c := range req.cookies {
		httpReq.AddCookie(c)
	}

	client := buildClient(req)
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, classifyErr(err)
	}

	wrapped := newResponse(resp)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return wrapped, &StatusError{StatusCode: resp.StatusCode, Response: wrapped}
	}
	return wrapped, nil
}

func buildClient(r *Request) *http.Client {
	c := &http.Client{}
	if r.timeout > 0 {
		c.Timeout = r.timeout
	}
	if r.proxy != nil {
		if base, ok := http.DefaultTransport.(*http.Transport); ok {
			tr := base.Clone()
			tr.Proxy = http.ProxyURL(r.proxy)
			c.Transport = tr
		} else {
			c.Transport = &http.Transport{Proxy: http.ProxyURL(r.proxy)}
		}
	}
	if r.redirectMax != nil {
		max := *r.redirectMax
		c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			if max <= 0 {
				return http.ErrUseLastResponse
			}
			if len(via) > max {
				return http.ErrUseLastResponse
			}
			return nil
		}
	}
	return c
}

func classifyErr(err error) error {
	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("%w: %v", ErrTimeout, err)
	}
	if ne, ok := err.(net.Error); ok && ne.Timeout() {
		return fmt.Errorf("%w: %v", ErrTimeout, err)
	}
	return fmt.Errorf("%w: %v", ErrNetwork, err)
}
