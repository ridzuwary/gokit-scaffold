# Contributing

Thanks for contributing to `gokit-scaffold`.

## Prerequisites

- Go 1.22+

## Local workflow

1. Make your change in a focused commit.
2. Format code:
   ```bash
   gofmt -w .
   ```
3. Run tests:
   ```bash
   go test ./...
   ```
4. If generation output changed intentionally, update goldens:
   ```bash
   UPDATE_GOLDEN=1 go test ./internal/generator -run TestGenerateGoldenHelloAPI
   go test ./...
   ```

## Guardrails

- Do not change CLI semantics without explicit agreement.
- Do not change marker schema without explicit agreement.
- Keep dependencies minimal; avoid adding new dependencies unless necessary.

## Pull requests

- Include a clear problem statement and solution summary.
- Include test updates for behavior changes.
- Note whether golden files were updated and why.
