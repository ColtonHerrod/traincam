package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend
var assets embed.FS

func main() {
	// Create an instance of our App struct from app.go
	app := NewApp()

	// Create the application
	err := wails.Run(&options.App{
		Title:       "Train Cam Tracker",
		Width:       1200,
		Height:      800,
		AssetServer: &assetserver.Options{Assets: assets},
		OnStartup:   app.startup,
		Bind: []interface{}{
			app, // This is the magic line that makes GetCameras() available in JS
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
