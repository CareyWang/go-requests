package requests

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"
)

// Response wraps http.Response with convenience helpers.
type Response struct {
	Raw        *http.Response
	StatusCode int
	Headers    http.Header

	once    sync.Once
	body    []byte
	bodyErr error
}

func newResponse(resp *http.Response) *Response {
	if resp == nil {
		return nil
	}
	return &Response{
		Raw:        resp,
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
	}
}

// Bytes reads and caches the response body and returns ErrResponseNil on nil responses.
func (r *Response) Bytes() ([]byte, error) {
	if r == nil || r.Raw == nil || r.Raw.Body == nil {
		return nil, ErrResponseNil
	}
	r.once.Do(func() {
		defer r.Raw.Body.Close()
		r.body, r.bodyErr = io.ReadAll(r.Raw.Body)
	})
	return r.body, r.bodyErr
}

// Text reads the response body as string.
func (r *Response) Text() (string, error) {
	b, err := r.Bytes()
	return string(b), err
}

// JSON decodes the response body into v and returns ErrNoContent on empty bodies.
func (r *Response) JSON(v any) error {
	b, err := r.Bytes()
	if err != nil {
		return err
	}
	if len(b) == 0 {
		return ErrNoContent
	}
	return json.Unmarshal(b, v)
}
