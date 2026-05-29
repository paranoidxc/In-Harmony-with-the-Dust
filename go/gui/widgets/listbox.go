package widgets

import (
	"classicui/event"
	"classicui/geom"
	"classicui/widget"
)

type ListBox struct {
	widget.BaseWidget
	items         []string
	selected      int
	topIndex      int
	pressedIndex  int
	focused       bool
	rowHeight     int
	scrollbarSize int
	scrollbar     *ScrollBar
	onChange      func(int, string)
	onActivate    func(int, string)
}

func NewListBox(id string, bounds geom.Rect) *ListBox {
	list := &ListBox{
		BaseWidget:    widget.NewBase(id, bounds),
		selected:      -1,
		pressedIndex:  -1,
		rowHeight:     16,
		scrollbarSize: 16,
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
		l.selected = -1
		l.topIndex = 0
		l.scrollbar.SetRange(0, 1)
		return
	}
	if l.selected >= len(l.items) {
		l.selected = len(l.items) - 1
	}
	l.syncScrollBar()
}

func (l *ListBox) Items() []string {
	return append([]string(nil), l.items...)
}

func (l *ListBox) SelectedIndex() int {
	return l.selected
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
		if index == l.selected {
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
	row := (local.Y - itemsRect.Y) / l.rowHeight
	index := l.topIndex + row
	l.pressedIndex = index
	l.setSelectedIndex(index, true)
	ctx.Invalidate(l)
}

func (l *ListBox) MouseUp(ctx EventContext, ev event.MouseButtonEvent, local geom.Point) {
	if l.scrollbar.dragging || l.scrollbar.Bounds().Contains(local) {
		l.scrollbar.MouseUp(ctx, ev, geom.Point{X: local.X - l.scrollbar.Bounds().X, Y: local.Y - l.scrollbar.Bounds().Y})
		l.topIndex = l.scrollbar.Value()
		ctx.Invalidate(l)
		l.pressedIndex = -1
		return
	}
	itemsRect := l.itemsRect(LocalRect(l))
	if itemsRect.Contains(local) {
		row := (local.Y - itemsRect.Y) / l.rowHeight
		index := l.topIndex + row
		if index == l.pressedIndex && index >= 0 && index < len(l.items) && l.onActivate != nil {
			l.onActivate(index, l.items[index])
		}
	}
	l.pressedIndex = -1
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
	}
}

func (l *ListBox) KeyDown(ctx EventContext, ev event.KeyEvent) bool {
	if !l.Enabled() {
		return false
	}
	if len(l.items) == 0 {
		return true
	}
	next := l.selected
	if next < 0 {
		next = 0
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

	l.setSelectedIndex(next, true)
	l.ensureVisible(l.selected)
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
		l.selected = -1
		return
	}
	if index < 0 {
		if l.selected == -1 {
			return
		}
		l.selected = -1
		return
	}
	index = clampInt(index, 0, len(l.items)-1)
	if index == l.selected {
		return
	}
	l.selected = index
	l.ensureVisible(index)
	if notify && l.onChange != nil {
		l.onChange(index, l.items[index])
	}
}
