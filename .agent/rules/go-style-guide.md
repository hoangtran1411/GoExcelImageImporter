---
trigger: always_on
---

# Go Style Guide - GoExcelImageImporter

> **Core Rules** - For full idioms reference, see `go-idioms-reference.md`

This project is an **Excel image insertion utility** built with:

- **Wails v3** for the Desktop GUI (Windows)
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

```text
GoExcelImageImporter/
  main.go              - Wails v3 entry point, window & single-instance configuration
  app.go               - Wails v3 Service (Go methods exposed to JS)
  updater.go           - Auto-update functionality from GitHub Releases
  Taskfile.yml         - Wails v3 project build tasks (native build runner)
  Makefile             - Dual-support developer build automation
  internal/
    engine/            - Core processing logic
      processor.go       - Image-to-Excel engine with Worker Pool & row range filtering
      processor_test.go  - Unit tests for processor
  frontend/            - Embedded UI assets
    dist/
      index.html       - HTML markup
      style.css        - Application styling (4-column responsive config grid)
      app.js           - Frontend controller (ES Module, typed bindings)
      runtime.js       - Standalone bundled Wails v3 runtime
      bindings/        - Generated JS/TS bindings for Go structs & services
  build/               - Wails v3 build assets, icons, configs, and outputs
    config.yml         - Wails v3 project metadata and dev-mode watcher config
    devserver.go       - Lightweight stdlib dev server for live-reload dev mode
    bin/               - Compiled binaries (tool_chen_anh.exe)
  docs/                - Documentation, architecture diagrams, and roadmap
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

## Wails v3 Integration

- **Service Pattern**: Register `App` as a Wails service using `application.NewService(app)`.
- **Lifecycle**: Methods must be exported (PascalCase) on `*App`. Use `ServiceStartup` for initialization.
- **Single Instance**: Guard desktop lifecycle using `application.Options.SingleInstance` with unique ID (`com.hoangtran.goexcelimageimporter`).
- **Dialogs**: Use `app.Dialog.OpenFile()` for native file and folder selection dialogs.
- **Events**: Emit updates using `app.Event.Emit("eventName", data)` and listen in JS via `Events.On("eventName", (e) => ...)` (data in `e.data`).
- **Data Transfer**: Return structs with `json` tags for frontend consumption (e.g., `Config`, `ProcessResult`).
- **Build Automation**: Maintain dual support: `Taskfile.yml` for canonical Wails v3 orchestration, `Makefile` for developer shortcuts (`make build`, `make dev`, `make test`).

## Excel Processing (Excelize)

- Use `excelize.OpenFile()` to read existing Excel files.
- Use `f.Rows()` iterator for memory-efficient reading of large files.
- Use `f.AddPictureFromBytes()` to insert images with scaling options.
- Use `f.SetRowHeight()` and `f.SetColWidth()` to adjust cell dimensions for images.
- Always call `defer f.Close()` after opening an Excel file.
- Support row range filtering (`StartRow` and `EndRow`, default 0 = all rows).

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

- Check `go.mod` for actual versions before suggesting APIs (Wails v3 beta)
- LLMs hallucinate APIs for newer features not in training data
- Prefer stable APIs over experimental features

### Context Engineering

- Right context at right time, not all docs at once
- Reference existing codebase patterns first
- State missing context rather than guessing

---

## Quick Reference Links

- [Effective Go](https://go.dev/doc/effective_go)
- [Wails v3](https://v3.wails.io/)
- [Excelize v2](https://github.com/xuri/excelize)
- [golangci-lint](https://github.com/golangci/golangci-lint)
- [Go Image (x/image)](https://pkg.go.dev/golang.org/x/image)

> **Full Reference:** See `.agent/rules/go-idioms-reference.md` for detailed idioms, code examples, and best practices.
