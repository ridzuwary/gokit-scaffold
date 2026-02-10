# gokit-scaffold

Minimal local-first Go scaffold generator.

## Usage

Create a new HTTP service scaffold:

```bash
go run ./cmd/gokit-scaffold new \
  --name hello-api \
  --module github.com/example/hello-api \
  --http-port 8080
```

Optional output directory (default `./<name>`):

```bash
go run ./cmd/gokit-scaffold new \
  --name hello-api \
  --module github.com/example/hello-api \
  --dir ./tmp/hello-api
```

Validate an existing scaffold:

```bash
go run ./cmd/gokit-scaffold validate --dir ./hello-api
```

## Troubleshooting

If `go test` fails with a Go build cache permission error in restricted environments, run with:

```bash
GOCACHE=/tmp/go-build go test ./...
```
