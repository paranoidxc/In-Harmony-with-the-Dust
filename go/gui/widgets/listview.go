package widgets

import (
	"classicui/event"
	"classicui/geom"
	"classicui/widget"
	"time"
)

type ListViewColumn = HeaderColumn

type ListViewItem struct {
	Texts []string
	Data  any
}

type ListViewContextMenuInfo struct {
	Index   int
	Item    ListViewItem
	HasItem bool
}

type ListViewSelectionOptions = SelectionBehaviorOptions

func DefaultListViewSelectionOptions() ListViewSelectionOptions {
	return DefaultSelectionBehaviorOptions()
}

type ListView struct {
	widget.BaseWidget
	columns             []ListViewColumn
	items               []ListViewItem
	header              *HeaderControl
	scrollbar           *ScrollBar
	selection           selectionModel[int]
	selectionOpts       ListViewSelectionOptions
	dragAnchor          int
	topIndex            int
	pressedIndex        int
	pressedBlank        bool
	pressedStart        geom.Point
	pressedPoint        geom.Point
	pressedMods         event.Modifiers
	dragSelecting       bool
	focused             bool
	hotIndex            int
	rowHeight           int
	headerHeight        int
	scrollbarSize       int
	onChange            func(int, ListViewItem)
	onActivate          func(int, ListViewItem)
	onColumnClick       func(int)
	contextMenu         *Menu
	contextMenuProvider func(ListViewContextMenuInfo) *Menu
	lastClickIndex      int
	lastClickAt         time.Time
	now                 func() time.Time
	doubleClick         time.Duration
	lastDragTick        time.Time
}

func NewListView(id string, bounds geom.Rect, columns ...ListViewColumn) *ListView {
	view := &ListView{
		BaseWidget:     widget.NewBase(id, bounds),
		header:         NewHeaderControl(id+".header", geom.Rect{}),
		scrollbar:      NewScrollBar(id+".scrollbar", geom.Rect{}),
		selection:      newSelectionModel[int](),
		selectionOpts:  DefaultListViewSelectionOptions(),
		dragAnchor:     -1,
		pressedIndex:   -1,
		hotIndex:       -1,
		rowHeight:      16,
		headerHeight:   18,
		scrollbarSize:  16,
		now:            time.Now,
		doubleClick:    500 * time.Millisecond,
		lastClickIndex: -1,
	}
	view.header.SetParent(view)
	view.scrollbar.SetParent(view)
	view.scrollbar.OnChange(func(value int) {
		view.topIndex = value
	})
	view.header.OnColumnClick(func(index int) {
		if view.onColumnClick != nil {
			view.onColumnClick(index)
		}
	})
	view.header.OnColumnResize(func(index, width int) {
		if index < 0 || index >= len(view.columns) {
			return
		}
		view.columns[index].Width = width
	})
	view.header.OnColumnAutoFit(func(ctx EventContext, index int) {
		view.autoFitColumn(index, ctx.MeasureText)
	})
	view.SetColumns(columns...)
	view.header.now = func() time.Time {
		if view.now != nil {
			return view.now()
		}
		return time.Now()
	}
	return view
}

func (l *ListView) SetColumns(columns ...ListViewColumn) {
	l.columns = append([]ListViewColumn(nil), columns...)
	l.header.SetColumns(columns...)
	l.syncChrome()
}

func (l *ListView) Columns() []ListViewColumn {
	return append([]ListViewColumn(nil), l.columns...)
}

func (l *ListView) SetItems(items []ListViewItem) {
	l.items = append([]ListViewItem(nil), items...)
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
	l.syncChrome()
}

func (l *ListView) Items() []ListViewItem {
	return append([]ListViewItem(nil), l.items...)
}

func (l *ListView) SelectedIndex() int {
	if selected, ok := l.selection.Lead(); ok {
		return selected
	}
	return -1
}

func (l *ListView) SelectedIndices() []int {
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

func (l *ListView) SetSelectedIndex(index int) {
	l.setSelectedIndex(index, true)
}

func (l *ListView) SetSelectedIndexSilent(index int) {
	l.setSelectedIndex(index, false)
}

func (l *ListView) SetSelectionOptions(options ListViewSelectionOptions) {
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

func (l *ListView) SelectionOptions() ListViewSelectionOptions {
	return l.selectionOpts
}

func (l *ListView) SetMultiSelect(enabled bool) {
	if l.selectionBehavior().allowsMultiSelect() == enabled {
		return
	}
	options := l.selectionOpts
	options.MultiSelect = enabled
	l.SetSelectionOptions(options)
}

func (l *ListView) SetSortIndicator(index int, descending bool) {
	l.header.SetSortIndicator(index, descending)
}

func (l *ListView) OnChange(fn func(int, ListViewItem)) {
	l.onChange = fn
}

func (l *ListView) OnActivate(fn func(int, ListViewItem)) {
	l.onActivate = fn
}

func (l *ListView) OnColumnClick(fn func(int)) {
	l.onColumnClick = fn
}

func (l *ListView) SetContextMenu(menu *Menu) {
	l.contextMenu = menu
}

func (l *ListView) SetContextMenuProvider(fn func(ListViewContextMenuInfo) *Menu) {
	l.contextMenuProvider = fn
}

func (l *ListView) Paint(ctx PaintContext) error {
	if !l.Visible() {
		return nil
	}

	lineHeight := 14
	if ctx.Text != nil {
		lineHeight = ctx.Text.LineHeight()
	}
	l.rowHeight = maxInt(lineHeight+2, 16)
	l.headerHeight = maxInt(lineHeight+8, 18)
	l.scrollbarSize = maxInt(ctx.Theme.Metrics.ScrollbarSize, 16)
	l.syncChrome()

	rect := ctx.BoundsFor(l)
	ctx.Canvas.FillRect(rect, ctx.Theme.Colors.Face)
	ctx.Canvas.DrawDoubleBevel(rect, ctx.Theme.Colors.Shadow, ctx.Theme.Colors.Lightest, ctx.Theme.Colors.DarkShadow, ctx.Theme.Colors.Light)

	childCtx := ctx.Child(l)
	if err := l.header.Paint(childCtx); err != nil {
		return err
	}

	itemsRect := l.itemsRect(rect)
	ctx.Canvas.FillRect(itemsRect, ctx.Theme.Colors.Window)
	ctx.Canvas.PushClip(itemsRect)
	for row := 0; row < l.visibleRows(); row++ {
		index := l.topIndex + row
		if index >= len(l.items) {
			break
		}
		rowRect := geom.Rect{
			X: itemsRect.X + 1,
			Y: itemsRect.Y + row*l.rowHeight,
			W: itemsRect.W - 2,
			H: l.rowHeight,
		}
		selected := l.isSelected(index)
		if selected {
			ctx.Canvas.FillRect(rowRect, ctx.Theme.Colors.Highlight)
		} else if index == l.hotIndex {
			ctx.Canvas.FillRect(rowRect, blendColor(ctx.Theme.Colors.Window, ctx.Theme.Colors.Light))
		}
		l.paintRow(ctx, rowRect, l.items[index], selected)
	}
	if rect, ok := l.marqueeRect(); ok {
		paintSelectionMarquee(ctx, rect)
	}
	ctx.Canvas.PopClip()

	if err := l.scrollbar.Paint(childCtx); err != nil {
		return err
	}
	if l.focused {
		ctx.Canvas.DrawFocusRect(itemsRect.Inset(1), ctx.Theme.Colors.DarkShadow)
	}
	return nil
}

func (l *ListView) MouseEnter(EventContext) {}

func (l *ListView) MouseLeave(ctx EventContext) {
	if l.hotIndex < 0 {
		return
	}
	l.hotIndex = -1
	ctx.Invalidate(l)
}

func (l *ListView) MouseDown(ctx EventContext, ev event.MouseButtonEvent, local geom.Point) {
	if !LocalContains(l, local) {
		return
	}
	if ev.Button == event.MouseButtonRight {
		l.handleContextMenu(ctx, local)
		return
	}
	if ev.Button != event.MouseButtonLeft {
		return
	}
	l.syncChrome()
	if l.header.Bounds().Contains(local) {
		l.header.MouseDown(ctx, ev, geom.Point{X: local.X - l.header.Bounds().X, Y: local.Y - l.header.Bounds().Y})
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
	l.lastDragTick = time.Time{}
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

func (l *ListView) MouseUp(ctx EventContext, ev event.MouseButtonEvent, local geom.Point) {
	l.syncChrome()
	if l.header.Bounds().Contains(local) || l.header.Focused() {
		l.header.MouseUp(ctx, ev, geom.Point{X: local.X - l.header.Bounds().X, Y: local.Y - l.header.Bounds().Y})
	}
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
		if ok && index == l.pressedIndex && l.onActivate != nil && l.isDoubleClick(index) {
			l.onActivate(index, l.items[index])
		}
	}
	l.clearPressed()
	clear(l.selection.dragBaseSet)
}

func (l *ListView) MouseWheel(ctx EventContext, ev event.MouseWheel, _ geom.Point) bool {
	if len(l.items) == 0 {
		return true
	}
	l.topIndex = clampInt(l.topIndex-ev.Delta, 0, l.maxTopIndex())
	l.scrollbar.SetValue(l.topIndex)
	ctx.Invalidate(l)
	return true
}

func (l *ListView) MouseMove(ctx EventContext, local geom.Point) {
	l.syncChrome()
	if l.header.Bounds().Contains(local) {
		l.header.MouseMove(ctx, geom.Point{X: local.X - l.header.Bounds().X, Y: local.Y - l.header.Bounds().Y})
		return
	}
	if l.scrollbar.dragging || l.scrollbar.Bounds().Contains(local) {
		l.scrollbar.MouseMove(ctx, geom.Point{X: local.X - l.scrollbar.Bounds().X, Y: local.Y - l.scrollbar.Bounds().Y})
		l.topIndex = l.scrollbar.Value()
		ctx.Invalidate(l)
		return
	}
	hotIndex, _ := l.itemIndexAtPoint(local)
	if hotIndex != l.hotIndex {
		l.hotIndex = hotIndex
		ctx.Invalidate(l)
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

func (l *ListView) KeyDown(ctx EventContext, ev event.KeyEvent) bool {
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
	case event.KeyEnter:
		if selected := l.SelectedIndex(); selected >= 0 && selected < len(l.items) && l.onActivate != nil {
			l.onActivate(selected, l.items[selected])
			return true
		}
		return false
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

func (l *ListView) CanFocus() bool {
	return l.Visible() && l.Enabled()
}

func (l *ListView) SetFocused(focused bool) {
	l.focused = focused
}

func (l *ListView) Focused() bool {
	return l.focused
}

func (l *ListView) Tick(ctx EventContext, now time.Time) bool {
	if !l.dragSelecting || len(l.items) == 0 {
		return false
	}
	if !l.lastDragTick.IsZero() && now.Sub(l.lastDragTick) < 50*time.Millisecond {
		return false
	}
	l.lastDragTick = now
	pointer := l.pressedPoint
	if pointer.Y < l.itemsRect(LocalRect(l)).Y {
		if l.topIndex == 0 {
			return false
		}
		l.topIndex = clampInt(l.topIndex-1, 0, l.maxTopIndex())
	} else if pointer.Y >= l.itemsRect(LocalRect(l)).Bottom() {
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

func (l *ListView) paintRow(ctx PaintContext, rowRect geom.Rect, item ListViewItem, selected bool) {
	textColor := ctx.Theme.Colors.WindowText
	if selected {
		textColor = ctx.Theme.Colors.HighlightText
	}
	x := rowRect.X
	for i, column := range l.columns {
		cellRect := geom.Rect{X: x, Y: rowRect.Y, W: maxInt(column.Width, 1), H: rowRect.H}
		textRect := cellRect.Inset(4)
		if ctx.Text != nil && i < len(item.Texts) && textRect.W > 0 && textRect.H > 0 {
			text := item.Texts[i]
			textSize := ctx.Text.MeasureString(text)
			textX := textRect.X
			switch column.Align {
			case HeaderAlignCenter:
				textX = textRect.X + maxInt((textRect.W-textSize.W)/2, 0)
			case HeaderAlignRight:
				textX = textRect.Right() - textSize.W
			}
			textY := textRect.Y + maxInt((textRect.H-textSize.H)/2, 0)
			ctx.Canvas.PushClip(textRect)
			_ = ctx.Text.DrawString(ctx.Canvas, geom.Point{X: textX, Y: textY}, text, textColor)
			ctx.Canvas.PopClip()
		}
		if i > 0 {
			ctx.Canvas.DrawVLine(cellRect.X, cellRect.Y, cellRect.H, blendColor(ctx.Theme.Colors.Face, ctx.Theme.Colors.Light))
		}
		x += maxInt(column.Width, 1)
	}
}

func (l *ListView) syncChrome() {
	rect := LocalRect(l)
	l.header.SetBounds(geom.Rect{X: 1, Y: 1, W: maxInt(rect.W-l.scrollbarSize-2, 0), H: maxInt(l.headerHeight, 1)})
	l.scrollbar.SetBounds(geom.Rect{
		X: maxInt(rect.W-l.scrollbarSize-1, 1),
		Y: maxInt(l.headerHeight+1, 1),
		W: maxInt(l.scrollbarSize, 1),
		H: maxInt(rect.H-l.headerHeight-2, 0),
	})
	l.scrollbar.SetRange(l.maxTopIndex(), maxInt(l.visibleRows(), 1))
	l.scrollbar.SetValue(l.topIndex)
}

func (l *ListView) itemsRect(rect geom.Rect) geom.Rect {
	return geom.Rect{
		X: rect.X + 2,
		Y: rect.Y + l.headerHeight + 2,
		W: maxInt(rect.W-l.scrollbarSize-4, 0),
		H: maxInt(rect.H-l.headerHeight-4, 0),
	}
}

func (l *ListView) visibleRows() int {
	rows := l.itemsRect(LocalRect(l)).H / maxInt(l.rowHeight, 1)
	return maxInt(rows, 1)
}

func (l *ListView) maxTopIndex() int {
	return maxInt(len(l.items)-l.visibleRows(), 0)
}

func (l *ListView) ensureVisible(index int) {
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

func (l *ListView) setSelectedIndex(index int, notify bool) {
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

func (l *ListView) isSelected(index int) bool {
	return l.selection.Contains(index)
}

func (l *ListView) toggleSelection(index int, notify bool) bool {
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

func (l *ListView) selectRangeTo(index int, notify bool) bool {
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

func (l *ListView) selectAll(notify bool) bool {
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

func (l *ListView) firstSelected() int {
	if index, ok := l.selection.firstSelectedInOrder(l.selectionOrder()); ok {
		return index
	}
	if len(l.items) > 0 {
		return 0
	}
	return -1
}

func (l *ListView) dropInvalidSelection() {
	l.selection.DropInvalid(func(index int) bool {
		return index >= 0 && index < len(l.items)
	})
}

func (l *ListView) selectionOrder() selectionOrder[int] {
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

func (l *ListView) visibleSelectionOrder() selectionOrder[int] {
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

func (l *ListView) anchorIndex() int {
	if anchor, ok := l.selection.Anchor(); ok {
		return anchor
	}
	return -1
}

func (l *ListView) recentIndex() int {
	if index, ok := l.selection.Recent(); ok && index >= 0 && index < len(l.items) {
		return index
	}
	return -1
}

func (l *ListView) recoveryIndexForKey(key event.Key) (int, bool) {
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

func (l *ListView) updateDragSelection(local geom.Point, clampToItems bool) bool {
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

func (l *ListView) selectDragMarquee() bool {
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

func (l *ListView) rowRect(index int) (geom.Rect, bool) {
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

func (l *ListView) indexAtPoint(local geom.Point, clampToItems bool) (int, bool) {
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

func (l *ListView) itemIndexAtPoint(local geom.Point) (int, bool) {
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

func (l *ListView) selectDragRange(index int) bool {
	anchor := l.pressedIndex
	if anchor < 0 || anchor >= len(l.items) {
		anchor = index
	}
	return l.selection.SelectDragRange(l.selectionOrder(), anchor, index)
}

func (l *ListView) selectDragUnion(index int) bool {
	startAnchor := l.dragAnchor
	if startAnchor < 0 || startAnchor >= len(l.items) {
		startAnchor = l.pressedIndex
	}
	if startAnchor < 0 || startAnchor >= len(l.items) {
		startAnchor = index
	}
	return l.selection.SelectDragUnion(l.selectionOrder(), startAnchor, index)
}

func (l *ListView) captureDragBaseSelection() {
	l.selection.CaptureDragBase()
}

func (l *ListView) marqueeRect() (geom.Rect, bool) {
	if !l.selectionBehavior().allowsBlankDrag() || !l.dragSelecting || (l.pressedIndex < 0 && !l.pressedBlank) {
		return geom.Rect{}, false
	}
	return dragMarqueeRect(l.itemsRect(LocalRect(l)), l.pressedStart, l.pressedPoint)
}

func (l *ListView) clearPressed() {
	l.pressedIndex = -1
	l.pressedBlank = false
	l.dragSelecting = false
}

func (l *ListView) selectionBehavior() selectionBehavior {
	return l.selectionOpts.behavior()
}

func (l *ListView) isDoubleClick(index int) bool {
	now := time.Now()
	if l.now != nil {
		now = l.now()
	}
	doubleClick := index >= 0 &&
		index == l.lastClickIndex &&
		!l.lastClickAt.IsZero() &&
		now.Sub(l.lastClickAt) <= l.doubleClick
	l.lastClickIndex = index
	l.lastClickAt = now
	if doubleClick {
		l.lastClickIndex = -1
		l.lastClickAt = time.Time{}
	}
	return doubleClick
}

func (l *ListView) autoFitColumn(index int, measure func(string) geom.Size) {
	if index < 0 || index >= len(l.columns) {
		return
	}
	width := 24
	if measure != nil {
		width = maxInt(width, measure(l.columns[index].Title).W+12)
		for _, item := range l.items {
			if index >= len(item.Texts) {
				continue
			}
			width = maxInt(width, measure(item.Texts[index]).W+12)
		}
	}
	l.columns[index].Width = width
}

func (l *ListView) handleContextMenu(ctx EventContext, local geom.Point) {
	l.syncChrome()
	if l.header.Bounds().Contains(local) || l.scrollbar.Bounds().Contains(local) {
		return
	}
	itemsRect := l.itemsRect(LocalRect(l))
	if !itemsRect.Contains(local) {
		if l.SelectedIndex() != -1 {
			l.setSelectedIndex(-1, false)
			ctx.Invalidate(l)
		}
		if menu := l.contextMenuFor(ListViewContextMenuInfo{Index: -1}); menu != nil {
			ctx.ShowContextMenu(l, geom.Rect{X: local.X, Y: local.Y, W: 1, H: 1}, menu)
		}
		return
	}
	if index, ok := l.itemIndexAtPoint(local); ok {
		if rowRect, rowOK := l.rowRect(index); rowOK && rowRect.Contains(local) {
			if !l.isSelected(index) {
				l.setSelectedIndex(index, true)
				ctx.Invalidate(l)
			}
			if menu := l.contextMenuFor(ListViewContextMenuInfo{
				Index:   index,
				Item:    l.items[index],
				HasItem: true,
			}); menu != nil {
				ctx.ShowContextMenu(l, geom.Rect{X: local.X, Y: local.Y, W: 1, H: 1}, menu)
			}
			return
		}
	}
	if l.SelectedIndex() != -1 {
		l.setSelectedIndex(-1, false)
		ctx.Invalidate(l)
	}
	if menu := l.contextMenuFor(ListViewContextMenuInfo{Index: -1}); menu != nil {
		ctx.ShowContextMenu(l, geom.Rect{X: local.X, Y: local.Y, W: 1, H: 1}, menu)
	}
}

func (l *ListView) contextMenuFor(info ListViewContextMenuInfo) *Menu {
	if l.contextMenuProvider != nil {
		return l.contextMenuProvider(info)
	}
	return l.contextMenu
}
