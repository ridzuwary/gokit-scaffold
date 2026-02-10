# gokit-scaffold

`gokit-scaffold` is a **local-first Go scaffold generator** for building small, production-ready HTTP services with **boring, explicit defaults**.

It exists to solve one specific problem:

> *“Every new Go service starts the same way, but everyone re-implements it slightly differently.”*

This tool generates a minimal, idiomatic Go service that:
- compiles and runs immediately
- has clear package boundaries
- exposes operational endpoints (`/healthz`, `/readyz`)
- fails safely instead of overwriting your work
- stays out of your way once generated

No frameworks. No runtime dependencies. No cloud coupling.

---

## What this project is (and is not)

### This **is**
- A **one-time scaffold generator**, not a framework
- Opinionated about *structure*, not about *business logic*
- Designed for **small to medium Go services**
- Safe, deterministic, and easy to audit
- Suitable for real projects, not demos

### This is **not**
- A code generator you run repeatedly
- A microservices platform
- A web framework
- A template zoo
- An AI-driven generator

Once the scaffold is generated, `gokit-scaffold` is no longer required.

---

## What gets generated

Running `gokit-scaffold new` produces a runnable Go service with:

- clear entrypoint (`cmd/server/main.go`)
- internal package boundaries:
  - `internal/config` — configuration loading (env-based)
  - `internal/logging` — thin logging wrapper
  - `internal/httpserver` — HTTP server + routes
- operational endpoints:
  - `GET /healthz`
  - `GET /readyz`
- graceful shutdown
- strict config validation
- a `.gokit-scaffold` marker for validation and drift detection

The generated service uses **only the Go standard library**.

---

## Example use cases

### 1. Starting a new internal service
You’re spinning up:
- a small API
- a webhook receiver
- a background service with an HTTP control plane

Instead of copying a previous repo or re-writing setup code:

```bash
gokit-scaffold new \
  --name orders-api \
  --module github.com/your-org/orders-api \
  --http-port 8080
```

You immediately get a clean, runnable baseline you can commit and build on.

---

### 2. Teaching or enforcing structure
If you work in a team and want:
- consistent project layout
- predictable startup/shutdown behavior
- clear separation of concerns

gokit-scaffold gives everyone the same starting point without enforcing a framework.

---

### 3. Replacing “copy-paste driven development”
Instead of:
- copying an old repo
- deleting files
- forgetting what’s safe to remove

You generate a fresh, minimal service every time, with only what’s required.

---

### 4. Validating existing scaffolds
If you generated a project earlier and want to ensure it still matches the expected structure:

```bash
gokit-scaffold validate --dir ./orders-api
```

This checks:
- the .gokit-scaffold marker
- required files and structure
- basic integrity of the scaffold

This is useful in CI or as a sanity check during refactors.

---

## How this fits into a development workflow

A typical flow looks like this:

1. Generate once  
   ```bash
   gokit-scaffold new --name my-service --module github.com/me/my-service
   ```

2. Commit immediately  
   ```bash
   git init  
   git add .  
   git commit -m "Initial scaffold"
   ```

3. Build real features  
   - add handlers  
   - add domain logic  
   - add storage  
   - delete anything you don’t need  

4. Ignore gokit-scaffold  
   - no runtime dependency  
   - no required updates  
   - no regeneration  

Optionally:
- use gokit-scaffold validate as a guardrail
- use gokit-scaffold print for documentation and inspection

---

## Install

Install from source:

```bash
go install github.com/ridzuwary/gokit-scaffold/cmd/gokit-scaffold@latest
```

Or run directly from the repository:

```bash
go run ./cmd/gokit-scaffold --help
```

---

## Quickstart

Generate a new HTTP service scaffold:

```bash
gokit-scaffold new --name hello-api --module github.com/example/hello-api --http-port 8080
```

Default output directory is ./<name>. You can override it:

```bash
gokit-scaffold new --name hello-api --module github.com/example/hello-api --dir ./tmp/hello-api
```

Run the generated service:

```bash
go run ./cmd/server
```

Verify endpoints:

```bash
curl http://localhost:8080/healthz  
curl http://localhost:8080/readyz
```

---

## Print

Show tool version, available template packs, generated output tree, marker schema summary, and example commands:

```bash
gokit-scaffold print
```

This is intended to be copy-pastable documentation, not debug output.

---

## Validate

Validate an existing scaffold directory:

```bash
gokit-scaffold validate --dir ./hello-api
```

Validation checks:
- .gokit-scaffold marker exists and matches required schema
- required scaffold files exist

This does not enforce how you write your business logic.

---

## Safety rule (important)

new is intentionally strict and will never overwrite existing work.

It refuses to run if:
- the target directory already exists and is non-empty
- the target path is not a directory
- the target is filesystem root
- parent traversal (..) is used

This is a deliberate design choice to prevent irreversible mistakes.

---

## Golden update policy (contributors)

Golden snapshots live under testdata/golden/hello-api and are enforced by TestGenerateGoldenHelloAPI.

Only update goldens when template output intentionally changes:

UPDATE_GOLDEN=1 go test ./internal/generator -run TestGenerateGoldenHelloAPI

Then run full tests before committing:

go test ./...

---

## Troubleshooting

### Go toolchain not found (VS Code / automation)
Ensure Go is on the system PATH, not just your shell config.

### Go build cache permission errors
In restricted environments:

GOCACHE=/tmp/go-build go test ./...
