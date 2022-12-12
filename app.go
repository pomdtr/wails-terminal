package main

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.design/x/hotkey"
)

// App struct
type App struct {
	ctx             context.Context
	windowIsVisible bool
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		windowIsVisible: true,
	}
}

func (a *App) startup(ctx context.Context) {
	server := NewServer(ctx)
	go server.ListenAndServe()
	go a.watchHotkey(ctx)
	a.ctx = ctx
}

func (a *App) domReady(ctx context.Context) {
}

func (a *App) HideWindow() {
	runtime.WindowHide(a.ctx)
	a.windowIsVisible = false
}

func (a *App) watchHotkey(ctx context.Context) {
	hk := hotkey.New([]hotkey.Modifier{hotkey.ModCmd, hotkey.ModShift}, hotkey.KeySpace)
	err := hk.Register()
	if err != nil {
		runtime.LogWarningf(ctx, "Error registering hotkey: %s", err.Error())
	}

	for {
		<-hk.Keydown()
		if a.windowIsVisible {
			runtime.WindowHide(ctx)
		} else {
			runtime.WindowShow(ctx)
		}
		a.windowIsVisible = !a.windowIsVisible
	}
}
