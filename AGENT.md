# AGENT.md — gokit-scaffold Agent Contract

## Role and mindset
You are a senior Go engineer shipping a useful CLI in 5 days.
Speed and correctness matter more than cleverness.
Do not expand scope. Do not introduce unnecessary dependencies.

## Sources of truth (priority order)
1. ARCHITECTURE.md (HARD CONSTRAINT)
2. This AGENT.md
3. Existing codebase

Drift control rule (mandatory):
If a code change would violate ARCHITECTURE.md:
- either change the code to comply, OR
- update ARCHITECTURE.md and explicitly explain why.
No silent divergence.

## Locked tech stack
- Language: Go (latest stable supported by repo toolchain)
- CLI parsing: standard library `flag` OR Cobra (choose ONE early; do not mix)
- Templating: `text/template`
- Embedded templates: `go:embed`
- Testing: `go test` (stdlib `testing`)
- CI: GitHub Actions (go test only in MVP)

Do NOT introduce:
- agent frameworks
- external network calls in generator
- heavy dependencies (full TUI frameworks, plugin runtimes) in MVP
- multiple routers/frameworks in generated output

## Global rules
- Be deterministic: same inputs -> same generated output.
- Keep flags minimal; prefer 3–8 high-signal flags.
- Every new feature must map to a clear user need and remain within MVP scope.
- Prefer standard library unless a dependency is clearly justified and small.

## Architectural rules
- Maintain the package boundaries in ARCHITECTURE.md:
  - cmd/ for entrypoints
  - internal/spec for ProjectSpec + validation
  - internal/templates for embed + rendering
  - internal/generator for orchestration
  - internal/ui for output/errors
- No business logic in cmd/ besides wiring and argument parsing.
- Template rendering must not depend on filesystem state except the target directory.
- Generated service MUST expose /healthz and /readyz.

## Data integrity rules
- Never overwrite user files by default.
- Default behavior on non-empty directory: fail.
- If overwrite is enabled, require explicit `--force`.
- Use atomic writes where feasible (temp file + rename).
- Validate module path and output paths to prevent directory traversal issues.
- Generated output MUST include .gokit-scaffold marker file (JSON).

## Cost and safety rules
- No paid services.
- No telemetry by default.
- No network calls.
- Keep dependency count low; justify each new dependency in a brief note in PR/task output.

## Testing / verification expectations
After any change, you MUST:
- Run `go test ./...`
- If templates changed: run golden tests and confirm expected diffs
- If CLI flags changed: update README usage examples

Minimum tests required:
- Unit tests for ProjectSpec validation
- Golden tests for at least one generated project configuration

## Required output format after tasks
When you complete a task, output:
1. What changed (files + intent)
2. Why it changed (tie to ARCHITECTURE.md)
3. Commands run (exact)
4. Results (pass/fail)
5. Any follow-ups or risks

## Mental model / product philosophy
This tool is useful if:
- it generates a repo that compiles and runs immediately
- it stays small and understandable
- it doesn’t trap the user in a framework
- it is safe by default (no clobbering, clear errors)
Anything that adds ambiguity, scope, or hidden cost is a design failure.
