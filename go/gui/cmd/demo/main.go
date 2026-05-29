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
	toolbarSortByName.Tooltip = "按名称排序"
	toolbarSortByLength := widgets.NewToolbarButton(cmdSortByLength, "By Length")
	toolbarSortByLength.Tooltip = "按名称长度排序"
	toolbarAddPath := widgets.NewToolbarButton(cmdAddPath, "Add Path")
	toolbarAddPath.Tooltip = "把输入框内容添加到列表"
	toolbarAbout := widgets.NewToolbarButton(cmdAbout, "About")
	toolbarAbout.Tooltip = "查看当前阶段说明"
	toolbar := widgets.NewToolbar("toolbar", classicui.Rect{
		X: 14,
		Y: 12,
		W: 392,
		H: 28,
	},
		toolbarAddPath,
		widgets.NewToolbarSeparator(),
		toolbarSortByName,
		toolbarSortByLength,
		widgets.NewToolbarSeparator(),
		toolbarAbout,
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
		Y: 146,
		W: 76,
		H: 24,
	})
	closeBtn := widgets.NewButton("cancel", "关闭", classicui.Rect{
		X: 292,
		Y: 146,
		W: 76,
		H: 24,
	})
	sortNameBtn := widgets.NewButton("sortName", "按名称", classicui.Rect{
		X: 10,
		Y: 78,
		W: 88,
		H: 24,
	})
	aboutBtn := widgets.NewButton("about", "关于", classicui.Rect{
		X: 104,
		Y: 78,
		W: 88,
		H: 24,
	})
	sortModeLabel := widgets.NewLabel("sortModeLabel", "排序方式：", classicui.Rect{
		X: 10,
		Y: 44,
		W: 72,
		H: 18,
	})
	sortCombo := widgets.NewComboBox("sortCombo", classicui.Rect{
		X: 76,
		Y: 40,
		W: 146,
		H: 24,
	})
	sortCombo.SetItems([]string{"按名称", "按长度"})
	sortCombo.SetEditable(true)
	sortCombo.SetTooltip("用下拉框切换排序命令")
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
	treeHint := widgets.NewLabel("treeHint", "TreeView 已支持双击、hover 热区，以及 F2 / 慢单击的 inline rename。", classicui.Rect{
		X: 10,
		Y: 12,
		W: 350,
		H: 18,
	})
	treeView := widgets.NewTreeView("tree", classicui.Rect{
		X: 10,
		Y: 40,
		W: 358,
		H: 138,
	},
		func() *widgets.TreeNode {
			root := widgets.NewFolderNode("桌面",
				widgets.NewFolderNode("我的电脑",
					widgets.NewFolderNode("Program Files",
						widgets.NewFolderNode("InHarmony",
							widgets.NewFileNode("dust.exe"),
							widgets.NewFileNode("readme.txt"),
						),
					),
					widgets.NewFolderNode("Windows",
						widgets.NewFileNode("explorer.exe"),
						widgets.NewFileNode("notepad.exe"),
					),
					widgets.NewFolderNode("Temp",
						widgets.NewFileNode("session.log"),
					),
				),
				widgets.NewFolderNode("文档",
					widgets.NewFolderNode("项目",
						widgets.NewFileNode("main.go"),
						widgets.NewFileNode("ui_spec.md"),
					),
					widgets.NewFolderNode("归档",
						widgets.NewFileNode("phase3-notes.txt"),
					),
				),
				widgets.NewFolderNode("下载",
					widgets.NewFolderNode("安装包",
						widgets.NewFileNode("classicui-setup.exe"),
					),
					widgets.NewFolderNode("压缩包",
						widgets.NewFileNode("assets.zip"),
					),
				),
				widgets.NewFileNode("README.txt"),
			)
			root.Expanded = true
			root.Children[0].Expanded = true
			root.Children[0].Children[0].Expanded = true
			return root
		}(),
	)
	browsePage := widgets.NewPanel("browsePage", classicui.Rect{})
	browsePage.Add(pathLabel)
	browsePage.Add(pathEdit)
	browsePage.Add(list)
	browsePage.Add(add)
	browsePage.Add(closeBtn)
	commandsPage := widgets.NewPanel("commandsPage", classicui.Rect{})
	commandsPage.Add(info1)
	commandsPage.Add(sortModeLabel)
	commandsPage.Add(sortCombo)
	commandsPage.Add(sortNameBtn)
	commandsPage.Add(aboutBtn)
	commandsPage.Add(info2)
	commandsPage.Add(info3)
	treePage := widgets.NewPanel("treePage", classicui.Rect{})
	treePage.Add(treeHint)
	treePage.Add(treeView)
	tabs := widgets.NewTabControl("tabs", classicui.Rect{
		X: 14,
		Y: 72,
		W: 392,
		H: 204,
	},
		widgets.NewTabPage("文件列表", browsePage),
		widgets.NewTabPage("命令演示", commandsPage),
		widgets.NewTabPage("目录树", treePage),
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
		if cmd == cmdSortByLength {
			sortCombo.SetSelectedIndexSilent(1)
		} else {
			sortCombo.SetSelectedIndexSilent(0)
		}
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
		if tabs.SelectedIndex() == 2 {
			win.SetDefaultButton(nil)
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
	treeView.OnChange(func(node *widgets.TreeNode) {
		if node != nil {
			count := len(treeView.SelectedNodes())
			if count > 1 {
				updateStatus(fmt.Sprintf("当前树节点：%s（已选 %d 项）", node.Text, count))
				return
			}
			updateStatus(fmt.Sprintf("当前树节点：%s", node.Text))
		}
	})
	treeView.OnActivate(func(node *widgets.TreeNode) {
		if node != nil {
			updateStatus(fmt.Sprintf("已激活树节点：%s", node.Text))
		}
	})
	treeView.OnBeginRename(func(node *widgets.TreeNode) {
		if node != nil {
			updateStatus(fmt.Sprintf("重命名预留接口：%s", node.Text))
		}
	})
	treeView.OnRenameCommit(func(node *widgets.TreeNode, oldText, newText string) {
		if node != nil {
			updateStatus(fmt.Sprintf("树节点已重命名：%s -> %s", oldText, newText))
		}
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
	aboutBtn.OnClick(func() {
		runCommand(cmdAbout)
	})
	sortCombo.OnCommit(func(index int, _ string) {
		if index < 0 {
			updateStatus("请输入“按名称”或“按长度”后再提交。")
			return
		}
		if index == 1 {
			runCommand(cmdSortByLength)
			return
		}
		runCommand(cmdSortByName)
	})
	add.SetTooltip("把当前路径加入列表")
	closeBtn.SetTooltip("关闭当前演示窗口")
	sortNameBtn.SetTooltip("触发按名称排序命令")
	aboutBtn.SetTooltip("在状态栏显示阶段说明")

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
