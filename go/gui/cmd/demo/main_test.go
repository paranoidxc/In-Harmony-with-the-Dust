package main

import (
	"testing"

	"classicui"
	"classicui/event"
	"classicui/widgets"
)

func TestBuildDemoRootBindsTreeAndParents(t *testing.T) {
	root, index := buildDemoRoot()
	if root == nil {
		t.Fatal("buildDemoRoot returned nil root")
	}
	if root.TreeNode == nil {
		t.Fatal("root tree node should be bound")
	}
	if index[root.TreeNode] != root {
		t.Fatal("tree index should map root node back to entry")
	}
	if !root.TreeNode.Expanded {
		t.Fatal("root node should start expanded")
	}
	if len(root.Children) < 2 {
		t.Fatalf("root child count = %d, want at least 2", len(root.Children))
	}
	if !root.Children[0].TreeNode.Expanded || !root.Children[1].TreeNode.Expanded {
		t.Fatal("first two top-level folders should start expanded for the demo")
	}

	docs := root.Children[1]
	projects := docs.Children[0]
	file := projects.Children[0]
	if docs.Parent != root {
		t.Fatal("folder parent link should be preserved")
	}
	if projects.TreeNode.Parent() != docs.TreeNode {
		t.Fatal("tree parent link should be preserved")
	}
	if index[file.TreeNode] != file {
		t.Fatal("tree index should include nested files")
	}
}

func TestSortDemoFolderKeepsModelAndTreeInSync(t *testing.T) {
	folder := newFolder("root",
		newFile("bbb.txt", 3),
		newFile("a.txt", 1),
		newFile("cc.txt", 2),
	)
	index := map[*widgets.TreeNode]*demoEntry{}
	bindTree(folder, index)

	sortDemoFolder(folder, cmdSortByName, false)
	assertDemoOrder(t, folder, []string{"a.txt", "bbb.txt", "cc.txt"})

	sortDemoFolder(folder, cmdSortBySize, false)
	assertDemoOrder(t, folder, []string{"a.txt", "cc.txt", "bbb.txt"})

	sortDemoFolder(folder, cmdSortBySize, true)
	assertDemoOrder(t, folder, []string{"bbb.txt", "cc.txt", "a.txt"})
}

func TestSortDemoFolderByTypeKeepsFoldersFirst(t *testing.T) {
	folder := newFolder("root",
		newFile("notes.md", 10),
		newFolder("assets"),
		newFile("data.bin", 20),
		newFile("readme.txt", 5),
	)
	index := map[*widgets.TreeNode]*demoEntry{}
	bindTree(folder, index)

	sortDemoFolder(folder, cmdSortByType, false)
	assertDemoOrder(t, folder, []string{"assets", "data.bin", "notes.md", "readme.txt"})

	sortDemoFolder(folder, cmdSortByType, true)
	assertDemoOrder(t, folder, []string{"assets", "readme.txt", "notes.md", "data.bin"})
}

func TestSortDemoFolderDescendingAlsoReordersFolderBlock(t *testing.T) {
	folder := newFolder("root",
		newFolder("alpha"),
		newFolder("beta"),
		newFolder("gamma"),
	)
	index := map[*widgets.TreeNode]*demoEntry{}
	bindTree(folder, index)

	sortDemoFolder(folder, cmdSortByName, false)
	assertDemoOrder(t, folder, []string{"alpha", "beta", "gamma"})

	sortDemoFolder(folder, cmdSortByName, true)
	assertDemoOrder(t, folder, []string{"gamma", "beta", "alpha"})
}

func TestAddEntryToFolderParsesLeafAndRejectsDuplicates(t *testing.T) {
	folder := newFolder("docs")
	index := map[*widgets.TreeNode]*demoEntry{}
	bindTree(folder, index)

	child, message := addEntryToFolder(folder, `C:\Temp\notes.txt\`, index)
	if message != "" {
		t.Fatalf("unexpected add error: %s", message)
	}
	if child == nil || child.Name != "notes.txt" || child.Kind != demoFile {
		t.Fatalf("added child = %#v, want file notes.txt", child)
	}
	if child.Parent != folder {
		t.Fatal("added child parent should point at folder")
	}
	if index[child.TreeNode] != child {
		t.Fatal("tree index should include added child")
	}
	if len(folder.TreeNode.Children) != 1 || folder.TreeNode.Children[0] != child.TreeNode {
		t.Fatal("tree node children should mirror folder children")
	}
	if !folder.TreeNode.Expanded {
		t.Fatal("adding a child should expand the parent folder")
	}

	if _, message := addEntryToFolder(folder, "NOTES.TXT", index); message == "" {
		t.Fatal("duplicate add should be rejected case-insensitively")
	}

	dir, message := addEntryToFolder(folder, "Projects", index)
	if message != "" {
		t.Fatalf("unexpected folder add error: %s", message)
	}
	if dir == nil || dir.Kind != demoFolder || dir.Name != "Projects" {
		t.Fatalf("added folder = %#v, want folder Projects", dir)
	}
}

func TestResolveDemoSelectionReturnsFolderAndSelectableFile(t *testing.T) {
	root, _ := buildDemoRoot()
	file := root.Children[len(root.Children)-1]
	folder, selectEntry := resolveDemoSelection(file)
	if folder != root {
		t.Fatal("file selection should resolve current folder to file parent")
	}
	if selectEntry != file {
		t.Fatal("file selection should keep the file as list selection target")
	}

	docs := root.Children[1]
	folder, selectEntry = resolveDemoSelection(docs)
	if folder != docs || selectEntry != nil {
		t.Fatal("folder selection should resolve to itself without list target")
	}
}

func TestFindDemoEntryIndexAndSelectionModeSync(t *testing.T) {
	root, _ := buildDemoRoot()
	entries := append([]*demoEntry(nil), root.Children...)
	target := root.Children[len(root.Children)-1]
	if index := findDemoEntryIndex(entries, target); index != len(entries)-1 {
		t.Fatalf("index = %d, want %d", index, len(entries)-1)
	}
	if index := findDemoEntryIndex(entries, newFile("missing.txt", 1)); index != -1 {
		t.Fatalf("missing entry index = %d, want -1", index)
	}

	list := widgets.NewListBox("files", classicui.Rect{})
	tree := widgets.NewTreeView("tree", classicui.Rect{}, root.TreeNode)
	applyDemoSelectionMode(list, tree, true)
	if list.SelectionOptions().MultiSelect {
		t.Fatal("list should switch to single-select mode")
	}
	if tree.SelectionOptions().MultiSelect {
		t.Fatal("tree should switch to single-select mode")
	}

	applyDemoSelectionMode(list, tree, false)
	if !list.SelectionOptions().MultiSelect {
		t.Fatal("list should switch back to multi-select mode")
	}
	if !tree.SelectionOptions().MultiSelect {
		t.Fatal("tree should switch back to multi-select mode")
	}
}

func TestDemoTypeTextAndListViewRows(t *testing.T) {
	folder := newFolder("docs")
	file := newFile("notes.txt", 12)
	bin := newFile("payload.bin", 64)
	entries := []*demoEntry{folder, file, bin}

	rows := demoListViewRows(entries)
	if len(rows) != 3 {
		t.Fatalf("row count = %d, want 3", len(rows))
	}
	if rows[0].Texts[2] != "文件夹" {
		t.Fatalf("folder type = %q, want 文件夹", rows[0].Texts[2])
	}
	if rows[1].Texts[1] != "12 KB" {
		t.Fatalf("size text = %q, want 12 KB", rows[1].Texts[1])
	}
	if rows[1].Texts[2] != "TXT 文件" {
		t.Fatalf("txt type = %q, want TXT 文件", rows[1].Texts[2])
	}
	if rows[2].Texts[2] != "BIN 文件" {
		t.Fatalf("bin type = %q, want BIN 文件", rows[2].Texts[2])
	}
}

func TestFirstDemoFileIndex(t *testing.T) {
	entries := []*demoEntry{
		newFolder("docs"),
		newFolder("assets"),
		newFile("notes.txt", 1),
	}
	if index := firstDemoFileIndex(entries); index != 2 {
		t.Fatalf("first file index = %d, want 2", index)
	}
	if index := firstDemoFileIndex([]*demoEntry{newFolder("only")}); index != 1 {
		t.Fatalf("first file index for all-folders = %d, want 1", index)
	}
}

func TestBuildDemoContextMenuBlankIncludesRefreshAndNoRename(t *testing.T) {
	sortByNameItem := widgets.NewMenuItem(cmdSortByName, "By &Name", nil)
	sortBySizeItem := widgets.NewMenuItem(cmdSortBySize, "By &Size", nil)
	sortByTypeItem := widgets.NewMenuItem(cmdSortByType, "By &Type", nil)
	renameItem := widgets.NewMenuItem(cmdFileRename, "&Rename", &widgets.Accelerator{Key: event.KeyF2})
	refreshItem := widgets.NewMenuItem(cmdViewRefresh, "&Refresh", nil)

	menu := buildDemoContextMenu(nil, newFolder("桌面"), renameItem, refreshItem, sortByNameItem, sortBySizeItem, sortByTypeItem)
	if menu == nil {
		t.Fatal("blank context menu should not be nil")
	}
	if !menuContainsCommand(menu, cmdViewRefresh) {
		t.Fatal("blank context menu should include refresh")
	}
	if !menuContainsCommand(menu, cmdFileNewFolder) || !menuContainsCommand(menu, cmdFileNewTextFile) {
		t.Fatal("blank context menu should include new-folder and new-text-file commands")
	}
	if menuContainsCommand(menu, cmdFileRename) {
		t.Fatal("blank context menu should not include rename")
	}
}

func TestBuildDemoContextMenuEntryIncludesRename(t *testing.T) {
	sortByNameItem := widgets.NewMenuItem(cmdSortByName, "By &Name", nil)
	sortBySizeItem := widgets.NewMenuItem(cmdSortBySize, "By &Size", nil)
	sortByTypeItem := widgets.NewMenuItem(cmdSortByType, "By &Type", nil)
	renameItem := widgets.NewMenuItem(cmdFileRename, "&Rename", &widgets.Accelerator{Key: event.KeyF2})
	refreshItem := widgets.NewMenuItem(cmdViewRefresh, "&Refresh", nil)
	entry := newFile("readme.txt", 8)

	menu := buildDemoContextMenu([]*demoEntry{entry}, newFolder("桌面"), renameItem, refreshItem, sortByNameItem, sortBySizeItem, sortByTypeItem)
	if menu == nil {
		t.Fatal("entry context menu should not be nil")
	}
	if !menuContainsCommand(menu, cmdFileRename) {
		t.Fatal("entry context menu should include rename")
	}
	if !menuContainsCommand(menu, cmdSelectionInfo) {
		t.Fatal("entry context menu should include properties command")
	}
	if !menuContainsCommand(menu, cmdViewRefresh) {
		t.Fatal("entry context menu should include refresh")
	}
}

func TestBuildDemoContextMenuMultiSelectionEnablesOpenAndDisablesRename(t *testing.T) {
	sortByNameItem := widgets.NewMenuItem(cmdSortByName, "By &Name", nil)
	sortBySizeItem := widgets.NewMenuItem(cmdSortBySize, "By &Size", nil)
	sortByTypeItem := widgets.NewMenuItem(cmdSortByType, "By &Type", nil)
	renameItem := widgets.NewMenuItem(cmdFileRename, "&Rename", &widgets.Accelerator{Key: event.KeyF2})
	refreshItem := widgets.NewMenuItem(cmdViewRefresh, "&Refresh", nil)
	selection := []*demoEntry{
		newFile("readme.txt", 8),
		newFolder("Docs"),
	}

	menu := buildDemoContextMenu(selection, newFolder("桌面"), renameItem, refreshItem, sortByNameItem, sortBySizeItem, sortByTypeItem)
	if menu == nil {
		t.Fatal("multi-selection context menu should not be nil")
	}
	openItem := findMenuItem(menu, cmdOpenSelection)
	if openItem == nil {
		t.Fatal("multi-selection context menu should include open item")
	}
	if !openItem.Enabled {
		t.Fatal("multi-selection open item should be enabled")
	}
	rename := findMenuItem(menu, cmdFileRename)
	if rename == nil {
		t.Fatal("multi-selection context menu should include rename item")
	}
	if rename.Enabled {
		t.Fatal("multi-selection rename item should be disabled")
	}
	if !menuContainsCommand(menu, cmdSelectionInfo) {
		t.Fatal("multi-selection context menu should include properties command")
	}
}

func TestBuildDemoContextMenuCurrentFolderAddsParentCommand(t *testing.T) {
	sortByNameItem := widgets.NewMenuItem(cmdSortByName, "By &Name", nil)
	sortBySizeItem := widgets.NewMenuItem(cmdSortBySize, "By &Size", nil)
	sortByTypeItem := widgets.NewMenuItem(cmdSortByType, "By &Type", nil)
	renameItem := widgets.NewMenuItem(cmdFileRename, "&Rename", &widgets.Accelerator{Key: event.KeyF2})
	refreshItem := widgets.NewMenuItem(cmdViewRefresh, "&Refresh", nil)
	root := newFolder("桌面")
	current := newFolder("文档")
	root.AddChild(current)

	menu := buildDemoContextMenu([]*demoEntry{current}, current, renameItem, refreshItem, sortByNameItem, sortBySizeItem, sortByTypeItem)
	if menu == nil {
		t.Fatal("current-folder context menu should not be nil")
	}
	if !menuContainsCommand(menu, cmdNavigateUp) {
		t.Fatal("current-folder context menu should include navigate-up command")
	}
	if menuContainsCommand(menu, cmdOpenSelection) {
		t.Fatal("current-folder context menu should not include open-selection command")
	}
}

func menuContainsCommand(menu *widgets.Menu, cmd classicui.CommandID) bool {
	return findMenuItem(menu, cmd) != nil
}

func findMenuItem(menu *widgets.Menu, cmd classicui.CommandID) *widgets.MenuItem {
	if menu == nil {
		return nil
	}
	for _, item := range menu.Items {
		if item == nil {
			continue
		}
		if item.ID == widgets.CommandID(cmd) {
			return item
		}
		if found := findMenuItem(item.Submenu, cmd); found != nil {
			return found
		}
	}
	return nil
}

func assertDemoOrder(t *testing.T, folder *demoEntry, want []string) {
	t.Helper()
	if len(folder.Children) != len(want) {
		t.Fatalf("child count = %d, want %d", len(folder.Children), len(want))
	}
	if folder.TreeNode == nil {
		t.Fatal("folder tree node should not be nil")
	}
	if len(folder.TreeNode.Children) != len(want) {
		t.Fatalf("tree child count = %d, want %d", len(folder.TreeNode.Children), len(want))
	}
	for i, name := range want {
		if folder.Children[i].Name != name {
			t.Fatalf("model child %d = %q, want %q", i, folder.Children[i].Name, name)
		}
		if folder.TreeNode.Children[i].Text != name {
			t.Fatalf("tree child %d = %q, want %q", i, folder.TreeNode.Children[i].Text, name)
		}
	}
}
