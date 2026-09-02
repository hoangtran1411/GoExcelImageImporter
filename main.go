package main

import (
	"embed"
	"log"

	"github.com/wailsapp/wails/v3/pkg/application"
)

//go:embed all:frontend/dist
var assets embed.FS

func init() {
	// Register custom events for frontend bindings
	application.RegisterEvent[float64]("progress")
	application.RegisterEvent[string]("updateProgress")
}

func main() {
	appService := NewApp()

	app := application.New(application.Options{
		Name:        "Image to Excel Importer",
		Description: "Batch import images into Excel sheets",
		Services: []application.Service{
			application.NewService(appService),
		},
		Assets: application.AssetOptions{
			Handler: application.AssetFileServerFS(assets),
		},
		SingleInstance: &application.SingleInstanceOptions{
			UniqueID: "com.hoangtran.goexcelimageimporter",
			OnSecondInstanceLaunch: func(data application.SecondInstanceData) {
				if win, ok := application.Get().Window.GetByName("main"); ok {
					win.Focus()
					win.Restore()
				}
			},
		},
	})

	appService.setApp(app)

	app.Window.NewWithOptions(application.WebviewWindowOptions{
		Name:             "main",
		Title:            "Image to Excel Importer",
		Width:            900,
		Height:           835,
		DisableResize:    true,
		MinWidth:         700,
		MinHeight:        600,
		BackgroundColour: application.NewRGB(27, 38, 54),
		URL:              "/",
	})

	err := app.Run()
	if err != nil {
		log.Fatal(err)
	}
}
