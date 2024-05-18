package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/runletapp/go-console"
	logs "github.com/tea4go/gh/log4go"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type consoleService struct {
	ctx context.Context

	termCmd string
	pty     console.Console

	ptyWrite io.WriteCloser
	ptyRead  io.Reader

	rows int
	cols int

	theme *Theme

	fontName   string
	fontWeight bool
	fontSize   int
}

// 获取Console实例
func NewConsole(term_cmd, font_name string, font_weight bool, font_size int, theme *Theme) *consoleService {
	if len(term_cmd) == 0 {
		term_cmd := os.Getenv("SHELL")
		if term_cmd == "" {
			term_cmd = "cmd"
		}
	}
	return &consoleService{
		termCmd:    term_cmd,
		fontName:   font_name,
		fontSize:   font_size,
		fontWeight: font_weight,
		theme:      theme,
	}
}

func (a *consoleService) Startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *consoleService) startTTY() error {
	logs.Debug("StartTTY() ......")

	var err error
	a.pty, err = console.New(80, 24)
	if err != nil {
		return fmt.Errorf("初始化本地终端失败，%s", err.Error())
	}
	a.ptyRead = a.pty
	a.ptyWrite = a.pty

	err = a.pty.Start([]string{a.termCmd})
	if err != nil {
		return fmt.Errorf("启动本地终端失败，%s", err.Error())
	}

	return nil
}

// 开始循环读取终端输出信息
func (a *consoleService) LoopRead() {
	err := a.startTTY()
	if err != nil {
		logs.Error(err)
		return
	}

	//从console读取输出信息，通过事件发送到前端
	go func() {
		for {
			buf := make([]byte, 20480)
			n, err := a.ptyRead.Read(buf)
			if err != nil {
				if !errors.Is(err, io.EOF) {
					logs.Error("读取Console输出错误，%s", err.Error())
					time.Sleep(1 * time.Second)
					continue
				}

				//退出应用程序
				runtime.Quit(a.ctx)
				continue
			}
			runtime.EventsEmit(a.ctx, "tty-data", buf[:n])
		}
	}()
}

func (a *consoleService) Close() error {
	var err error
	if a.pty != nil {
		err = a.pty.Close()
		a.pty = nil
	}
	return err
}

func (a *consoleService) Resize(rows, cols int) {
	if a.pty != nil && rows > 0 && cols >= 10 {
		logs.Debug("SetTTYSize: %d, %d", rows, cols)
		a.rows = rows
		a.cols = cols
		a.pty.SetSize(cols, rows)
	}
}

func (a *consoleService) SendText(text string) {
	a.ptyWrite.Write([]byte(text))
}

func (a *consoleService) GetTheme() *Theme {
	return a.theme
}

func (a *consoleService) GetFontName() string {
	if len(a.fontName) > 0 {
		a.fontName = a.fontName + ","
	}
	logs.Notice("FontName : %s", a.fontName)
	return a.fontName
}

func (a *consoleService) GetFontSize() int {
	if a.fontSize > 28 || a.fontSize < 9 {
		a.fontSize = 18
	}
	logs.Notice("FontSize : %d", a.fontSize)
	return a.fontSize
}

func (a *consoleService) GetFontWeight() bool {
	logs.Notice("FontWeight : %s", a.fontWeight)
	return a.fontWeight
}
