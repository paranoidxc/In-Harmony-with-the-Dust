package desktop

import (
	"testing"

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
