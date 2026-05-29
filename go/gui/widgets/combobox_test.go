package widgets

import (
	"testing"

	"classicui/event"
	"classicui/geom"
)

func TestComboBoxKeyboardSelectionAndPopupLifecycle(t *testing.T) {
	combo := NewComboBox("sort", geom.Rect{X: 0, Y: 0, W: 120, H: 24})
	combo.SetItems([]string{"Name", "Length", "Date"})
	ctx := &fakeContext{}

	if combo.SelectedIndex() != 0 {
		t.Fatalf("initial selection = %d, want 0", combo.SelectedIndex())
	}
	if !combo.KeyDown(ctx, event.KeyEvent{Key: event.KeyDown}) {
		t.Fatal("down key should be handled")
	}
	if combo.SelectedIndex() != 1 {
		t.Fatalf("selection after down = %d, want 1", combo.SelectedIndex())
	}

	combo.KeyDown(ctx, event.KeyEvent{Key: event.KeyEnter})
	if !combo.dropped || ctx.overlay == nil {
		t.Fatal("enter should open popup overlay")
	}

	combo.KeyDown(ctx, event.KeyEvent{Key: event.KeyDown})
	popup, ok := ctx.overlay.Content.(*ListBox)
	if !ok {
		t.Fatalf("overlay content = %T, want *ListBox", ctx.overlay.Content)
	}
	if popup.SelectedIndex() != 2 {
		t.Fatalf("popup selection = %d, want 2", popup.SelectedIndex())
	}

	combo.KeyDown(ctx, event.KeyEvent{Key: event.KeyEnter})
	if combo.SelectedIndex() != 2 {
		t.Fatalf("selection after popup commit = %d, want 2", combo.SelectedIndex())
	}
	if combo.dropped || ctx.overlay != nil {
		t.Fatal("popup should close after commit")
	}
}

func TestComboBoxMouseToggleAndEscapeClose(t *testing.T) {
	combo := NewComboBox("sort", geom.Rect{X: 0, Y: 0, W: 120, H: 24})
	combo.SetItems([]string{"Name", "Length"})
	ctx := &fakeContext{}

	combo.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, geom.Point{X: 8, Y: 8})
	if !combo.dropped || ctx.overlay == nil {
		t.Fatal("mouse down should open popup")
	}

	combo.KeyDown(ctx, event.KeyEvent{Key: event.KeyEscape})
	if combo.dropped || ctx.overlay != nil {
		t.Fatal("escape should close popup")
	}
}

func TestEditableComboBoxTextInputAndPopupCommit(t *testing.T) {
	combo := NewComboBox("path", geom.Rect{X: 0, Y: 0, W: 160, H: 24})
	combo.SetItems([]string{"Alpha", "Beta", "Gamma"})
	combo.SetEditable(true)
	ctx := &fakeContext{}
	commits := 0
	combo.OnCommit(func(int, string) {
		commits++
	})

	if !combo.TextInput(ctx, event.TextInput{Text: "D"}) {
		t.Fatal("text input should be handled in editable mode")
	}
	if got := combo.Text(); got != "AlphaD" {
		t.Fatalf("text after input = %q, want %q", got, "AlphaD")
	}
	if combo.SelectedIndex() != -1 {
		t.Fatalf("selection after freeform input = %d, want -1", combo.SelectedIndex())
	}
	if commits != 0 {
		t.Fatalf("commit count after freeform input = %d, want 0", commits)
	}

	combo.SetText("B")

	if !combo.KeyDown(ctx, event.KeyEvent{Key: event.KeyDown}) {
		t.Fatal("down key should open popup in editable mode")
	}
	if !combo.dropped || ctx.overlay == nil {
		t.Fatal("down key should open popup overlay")
	}

	popup, ok := ctx.overlay.Content.(*ListBox)
	if !ok {
		t.Fatalf("overlay content = %T, want *ListBox", ctx.overlay.Content)
	}
	if popup.SelectedIndex() != 1 {
		t.Fatalf("popup selection = %d, want 1", popup.SelectedIndex())
	}

	itemPoint := geom.Point{X: 10, Y: 20}
	popup.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, itemPoint)
	if commits != 0 {
		t.Fatalf("commit count after mouse down = %d, want 0", commits)
	}
	popup.MouseUp(ctx, event.MouseButtonEvent{Button: event.MouseButtonLeft}, itemPoint)
	if got := combo.Text(); got != "Beta" {
		t.Fatalf("text after popup commit = %q, want %q", got, "Beta")
	}
	if combo.SelectedIndex() != 1 {
		t.Fatalf("selection after popup commit = %d, want 1", combo.SelectedIndex())
	}
	if commits != 1 {
		t.Fatalf("commit count after popup commit = %d, want 1", commits)
	}
}

func TestEditableComboBoxEscapeRestoresCommittedValue(t *testing.T) {
	combo := NewComboBox("path", geom.Rect{X: 0, Y: 0, W: 160, H: 24})
	combo.SetItems([]string{"Alpha", "Beta"})
	combo.SetEditable(true)
	ctx := &fakeContext{}

	if !combo.TextInput(ctx, event.TextInput{Text: "D"}) {
		t.Fatal("text input should be handled")
	}
	if got := combo.Text(); got != "AlphaD" {
		t.Fatalf("text after input = %q, want %q", got, "AlphaD")
	}

	if !combo.KeyDown(ctx, event.KeyEvent{Key: event.KeyEscape}) {
		t.Fatal("escape should be handled")
	}
	if got := combo.Text(); got != "Alpha" {
		t.Fatalf("text after escape = %q, want %q", got, "Alpha")
	}
	if combo.SelectedIndex() != 0 {
		t.Fatalf("selection after escape = %d, want 0", combo.SelectedIndex())
	}
}
