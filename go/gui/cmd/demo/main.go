package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"classicui"
	"classicui/widgets"
)

func main() {
	autoQuit := flag.Duration("auto-quit", 0, "automatically exit after the given duration")
	flag.Parse()

	app := classicui.NewApp(classicui.Config{
		Title:        "Classic UI Phase 0",
		LogicalSize:  classicui.Size{W: 640, H: 480},
		PresentScale: 2,
		Theme:        classicui.DefaultClassicTheme(),
	})

	win := classicui.NewWindow("phase0", classicui.Rect{
		X: 72,
		Y: 48,
		W: 360,
		H: 220,
	})
	win.SetTitle("Phase 1 - Windows Classic")

	title := widgets.NewLabel("intro", "按 Tab 切换焦点，Enter 触发默认按钮。", classicui.Rect{
		X: 14,
		Y: 16,
		W: 300,
		H: 18,
	})
	status := widgets.NewLabel("status", "默认按钮还没有被触发。", classicui.Rect{
		X: 14,
		Y: 44,
		W: 300,
		H: 18,
	})
	ok := widgets.NewButton("ok", "确定", classicui.Rect{
		X: 168,
		Y: 150,
		W: 76,
		H: 24,
	})
	cancel := widgets.NewButton("cancel", "关闭", classicui.Rect{
		X: 252,
		Y: 150,
		W: 76,
		H: 24,
	})

	clicks := 0
	ok.OnClick(func() {
		clicks++
		status.SetText(fmt.Sprintf("默认按钮已触发 %d 次。", clicks))
		win.SetTitle(fmt.Sprintf("Phase 1 - 点击 %d", clicks))
		app.Desktop().InvalidateRect(win.Bounds())
	})
	cancel.OnClick(func() {
		app.Quit()
	})

	win.Content().Add(title)
	win.Content().Add(status)
	win.Content().Add(ok)
	win.Content().Add(cancel)
	win.SetDefaultButton(ok)
	app.Desktop().AddWindow(win)

	if *autoQuit > 0 {
		time.AfterFunc(*autoQuit, func() {
			app.Post(func() {
				app.Quit()
			})
		})
	}

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
