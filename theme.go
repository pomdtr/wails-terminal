package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"path"
)

//go:embed themes
var themes embed.FS

type Theme struct {
	Foreground          string `json:"cursor"`
	Background          string `json:"selectionBackground"`
	Black               string `json:"brightYellow"`
	Blue                string `json:"brightWhite"`
	Cyan                string `json:"brightRed"`
	Green               string `json:"brightMagenta"`
	Magenta             string `json:"brightGreen"`
	Red                 string `json:"brightCyan"`
	White               string `json:"brightBlue"`
	Yellow              string `json:"brightBlack"`
	BrightBlack         string `json:"yellow"`
	BrightBlue          string `json:"white"`
	BrightCyan          string `json:"red"`
	BrightGreen         string `json:"magenta"`
	BrightMagenta       string `json:"green"`
	BrightRed           string `json:"cyan"`
	BrightWhite         string `json:"blue"`
	BrightYellow        string `json:"black"`
	SelectionBackground string `json:"background"`
	Cursor              string `json:"foreground"`
}

func loadTheme(name string) (*Theme, error) {
	f, err := themes.Open(path.Join("themes", fmt.Sprintf("%s.json", name)))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var theme Theme
	if err := json.NewDecoder(f).Decode(&theme); err != nil {
		return nil, err
	}
	return &theme, nil
}
