package requests

import "fmt"

var (
	// ErrRequest indicates a request build or configuration error.
	ErrRequest = fmt.Errorf("request error")
	// ErrNetwork indicates a transport/network error.
	ErrNetwork = fmt.Errorf("network error")
	// ErrTimeout indicates a timeout error.
	ErrTimeout = fmt.Errorf("timeout")
	// ErrStatus indicates a non-2xx HTTP status.
	ErrStatus = fmt.Errorf("unexpected status")
	// ErrResponse indicates a response read or decode error.
	ErrResponse = fmt.Errorf("response error")
	// ErrResponseNil indicates a nil response or body.
	ErrResponseNil = fmt.Errorf("nil response")
	// ErrNoContent indicates an empty response body.
	ErrNoContent = fmt.Errorf("empty response body")
)

// StatusError is returned for non-2xx responses.
type StatusError struct {
	StatusCode int
	Response   *Response
}

func (e *StatusError) Error() string {
	return fmt.Sprintf("unexpected status: %d", e.StatusCode)
}

func (e *StatusError) Unwrap() error {
	return ErrStatus
}
