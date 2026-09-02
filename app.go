package main

import (
	"context"
	"fmt"
	"imagetoexcel/internal/engine"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	stdruntime "runtime"

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/xuri/excelize/v2"
)

// App struct
type App struct {
	app *application.App
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// SetApp sets the Wails v3 application instance
func (a *App) SetApp(app *application.App) {
	a.app = app
}

// ServiceStartup is called when Wails v3 initializes the service
func (a *App) ServiceStartup(ctx context.Context, _ application.ServiceOptions) error {
	a.ctx = ctx
	return nil
}

func (a *App) getApp() *application.App {
	if a.app != nil {
		return a.app
	}
	return application.Get()
}

func (a *App) getContext() context.Context {
	if a.ctx != nil {
		return a.ctx
	}
	if a.app != nil {
		return a.app.Context()
	}
	return context.Background()
}

// Config holds the processing configuration
type Config struct {
	ExcelPath   string  `json:"excelPath"`
	ImageDir    string  `json:"imageDir"`
	CodeCol     string  `json:"codeCol"`
	ImageCol    string  `json:"imageCol"`
	SheetName   string  `json:"sheetName"`
	RowHeight   float64 `json:"rowHeight"`
	ColWidth    float64 `json:"colWidth"`
	WorkerCount int     `json:"workerCount"`
}

// ProcessResult holds the result of processing
type ProcessResult struct {
	Success      bool     `json:"success"`
	Message      string   `json:"message"`
	MissingCodes []string `json:"missingCodes"`
	OutputPath   string   `json:"outputPath"`
}

// SelectExcelFile opens a file dialog to select an Excel file
func (a *App) SelectExcelFile() (string, error) {
	app := a.getApp()
	dialog := app.Dialog.OpenFile()
	dialog.SetTitle("Select Excel File")
	dialog.AddFilter("Excel Files (*.xlsx)", "*.xlsx")
	return dialog.PromptForSingleSelection()
}

// SelectImageFolder opens a folder dialog to select the image directory
func (a *App) SelectImageFolder() (string, error) {
	app := a.getApp()
	dialog := app.Dialog.OpenFile()
	dialog.SetTitle("Select Image Folder")
	dialog.CanChooseDirectories(true)
	dialog.CanChooseFiles(false)
	return dialog.PromptForSingleSelection()
}

// GetSheets returns sheet names from an Excel file
func (a *App) GetSheets(excelPath string) ([]string, error) {
	if excelPath == "" {
		return []string{}, nil
	}

	f, err := excelize.OpenFile(excelPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	return f.GetSheetList(), nil
}

// Process runs the image importing process
func (a *App) Process(config Config) ProcessResult {
	// Validate inputs
	if config.ExcelPath == "" {
		return ProcessResult{Success: false, Message: "Please select an Excel file"}
	}
	if config.ImageDir == "" {
		return ProcessResult{Success: false, Message: "Please select an image folder"}
	}

	// Defaults
	if config.CodeCol == "" {
		config.CodeCol = "A"
	}
	if config.ImageCol == "" {
		config.ImageCol = "F"
	}
	if config.RowHeight <= 0 {
		config.RowHeight = 105
	}
	if config.ColWidth <= 0 {
		config.ColWidth = 20
	}
	if config.WorkerCount <= 0 {
		config.WorkerCount = 10
	}

	// Create processor
	p := engine.NewProcessor(
		config.ExcelPath,
		config.ImageDir,
		config.CodeCol,
		config.ImageCol,
		config.SheetName,
		config.WorkerCount,
		config.RowHeight,
		config.ColWidth,
	)

	// Progress channel for real-time updates
	progressChan := make(chan float64, 100)
	p.SetProgressChan(progressChan)

	// Send progress updates to frontend
	go func() {
		app := a.getApp()
		for progress := range progressChan {
			app.Event.Emit("progress", progress*100)
		}
	}()

	// Run processing
	err := p.Run(a.getContext())
	if err != nil {
		return ProcessResult{
			Success: false,
			Message: fmt.Sprintf("Processing failed: %v", err),
		}
	}

	// Get output file path
	outputPath := findOutputFile(config.ExcelPath)

	return ProcessResult{
		Success:      true,
		Message:      fmt.Sprintf("Processing completed! %d images processed, %d missing", p.ProcessedCount, len(p.MissingCodes)),
		MissingCodes: p.MissingCodes,
		OutputPath:   outputPath,
	}
}

// findOutputFile finds the most recent output file
func findOutputFile(excelPath string) string {
	dir := filepath.Dir(excelPath)
	base := strings.TrimSuffix(filepath.Base(excelPath), filepath.Ext(excelPath))

	files, err := os.ReadDir(dir)
	if err != nil {
		return ""
	}

	var latestFile string
	for _, f := range files {
		if strings.HasPrefix(f.Name(), base+"_output_") && strings.HasSuffix(f.Name(), ".xlsx") {
			latestFile = filepath.Join(dir, f.Name())
		}
	}

	return latestFile
}

// OpenFileLocation opens the file explorer to the output file location
func (a *App) OpenFileLocation(path string) error {
	if path == "" {
		return fmt.Errorf("no file path provided")
	}
	// Use Windows explorer to show the file
	cmd := exec.Command("explorer", "/select,", path)
	return cmd.Start()
}

// GetCPUCount returns the number of logical CPUs
func (a *App) GetCPUCount() int {
	return stdruntime.NumCPU()
}
