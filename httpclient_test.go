package requests

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetWithQueryAndHeader(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "x", r.URL.Query().Get("q"))
		assert.Equal(t, "1", r.Header.Get("X-Test"))
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.WriteString(w, `{"ok":true}`)
	}))
	defer srv.Close()

	resp, err := Get(srv.URL, WithQuery(map[string]string{"q": "x"}), WithHeader("X-Test", "1"))
	assert.NoError(t, err)

	var out struct {
		Ok bool `json:"ok"`
	}
	assert.NoError(t, resp.JSON(&out))
	assert.True(t, out.Ok)
}

func TestPostWithJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		var payload map[string]any
		assert.NoError(t, json.NewDecoder(r.Body).Decode(&payload))
		assert.Equal(t, "alice", payload["name"])
		w.WriteHeader(http.StatusCreated)
	}))
	defer srv.Close()

	resp, err := Post(srv.URL, WithJSON(map[string]any{"name": "alice"}))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
}

func TestStatusError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = io.WriteString(w, "bad")
	}))
	defer srv.Close()

	resp, err := Get(srv.URL)
	assert.Error(t, err)
	var se *StatusError
	assert.ErrorAs(t, err, &se)
	assert.ErrorIs(t, err, ErrStatus)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestTimeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		_, _ = io.WriteString(w, "late")
	}))
	defer srv.Close()

	_, err := Get(srv.URL, WithTimeout(50*time.Millisecond))
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrTimeout)
}

func TestJSONEmptyBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	resp, err := Get(srv.URL)
	assert.NoError(t, err)
	var out map[string]any
	assert.ErrorIs(t, resp.JSON(&out), ErrNoContent)
}

func TestBytesNilResponse(t *testing.T) {
	var resp *Response
	_, err := resp.Bytes()
	assert.ErrorIs(t, err, ErrResponseNil)
}

func TestText(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		_, _ = io.WriteString(w, "hello world")
	}))
	defer srv.Close()

	resp, err := Get(srv.URL)
	assert.NoError(t, err)
	text, err := resp.Text()
	assert.NoError(t, err)
	assert.Equal(t, "hello world", text)
}

func TestGzipAutoDecompress(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		gz := gzip.NewWriter(&buf)
		_, _ = gz.Write([]byte("hello gzip"))
		assert.NoError(t, gz.Close())
		w.Header().Set("Content-Encoding", "gzip")
		_, _ = w.Write(buf.Bytes())
	}))
	defer srv.Close()

	resp, err := Get(srv.URL, WithHeader("Accept-Encoding", "gzip"), WithDecompressGzip())
	assert.NoError(t, err)
	text, err := resp.Text()
	assert.NoError(t, err)
	assert.Equal(t, "hello gzip", text)
}

func TestGzipWithoutDecompressOption(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		gz := gzip.NewWriter(&buf)
		_, _ = gz.Write([]byte("hello gzip"))
		assert.NoError(t, gz.Close())
		w.Header().Set("Content-Encoding", "gzip")
		_, _ = w.Write(buf.Bytes())
	}))
	defer srv.Close()

	resp, err := Get(srv.URL, WithHeader("Accept-Encoding", "gzip"))
	assert.NoError(t, err)
	b, err := resp.Bytes()
	assert.NoError(t, err)
	assert.True(t, len(b) >= 2)
	assert.Equal(t, byte(0x1f), b[0])
	assert.Equal(t, byte(0x8b), b[1])
}

func TestPut(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPut, r.Method)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	_, err := Put(srv.URL)
	assert.NoError(t, err)
}

func TestPatch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPatch, r.Method)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	_, err := Patch(srv.URL)
	assert.NoError(t, err)
}

func TestDelete(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodDelete, r.Method)
		w.WriteHeader(http.StatusNoContent)
	}))
	defer srv.Close()

	_, err := Delete(srv.URL)
	assert.NoError(t, err)
}

func TestHead(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodHead, r.Method)
		w.Header().Set("Content-Length", "5")
	}))
	defer srv.Close()

	resp, err := Head(srv.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestOptions(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodOptions, r.Method)
		w.Header().Set("Allow", "GET, POST")
	}))
	defer srv.Close()

	_, err := Options(srv.URL)
	assert.NoError(t, err)
}

func TestWithHeaders(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "value", r.Header.Get("X-Custom"))
		assert.Equal(t, "test", r.Header.Get("X-Another"))
	}))
	defer srv.Close()

	_, err := Get(srv.URL, WithHeaders(map[string]string{
		"X-Custom":  "value",
		"X-Another": "test",
	}))
	assert.NoError(t, err)
}

func TestWithHeaderOverwrite(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "second", r.Header.Get("X-Key"))
	}))
	defer srv.Close()

	_, err := Get(srv.URL, WithHeader("X-Key", "first"), WithHeader("X-Key", "second"))
	assert.NoError(t, err)
}

func TestWithForm(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
		assert.Equal(t, "bob", r.FormValue("name"))
	}))
	defer srv.Close()

	_, err := Post(srv.URL, WithForm(map[string]string{"name": "bob"}))
	assert.NoError(t, err)
}

func TestWithFormOverridesContentType(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))
	}))
	defer srv.Close()

	_, err := Post(srv.URL, WithHeader("Content-Type", "text/plain"), WithForm(map[string]string{"name": "bob"}))
	assert.NoError(t, err)
}

func TestWithJSONRespectsExistingContentType(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "text/plain", r.Header.Get("Content-Type"))
	}))
	defer srv.Close()

	_, err := Post(srv.URL, WithHeader("Content-Type", "text/plain"), WithJSON(map[string]any{"name": "alice"}))
	assert.NoError(t, err)
}

func TestWithBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		assert.Equal(t, "raw body", string(b))
	}))
	defer srv.Close()

	_, err := Post(srv.URL, WithBody(strings.NewReader("raw body")))
	assert.NoError(t, err)
}

func TestWithCookies(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session")
		assert.NoError(t, err)
		assert.Equal(t, "abc123", c.Value)
	}))
	defer srv.Close()

	_, err := Get(srv.URL, WithCookies(&http.Cookie{Name: "session", Value: "abc123"}))
	assert.NoError(t, err)
}

func TestWithProxy(t *testing.T) {
	proxy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))
	defer proxy.Close()

	target := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer target.Close()

	_, err := Get(target.URL, WithProxy(proxy.URL))
	assert.Error(t, err)
	var se *StatusError
	assert.ErrorAs(t, err, &se)
	assert.Equal(t, http.StatusTeapot, se.StatusCode)
}

func TestWithRedirect(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/redirect" {
			http.Redirect(w, r, "/target", http.StatusFound)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer srv.Close()

	_, err := Get(srv.URL+"/redirect", WithRedirect(0))
	assert.Error(t, err)
	var se *StatusError
	assert.ErrorAs(t, err, &se)
	assert.Equal(t, http.StatusFound, se.StatusCode)
}

func TestWithRedirectLimit(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/a":
			http.Redirect(w, r, "/b", http.StatusFound)
		case "/b":
			http.Redirect(w, r, "/c", http.StatusFound)
		default:
			w.WriteHeader(http.StatusOK)
		}
	}))
	defer srv.Close()

	_, err := Get(srv.URL+"/a", WithRedirect(1))
	assert.Error(t, err)
	var se *StatusError
	assert.ErrorAs(t, err, &se)
	assert.Equal(t, http.StatusFound, se.StatusCode)
}

func TestWithRedirectAllow(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/redirect" {
			http.Redirect(w, r, "/ok", http.StatusFound)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	resp, err := Get(srv.URL+"/redirect", WithRedirect(1))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestSession(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "value", r.Header.Get("X-Session"))
		assert.Equal(t, http.MethodPost, r.Method)
	}))
	defer srv.Close()

	session := NewSession(WithHeader("X-Session", "value"))
	_, err := session.Post(srv.URL)
	assert.NoError(t, err)
}

func TestSessionMethods(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "value", r.Header.Get("X-Session"))
		assert.Contains(t, []string{http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodHead, http.MethodOptions}, r.Method)
	}))
	defer srv.Close()

	session := NewSession(WithHeader("X-Session", "value"))

	_, err := session.Put(srv.URL)
	assert.NoError(t, err)
	_, err = session.Patch(srv.URL)
	assert.NoError(t, err)
	_, err = session.Delete(srv.URL)
	assert.NoError(t, err)
	_, err = session.Head(srv.URL)
	assert.NoError(t, err)
	_, err = session.Options(srv.URL)
	assert.NoError(t, err)
}

func TestSessionOverrideOptions(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "override", r.Header.Get("X-Key"))
	}))
	defer srv.Close()

	session := NewSession(WithHeader("X-Key", "default"))
	_, err := session.Get(srv.URL, WithHeader("X-Key", "override"))
	assert.NoError(t, err)
}

func TestBuildURLWithExistingQuery(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		assert.Equal(t, 2, len(q["a"]))
		assert.Equal(t, "2", q.Get("b"))
	}))
	defer srv.Close()

	_, err := Get(srv.URL+"?a=1", WithQuery(map[string]string{"a": "3", "b": "2"}))
	assert.NoError(t, err)
}

func TestInvalidURL(t *testing.T) {
	_, err := Get("://invalid")
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrRequest)
}

func TestNetworkError(t *testing.T) {
	_, err := Get("http://localhost:9999")
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrNetwork)
}

type timeoutErr struct{}

func (timeoutErr) Error() string   { return "timeout" }
func (timeoutErr) Timeout() bool   { return true }
func (timeoutErr) Temporary() bool { return false }

type timeoutTransport struct{}

func (timeoutTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, timeoutErr{}
}

func TestNetworkTimeoutError(t *testing.T) {
	prev := http.DefaultTransport
	http.DefaultTransport = timeoutTransport{}
	t.Cleanup(func() {
		http.DefaultTransport = prev
	})

	_, err := Get("http://example.com")
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrTimeout)
}

func TestWithJSONError(t *testing.T) {
	ch := make(chan int)
	defer close(ch)

	_, err := Post("http://example.com", WithJSON(ch))
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrRequest)
}

func TestWithProxyError(t *testing.T) {
	_, err := Get("http://example.com", WithProxy("://invalid"))
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrRequest)
}

type countingReadCloser struct {
	readCount int
	closed    bool
	data      []byte
}

func (c *countingReadCloser) Read(p []byte) (int, error) {
	if len(c.data) == 0 {
		return 0, io.EOF
	}
	n := copy(p, c.data)
	c.data = c.data[n:]
	c.readCount++
	return n, nil
}

func (c *countingReadCloser) Close() error {
	c.closed = true
	return nil
}

func TestBytesCachesBody(t *testing.T) {
	rc := &countingReadCloser{data: []byte("hello")}
	resp := &Response{Raw: &http.Response{Body: rc}}

	first, err := resp.Bytes()
	assert.NoError(t, err)
	second, err := resp.Bytes()
	assert.NoError(t, err)
	assert.Equal(t, "hello", string(first))
	assert.Equal(t, "hello", string(second))
	assert.Equal(t, 1, rc.readCount)
	assert.True(t, rc.closed)
}

func TestJSONInvalidBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, "not json")
	}))
	defer srv.Close()

	resp, err := Get(srv.URL)
	assert.NoError(t, err)
	var out map[string]any
	err = resp.JSON(&out)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrResponse)
}

func TestTextNilResponse(t *testing.T) {
	var resp *Response
	_, err := resp.Text()
	assert.ErrorIs(t, err, ErrResponseNil)
}
