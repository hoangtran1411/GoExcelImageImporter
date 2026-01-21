# Go Excel Image Importer

A powerful Excel manipulation tool built with **Go** and **Wails**, designed to automatically insert images into spreadsheets based on product codes. This tool combines the raw processing power of Go with a modern, responsive Web interface (HTML/CSS/JS).

---

## ✨ Key Features

- **🚀 High Performance**: Go backend processes images and Excel files extremely fast using Worker Pools.
- **🎨 Modern UI**: Premium Dark Mode, smooth **Toast Notifications**, and fully responsive design.
- **🔄 Auto Update**: Automatically checks for and installs the latest versions from GitHub Releases.
- **💾 Memory Optimized**: Streams Excel data to handle large files without excessive RAM usage.
- **🔍 Smart Search**: Flexible matching of image filenames to product codes in the spreadsheet.
- **📦 Lightweight**: Native Windows application (~10MB) leveraging the built-in WebView2 runtime.

## 🛠️ Built With

*   [![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org/)
*   [![Wails](https://img.shields.io/badge/Wails-E34F26?style=for-the-badge&logo=wails&logoColor=white)](https://wails.io/)
*   [![Excelize](https://img.shields.io/badge/Excelize-217346?style=for-the-badge&logo=microsoft-excel&logoColor=white)](https://github.com/xuri/excelize)
*   [![JavaScript](https://img.shields.io/badge/JavaScript-F7DF1E?style=for-the-badge&logo=javascript&logoColor=black)](https://developer.mozilla.org/en-US/docs/Web/JavaScript)
*   [![HTML5](https://img.shields.io/badge/HTML5-E34F26?style=for-the-badge&logo=html5&logoColor=white)](https://developer.mozilla.org/en-US/docs/Web/Guide/HTML/HTML5)
*   [![CSS3](https://img.shields.io/badge/CSS3-1572B6?style=for-the-badge&logo=css3&logoColor=white)](https://developer.mozilla.org/en-US/docs/Web/CSS)

## 💻 Compatibility & Requirements

This tool is optimized for **Windows**.

| OS Version | Status | Notes |
| :--- | :--- | :--- |
| **Windows 11** | ✅ Best | Works perfectly (WebView2 included). |
| **Windows 10** | ✅ Best | Works perfectly (WebView2 usually included). |
| **Windows 7 / 8** | ⚠️ Limited | Requires manual installation of [WebView2 Runtime](https://developer.microsoft.com/en-us/microsoft-edge/webview2/). |
| **Windows Server** | ✅ Supported | Supported on 2016+ (requires WebView2). |

### Hardware Requirements
- **Architecture**: Windows 64-bit (x64) is required.
- **RAM**: Minimum 2GB (4GB+ recommended for processing thousands of images).

## 🚀 Installation & Setup

### Prerequisites
- [Go](https://go.dev/dl/) (v1.20+)
- [Wails CLI](https://wails.io/docs/gettingstarted/installation):
  ```bash
  go install github.com/wailsapp/wails/v2/cmd/wails@latest
  ```

### Development
1. **Clone the repository**:
   ```bash
   git clone https://github.com/hoangtran1411/GoExcelImageImporter.git
   cd GoExcelImageImporter
   ```
2. **Run in Dev Mode**:
   ```bash
   wails dev
   ```
   This will install dependencies, build the frontend, and launch the application.

### Build via Release
To generate the `.exe` file for distribution:
```bash
wails build
```
The output file will be located in `build/bin/`.

## 📖 Usage Guide

1.  **Select Excel File**: Choose the source Excel file containing your product list.
2.  **Select Image Folder**: Choose the folder containing your product images (supports .jpg, .png, .webp).
3.  **Configuration**:
    *   **Sheet Name**: Select the target sheet.
    *   **Code Column**: The column containing product codes (e.g., A).
    *   **Image Column**: The column where images should be inserted (e.g., F).
    *   **Dimensions**: Adjust Row Height and Column Width.
4.  **Start**: Click **Start Processing** and watch the progress.

## 🧪 Testing

The core logic has >80% test coverage.

```bash
# Run Unit Tests
go test ./... -v

# Or using Makefile (if available)
make test
```

## 🤝 Contributing

We welcome contributions! Please iterate through our [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on code style, formatting, and the pull request process.

## 📂 Project Structure

- `main.go`: Wails entry point and window configuration.
- `app.go`: Backend logic exposed to the frontend.
- `frontend/`: UI source code (HTML/CSS/JS).
- `internal/engine`: Core logic for Image and Excel processing.
- `wails.json`: Wails project configuration.
- `docs/`: Additional documentation (Architecture, Roadmap).

## 📝 Documentation
For more details on the internal workings and future plans, check the [docs/](docs/) folder:
- [Architecture](docs/ARCHITECTURE.md)
- [Future Improvements](docs/FUTURE_IMPROVEMENTS.md)
- [Learning Notes](docs/LEARNING.md)

---
Created by [Antigravity] - January 2026.
