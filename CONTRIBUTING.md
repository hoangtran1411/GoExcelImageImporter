# Contributing to GoExcelImageImporter

Thank you for your interest in contributing to **GoExcelImageImporter**! We welcome contributions from everyone. By participating in this project, you agree to abide by the code of conduct and follow the guidelines below.

## 🤝 How to Contribute

We welcome all types of contributions, including:
- **Bug Fixes**: Fixing issues reported in the issue tracker.
- **New Features**: Implement new functionality (please discuss in an issue first).
- **Documentation**: Improving README, guides, or code comments.
- **Refactoring**: Improving code quality and performance (adhering to our style guide).

### Workflow

1.  **Fork the repository** to your GitHub account.
2.  **Clone the project** locally:
    ```bash
    git clone https://github.com/hoangtran1411/GoExcelImageImporter.git
    ```
3.  **Create a new branch** for your feature or fix:
    ```bash
    git checkout -b feature/your-feature-name
    ```
    or
    ```bash
    git checkout -b fix/your-bug-fix
    ```
4.  **Make your changes**. Ensure you follow the specific **Code Style** below.
5.  **Commit your changes** with a descriptive message:
    ```bash
    git commit -m "feat: add support for dynamic image resizing"
    ```
6.  **Push to your fork**:
    ```bash
    git push origin feature/your-feature-name
    ```
7.  **Submit a Pull Request (PR)** to the `main` branch of this repository.

---

## 🎨 Code Style & Guidelines

To maintain improved code quality, please strictly follow these rules.

### Go Code Style
- **Formatting**: All Go code must be formatted using `gofmt` or `goimports`.
- **Linting**: We use `golangci-lint`. Please run `golangci-lint run ./...` before committing.
- **Structure**:
    - Core logic goes into `internal/engine/`.
    - Wails bindings go into `app.go`.
    - Main entry point is `main.go`.
- **Error Handling**: Use `fmt.Errorf("context: %w", err)` to wrap errors. Do not ignore errors.
- **Concurrency**: Use `context.Context` for long-running operations. Use Worker Pools for heavy processing.

### Frontend Style
- **Tech**: HTML, Vanilla CSS, JavaScript.
- **Design**:
    - Use variables in `style.css` for colors and spacing.
    - Ensure the UI remains responsive and supports Dark Mode.
    - Avoid external CSS frameworks (like Bootstrap) to keep the app lightweight, unless approved.

### Documentation
- Comment exported functions.
- Update `README.md` if you change behavior or add features.

---

## 🧪 Testing

- **Unit Tests**: We encourage Table-Driven Tests.
- **Coverage**: Aim for decent coverage, especially in `internal/engine/`.
- Run tests:
    ```bash
    go test ./... -v
    ```

---

## 🐛 Reporting Issues

If you find a bug, please create an Issue using the following template:
- **Description**: What happened?
- **Steps to Reproduce**: How can we see it?
- **Expected Behavior**: What should have happened?
- **Environment**: OS (Windows 10/11?), Excel Version, App Version.

---

Thank you for helping make **GoExcelImageImporter** better!
