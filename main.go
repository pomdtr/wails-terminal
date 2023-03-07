package main

import (
	"embed"
	"flag"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	var err error
	flags := struct {
		lightTheme string
		darkTheme  string
	}{}

	flag.StringVar(&flags.lightTheme, "light-theme", "tomorrow", "Theme to use")
	flag.StringVar(&flags.darkTheme, "dark-theme", "tomorrow-night", "Theme to use")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		shell := os.Getenv("SHELL")
		args = []string{shell, "-li"}
	}

	lightTheme, err := loadTheme(flags.lightTheme)
	if err != nil {
		println("Error:", err.Error())
		os.Exit(1)
	}

	darkTheme, err := loadTheme(flags.darkTheme)
	if err != nil {
		println("Error:", err.Error())
		os.Exit(1)
	}

	// Create an instance of the app structure
	app := NewApp(TerminalOptions{
		args:       args,
		lightTheme: lightTheme,
		darkTheme:  darkTheme,
	})

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "Wails Terminal",
		Width:  750,
		Height: 475,
		AssetServer: &assetserver.Options{
			Assets: assets,
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
