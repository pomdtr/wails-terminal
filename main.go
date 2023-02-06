package main

import (
	"embed"
	"fmt"
	"os"
	"os/exec"

	"github.com/creack/pty"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	var err error

	// Create an instance of the app structure
	app := NewApp()
	cmd := exec.Command("fish")
	pty, err := pty.Start(cmd)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	app.tty = pty

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "sunbeam",
		Width:  750,
		Height: 475,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		AlwaysOnTop:   true,
		DisableResize: true,
		Frameless:     true,
		Mac: &mac.Options{
			WebviewIsTransparent: true,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 100},
		OnStartup:        app.startup,
		OnDomReady:       app.domReady,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
