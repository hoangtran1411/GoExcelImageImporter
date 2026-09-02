# Go Excel Image Importer

[![Go](https://img.shields.io/badge/Go-1.25+-00ADD8?style=flat-square&logo=go&logoColor=white)](https://golang.org/)
[![Wails v3](https://img.shields.io/badge/Wails-v3_Beta-DF1B12?style=flat-square&logo=wails&logoColor=white)](https://v3.wails.io/)
[![Excelize](https://img.shields.io/badge/Excelize-v2-217346?style=flat-square&logo=microsoft-excel&logoColor=white)](https://github.com/xuri/excelize)
[![Platform](https://img.shields.io/badge/Platform-Windows_x64-0078D6?style=flat-square&logo=windows&logoColor=white)](https://microsoft.com/windows)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square)](LICENSE)

A high-performance desktop utility built with **Go** and **Wails v3**, designed to batch import and scale product images into Excel spreadsheets based on matching product codes.

---

## ✨ Key Features

- **🚀 Concurrent Worker Pool**: Multi-threaded image loading, decoding, and proportional scaling utilizing logical CPU cores.
- **📊 Row Range Filtering**: Process custom slices (`StartRow` to `EndRow`), enabling chunked batches or testing without modifying the original spreadsheet.
- **🖼️ Broad Image Format Support**: Native decoders for `.jpg`, `.jpeg`, `.png`, `.gif`, and `.webp` with automatic aspect-ratio cell fitting.
- **🔒 Single Instance Guard**: Native desktop process guarding (`com.hoangtran.goexcelimageimporter`) prevents multiple instances from corrupting open Excel workbooks.
- **🎨 Lightweight Modern UI**: Responsive Glassmorphism dark theme, real-time progress bar, sliding toast alerts, and a 1-click **"Open Output File"** action in Windows Explorer.
- **⚡ Zero NPM/Node.js Bloat**: Frontend built with pure ES Modules, standard HTML5/CSS3, and typed Wails v3 runtime bindings.
- **🔄 In-App Auto Update**: One-click update checking and in-place executable replacement directly from GitHub Releases.
- **💾 Memory Optimized**: Excelize row stream iterators and direct byte buffer decoding keep RAM consumption minimal, even on large workbooks.

---

## 💻 Compatibility & System Requirements

This application is optimized for **Windows 64-bit (x64)**.

| OS Version         | Compatibility Status | Notes                                                                                                            |
| :----------------- | :------------------- | :--------------------------------------------------------------------------------------------------------------- |
| **Windows 11**     | ✅ Recommended       | Native support (WebView2 pre-installed).                                                                         |
| **Windows 10**     | ✅ Recommended       | Native support (WebView2 pre-installed).                                                                         |
| **Windows 7 / 8**  | ⚠️ Supported         | Requires manual [WebView2 Runtime](https://developer.microsoft.com/en-us/microsoft-edge/webview2/) installation. |
| **Windows Server** | ✅ Supported         | Windows Server 2016+ with WebView2 installed.                                                                    |

---

## 🛠️ Developer Getting Started

### Prerequisites

- **Go**: `v1.25+` ([Download Go](https://go.dev/dl/))
- **Wails v3 CLI**:

  ```bash
  go install github.com/wailsapp/wails/v3/cmd/wails3@latest
  ```

- **Task** (optional, recommended) or **Make**:

  ```bash
  go install github.com/go-task/task/v3/cmd/task@latest
  ```

### Development & Live Reload

Run the app in development mode with live-reloading:

```bash
# Using Makefile
make dev

# Or using Taskfile
task dev

# Or directly via Wails v3 CLI
wails3 dev -config ./build/config.yml
```

### Building for Production

Compile a standalone optimized Windows executable:

```bash
# Using Makefile
make build-windows

# Release build with version injection (e.g., v2.1.0)
make build-release VERSION=v2.1.0

# Using Taskfile
task windows:build
```

The compiled binary (`tool_chen_anh.exe`) will be generated in `build/bin/`.

---

## 🧪 Testing & Linting

The core engine maintains **>80% test coverage**.

```bash
# Run unit tests
make test
# or
go test ./... -v

# Run tests with HTML coverage report
make coverage

# Run golangci-lint
make lint
# or
golangci-lint run
```

---

## 📖 Usage Guide

1. **Select Excel File**: Pick your source `.xlsx` spreadsheet containing product codes.
2. **Select Image Folder**: Choose the folder containing your product images (filenames should match product codes, e.g., `SKU123.jpg`).
3. **Configure Options**:
   - **Sheet Name**: Choose the target sheet from the auto-populated dropdown.
   - **Product Code Column**: Column containing codes (default: `A`).
   - **Image Target Column**: Column where images will be inserted (default: `F`).
   - **Worker Count**: Number of concurrent workers (defaults to your CPU thread count).
   - **Start / End Row**: Row boundaries (use `0` for all rows).
   - **Dimensions**: Set custom Row Height and Column Width to fit your layout.
4. **Process**: Click **Start Processing**.
5. **Output**: When complete, click **Open Output File** to immediately inspect the new timestamped file in Windows Explorer.

---

## 📂 Project Structure

```text
GoExcelImageImporter/
├── main.go               # Wails v3 entry point, window & single-instance configuration
├── app.go                # Wails v3 service methods exposed to JavaScript
├── updater.go            # GitHub Releases update checker & self-installer
├── updater_test.go       # Version comparison unit tests
├── Taskfile.yml          # Canonical Wails v3 orchestration tasks
├── Makefile              # Developer shortcuts (build, dev, test, lint)
├── internal/
│   └── engine/
│       ├── processor.go       # Core worker pool, image decoding, scaling & Excel insertion
│       └── processor_test.go  # Unit and integration tests for processor (>80% coverage)
├── frontend/
│   └── dist/
│       ├── index.html    # Application markup
│       ├── style.css     # Glassmorphism dark theme styling
│       ├── app.js        # Controller with typed Wails v3 bindings
│       ├── runtime.js    # Bundled Wails v3 runtime
│       └── bindings/     # Auto-generated JS/TS bindings for Go models
├── build/
│   ├── config.yml        # Wails v3 project metadata and file watchers
│   ├── devserver.go      # Lightweight dev server for hot reload
│   ├── bin/              # Compiled binary outputs (tool_chen_anh.exe)
│   └── windows/          # Windows icons, manifest, and installer scripts
└── docs/                 # Architectural documentation & developer guides
    ├── ARCHITECTURE.md          # Architecture & data flow diagrams
    ├── MIGRATION_WAILS_V3.md    # Wails v2 to v3 migration details
    ├── FUTURE_IMPROVEMENTS.md   # Feature roadmap
    └── LEARNING.md              # Technical insights and notes
```

---

## 📝 Documentation & References

- [Architecture Guide](docs/ARCHITECTURE.md)
- [Wails v3 Migration Notes](docs/MIGRATION_WAILS_V3.md)
- [Future Improvements & Roadmap](docs/FUTURE_IMPROVEMENTS.md)
- [Learning & Design Notes](docs/LEARNING.md)
- [Contributing Guidelines](CONTRIBUTING.md)
- [Code of Conduct](CODE_OF_CONDUCT.md)

---

## 📜 License

Distributed under the **MIT License**. See [LICENSE](LICENSE) for more information.

<p align="center">
  Made with ❤️ by <a href="https://github.com/hoangtran1411">Hoang Tran</a>
</p>
