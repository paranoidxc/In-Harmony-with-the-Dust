package widgets

import (
	"testing"
	"time"

	"classicui/event"
	"classicui/geom"
)

func TestListViewCtrlClickTogglesMultiSelection(t *testing.T) {
	view := newTestListView()
	ctx := &fakeContext{}

	view.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft, Modifiers: event.ModCtrl}, geom.Point{X: 8, Y: 46})
	view.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft, Modifiers: event.ModCtrl}, geom.Point{X: 8, Y: 62})

	selected := view.SelectedIndices()
	if len(selected) != 3 || selected[0] != 0 || selected[1] != 1 || selected[2] != 2 {
		t.Fatalf("selected indices = %#v, want [0 1 2]", selected)
	}
}

func TestListViewSingleSelectModeIgnoresMultiSelectModifiers(t *testing.T) {
	view := newTestListView()
	view.SetMultiSelect(false)
	ctx := &fakeContext{}

	view.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft, Modifiers: event.ModCtrl}, geom.Point{X: 8, Y: 46})
	view.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft, Modifiers: event.ModShift}, geom.Point{X: 8, Y: 62})

	selected := view.SelectedIndices()
	if len(selected) != 1 || selected[0] != 2 {
		t.Fatalf("selected indices = %#v, want [2]", selected)
	}
}

func TestListViewColumnClickNotifies(t *testing.T) {
	view := newTestListView()
	ctx := &fakeContext{}

	clicked := -1
	view.OnColumnClick(func(index int) {
		clicked = index
	})

	view.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, geom.Point{X: 130, Y: 8})
	view.MouseUp(ctx, event.MouseButtonEvent{Button: event.MouseButtonLeft}, geom.Point{X: 130, Y: 8})

	if clicked != 1 {
		t.Fatalf("clicked column = %d, want 1", clicked)
	}
}

func TestListViewHeaderDividerDragResizesColumn(t *testing.T) {
	view := newTestListView()
	ctx := &fakeContext{}

	view.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, geom.Point{X: 120, Y: 8})
	view.MouseMove(ctx, geom.Point{X: 144, Y: 8})
	view.MouseUp(ctx, event.MouseButtonEvent{Button: event.MouseButtonLeft}, geom.Point{X: 144, Y: 8})

	columns := view.Columns()
	if columns[0].Width != 144 {
		t.Fatalf("column 0 width = %d, want 144", columns[0].Width)
	}
}

func TestListViewHeaderDividerDoubleClickAutoFitsColumn(t *testing.T) {
	view := newTestListView()
	ctx := &fakeContext{}

	view.now = func() time.Time {
		return time.Date(2026, time.May, 30, 10, 0, 0, 0, time.UTC)
	}
	point := geom.Point{X: 120, Y: 8}
	view.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, point)
	view.MouseUp(ctx, event.MouseButtonEvent{Button: event.MouseButtonLeft}, point)
	view.now = func() time.Time {
		return time.Date(2026, time.May, 30, 10, 0, 0, int(200*time.Millisecond), time.UTC)
	}
	view.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, point)

	columns := view.Columns()
	if columns[0].Width != 52 {
		t.Fatalf("column 0 width = %d, want auto-fit width 52", columns[0].Width)
	}
}

func TestListViewSingleClickSelectsButDoesNotActivate(t *testing.T) {
	view := newTestListView()
	ctx := &fakeContext{}

	activations := 0
	view.OnActivate(func(int, ListViewItem) {
		activations++
	})

	view.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, geom.Point{X: 8, Y: 46})
	view.MouseUp(ctx, event.MouseButtonEvent{Button: event.MouseButtonLeft}, geom.Point{X: 8, Y: 46})

	if activations != 0 {
		t.Fatalf("activations = %d, want 0 on single click", activations)
	}
	if view.SelectedIndex() != 1 {
		t.Fatalf("selected index = %d, want 1", view.SelectedIndex())
	}
}

func TestListViewDoubleClickActivates(t *testing.T) {
	view := newTestListView()
	ctx := &fakeContext{}

	activations := 0
	view.now = func() time.Time {
		return time.Date(2026, time.May, 30, 10, 0, 0, 0, time.UTC)
	}
	view.OnActivate(func(int, ListViewItem) {
		activations++
	})

	point := geom.Point{X: 8, Y: 46}
	view.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, point)
	view.MouseUp(ctx, event.MouseButtonEvent{Button: event.MouseButtonLeft}, point)
	view.now = func() time.Time {
		return time.Date(2026, time.May, 30, 10, 0, 0, int(200*time.Millisecond), time.UTC)
	}
	view.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, point)
	view.MouseUp(ctx, event.MouseButtonEvent{Button: event.MouseButtonLeft}, point)

	if activations != 1 {
		t.Fatalf("activations = %d, want 1 on double click", activations)
	}
}

func TestListViewEnterActivatesLeadItem(t *testing.T) {
	view := newTestListView()
	ctx := &fakeContext{}

	activations := 0
	view.OnActivate(func(index int, _ ListViewItem) {
		if index != 2 {
			t.Fatalf("activated index = %d, want 2", index)
		}
		activations++
	})

	view.SetSelectedIndex(2)
	if !view.KeyDown(ctx, event.KeyEvent{Key: event.KeyEnter}) {
		t.Fatal("enter should be handled")
	}
	if activations != 1 {
		t.Fatalf("activations = %d, want 1", activations)
	}
}

func TestListViewRightClickSelectsItemAndShowsContextMenu(t *testing.T) {
	view := newTestListView()
	ctx := &fakeContext{}
	menu := NewMenu(NewMenuItem("cmd.open", "&Open", nil))
	view.SetContextMenu(menu)

	view.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonRight}, geom.Point{X: 8, Y: 46})

	if view.SelectedIndex() != 1 {
		t.Fatalf("selected index = %d, want 1", view.SelectedIndex())
	}
	if ctx.contextOwner != view {
		t.Fatal("context menu owner should be the list view")
	}
	if ctx.contextMenu != menu {
		t.Fatal("context menu should be shown")
	}
}

func TestListViewContextMenuProviderReceivesBlankAndItemContext(t *testing.T) {
	view := newTestListView()
	ctx := &fakeContext{}

	var gotBlank, gotItem bool
	view.SetContextMenuProvider(func(info ListViewContextMenuInfo) *Menu {
		if info.HasItem {
			gotItem = info.Index == 1 && info.Item.Texts[0] == "Two"
		} else {
			gotBlank = true
		}
		return NewMenu(NewMenuItem("cmd.open", "&Open", nil))
	})

	lastRow, ok := view.rowRect(3)
	if !ok {
		t.Fatal("expected row rect for last item")
	}
	blankY := lastRow.Bottom() + 2
	if blankY >= view.itemsRect(LocalRect(view)).Bottom() {
		t.Fatal("test list view should expose blank area below last row")
	}

	view.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonRight}, geom.Point{X: 8, Y: 46})
	view.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonRight}, geom.Point{X: 24, Y: blankY})

	if !gotItem {
		t.Fatal("context menu provider should receive item context")
	}
	if !gotBlank {
		t.Fatal("context menu provider should receive blank context")
	}
}

func TestListViewBlankDragSelectsRows(t *testing.T) {
	view := newTestListView()
	ctx := &fakeContext{}

	view.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, geom.Point{X: 24, Y: 80})
	view.MouseMove(ctx, geom.Point{X: 26, Y: 32})

	if !view.dragSelecting {
		t.Fatal("blank drag should enter drag-select mode")
	}
	if _, ok := view.marqueeRect(); !ok {
		t.Fatal("blank drag should expose a marquee rect")
	}
	selected := view.SelectedIndices()
	if len(selected) != 4 || selected[0] != 0 || selected[1] != 1 || selected[2] != 2 || selected[3] != 3 {
		t.Fatalf("selected indices after blank drag = %#v, want [0 1 2 3]", selected)
	}
}

func TestListViewDragAutoScroll(t *testing.T) {
	view := NewListView("files", geom.Rect{X: 0, Y: 0, W: 220, H: 72},
		ListViewColumn{Title: "Name", Width: 120},
		ListViewColumn{Title: "Size", Width: 80, Align: HeaderAlignRight},
	)
	view.SetItems([]ListViewItem{
		{Texts: []string{"One", "1 KB"}},
		{Texts: []string{"Two", "2 KB"}},
		{Texts: []string{"Three", "3 KB"}},
		{Texts: []string{"Four", "4 KB"}},
		{Texts: []string{"Five", "5 KB"}},
		{Texts: []string{"Six", "6 KB"}},
	})
	ctx := &fakeContext{}

	view.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, geom.Point{X: 8, Y: 28})
	view.MouseMove(ctx, geom.Point{X: 16, Y: 90})
	if !view.dragSelecting {
		t.Fatal("drag should be active")
	}
	if !view.Tick(ctx, time.Date(2026, time.May, 30, 10, 0, 1, 0, time.UTC)) {
		t.Fatal("tick should auto-scroll while dragging past bottom")
	}
	if view.topIndex == 0 {
		t.Fatal("auto-scroll should advance top index")
	}
}

func newTestListView() *ListView {
	view := NewListView("files", geom.Rect{X: 0, Y: 0, W: 220, H: 140},
		ListViewColumn{Title: "Name", Width: 120},
		ListViewColumn{Title: "Size", Width: 80, Align: HeaderAlignRight},
	)
	view.SetItems([]ListViewItem{
		{Texts: []string{"One", "1 KB"}},
		{Texts: []string{"Two", "2 KB"}},
		{Texts: []string{"Three", "3 KB"}},
		{Texts: []string{"Four", "4 KB"}},
	})
	return view
}
