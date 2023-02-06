package main

import (
	"context"
	"errors"
	"io"
	"os"
	"os/exec"

	"github.com/creack/pty"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.design/x/hotkey"
)

// App struct
type App struct {
	ctx          context.Context
	tty          *os.File
	windowHidden bool
	rows         uint16
	cols         uint16
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	go a.watchHotkey(ctx)
}

func (a *App) domReady(ctx context.Context) {
	go func() {
		for {
			buf := make([]byte, 20480)
			n, err := a.tty.Read(buf)
			if err != nil {
				if !errors.Is(err, io.EOF) {
					runtime.LogErrorf(a.ctx, "Read error: %s", err)
					panic(err)
				}

				a.HideWindow()
				cmd := exec.Command("fish")
				tty, err := pty.Start(cmd)
				if err != nil {
					runtime.LogErrorf(a.ctx, "Read error: %s", err)
					break
				}
				pty.Setsize(tty, &pty.Winsize{Rows: a.rows, Cols: a.cols})
				a.tty = tty
				runtime.EventsEmit(ctx, "clear-terminal")
				continue
			}
			runtime.EventsEmit(ctx, "ttyData", buf[:n])
		}
	}()
}

func (a *App) HideWindow() {
	runtime.WindowHide(a.ctx)
	a.windowHidden = true
}

func (a *App) SetTTYSize(rows, cols uint16) {
	a.rows = rows
	a.cols = cols
	runtime.LogDebugf(a.ctx, "SetTTYSize: %d, %d", rows, cols)
	pty.Setsize(a.tty, &pty.Winsize{Rows: rows, Cols: cols})
}

func (a *App) SendText(text string) {
	a.tty.Write([]byte(text))
}

func (a *App) watchHotkey(ctx context.Context) {
	hk := hotkey.New([]hotkey.Modifier{hotkey.ModCmd, hotkey.ModShift}, hotkey.KeySpace)
	err := hk.Register()
	if err != nil {
		runtime.LogWarningf(ctx, "Error registering hotkey: %s", err.Error())
	}

	for {
		<-hk.Keydown()
		if a.windowHidden {
			runtime.WindowShow(ctx)
		} else {
			runtime.WindowHide(ctx)
		}
		a.windowHidden = !a.windowHidden
	}
}
