package engine

import (
	"context"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	_ "golang.org/x/image/webp"

	"github.com/xuri/excelize/v2"
)

type Job struct {
	ProductCode string
	ImagePath   string
	RowIndex    int
}

type Result struct {
	Job      Job
	ImgBytes []byte
	Width    int
	Height   int
	Err      error
}

type Processor struct {
	ExcelPath   string
	ImageDir    string
	CodeCol     string
	ImageCol    string
	SheetName   string
	WorkerCount int
	RowHeight   float64
	ColWidth    float64

	f              *excelize.File
	productMap     map[string]int
	jobs           chan Job
	results        chan Result
	progressChan   chan float64
	missingMu      sync.Mutex // Protects MissingCodes
	MissingCodes   []string
	ProcessedCount int // Number of successfully processed images
}

func NewProcessor(excelPath, imageDir, codeCol, imageCol, sheetName string, workerCount int, rowHeight, colWidth float64) *Processor {
	return &Processor{
		ExcelPath:    excelPath,
		ImageDir:     imageDir,
		CodeCol:      codeCol,
		ImageCol:     imageCol,
		SheetName:    sheetName,
		WorkerCount:  workerCount,
		RowHeight:    rowHeight,
		ColWidth:     colWidth,
		productMap:   make(map[string]int),
		jobs:         make(chan Job, 100),
		results:      make(chan Result, 100),
		MissingCodes: make([]string, 0, 100),
	}
}

func (p *Processor) SetProgressChan(ch chan float64) {
	p.progressChan = ch
}

func (p *Processor) Run(ctx context.Context) error {
	var err error
	p.f, err = excelize.OpenFile(p.ExcelPath)
	if err != nil {
		return fmt.Errorf("failed to open excel: %w", err)
	}
	defer p.f.Close()

	if p.SheetName == "" {
		p.SheetName = p.f.GetSheetName(0)
	}

	// 1. Mapping: Read all product codes using iterator
	rows, err := p.f.Rows(p.SheetName)
	if err != nil {
		return fmt.Errorf("failed to get rows: %w", err)
	}
	defer rows.Close()

	codeColIdx, err := excelize.ColumnNameToNumber(p.CodeCol)
	if err != nil {
		return fmt.Errorf("invalid code column: %w", err)
	}
	codeColIdx-- // 0-indexed

	rowIdx := 0
	for rows.Next() {
		// Check for cancellation during row processing
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		rowIdx++
		row, err := rows.Columns()
		if err != nil {
			log.Printf("Warning: failed to read columns for row %d: %v", rowIdx, err)
			continue
		}
		if len(row) > codeColIdx {
			code := strings.TrimSpace(row[codeColIdx])
			if code != "" {
				p.productMap[code] = rowIdx
			}
		}
	}

	if err := rows.Error(); err != nil {
		return fmt.Errorf("error reading rows: %w", err)
	}

	// 2. Start Workers for Image Loading/Scaling
	var wg sync.WaitGroup
	for i := 0; i < p.WorkerCount; i++ {
		wg.Add(1)
		go p.worker(ctx, &wg)
	}

	// 3. Dispatcher: Scan images and send jobs
	// 3. Dispatcher: Scan images and send jobs
	go func() {
		defer close(p.jobs)

		files, err := os.ReadDir(p.ImageDir)
		if err != nil {
			log.Printf("Error reading image directory: %v", err)
			return
		}

		// Map files for quick lookup
		availableImages := make(map[string]string)
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			ext := strings.ToLower(filepath.Ext(file.Name()))
			if ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".webp" {
				name := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
				availableImages[name] = file.Name()
			}
		}

		// Dispatch jobs and track missing
		for code, rowIndex := range p.productMap {
			// Check for cancellation
			select {
			case <-ctx.Done():
				return
			default:
			}

			if fileName, ok := availableImages[code]; ok {
				select {
				case p.jobs <- Job{
					ProductCode: code,
					ImagePath:   filepath.Join(p.ImageDir, fileName),
					RowIndex:    rowIndex,
				}:
				case <-ctx.Done():
					return
				}
			} else {
				p.missingMu.Lock()
				p.MissingCodes = append(p.MissingCodes, code)
				p.missingMu.Unlock()
			}
		}
	}()

	// Wait for workers and close results
	go func() {
		wg.Wait()
		close(p.results)
	}()

	// 4. Main Loop: Receive results and modify Excel
	p.ProcessedCount = 0

	// We'll update progress based on results received
resultLoop:
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case res, ok := <-p.results:
			if !ok {
				break resultLoop // Results channel closed, finishing up
			}

			if res.Err != nil {
				log.Printf("Error processing %s: %v", res.Job.ProductCode, res.Err)
				continue
			}

			if err := p.insertImageToExcel(res); err != nil {
				log.Printf("Error inserting %s: %v", res.Job.ProductCode, err)
				continue
			}

			p.ProcessedCount++
			if p.progressChan != nil {
				// Non-blocking send to progress channel to prevent stalling if frontend is slow
				select {
				case p.progressChan <- float64(p.ProcessedCount) / float64(len(p.productMap)):
				default:
				}
			}
		}
	}

	// 5. Save result
	timestamp := time.Now().Format("20060102_150405")
	outputName := fmt.Sprintf("%s_output_%s.xlsx", strings.TrimSuffix(p.ExcelPath, filepath.Ext(p.ExcelPath)), timestamp)
	if err := p.f.SaveAs(outputName); err != nil {
		return fmt.Errorf("failed to save excel: %w", err)
	}

	// 6. Write log for missing codes
	if len(p.MissingCodes) > 0 {
		logPath := fmt.Sprintf("%s_missing_%s.log", strings.TrimSuffix(p.ExcelPath, filepath.Ext(p.ExcelPath)), timestamp)
		// We ignore error here as it's secondary
		_ = os.WriteFile(logPath, []byte(strings.Join(p.MissingCodes, "\n")), 0644)
	}

	return nil
}

func (p *Processor) worker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-p.jobs:
			if !ok {
				return
			}
			imgBytes, w, h, err := p.loadImageData(job.ImagePath)
			select {
			case p.results <- Result{
				Job:      job,
				ImgBytes: imgBytes,
				Width:    w,
				Height:   h,
				Err:      err,
			}:
			case <-ctx.Done():
				return
			}
		}
	}
}

func (p *Processor) loadImageData(path string) ([]byte, int, int, error) {
	imgFile, err := os.Open(path)
	if err != nil {
		return nil, 0, 0, err
	}
	defer imgFile.Close()

	imgConfig, _, err := image.DecodeConfig(imgFile)
	if err != nil {
		return nil, 0, 0, err
	}

	_, _ = imgFile.Seek(0, 0)
	imgBytes, err := os.ReadFile(path)
	if err != nil {
		return nil, 0, 0, err
	}

	return imgBytes, imgConfig.Width, imgConfig.Height, nil
}

func (p *Processor) insertImageToExcel(res Result) error {
	colIdx, err := excelize.ColumnNameToNumber(p.ImageCol)
	if err != nil {
		return fmt.Errorf("invalid image column '%s': %w", p.ImageCol, err)
	}
	cellName, err := excelize.CoordinatesToCellName(colIdx, res.Job.RowIndex)
	if err != nil {
		return fmt.Errorf("failed to get cell name: %w", err)
	}

	// Set Row Height and Col Width from Processor settings
	if err = p.f.SetRowHeight(p.SheetName, res.Job.RowIndex, p.RowHeight); err != nil {
		return fmt.Errorf("failed to set row height: %w", err)
	}
	if err = p.f.SetColWidth(p.SheetName, p.ImageCol, p.ImageCol, p.ColWidth); err != nil {
		return fmt.Errorf("failed to set col width: %w", err)
	}

	// Calculate scale to fit. Default rowHeight 105pt is ~140px. colWidth 20 is ~145px
	// We'll use a margin of 10px
	targetH := p.RowHeight*1.333 - 10
	targetW := p.ColWidth*7.0 - 10 // Approximation for column width in pixels

	scaleX := targetW / float64(res.Width)
	scaleY := targetH / float64(res.Height)

	scale := scaleX
	if scaleY < scale {
		scale = scaleY
	}

	err = p.f.AddPictureFromBytes(p.SheetName, cellName, &excelize.Picture{
		Extension: filepath.Ext(res.Job.ImagePath),
		File:      res.ImgBytes,
		Format: &excelize.GraphicOptions{
			ScaleX:      scale,
			ScaleY:      scale,
			OffsetX:     5,
			OffsetY:     5,
			Positioning: "oneCell",
		},
	})

	return err
}
