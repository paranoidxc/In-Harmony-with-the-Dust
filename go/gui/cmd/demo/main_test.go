package main

import (
	"testing"

	"classicui"
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
		newFile("bbb.txt", 1),
		newFile("a.txt", 1),
		newFile("cc.txt", 1),
	)
	index := map[*widgets.TreeNode]*demoEntry{}
	bindTree(folder, index)

	sortDemoFolder(folder, cmdSortByName)
	assertDemoOrder(t, folder, []string{"a.txt", "bbb.txt", "cc.txt"})

	sortDemoFolder(folder, cmdSortByLength)
	assertDemoOrder(t, folder, []string{"a.txt", "cc.txt", "bbb.txt"})
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
