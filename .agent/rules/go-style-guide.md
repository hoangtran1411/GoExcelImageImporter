---
trigger: always_on
---

# Go Style Guide - GoExcelImageImporter

> **Core Rules** - For full idioms reference, see `go-idioms-reference.md`

This project is an **Excel image insertion utility** built with:
- **Wails v2** for the Desktop GUI (Windows)
- **Excelize v2** for Excel file manipulation
- **Worker Pool** for parallel image processing

---

## Code Style

- Format with `gofmt`/`goimports`. Run `golangci-lint` (v2.8.0+) `run ./...` before commit.
- **Linting Configuration**: MUST use `golangci-lint` v2 configuration schema (v2.8.x+). 
  - Top-level `version: "2"` is mandatory.
  - Use kebab-case for all linter settings.
  - Exclusions move to `linters: exclusions: rules`.
- Adhere to [Effective Go](https://go.dev/doc/effective_go) and [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments).
- Organize core processing logic in `internal/engine/`.
- Keep Wails-specific code (`main.go`, `app.go`, `updater.go`) in the root package.

## Project Structure

```
GoExcelImageImporter/
  main.go              - Wails entry point
  app.go               - Wails app bindings (Go methods exposed to JS)
  updater.go           - Auto-update functionality from GitHub Releases
  internal/
    engine/            - Core processing logic
      processor.go       - Image-to-Excel engine with Worker Pool
      processor_test.go  - Unit tests for processor
  frontend/            - Wails frontend (HTML/CSS/JS)
  build/               - Build assets and output binaries
  docs/                - Documentation and roadmap
  wails.json           - Wails project configuration
```

## Error Handling

- Wrap errors: `fmt.Errorf("context: %w", err)`. **Critical** for tracing Excel/Image I/O errors.
- Implement fail-fast logic using guard clauses to minimize indentation.
- Handle close errors in defer statements: use `defer f.Close()` pattern for `excelize.File`.
- Do not log and return the same error.

## Context and Concurrency

- Functions performing I/O or long-running operations MUST accept `context.Context` as the first argument.
- Use context to manage timeouts and cancellations (e.g. `Processor.Run`).
- Use Worker Pool pattern with `chan Job` and `chan Result` for parallel image loading.
- Use `sync.WaitGroup` for coordinating workers.

## Wails Integration

- All Wails-bound methods must be on the `*App` struct and be exported (PascalCase).
- Use `runtime.EventsEmit()` for progress updates to frontend.
- Return structs with `json` tags for frontend consumption (e.g., `Config`, `ProcessResult`).
- Use `runtime.OpenFileDialog()` and `runtime.OpenDirectoryDialog()` for native file/folder selection.

## Excel Processing (Excelize)

- Use `excelize.OpenFile()` to read existing Excel files.
- Use `f.Rows()` iterator for memory-efficient reading of large files.
- Use `f.AddPictureFromBytes()` to insert images with scaling options.
- Use `f.SetRowHeight()` and `f.SetColWidth()` to adjust cell dimensions for images.
- Always call `defer f.Close()` after opening an Excel file.

## Image Handling

- Support common image formats: `.jpg`, `.jpeg`, `.png`, `.gif`, `.webp`.
- Import blank image decoders for jpeg, png, and `golang.org/x/image/webp`.
- Use `image.DecodeConfig()` to get image dimensions without fully decoding.
- Calculate scale factors to fit images within cell dimensions.

## Testing & Linting

- Prioritize Table-driven tests combined with `t.Run`.
- Target 70% coverage for `internal/` (Minimum 40% enforced).
- Use `go test ./... -v` or `make test`.
- **Linting:** Recommended linters: `errcheck`, `gosimple`, `govet`, `ineffassign`, `staticcheck`, `unused`.
- Exclude `frontend/` and `build/` directories from linting.

---

## AI Agent Rules (Critical)

### Enforcement

- Prefer clarity over cleverness
- Prefer idiomatic Go over Java/C#/JS patterns
- If unsure, follow Effective Go first

### Context Accuracy

- Documentation links ≠ guarantees of correctness
- For external APIs: prefer explicit function signatures in context
- State assumptions when context is missing

### Library Version Awareness

- Check `go.mod` for actual versions before suggesting APIs
- LLMs hallucinate APIs for newer features not in training data
- Prefer stable APIs over experimental features

### Context Engineering

- Right context at right time, not all docs at once
- Reference existing codebase patterns first
- State missing context rather than guessing

---

## Quick Reference Links

- [Effective Go](https://go.dev/doc/effective_go)
- [Wails v2](https://github.com/wailsapp/wails)
- [Excelize v2](https://github.com/xuri/excelize)
- [golangci-lint](https://github.com/golangci/golangci-lint)
- [Go Image (x/image)](https://pkg.go.dev/golang.org/x/image)

> **Full Reference:** See `.agent/rules/go-idioms-reference.md` for detailed idioms, code examples, and best practices.
