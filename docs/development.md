# Development

Use the Go version declared in `go.mod` and format Go files with `gofmt`.

```sh
gofmt -w $(find . -name '*.go' -type f)
go vet ./...
go test ./...
go build ./cmd/gomal-scan
```

The CLI foundation can be run with `go run ./cmd/gomal-scan --help` or `go run ./cmd/gomal-scan version`.
