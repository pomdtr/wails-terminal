package main

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/creack/pty"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.design/x/hotkey"
)

//go:embed themes
var themes embed.FS

// App struct
type App struct {
	ctx context.Context

	tty *os.File

	theme        string
	args         []string
	windowHidden bool

	rows uint16
	cols uint16
}

// NewApp creates a new App application struct
func NewApp(theme string, args []string) *App {
	return &App{theme: theme, args: args, windowHidden: true}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.StartTTY()
	go a.watchHotkey(ctx)
}

func (a *App) StartTTY() error {
	var cmd *exec.Cmd
	switch len(a.args) {
	case 0:
		return fmt.Errorf("no command specified")
	case 1:
		cmd = exec.Command(a.args[0])
	default:
		cmd = exec.Command(a.args[0], a.args[1:]...)
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
	a.ShowWindow()
	go func() {
		for {
			buf := make([]byte, 20480)
			n, err := a.tty.Read(buf)
			if err != nil {
				if !errors.Is(err, io.EOF) {
					runtime.LogErrorf(a.ctx, "Read error: %s", err)
				}

				// Restart the TTY
				a.HideWindow()
				runtime.EventsEmit(a.ctx, "clear-terminal")
				a.StartTTY()
				continue
			}
			runtime.EventsEmit(a.ctx, "tty-data", buf[:n])
		}
	}()
}

func (a *App) GetTheme() map[string]string {
	themePath := fmt.Sprintf("themes/%s.json", a.theme)

	bytes, err := themes.ReadFile(themePath)
	if err != nil {
		runtime.LogWarningf(a.ctx, "Error reading theme: %s", err)
		return nil
	}

	var theme map[string]string
	if err = json.Unmarshal(bytes, &theme); err != nil {
		runtime.LogWarningf(a.ctx, "Error parsing theme: %s", err)
		return nil
	}

	return theme
}

func (a *App) ShowWindow() {
	runtime.WindowShow(a.ctx)
	a.windowHidden = false
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
