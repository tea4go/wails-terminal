package main

import (
	"embed"
	"flag"
	"os"

	logs "github.com/tea4go/gh/log4go"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	//logs.SetLogger("conn", fmt.Sprintf(`{"addr":"%s:%d","level":7,"name":"Mate60"}`, "192.168.50.1", 9514))
	//logs.SetLogger("file", `{"filename":"ulog_whaleterm.log", "perm": "0666","level":7}`)
	logs.SetLogger("console", `{"color":true,"level":7}`)
	logs.SetLevel(logs.LevelDebug)
	logs.Debug("=================================")

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
		if shell == "" {
			args = []string{"cmd"}
		} else {
			args = []string{shell}
		}
	}
	logs.Notice(args)
	// Create an instance of the app structure
	app := NewConsole(args)

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "Wails Terminal",
		Width:  750,
		Height: 475,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 0, G: 0, B: 0, A: 1},
		OnStartup:        app.Startup,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
