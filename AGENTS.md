# Agent Guidelines

## Commands

- **Build**: `go build`
- **Run tests**: `go test -v ./...` (CI uses `-v`; all tests live in `httpclient_test.go`)
- **Single test**: `go test -run TestFunctionName` (e.g., `go test -run TestGetWithQueryAndHeader`)
- **Lint**: `go fmt ./...` then `go vet ./...`

## Project Facts

- **Go version**: 1.25 (set in `go.mod` and CI)
- **Module**: `github.com/CareyWang/go-requests`, package name `requests`
- **Test dep**: `github.com/stretchr/testify` — use `assert.Equal/Error/ErrorIs/ErrorAs`, not `t.Error`
- **CI**: runs `go test -v ./...` on push/PR to `main` or `master`
- **No benchmarks, no linter config, no codegen**

## Code Style

- **Option pattern**: `type Option func(*Request)` — all request config is a variadic `...Option`
- **Naming**: exported = PascalCase (Get, Post, WithHeader), unexported = camelCase (do, newRequest, buildClient)
- **Variables**: concise (opts, q, srv, gz); `err` for errors
- **Comments**: doc comment on every exported symbol; minimal inline comments
- **File layout**: one responsibility per file — errors.go, options.go, request.go, response.go, session.go, httpclient.go

## Error Patterns

- **Sentinel vars**: `Err*` prefix (`ErrRequest`, `ErrNetwork`, `ErrTimeout`, `ErrStatus`, `ErrResponse`, `ErrResponseNil`, `ErrNoContent`) — check with `errors.Is(err, ErrXxx)`
- **Wrapping**: `fmt.Errorf("%w: %v", ErrType, originalErr)` — the sentinel is always the `%w` target
- **StatusError**: struct with `StatusCode` + `Response`; implements `Error()` and `Unwrap() → ErrStatus`; check with `errors.As(err, &se)`
- **Option errors**: options that can fail (WithJSON, WithProxy) set `r.err` and early-return in subsequent options; surfaced as `ErrRequest` in `do()`
- **Error classification**: `classifyErr()` maps `context.DeadlineExceeded` or `net.Error.Timeout()` → `ErrTimeout`, everything else → `ErrNetwork`

## Testing Patterns

- All tests use `httptest.NewServer` with inline handler; no external services
- Test names follow `Test<Feature>` (e.g., `TestGetWithQueryAndHeader`, `TestStatusError`, `TestTimeout`)
- Non-2xx: both `*Response` and `error` are returned — test checks both `assert.Error` and `assert.NotNil(resp)`
- `Response.Bytes()` caches body via `sync.Once` — tests verify single-read via `countingReadCloser`

## Architecture Notes

- **Context**: all public methods and Session methods take `context.Context` as the first parameter; passed through to `http.NewRequestWithContext` so caller cancellation/timeout is always respected
- **Non-2xx responses**: `do()` returns `(resp, &StatusError{...})` — caller gets the response to inspect body/headers AND an error. This is intentional.
- **Session**: `NewSession` copies the option slice; per-call options are appended after session defaults, so call-site options override session defaults (later option wins for `WithHeader`/`WithHeaders`)
- **WithJSON**: only sets `Content-Type: application/json` if header is empty; respects explicit prior set
- **WithForm**: always sets `Content-Type: application/x-www-form-urlencoded`, overriding any prior value
- **Gzip**: opt-in via `WithDecompressGzip()`; without it, raw gzip bytes are returned as-is
- **Redirect**: `WithRedirect(0)` disables redirects; `WithRedirect(N)` allows up to N hops; default follows Go's built-in behavior
- **Response body**: `sync.Once` ensures body is read exactly once and `Close()` is called; `Bytes()/Text()/JSON()` all share the cached data
