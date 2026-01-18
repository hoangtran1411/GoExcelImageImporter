package main

import (
	"context"
	"fmt"
	"imagetoexcel/internal/engine"
	"os"
	"path/filepath"
	"strings"

	stdruntime "runtime"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/xuri/excelize/v2"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
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
	file, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Excel File",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Excel Files (*.xlsx)",
				Pattern:     "*.xlsx",
			},
		},
	})
	return file, err
}

// SelectImageFolder opens a folder dialog to select the image directory
func (a *App) SelectImageFolder() (string, error) {
	folder, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select Image Folder",
	})
	return folder, err
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
		for progress := range progressChan {
			runtime.EventsEmit(a.ctx, "progress", progress*100)
		}
	}()

	// Run processing
	err := p.Run(context.Background())
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
	return nil
}

// GetCPUCount returns the number of logical CPUs
func (a *App) GetCPUCount() int {
	return stdruntime.NumCPU()
}
