package main

import (
	"embed"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/adrg/sysfont"
	logs "github.com/tea4go/gh/log4go"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

//go:embed all:frontend/dist
var assets embed.FS

type TFontItem struct {
	Name   string `json:"name"`
	Family string `json:"family"`
	Path   string `json:"path"`
}

// 获取字体列表
func ShowFontList() {
	finder := sysfont.NewFinder(nil)
	var fontList []TFontItem
	fontSet := NewSet[string]() //字体排重
	for _, font := range finder.List() {
		if len(font.Family) == 0 && strings.Contains(font.Filename, "MicrosoftYaHeiMono") {
			font.Family = "Microsoft YaHei Mono"
			font.Name = "Microsoft YaHei Mono"
		}
		if len(font.Family) > 0 && !strings.HasPrefix(font.Family, ".") && fontSet.Add(font.Family) {
			fontList = append(fontList, TFontItem{
				Name:   font.Name,
				Family: font.Family,
				Path:   font.Filename,
			})
			//fmt.Println(font.Name, "(", font.Family, ")", font.Filename)
		}
	}
	sort.Slice(fontList, func(i, j int) bool {
		return fontList[i].Name < fontList[j].Name
	})

	for _, v := range fontList {
		if strings.Contains(v.Name, "Ubuntu") {
			fmt.Printf("%-40s %-30s %s\n", v.Name, v.Family, v.Path)
		}
	}
}

func main() {
	logs.SetLogger("console", `{"color":true,"level":7}`)
	logs.SetLevel(logs.LevelDebug)
	logs.Debug("=================================")

	var err error
	var cmd_name string
	var font_name string
	var font_weight bool
	var font_size int
	var theme_name string

	flag.StringVar(&theme_name, "theme", "tomorrow", "主题名称")
	flag.StringVar(&font_name, "font_name", "", "字体名称")
	flag.BoolVar(&font_weight, "font_weight", false, "是否粗体")
	flag.IntVar(&font_size, "font_size", 18, "字体大小")
	flag.Parse()

	cmd_name = "/bin/bash"
	args := flag.Args()
	if len(args) > 0 {
		cmd_name = args[0]
	}

	theme, err := loadTheme(theme_name)
	if err != nil {
		println("Error:", err.Error())
		os.Exit(1)
	}

	ShowFontList()

	// Create an instance of the app structure
	app := NewConsole(cmd_name, font_name, font_weight, font_size, theme)

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "Wails Terminal",
		Width:  1440,
		Height: 800,
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
