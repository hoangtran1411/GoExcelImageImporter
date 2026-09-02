---
name: wails-go-app
description: Create desktop applications using Go backend with Wails v3 framework and modern HTML/CSS/JS frontend. Lightweight alternative to Electron and Fyne.
---

# Wails v3 Desktop Application Skill

This skill provides instructions for building modern, lightweight desktop applications using **Wails v3** with a Go backend and HTML/CSS/JS frontend.

## When to Use This Skill

Use this skill when:
- Building a Go desktop application with a modern GUI
- Need a lightweight, memory-efficient alternative to Electron (~15MB vs ~150MB) or Fyne (which requires OpenGL toolchains)
- Want native Windows WebView2 integration with HTML/CSS/JS
- Requiring native desktop capabilities: file dialogs, system tray, single-instance lifecycle, and auto-updates

## Prerequisites

### 1. Install Wails v3 CLI
```bash
go install github.com/wailsapp/wails/v3/cmd/wails3@latest
```

### 2. Verify Installation
```bash
wails3 doctor
```

This checks:
- Go version (Go 1.22+)
- WebView2 runtime (Windows 10/11 native)
- C compiler / build environment

## Project Structure

Standard Wails v3 project structure (configured for Vanilla JS without bundler dependencies):

```text
project/
├── main.go              # Wails v3 entry point, window & single-instance options
├── app.go               # Wails Service (Go methods exposed to JS bindings)
├── Taskfile.yml         # Wails v3 task runner (standard build automation)
├── Makefile             # Optional dual-support Makefile wrapping Taskfile
├── go.mod               # Go module definition
├── frontend/
│   └── dist/
│       ├── index.html   # Main HTML file (loads app.js as module)
│       ├── style.css    # CSS styling
│       ├── app.js       # Frontend controller (ES Module)
│       ├── runtime.js   # Bundled standalone Wails v3 runtime
│       └── bindings/    # Generated JS bindings for Go services & models
└── build/               # Build assets, configuration, and artifacts
    ├── config.yml       # Wails v3 project config & dev watcher
    ├── devserver.go     # Lightweight stdlib dev server for live-reloading
    └── windows/         # Windows manifest, icons, and packaging scripts
```

## Core Files

### 1. main.go - Application Entry Point

```go
package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	appService := NewApp()

	app := application.New(application.Options{
		Name:        "MyApp",
		Description: "My Desktop Application",
		Services: []application.Service{
			application.NewService(appService),
		},
		Assets: application.AssetFileServerFS(assets),
		SingleInstance: &application.SingleInstanceOptions{
			UniqueID: "com.example.myapp",
			OnSecondInstanceLaunch: func(data application.SecondInstanceData) {
				if win, ok := application.Get().Window.GetByName("main"); ok {
					win.Restore()
					win.Focus()
				}
			},
		},
	})

	appService.setApp(app)

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:             "main",
		Title:            "My App",
		Width:            900,
		Height:           720,
		MinWidth:         700,
		MinHeight:        600,
		BackgroundColour: application.NewRGB(27, 38, 54),
	})

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
```

### 2. app.go - Service & Backend Logic

```go
package main

import (
	"fmt"

	"github.com/wailsapp/wails/v3/pkg/application"
)

type App struct {
	app *application.App
}

func NewApp() *App {
	return &App{}
}

func (a *App) setApp(app *application.App) {
	a.app = app
}

func (a *App) getApp() *application.App {
	if a.app != nil {
		return a.app
	}
	return application.Get()
}

// ServiceStartup is called by Wails v3 when the service starts
func (a *App) ServiceStartup(ctx application.ServiceContext) error {
	return nil
}

// Exposed methods - automatically bound to JavaScript

func (a *App) SelectFile() (string, error) {
	dialog := a.getApp().Dialog.OpenFile()
	dialog.SetTitle("Select File")
	dialog.AddFilter("All Files (*.*)", "*.*")
	return dialog.PromptForSingleSelection()
}

func (a *App) SelectFolder() (string, error) {
	dialog := a.getApp().Dialog.OpenFile()
	dialog.SetTitle("Select Folder")
	dialog.CanChooseDirectories(true)
	dialog.CanChooseFiles(false)
	return dialog.PromptForSingleSelection()
}

func (a *App) EmitProgress(percent float64) {
	a.getApp().Event.Emit("progress", percent)
}
```

### 3. frontend/dist/index.html

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>My Wails App</title>
    <link rel="stylesheet" href="style.css">
</head>
<body>
    <div id="app">
        <!-- UI Elements -->
        <button id="selectBtn" onclick="selectFile()">Browse</button>
        <button id="startBtn" onclick="startProcess()">Start</button>
    </div>

    <!-- Wails v3 ES Module entrypoint -->
    <script type="module" src="app.js"></script>
</body>
</html>
```

### 4. frontend/dist/app.js - Frontend Logic & Bindings

```javascript
import * as App from './bindings/main/app.js';
import { Events } from '/wails/runtime.js';

// Listen to backend events (payload is in e.data)
Events.On('progress', (e) => {
    console.log('Progress:', e.data);
});

// Call exposed Go methods
async function selectFile() {
    try {
        const path = await App.SelectFile();
        console.log('Selected:', path);
    } catch (err) {
        console.error('Failed to select file:', err);
    }
}

// Attach to window for inline HTML onclick handlers
window.selectFile = selectFile;
```

---

## Frontend Architecture Options

### Option A: Vanilla JS (Default & Recommended for Utility Apps)
- **Zero build dependencies**: No Node.js, npm, webpack, or Vite needed.
- Assets live directly in `frontend/dist`.
- Generates bundled runtime: `wails3 generate runtime -d frontend/dist`
- Generates ES module bindings: `wails3 generate bindings -b -d frontend/dist/bindings`
- Uses `build/devserver.go` for live-reload dev mode.

### Option B: Frontend Framework (React, Vue, Svelte, Vite)
- Uses Node/npm package manager in `frontend/`.
- Configures `frontend:run:npm` in `build/Taskfile.yml`.
- Wails proxies requests to Vite dev server (`http://localhost:5173`) in dev mode.

---

## Build Automation & Commands

### 1. Generating Assets & Bindings
```bash
# Generate bundled runtime for vanilla JS
wails3 generate runtime -d frontend/dist

# Generate typed bindings with bundled runtime references
wails3 generate bindings -b -d frontend/dist/bindings
```

### 2. Development Mode
```bash
wails3 dev
# or via Makefile
make dev
```
- Automatically builds dev binary with `-gcflags=all="-l"`.
- Runs dev server in background and watches `*.go`, `*.js`, `*.html`, `*.css`.

### 3. Production Build
```bash
# Canonical Wails v3 build
wails3 build

# Windows targeted build
wails3 task windows:build

# Release build with version injection
wails3 task windows:build VERSION=v1.0.0

# Using Makefile
make build
make build-release VERSION=v1.0.0
```

---

## Wails v3 Runtime API Reference

### Dialogs
```go
// Open file
file, err := app.Dialog.OpenFile().
    SetTitle("Select Document").
    AddFilter("Excel Files (*.xlsx)", "*.xlsx").
    PromptForSingleSelection()

// Open folder
folder, err := app.Dialog.OpenFile().
    SetTitle("Select Directory").
    CanChooseDirectories(true).
    CanChooseFiles(false).
    PromptForSingleSelection()

// Message dialog
app.Dialog.Info().SetTitle("Done").SetMessage("Completed successfully").Show()
```

### Events
```go
// Backend Emit
app.Event.Emit("eventName", payload)
```
```javascript
// Frontend Listen
import { Events } from '/wails/runtime.js';
Events.On('eventName', (e) => {
    console.log('Received:', e.data);
});
```

### Window Management
```go
win, ok := app.Window.GetByName("main")
if ok {
    win.Center()
    win.Focus()
    win.SetTitle("Updated Title")
}
```

---

## Best Practices

1. **Single Instance**: Always configure `SingleInstanceOptions` with a unique ID to prevent database/file locking conflicts.
2. **Context & Goroutines**: Heavy processing belongs in separate goroutines coordinated via channels. Do not block Wails service methods.
3. **Non-blocking Event Emission**: When emitting frequent progress updates, use non-blocking channel sends or throttle event emits to avoid flooding the webview.
4. **Clean Binding Generation**: Keep internal helper methods unexported or on unexported receiver types so they are not exposed in generated JS bindings.
