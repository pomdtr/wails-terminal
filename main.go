package main

import (
	"embed"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	app := NewApp()

	server := NewServer()
	go server.ListenAndServe()

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "wails-xterm",
		Width:  750,
		Height: 475,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		Frameless: true,
		// AlwaysOnTop:      true,
		DisableResize:    true,
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 50},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
