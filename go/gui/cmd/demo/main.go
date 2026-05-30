package main

import (
	"flag"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	"classicui"
	"classicui/event"
	"classicui/widgets"
)

const (
	cmdAddPath            classicui.CommandID = "cmd.file.add_path"
	cmdOpenSelection      classicui.CommandID = "cmd.file.open_selection"
	cmdFileRename         classicui.CommandID = "cmd.file.rename"
	cmdNavigateUp         classicui.CommandID = "cmd.file.navigate_up"
	cmdViewRefresh        classicui.CommandID = "cmd.view.refresh"
	cmdExit               classicui.CommandID = "cmd.file.exit"
	cmdSortByName         classicui.CommandID = "cmd.view.sort.name"
	cmdSortBySize         classicui.CommandID = "cmd.view.sort.size"
	cmdSortByType         classicui.CommandID = "cmd.view.sort.type"
	cmdToggleSingleSelect classicui.CommandID = "cmd.view.single_select"
	cmdAbout              classicui.CommandID = "cmd.help.about"
)

type demoEntryKind int

const (
	demoFolder demoEntryKind = iota
	demoFile
)

type demoEntry struct {
	Name     string
	Kind     demoEntryKind
	Size     int
	Children []*demoEntry
	Parent   *demoEntry
	TreeNode *widgets.TreeNode
}

func newFolder(name string, children ...*demoEntry) *demoEntry {
	entry := &demoEntry{Name: name, Kind: demoFolder}
	for _, child := range children {
		entry.AddChild(child)
	}
	return entry
}

func newFile(name string, size int) *demoEntry {
	return &demoEntry{Name: name, Kind: demoFile, Size: size}
}

func (e *demoEntry) AddChild(child *demoEntry) {
	if e == nil || child == nil {
		return
	}
	child.Parent = e
	e.Children = append(e.Children, child)
}

func (e *demoEntry) IsFolder() bool {
	return e != nil && e.Kind == demoFolder
}

func (e *demoEntry) Path() string {
	if e == nil {
		return ""
	}
	var parts []string
	for current := e; current != nil; current = current.Parent {
		parts = append(parts, current.Name)
	}
	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}
	return strings.Join(parts, `\`)
}

func (e *demoEntry) DisplayName() string {
	if e == nil {
		return ""
	}
	if e.IsFolder() {
		return "[" + e.Name + "]"
	}
	return e.Name
}

func buildDemoRoot() (*demoEntry, map[*widgets.TreeNode]*demoEntry) {
	root := newFolder("桌面",
		newFolder("我的电脑",
			newFolder("Program Files",
				newFolder("InHarmony",
					newFile("dust.exe", 4120),
					newFile("readme.txt", 14),
					newFile("changelog.md", 32),
				),
			),
			newFolder("Windows",
				newFile("explorer.exe", 1452),
				newFile("notepad.exe", 188),
				newFile("system.ini", 4),
			),
			newFolder("Temp",
				newFile("session.log", 27),
				newFile("cache.bin", 96),
			),
		),
		newFolder("文档",
			newFolder("项目",
				newFile("main.go", 73),
				newFile("ui_spec.md", 28),
				newFile("todo.txt", 3),
			),
			newFolder("归档",
				newFile("phase3-notes.txt", 12),
				newFile("retro-assets.zip", 904),
			),
		),
		newFolder("下载",
			newFolder("安装包",
				newFile("classicui-setup.exe", 2148),
			),
			newFolder("压缩包",
				newFile("assets.zip", 640),
				newFile("textures.zip", 522),
			),
		),
		newFolder("图片",
			newFile("splash.bmp", 768),
			newFile("icon-set.ico", 64),
		),
		newFile("README.txt", 8),
	)

	treeIndex := make(map[*widgets.TreeNode]*demoEntry)
	bindTree(root, treeIndex)
	root.TreeNode.Expanded = true
	if len(root.TreeNode.Children) > 0 {
		root.TreeNode.Children[0].Expanded = true
	}
	if len(root.TreeNode.Children) > 1 {
		root.TreeNode.Children[1].Expanded = true
	}
	return root, treeIndex
}

func sortDemoFolder(folder *demoEntry, mode classicui.CommandID, descending bool) {
	if folder == nil || !folder.IsFolder() {
		return
	}
	sort.SliceStable(folder.Children, func(i, j int) bool {
		left := folder.Children[i]
		right := folder.Children[j]
		if left.IsFolder() != right.IsFolder() {
			return left.IsFolder()
		}
		switch mode {
		case cmdSortBySize:
			if left.IsFolder() {
				return left.Name < right.Name
			}
			if left.Size == right.Size {
				return left.Name < right.Name
			}
			return left.Size < right.Size
		case cmdSortByType:
			leftType := demoTypeText(left)
			rightType := demoTypeText(right)
			if leftType == rightType {
				return left.Name < right.Name
			}
			return leftType < rightType
		default:
			return left.Name < right.Name
		}
	})
	if descending {
		firstFile := firstDemoFileIndex(folder.Children)
		for i, j := 0, firstFile-1; i < j; i, j = i+1, j-1 {
			folder.Children[i], folder.Children[j] = folder.Children[j], folder.Children[i]
		}
		for i, j := firstFile, len(folder.Children)-1; i < j; i, j = i+1, j-1 {
			folder.Children[i], folder.Children[j] = folder.Children[j], folder.Children[i]
		}
	}
	if folder.TreeNode != nil {
		folder.TreeNode.Children = folder.TreeNode.Children[:0]
		for _, child := range folder.Children {
			folder.TreeNode.Children = append(folder.TreeNode.Children, child.TreeNode)
		}
	}
}

func findDemoEntryIndex(entries []*demoEntry, target *demoEntry) int {
	for i, entry := range entries {
		if entry == target {
			return i
		}
	}
	return -1
}

func resolveDemoSelection(entry *demoEntry) (folder *demoEntry, selectEntry *demoEntry) {
	if entry == nil {
		return nil, nil
	}
	if entry.IsFolder() {
		return entry, nil
	}
	if entry.Parent != nil {
		return entry.Parent, entry
	}
	return entry, nil
}

func addEntryToFolder(folder *demoEntry, raw string, treeIndex map[*widgets.TreeNode]*demoEntry) (*demoEntry, string) {
	if folder == nil || !folder.IsFolder() {
		return nil, "当前目录不可写入。"
	}
	name := strings.TrimSpace(raw)
	name = strings.TrimRight(name, `\/`)
	if name == "" {
		return nil, "路径栏不能为空。"
	}
	if strings.Contains(name, `\`) || strings.Contains(name, "/") {
		parts := strings.FieldsFunc(name, func(r rune) bool {
			return r == '\\' || r == '/'
		})
		if len(parts) == 0 {
			return nil, "路径栏不能为空。"
		}
		name = parts[len(parts)-1]
	}
	for _, child := range folder.Children {
		if strings.EqualFold(child.Name, name) {
			return nil, "当前目录已存在同名项目。"
		}
	}
	var child *demoEntry
	if strings.ContainsRune(name, '.') {
		child = newFile(name, maxInt(len([]rune(name))*64, 8))
	} else {
		child = newFolder(name)
	}
	folder.AddChild(child)
	bindTree(child, treeIndex)
	if folder.TreeNode != nil {
		folder.TreeNode.AddChild(child.TreeNode)
		folder.TreeNode.Expanded = true
	}
	return child, ""
}

func applyDemoSelectionMode(list *widgets.ListBox, tree *widgets.TreeView, single bool) {
	if list != nil {
		options := list.SelectionOptions()
		options.MultiSelect = !single
		list.SetSelectionOptions(options)
	}
	if tree != nil {
		options := tree.SelectionOptions()
		options.MultiSelect = !single
		tree.SetSelectionOptions(options)
	}
}

func applyDemoListViewSelectionMode(list *widgets.ListView, tree *widgets.TreeView, single bool) {
	if list != nil {
		options := list.SelectionOptions()
		options.MultiSelect = !single
		list.SetSelectionOptions(options)
	}
	if tree != nil {
		options := tree.SelectionOptions()
		options.MultiSelect = !single
		tree.SetSelectionOptions(options)
	}
}

func demoSizeText(entry *demoEntry) string {
	if entry == nil || entry.IsFolder() {
		return ""
	}
	return fmt.Sprintf("%d KB", entry.Size)
}

func demoTypeText(entry *demoEntry) string {
	if entry == nil {
		return ""
	}
	if entry.IsFolder() {
		return "文件夹"
	}
	dot := strings.LastIndex(entry.Name, ".")
	if dot < 0 || dot == len(entry.Name)-1 {
		return "文件"
	}
	return strings.ToUpper(entry.Name[dot+1:]) + " 文件"
}

func demoListViewItem(entry *demoEntry) widgets.ListViewItem {
	return widgets.ListViewItem{
		Texts: []string{entry.DisplayName(), demoSizeText(entry), demoTypeText(entry)},
		Data:  entry,
	}
}

func demoListViewRows(entries []*demoEntry) []widgets.ListViewItem {
	rows := make([]widgets.ListViewItem, len(entries))
	for i, entry := range entries {
		rows[i] = demoListViewItem(entry)
	}
	return rows
}

func buildDemoContextMenu(
	hasEntry bool,
	entry *demoEntry,
	renameItem *widgets.MenuItem,
	refreshItem *widgets.MenuItem,
	sortByNameItem *widgets.MenuItem,
	sortBySizeItem *widgets.MenuItem,
	sortByTypeItem *widgets.MenuItem,
) *widgets.Menu {
	if !hasEntry || entry == nil {
		return widgets.NewMenu(
			widgets.NewMenuItem(cmdAddPath, "&Add Item", nil),
			refreshItem,
			widgets.NewSeparator(),
			widgets.NewSubmenuItem("&Sort", widgets.NewMenu(
				sortByNameItem,
				sortBySizeItem,
				sortByTypeItem,
			)),
		)
	}
	openText := "&Open"
	if !entry.IsFolder() {
		openText = "&Open File"
	}
	items := []*widgets.MenuItem{
		widgets.NewMenuItem(cmdOpenSelection, openText, nil),
		renameItem,
	}
	if entry.IsFolder() {
		items = append(items, widgets.NewMenuItem(cmdNavigateUp, "Open &Parent", nil))
	}
	items = append(items,
		widgets.NewSeparator(),
		widgets.NewMenuItem(cmdAddPath, "&Add Item", nil),
		refreshItem,
		widgets.NewSubmenuItem("&Sort", widgets.NewMenu(
			sortByNameItem,
			sortBySizeItem,
			sortByTypeItem,
		)),
	)
	return widgets.NewMenu(items...)
}

func firstDemoFileIndex(entries []*demoEntry) int {
	for i, entry := range entries {
		if entry == nil || !entry.IsFolder() {
			return i
		}
	}
	return len(entries)
}

func main() {
	autoQuit := flag.Duration("auto-quit", 0, "automatically exit after the given duration")
	flag.Parse()

	app := classicui.NewApp(classicui.Config{
		Title:        "Classic UI Explorer Demo",
		LogicalSize:  classicui.Size{W: 800, H: 560},
		PresentScale: 2,
		Theme:        classicui.DefaultClassicTheme(),
	})

	root, treeIndex := buildDemoRoot()

	appWin := classicui.NewWindow("explorer", classicui.Rect{
		X: 52,
		Y: 28,
		W: 580,
		H: 426,
	})
	appWin.SetTitle("Classic Explorer")

	toolbarUp := widgets.NewToolbarButton(cmdNavigateUp, "Up")
	toolbarUp.Tooltip = "切换到上级目录"
	toolbarAddPath := widgets.NewToolbarButton(cmdAddPath, "Add")
	toolbarAddPath.Tooltip = "把路径栏内容加入当前目录"
	toolbarSingle := widgets.NewToolbarButton(cmdToggleSingleSelect, "Single")
	toolbarSingle.Tooltip = "切换单选/多选模式"
	toolbarSortByName := widgets.NewToolbarButton(cmdSortByName, "By Name")
	toolbarSortByName.Checked = true
	toolbarSortByName.Tooltip = "按名称排序当前目录"
	toolbarSortBySize := widgets.NewToolbarButton(cmdSortBySize, "By Size")
	toolbarSortBySize.Tooltip = "按大小排序当前目录"
	toolbarSortByType := widgets.NewToolbarButton(cmdSortByType, "By Type")
	toolbarSortByType.Tooltip = "按类型排序当前目录"
	toolbarAbout := widgets.NewToolbarButton(cmdAbout, "About")
	toolbarAbout.Tooltip = "显示当前 demo 的验证重点"
	toolbar := widgets.NewToolbar("toolbar", classicui.Rect{
		X: 12,
		Y: 12,
		W: 548,
		H: 28,
	},
		toolbarUp,
		toolbarAddPath,
		widgets.NewToolbarSeparator(),
		toolbarSingle,
		widgets.NewToolbarSeparator(),
		toolbarSortByName,
		toolbarSortBySize,
		toolbarSortByType,
		widgets.NewToolbarSeparator(),
		toolbarAbout,
	)

	title := widgets.NewLabel("intro", "集成 demo：目录树、文件列表、命令、状态栏和页签全部联动。", classicui.Rect{
		X: 12,
		Y: 48,
		W: 548,
		H: 18,
	})

	pathLabel := widgets.NewLabel("pathLabel", "路径：", classicui.Rect{X: 10, Y: 12, W: 40, H: 18})
	pathEdit := widgets.NewEdit("path", classicui.Rect{X: 48, Y: 8, W: 292, H: 24})
	sortLabel := widgets.NewLabel("sortLabel", "排序：", classicui.Rect{X: 346, Y: 12, W: 36, H: 18})
	sortCombo := widgets.NewComboBox("sortCombo", classicui.Rect{X: 382, Y: 8, W: 128, H: 24})
	sortCombo.SetItems([]string{"按名称", "按大小", "按类型"})
	sortCombo.SetEditable(true)
	sortCombo.SetTooltip("切换当前目录的排序方式")

	treeView := widgets.NewTreeView("tree", classicui.Rect{X: 10, Y: 42, W: 196, H: 172}, root.TreeNode)
	list := widgets.NewListView("files", classicui.Rect{X: 214, Y: 42, W: 296, H: 172},
		widgets.ListViewColumn{Title: "名称", Width: 144},
		widgets.ListViewColumn{Title: "大小", Width: 56, Align: widgets.HeaderAlignRight},
		widgets.ListViewColumn{Title: "类型", Width: 76},
	)
	list.SetSortIndicator(0, false)
	upBtn := widgets.NewButton("up", "上级", classicui.Rect{X: 272, Y: 222, W: 72, H: 24})
	addBtn := widgets.NewButton("add", "添加", classicui.Rect{X: 350, Y: 222, W: 72, H: 24})
	closeBtn := widgets.NewButton("close", "关闭", classicui.Rect{X: 438, Y: 222, W: 72, H: 24})
	upBtn.SetTooltip("返回当前目录的父目录")
	addBtn.SetTooltip("使用路径栏文本创建项目")
	closeBtn.SetTooltip("关闭当前演示窗口")

	browsePage := widgets.NewPanel("browsePage", classicui.Rect{})
	browsePage.Add(pathLabel)
	browsePage.Add(pathEdit)
	browsePage.Add(sortLabel)
	browsePage.Add(sortCombo)
	browsePage.Add(treeView)
	browsePage.Add(list)
	browsePage.Add(upBtn)
	browsePage.Add(addBtn)
	browsePage.Add(closeBtn)

	statePath := widgets.NewLabel("statePath", "", classicui.Rect{X: 10, Y: 12, W: 500, H: 18})
	stateTree := widgets.NewLabel("stateTree", "", classicui.Rect{X: 10, Y: 40, W: 500, H: 18})
	stateList := widgets.NewLabel("stateList", "", classicui.Rect{X: 10, Y: 68, W: 500, H: 18})
	stateMode := widgets.NewLabel("stateMode", "", classicui.Rect{X: 10, Y: 96, W: 500, H: 18})
	stateHint1 := widgets.NewLabel("stateHint1", "验证重点：树选择驱动列表，列表激活驱动树导航。", classicui.Rect{X: 10, Y: 132, W: 500, H: 18})
	stateHint2 := widgets.NewLabel("stateHint2", "再试试 Ctrl+A、Ctrl+Space、空白点击、F2、慢单击重命名。", classicui.Rect{X: 10, Y: 156, W: 500, H: 18})
	stateHint3 := widgets.NewLabel("stateHint3", "菜单、工具栏、按钮、下拉框都走同一套命令入口。", classicui.Rect{X: 10, Y: 180, W: 500, H: 18})
	statePage := widgets.NewPanel("statePage", classicui.Rect{})
	statePage.Add(statePath)
	statePage.Add(stateTree)
	statePage.Add(stateList)
	statePage.Add(stateMode)
	statePage.Add(stateHint1)
	statePage.Add(stateHint2)
	statePage.Add(stateHint3)

	help1 := widgets.NewLabel("help1", "帮助：", classicui.Rect{X: 10, Y: 12, W: 80, H: 18})
	help2 := widgets.NewLabel("help2", "1. 左侧 TreeView 选目录，右侧 ListView 自动切换内容。", classicui.Rect{X: 10, Y: 40, W: 520, H: 18})
	help3 := widgets.NewLabel("help3", "2. 双击或激活右侧目录项，会把树选择切到该目录。", classicui.Rect{X: 10, Y: 64, W: 520, H: 18})
	help4 := widgets.NewLabel("help4", "3. 切到单选模式后，Ctrl/Shift 扩选和空白框选会被策略层统一关闭。", classicui.Rect{X: 10, Y: 88, W: 520, H: 18})
	help5 := widgets.NewLabel("help5", "4. 点击列头或切换排序命令，只会重排当前目录，但树和列表会同步。", classicui.Rect{X: 10, Y: 112, W: 520, H: 18})
	help6 := widgets.NewLabel("help6", "5. Ctrl+N 添加项目，Ctrl+Q 退出，Alt 菜单、工具栏和状态栏共享状态。", classicui.Rect{X: 10, Y: 136, W: 520, H: 18})
	helpPage := widgets.NewPanel("helpPage", classicui.Rect{})
	helpPage.Add(help1)
	helpPage.Add(help2)
	helpPage.Add(help3)
	helpPage.Add(help4)
	helpPage.Add(help5)
	helpPage.Add(help6)

	tabs := widgets.NewTabControl("tabs", classicui.Rect{
		X: 12,
		Y: 72,
		W: 548,
		H: 274,
	},
		widgets.NewTabPage("浏览器", browsePage),
		widgets.NewTabPage("状态", statePage),
		widgets.NewTabPage("帮助", helpPage),
	)

	statusBar := widgets.NewStatusBar("status", classicui.Rect{
		X: 12,
		Y: 356,
		W: 548,
		H: 22,
	})

	sortByNameItem := widgets.NewMenuItem(cmdSortByName, "By &Name", nil)
	sortByNameItem.Checked = true
	sortBySizeItem := widgets.NewMenuItem(cmdSortBySize, "By &Size", nil)
	sortByTypeItem := widgets.NewMenuItem(cmdSortByType, "By &Type", nil)
	renameItem := widgets.NewMenuItem(cmdFileRename, "&Rename", &widgets.Accelerator{
		Key: event.KeyF2,
	})
	refreshItem := widgets.NewMenuItem(cmdViewRefresh, "&Refresh", nil)
	singleSelectItem := widgets.NewMenuItem(cmdToggleSingleSelect, "&Single Selection", nil)

	appWin.SetMenuBar(widgets.NewMenuBar(
		widgets.NewSubmenuItem("&File", widgets.NewMenu(
			widgets.NewMenuItem(cmdAddPath, "&Add Item", &widgets.Accelerator{
				Key:       event.KeyN,
				Modifiers: event.ModCtrl,
			}),
			renameItem,
			widgets.NewMenuItem(cmdNavigateUp, "&Up", &widgets.Accelerator{
				Key:       event.KeyBackspace,
				Modifiers: event.ModAlt,
			}),
			widgets.NewSeparator(),
			widgets.NewMenuItem(cmdExit, "E&xit", &widgets.Accelerator{
				Key:       event.KeyQ,
				Modifiers: event.ModCtrl,
			}),
		)),
		widgets.NewSubmenuItem("&View", widgets.NewMenu(
			refreshItem,
			widgets.NewSeparator(),
			widgets.NewSubmenuItem("&Sort", widgets.NewMenu(
				sortByNameItem,
				sortBySizeItem,
				sortByTypeItem,
			)),
			singleSelectItem,
		)),
		widgets.NewSubmenuItem("&Help", widgets.NewMenu(
			widgets.NewMenuItem(cmdAbout, "&About", nil),
		)),
	))
	treeView.SetContextMenuProvider(func(info widgets.TreeViewContextMenuInfo) *widgets.Menu {
		if !info.HasNode || info.Node == nil {
			return buildDemoContextMenu(false, nil, renameItem, refreshItem, sortByNameItem, sortBySizeItem, sortByTypeItem)
		}
		return buildDemoContextMenu(true, treeIndex[info.Node], renameItem, refreshItem, sortByNameItem, sortBySizeItem, sortByTypeItem)
	})
	list.SetContextMenuProvider(func(info widgets.ListViewContextMenuInfo) *widgets.Menu {
		if !info.HasItem {
			return buildDemoContextMenu(false, nil, renameItem, refreshItem, sortByNameItem, sortBySizeItem, sortByTypeItem)
		}
		entry, _ := info.Item.Data.(*demoEntry)
		return buildDemoContextMenu(entry != nil, entry, renameItem, refreshItem, sortByNameItem, sortBySizeItem, sortByTypeItem)
	})

	currentSort := cmdSortByName
	sortDescending := false
	currentFolder := root
	currentEntries := []*demoEntry(nil)
	statusMessage := ""

	sortLabelText := func(cmd classicui.CommandID) string {
		switch cmd {
		case cmdSortBySize:
			return "Size"
		case cmdSortByType:
			return "Type"
		default:
			return "Name"
		}
	}
	sortColumnIndex := func(cmd classicui.CommandID) int {
		switch cmd {
		case cmdSortBySize:
			return 1
		case cmdSortByType:
			return 2
		default:
			return 0
		}
	}
	describeNodes := func(nodes []*widgets.TreeNode) string {
		if len(nodes) == 0 {
			return "(无)"
		}
		names := make([]string, 0, minInt(len(nodes), 3))
		for i, node := range nodes {
			if i == 3 {
				names = append(names, "...")
				break
			}
			names = append(names, node.Text)
		}
		return strings.Join(names, ", ")
	}
	describeListSelection := func() string {
		indices := list.SelectedIndices()
		if len(indices) == 0 {
			return "(无)"
		}
		names := make([]string, 0, minInt(len(indices), 3))
		for i, index := range indices {
			if i == 3 {
				names = append(names, "...")
				break
			}
			if index >= 0 && index < len(currentEntries) {
				names = append(names, currentEntries[index].DisplayName())
			}
		}
		return strings.Join(names, ", ")
	}
	updateStatus := func(message string) {
		if message != "" {
			statusMessage = message
		}
		selectedCount := len(list.SelectedIndices())
		statusBar.SetPanes([]widgets.StatusPane{
			{Text: statusMessage},
			{Text: fmt.Sprintf("%d items", len(currentEntries)), Width: 72},
			{Text: fmt.Sprintf("%d sel", selectedCount), Width: 56},
			{Text: "Sort: " + sortLabelText(currentSort), Width: 88},
		})
		statePath.SetText("当前路径: " + currentFolder.Path())
		stateTree.SetText("树选择: " + describeNodes(treeView.SelectedNodes()))
		stateList.SetText("列表选择: " + describeListSelection())
		modeText := "多选"
		if !list.SelectionOptions().MultiSelect {
			modeText = "单选"
		}
		stateMode.SetText("选择模式: " + modeText)
		appWin.SetTitle("Classic Explorer - " + currentFolder.Path())
		app.Desktop().InvalidateRect(appWin.Bounds())
	}

	refreshCurrentFolder := func(selectEntry *demoEntry) {
		sortDemoFolder(currentFolder, currentSort, sortDescending)
		currentEntries = append(currentEntries[:0], currentFolder.Children...)
		list.SetItems(demoListViewRows(currentEntries))
		list.SetSelectionOptions(list.SelectionOptions())
		list.SetSelectedIndexSilent(-1)
		if selectEntry != nil {
			if index := findDemoEntryIndex(currentEntries, selectEntry); index >= 0 {
				list.SetSelectedIndexSilent(index)
			}
		}
		pathEdit.SetText(currentFolder.Path())
	}

	applySelectionMode := func(single bool) {
		applyDemoListViewSelectionMode(list, treeView, single)
		toolbar.SetChecked(cmdToggleSingleSelect, single)
		singleSelectItem.Checked = single
		if single {
			updateStatus("已切换为单选模式。")
			return
		}
		updateStatus("已切换为多选模式。")
	}

	applySort := func(cmd classicui.CommandID, message string) {
		currentSort = cmd
		sortByNameItem.Checked = cmd == cmdSortByName
		sortBySizeItem.Checked = cmd == cmdSortBySize
		sortByTypeItem.Checked = cmd == cmdSortByType
		toolbar.SetChecked(cmdSortByName, cmd == cmdSortByName)
		toolbar.SetChecked(cmdSortBySize, cmd == cmdSortBySize)
		toolbar.SetChecked(cmdSortByType, cmd == cmdSortByType)
		list.SetSortIndicator(sortColumnIndex(cmd), sortDescending)
		switch cmd {
		case cmdSortBySize:
			sortCombo.SetSelectedIndexSilent(1)
		case cmdSortByType:
			sortCombo.SetSelectedIndexSilent(2)
		default:
			sortCombo.SetSelectedIndexSilent(0)
		}
		refreshCurrentFolder(nil)
		updateStatus(message)
	}

	syncDefaultButton := func() {
		if tabs.SelectedIndex() == 0 {
			appWin.SetDefaultButton(addBtn)
			return
		}
		appWin.SetDefaultButton(nil)
	}

	runCommand := func(cmd classicui.CommandID) {
		switch cmd {
		case cmdAddPath:
			child, errMessage := addEntryToFolder(currentFolder, pathEdit.Text(), treeIndex)
			if errMessage != "" {
				updateStatus(errMessage)
				return
			}
			refreshCurrentFolder(child)
			updateStatus("已添加到当前目录: " + child.DisplayName())
		case cmdOpenSelection:
			if index := list.SelectedIndex(); index >= 0 && index < len(currentEntries) {
				entry := currentEntries[index]
				treeView.SetSelectedNode(entry.TreeNode)
				if entry.IsFolder() {
					updateStatus("已打开目录: " + entry.Path())
				} else {
					updateStatus(fmt.Sprintf("已打开文件: %s (%d KB)", entry.Path(), entry.Size))
				}
				return
			}
			if node := treeView.SelectedNode(); node != nil {
				if entry := treeIndex[node]; entry != nil {
					if entry.IsFolder() {
						updateStatus("已打开目录: " + entry.Path())
					} else {
						updateStatus(fmt.Sprintf("已打开文件: %s (%d KB)", entry.Path(), entry.Size))
					}
				}
				return
			}
			updateStatus("当前没有可打开的项目。")
		case cmdFileRename:
			if index := list.SelectedIndex(); index >= 0 && index < len(currentEntries) {
				if !list.BeginRename(index) {
					updateStatus("当前项目无法重命名。")
					return
				}
				if entry := currentEntries[index]; entry != nil {
					updateStatus("开始重命名: " + entry.Path())
				}
				return
			}
			node := treeView.SelectedNode()
			entry := treeIndex[node]
			if entry == nil || entry.TreeNode == nil {
				updateStatus("当前没有可重命名的项目。")
				return
			}
			if !treeView.BeginRename(entry.TreeNode) {
				updateStatus("当前项目无法重命名。")
				return
			}
			updateStatus("开始重命名: " + entry.Path())
		case cmdNavigateUp:
			if currentFolder.Parent == nil {
				updateStatus("已经在根目录。")
				return
			}
			treeView.SetSelectedNode(currentFolder.Parent.TreeNode)
			updateStatus("已切换到上级目录。")
		case cmdViewRefresh:
			refreshCurrentFolder(nil)
			updateStatus("已刷新当前目录。")
		case cmdSortByName:
			sortDescending = false
			applySort(cmd, "已切换为按名称排序。")
		case cmdSortBySize:
			sortDescending = false
			applySort(cmd, "已切换为按大小排序。")
		case cmdSortByType:
			sortDescending = false
			applySort(cmd, "已切换为按类型排序。")
		case cmdToggleSingleSelect:
			applySelectionMode(!list.SelectionOptions().MultiSelect)
		case cmdAbout:
			updateStatus("Explorer Demo：验证树/列表/命令/状态同步，而不是单个控件展示。")
		case cmdExit:
			app.Quit()
		}
	}
	app.OnCommand(runCommand)

	treeView.OnChange(func(node *widgets.TreeNode) {
		entry := treeIndex[node]
		if entry == nil {
			return
		}
		nextFolder, selectEntry := resolveDemoSelection(entry)
		if nextFolder != nil {
			currentFolder = nextFolder
		}
		refreshCurrentFolder(selectEntry)
		count := len(treeView.SelectedNodes())
		if count > 1 {
			updateStatus(fmt.Sprintf("树选择已更新：%s（已选 %d 项）", entry.Name, count))
			return
		}
		if entry.IsFolder() {
			updateStatus("当前目录: " + currentFolder.Path())
			return
		}
		updateStatus("当前文件: " + entry.Path())
	})
	treeView.OnActivate(func(node *widgets.TreeNode) {
		if entry := treeIndex[node]; entry != nil {
			if entry.IsFolder() {
				updateStatus("已激活目录: " + entry.Path())
			} else {
				updateStatus("已激活文件: " + entry.Path())
			}
		}
	})
	treeView.OnBeginRename(func(node *widgets.TreeNode) {
		if entry := treeIndex[node]; entry != nil {
			updateStatus("开始重命名: " + entry.Path())
		}
	})
	treeView.OnRenameCommit(func(node *widgets.TreeNode, oldText, newText string) {
		entry := treeIndex[node]
		if entry == nil {
			return
		}
		entry.Name = newText
		refreshCurrentFolder(entry)
		updateStatus(fmt.Sprintf("已重命名: %s -> %s", oldText, newText))
	})

	list.OnChange(func(index int, item widgets.ListViewItem) {
		if index < 0 || index >= len(currentEntries) {
			updateStatus("列表选择已清空。")
			return
		}
		entry := currentEntries[index]
		count := len(list.SelectedIndices())
		if count > 1 {
			updateStatus(fmt.Sprintf("列表选择已更新：%s（已选 %d 项）", item.Texts[0], count))
			return
		}
		if entry.IsFolder() {
			updateStatus("准备打开目录: " + entry.Path())
			return
		}
		updateStatus(fmt.Sprintf("准备打开文件: %s (%d KB)", entry.Path(), entry.Size))
	})
	list.OnActivate(func(index int, _ widgets.ListViewItem) {
		if index < 0 || index >= len(currentEntries) {
			return
		}
		entry := currentEntries[index]
		if entry.IsFolder() {
			treeView.SetSelectedNode(entry.TreeNode)
			updateStatus("已打开目录: " + entry.Path())
			return
		}
		treeView.SetSelectedNode(entry.TreeNode)
		updateStatus(fmt.Sprintf("已打开文件: %s (%d KB)", entry.Path(), entry.Size))
	})
	list.OnRenameRequest(func(_ widgets.EventContext, index int, item widgets.ListViewItem) bool {
		if index < 0 || index >= len(currentEntries) {
			return false
		}
		entry, _ := item.Data.(*demoEntry)
		if entry == nil {
			entry = currentEntries[index]
		}
		if entry == nil {
			return false
		}
		updateStatus("开始重命名: " + entry.Path())
		return false
	})
	list.OnRenameCommit(func(index int, item widgets.ListViewItem, oldText, newText string) {
		if index < 0 || index >= len(currentEntries) {
			return
		}
		entry, _ := item.Data.(*demoEntry)
		if entry == nil {
			entry = currentEntries[index]
		}
		if entry == nil {
			return
		}
		entry.Name = newText
		if entry.TreeNode != nil {
			entry.TreeNode.Text = newText
		}
		refreshCurrentFolder(entry)
		treeView.SetSelectedNode(entry.TreeNode)
		updateStatus(fmt.Sprintf("已重命名: %s -> %s", oldText, newText))
	})
	list.OnColumnClick(func(index int) {
		nextSort := cmdSortByName
		switch index {
		case 1:
			nextSort = cmdSortBySize
		case 2:
			nextSort = cmdSortByType
		}
		if currentSort == nextSort {
			sortDescending = !sortDescending
			applySort(nextSort, "已切换排序方向。")
			return
		}
		sortDescending = false
		runCommand(nextSort)
	})

	tabs.OnSelectionChange(func(_ int, page *widgets.TabPage) {
		syncDefaultButton()
		updateStatus("当前页签: " + page.Title)
	})
	sortCombo.OnCommit(func(index int, _ string) {
		if index < 0 {
			updateStatus("请输入“按名称”“按大小”或“按类型”后再提交。")
			return
		}
		if index == 1 {
			runCommand(cmdSortBySize)
			return
		}
		if index == 2 {
			runCommand(cmdSortByType)
			return
		}
		runCommand(cmdSortByName)
	})
	addBtn.OnClick(func() { runCommand(cmdAddPath) })
	upBtn.OnClick(func() { runCommand(cmdNavigateUp) })
	closeBtn.OnClick(func() { runCommand(cmdExit) })

	appWin.Content().Add(title)
	appWin.Content().Add(toolbar)
	appWin.Content().Add(tabs)
	appWin.Content().Add(statusBar)
	syncDefaultButton()
	app.Desktop().AddWindow(appWin)

	currentFolder = root
	refreshCurrentFolder(nil)
	sortCombo.SetSelectedIndexSilent(0)
	applySelectionMode(false)
	treeView.SetSelectedNode(root.TreeNode)
	updateStatus("试试树/列表联动、Ctrl+N、Alt+Backspace、F2 和单选模式切换。")

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

func bindTree(entry *demoEntry, index map[*widgets.TreeNode]*demoEntry) {
	if entry == nil {
		return
	}
	if entry.IsFolder() {
		entry.TreeNode = widgets.NewFolderNode(entry.Name)
	} else {
		entry.TreeNode = widgets.NewFileNode(entry.Name)
	}
	index[entry.TreeNode] = entry
	for _, child := range entry.Children {
		child.Parent = entry
		bindTree(child, index)
		entry.TreeNode.AddChild(child.TreeNode)
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
