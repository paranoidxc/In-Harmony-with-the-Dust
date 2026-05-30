package desktop

import (
	"classicui/event"
	"classicui/geom"
	"classicui/paint"
	"classicui/uicolor"
	"classicui/widgets"
)

const (
	menuPopupHPadding   = 6
	menuPopupCheckWidth = 16
	menuPopupArrowWidth = 12
	menuPopupGap        = 16
	menuSeparatorHeight = 7
)

type popupMenuState struct {
	menu     *widgets.Menu
	rect     geom.Rect
	items    []popupMenuLayout
	selected int
}

type popupMenuLayout struct {
	item *widgets.MenuItem
	rect geom.Rect
}

func (p *popupMenuState) itemAt(point geom.Point) int {
	for i, layout := range p.items {
		if layout.rect.Contains(point) {
			return i
		}
	}
	return -1
}

func (p *popupMenuState) Bounds() geom.Rect {
	return p.rect
}

func (p *popupMenuState) Paint(d *Desktop, canvas *paint.Canvas) error {
	return d.paintPopupMenu(canvas, p)
}

func (d *Desktop) BindCommandHandler(handler func(*Window, widgets.CommandID)) {
	d.commandHandler = handler
}

func (d *Desktop) handleMenuMouseMove(point geom.Point) bool {
	if !d.menuMode || d.menuWindow == nil {
		return false
	}

	if d.menuWindow.HitTest(point, d.theme) == HitMenuBar {
		index := d.menuWindow.MenuBarItemAt(point, d.theme, d.text)
		if index >= 0 {
			d.activateMenuBar(d.menuWindow, index, true)
		}
		return true
	}

	for level := len(d.menuPopups) - 1; level >= 0; level-- {
		popup := d.menuPopups[level]
		if !popup.rect.Contains(point) {
			continue
		}

		index := popup.itemAt(point)
		if index < 0 || popup.items[index].item == nil || popup.items[index].item.Separator {
			d.setPopupSelection(level, -1)
			d.trimPopupStack(level + 1)
			return true
		}

		d.setPopupSelection(level, index)
		item := popup.items[index].item
		if item.Selectable() && item.Submenu != nil {
			d.openSubmenu(level, popup.items[index].rect, item.Submenu)
			return true
		}

		d.trimPopupStack(level + 1)
		return true
	}

	return true
}

func (d *Desktop) handleMenuMouseDown(point geom.Point) bool {
	if d.menuMode && d.menuWindow != nil {
		if d.menuWindow.HitTest(point, d.theme) == HitMenuBar {
			index := d.menuWindow.MenuBarItemAt(point, d.theme, d.text)
			if index >= 0 {
				d.activateMenuBar(d.menuWindow, index, true)
				return true
			}
		}

		for level := len(d.menuPopups) - 1; level >= 0; level-- {
			popup := d.menuPopups[level]
			if !popup.rect.Contains(point) {
				continue
			}
			index := popup.itemAt(point)
			if index >= 0 && popup.items[index].item != nil && !popup.items[index].item.Separator {
				d.setPopupSelection(level, index)
				if submenu := popup.items[index].item.Submenu; popup.items[index].item.Selectable() && submenu != nil {
					d.openSubmenu(level, popup.items[index].rect, submenu)
				} else {
					d.trimPopupStack(level + 1)
				}
			} else {
				d.setPopupSelection(level, -1)
				d.trimPopupStack(level + 1)
			}
			return true
		}

		d.closeMenus()
	}

	win, _ := d.windowAt(point)
	if win == nil || win.HitTest(point, d.theme) != HitMenuBar {
		return false
	}

	index := win.MenuBarItemAt(point, d.theme, d.text)
	if index < 0 {
		return false
	}

	d.activateWindow(win)
	d.activateMenuBar(win, index, true)
	d.setHoveredControl(nil, nil)
	return true
}

func (d *Desktop) handleMenuMouseUp(e event.MouseButtonEvent) bool {
	if !d.menuMode || d.menuWindow == nil {
		return false
	}
	if e.Button != event.MouseButtonLeft {
		return true
	}
	point := e.Position

	if d.menuWindow.HitTest(point, d.theme) == HitMenuBar {
		return true
	}

	for level := len(d.menuPopups) - 1; level >= 0; level-- {
		popup := d.menuPopups[level]
		if !popup.rect.Contains(point) {
			continue
		}

		index := popup.itemAt(point)
		if index < 0 {
			return true
		}

		item := popup.items[index].item
		if item == nil || !item.Selectable() {
			return true
		}
		if item.Submenu != nil {
			d.openSubmenu(level, popup.items[index].rect, item.Submenu)
			return true
		}

		win := d.menuWindow
		cmd := item.ID
		d.closeMenus()
		d.dispatchCommand(win, cmd)
		return true
	}

	d.closeMenus()
	return true
}

func (d *Desktop) handleMenuKeyDown(active *Window, e event.KeyEvent) bool {
	if !e.Repeat && (e.Key == event.KeyLeftAlt || e.Key == event.KeyRightAlt) {
		handled := false
		if d.menuMode {
			d.closeMenus()
			handled = true
		} else if active != nil && active.MenuBar() != nil && len(active.MenuBar().Items) > 0 {
			index := nextSelectableMenuBarIndex(active.MenuBar().Items, -1, 1)
			if index >= 0 {
				d.activateMenuBar(active, index, false)
				handled = true
			}
		}
		return handled
	}

	if active != nil && e.Modifiers&event.ModAlt != 0 && e.Modifiers&event.ModCtrl == 0 {
		if bar := active.MenuBar(); bar != nil {
			if index := bar.FindTopLevelByMnemonic(e.Key); index >= 0 {
				d.activateMenuBar(active, index, true)
				return true
			}
		}
	}

	if !d.menuMode || d.menuWindow == nil || d.menuWindow.MenuBar() == nil {
		return false
	}

	if index, level := d.findPopupMnemonic(e.Key); index >= 0 {
		d.setPopupSelection(level, index)
		item := d.menuPopups[level].items[index].item
		if item.Submenu != nil {
			d.openSubmenu(level, d.menuPopups[level].items[index].rect, item.Submenu)
		} else {
			win := d.menuWindow
			cmd := item.ID
			d.closeMenus()
			d.dispatchCommand(win, cmd)
		}
		return true
	}

	switch e.Key {
	case event.KeyEscape:
		d.closeMenus()
	case event.KeyLeft:
		if len(d.menuPopups) > 1 {
			d.trimPopupStack(len(d.menuPopups) - 1)
			return true
		}
		d.selectAdjacentTopLevel(-1)
	case event.KeyRight:
		if d.tryOpenSelectedSubmenu() {
			return true
		}
		d.selectAdjacentTopLevel(1)
	case event.KeyDown:
		if len(d.menuPopups) == 0 {
			d.openCurrentTopLevelPopup()
			return true
		}
		d.movePopupSelection(1)
	case event.KeyUp:
		if len(d.menuPopups) == 0 {
			d.openCurrentTopLevelPopup()
			return true
		}
		d.movePopupSelection(-1)
	case event.KeyEnter, event.KeySpace:
		if len(d.menuPopups) == 0 {
			d.openCurrentTopLevelPopup()
			return true
		}
		d.activateSelectedMenuItem()
	default:
		return true
	}
	return true
}

func (d *Desktop) dispatchAccelerator(win *Window, e event.KeyEvent) bool {
	if e.Repeat || win == nil || win.MenuBar() == nil {
		return false
	}

	item, ok := win.MenuBar().FindAccelerator(e.Key, e.Modifiers)
	if !ok {
		return false
	}

	d.dispatchCommand(win, item.ID)
	return true
}

func (d *Desktop) showContextMenu(win *Window, owner widgets.Control, anchor geom.Rect, menu *widgets.Menu) bool {
	if win == nil || owner == nil || menu == nil {
		return false
	}
	if d.menuMode {
		d.closeMenus()
	}
	ownerRect := d.controlScreenRect(win, owner)
	anchorRect := anchor.Move(ownerRect.X, ownerRect.Y)
	d.menuWindow = win
	d.menuMode = true
	d.trimPopupStack(0)
	popup := d.buildPopupMenu(menu, anchorRect, false)
	d.menuPopups = append(d.menuPopups, popup)
	d.pushOverlay(popup)
	return true
}

func (d *Desktop) activateMenuBar(win *Window, index int, openPopup bool) {
	if win == nil || win.MenuBar() == nil || index < 0 || index >= len(win.MenuBar().Items) {
		return
	}

	if d.menuWindow != nil && d.menuWindow != win {
		d.setMenuBarHighlight(d.menuWindow, -1)
	}

	d.menuWindow = win
	d.menuMode = true
	d.setMenuBarHighlight(win, index)

	if openPopup {
		d.openCurrentTopLevelPopup()
		return
	}
	d.trimPopupStack(0)
}

func (d *Desktop) openCurrentTopLevelPopup() {
	if d.menuWindow == nil || d.menuWindow.MenuBar() == nil {
		return
	}

	index := d.menuWindow.MenuBarActiveIndex()
	if index < 0 || index >= len(d.menuWindow.MenuBar().Items) {
		return
	}

	item := d.menuWindow.MenuBar().Items[index]
	d.trimPopupStack(0)
	if item == nil || !item.Selectable() {
		return
	}
	if item.Submenu == nil {
		win := d.menuWindow
		cmd := item.ID
		d.closeMenus()
		d.dispatchCommand(win, cmd)
		return
	}

	popup := d.buildPopupMenu(item.Submenu, d.menuWindow.MenuBarItemRect(index, d.theme, d.text), false)
	d.menuPopups = append(d.menuPopups, popup)
	d.pushOverlay(popup)
}

func (d *Desktop) openSubmenu(level int, anchorRect geom.Rect, menu *widgets.Menu) {
	if menu == nil {
		d.trimPopupStack(level + 1)
		return
	}

	popup := d.buildPopupMenu(menu, anchorRect, true)
	if len(d.menuPopups) > level+1 {
		existing := d.menuPopups[level+1]
		if existing.menu == popup.menu && existing.rect == popup.rect {
			return
		}
	}

	d.trimPopupStack(level + 1)
	d.menuPopups = append(d.menuPopups, popup)
	d.pushOverlay(popup)
}

func (d *Desktop) movePopupSelection(delta int) {
	level := len(d.menuPopups) - 1
	if level < 0 {
		return
	}

	popup := d.menuPopups[level]
	index := nextSelectablePopupIndex(popup.items, popup.selected, delta)
	if index < 0 {
		return
	}

	d.setPopupSelection(level, index)
	d.trimPopupStack(level + 1)
}

func (d *Desktop) activateSelectedMenuItem() {
	level := len(d.menuPopups) - 1
	if level < 0 {
		return
	}

	popup := d.menuPopups[level]
	index := popup.selected
	if index < 0 {
		index = nextSelectablePopupIndex(popup.items, -1, 1)
		if index < 0 {
			return
		}
		d.setPopupSelection(level, index)
	}

	item := popup.items[index].item
	if item == nil || !item.Selectable() {
		return
	}
	if item.Submenu != nil {
		d.openSubmenu(level, popup.items[index].rect, item.Submenu)
		return
	}

	win := d.menuWindow
	cmd := item.ID
	d.closeMenus()
	d.dispatchCommand(win, cmd)
}

func (d *Desktop) tryOpenSelectedSubmenu() bool {
	level := len(d.menuPopups) - 1
	if level < 0 {
		return false
	}

	popup := d.menuPopups[level]
	if popup.selected < 0 || popup.selected >= len(popup.items) {
		return false
	}

	item := popup.items[popup.selected].item
	if item == nil || !item.Selectable() || item.Submenu == nil {
		return false
	}

	d.openSubmenu(level, popup.items[popup.selected].rect, item.Submenu)
	return true
}

func (d *Desktop) selectAdjacentTopLevel(delta int) {
	if d.menuWindow == nil || d.menuWindow.MenuBar() == nil {
		return
	}

	index := nextSelectableMenuBarIndex(d.menuWindow.MenuBar().Items, d.menuWindow.MenuBarActiveIndex(), delta)
	if index < 0 {
		return
	}
	d.activateMenuBar(d.menuWindow, index, true)
}

func (d *Desktop) closeMenus() {
	if d.menuWindow != nil {
		d.InvalidateRect(d.menuWindow.MenuBarRect(d.theme))
		d.menuWindow.SetMenuBarActiveIndex(-1)
	}

	for _, popup := range d.menuPopups {
		d.removeOverlay(popup)
	}

	d.menuWindow = nil
	d.menuMode = false
	d.menuPopups = nil
}

func (d *Desktop) trimPopupStack(keep int) {
	if keep < 0 {
		keep = 0
	}
	if keep >= len(d.menuPopups) {
		return
	}

	for _, popup := range d.menuPopups[keep:] {
		d.removeOverlay(popup)
	}
	d.menuPopups = d.menuPopups[:keep]
}

func (d *Desktop) setPopupSelection(level, index int) {
	if level < 0 || level >= len(d.menuPopups) {
		return
	}
	if d.menuPopups[level].selected == index {
		return
	}
	d.menuPopups[level].selected = index
	d.InvalidateRect(d.menuPopups[level].rect)
}

func (d *Desktop) setMenuBarHighlight(win *Window, index int) {
	if win == nil || win.MenuBar() == nil {
		return
	}
	if win.MenuBarActiveIndex() == index {
		return
	}
	d.InvalidateRect(win.MenuBarRect(d.theme))
	win.SetMenuBarActiveIndex(index)
	d.InvalidateRect(win.MenuBarRect(d.theme))
}

func (d *Desktop) findPopupMnemonic(key event.Key) (int, int) {
	for level := len(d.menuPopups) - 1; level >= 0; level-- {
		index := d.menuPopups[level].menu.FindByMnemonic(key)
		if index >= 0 {
			return index, level
		}
	}
	if d.menuWindow == nil || d.menuWindow.MenuBar() == nil {
		return -1, -1
	}
	index := d.menuWindow.MenuBar().FindTopLevelByMnemonic(key)
	if index >= 0 {
		d.activateMenuBar(d.menuWindow, index, true)
		return -1, -1
	}
	return -1, -1
}

func (d *Desktop) buildPopupMenu(menu *widgets.Menu, anchorRect geom.Rect, submenu bool) *popupMenuState {
	lineHeight := d.menuLineHeight()
	rowHeight := max(d.theme.Metrics.MenuHeight, lineHeight+6)
	labelWidth := 0
	shortcutWidth := 0
	hasSubmenu := false
	height := 2

	for _, item := range menu.Items {
		if item == nil {
			continue
		}
		if item.Separator {
			height += menuSeparatorHeight
			continue
		}

		labelWidth = max(labelWidth, d.measureText(item.DisplayText()).W)
		shortcutWidth = max(shortcutWidth, d.measureText(item.ShortcutLabel()).W)
		hasSubmenu = hasSubmenu || item.Submenu != nil
		height += rowHeight
	}

	width := menuPopupHPadding*2 + menuPopupCheckWidth + labelWidth
	if shortcutWidth > 0 {
		width += menuPopupGap + shortcutWidth
	}
	if hasSubmenu {
		width += menuPopupArrowWidth
	}
	width = max(width, 96)
	height = max(height+2, rowHeight+4)

	placement := widgets.OverlayBelowStart
	if submenu {
		placement = widgets.OverlayRightTop
	}
	origin := d.fitOverlayOrigin(anchorRect, geom.Size{W: width, H: height}, placement)
	popup := &popupMenuState{
		menu:     menu,
		rect:     geom.Rect{X: origin.X, Y: origin.Y, W: width, H: height},
		selected: nextSelectableMenuBarIndex(menu.Items, -1, 1),
	}

	y := popup.rect.Y + 2
	for _, item := range menu.Items {
		if item == nil {
			continue
		}

		itemHeight := rowHeight
		if item.Separator {
			itemHeight = menuSeparatorHeight
		}

		popup.items = append(popup.items, popupMenuLayout{
			item: item,
			rect: geom.Rect{X: popup.rect.X + 2, Y: y, W: popup.rect.W - 4, H: itemHeight},
		})
		y += itemHeight
	}

	return popup
}

func (d *Desktop) paintPopupMenu(canvas *paint.Canvas, popup *popupMenuState) error {
	canvas.FillRect(popup.rect, d.theme.Colors.Face)
	canvas.DrawDoubleBevel(popup.rect, d.theme.Colors.Lightest, d.theme.Colors.DarkShadow, d.theme.Colors.Light, d.theme.Colors.Shadow)

	for i, layout := range popup.items {
		if layout.item == nil {
			continue
		}
		if layout.item.Separator {
			lineY := layout.rect.Y + layout.rect.H/2
			canvas.DrawHLine(layout.rect.X+2, lineY, max(layout.rect.W-4, 0), d.theme.Colors.Shadow)
			canvas.DrawHLine(layout.rect.X+2, lineY+1, max(layout.rect.W-4, 0), d.theme.Colors.Lightest)
			continue
		}

		selected := i == popup.selected && layout.item.Selectable()
		if selected {
			canvas.FillRect(layout.rect, d.theme.Colors.Highlight)
		}

		textColor := d.theme.Colors.WindowText
		if !layout.item.Enabled {
			textColor = d.theme.Colors.GrayText
		} else if selected {
			textColor = d.theme.Colors.HighlightText
		}

		checkRect := geom.Rect{
			X: layout.rect.X + 2,
			Y: layout.rect.Y + max((layout.rect.H-9)/2, 0),
			W: 9,
			H: 9,
		}
		if layout.item.Checked {
			drawMenuCheck(canvas, checkRect, textColor)
		}

		if d.text != nil {
			textY := layout.rect.Y + max((layout.rect.H-d.text.LineHeight())/2, 0)
			labelX := layout.rect.X + menuPopupCheckWidth + menuPopupHPadding
			if err := d.text.DrawString(canvas, geom.Point{X: labelX, Y: textY}, layout.item.DisplayText(), textColor); err != nil {
				return err
			}

			shortcut := layout.item.ShortcutLabel()
			if shortcut != "" {
				shortWidth := d.text.MeasureString(shortcut).W
				shortcutX := layout.rect.Right() - menuPopupHPadding - shortWidth
				if layout.item.Submenu != nil {
					shortcutX -= menuPopupArrowWidth
				}
				if err := d.text.DrawString(canvas, geom.Point{X: shortcutX, Y: textY}, shortcut, textColor); err != nil {
					return err
				}
			}
		}

		if layout.item.Submenu != nil {
			arrowColor := textColor
			drawMenuArrow(canvas, geom.Rect{
				X: layout.rect.Right() - menuPopupArrowWidth + 2,
				Y: layout.rect.Y + max((layout.rect.H-7)/2, 0),
				W: 7,
				H: 7,
			}, arrowColor)
		}
	}

	return nil
}

func (d *Desktop) measureText(text string) geom.Size {
	if d.text == nil {
		return geom.Size{W: len([]rune(text)) * 7, H: 14}
	}
	return d.text.MeasureString(text)
}

func (d *Desktop) menuLineHeight() int {
	if d.text == nil {
		return 14
	}
	return d.text.LineHeight()
}

func (d *Desktop) dispatchCommand(win *Window, cmd widgets.CommandID) {
	if cmd == "" || d.commandHandler == nil {
		return
	}
	d.commandHandler(win, cmd)
}

func nextSelectableMenuBarIndex(items []*widgets.MenuItem, start, delta int) int {
	if len(items) == 0 {
		return -1
	}

	if delta == 0 {
		delta = 1
	}

	index := start
	if index < 0 {
		if delta > 0 {
			index = 0
		} else {
			index = len(items) - 1
		}
	} else {
		index = (index + delta + len(items)) % len(items)
	}

	for tried := 0; tried < len(items); tried++ {
		if items[index] != nil && items[index].Selectable() {
			return index
		}
		index = (index + delta + len(items)) % len(items)
	}
	return -1
}

func nextSelectablePopupIndex(items []popupMenuLayout, start, delta int) int {
	if len(items) == 0 {
		return -1
	}

	if delta == 0 {
		delta = 1
	}

	index := start
	if index < 0 {
		if delta > 0 {
			index = 0
		} else {
			index = len(items) - 1
		}
	} else {
		index = (index + delta + len(items)) % len(items)
	}

	for tried := 0; tried < len(items); tried++ {
		item := items[index].item
		if item != nil && item.Selectable() {
			return index
		}
		index = (index + delta + len(items)) % len(items)
	}
	return -1
}

func drawMenuArrow(canvas *paint.Canvas, rect geom.Rect, color uicolor.RGBA) {
	midY := rect.Y + rect.H/2
	for col := 0; col < rect.W; col++ {
		x := rect.X + col
		top := midY - col
		bottom := midY + col
		for y := top; y <= bottom; y++ {
			if y >= rect.Y && y < rect.Bottom() {
				canvas.DrawPixel(x, y, color)
			}
		}
	}
}

func drawMenuCheck(canvas *paint.Canvas, rect geom.Rect, color uicolor.RGBA) {
	points := []geom.Point{
		{X: rect.X + 1, Y: rect.Y + 4},
		{X: rect.X + 2, Y: rect.Y + 5},
		{X: rect.X + 3, Y: rect.Y + 6},
		{X: rect.X + 4, Y: rect.Y + 5},
		{X: rect.X + 5, Y: rect.Y + 4},
		{X: rect.X + 6, Y: rect.Y + 3},
	}
	for _, point := range points {
		canvas.DrawPixel(point.X, point.Y, color)
		canvas.DrawPixel(point.X, point.Y+1, color)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
