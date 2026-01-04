package gui

import (
	"context"
	"fmt"
	"imagetoexcel/internal/engine"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/xuri/excelize/v2"
)

type App struct {
	window fyne.Window

	excelPath  *widget.Entry
	imageDir   *widget.Entry
	codeCol    *widget.Entry
	imageCol   *widget.Entry
	rowHeight  *widget.Entry
	colWidth   *widget.Entry
	workerCnt  *widget.Entry
	sheetEntry *widget.Select
	progress   *widget.ProgressBar
	status     *widget.Label
	runBtn     *widget.Button
}

func NewApp() *App {
	a := app.New()
	w := a.NewWindow("Image to Excel Tool")
	w.Resize(fyne.NewSize(500, 400))

	return &App{
		window: w,
	}
}

func (a *App) Run() {
	a.excelPath = widget.NewEntry()
	a.excelPath.SetPlaceHolder("Path to Excel file...")

	excelBtn := widget.NewButton("Browse Excel", func() {
		d := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if reader != nil {
				path := reader.URI().Path()
				a.excelPath.SetText(path)

				// Automatically fetch sheets
				f, err := excelize.OpenFile(path)
				if err == nil {
					sheets := f.GetSheetList()
					a.sheetEntry.Options = sheets
					if len(sheets) > 0 {
						a.sheetEntry.SetSelected(sheets[0])
					}
					f.Close()
				}
			}
		}, a.window)
		d.SetFilter(storage.NewExtensionFileFilter([]string{".xlsx"}))
		d.Show()
	})

	a.imageDir = widget.NewEntry()
	a.imageDir.SetPlaceHolder("Path to image folder...")

	imageBtn := widget.NewButton("Browse Folder", func() {
		d := dialog.NewFolderOpen(func(list fyne.ListableURI, err error) {
			if list != nil {
				a.imageDir.SetText(list.Path())
			}
		}, a.window)
		d.Show()
	})

	a.codeCol = widget.NewEntry()
	a.codeCol.SetText("A")
	a.imageCol = widget.NewEntry()
	a.imageCol.SetText("F")

	a.rowHeight = widget.NewEntry()
	a.rowHeight.SetText("105")
	a.colWidth = widget.NewEntry()
	a.colWidth.SetText("20")
	a.workerCnt = widget.NewEntry()
	a.workerCnt.SetText("10")

	a.sheetEntry = widget.NewSelect([]string{}, nil)
	a.sheetEntry.PlaceHolder = "Select Sheet..."

	a.progress = widget.NewProgressBar()
	a.status = widget.NewLabel("Ready")

	a.runBtn = widget.NewButton("Start Processing", a.handleRun)

	form := container.NewVBox(
		widget.NewLabel("Input Configuration"),
		container.NewBorder(nil, nil, nil, excelBtn, a.excelPath),
		container.NewBorder(nil, nil, nil, imageBtn, a.imageDir),
		container.NewGridWithColumns(3,
			container.NewVBox(widget.NewLabel("Product Code Col:"), a.codeCol),
			container.NewVBox(widget.NewLabel("Image Target Col:"), a.imageCol),
			container.NewVBox(widget.NewLabel("Sheet Name:"), a.sheetEntry),
		),
		container.NewGridWithColumns(3,
			container.NewVBox(widget.NewLabel("Row Height:"), a.rowHeight),
			container.NewVBox(widget.NewLabel("Col Width:"), a.colWidth),
			container.NewVBox(widget.NewLabel("Worker Threads:"), a.workerCnt),
		),
		widget.NewSeparator(),
		a.status,
		a.progress,
		a.runBtn,
	)

	a.window.SetContent(container.NewPadded(form))
	a.window.ShowAndRun()
}

func (a *App) handleRun() {
	if a.excelPath.Text == "" || a.imageDir.Text == "" {
		dialog.ShowError(fmt.Errorf("please select both Excel file and image folder"), a.window)
		return
	}

	a.runBtn.Disable()
	a.status.SetText("Processing...")
	a.progress.SetValue(0)

	go func() {
		var rH, cW float64
		var wC int
		fmt.Sscanf(a.rowHeight.Text, "%f", &rH)
		fmt.Sscanf(a.colWidth.Text, "%f", &cW)
		fmt.Sscanf(a.workerCnt.Text, "%d", &wC)

		if rH <= 0 {
			rH = 105
		}
		if cW <= 0 {
			cW = 20
		}
		if wC <= 0 {
			wC = 10
		}

		p := engine.NewProcessor(a.excelPath.Text, a.imageDir.Text, a.codeCol.Text, a.imageCol.Text, a.sheetEntry.Selected, wC, rH, cW)

		progressChan := make(chan float64)
		p.SetProgressChan(progressChan)

		go func() {
			for val := range progressChan {
				a.progress.SetValue(val)
			}
		}()

		err := p.Run(context.Background())

		a.window.Canvas().Refresh(a.runBtn)
		if err != nil {
			a.status.SetText("Error: " + err.Error())
			dialog.ShowError(err, a.window)
		} else {
			msg := fmt.Sprintf("Processing completed! Output saved.\n- Missing images: %d (log saved)", len(p.MissingCodes))
			a.status.SetText("Done! Check missing.log if needed.")
			dialog.ShowInformation("Success", msg, a.window)
		}
		a.runBtn.Enable()
	}()
}
