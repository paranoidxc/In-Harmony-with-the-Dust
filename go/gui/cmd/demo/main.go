package main

import (
	"flag"
	"fmt"
	"log"
	"sort"
	"time"

	"classicui"
	"classicui/event"
	"classicui/widgets"
)

const (
	cmdAddPath      classicui.CommandID = "cmd.file.add_path"
	cmdExit         classicui.CommandID = "cmd.file.exit"
	cmdSortByName   classicui.CommandID = "cmd.view.sort.name"
	cmdSortByLength classicui.CommandID = "cmd.view.sort.length"
	cmdAbout        classicui.CommandID = "cmd.help.about"
)

func main() {
	autoQuit := flag.Duration("auto-quit", 0, "automatically exit after the given duration")
	flag.Parse()

	app := classicui.NewApp(classicui.Config{
		Title:        "Classic UI Phase 3",
		LogicalSize:  classicui.Size{W: 640, H: 480},
		PresentScale: 2,
		Theme:        classicui.DefaultClassicTheme(),
	})

	win := classicui.NewWindow("phase0", classicui.Rect{
		X: 64,
		Y: 40,
		W: 436,
		H: 350,
	})
	win.SetTitle("Phase 3 - Windows Classic")

	title := widgets.NewLabel("intro", "菜单栏、命令、快捷键和子菜单已经接进来了。", classicui.Rect{
		X: 14,
		Y: 16,
		W: 392,
		H: 18,
	})
	status := widgets.NewLabel("status", "试试 Alt、方向键、Ctrl+N 和 Ctrl+Q。", classicui.Rect{
		X: 14,
		Y: 244,
		W: 392,
		H: 18,
	})
	pathLabel := widgets.NewLabel("pathLabel", "路径：", classicui.Rect{
		X: 14,
		Y: 48,
		W: 50,
		H: 18,
	})
	pathEdit := widgets.NewEdit("path", classicui.Rect{
		X: 56,
		Y: 44,
		W: 350,
		H: 24,
	})
	pathEdit.SetText(`C:\我的文档\新建文件夹`)

	list := widgets.NewListBox("files", classicui.Rect{
		X: 14,
		Y: 80,
		W: 392,
		H: 148,
	})
	items := []string{
		`桌面`,
		`我的电脑`,
		`网络邻居`,
		`Program Files`,
		`Windows`,
		`Temp`,
		`Documents`,
		`Downloads`,
		`Music`,
		`Pictures`,
		`Videos`,
		`Games`,
		`Logs`,
		`Backups`,
		`Old Projects`,
	}
	list.SetItems(items)

	add := widgets.NewButton("add", "添加", classicui.Rect{
		X: 236,
		Y: 274,
		W: 76,
		H: 24,
	})
	closeBtn := widgets.NewButton("cancel", "关闭", classicui.Rect{
		X: 320,
		Y: 274,
		W: 76,
		H: 24,
	})

	sortByNameItem := widgets.NewMenuItem(cmdSortByName, "By &Name", nil)
	sortByNameItem.Checked = true
	sortByLengthItem := widgets.NewMenuItem(cmdSortByLength, "By &Length", nil)
	win.SetMenuBar(widgets.NewMenuBar(
		widgets.NewSubmenuItem("&File", widgets.NewMenu(
			widgets.NewMenuItem(cmdAddPath, "&Add Path", &widgets.Accelerator{
				Key:       event.KeyN,
				Modifiers: event.ModCtrl,
			}),
			widgets.NewSeparator(),
			widgets.NewMenuItem(cmdExit, "E&xit", &widgets.Accelerator{
				Key:       event.KeyQ,
				Modifiers: event.ModCtrl,
			}),
		)),
		widgets.NewSubmenuItem("&View", widgets.NewMenu(
			widgets.NewSubmenuItem("&Sort", widgets.NewMenu(
				sortByNameItem,
				sortByLengthItem,
			)),
		)),
		widgets.NewSubmenuItem("&Help", widgets.NewMenu(
			widgets.NewMenuItem(cmdAbout, "&About", nil),
		)),
	))

	currentSort := cmdSortByName
	applySort := func(cmd classicui.CommandID) {
		currentSort = cmd
		sortByNameItem.Checked = cmd == cmdSortByName
		sortByLengthItem.Checked = cmd == cmdSortByLength
		switch cmd {
		case cmdSortByLength:
			sort.SliceStable(items, func(i, j int) bool {
				left := len([]rune(items[i]))
				right := len([]rune(items[j]))
				if left == right {
					return items[i] < items[j]
				}
				return left < right
			})
			status.SetText("已切换为按名称长度排序。")
		default:
			sort.SliceStable(items, func(i, j int) bool {
				return items[i] < items[j]
			})
			status.SetText("已切换为按名称排序。")
		}
		list.SetItems(items)
		app.Desktop().InvalidateRect(win.Bounds())
	}

	runCommand := func(cmd classicui.CommandID) {
		switch cmd {
		case cmdAddPath:
			next := pathEdit.Text()
			if next == "" {
				status.SetText("输入框不能为空。")
				app.Desktop().InvalidateRect(win.Bounds())
				return
			}
			items = append(items, next)
			applySort(currentSort)
			status.SetText(fmt.Sprintf("已添加：%s", next))
			win.SetTitle(fmt.Sprintf("Phase 3 - 共 %d 项", len(items)))
		case cmdExit:
			app.Quit()
		case cmdSortByName, cmdSortByLength:
			applySort(cmd)
		case cmdAbout:
			status.SetText("Phase 3: MenuBar、PopupMenu、快捷键和命令路由已就位。")
		}
		app.Desktop().InvalidateRect(win.Bounds())
	}
	app.OnCommand(runCommand)

	list.OnChange(func(index int, value string) {
		status.SetText(fmt.Sprintf("当前选中：%s", value))
		app.Desktop().InvalidateRect(win.Bounds())
	})

	add.OnClick(func() {
		runCommand(cmdAddPath)
	})
	closeBtn.OnClick(func() {
		runCommand(cmdExit)
	})

	win.Content().Add(title)
	win.Content().Add(pathLabel)
	win.Content().Add(pathEdit)
	win.Content().Add(list)
	win.Content().Add(status)
	win.Content().Add(add)
	win.Content().Add(closeBtn)
	win.SetDefaultButton(add)
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
