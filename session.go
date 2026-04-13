package requests

import (
	"context"
	"net/http"
)

// Session holds default options for requests.
type Session struct {
	opts []Option
}

// NewSession creates a new session with default options.
func NewSession(opts ...Option) *Session {
	copied := make([]Option, len(opts))
	copy(copied, opts)
	return &Session{opts: copied}
}

// Get sends a GET request using session defaults.
func (s *Session) Get(ctx context.Context, url string, opts ...Option) (*Response, error) {
	return s.do(ctx, http.MethodGet, url, opts...)
}

// Post sends a POST request using session defaults.
func (s *Session) Post(ctx context.Context, url string, opts ...Option) (*Response, error) {
	return s.do(ctx, http.MethodPost, url, opts...)
}

// Put sends a PUT request using session defaults.
func (s *Session) Put(ctx context.Context, url string, opts ...Option) (*Response, error) {
	return s.do(ctx, http.MethodPut, url, opts...)
}

// Patch sends a PATCH request using session defaults.
func (s *Session) Patch(ctx context.Context, url string, opts ...Option) (*Response, error) {
	return s.do(ctx, http.MethodPatch, url, opts...)
}

// Delete sends a DELETE request using session defaults.
func (s *Session) Delete(ctx context.Context, url string, opts ...Option) (*Response, error) {
	return s.do(ctx, http.MethodDelete, url, opts...)
}

// Head sends a HEAD request using session defaults.
func (s *Session) Head(ctx context.Context, url string, opts ...Option) (*Response, error) {
	return s.do(ctx, http.MethodHead, url, opts...)
}

// Options sends an OPTIONS request using session defaults.
func (s *Session) Options(ctx context.Context, url string, opts ...Option) (*Response, error) {
	return s.do(ctx, http.MethodOptions, url, opts...)
}

func (s *Session) do(ctx context.Context, method, url string, opts ...Option) (*Response, error) {
	all := make([]Option, 0, len(s.opts)+len(opts))
	all = append(all, s.opts...)
	all = append(all, opts...)
	return do(ctx, method, url, all...)
}
