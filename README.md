# go-requests

一个轻量级的、受 Python requests 启发的 Go HTTP 客户端。API 采用函数优先设计，使用小型可组合选项。旨在覆盖最常见的用例（核心 80%），避免过度抽象。

## 特性

- 简单的 `Get/Post/Put/Patch/Delete/Head/Options` 辅助函数
- 基于选项的请求配置（headers, query, body, timeout, proxy, redirect）
- 响应辅助函数：`Bytes`, `Text`, `JSON`
- 清晰的错误语义，带有类型化错误
- 可选的 `Session` 用于设置默认值
- 可选的 gzip 自动解压

## 安装

```bash
go get github.com/CareyWang/go-requests
```

## 快速开始

```go
package main

import (
	"fmt"
	"log"

	requests "github.com/CareyWang/go-requests"
)

func main() {
	resp, err := requests.Get("https://httpbin.org/get",
		requests.WithQuery(map[string]string{"q": "golang"}),
		requests.WithHeader("X-Trace", "abc"),
	)
	if err != nil {
		log.Fatal(err)
	}
	text, _ := resp.Text()
	fmt.Println(text)
}
```

## 使用示例

### POST JSON

```go
resp, err := requests.Post("https://httpbin.org/post",
	requests.WithJSON(map[string]any{"name": "alice"}),
)
if err != nil {
	log.Fatal(err)
}
fmt.Println(resp.StatusCode)
```

### 表单数据

```go
resp, err := requests.Post("https://httpbin.org/post",
	requests.WithForm(map[string]string{"a": "1", "b": "2"}),
)
```

### 原始请求体

```go
body := strings.NewReader("raw payload")
resp, err := requests.Post("https://httpbin.org/post",
	requests.WithBody(body),
	requests.WithHeader("Content-Type", "text/plain"),
)
```

### 超时

```go
resp, err := requests.Get("https://httpbin.org/delay/2",
	requests.WithTimeout(1*time.Second),
)
if errors.Is(err, requests.ErrTimeout) {
	log.Println("timeout")
}
```

### gzip 自动解压

```go
resp, err := requests.Get("https://example.com",
	requests.WithHeader("Accept-Encoding", "gzip"),
	requests.WithDecompressGzip(),
)
text, _ := resp.Text()
```

### 重定向控制

```go
// max=0 禁用重定向
resp, err := requests.Get("https://httpbin.org/redirect/1",
	requests.WithRedirect(0),
)
```

### 代理

```go
resp, err := requests.Get("https://example.com",
	requests.WithProxy("http://127.0.0.1:8080"),
)
```

### Session 默认值

```go
s := requests.NewSession(
	requests.WithHeader("User-Agent", "httpclient"),
	requests.WithTimeout(5*time.Second),
)

resp, err := s.Get("https://httpbin.org/get")
```

## API 文档

### 顶级方法

```go
func Get(url string, opts ...Option) (*Response, error)
func Post(url string, opts ...Option) (*Response, error)
func Put(url string, opts ...Option) (*Response, error)
func Patch(url string, opts ...Option) (*Response, error)
func Delete(url string, opts ...Option) (*Response, error)
func Head(url string, opts ...Option) (*Response, error)
func Options(url string, opts ...Option) (*Response, error)
```

### Session

```go
type Session struct{}

func NewSession(opts ...Option) *Session
func (s *Session) Get(url string, opts ...Option) (*Response, error)
func (s *Session) Post(url string, opts ...Option) (*Response, error)
func (s *Session) Put(url string, opts ...Option) (*Response, error)
func (s *Session) Patch(url string, opts ...Option) (*Response, error)
func (s *Session) Delete(url string, opts ...Option) (*Response, error)
func (s *Session) Head(url string, opts ...Option) (*Response, error)
func (s *Session) Options(url string, opts ...Option) (*Response, error)
```

### Options

```go
type Option func(*Request)

func WithHeader(key, value string) Option
func WithHeaders(h map[string]string) Option
func WithQuery(q map[string]string) Option
func WithTimeout(d time.Duration) Option
func WithDecompressGzip() Option
func WithJSON(v any) Option
func WithForm(values map[string]string) Option
func WithBody(body io.Reader) Option
func WithCookies(cookies ...*http.Cookie) Option
func WithProxy(rawURL string) Option
func WithRedirect(max int) Option
```

### Response

```go
type Response struct {
	Raw        *http.Response
	StatusCode int
	Headers    http.Header
}

func (r *Response) Bytes() ([]byte, error)
func (r *Response) Text() (string, error)
func (r *Response) JSON(v any) error
```

### 错误

```go
var (
	ErrRequest = fmt.Errorf("request error")
	ErrNetwork = fmt.Errorf("network error")
	ErrTimeout = fmt.Errorf("timeout")
	ErrStatus  = fmt.Errorf("unexpected status")
	ErrResponse = fmt.Errorf("response error")
	ErrResponseNil = fmt.Errorf("nil response")
	ErrNoContent   = fmt.Errorf("empty response body")
)

type StatusError struct {
	StatusCode int
	Response   *Response
}
```

## 错误处理说明

- 非 2xx 响应返回 `*StatusError`，`errors.Is(err, ErrStatus)` 为 true
- 超时返回 `errors.Is(err, ErrTimeout)`
- 其他传输故障返回 `errors.Is(err, ErrNetwork)`
- `Response.JSON` 在空响应体时返回 `ErrNoContent`
- `Response.Bytes` 在响应或响应体为 nil 时返回 `ErrResponseNil`
- `Response.Bytes` 在读取或解压失败时返回 `ErrResponse`
- `Response.JSON` 在解码失败时返回 `ErrResponse`

## 许可证

MIT
