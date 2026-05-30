package widgets

import (
	"testing"

	"classicui/event"
)

func TestSelectionBehaviorSingleSelectSuppressesRangeAndToggleModifiers(t *testing.T) {
	behavior := SelectionBehaviorOptions{
		MultiSelect:       false,
		RecoverFromRecent: true,
		BlankDragSelect:   true,
	}.behavior()

	mods := behavior.normalizeModifiers(event.ModCtrl | event.ModShift)
	if mods != 0 {
		t.Fatalf("normalized modifiers = %v, want 0", mods)
	}
	if behavior.toggleLeadShortcut(event.KeyEvent{Key: event.KeySpace, Modifiers: event.ModCtrl}) {
		t.Fatal("single-select behavior should not expose ctrl+space toggle shortcut")
	}
	if behavior.allowsBlankDrag() {
		t.Fatal("single-select behavior should not allow blank drag marquee")
	}
}

func TestSelectionBehaviorMultiSelectRetainsConfiguredGestures(t *testing.T) {
	behavior := DefaultSelectionBehaviorOptions().behavior()

	if !behavior.extendRange(event.ModShift) {
		t.Fatal("multi-select behavior should extend range with shift")
	}
	if !behavior.toggleLeadShortcut(event.KeyEvent{Key: event.KeySpace, Modifiers: event.ModCtrl}) {
		t.Fatal("multi-select behavior should expose ctrl+space toggle shortcut")
	}
	if !behavior.selectAllShortcut(event.KeyEvent{Key: event.KeyA, Modifiers: event.ModCtrl}) {
		t.Fatal("multi-select behavior should expose ctrl+a select-all shortcut")
	}
}
