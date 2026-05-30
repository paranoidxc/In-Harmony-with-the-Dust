package widgets

import "classicui/event"

type SelectionBehaviorOptions struct {
	MultiSelect       bool
	RecoverFromRecent bool
	BlankDragSelect   bool
}

func DefaultSelectionBehaviorOptions() SelectionBehaviorOptions {
	return SelectionBehaviorOptions{
		MultiSelect:       true,
		RecoverFromRecent: true,
		BlankDragSelect:   true,
	}
}

type ListBoxSelectionOptions = SelectionBehaviorOptions
type TreeViewSelectionOptions = SelectionBehaviorOptions

func DefaultListBoxSelectionOptions() ListBoxSelectionOptions {
	return DefaultSelectionBehaviorOptions()
}

func DefaultTreeViewSelectionOptions() TreeViewSelectionOptions {
	return DefaultSelectionBehaviorOptions()
}

type selectionBehavior struct {
	options SelectionBehaviorOptions
}

func (o SelectionBehaviorOptions) behavior() selectionBehavior {
	return selectionBehavior{options: o}
}

func (b selectionBehavior) normalizeModifiers(mods event.Modifiers) event.Modifiers {
	if !b.options.MultiSelect {
		return mods &^ (event.ModCtrl | event.ModShift)
	}
	return mods
}

func (b selectionBehavior) allowsMultiSelect() bool {
	return b.options.MultiSelect
}

func (b selectionBehavior) allowsBlankDrag() bool {
	return b.options.MultiSelect && b.options.BlankDragSelect
}

func (b selectionBehavior) allowsRecentRecovery() bool {
	return b.options.RecoverFromRecent
}

func (b selectionBehavior) extendRange(mods event.Modifiers) bool {
	return b.options.MultiSelect && mods&event.ModShift != 0
}

func (b selectionBehavior) toggleLeadShortcut(ev event.KeyEvent) bool {
	return b.options.MultiSelect && ev.Modifiers&event.ModCtrl != 0 && ev.Key == event.KeySpace
}

func (b selectionBehavior) selectAllShortcut(ev event.KeyEvent) bool {
	return ev.Modifiers&event.ModCtrl != 0 && ev.Key == event.KeyA
}
