package widgets

import (
	"classicui/event"
	"strings"
	"unicode"
)

type CommandID string

type Accelerator struct {
	Key       event.Key
	Modifiers event.Modifiers
	Label     string
}

func (a Accelerator) Matches(key event.Key, modifiers event.Modifiers) bool {
	return a.Key == key && normalizeMenuModifiers(a.Modifiers) == normalizeMenuModifiers(modifiers)
}

func (a Accelerator) DisplayLabel() string {
	if a.Label != "" {
		return a.Label
	}

	parts := make([]string, 0, 4)
	modifiers := normalizeMenuModifiers(a.Modifiers)
	if modifiers&event.ModCtrl != 0 {
		parts = append(parts, "Ctrl")
	}
	if modifiers&event.ModAlt != 0 {
		parts = append(parts, "Alt")
	}
	if modifiers&event.ModShift != 0 {
		parts = append(parts, "Shift")
	}

	keyName := menuKeyName(a.Key)
	if keyName == "" {
		return strings.Join(parts, "+")
	}
	parts = append(parts, keyName)
	return strings.Join(parts, "+")
}

type MenuItem struct {
	ID        CommandID
	Text      string
	Enabled   bool
	Checked   bool
	Separator bool
	Submenu   *Menu
	Shortcut  *Accelerator
}

func NewMenuItem(id CommandID, text string, shortcut *Accelerator) *MenuItem {
	return &MenuItem{
		ID:       id,
		Text:     text,
		Enabled:  true,
		Shortcut: shortcut,
	}
}

func NewSubmenuItem(text string, submenu *Menu) *MenuItem {
	return &MenuItem{
		Text:    text,
		Enabled: true,
		Submenu: submenu,
	}
}

func NewSeparator() *MenuItem {
	return &MenuItem{Separator: true}
}

func (i *MenuItem) DisplayText() string {
	if i == nil {
		return ""
	}
	return MenuDisplayText(i.Text)
}

func (i *MenuItem) Mnemonic() (rune, bool) {
	if i == nil {
		return 0, false
	}
	return menuMnemonic(i.Text)
}

func (i *MenuItem) ShortcutLabel() string {
	if i == nil || i.Shortcut == nil {
		return ""
	}
	return i.Shortcut.DisplayLabel()
}

func (i *MenuItem) Selectable() bool {
	return i != nil && !i.Separator && i.Enabled
}

type Menu struct {
	Items []*MenuItem
}

func NewMenu(items ...*MenuItem) *Menu {
	return &Menu{Items: append([]*MenuItem(nil), items...)}
}

func (m *Menu) FindByMnemonic(key event.Key) int {
	target, ok := menuKeyRune(key)
	if !ok {
		return -1
	}

	for i, item := range m.Items {
		if item == nil || !item.Selectable() {
			continue
		}
		mnemonic, ok := item.Mnemonic()
		if ok && unicode.ToUpper(mnemonic) == target {
			return i
		}
	}
	return -1
}

func (m *Menu) FindAccelerator(key event.Key, modifiers event.Modifiers) (*MenuItem, bool) {
	return findMenuAccelerator(m.Items, key, modifiers)
}

type MenuBar struct {
	Items []*MenuItem
}

func NewMenuBar(items ...*MenuItem) *MenuBar {
	return &MenuBar{Items: append([]*MenuItem(nil), items...)}
}

func (m *MenuBar) FindTopLevelByMnemonic(key event.Key) int {
	target, ok := menuKeyRune(key)
	if !ok {
		return -1
	}

	for i, item := range m.Items {
		if item == nil || !item.Selectable() {
			continue
		}
		mnemonic, ok := item.Mnemonic()
		if ok && unicode.ToUpper(mnemonic) == target {
			return i
		}
	}
	return -1
}

func (m *MenuBar) FindAccelerator(key event.Key, modifiers event.Modifiers) (*MenuItem, bool) {
	return findMenuAccelerator(m.Items, key, modifiers)
}

func MenuDisplayText(text string) string {
	runes := []rune(text)
	out := make([]rune, 0, len(runes))
	for i := 0; i < len(runes); i++ {
		if runes[i] != '&' {
			out = append(out, runes[i])
			continue
		}
		if i+1 < len(runes) && runes[i+1] == '&' {
			out = append(out, '&')
			i++
		}
	}
	return string(out)
}

func findMenuAccelerator(items []*MenuItem, key event.Key, modifiers event.Modifiers) (*MenuItem, bool) {
	for _, item := range items {
		if item == nil || item.Separator {
			continue
		}
		if item.Submenu != nil {
			if match, ok := findMenuAccelerator(item.Submenu.Items, key, modifiers); ok {
				return match, true
			}
		}
		if item.Selectable() && item.Shortcut != nil && item.Shortcut.Matches(key, modifiers) {
			return item, true
		}
	}
	return nil, false
}

func menuMnemonic(text string) (rune, bool) {
	runes := []rune(text)
	for i := 0; i+1 < len(runes); i++ {
		if runes[i] != '&' {
			continue
		}
		if runes[i+1] == '&' {
			i++
			continue
		}
		return unicode.ToUpper(runes[i+1]), true
	}
	return 0, false
}

func menuKeyRune(key event.Key) (rune, bool) {
	switch key {
	case event.KeyA:
		return 'A', true
	case event.KeyB:
		return 'B', true
	case event.KeyC:
		return 'C', true
	case event.KeyD:
		return 'D', true
	case event.KeyE:
		return 'E', true
	case event.KeyF:
		return 'F', true
	case event.KeyG:
		return 'G', true
	case event.KeyH:
		return 'H', true
	case event.KeyI:
		return 'I', true
	case event.KeyJ:
		return 'J', true
	case event.KeyK:
		return 'K', true
	case event.KeyL:
		return 'L', true
	case event.KeyM:
		return 'M', true
	case event.KeyN:
		return 'N', true
	case event.KeyO:
		return 'O', true
	case event.KeyP:
		return 'P', true
	case event.KeyQ:
		return 'Q', true
	case event.KeyR:
		return 'R', true
	case event.KeyS:
		return 'S', true
	case event.KeyT:
		return 'T', true
	case event.KeyU:
		return 'U', true
	case event.KeyV:
		return 'V', true
	case event.KeyW:
		return 'W', true
	case event.KeyX:
		return 'X', true
	case event.KeyY:
		return 'Y', true
	case event.KeyZ:
		return 'Z', true
	default:
		return 0, false
	}
}

func menuKeyName(key event.Key) string {
	if letter, ok := menuKeyRune(key); ok {
		return string(letter)
	}

	switch key {
	case event.KeyEscape:
		return "Esc"
	case event.KeyEnter:
		return "Enter"
	case event.KeySpace:
		return "Space"
	case event.KeyTab:
		return "Tab"
	case event.KeyBackspace:
		return "Backspace"
	case event.KeyDelete:
		return "Del"
	case event.KeyLeft:
		return "Left"
	case event.KeyRight:
		return "Right"
	case event.KeyUp:
		return "Up"
	case event.KeyDown:
		return "Down"
	case event.KeyHome:
		return "Home"
	case event.KeyEnd:
		return "End"
	case event.KeyPageUp:
		return "PgUp"
	case event.KeyPageDown:
		return "PgDn"
	default:
		return ""
	}
}

func normalizeMenuModifiers(modifiers event.Modifiers) event.Modifiers {
	return modifiers & (event.ModShift | event.ModCtrl | event.ModAlt)
}
