package widgets

import (
	"testing"

	"classicui/event"
)

func TestMenuDisplayTextAndMnemonic(t *testing.T) {
	if got := MenuDisplayText("&File && E&xit"); got != "File & Exit" {
		t.Fatalf("display text = %q, want %q", got, "File & Exit")
	}

	item := NewMenuItem("cmd.exit", "&Exit", nil)
	mnemonic, ok := item.Mnemonic()
	if !ok || mnemonic != 'E' {
		t.Fatalf("mnemonic = %q, %v; want %q, true", mnemonic, ok, 'E')
	}
}

func TestAcceleratorDisplayAndMatch(t *testing.T) {
	accel := Accelerator{Key: event.KeyQ, Modifiers: event.ModCtrl | event.ModShift}
	if !accel.Matches(event.KeyQ, event.ModCtrl|event.ModShift) {
		t.Fatal("accelerator should match key/modifier pair")
	}
	if got := accel.DisplayLabel(); got != "Ctrl+Shift+Q" {
		t.Fatalf("label = %q, want %q", got, "Ctrl+Shift+Q")
	}
}

func TestMenuBarFindAccelerator(t *testing.T) {
	bar := NewMenuBar(
		NewSubmenuItem("&File", NewMenu(
			NewMenuItem("cmd.open", "&Open", &Accelerator{Key: event.KeyO, Modifiers: event.ModCtrl}),
			NewSubmenuItem("So&rt", NewMenu(
				NewMenuItem("cmd.sort.name", "By &Name", &Accelerator{Key: event.KeyN, Modifiers: event.ModCtrl}),
			)),
		)),
	)

	item, ok := bar.FindAccelerator(event.KeyN, event.ModCtrl)
	if !ok {
		t.Fatal("expected accelerator match")
	}
	if item.ID != "cmd.sort.name" {
		t.Fatalf("matched command = %q, want %q", item.ID, "cmd.sort.name")
	}
}
