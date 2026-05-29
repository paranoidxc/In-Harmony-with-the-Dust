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
		Title:        "Classic UI Phase 4",
		LogicalSize:  classicui.Size{W: 640, H: 480},
		PresentScale: 2,
		Theme:        classicui.DefaultClassicTheme(),
	})

	win := classicui.NewWindow("phase0", classicui.Rect{
		X: 64,
		Y: 40,
		W: 436,
		H: 362,
	})
	win.SetTitle("Phase 4 - Windows Classic")

	toolbarSortByName := widgets.NewToolbarButton(cmdSortByName, "By Name")
	toolbarSortByName.Checked = true
	toolbarSortByLength := widgets.NewToolbarButton(cmdSortByLength, "By Length")
	toolbar := widgets.NewToolbar("toolbar", classicui.Rect{
		X: 14,
		Y: 12,
		W: 392,
		H: 28,
	},
		widgets.NewToolbarButton(cmdAddPath, "Add Path"),
		widgets.NewToolbarSeparator(),
		toolbarSortByName,
		toolbarSortByLength,
		widgets.NewToolbarSeparator(),
		widgets.NewToolbarButton(cmdAbout, "About"),
	)

	title := widgets.NewLabel("intro", "菜单、工具栏和状态栏共用同一套命令。", classicui.Rect{
		X: 14,
		Y: 48,
		W: 392,
		H: 18,
	})
	pathLabel := widgets.NewLabel("pathLabel", "路径：", classicui.Rect{
		X: 10,
		Y: 12,
		W: 50,
		H: 18,
	})
	pathEdit := widgets.NewEdit("path", classicui.Rect{
		X: 52,
		Y: 8,
		W: 316,
		H: 24,
	})
	pathEdit.SetText(`C:\我的文档\新建文件夹`)

	list := widgets.NewListBox("files", classicui.Rect{
		X: 10,
		Y: 42,
		W: 358,
		H: 104,
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
		X: 208,
		Y: 154,
		W: 76,
		H: 24,
	})
	closeBtn := widgets.NewButton("cancel", "关闭", classicui.Rect{
		X: 292,
		Y: 154,
		W: 76,
		H: 24,
	})
	sortNameBtn := widgets.NewButton("sortName", "按名称", classicui.Rect{
		X: 10,
		Y: 44,
		W: 88,
		H: 24,
	})
	sortLengthBtn := widgets.NewButton("sortLength", "按长度", classicui.Rect{
		X: 104,
		Y: 44,
		W: 88,
		H: 24,
	})
	aboutBtn := widgets.NewButton("about", "关于", classicui.Rect{
		X: 10,
		Y: 78,
		W: 88,
		H: 24,
	})
	info1 := widgets.NewLabel("info1", "TabControl 是 Phase 4 的下一步。", classicui.Rect{
		X: 10,
		Y: 12,
		W: 320,
		H: 18,
	})
	info2 := widgets.NewLabel("info2", "切换页签后，隐藏页不会再参与焦点循环。", classicui.Rect{
		X: 10,
		Y: 112,
		W: 340,
		H: 18,
	})
	info3 := widgets.NewLabel("info3", "当前排序仍然复用菜单和工具栏命令。", classicui.Rect{
		X: 10,
		Y: 136,
		W: 340,
		H: 18,
	})
	browsePage := widgets.NewPanel("browsePage", classicui.Rect{})
	browsePage.Add(pathLabel)
	browsePage.Add(pathEdit)
	browsePage.Add(list)
	browsePage.Add(add)
	browsePage.Add(closeBtn)
	commandsPage := widgets.NewPanel("commandsPage", classicui.Rect{})
	commandsPage.Add(info1)
	commandsPage.Add(sortNameBtn)
	commandsPage.Add(sortLengthBtn)
	commandsPage.Add(aboutBtn)
	commandsPage.Add(info2)
	commandsPage.Add(info3)
	tabs := widgets.NewTabControl("tabs", classicui.Rect{
		X: 14,
		Y: 72,
		W: 392,
		H: 204,
	},
		widgets.NewTabPage("文件列表", browsePage),
		widgets.NewTabPage("命令演示", commandsPage),
	)
	statusBar := widgets.NewStatusBar("status", classicui.Rect{
		X: 14,
		Y: 286,
		W: 392,
		H: 22,
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
	sortLabel := func(cmd classicui.CommandID) string {
		switch cmd {
		case cmdSortByLength:
			return "Length"
		default:
			return "Name"
		}
	}
	updateStatus := func(message string) {
		statusBar.SetPanes([]widgets.StatusPane{
			{Text: message},
			{Text: fmt.Sprintf("%d items", len(items)), Width: 74},
			{Text: "Sort: " + sortLabel(currentSort), Width: 92},
		})
		app.Desktop().InvalidateRect(win.Bounds())
	}
	applySort := func(cmd classicui.CommandID, message string) {
		currentSort = cmd
		sortByNameItem.Checked = cmd == cmdSortByName
		sortByLengthItem.Checked = cmd == cmdSortByLength
		toolbar.SetChecked(cmdSortByName, cmd == cmdSortByName)
		toolbar.SetChecked(cmdSortByLength, cmd == cmdSortByLength)
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
		default:
			sort.SliceStable(items, func(i, j int) bool {
				return items[i] < items[j]
			})
		}
		list.SetItems(items)
		updateStatus(message)
	}
	syncDefaultButton := func() {
		if tabs.SelectedIndex() == 1 {
			win.SetDefaultButton(aboutBtn)
			return
		}
		win.SetDefaultButton(add)
	}

	runCommand := func(cmd classicui.CommandID) {
		switch cmd {
		case cmdAddPath:
			next := pathEdit.Text()
			if next == "" {
				updateStatus("输入框不能为空。")
				return
			}
			items = append(items, next)
			applySort(currentSort, fmt.Sprintf("已添加：%s", next))
			win.SetTitle(fmt.Sprintf("Phase 4 - 共 %d 项", len(items)))
		case cmdExit:
			app.Quit()
		case cmdSortByName, cmdSortByLength:
			message := "已切换为按名称排序。"
			if cmd == cmdSortByLength {
				message = "已切换为按名称长度排序。"
			}
			applySort(cmd, message)
		case cmdAbout:
			updateStatus("Phase 4: TabControl 已接上现有命令流。")
		}
		app.Desktop().InvalidateRect(win.Bounds())
	}
	app.OnCommand(runCommand)

	list.OnChange(func(index int, value string) {
		updateStatus(fmt.Sprintf("当前选中：%s", value))
	})
	tabs.OnSelectionChange(func(index int, page *widgets.TabPage) {
		syncDefaultButton()
		updateStatus(fmt.Sprintf("当前页签：%s", page.Title))
	})

	add.OnClick(func() {
		runCommand(cmdAddPath)
	})
	closeBtn.OnClick(func() {
		runCommand(cmdExit)
	})
	sortNameBtn.OnClick(func() {
		runCommand(cmdSortByName)
	})
	sortLengthBtn.OnClick(func() {
		runCommand(cmdSortByLength)
	})
	aboutBtn.OnClick(func() {
		runCommand(cmdAbout)
	})

	win.Content().Add(title)
	win.Content().Add(toolbar)
	win.Content().Add(tabs)
	win.Content().Add(statusBar)
	syncDefaultButton()
	app.Desktop().AddWindow(win)
	applySort(currentSort, "试试页签、Alt、方向键、Ctrl+N 和 Ctrl+Q。")

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
