# Agent Guidelines

## Commands

- **Build**: `go build`
- **Run tests**: `go test` (all tests)
- **Single test**: `go test -run TestFunctionName` (e.g., `go test -run TestGetWithQueryAndHeader`)
- **Lint**: `go fmt ./...` and `go vet ./...`

## Code Style

- **Imports**: Group standard library imports alphabetically with blank lines between groups
- **Formatting**: Use standard `go fmt`, structure fields in logical order
- **Types**: Use Option pattern `type Option func(*Request)` for functional configuration
- **Naming**: Exported = PascalCase (Get, Post, WithHeader), internal = camelCase (do, newRequest)
- **Variables**: Concise names (opts, q, srv), err for errors
- **Errors**: Wrap with `fmt.Errorf("%w: %v", ErrType, err)`, define typed errors with Error()/Unwrap()
- **Error vars**: Use Err prefix (ErrRequest, ErrNetwork, ErrTimeout), check with errors.Is/As
- **Tests**: Use httptest.NewServer, test names prefix with Test (e.g., TestPostWithJSON)
- **Comments**: Exported symbols have doc comments, code comments are minimal
- **Files**: Single responsibility per file (errors.go, options.go, request.go, response.go, session.go)
