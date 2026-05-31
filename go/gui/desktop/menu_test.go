package desktop

import (
	"testing"
	"time"

	"classicui/event"
	"classicui/geom"
	"classicui/theme"
	"classicui/widgets"
)

type popupProbeControl struct {
	*widgets.Panel
	downCount int
	lastDown  geom.Point
}

func newPopupProbeControl(bounds geom.Rect) *popupProbeControl {
	return &popupProbeControl{
		Panel: widgets.NewPanel("popup-probe", bounds),
	}
}

func (p *popupProbeControl) MouseDown(ctx widgets.EventContext, ev event.MouseButtonEvent, local geom.Point) {
	p.downCount++
	p.lastDown = local
	p.Panel.MouseDown(ctx, ev, local)
}

func TestDesktopDispatchesMenuAccelerator(t *testing.T) {
	d := New(geom.Size{W: 320, H: 200}, theme.DefaultClassic())
	win := NewWindow("main", geom.Rect{X: 20, Y: 20, W: 220, H: 140})
	win.SetMenuBar(widgets.NewMenuBar(
		widgets.NewSubmenuItem("&File", widgets.NewMenu(
			widgets.NewMenuItem("cmd.exit", "E&xit", &widgets.Accelerator{
				Key:       event.KeyQ,
				Modifiers: event.ModCtrl,
			}),
		)),
	))
	d.AddWindow(win)

	var got widgets.CommandID
	d.BindCommandHandler(func(_ *Window, cmd widgets.CommandID) {
		got = cmd
	})

	d.HandleEvent(event.KeyEvent{
		Down:      true,
		Key:       event.KeyQ,
		Modifiers: event.ModCtrl,
	})

	if got != "cmd.exit" {
		t.Fatalf("command = %q, want %q", got, "cmd.exit")
	}
}

func TestDesktopMenuKeyboardActivation(t *testing.T) {
	d := New(geom.Size{W: 320, H: 200}, theme.DefaultClassic())
	win := NewWindow("main", geom.Rect{X: 20, Y: 20, W: 220, H: 140})
	win.SetMenuBar(widgets.NewMenuBar(
		widgets.NewSubmenuItem("&File", widgets.NewMenu(
			widgets.NewMenuItem("cmd.open", "&Open", nil),
		)),
	))
	d.AddWindow(win)

	var got widgets.CommandID
	d.BindCommandHandler(func(_ *Window, cmd widgets.CommandID) {
		got = cmd
	})

	d.HandleEvent(event.KeyEvent{Down: true, Key: event.KeyLeftAlt})
	if !d.menuMode {
		t.Fatal("menu mode should be active after Alt")
	}
	if win.MenuBarActiveIndex() != 0 {
		t.Fatalf("active menu index = %d, want 0", win.MenuBarActiveIndex())
	}

	d.HandleEvent(event.KeyEvent{Down: true, Key: event.KeyDown})
	if len(d.menuPopups) != 1 {
		t.Fatalf("popup count = %d, want 1", len(d.menuPopups))
	}

	d.HandleEvent(event.KeyEvent{Down: true, Key: event.KeyEnter})
	if got != "cmd.open" {
		t.Fatalf("command = %q, want %q", got, "cmd.open")
	}
	if d.menuMode {
		t.Fatal("menu mode should close after command dispatch")
	}
}

func TestDesktopComboBoxOverlayClosesOnOutsideClick(t *testing.T) {
	d := New(geom.Size{W: 320, H: 220}, theme.DefaultClassic())
	win := NewWindow("main", geom.Rect{X: 20, Y: 20, W: 220, H: 160})
	combo := widgets.NewComboBox("sort", geom.Rect{X: 12, Y: 12, W: 120, H: 24})
	combo.SetItems([]string{"Name", "Length"})
	win.Content().Add(combo)
	d.AddWindow(win)

	client := win.ClientRect(d.theme)
	point := geom.Point{X: client.X + 20, Y: client.Y + 20}
	d.HandleEvent(event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft, Position: point})
	d.HandleEvent(event.MouseButtonEvent{Button: event.MouseButtonLeft, Position: point})

	if !d.overlayVisible(combo) {
		t.Fatal("combo box overlay should be visible after click")
	}

	outside := geom.Point{X: client.X + 180, Y: client.Y + 100}
	d.HandleEvent(event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft, Position: outside})
	if d.overlayVisible(combo) {
		t.Fatal("combo box overlay should close after outside click")
	}
}

func TestDesktopShowsTooltipAfterHoverDelay(t *testing.T) {
	d := New(geom.Size{W: 320, H: 220}, theme.DefaultClassic())
	win := NewWindow("main", geom.Rect{X: 20, Y: 20, W: 220, H: 160})
	button := widgets.NewButton("apply", "Apply", geom.Rect{X: 12, Y: 12, W: 80, H: 24})
	button.SetTooltip("Apply changes")
	win.Content().Add(button)
	d.AddWindow(win)

	client := win.ClientRect(d.theme)
	point := geom.Point{X: client.X + 20, Y: client.Y + 20}
	d.HandleEvent(event.MouseMove{Position: point})
	now := time.Now()
	d.Update(now.Add(tooltipDelay + 50*time.Millisecond))

	if d.tooltipOverlay == nil {
		t.Fatal("tooltip should become visible after hover delay")
	}
	if d.tooltipOverlay.text != "Apply changes" {
		t.Fatalf("tooltip text = %q, want %q", d.tooltipOverlay.text, "Apply changes")
	}
}

func TestDesktopTabLetsEditableComboBoxCommitCandidateThenMoveFocus(t *testing.T) {
	d := New(geom.Size{W: 320, H: 220}, theme.DefaultClassic())
	win := NewWindow("main", geom.Rect{X: 20, Y: 20, W: 220, H: 160})
	combo := widgets.NewComboBox("sort", geom.Rect{X: 12, Y: 12, W: 140, H: 24})
	combo.SetItems([]string{"Alpha", "Beta"})
	combo.SetEditable(true)
	next := widgets.NewButton("next", "Next", geom.Rect{X: 12, Y: 48, W: 80, H: 24})
	win.Content().Add(combo)
	win.Content().Add(next)
	d.AddWindow(win)
	d.setFocus(win, combo)

	combo.SetText("B")
	d.HandleEvent(event.KeyEvent{Down: true, Key: event.KeyTab})

	if got := combo.Text(); got != "Beta" {
		t.Fatalf("combo text after tab = %q, want %q", got, "Beta")
	}
	if d.focusedControl != next {
		t.Fatal("focus should advance to the next control after tab")
	}
}

func TestDesktopEditableComboBoxKeepsPopupOpenWhenClickingEditArea(t *testing.T) {
	d := New(geom.Size{W: 320, H: 220}, theme.DefaultClassic())
	win := NewWindow("main", geom.Rect{X: 20, Y: 20, W: 220, H: 160})
	combo := widgets.NewComboBox("sort", geom.Rect{X: 12, Y: 12, W: 140, H: 24})
	combo.SetItems([]string{"Alpha", "Beta"})
	combo.SetEditable(true)
	win.Content().Add(combo)
	d.AddWindow(win)

	client := win.ClientRect(d.theme)
	buttonPoint := geom.Point{X: client.X + combo.Bounds().X + combo.Bounds().W - 8, Y: client.Y + combo.Bounds().Y + 8}
	d.HandleEvent(event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft, Position: buttonPoint})
	d.HandleEvent(event.MouseButtonEvent{Button: event.MouseButtonLeft, Position: buttonPoint})
	if !d.overlayVisible(combo) {
		t.Fatal("combo popup should be visible after button click")
	}

	editPoint := geom.Point{X: client.X + combo.Bounds().X + 12, Y: client.Y + combo.Bounds().Y + 8}
	d.HandleEvent(event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft, Position: editPoint})
	if !d.overlayVisible(combo) {
		t.Fatal("combo popup should stay visible when clicking edit area")
	}
}

func TestDesktopRightMouseUpDoesNotImmediatelyCloseContextMenu(t *testing.T) {
	d := New(geom.Size{W: 320, H: 220}, theme.DefaultClassic())
	win := NewWindow("main", geom.Rect{X: 20, Y: 20, W: 220, H: 160})
	list := widgets.NewListView("files", geom.Rect{X: 12, Y: 12, W: 160, H: 100},
		widgets.ListViewColumn{Title: "Name", Width: 120},
	)
	list.SetItems([]widgets.ListViewItem{{Texts: []string{"One"}}})
	list.SetContextMenu(widgets.NewMenu(
		widgets.NewMenuItem("cmd.open", "&Open", nil),
	))
	win.Content().Add(list)
	d.AddWindow(win)

	client := win.ClientRect(d.theme)
	rowPoint := geom.Point{
		X: client.X + list.Bounds().X + 20,
		Y: client.Y + list.Bounds().Y + 28,
	}
	d.HandleEvent(event.MouseButtonEvent{Down: true, Button: event.MouseButtonRight, Position: rowPoint})
	if !d.menuMode || len(d.menuPopups) != 1 {
		t.Fatal("context menu should be open after right mouse down")
	}

	d.HandleEvent(event.MouseButtonEvent{Button: event.MouseButtonRight, Position: rowPoint})
	if !d.menuMode || len(d.menuPopups) != 1 {
		t.Fatal("right mouse up should not immediately close the context menu")
	}
}

func TestDesktopRoutesMouseToTopmostControlPopupBeforeMenu(t *testing.T) {
	d := New(geom.Size{W: 320, H: 220}, theme.DefaultClassic())
	win := NewWindow("main", geom.Rect{X: 20, Y: 20, W: 220, H: 160})
	owner := widgets.NewButton("owner", "Owner", geom.Rect{X: 12, Y: 12, W: 80, H: 24})
	win.Content().Add(owner)
	d.AddWindow(win)

	client := win.ClientRect(d.theme)
	menuPoint := geom.Point{X: client.X + 30, Y: client.Y + 30}
	if !d.showContextMenu(win, owner, geom.Rect{X: 8, Y: 8, W: 1, H: 1}, widgets.NewMenu(
		widgets.NewMenuItem("cmd.open", "&Open", nil),
	)) {
		t.Fatal("showContextMenu should succeed")
	}
	if !d.menuMode || len(d.menuPopups) != 1 {
		t.Fatal("context menu should be open")
	}

	popupContent := newPopupProbeControl(geom.Rect{X: 0, Y: 0, W: 96, H: 48})
	if !d.showControlOverlay(win, widgets.PopupRequest{
		Owner:          owner,
		Content:        popupContent,
		Anchor:         geom.Rect{X: 16, Y: 16, W: 1, H: 1},
		Placement:      widgets.PopupBelowStart,
		CloseOnOutside: true,
		Kind:           widgets.PopupKindInteractive,
	}) {
		t.Fatal("showControlOverlay should succeed")
	}

	overlay := d.topControlOverlay()
	if overlay == nil {
		t.Fatal("control overlay should be topmost")
	}

	click := geom.Point{X: overlay.rect.X + 10, Y: overlay.rect.Y + 10}
	d.HandleEvent(event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft, Position: click})

	if popupContent.downCount != 1 {
		t.Fatalf("popup down count = %d, want 1", popupContent.downCount)
	}
	if !d.menuMode || len(d.menuPopups) != 1 {
		t.Fatal("menu popup should remain open when clicking top control popup")
	}
	if popupContent.lastDown != (geom.Point{X: 10, Y: 10}) {
		t.Fatalf("popup local point = %+v, want {X:10 Y:10}", popupContent.lastDown)
	}

	_ = menuPoint
}
