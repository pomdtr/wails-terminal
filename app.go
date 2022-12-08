package main

import (
	"context"
	"encoding/base64"
	"os"
	"os/exec"

	"github.com/creack/pty"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx  context.Context
	tty  *os.File
	rows uint16
	cols uint16
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) domReady(ctx context.Context) {
	go func() {
		for {
			buf := make([]byte, 10240)
			n, err := a.tty.Read(buf)
			if err != nil {
				runtime.LogWarningf(a.ctx, "Read error: %s", err)
				cmd := exec.Command("fish")
				tty, err := pty.Start(cmd)
				if err != nil {
					runtime.LogErrorf(a.ctx, "Read error: %s", err)
					break
				}
				pty.Setsize(tty, &pty.Winsize{Rows: a.rows, Cols: a.cols})
				a.tty = tty
				runtime.WindowMinimise(ctx)
				runtime.EventsEmit(ctx, "clearTerminal")
				continue
			}
			runtime.EventsEmit(ctx, "ttyData", base64.StdEncoding.EncodeToString(buf[:n]))
		}
	}()
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
