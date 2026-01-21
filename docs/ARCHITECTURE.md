# Project Architecture

The project is organized following the **Wails Architecture**, combining Go (Backend) with Web Technologies (Frontend).

## 🏗️ Structure Diagram

```text
ImageToExcel/
├── main.go               # Wails Entry point, window configuration
├── app.go                # Backend Logic (Methods exposed to JS)
├── wails.json            # Wails project configuration
├── frontend/             # User Interface
│   └── dist/             # HTML/CSS/JS Assets (embedded into binary)
│       ├── index.html
│       ├── style.css
│       └── app.js
├── internal/             # Core Business Logic
│   └── engine/           # Processing Logic
│       ├── processor.go  # Excel mapping, worker pool, image insertion
│       └── processor_test.go
├── build/                # Build output directory
└── go.mod                # Dependency management (Go)
```

## ⚙️ Main Flow

1.  **Frontend (JS)**: Users interact with the HTML/CSS interface. When "Start" is clicked, JS calls the `Process()` method exposed by the Backend.
2.  **Bridge**: The Wails Bridge routes the call from JS to the Go method `Process` in `app.go`.
3.  **App Logic**: `app.go` receives the configuration and initializes the `Processor` from `internal/engine`.
4.  **Processor Phase**:
    - **Mapping**: Reads the product code column from Excel -> Map.
    - **Dispatching**: Scans the image directory and creates Jobs.
    - **Workers**: Processes images in parallel (Scaling, Decoding).
    - **Collection**: Collects results and inserts them into Excel (Single Thread safe).
5.  **Feedback**: During the process, the Backend emits `progress` events back to the Frontend. Upon completion, the Frontend displays a **Toast Notification** with detailed results.

## 🔄 Auto Update Mechanism

The auto-update system works as follows:
1.  **Check**: On startup, the Backend calls the GitHub API to check for the latest release.
2.  **Notify**: If a new version exists, it signals the Frontend to show the Update button.
3.  **Update Action**: User clicks Update -> Backend downloads the new `.exe` to a temporary folder.
4.  **Swap**: Runs a batch script to: Kill current app -> Delete old exe -> Move new exe to position -> Run new app.

## 🔒 Technical Notes
- **Wails Bridge**: Communication between JS and Go is asynchronous (Promise-based).
- **Concurrency**: Goroutines are used for heavy image processing, but writing to the Excel file must be sequential.
