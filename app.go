package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/creack/pty"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context

	options TerminalOptions

	tty  *os.File
	rows uint16
	cols uint16
}

type TerminalOptions struct {
	args       []string
	lightTheme *Theme
	darkTheme  *Theme
}

// NewApp creates a new App application struct
func NewApp(options TerminalOptions) *App {
	return &App{options: options}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) StartTTY() error {
	var cmd *exec.Cmd
	switch len(a.options.args) {
	case 0:
		return fmt.Errorf("no command specified")
	case 1:
		cmd = exec.Command(a.options.args[0])
	default:
		cmd = exec.Command(a.options.args[0], a.options.args[1:]...)
	}
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, "TERM=xterm-256color")

	tty, err := pty.Start(cmd)
	if err != nil {
		return fmt.Errorf("failed to start pty: %w", err)
	}
	if a.rows != 0 && a.cols != 0 {
		pty.Setsize(tty, &pty.Winsize{Rows: a.rows, Cols: a.cols})
	}

	a.tty = tty
	return nil
}

func (a *App) Start() {
	a.StartTTY()
	go func() {
		for {
			buf := make([]byte, 20480)
			n, err := a.tty.Read(buf)
			if err != nil {
				if !errors.Is(err, io.EOF) {
					runtime.LogErrorf(a.ctx, "Read error: %s", err)
					continue
				}

				runtime.Quit(a.ctx)
				continue
			}
			runtime.EventsEmit(a.ctx, "tty-data", buf[:n])
		}
	}()
}

func (a *App) GetDarkTheme() *Theme {
	return a.options.darkTheme
}

func (a *App) GetLightTheme() *Theme {
	return a.options.lightTheme
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
