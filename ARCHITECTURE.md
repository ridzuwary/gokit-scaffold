# ARCHITECTURE.md — gokit-scaffold

## System overview
gokit-scaffold is a local-first Go CLI that generates a minimal, production-grade Go service skeleton (HTTP API) with sane defaults. The goal is to reduce time-to-first-correct-service while keeping the output understandable and easy to modify.

The tool is intentionally opinionated but scope-limited: it focuses on generating a small set of templates that demonstrate clean package boundaries, predictable configuration, basic operational endpoints, and CI-friendly structure.

## Problem statement
Developers repeatedly bootstrap new Go services and re-implement the same baseline concerns (project layout, config, logging, health endpoints, graceful shutdown, Docker/CI). This causes inconsistency, wasted time, and avoidable mistakes.

gokit-scaffold provides a deterministic generator that produces a ready-to-run service skeleton in seconds.

## Non-goals
- Not a general-purpose framework or runtime dependency.
- Not a monorepo manager.
- Not a full microservices platform (no service mesh, no k8s charts by default).
- Not an AI-powered code generator.
- No cloud vendor coupling.

## Users and scale assumptions
- Primary users: developers creating small-to-medium Go services.
- Invocation frequency: occasional (per new service), not high-throughput.
- Generated repos: typical single service repo, < 10k LOC initial scaffold.

## High-level architecture

User -> CLI (gokit-scaffold)
         |
         |-- loads embedded templates (go:embed)
         |-- applies a small config model (ProjectSpec)
         |-- renders templates to target directory
         |-- writes scaffold marker (.gokit-scaffold)
         |-- runs validations (paths, module name, go version)
         '-- prints next steps (go test, go run, docker build)

## Core components and responsibilities

### cmd/gokit-scaffold (CLI entry)
- Argument parsing and subcommands:
  - `new` (create a new project)
  - `print` (show template tree / versions)
  - `validate` (validate a directory matches expected scaffold markers)
- Handles flags, reads minimal config, calls generator.

### internal/spec
- Defines `ProjectSpec` (name, module, features, ports, go version).
- Validation rules:
  - module format
  - target directory state (empty or allowed)
  - feature compatibility

### templates (embedded assets)
- Template registry
  - templates live under templates/ in the repo
  - embedded into the binary via go:embed
  - mapping of template paths -> output paths
  - generator reads from the embedded FS using stable paths like service-http/...
- Rendering:
  - Go `text/template` with safe helper funcs
  - deterministic ordering
- File writing:
  - creates dirs, writes files, sets permissions as needed
  - atomic writes where feasible (temp + rename)

### internal/generator
- Orchestrates:
  - build render plan
  - render templates
  - write outputs
  - post-generation checks (go.mod exists, main package compiles if possible)

### internal/ui
- Console output formatting, error wrapping, user guidance.

## Data model (conceptual)
The system stores no persistent data.

Conceptual structures:
- ProjectSpec:
  - ProjectName (string)
  - ModulePath (string)
  - Features (set: http, docker, ci, lint, optional db)
  - HttpPort (int)
  - GoVersion (string)
- RenderPlan:
  - Entries[]:
    - TemplatePath
    - OutputPath
    - Mode (overwrite/skip/fail)
    - ContentHash (optional for drift detection)
- ScaffoldMarker (.gokit-scaffold):
  - Tool name
  - Scaffold version
  - Template pack identifier
  - Minimal ProjectSpec snapshot

## Execution flows

### Flow: create new project
1. User runs `gokit-scaffold new ...`
2. CLI parses args -> creates ProjectSpec
3. Spec validation (module path, directory safety)
4. Generator builds render plan from embedded templates + features
5. Templates rendered -> files written
6. .gokit-scaffold marker file written to project root
7. Tool prints "next steps" and optionally runs `go test ./...` if enabled

### Flow: validate existing directory
1. User runs `gokit-scaffold validate`
2. Tool checks presence and validity of .gokit-scaffold marker
3. Tool checks presence of marker files and required structure
4. Prints actionable errors and recommended fixes

## Failure modes and safeguards
- Writing into non-empty directory:
  - default: fail unless `--force` or `--allow-nonempty` is set
- Invalid module path:
  - fail fast with clear error
- Missing or malformed .gokit-scaffold marker:
  - validate fails with explanation
- Partial write due to crash:
  - write files atomically where feasible (temp + rename)
- Template drift/bugs:
  - template rendering tests
  - golden-file tests for generated output

## Cost-control mechanisms
- No external APIs.
- Local execution only.
- Minimal dependencies.
- Deterministic outputs reduce maintenance cost.

## Security and compliance assumptions
- No PII collected or stored.
- No network calls by default.
- Generated code should include secure defaults:
  - timeouts on HTTP server
  - graceful shutdown
  - explicit /healthz and /readyz endpoints returning HTTP 200
  - safe logging (no request bodies by default)
  - clear config boundaries
- No claims of compliance (SOC2/ISO/etc).

## MVP scope vs deferred work

### MVP (5-day scope)
- CLI with `new` command
- Generate a runnable HTTP service skeleton:
  - cmd/server/main.go
  - internal/httpserver (router + handlers placeholder)
  - internal/config
  - internal/logging
  - health/readiness endpoints
  - graceful shutdown
- .gokit-scaffold marker file (JSON)
- Dockerfile (optional flag)
- GitHub Actions CI (go test) (optional flag)
- Unit tests for generator + golden tests for outputs
- README with quickstart

### Deferred
- DB wiring templates (postgres/sqlc/migrations)
- OpenTelemetry integration
- Versioned template packs / plugin system
- Project upgrade mechanism (scaffold diff/apply)

## Open questions / risks
- Windows support requirements (paths, permissions).
- Preferred router (net/http only vs chi/gorilla). MVP should default to stdlib to avoid lock-in.
- Whether to include linting/staticcheck in MVP CI (adds deps/time).
- Naming: keep output layout idiomatic and small, avoid over-structuring.

## Initial Repo Tree
.
├── ARCHITECTURE.md
├── AGENT.md
├── README.md
├── go.mod
├── cmd/
│   └── gokit-scaffold/
│       └── main.go
├── internal/
│   ├── generator/
│   ├── spec/
│   └── ui/
├── templates/                 # embedded input templates
│   └── service-http/          # template pack name
│       ├── cmd/server/main.go.tmpl
│       ├── internal/config/config.go.tmpl
│       ├── internal/httpserver/server.go.tmpl
│       ├── internal/httpserver/health.go.tmpl
│       ├── internal/httpserver/handlers.go.tmpl
│       ├── README.md.tmpl
│       └── .gitignore.tmpl
└── testdata/
    └── golden/
        └── hello-api/         # expected output tree + contents

