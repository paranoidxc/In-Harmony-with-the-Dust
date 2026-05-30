package widgets

import (
	"testing"
	"time"

	"classicui/event"
	"classicui/geom"
	"classicui/paint"
	"classicui/theme"
)

func TestTreeViewKeyboardNavigationAndCollapse(t *testing.T) {
	root := NewTreeNode("Root",
		NewTreeNode("Child 1"),
		NewTreeNode("Child 2"),
	)
	root.Expanded = true
	tree := NewTreeView("tree", geom.Rect{X: 0, Y: 0, W: 180, H: 120}, root)
	ctx := &fakeContext{}

	if tree.SelectedNode() != root {
		t.Fatal("tree should select the first visible node by default")
	}

	if !tree.KeyDown(ctx, event.KeyEvent{Key: event.KeyRight}) {
		t.Fatal("right key should be handled")
	}
	if got := tree.SelectedNode(); got == nil || got.Text != "Child 1" {
		t.Fatalf("selected after right = %#v, want Child 1", got)
	}

	tree.KeyDown(ctx, event.KeyEvent{Key: event.KeyLeft})
	if tree.SelectedNode() != root {
		t.Fatal("left key should move selection back to parent")
	}

	tree.KeyDown(ctx, event.KeyEvent{Key: event.KeyLeft})
	if root.Expanded {
		t.Fatal("left key on expanded node should collapse it")
	}
	if len(tree.visibleNodes()) != 1 {
		t.Fatalf("visible node count after collapse = %d, want 1", len(tree.visibleNodes()))
	}
}

func TestTreeViewMouseToggleExpanderAndSelection(t *testing.T) {
	root := NewTreeNode("Root", NewTreeNode("Child"))
	tree := NewTreeView("tree", geom.Rect{X: 0, Y: 0, W: 180, H: 120}, root)
	ctx := &fakeContext{}

	rootRect := geom.Rect{X: 3, Y: 2, W: 160, H: tree.rowHeight}
	expander := tree.expanderRect(rootRect, 0)
	click := geom.Point{X: expander.X + 1, Y: expander.Y + 1}
	tree.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, click)

	if !root.Expanded {
		t.Fatal("clicking expander should expand the node")
	}

	rowClick := geom.Point{X: expander.Right() + 12, Y: expander.Y + 1}
	tree.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, rowClick)
	if tree.SelectedNode() != root {
		t.Fatal("clicking a row should select that node")
	}
}

func TestTreeViewCtrlClickTogglesMultiSelection(t *testing.T) {
	root := NewFolderNode("Root",
		NewFileNode("Child 1"),
		NewFileNode("Child 2"),
	)
	root.Expanded = true
	tree := NewTreeView("tree", geom.Rect{X: 0, Y: 0, W: 180, H: 120}, root)
	ctx := &fakeContext{}

	child1Rect := geom.Rect{X: 3, Y: 18, W: 160, H: tree.rowHeight}
	child1Click := geom.Point{X: tree.iconRect(child1Rect, 1).Right() + 6, Y: child1Rect.Y + child1Rect.H/2}
	tree.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft, Modifiers: event.ModCtrl}, child1Click)

	child2Rect := geom.Rect{X: 3, Y: 34, W: 160, H: tree.rowHeight}
	child2Click := geom.Point{X: tree.iconRect(child2Rect, 1).Right() + 6, Y: child2Rect.Y + child2Rect.H/2}
	tree.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft, Modifiers: event.ModCtrl}, child2Click)

	selected := tree.SelectedNodes()
	if len(selected) != 3 {
		t.Fatalf("selected nodes = %d, want 3", len(selected))
	}
	if selected[1].Text != "Child 1" || selected[2].Text != "Child 2" {
		t.Fatalf("selected order = [%s %s %s], want Root/Child 1/Child 2", selected[0].Text, selected[1].Text, selected[2].Text)
	}
	if tree.SelectedNode() == nil || tree.SelectedNode().Text != "Child 2" {
		t.Fatalf("lead selection = %#v, want Child 2", tree.SelectedNode())
	}
}

func TestTreeViewSingleSelectModeIgnoresMultiSelectModifiers(t *testing.T) {
	root := NewFolderNode("Root",
		NewFileNode("Child 1"),
		NewFileNode("Child 2"),
	)
	root.Expanded = true
	tree := NewTreeView("tree", geom.Rect{X: 0, Y: 0, W: 180, H: 120}, root)
	tree.SetMultiSelect(false)
	ctx := &fakeContext{}

	child1Rect := geom.Rect{X: 3, Y: 18, W: 160, H: tree.rowHeight}
	child1Click := geom.Point{X: tree.iconRect(child1Rect, 1).Right() + 6, Y: child1Rect.Y + child1Rect.H/2}
	tree.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft, Modifiers: event.ModCtrl}, child1Click)

	child2Rect := geom.Rect{X: 3, Y: 34, W: 160, H: tree.rowHeight}
	child2Click := geom.Point{X: tree.iconRect(child2Rect, 1).Right() + 6, Y: child2Rect.Y + child2Rect.H/2}
	tree.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft, Modifiers: event.ModShift}, child2Click)

	selected := tree.SelectedNodes()
	if len(selected) != 1 || selected[0].Text != "Child 2" {
		t.Fatalf("selected nodes = %#v, want only Child 2", selected)
	}
	if !tree.KeyDown(ctx, event.KeyEvent{Key: event.KeyA, Modifiers: event.ModCtrl}) {
		t.Fatal("ctrl+a should still be handled in single-select mode")
	}
	selected = tree.SelectedNodes()
	if len(selected) != 1 || selected[0].Text != "Child 2" {
		t.Fatalf("selected nodes after ctrl+a = %#v, want still only Child 2", selected)
	}
}

func TestTreeViewShiftClickSelectsVisibleRange(t *testing.T) {
	root := NewFolderNode("Root",
		NewFileNode("Child 1"),
		NewFileNode("Child 2"),
		NewFileNode("Child 3"),
	)
	root.Expanded = true
	tree := NewTreeView("tree", geom.Rect{X: 0, Y: 0, W: 180, H: 120}, root)
	ctx := &fakeContext{}

	firstRect := geom.Rect{X: 3, Y: 18, W: 160, H: tree.rowHeight}
	firstClick := geom.Point{X: tree.iconRect(firstRect, 1).Right() + 6, Y: firstRect.Y + firstRect.H/2}
	tree.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, firstClick)

	lastRect := geom.Rect{X: 3, Y: 50, W: 160, H: tree.rowHeight}
	lastClick := geom.Point{X: tree.iconRect(lastRect, 1).Right() + 6, Y: lastRect.Y + lastRect.H/2}
	tree.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft, Modifiers: event.ModShift}, lastClick)

	selected := tree.SelectedNodes()
	if len(selected) != 3 {
		t.Fatalf("selected nodes = %d, want 3", len(selected))
	}
	if selected[0].Text != "Child 1" || selected[1].Text != "Child 2" || selected[2].Text != "Child 3" {
		t.Fatalf("selected order = [%s %s %s], want Child 1..3", selected[0].Text, selected[1].Text, selected[2].Text)
	}
}

func TestTreeViewCollapseMovesSelectionToCollapsedAncestor(t *testing.T) {
	child := NewTreeNode("Child")
	root := NewTreeNode("Root", child)
	root.Expanded = true
	tree := NewTreeView("tree", geom.Rect{X: 0, Y: 0, W: 180, H: 120}, root)
	tree.SetSelectedNode(child)
	ctx := &fakeContext{}

	tree.KeyDown(ctx, event.KeyEvent{Key: event.KeyLeft})
	tree.KeyDown(ctx, event.KeyEvent{Key: event.KeyLeft})

	if tree.SelectedNode() != root {
		t.Fatal("collapsing the parent should move selection back to that parent")
	}
}

func TestTreeViewSelectingHiddenNodeExpandsAncestors(t *testing.T) {
	child := NewTreeNode("Child")
	root := NewTreeNode("Root", child)
	tree := NewTreeView("tree", geom.Rect{X: 0, Y: 0, W: 180, H: 120}, root)

	if !tree.SetSelectedNode(child) {
		t.Fatal("selecting hidden descendant should succeed")
	}
	if !root.Expanded {
		t.Fatal("selecting hidden descendant should expand its ancestors")
	}
}

func TestTreeViewKeyboardShiftExtendsSelection(t *testing.T) {
	root := NewFolderNode("Root",
		NewFileNode("Child 1"),
		NewFileNode("Child 2"),
		NewFileNode("Child 3"),
	)
	root.Expanded = true
	tree := NewTreeView("tree", geom.Rect{X: 0, Y: 0, W: 180, H: 120}, root)
	ctx := &fakeContext{}

	tree.SetSelectedNode(root.Children[0])
	if !tree.KeyDown(ctx, event.KeyEvent{Key: event.KeyDown, Modifiers: event.ModShift}) {
		t.Fatal("shift+down should be handled")
	}
	if !tree.KeyDown(ctx, event.KeyEvent{Key: event.KeyDown, Modifiers: event.ModShift}) {
		t.Fatal("second shift+down should be handled")
	}

	selected := tree.SelectedNodes()
	if len(selected) != 3 {
		t.Fatalf("selected nodes = %d, want 3", len(selected))
	}
	if selected[0].Text != "Child 1" || selected[2].Text != "Child 3" {
		t.Fatalf("selected range = [%s ... %s], want Child 1..3", selected[0].Text, selected[2].Text)
	}
}

func TestTreeViewCtrlSpaceTogglesLeadSelection(t *testing.T) {
	root := NewFolderNode("Root",
		NewFileNode("Child 1"),
		NewFileNode("Child 2"),
	)
	root.Expanded = true
	tree := NewTreeView("tree", geom.Rect{X: 0, Y: 0, W: 180, H: 120}, root)
	ctx := &fakeContext{}

	tree.toggleSelection(root.Children[0])
	if !tree.KeyDown(ctx, event.KeyEvent{Key: event.KeySpace, Modifiers: event.ModCtrl}) {
		t.Fatal("ctrl+space should be handled")
	}
	selected := tree.SelectedNodes()
	if len(selected) != 1 || selected[0] != root {
		t.Fatalf("selected after ctrl+space = %#v, want only root", selected)
	}
}

func TestTreeViewCtrlASelectsAllVisibleNodes(t *testing.T) {
	root := NewFolderNode("Root",
		NewFolderNode("Collapsed", NewFileNode("Hidden.txt")),
		NewFileNode("Visible.txt"),
	)
	root.Expanded = true
	tree := NewTreeView("tree", geom.Rect{X: 0, Y: 0, W: 180, H: 120}, root)
	ctx := &fakeContext{}

	if !tree.KeyDown(ctx, event.KeyEvent{Key: event.KeyA, Modifiers: event.ModCtrl}) {
		t.Fatal("ctrl+a should be handled")
	}

	selected := tree.SelectedNodes()
	if len(selected) != 3 {
		t.Fatalf("selected nodes = %d, want 3 visible nodes", len(selected))
	}
	if selected[0].Text != "Root" || selected[1].Text != "Collapsed" || selected[2].Text != "Visible.txt" {
		t.Fatalf("selected order = [%s %s %s], want only visible nodes", selected[0].Text, selected[1].Text, selected[2].Text)
	}
}

func TestTreeViewBlankClickClearsSelectionAndFocuses(t *testing.T) {
	root := NewFolderNode("Root",
		NewFileNode("Child 1"),
		NewFileNode("Child 2"),
	)
	root.Expanded = true
	tree := NewTreeView("tree", geom.Rect{X: 0, Y: 0, W: 180, H: 120}, root)
	ctx := &fakeContext{}

	tree.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, geom.Point{X: 24, Y: 72})

	if ctx.focused != tree {
		t.Fatal("blank click should focus the tree")
	}
	if tree.SelectedNode() != nil {
		t.Fatalf("selected node = %#v, want nil", tree.SelectedNode())
	}
	if selected := tree.SelectedNodes(); len(selected) != 0 {
		t.Fatalf("selected nodes = %#v, want none", selected)
	}
}

func TestTreeViewCtrlBlankClickPreservesSelection(t *testing.T) {
	root := NewFolderNode("Root",
		NewFileNode("Child 1"),
		NewFileNode("Child 2"),
	)
	root.Expanded = true
	tree := NewTreeView("tree", geom.Rect{X: 0, Y: 0, W: 180, H: 120}, root)
	ctx := &fakeContext{}

	childRect := geom.Rect{X: 3, Y: 18, W: 160, H: tree.rowHeight}
	childClick := geom.Point{X: tree.iconRect(childRect, 1).Right() + 6, Y: childRect.Y + childRect.H/2}
	tree.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, childClick)
	tree.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft, Modifiers: event.ModCtrl}, geom.Point{X: 24, Y: 72})

	if ctx.focused != tree {
		t.Fatal("ctrl+blank click should focus the tree")
	}
	if tree.SelectedNode() == nil || tree.SelectedNode().Text != "Child 1" {
		t.Fatalf("selected node = %#v, want Child 1", tree.SelectedNode())
	}
	selected := tree.SelectedNodes()
	if len(selected) != 1 || selected[0].Text != "Child 1" {
		t.Fatalf("selected nodes = %#v, want only Child 1", selected)
	}
}

func TestTreeViewKeyboardRecoversRecentSelectionAfterBlankClear(t *testing.T) {
	root := NewFolderNode("Root",
		NewFileNode("Child 1"),
		NewFileNode("Child 2"),
		NewFileNode("Child 3"),
	)
	root.Expanded = true
	tree := NewTreeView("tree", geom.Rect{X: 0, Y: 0, W: 180, H: 120}, root)
	ctx := &fakeContext{}

	childRect := geom.Rect{X: 3, Y: 34, W: 160, H: tree.rowHeight}
	childClick := geom.Point{X: tree.iconRect(childRect, 1).Right() + 6, Y: childRect.Y + childRect.H/2}
	tree.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, childClick)
	tree.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, geom.Point{X: 24, Y: 96})

	if !tree.KeyDown(ctx, event.KeyEvent{Key: event.KeyDown}) {
		t.Fatal("down should recover from empty selection")
	}
	if tree.SelectedNode() == nil || tree.SelectedNode().Text != "Child 2" {
		t.Fatalf("selected node after recovery = %#v, want Child 2", tree.SelectedNode())
	}
}

func TestTreeViewCtrlSpaceRecoversRecentSelectionAfterBlankClear(t *testing.T) {
	root := NewFolderNode("Root",
		NewFileNode("Child 1"),
		NewFileNode("Child 2"),
		NewFileNode("Child 3"),
	)
	root.Expanded = true
	tree := NewTreeView("tree", geom.Rect{X: 0, Y: 0, W: 180, H: 120}, root)
	ctx := &fakeContext{}

	childRect := geom.Rect{X: 3, Y: 50, W: 160, H: tree.rowHeight}
	childClick := geom.Point{X: tree.iconRect(childRect, 1).Right() + 6, Y: childRect.Y + childRect.H/2}
	tree.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, childClick)
	tree.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, geom.Point{X: 24, Y: 96})

	if !tree.KeyDown(ctx, event.KeyEvent{Key: event.KeySpace, Modifiers: event.ModCtrl}) {
		t.Fatal("ctrl+space should recover from empty selection")
	}
	if tree.SelectedNode() == nil || tree.SelectedNode().Text != "Child 3" {
		t.Fatalf("selected node after ctrl+space recovery = %#v, want Child 3", tree.SelectedNode())
	}
}

func TestTreeViewHomeEndOverrideRecentSelectionRecovery(t *testing.T) {
	root := NewFolderNode("Root",
		NewFileNode("Child 1"),
		NewFileNode("Child 2"),
		NewFileNode("Child 3"),
	)
	root.Expanded = true
	tree := NewTreeView("tree", geom.Rect{X: 0, Y: 0, W: 180, H: 120}, root)
	ctx := &fakeContext{}

	childRect := geom.Rect{X: 3, Y: 34, W: 160, H: tree.rowHeight}
	childClick := geom.Point{X: tree.iconRect(childRect, 1).Right() + 6, Y: childRect.Y + childRect.H/2}
	tree.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, childClick)
	tree.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, geom.Point{X: 24, Y: 96})

	if !tree.KeyDown(ctx, event.KeyEvent{Key: event.KeyEnd}) {
		t.Fatal("end should recover to last visible node")
	}
	if tree.SelectedNode() == nil || tree.SelectedNode().Text != "Child 3" {
		t.Fatalf("selected node after end = %#v, want Child 3", tree.SelectedNode())
	}

	tree.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, geom.Point{X: 24, Y: 96})
	if !tree.KeyDown(ctx, event.KeyEvent{Key: event.KeyHome}) {
		t.Fatal("home should recover to first visible node")
	}
	if tree.SelectedNode() == nil || tree.SelectedNode().Text != "Root" {
		t.Fatalf("selected node after home = %#v, want Root", tree.SelectedNode())
	}
}

func TestTreeViewSelectionOptionsCanDisableRecentRecovery(t *testing.T) {
	root := NewFolderNode("Root",
		NewFileNode("Child 1"),
		NewFileNode("Child 2"),
		NewFileNode("Child 3"),
	)
	root.Expanded = true
	tree := NewTreeView("tree", geom.Rect{X: 0, Y: 0, W: 180, H: 120}, root)
	options := tree.SelectionOptions()
	options.RecoverFromRecent = false
	tree.SetSelectionOptions(options)
	ctx := &fakeContext{}

	childRect := geom.Rect{X: 3, Y: 34, W: 160, H: tree.rowHeight}
	childClick := geom.Point{X: tree.iconRect(childRect, 1).Right() + 6, Y: childRect.Y + childRect.H/2}
	tree.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, childClick)
	tree.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, geom.Point{X: 24, Y: 96})

	if !tree.KeyDown(ctx, event.KeyEvent{Key: event.KeyDown}) {
		t.Fatal("down should recover from empty selection")
	}
	if tree.SelectedNode() == nil || tree.SelectedNode().Text != "Root" {
		t.Fatalf("selected node with recent recovery disabled = %#v, want Root", tree.SelectedNode())
	}
}

func TestTreeViewSelectionOptionsCanDisableBlankDragSelect(t *testing.T) {
	root := NewFolderNode("Root",
		NewFileNode("Child 1"),
		NewFileNode("Child 2"),
		NewFileNode("Child 3"),
	)
	root.Expanded = true
	tree := NewTreeView("tree", geom.Rect{X: 0, Y: 0, W: 180, H: 120}, root)
	options := tree.SelectionOptions()
	options.BlankDragSelect = false
	tree.SetSelectionOptions(options)
	ctx := &fakeContext{}

	tree.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, geom.Point{X: 24, Y: 88})
	tree.MouseMove(ctx, geom.Point{X: 28, Y: 20})

	if !tree.dragSelecting {
		t.Fatal("blank drag should still cross the drag threshold")
	}
	if _, ok := tree.marqueeRect(); ok {
		t.Fatal("blank drag marquee should be disabled")
	}
	if selected := tree.SelectedNodes(); len(selected) != 0 {
		t.Fatalf("selected nodes after disabled blank drag = %#v, want none", selected)
	}
}

func TestTreeViewDragSelectionStartsAfterThreshold(t *testing.T) {
	root := NewFolderNode("Root",
		NewFileNode("Child 1"),
		NewFileNode("Child 2"),
		NewFileNode("Child 3"),
	)
	root.Expanded = true
	tree := NewTreeView("tree", geom.Rect{X: 0, Y: 0, W: 180, H: 120}, root)
	ctx := &fakeContext{}

	start := geom.Point{X: 24, Y: 24}
	tree.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, start)
	tree.MouseMove(ctx, geom.Point{X: 25, Y: 25})
	if tree.dragSelecting {
		t.Fatal("small move should not start drag selection")
	}
	tree.MouseMove(ctx, geom.Point{X: 24, Y: 56})
	if !tree.dragSelecting {
		t.Fatal("move beyond threshold should start drag selection")
	}
	selected := tree.SelectedNodes()
	if len(selected) != 3 || selected[0].Text != "Child 1" || selected[2].Text != "Child 3" {
		t.Fatalf("selected nodes after drag = %#v, want Child 1..3", selected)
	}
}

func TestTreeViewBlankDragSelectsIntersectingRows(t *testing.T) {
	root := NewFolderNode("Root",
		NewFileNode("Child 1"),
		NewFileNode("Child 2"),
		NewFileNode("Child 3"),
	)
	root.Expanded = true
	tree := NewTreeView("tree", geom.Rect{X: 0, Y: 0, W: 180, H: 120}, root)
	ctx := &fakeContext{}

	tree.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, geom.Point{X: 24, Y: 88})
	tree.MouseMove(ctx, geom.Point{X: 28, Y: 20})

	if !tree.dragSelecting {
		t.Fatal("blank drag should enter drag-select mode")
	}
	if _, ok := tree.marqueeRect(); !ok {
		t.Fatal("blank drag should expose a marquee rect")
	}
	selected := tree.SelectedNodes()
	if len(selected) != 3 || selected[0].Text != "Child 1" || selected[1].Text != "Child 2" || selected[2].Text != "Child 3" {
		t.Fatalf("selected nodes after blank drag = %#v, want Child 1..3", selected)
	}
	if tree.SelectedNode() == nil || tree.SelectedNode().Text != "Child 3" {
		t.Fatalf("lead selection = %#v, want Child 3", tree.SelectedNode())
	}
}

func TestTreeViewCtrlBlankDragUnionsSelection(t *testing.T) {
	root := NewFolderNode("Root",
		NewFileNode("Child 1"),
		NewFileNode("Child 2"),
		NewFileNode("Child 3"),
	)
	root.Expanded = true
	tree := NewTreeView("tree", geom.Rect{X: 0, Y: 0, W: 180, H: 120}, root)
	ctx := &fakeContext{}

	tree.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft, Modifiers: event.ModCtrl}, geom.Point{X: 24, Y: 88})
	tree.MouseMove(ctx, geom.Point{X: 28, Y: 20})

	selected := tree.SelectedNodes()
	if len(selected) != 4 || selected[0].Text != "Root" || selected[3].Text != "Child 3" {
		t.Fatalf("selected nodes after ctrl+blank drag = %#v, want Root plus Child 1..3", selected)
	}
}

func TestTreeViewDragAutoScroll(t *testing.T) {
	root := NewFolderNode("Root",
		NewFileNode("Child 1"),
		NewFileNode("Child 2"),
		NewFileNode("Child 3"),
		NewFileNode("Child 4"),
		NewFileNode("Child 5"),
	)
	root.Expanded = true
	tree := NewTreeView("tree", geom.Rect{X: 0, Y: 0, W: 180, H: 52}, root)
	ctx := &fakeContext{}

	tree.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, geom.Point{X: 24, Y: 8})
	tree.MouseMove(ctx, geom.Point{X: 24, Y: 70})
	if !tree.dragSelecting {
		t.Fatal("drag should be active")
	}
	if !tree.Tick(ctx, time.Date(2026, time.May, 30, 10, 0, 1, 0, time.UTC)) {
		t.Fatal("tick should auto-scroll while dragging past bottom")
	}
	if tree.topIndex == 0 {
		t.Fatal("auto-scroll should advance top index")
	}
}

func TestTreeViewMouseDoubleClickTogglesExpanded(t *testing.T) {
	root := NewFolderNode("Root", NewFileNode("Readme.txt"))
	tree := NewTreeView("tree", geom.Rect{X: 0, Y: 0, W: 180, H: 120}, root)
	ctx := &fakeContext{}

	now := time.Date(2026, time.May, 29, 12, 0, 0, 0, time.UTC)
	tree.now = func() time.Time {
		return now
	}

	rootRect := geom.Rect{X: 3, Y: 2, W: 160, H: tree.rowHeight}
	click := geom.Point{X: tree.iconRect(rootRect, 0).Right() + 6, Y: rootRect.Y + rootRect.H/2}

	tree.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, click)
	now = now.Add(200 * time.Millisecond)
	tree.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, click)

	if !root.Expanded {
		t.Fatal("double-clicking a branch row should expand it")
	}

	now = now.Add(700 * time.Millisecond)
	tree.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, click)
	now = now.Add(200 * time.Millisecond)
	tree.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, click)

	if root.Expanded {
		t.Fatal("double-clicking an expanded branch row should collapse it")
	}
	if ctx.focused != tree {
		t.Fatal("tree should request focus on mouse interaction")
	}
}

func TestTreeNodeEffectiveKind(t *testing.T) {
	if got := NewTreeNode("folder", NewTreeNode("child")).EffectiveKind(); got != TreeNodeFolder {
		t.Fatalf("kind for node with children = %v, want folder", got)
	}
	if got := NewTreeNode("file.txt").EffectiveKind(); got != TreeNodeFile {
		t.Fatalf("kind for leaf node = %v, want file", got)
	}
	if got := NewFolderNode("empty").EffectiveKind(); got != TreeNodeFolder {
		t.Fatalf("kind for explicit folder = %v, want folder", got)
	}
}

func TestTreeViewVisibleNodesCarryGuideState(t *testing.T) {
	leaf := NewFileNode("main.go")
	folderA := NewFolderNode("src", leaf)
	folderA.Expanded = true
	folderB := NewFolderNode("assets")
	root := NewFolderNode("project", folderA, folderB)
	root.Expanded = true

	tree := NewTreeView("tree", geom.Rect{X: 0, Y: 0, W: 180, H: 120}, root)
	visible := tree.visibleNodes()
	if len(visible) != 4 {
		t.Fatalf("visible nodes = %d, want 4", len(visible))
	}

	entry := visible[2]
	if entry.node != leaf {
		t.Fatalf("visible[2] = %q, want %q", entry.node.Text, leaf.Text)
	}
	if entry.depth != 2 {
		t.Fatalf("leaf depth = %d, want 2", entry.depth)
	}
	if len(entry.guides) != 2 || entry.guides[0] || !entry.guides[1] {
		t.Fatalf("leaf guides = %#v, want [false true]", entry.guides)
	}
}

func TestTreeViewMouseMoveTracksHotPart(t *testing.T) {
	root := NewFolderNode("Root", NewFileNode("Child"))
	tree := NewTreeView("tree", geom.Rect{X: 0, Y: 0, W: 180, H: 120}, root)
	ctx := &fakeContext{}

	rowRect := geom.Rect{X: 3, Y: 2, W: 160, H: tree.rowHeight}
	expander := tree.expanderRect(rowRect, 0)
	tree.MouseMove(ctx, geom.Point{X: expander.X + 1, Y: expander.Y + 1})
	if tree.hotNode != root || tree.hotPart != treeHotPartExpander {
		t.Fatalf("hot state = (%v, %v), want expander hot on root", tree.hotNode, tree.hotPart)
	}

	icon := tree.iconRect(rowRect, 0)
	tree.MouseMove(ctx, geom.Point{X: icon.Right() + 10, Y: rowRect.Y + rowRect.H/2})
	if tree.hotNode != root || tree.hotPart != treeHotPartRow {
		t.Fatalf("hot state = (%v, %v), want row hot on root", tree.hotNode, tree.hotPart)
	}

	tree.MouseLeave(ctx)
	if tree.hotNode != nil || tree.hotPart != treeHotPartNone {
		t.Fatal("mouse leave should clear hot state")
	}
}

func TestTreeViewF2BeginsRename(t *testing.T) {
	root := NewFileNode("Readme.txt")
	tree := NewTreeView("tree", geom.Rect{X: 0, Y: 0, W: 180, H: 120}, root)
	ctx := &fakeContext{}

	var renamed *TreeNode
	tree.OnBeginRename(func(node *TreeNode) {
		renamed = node
	})

	if !tree.KeyDown(ctx, event.KeyEvent{Key: event.KeyF2}) {
		t.Fatal("F2 should be handled by treeview")
	}
	if renamed != root {
		t.Fatalf("renamed node = %#v, want root", renamed)
	}
}

func TestTreeViewDelayedRenameAfterClickOnSelectedNode(t *testing.T) {
	root := NewFileNode("Readme.txt")
	tree := NewTreeView("tree", geom.Rect{X: 0, Y: 0, W: 180, H: 120}, root)
	ctx := &fakeContext{}

	now := time.Date(2026, time.May, 29, 12, 0, 0, 0, time.UTC)
	tree.now = func() time.Time {
		return now
	}

	var renamed *TreeNode
	tree.OnBeginRename(func(node *TreeNode) {
		renamed = node
	})

	rowRect := geom.Rect{X: 3, Y: 2, W: 160, H: tree.rowHeight}
	click := geom.Point{X: tree.iconRect(rowRect, 0).Right() + 6, Y: rowRect.Y + rowRect.H/2}

	tree.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, click)
	tree.MouseUp(ctx, event.MouseButtonEvent{Button: event.MouseButtonLeft}, click)
	if tree.renameNode != root {
		t.Fatal("clicking an already selected node should arm delayed rename")
	}

	tree.Tick(ctx, now.Add(300*time.Millisecond))
	if renamed != nil {
		t.Fatal("rename should not fire before delay elapses")
	}

	tree.Tick(ctx, now.Add(600*time.Millisecond))
	if renamed != root {
		t.Fatalf("renamed node = %#v, want root after delay", renamed)
	}
}

func TestTreeViewDoubleClickCancelsDelayedRename(t *testing.T) {
	root := NewFolderNode("Root", NewFileNode("Child.txt"))
	tree := NewTreeView("tree", geom.Rect{X: 0, Y: 0, W: 180, H: 120}, root)
	ctx := &fakeContext{}

	now := time.Date(2026, time.May, 29, 12, 0, 0, 0, time.UTC)
	tree.now = func() time.Time {
		return now
	}

	renamed := false
	tree.OnBeginRename(func(*TreeNode) {
		renamed = true
	})

	rowRect := geom.Rect{X: 3, Y: 2, W: 160, H: tree.rowHeight}
	click := geom.Point{X: tree.iconRect(rowRect, 0).Right() + 6, Y: rowRect.Y + rowRect.H/2}

	tree.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, click)
	tree.MouseUp(ctx, event.MouseButtonEvent{Button: event.MouseButtonLeft}, click)
	if tree.renameNode != root {
		t.Fatal("first click should arm delayed rename")
	}

	now = now.Add(200 * time.Millisecond)
	tree.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, click)
	tree.MouseUp(ctx, event.MouseButtonEvent{Button: event.MouseButtonLeft}, click)

	if tree.renameNode != nil {
		t.Fatal("double-click should cancel delayed rename")
	}
	tree.Tick(ctx, now.Add(600*time.Millisecond))
	if renamed {
		t.Fatal("rename should not fire after double-click")
	}
	if !root.Expanded {
		t.Fatal("double-click should still perform branch activation behavior")
	}
}

func TestTreeViewLayoutKeepsTextAndIconVerticallyAligned(t *testing.T) {
	tree := NewTreeView("tree", geom.Rect{X: 0, Y: 0, W: 180, H: 120}, NewFileNode("Readme.txt"))
	rowRect := geom.Rect{X: 3, Y: 2, W: 160, H: 16}
	layout := tree.layoutEntry(rowRect, 0, 14)

	iconCenter := layout.icon.Y + layout.icon.H/2
	textCenter := layout.textY + 14/2
	if delta := iconCenter - textCenter; delta < -1 || delta > 1 {
		t.Fatalf("icon/text centers differ too much: icon=%d text=%d", iconCenter, textCenter)
	}
}

func TestTreeViewInlineRenameCommit(t *testing.T) {
	root := NewFileNode("Readme.txt")
	tree := NewTreeView("tree", geom.Rect{X: 0, Y: 0, W: 220, H: 120}, root)
	ctx := &fakeContext{}
	tree.SetFocused(true)

	var oldName, newName string
	tree.OnRenameCommit(func(_ *TreeNode, oldText, newText string) {
		oldName = oldText
		newName = newText
	})

	if !tree.beginRenameWithContext(ctx, root) {
		t.Fatal("begin rename should succeed")
	}
	if tree.renamingNode != root || tree.renameEdit == nil || !tree.renameEdit.Visible() {
		t.Fatal("begin rename should show inline edit")
	}

	if !tree.TextInput(ctx, event.TextInput{Text: "Notes"}) {
		t.Fatal("text input should be forwarded to inline edit")
	}
	if !tree.KeyDown(ctx, event.KeyEvent{Key: event.KeyEnter}) {
		t.Fatal("enter should commit inline rename")
	}
	if got := root.Text; got != "Notes.txt" {
		t.Fatalf("text after commit = %q, want %q", got, "Notes.txt")
	}
	if oldName != "Readme.txt" || newName != "Notes.txt" {
		t.Fatalf("rename callback = (%q, %q), want (%q, %q)", oldName, newName, "Readme.txt", "Notes.txt")
	}
	if tree.renamingNode != nil {
		t.Fatal("renaming state should clear after commit")
	}
}

func TestTreeViewInlineRenameCancel(t *testing.T) {
	root := NewFileNode("Readme.txt")
	tree := NewTreeView("tree", geom.Rect{X: 0, Y: 0, W: 220, H: 120}, root)
	ctx := &fakeContext{}
	tree.SetFocused(true)

	if !tree.beginRenameWithContext(ctx, root) {
		t.Fatal("begin rename should succeed")
	}
	tree.TextInput(ctx, event.TextInput{Text: "Draft.txt"})
	if !tree.KeyDown(ctx, event.KeyEvent{Key: event.KeyEscape}) {
		t.Fatal("escape should cancel inline rename")
	}
	if got := root.Text; got != "Readme.txt" {
		t.Fatalf("text after cancel = %q, want original", got)
	}
	if tree.renamingNode != nil {
		t.Fatal("renaming state should clear after cancel")
	}
}

func TestTreeViewInlineRenameSelectsFileStemOnly(t *testing.T) {
	root := NewFileNode("archive.tar.gz")
	tree := NewTreeView("tree", geom.Rect{X: 0, Y: 0, W: 220, H: 120}, root)
	ctx := &fakeContext{}
	tree.SetFocused(true)

	if !tree.beginRenameWithContext(ctx, root) {
		t.Fatal("begin rename should succeed")
	}
	if tree.renameEdit == nil {
		t.Fatal("rename edit should be created")
	}
	start, end, ok := tree.renameEdit.selectionRange()
	if !ok {
		t.Fatal("file rename should start with a selection")
	}
	if start != 0 || end != len([]rune("archive.tar")) {
		t.Fatalf("selection = (%d, %d), want (0, %d)", start, end, len([]rune("archive.tar")))
	}
}

func TestTreeViewInlineRenameSelectsWholeFolderName(t *testing.T) {
	root := NewFolderNode("Documents")
	tree := NewTreeView("tree", geom.Rect{X: 0, Y: 0, W: 220, H: 120}, root)
	ctx := &fakeContext{}
	tree.SetFocused(true)

	if !tree.beginRenameWithContext(ctx, root) {
		t.Fatal("begin rename should succeed")
	}
	start, end, ok := tree.renameEdit.selectionRange()
	if !ok {
		t.Fatal("folder rename should start with a selection")
	}
	if start != 0 || end != len([]rune("Documents")) {
		t.Fatalf("selection = (%d, %d), want full folder name", start, end)
	}
}

func TestFileStemEndKeepsDotfilesAndExtensionlessNamesWhole(t *testing.T) {
	if got := fileStemEnd(".gitignore"); got != len([]rune(".gitignore")) {
		t.Fatalf("stem for dotfile = %d, want full length", got)
	}
	if got := fileStemEnd("README"); got != len([]rune("README")) {
		t.Fatalf("stem for extensionless file = %d, want full length", got)
	}
	if got := fileStemEnd("name."); got != len([]rune("name.")) {
		t.Fatalf("stem for trailing-dot name = %d, want full length", got)
	}
}

func TestTreeViewSelectedExpanderKeepsGlyphVisible(t *testing.T) {
	tree := NewTreeView("tree", geom.Rect{X: 0, Y: 0, W: 180, H: 120}, NewFolderNode("Root"))
	canvas := paint.NewCanvas(24, 24)
	th := theme.DefaultClassic()
	rect := geom.Rect{X: 5, Y: 5, W: 9, H: 9}

	tree.paintExpander(PaintContext{Canvas: canvas, Theme: th}, rect, false, true, false)

	midX := rect.X + rect.W/2
	midY := rect.Y + rect.H/2
	idx := (midY*canvas.Width + midX) * 4
	got := [4]byte{canvas.Pix[idx], canvas.Pix[idx+1], canvas.Pix[idx+2], canvas.Pix[idx+3]}
	want := [4]byte{th.Colors.DarkShadow.R, th.Colors.DarkShadow.G, th.Colors.DarkShadow.B, th.Colors.DarkShadow.A}
	if got != want {
		t.Fatalf("selected expander glyph pixel = %#v, want %#v", got, want)
	}
}

func TestTreeViewMultiSelectedRowsUseHighlightText(t *testing.T) {
	root := NewFolderNode("Root",
		NewFileNode("Child 1"),
		NewFileNode("Child 2"),
	)
	root.Expanded = true
	tree := NewTreeView("tree", geom.Rect{X: 0, Y: 0, W: 220, H: 120}, root)
	tree.toggleSelection(root.Children[0])
	tree.toggleSelection(root.Children[1])

	entry := tree.visibleNodes()[1]
	rowRect := geom.Rect{X: 3, Y: 18, W: 180, H: 16}
	canvas := paint.NewCanvas(240, 80)
	err := tree.paintEntry(PaintContext{
		Canvas: canvas,
		Theme:  theme.DefaultClassic(),
	}, entry, rowRect, 14)
	if err != nil {
		t.Fatalf("paint entry failed: %v", err)
	}
}
