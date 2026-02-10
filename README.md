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

Planned marker validation command:

```bash
go run ./cmd/gokit-scaffold validate --dir ./hello-api
```
