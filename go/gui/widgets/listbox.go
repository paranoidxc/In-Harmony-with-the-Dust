package widgets

import (
	"classicui/event"
	"classicui/geom"
	"classicui/widget"
	"time"
)

type ListBox struct {
	widget.BaseWidget
	items         []string
	selection     selectionModel[int]
	selectionOpts ListBoxSelectionOptions
	dragAnchor    int
	topIndex      int
	pressedIndex  int
	pressedBlank  bool
	pressedStart  geom.Point
	pressedPoint  geom.Point
	pressedMods   event.Modifiers
	dragSelecting bool
	focused       bool
	rowHeight     int
	scrollbarSize int
	scrollbar     *ScrollBar
	onChange      func(int, string)
	onActivate    func(int, string)
	now           func() time.Time
	lastDragTick  time.Time
}

func NewListBox(id string, bounds geom.Rect) *ListBox {
	list := &ListBox{
		BaseWidget:    widget.NewBase(id, bounds),
		selection:     newSelectionModel[int](),
		selectionOpts: DefaultListBoxSelectionOptions(),
		dragAnchor:    -1,
		pressedIndex:  -1,
		rowHeight:     16,
		scrollbarSize: 16,
		now:           time.Now,
	}
	list.scrollbar = NewScrollBar(id+".scrollbar", geom.Rect{})
	list.scrollbar.OnChange(func(value int) {
		list.topIndex = value
	})
	return list
}

func (l *ListBox) SetItems(items []string) {
	l.items = append([]string(nil), items...)
	if len(l.items) == 0 {
		l.selection.Clear()
		l.topIndex = 0
		l.scrollbar.SetRange(0, 1)
		return
	}
	selected := l.SelectedIndex()
	if selected >= len(l.items) {
		selected = len(l.items) - 1
	}
	if selected < 0 {
		selected = 0
	}
	l.selection.SetLead(selected)
	if anchor := l.anchorIndex(); anchor >= 0 {
		l.selection.SetAnchor(clampInt(anchor, 0, len(l.items)-1))
	}
	l.dropInvalidSelection()
	if l.selection.Count() == 0 {
		l.selection.SelectOnly(selected)
	}
	l.syncScrollBar()
}

func (l *ListBox) Items() []string {
	return append([]string(nil), l.items...)
}

func (l *ListBox) SelectedIndex() int {
	if selected, ok := l.selection.Lead(); ok {
		return selected
	}
	return -1
}

func (l *ListBox) SelectedIndices() []int {
	if l.selection.Count() == 0 {
		return nil
	}
	out := make([]int, 0, l.selection.Count())
	for index := range l.selection.selectedSet {
		if index >= 0 && index < len(l.items) {
			out = append(out, index)
		}
	}
	for i := 0; i < len(out); i++ {
		for j := i + 1; j < len(out); j++ {
			if out[j] < out[i] {
				out[i], out[j] = out[j], out[i]
			}
		}
	}
	return out
}

func (l *ListBox) SetSelectedIndex(index int) {
	l.setSelectedIndex(index, true)
}

func (l *ListBox) SetSelectedIndexSilent(index int) {
	l.setSelectedIndex(index, false)
}

func (l *ListBox) OnChange(fn func(int, string)) {
	l.onChange = fn
}

func (l *ListBox) OnActivate(fn func(int, string)) {
	l.onActivate = fn
}

func (l *ListBox) SelectionOptions() ListBoxSelectionOptions {
	return l.selectionOpts
}

func (l *ListBox) SetSelectionOptions(options ListBoxSelectionOptions) {
	l.selectionOpts = options
	if !l.selectionBehavior().allowsMultiSelect() {
		selected := l.SelectedIndex()
		if selected >= 0 {
			l.selection.SelectOnly(selected)
			return
		}
		if first := l.firstSelected(); first >= 0 {
			l.selection.SelectOnly(first)
			return
		}
		l.selection.Clear()
	}
}

func (l *ListBox) SetMultiSelect(enabled bool) {
	if l.selectionBehavior().allowsMultiSelect() == enabled {
		return
	}
	options := l.selectionOpts
	options.MultiSelect = enabled
	l.SetSelectionOptions(options)
}

func (l *ListBox) Paint(ctx PaintContext) error {
	if !l.Visible() {
		return nil
	}

	lineHeight := 14
	if ctx.Text != nil {
		lineHeight = ctx.Text.LineHeight()
	}
	l.rowHeight = maxInt(lineHeight+2, 16)
	l.scrollbarSize = maxInt(ctx.Theme.Metrics.ScrollbarSize, 16)
	l.syncScrollBar()

	rect := ctx.BoundsFor(l)
	ctx.Canvas.FillRect(rect, ctx.Theme.Colors.Face)
	ctx.Canvas.DrawDoubleBevel(rect, ctx.Theme.Colors.Shadow, ctx.Theme.Colors.Lightest, ctx.Theme.Colors.DarkShadow, ctx.Theme.Colors.Light)

	itemsRect := l.itemsRect(rect)
	ctx.Canvas.FillRect(itemsRect, ctx.Theme.Colors.Window)
	ctx.Canvas.PushClip(itemsRect)
	for row := 0; row < l.visibleRows(); row++ {
		index := l.topIndex + row
		if index >= len(l.items) {
			break
		}
		itemRect := geom.Rect{
			X: itemsRect.X + 1,
			Y: itemsRect.Y + row*l.rowHeight,
			W: itemsRect.W - 2,
			H: l.rowHeight,
		}
		if l.isSelected(index) {
			ctx.Canvas.FillRect(itemRect, ctx.Theme.Colors.Highlight)
			if ctx.Text != nil {
				if err := ctx.Text.DrawString(ctx.Canvas, geom.Point{X: itemRect.X + 2, Y: itemRect.Y + 1}, l.items[index], ctx.Theme.Colors.HighlightText); err != nil {
					ctx.Canvas.PopClip()
					return err
				}
			}
		} else if ctx.Text != nil {
			if err := ctx.Text.DrawString(ctx.Canvas, geom.Point{X: itemRect.X + 2, Y: itemRect.Y + 1}, l.items[index], ctx.Theme.Colors.WindowText); err != nil {
				ctx.Canvas.PopClip()
				return err
			}
		}
	}
	if rect, ok := l.marqueeRect(); ok {
		paintSelectionMarquee(ctx, rect)
	}
	ctx.Canvas.PopClip()

	sbCtx := ctx
	sbCtx.Origin = geom.Point{X: rect.X, Y: rect.Y}
	if err := l.scrollbar.Paint(sbCtx); err != nil {
		return err
	}

	if l.focused {
		ctx.Canvas.DrawFocusRect(itemsRect.Inset(1), ctx.Theme.Colors.DarkShadow)
	}
	return nil
}

func (l *ListBox) MouseEnter(EventContext) {}
func (l *ListBox) MouseLeave(EventContext) {}

func (l *ListBox) MouseDown(ctx EventContext, ev event.MouseButtonEvent, local geom.Point) {
	if ev.Button != event.MouseButtonLeft || !LocalContains(l, local) {
		return
	}
	if l.scrollbar.Bounds().Contains(local) {
		l.scrollbar.MouseDown(ctx, ev, geom.Point{X: local.X - l.scrollbar.Bounds().X, Y: local.Y - l.scrollbar.Bounds().Y})
		l.topIndex = l.scrollbar.Value()
		ctx.Invalidate(l)
		return
	}
	itemsRect := l.itemsRect(LocalRect(l))
	if !itemsRect.Contains(local) {
		return
	}
	ctx.SetFocus(l)
	mods := l.selectionBehavior().normalizeModifiers(ev.Modifiers)
	index, ok := l.itemIndexAtPoint(local)
	if !ok {
		hadSelection := l.SelectedIndex() != -1 || l.selection.Count() != 0
		l.pressedIndex = -1
		l.pressedBlank = true
		l.pressedStart = local
		l.pressedPoint = local
		l.pressedMods = mods
		l.dragSelecting = false
		l.dragAnchor = -1
		l.captureDragBaseSelection()
		if mods&(event.ModCtrl|event.ModShift) == 0 {
			l.setSelectedIndex(-1, false)
			if hadSelection {
				ctx.Invalidate(l)
			}
		}
		return
	}
	l.pressedIndex = index
	l.pressedBlank = false
	l.pressedStart = local
	l.pressedPoint = local
	l.pressedMods = mods
	l.dragSelecting = false
	l.dragAnchor = l.anchorIndex()
	l.captureDragBaseSelection()
	switch {
	case mods&event.ModShift != 0:
		l.selectRangeTo(index, true)
	case mods&event.ModCtrl != 0:
		l.toggleSelection(index, true)
	default:
		l.setSelectedIndex(index, true)
	}
	ctx.Invalidate(l)
}

func (l *ListBox) MouseUp(ctx EventContext, ev event.MouseButtonEvent, local geom.Point) {
	if l.scrollbar.dragging || l.scrollbar.Bounds().Contains(local) {
		l.scrollbar.MouseUp(ctx, ev, geom.Point{X: local.X - l.scrollbar.Bounds().X, Y: local.Y - l.scrollbar.Bounds().Y})
		l.topIndex = l.scrollbar.Value()
		ctx.Invalidate(l)
		l.clearPressed()
		clear(l.selection.dragBaseSet)
		return
	}
	itemsRect := l.itemsRect(LocalRect(l))
	if !l.dragSelecting && !l.pressedBlank && itemsRect.Contains(local) {
		index, ok := l.itemIndexAtPoint(local)
		if ok && index == l.pressedIndex && l.onActivate != nil {
			l.onActivate(index, l.items[index])
		}
	}
	l.clearPressed()
	clear(l.selection.dragBaseSet)
}

func (l *ListBox) MouseWheel(ctx EventContext, ev event.MouseWheel, _ geom.Point) bool {
	if len(l.items) == 0 {
		return true
	}
	l.topIndex = clampInt(l.topIndex-ev.Delta, 0, l.maxTopIndex())
	l.scrollbar.SetValue(l.topIndex)
	ctx.Invalidate(l)
	return true
}

func (l *ListBox) MouseMove(ctx EventContext, local geom.Point) {
	if l.scrollbar.dragging || l.scrollbar.Bounds().Contains(local) {
		l.scrollbar.MouseMove(ctx, geom.Point{X: local.X - l.scrollbar.Bounds().X, Y: local.Y - l.scrollbar.Bounds().Y})
		l.topIndex = l.scrollbar.Value()
		ctx.Invalidate(l)
		return
	}
	if l.pressedIndex >= 0 || l.pressedBlank {
		startedDrag := false
		if !l.dragSelecting {
			dx := local.X - l.pressedPoint.X
			dy := local.Y - l.pressedPoint.Y
			if dx*dx+dy*dy >= 9 {
				l.dragSelecting = true
				startedDrag = true
			}
		}
		l.pressedPoint = local
		if l.dragSelecting && (l.updateDragSelection(local, true) || startedDrag) {
			ctx.Invalidate(l)
		}
	}
}

func (l *ListBox) KeyDown(ctx EventContext, ev event.KeyEvent) bool {
	if !l.Enabled() {
		return false
	}
	if len(l.items) == 0 {
		return true
	}
	if l.selectionBehavior().toggleLeadShortcut(ev) {
		if selected := l.SelectedIndex(); selected >= 0 {
			if l.toggleSelection(selected, true) {
				ctx.Invalidate(l)
			}
			return true
		}
		if index, ok := l.recoveryIndexForKey(ev.Key); ok && l.toggleSelection(index, true) {
			ctx.Invalidate(l)
		}
		return true
	}
	if l.selectionBehavior().selectAllShortcut(ev) {
		if l.selectionBehavior().allowsMultiSelect() {
			l.selectAll(true)
			ctx.Invalidate(l)
		}
		return true
	}

	next := l.SelectedIndex()
	if next < 0 {
		if index, ok := l.recoveryIndexForKey(ev.Key); ok {
			if l.selectionBehavior().extendRange(ev.Modifiers) {
				l.selectRangeTo(index, true)
			} else {
				l.setSelectedIndex(index, true)
			}
			l.ensureVisible(l.SelectedIndex())
			l.scrollbar.SetValue(l.topIndex)
			ctx.Invalidate(l)
			return true
		}
		return false
	}

	switch ev.Key {
	case event.KeyUp:
		next--
	case event.KeyDown:
		next++
	case event.KeyHome:
		next = 0
	case event.KeyEnd:
		next = len(l.items) - 1
	case event.KeyPageUp:
		next -= maxInt(l.visibleRows()-1, 1)
	case event.KeyPageDown:
		next += maxInt(l.visibleRows()-1, 1)
	default:
		return false
	}

	if l.selectionBehavior().extendRange(ev.Modifiers) {
		l.selectRangeTo(next, true)
	} else {
		l.setSelectedIndex(next, true)
	}
	l.ensureVisible(l.SelectedIndex())
	l.scrollbar.SetValue(l.topIndex)
	ctx.Invalidate(l)
	return true
}

func (l *ListBox) CanFocus() bool {
	return l.Visible() && l.Enabled()
}

func (l *ListBox) SetFocused(focused bool) {
	l.focused = focused
}

func (l *ListBox) Focused() bool {
	return l.focused
}

func (l *ListBox) Tick(ctx EventContext, now time.Time) bool {
	if !l.dragSelecting || len(l.items) == 0 {
		return false
	}
	if !l.lastDragTick.IsZero() && now.Sub(l.lastDragTick) < 50*time.Millisecond {
		return false
	}
	l.lastDragTick = now
	pointer := l.pressedPoint
	if pointer.Y < 0 {
		if l.topIndex == 0 {
			return false
		}
		l.topIndex = clampInt(l.topIndex-1, 0, l.maxTopIndex())
	} else if pointer.Y >= LocalRect(l).H {
		if l.topIndex == l.maxTopIndex() {
			return false
		}
		l.topIndex = clampInt(l.topIndex+1, 0, l.maxTopIndex())
	} else {
		return false
	}
	l.scrollbar.SetValue(l.topIndex)
	l.updateDragSelection(pointer, false)
	ctx.Invalidate(l)
	return true
}

func (l *ListBox) itemsRect(rect geom.Rect) geom.Rect {
	return geom.Rect{
		X: rect.X + 2,
		Y: rect.Y + 2,
		W: maxInt(rect.W-l.scrollbarSize-4, 0),
		H: maxInt(rect.H-4, 0),
	}
}

func (l *ListBox) syncScrollBar() {
	rect := LocalRect(l)
	l.scrollbar.SetBounds(geom.Rect{
		X: maxInt(rect.W-l.scrollbarSize-1, 1),
		Y: 1,
		W: maxInt(l.scrollbarSize, 1),
		H: maxInt(rect.H-2, 0),
	})
	l.scrollbar.SetRange(l.maxTopIndex(), maxInt(l.visibleRows(), 1))
	l.scrollbar.SetValue(l.topIndex)
}

func (l *ListBox) visibleRows() int {
	rows := l.itemsRect(LocalRect(l)).H / maxInt(l.rowHeight, 1)
	return maxInt(rows, 1)
}

func (l *ListBox) maxTopIndex() int {
	return maxInt(len(l.items)-l.visibleRows(), 0)
}

func (l *ListBox) ensureVisible(index int) {
	if index < 0 {
		return
	}
	if index < l.topIndex {
		l.topIndex = index
	}
	bottom := l.topIndex + l.visibleRows() - 1
	if index > bottom {
		l.topIndex = index - l.visibleRows() + 1
	}
	l.topIndex = clampInt(l.topIndex, 0, l.maxTopIndex())
}

func (l *ListBox) setSelectedIndex(index int, notify bool) {
	if len(l.items) == 0 {
		l.selection.Clear()
		return
	}
	if index < 0 {
		if l.SelectedIndex() == -1 {
			return
		}
		l.selection.Clear()
		return
	}
	index = clampInt(index, 0, len(l.items)-1)
	if index == l.SelectedIndex() {
		if !l.isSelected(index) {
			l.selection.EnsureLeadSelected()
		}
		return
	}
	l.selection.SelectOnly(index)
	l.ensureVisible(index)
	if notify && l.onChange != nil {
		l.onChange(index, l.items[index])
	}
}

func (l *ListBox) isSelected(index int) bool {
	return l.selection.Contains(index)
}

func (l *ListBox) selectOnly(index int) {
	if index >= 0 && index < len(l.items) {
		l.selection.SelectOnly(index)
	}
}

func (l *ListBox) toggleSelection(index int, notify bool) bool {
	if index < 0 || index >= len(l.items) {
		return false
	}
	if !l.selection.Toggle(index, l.selectionOrder()) {
		return false
	}
	if selected := l.SelectedIndex(); selected >= 0 {
		l.ensureVisible(selected)
		if notify && l.onChange != nil {
			l.onChange(selected, l.items[selected])
		}
	}
	return true
}

func (l *ListBox) selectRangeTo(index int, notify bool) bool {
	if len(l.items) == 0 {
		return false
	}
	index = clampInt(index, 0, len(l.items)-1)
	if !l.selection.SelectRange(l.selectionOrder(), index) {
		return false
	}
	l.ensureVisible(index)
	if notify && l.onChange != nil {
		l.onChange(index, l.items[index])
	}
	return true
}

func (l *ListBox) selectAll(notify bool) bool {
	if len(l.items) == 0 {
		return false
	}
	changed := l.selection.SelectAll(l.selectionOrder())
	if notify && changed {
		if selected := l.SelectedIndex(); selected >= 0 && l.onChange != nil {
			l.onChange(selected, l.items[selected])
		}
	}
	return changed
}

func (l *ListBox) firstSelected() int {
	if index, ok := l.selection.firstSelectedInOrder(l.selectionOrder()); ok {
		return index
	}
	if len(l.items) > 0 {
		return 0
	}
	return -1
}

func (l *ListBox) dropInvalidSelection() {
	l.selection.DropInvalid(func(index int) bool {
		return index >= 0 && index < len(l.items)
	})
}

func (l *ListBox) selectionOrder() selectionOrder[int] {
	return selectionOrder[int]{
		Len:    func() int { return len(l.items) },
		ItemAt: func(index int) int { return index },
		IndexOf: func(index int) int {
			if index < 0 || index >= len(l.items) {
				return -1
			}
			return index
		},
	}
}

func (l *ListBox) visibleSelectionOrder() selectionOrder[int] {
	start := l.topIndex
	count := minInt(l.visibleRows(), maxInt(len(l.items)-start, 0))
	return selectionOrder[int]{
		Len:    func() int { return count },
		ItemAt: func(offset int) int { return start + offset },
		IndexOf: func(index int) int {
			if index < start || index >= start+count {
				return -1
			}
			return index - start
		},
	}
}

func (l *ListBox) anchorIndex() int {
	if anchor, ok := l.selection.Anchor(); ok {
		return anchor
	}
	return -1
}

func (l *ListBox) recentIndex() int {
	if index, ok := l.selection.Recent(); ok && index >= 0 && index < len(l.items) {
		return index
	}
	return -1
}

func (l *ListBox) recoveryIndexForKey(key event.Key) (int, bool) {
	if len(l.items) == 0 {
		return -1, false
	}
	switch key {
	case event.KeyHome:
		return 0, true
	case event.KeyEnd:
		return len(l.items) - 1, true
	case event.KeyUp, event.KeyDown, event.KeyPageUp, event.KeyPageDown, event.KeySpace:
		if l.selectionBehavior().allowsRecentRecovery() {
			if recent := l.recentIndex(); recent >= 0 {
				return recent, true
			}
		}
		return 0, true
	default:
		return -1, false
	}
}

func (l *ListBox) updateDragSelection(local geom.Point, clampToItems bool) bool {
	if len(l.items) == 0 {
		return false
	}
	if l.pressedBlank {
		if !l.selectionBehavior().allowsBlankDrag() {
			return false
		}
		return l.selectDragMarquee()
	}
	if l.pressedIndex < 0 {
		return false
	}
	currentPointer := l.pressedPoint
	index, ok := l.indexAtPoint(local, clampToItems)
	if !ok {
		return false
	}
	l.pressedPoint = currentPointer
	if !l.selectionBehavior().allowsMultiSelect() {
		previous := l.SelectedIndex()
		l.setSelectedIndex(index, false)
		return previous != l.SelectedIndex()
	}
	switch {
	case l.pressedMods&event.ModShift != 0:
		return l.selectRangeTo(index, true)
	case l.pressedMods&event.ModCtrl != 0:
		return l.selectDragUnion(index)
	default:
		return l.selectDragRange(index)
	}
}

func (l *ListBox) selectDragMarquee() bool {
	rect, ok := l.marqueeRect()
	if !ok {
		return false
	}
	return l.selection.ApplyMarquee(l.visibleSelectionOrder(), func(index int) bool {
		rowRect, ok := l.rowRect(index)
		if !ok || rowRect.Empty() {
			return false
		}
		_, hit := geom.Intersect(rect, rowRect)
		return hit
	}, l.pressedMods&(event.ModCtrl|event.ModShift) != 0)
}

func (l *ListBox) rowRect(index int) (geom.Rect, bool) {
	if index < l.topIndex || index >= len(l.items) {
		return geom.Rect{}, false
	}
	row := index - l.topIndex
	if row < 0 || row >= l.visibleRows() {
		return geom.Rect{}, false
	}
	itemsRect := l.itemsRect(LocalRect(l))
	return geom.Rect{
		X: itemsRect.X + 1,
		Y: itemsRect.Y + row*l.rowHeight,
		W: itemsRect.W - 2,
		H: l.rowHeight,
	}, true
}

func (l *ListBox) indexAtPoint(local geom.Point, clampToItems bool) (int, bool) {
	itemsRect := l.itemsRect(LocalRect(l))
	if itemsRect.W <= 0 || itemsRect.H <= 0 {
		return -1, false
	}
	if !clampToItems && !itemsRect.Contains(local) && local.Y >= itemsRect.Y && local.Y < itemsRect.Bottom() {
		return -1, false
	}
	y := local.Y
	if clampToItems {
		y = clampInt(y, itemsRect.Y, itemsRect.Bottom()-1)
	} else if y < itemsRect.Y || y >= itemsRect.Bottom() {
		return -1, false
	}
	row := (y - itemsRect.Y) / maxInt(l.rowHeight, 1)
	index := clampInt(l.topIndex+row, 0, len(l.items)-1)
	return index, true
}

func (l *ListBox) itemIndexAtPoint(local geom.Point) (int, bool) {
	itemsRect := l.itemsRect(LocalRect(l))
	if !itemsRect.Contains(local) {
		return -1, false
	}
	row := (local.Y - itemsRect.Y) / maxInt(l.rowHeight, 1)
	index := l.topIndex + row
	if index < 0 || index >= len(l.items) {
		return -1, false
	}
	return index, true
}

func (l *ListBox) selectDragRange(index int) bool {
	anchor := l.pressedIndex
	if anchor < 0 || anchor >= len(l.items) {
		anchor = index
	}
	return l.selection.SelectDragRange(l.selectionOrder(), anchor, index)
}

func (l *ListBox) selectDragUnion(index int) bool {
	startAnchor := l.dragAnchor
	if startAnchor < 0 || startAnchor >= len(l.items) {
		startAnchor = l.pressedIndex
	}
	if startAnchor < 0 || startAnchor >= len(l.items) {
		startAnchor = index
	}
	return l.selection.SelectDragUnion(l.selectionOrder(), startAnchor, index)
}

func (l *ListBox) captureDragBaseSelection() {
	l.selection.CaptureDragBase()
}

func (l *ListBox) marqueeRect() (geom.Rect, bool) {
	if !l.selectionBehavior().allowsBlankDrag() || !l.dragSelecting || (l.pressedIndex < 0 && !l.pressedBlank) {
		return geom.Rect{}, false
	}
	return dragMarqueeRect(l.itemsRect(LocalRect(l)), l.pressedStart, l.pressedPoint)
}

func (l *ListBox) clearPressed() {
	l.pressedIndex = -1
	l.pressedBlank = false
	l.dragSelecting = false
}

func (l *ListBox) selectionBehavior() selectionBehavior {
	return l.selectionOpts.behavior()
}
