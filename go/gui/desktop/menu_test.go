package desktop

import (
	"testing"
	"time"

	"classicui/event"
	"classicui/geom"
	"classicui/theme"
	"classicui/widgets"
)

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
