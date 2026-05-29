package widgets

import (
	"testing"

	"classicui/event"
	"classicui/geom"
)

func TestEditTextInputClipboardAndDelete(t *testing.T) {
	edit := NewEdit("path", geom.Rect{X: 0, Y: 0, W: 160, H: 24})
	ctx := &fakeContext{}

	edit.SetFocused(true)
	edit.FocusGained(ctx)

	if !edit.TextInput(ctx, event.TextInput{Text: "hello"}) {
		t.Fatal("text input should be handled")
	}
	if got := edit.Text(); got != "hello" {
		t.Fatalf("text after input = %q, want %q", got, "hello")
	}

	edit.KeyDown(ctx, event.KeyEvent{Key: event.KeyA, Modifiers: event.ModCtrl})
	edit.KeyDown(ctx, event.KeyEvent{Key: event.KeyC, Modifiers: event.ModCtrl})
	if ctx.clipboard != "hello" {
		t.Fatalf("clipboard = %q, want %q", ctx.clipboard, "hello")
	}

	edit.KeyDown(ctx, event.KeyEvent{Key: event.KeyBackspace})
	if got := edit.Text(); got != "" {
		t.Fatalf("text after delete = %q, want empty", got)
	}

	edit.KeyDown(ctx, event.KeyEvent{Key: event.KeyV, Modifiers: event.ModCtrl})
	if got := edit.Text(); got != "hello" {
		t.Fatalf("text after paste = %q, want %q", got, "hello")
	}
}

func TestEditShiftSelectionAndCut(t *testing.T) {
	edit := NewEdit("name", geom.Rect{X: 0, Y: 0, W: 160, H: 24})
	ctx := &fakeContext{}

	edit.SetText("dust")
	edit.SetFocused(true)
	edit.FocusGained(ctx)

	edit.KeyDown(ctx, event.KeyEvent{Key: event.KeyHome})
	edit.KeyDown(ctx, event.KeyEvent{Key: event.KeyRight, Modifiers: event.ModShift})
	edit.KeyDown(ctx, event.KeyEvent{Key: event.KeyRight, Modifiers: event.ModShift})
	edit.KeyDown(ctx, event.KeyEvent{Key: event.KeyX, Modifiers: event.ModCtrl})

	if got := ctx.clipboard; got != "du" {
		t.Fatalf("clipboard after cut = %q, want %q", got, "du")
	}
	if got := edit.Text(); got != "st" {
		t.Fatalf("text after cut = %q, want %q", got, "st")
	}
}
