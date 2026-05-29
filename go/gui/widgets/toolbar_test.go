package widgets

import (
	"testing"

	"classicui/event"
	"classicui/geom"
)

func TestToolbarDispatchesCommandOnClick(t *testing.T) {
	toolbar := NewToolbar("tools", geom.Rect{X: 0, Y: 0, W: 180, H: 28},
		NewToolbarButton("cmd.add", "Add"),
		NewToolbarSeparator(),
		NewToolbarButton("cmd.exit", "Exit"),
	)
	ctx := &fakeContext{}

	toolbar.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, geom.Point{X: 10, Y: 10})
	toolbar.MouseUp(ctx, event.MouseButtonEvent{Button: event.MouseButtonLeft}, geom.Point{X: 10, Y: 10})

	if len(ctx.commands) != 1 || ctx.commands[0] != "cmd.add" {
		t.Fatalf("commands = %#v, want [cmd.add]", ctx.commands)
	}
}

func TestToolbarCheckedStateAndHitTesting(t *testing.T) {
	toolbar := NewToolbar("tools", geom.Rect{X: 0, Y: 0, W: 180, H: 28},
		NewToolbarButton("cmd.sort.name", "Name"),
		NewToolbarButton("cmd.sort.length", "Length"),
	)
	toolbar.SetChecked("cmd.sort.length", true)

	if !toolbar.isPressed(1) {
		t.Fatal("checked toolbar item should paint as pressed")
	}
	if toolbar.hitIndex(geom.Point{X: 70, Y: 10}, func(text string) geom.Size {
		return geom.Size{W: len([]rune(text)) * 8, H: 14}
	}) < 0 {
		t.Fatal("expected valid hit in toolbar button area")
	}
}
