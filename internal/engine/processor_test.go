package engine

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"

	"github.com/xuri/excelize/v2"
)

// Helper to create a dummy image
func createDummyImage(path string, width, height int) error {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	// Fill with blue
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.Set(x, y, color.RGBA{0, 0, 255, 255})
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return png.Encode(f, img)
}

// Helper to create a dummy excel file
func createDummyExcel(path, sheetName, codeCol string, products []string) error {
	f := excelize.NewFile()
	index, _ := f.NewSheet(sheetName)
	f.SetActiveSheet(index)

	// Write product codes
	for i, p := range products {
		cell := fmt.Sprintf("%s%d", codeCol, i+1) // Start from row 1
		_ = f.SetCellValue(sheetName, cell, p)
	}

	return f.SaveAs(path)
}

func TestProcessor_Run(t *testing.T) {
	// Setup temporary directory
	tempDir, err := os.MkdirTemp("", "imagetoexcel_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Paths
	excelPath := filepath.Join(tempDir, "test.xlsx")
	imageDir := filepath.Join(tempDir, "images")
	_ = os.Mkdir(imageDir, 0755)

	// Create dummy data
	products := []string{"P001", "P002", "P003"}
	err = createDummyExcel(excelPath, "Sheet1", "A", products)
	if err != nil {
		t.Fatal(err)
	}

	// Create dummy images
	_ = createDummyImage(filepath.Join(imageDir, "P001.png"), 100, 100)
	_ = createDummyImage(filepath.Join(imageDir, "P002.jpg"), 200, 200) // Test jpg extension (even though content is png, go image decoder handles magic numbers)
	// P003 is missing

	// Initialize Processor
	p := NewProcessor(
		excelPath,
		imageDir,
		"A", // Code column
		"B", // Image column
		"Sheet1",
		2,   // Workers
		100, // Row Height
		20,  // Col Width
	)

	// Capture progress
	progressChan := make(chan float64, 10)
	p.SetProgressChan(progressChan)

	// Run in background to consume progress
	go func() {
		for range progressChan {
			// Drain channel
		}
	}()

	// Run processor
	err = p.Run(context.Background())
	if err != nil {
		t.Errorf("Processor.Run() returned error: %v", err)
	}

	// Assertions
	if len(p.MissingCodes) != 1 {
		t.Errorf("Expected 1 missing code, got %d", len(p.MissingCodes))
	}
	if len(p.MissingCodes) > 0 && p.MissingCodes[0] != "P003" {
		t.Errorf("Expected missing code P003, got %s", p.MissingCodes[0])
	}

	// Check output file
	files, _ := os.ReadDir(tempDir)
	outputFile := ""
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".xlsx" && f.Name() != "test.xlsx" {
			outputFile = filepath.Join(tempDir, f.Name())
			break
		}
	}

	if outputFile == "" {
		t.Error("Output file not generated")
	}

	// Verify output Excel content
	f, err := excelize.OpenFile(outputFile)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	// Check if pictures exist (Excelize doesn't have easy API to list pictures in cell,
	// but we can check if file is valid and rows are resized)
	h, _ := f.GetRowHeight("Sheet1", 1)
	if h != 100 {
		t.Errorf("Expected row height 100, got %f", h)
	}
}

func TestProcessor_LoadImageData(t *testing.T) {
	tempDir, _ := os.MkdirTemp("", "img_test")
	defer os.RemoveAll(tempDir)

	imgPath := filepath.Join(tempDir, "test.png")
	_ = createDummyImage(imgPath, 50, 60)

	p := &Processor{}
	bytes, w, h, err := p.loadImageData(imgPath)

	if err != nil {
		t.Errorf("loadImageData failed: %v", err)
	}
	if w != 50 || h != 60 {
		t.Errorf("Expected size 50x60, got %dx%d", w, h)
	}
	if len(bytes) == 0 {
		t.Error("Returned empty bytes")
	}
}
