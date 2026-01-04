package engine

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/xuri/excelize/v2"
)

func TestProcessor_Mapping(t *testing.T) {
	// Create dummy excel
	f := excelize.NewFile()
	sheet := "Sheet1"
	f.SetCellValue(sheet, "A1", "CODE001")
	f.SetCellValue(sheet, "A2", "CODE002")
	f.SetCellValue(sheet, "A3", "")
	f.SetCellValue(sheet, "A4", "CODE003")

	tmpExcel := filepath.Join(t.TempDir(), "test.xlsx")
	if err := f.SaveAs(tmpExcel); err != nil {
		t.Fatal(err)
	}

	// Create dummy image folder with realistic images (minimum possible size)
	// Actually, just small files are enough if we mock loadImageData or use real small images.
	// For testing, let's create a 1x1 pixel PNG
	tmpDir := t.TempDir()
	pngData := []byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR\x00\x00\x00\x01\x00\x00\x00\x01\x08\x06\x00\x00\x00\x1f\x15\xc4\x89\x00\x00\x00\nIDATx\x9cc\x00\x01\x00\x00\x05\x00\x01\r\n-\xb4\x00\x00\x00\x00IEND\xaeB`\x82")
	os.WriteFile(filepath.Join(tmpDir, "CODE001.png"), pngData, 0644)
	os.WriteFile(filepath.Join(tmpDir, "CODE003.png"), pngData, 0644)

	p := NewProcessor(tmpExcel, tmpDir, "A", "F", "Sheet1", 2, 105, 20)

	err := p.Run(context.Background())
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	// Verify output file exists (search for timestamped file)
	matches, err := filepath.Glob(filepath.Join(filepath.Dir(tmpExcel), "test_output_*.xlsx"))
	if err != nil || len(matches) == 0 {
		t.Fatal("Output file not created")
	}

	// Verify productMap
	if p.productMap["CODE001"] != 1 {
		t.Errorf("Expected CODE001 row 1, got %d", p.productMap["CODE001"])
	}
	if p.productMap["CODE003"] != 4 {
		t.Errorf("Expected CODE003 row 4, got %d", p.productMap["CODE003"])
	}
}
