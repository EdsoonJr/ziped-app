package main

import (
	"context"
	"embed"
	"ziped-app/backend/services"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()
	archiveService := services.NewArchiveService()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "ziped-app",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		DragAndDrop: &options.DragAndDrop{
			EnableFileDrop:     true,
			DisableWebViewDrop: true,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			archiveService.Startup(ctx)
		},
		Bind: []interface{}{
			app,
			archiveService,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
