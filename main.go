package main

import (
	"embed"
	"flag"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	var err error

	lightTheme := flag.String("light-theme", "tomorrow", "Theme to use")
	darkTheme := flag.String("dark-theme", "tomorrow-night", "Theme to use")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		shell := os.Getenv("SHELL")
		args = []string{shell, "-li"}
	}

	// Create an instance of the app structure
	app := NewApp(args)
	app.darkTheme = *darkTheme
	app.lightTheme = *lightTheme

	// Create application with options
	err = wails.Run(&options.App{
		Title:       "Wails Terminal",
		Width:       750,
		StartHidden: true,
		Height:      475,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		AlwaysOnTop:   true,
		DisableResize: true,
		Frameless:     true,
		Mac: &mac.Options{
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
		},
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
