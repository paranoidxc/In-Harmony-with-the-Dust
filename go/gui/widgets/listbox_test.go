package widgets

import (
	"testing"
	"time"

	"classicui/event"
	"classicui/geom"
)

func TestListBoxCtrlClickTogglesMultiSelection(t *testing.T) {
	list := NewListBox("files", geom.Rect{X: 0, Y: 0, W: 180, H: 120})
	list.SetItems([]string{"One", "Two", "Three"})
	ctx := &fakeContext{}

	list.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft, Modifiers: event.ModCtrl}, geom.Point{X: 8, Y: 24})
	list.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft, Modifiers: event.ModCtrl}, geom.Point{X: 8, Y: 40})

	selected := list.SelectedIndices()
	if len(selected) != 3 || selected[0] != 0 || selected[1] != 1 || selected[2] != 2 {
		t.Fatalf("selected indices = %#v, want [0 1 2]", selected)
	}
	if list.SelectedIndex() != 2 {
		t.Fatalf("lead selection = %d, want 2", list.SelectedIndex())
	}
}

func TestListBoxSingleSelectModeIgnoresMultiSelectModifiers(t *testing.T) {
	list := NewListBox("files", geom.Rect{X: 0, Y: 0, W: 180, H: 120})
	list.SetItems([]string{"One", "Two", "Three"})
	list.SetMultiSelect(false)
	ctx := &fakeContext{}

	list.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft, Modifiers: event.ModCtrl}, geom.Point{X: 8, Y: 24})
	list.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft, Modifiers: event.ModShift}, geom.Point{X: 8, Y: 40})

	selected := list.SelectedIndices()
	if len(selected) != 1 || selected[0] != 2 {
		t.Fatalf("selected indices = %#v, want [2]", selected)
	}
	if !list.KeyDown(ctx, event.KeyEvent{Key: event.KeyA, Modifiers: event.ModCtrl}) {
		t.Fatal("ctrl+a should still be handled in single-select mode")
	}
	selected = list.SelectedIndices()
	if len(selected) != 1 || selected[0] != 2 {
		t.Fatalf("selected indices after ctrl+a = %#v, want still [2]", selected)
	}
}

func TestListBoxSelectionOptionsCanDisableRecentRecovery(t *testing.T) {
	list := NewListBox("files", geom.Rect{X: 0, Y: 0, W: 180, H: 120})
	list.SetItems([]string{"One", "Two", "Three", "Four"})
	options := list.SelectionOptions()
	options.RecoverFromRecent = false
	list.SetSelectionOptions(options)
	ctx := &fakeContext{}

	list.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, geom.Point{X: 8, Y: 56})
	list.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, geom.Point{X: 24, Y: 96})

	if !list.KeyDown(ctx, event.KeyEvent{Key: event.KeyDown}) {
		t.Fatal("down should recover from empty selection")
	}
	if list.SelectedIndex() != 0 {
		t.Fatalf("selected index with recent recovery disabled = %d, want 0", list.SelectedIndex())
	}
}

func TestListBoxSelectionOptionsCanDisableBlankDragSelect(t *testing.T) {
	list := NewListBox("files", geom.Rect{X: 0, Y: 0, W: 180, H: 120})
	list.SetItems([]string{"One", "Two", "Three", "Four"})
	options := list.SelectionOptions()
	options.BlankDragSelect = false
	list.SetSelectionOptions(options)
	ctx := &fakeContext{}

	list.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, geom.Point{X: 24, Y: 72})
	list.MouseMove(ctx, geom.Point{X: 26, Y: 20})

	if !list.dragSelecting {
		t.Fatal("blank drag should still cross the drag threshold")
	}
	if _, ok := list.marqueeRect(); ok {
		t.Fatal("blank drag marquee should be disabled")
	}
	if selected := list.SelectedIndices(); len(selected) != 0 {
		t.Fatalf("selected indices after disabled blank drag = %#v, want none", selected)
	}
}

func TestListBoxShiftClickSelectsRange(t *testing.T) {
	list := NewListBox("files", geom.Rect{X: 0, Y: 0, W: 180, H: 120})
	list.SetItems([]string{"One", "Two", "Three", "Four"})
	ctx := &fakeContext{}

	list.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, geom.Point{X: 8, Y: 24})
	list.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft, Modifiers: event.ModShift}, geom.Point{X: 8, Y: 56})

	selected := list.SelectedIndices()
	if len(selected) != 3 || selected[0] != 1 || selected[1] != 2 || selected[2] != 3 {
		t.Fatalf("selected indices = %#v, want [1 2 3]", selected)
	}
}

func TestListBoxKeyboardShiftExtendsSelection(t *testing.T) {
	list := NewListBox("files", geom.Rect{X: 0, Y: 0, W: 180, H: 120})
	list.SetItems([]string{"One", "Two", "Three", "Four"})
	ctx := &fakeContext{}

	list.SetSelectedIndex(1)
	if !list.KeyDown(ctx, event.KeyEvent{Key: event.KeyDown, Modifiers: event.ModShift}) {
		t.Fatal("shift+down should be handled")
	}
	if !list.KeyDown(ctx, event.KeyEvent{Key: event.KeyDown, Modifiers: event.ModShift}) {
		t.Fatal("second shift+down should be handled")
	}

	selected := list.SelectedIndices()
	if len(selected) != 3 || selected[0] != 1 || selected[1] != 2 || selected[2] != 3 {
		t.Fatalf("selected indices = %#v, want [1 2 3]", selected)
	}
}

func TestListBoxCtrlASelectsAll(t *testing.T) {
	list := NewListBox("files", geom.Rect{X: 0, Y: 0, W: 180, H: 120})
	list.SetItems([]string{"One", "Two", "Three"})
	ctx := &fakeContext{}

	if !list.KeyDown(ctx, event.KeyEvent{Key: event.KeyA, Modifiers: event.ModCtrl}) {
		t.Fatal("ctrl+a should be handled")
	}
	selected := list.SelectedIndices()
	if len(selected) != 3 || selected[0] != 0 || selected[1] != 1 || selected[2] != 2 {
		t.Fatalf("selected indices = %#v, want [0 1 2]", selected)
	}
}

func TestListBoxCtrlSpaceTogglesLeadSelection(t *testing.T) {
	list := NewListBox("files", geom.Rect{X: 0, Y: 0, W: 180, H: 120})
	list.SetItems([]string{"One", "Two", "Three"})
	ctx := &fakeContext{}

	list.toggleSelection(1, true)
	if !list.KeyDown(ctx, event.KeyEvent{Key: event.KeySpace, Modifiers: event.ModCtrl}) {
		t.Fatal("ctrl+space should be handled")
	}

	selected := list.SelectedIndices()
	if len(selected) != 1 || selected[0] != 0 {
		t.Fatalf("selected indices = %#v, want [0]", selected)
	}
}

func TestListBoxBlankClickClearsSelectionAndFocuses(t *testing.T) {
	list := NewListBox("files", geom.Rect{X: 0, Y: 0, W: 180, H: 120})
	list.SetItems([]string{"One", "Two", "Three"})
	ctx := &fakeContext{}

	list.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, geom.Point{X: 8, Y: 72})

	if ctx.focused != list {
		t.Fatal("blank click should focus the list")
	}
	if list.SelectedIndex() != -1 {
		t.Fatalf("selected index = %d, want -1", list.SelectedIndex())
	}
	if selected := list.SelectedIndices(); len(selected) != 0 {
		t.Fatalf("selected indices = %#v, want none", selected)
	}
}

func TestListBoxCtrlBlankClickPreservesSelection(t *testing.T) {
	list := NewListBox("files", geom.Rect{X: 0, Y: 0, W: 180, H: 120})
	list.SetItems([]string{"One", "Two", "Three"})
	ctx := &fakeContext{}

	list.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, geom.Point{X: 8, Y: 24})
	list.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft, Modifiers: event.ModCtrl}, geom.Point{X: 8, Y: 72})

	if ctx.focused != list {
		t.Fatal("ctrl+blank click should focus the list")
	}
	selected := list.SelectedIndices()
	if len(selected) != 1 || selected[0] != 1 {
		t.Fatalf("selected indices = %#v, want [1]", selected)
	}
	if list.SelectedIndex() != 1 {
		t.Fatalf("lead selection = %d, want 1", list.SelectedIndex())
	}
}

func TestListBoxKeyboardRecoversRecentSelectionAfterBlankClear(t *testing.T) {
	list := NewListBox("files", geom.Rect{X: 0, Y: 0, W: 180, H: 120})
	list.SetItems([]string{"One", "Two", "Three", "Four"})
	ctx := &fakeContext{}

	list.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, geom.Point{X: 8, Y: 40})
	list.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, geom.Point{X: 24, Y: 96})

	if !list.KeyDown(ctx, event.KeyEvent{Key: event.KeyDown}) {
		t.Fatal("down should recover from empty selection")
	}
	if list.SelectedIndex() != 2 {
		t.Fatalf("selected index after recovery = %d, want 2", list.SelectedIndex())
	}
}

func TestListBoxCtrlSpaceRecoversRecentSelectionAfterBlankClear(t *testing.T) {
	list := NewListBox("files", geom.Rect{X: 0, Y: 0, W: 180, H: 120})
	list.SetItems([]string{"One", "Two", "Three", "Four"})
	ctx := &fakeContext{}

	list.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, geom.Point{X: 8, Y: 56})
	list.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, geom.Point{X: 24, Y: 96})

	if !list.KeyDown(ctx, event.KeyEvent{Key: event.KeySpace, Modifiers: event.ModCtrl}) {
		t.Fatal("ctrl+space should recover from empty selection")
	}
	if list.SelectedIndex() != 3 {
		t.Fatalf("selected index after ctrl+space recovery = %d, want 3", list.SelectedIndex())
	}
}

func TestListBoxHomeEndOverrideRecentSelectionRecovery(t *testing.T) {
	list := NewListBox("files", geom.Rect{X: 0, Y: 0, W: 180, H: 120})
	list.SetItems([]string{"One", "Two", "Three", "Four"})
	ctx := &fakeContext{}

	list.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, geom.Point{X: 8, Y: 40})
	list.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, geom.Point{X: 24, Y: 96})

	if !list.KeyDown(ctx, event.KeyEvent{Key: event.KeyEnd}) {
		t.Fatal("end should recover to last item")
	}
	if list.SelectedIndex() != 3 {
		t.Fatalf("selected index after end = %d, want 3", list.SelectedIndex())
	}

	list.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, geom.Point{X: 24, Y: 96})
	if !list.KeyDown(ctx, event.KeyEvent{Key: event.KeyHome}) {
		t.Fatal("home should recover to first item")
	}
	if list.SelectedIndex() != 0 {
		t.Fatalf("selected index after home = %d, want 0", list.SelectedIndex())
	}
}

func TestListBoxDragSelectionStartsAfterThreshold(t *testing.T) {
	list := NewListBox("files", geom.Rect{X: 0, Y: 0, W: 180, H: 120})
	list.SetItems([]string{"One", "Two", "Three", "Four"})
	ctx := &fakeContext{}

	list.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, geom.Point{X: 8, Y: 8})
	list.MouseMove(ctx, geom.Point{X: 9, Y: 9})
	if list.dragSelecting {
		t.Fatal("small move should not start drag selection")
	}
	list.MouseMove(ctx, geom.Point{X: 16, Y: 40})
	if !list.dragSelecting {
		t.Fatal("move beyond threshold should start drag selection")
	}
	selected := list.SelectedIndices()
	if len(selected) != 3 || selected[0] != 0 || selected[2] != 2 {
		t.Fatalf("selected indices after drag = %#v, want [0 1 2]", selected)
	}
}

func TestListBoxBlankDragSelectsIntersectingRows(t *testing.T) {
	list := NewListBox("files", geom.Rect{X: 0, Y: 0, W: 180, H: 120})
	list.SetItems([]string{"One", "Two", "Three", "Four"})
	ctx := &fakeContext{}

	list.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, geom.Point{X: 24, Y: 72})
	list.MouseMove(ctx, geom.Point{X: 26, Y: 20})

	if !list.dragSelecting {
		t.Fatal("blank drag should enter drag-select mode")
	}
	if _, ok := list.marqueeRect(); !ok {
		t.Fatal("blank drag should expose a marquee rect")
	}
	selected := list.SelectedIndices()
	if len(selected) != 3 || selected[0] != 1 || selected[1] != 2 || selected[2] != 3 {
		t.Fatalf("selected indices after blank drag = %#v, want [1 2 3]", selected)
	}
	if list.SelectedIndex() != 3 {
		t.Fatalf("lead selection = %d, want 3", list.SelectedIndex())
	}
}

func TestListBoxCtrlBlankDragUnionsSelection(t *testing.T) {
	list := NewListBox("files", geom.Rect{X: 0, Y: 0, W: 180, H: 120})
	list.SetItems([]string{"One", "Two", "Three", "Four"})
	ctx := &fakeContext{}

	list.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, geom.Point{X: 8, Y: 8})
	list.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft, Modifiers: event.ModCtrl}, geom.Point{X: 24, Y: 72})
	list.MouseMove(ctx, geom.Point{X: 26, Y: 20})

	selected := list.SelectedIndices()
	if len(selected) != 4 || selected[0] != 0 || selected[3] != 3 {
		t.Fatalf("selected indices after ctrl+blank drag = %#v, want [0 1 2 3]", selected)
	}
}

func TestListBoxDragAutoScroll(t *testing.T) {
	list := NewListBox("files", geom.Rect{X: 0, Y: 0, W: 180, H: 52})
	list.SetItems([]string{"One", "Two", "Three", "Four", "Five", "Six"})
	ctx := &fakeContext{}

	list.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, geom.Point{X: 8, Y: 8})
	list.MouseMove(ctx, geom.Point{X: 16, Y: 70})
	if !list.dragSelecting {
		t.Fatal("drag should be active")
	}
	if !list.Tick(ctx, time.Date(2026, time.May, 30, 10, 0, 1, 0, time.UTC)) {
		t.Fatal("tick should auto-scroll while dragging past bottom")
	}
	if list.topIndex == 0 {
		t.Fatal("auto-scroll should advance top index")
	}
}
