package requests

import "net/http"

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
func (s *Session) Get(url string, opts ...Option) (*Response, error) {
	return s.do(http.MethodGet, url, opts...)
}

// Post sends a POST request using session defaults.
func (s *Session) Post(url string, opts ...Option) (*Response, error) {
	return s.do(http.MethodPost, url, opts...)
}

// Put sends a PUT request using session defaults.
func (s *Session) Put(url string, opts ...Option) (*Response, error) {
	return s.do(http.MethodPut, url, opts...)
}

// Patch sends a PATCH request using session defaults.
func (s *Session) Patch(url string, opts ...Option) (*Response, error) {
	return s.do(http.MethodPatch, url, opts...)
}

// Delete sends a DELETE request using session defaults.
func (s *Session) Delete(url string, opts ...Option) (*Response, error) {
	return s.do(http.MethodDelete, url, opts...)
}

// Head sends a HEAD request using session defaults.
func (s *Session) Head(url string, opts ...Option) (*Response, error) {
	return s.do(http.MethodHead, url, opts...)
}

// Options sends an OPTIONS request using session defaults.
func (s *Session) Options(url string, opts ...Option) (*Response, error) {
	return s.do(http.MethodOptions, url, opts...)
}

func (s *Session) do(method, url string, opts ...Option) (*Response, error) {
	all := make([]Option, 0, len(s.opts)+len(opts))
	all = append(all, s.opts...)
	all = append(all, opts...)
	return do(method, url, all...)
}
