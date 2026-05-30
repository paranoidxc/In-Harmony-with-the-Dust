package widgets

import (
	"classicui/event"
	"classicui/geom"
	"classicui/widget"
	"time"
)

type HeaderColumn struct {
	Title string
	Width int
	Align HeaderAlign
}

type HeaderAlign int

const (
	HeaderAlignLeft HeaderAlign = iota
	HeaderAlignCenter
	HeaderAlignRight
)

type HeaderControl struct {
	widget.BaseWidget
	columns         []HeaderColumn
	focused         bool
	hotIndex        int
	hotDivider      int
	pressedIndex    int
	resizingIndex   int
	resizeStartX    int
	resizeWidth     int
	sortedIndex     int
	sortedDesc      bool
	onColumnClick   func(int)
	onColumnResize  func(int, int)
	onColumnAutoFit func(EventContext, int)
	now             func() time.Time
	doubleClick     time.Duration
	lastDivider     int
	lastDividerAt   time.Time
}

func NewHeaderControl(id string, bounds geom.Rect, columns ...HeaderColumn) *HeaderControl {
	header := &HeaderControl{
		BaseWidget:    widget.NewBase(id, bounds),
		hotIndex:      -1,
		hotDivider:    -1,
		pressedIndex:  -1,
		resizingIndex: -1,
		sortedIndex:   -1,
		now:           time.Now,
		doubleClick:   500 * time.Millisecond,
		lastDivider:   -1,
	}
	header.SetColumns(columns...)
	return header
}

func (h *HeaderControl) SetColumns(columns ...HeaderColumn) {
	h.columns = append([]HeaderColumn(nil), columns...)
}

func (h *HeaderControl) Columns() []HeaderColumn {
	return append([]HeaderColumn(nil), h.columns...)
}

func (h *HeaderControl) SetSortIndicator(index int, descending bool) {
	if index < 0 || index >= len(h.columns) {
		h.sortedIndex = -1
		h.sortedDesc = false
		return
	}
	h.sortedIndex = index
	h.sortedDesc = descending
}

func (h *HeaderControl) OnColumnClick(fn func(int)) {
	h.onColumnClick = fn
}

func (h *HeaderControl) OnColumnResize(fn func(int, int)) {
	h.onColumnResize = fn
}

func (h *HeaderControl) OnColumnAutoFit(fn func(EventContext, int)) {
	h.onColumnAutoFit = fn
}

func (h *HeaderControl) Paint(ctx PaintContext) error {
	if !h.Visible() {
		return nil
	}
	rect := ctx.BoundsFor(h)
	ctx.Canvas.FillRect(rect, ctx.Theme.Colors.Face)
	ctx.Canvas.DrawDoubleBevel(rect, ctx.Theme.Colors.Lightest, ctx.Theme.Colors.DarkShadow, ctx.Theme.Colors.Light, ctx.Theme.Colors.Shadow)

	for i := range h.columns {
		colRect, ok := h.columnRect(i)
		if !ok {
			continue
		}
		colRect = colRect.Move(rect.X, rect.Y)
		pressed := i == h.pressedIndex
		fill := ctx.Theme.Colors.Face
		if i == h.hotIndex && h.hotDivider < 0 && !pressed && h.resizingIndex < 0 {
			fill = blendColor(ctx.Theme.Colors.Face, ctx.Theme.Colors.Lightest)
		}
		ctx.Canvas.FillRect(colRect, fill)
		if pressed {
			ctx.Canvas.DrawDoubleBevel(colRect, ctx.Theme.Colors.Shadow, ctx.Theme.Colors.Lightest, ctx.Theme.Colors.DarkShadow, ctx.Theme.Colors.Light)
		} else {
			ctx.Canvas.DrawDoubleBevel(colRect, ctx.Theme.Colors.Lightest, ctx.Theme.Colors.DarkShadow, ctx.Theme.Colors.Light, ctx.Theme.Colors.Shadow)
		}
		if i > 0 {
			ctx.Canvas.DrawVLine(colRect.X, colRect.Y+2, maxInt(colRect.H-4, 0), ctx.Theme.Colors.Shadow)
		}
		if ctx.Text != nil && h.columns[i].Title != "" {
			textSize := ctx.Text.MeasureString(h.columns[i].Title)
			textRect := colRect.Inset(4)
			if textRect.W <= 0 || textRect.H <= 0 {
				continue
			}
			textX := textRect.X
			switch h.columns[i].Align {
			case HeaderAlignCenter:
				textX = textRect.X + maxInt((textRect.W-textSize.W)/2, 0)
			case HeaderAlignRight:
				textX = textRect.Right() - textSize.W
			}
			if h.sortedIndex == i {
				textRect.W -= 10
				if h.columns[i].Align == HeaderAlignRight {
					textX = textRect.Right() - textSize.W
				}
			}
			offset := 0
			if pressed {
				offset = 1
			}
			textY := textRect.Y + maxInt((textRect.H-textSize.H)/2, 0) + offset
			ctx.Canvas.PushClip(textRect)
			if err := ctx.Text.DrawString(ctx.Canvas, geom.Point{X: textX + offset, Y: textY}, h.columns[i].Title, ctx.Theme.Colors.WindowText); err != nil {
				ctx.Canvas.PopClip()
				return err
			}
			ctx.Canvas.PopClip()
		}
		if h.sortedIndex == i {
			h.paintSortArrow(ctx, colRect, pressed)
		}
	}
	if dividerX, ok := h.hotDividerX(); ok {
		ctx.Canvas.DrawVLine(rect.X+dividerX, rect.Y+2, maxInt(rect.H-4, 0), ctx.Theme.Colors.DarkShadow)
	}
	return nil
}

func (h *HeaderControl) MouseEnter(EventContext) {}

func (h *HeaderControl) MouseLeave(ctx EventContext) {
	if h.resizingIndex >= 0 {
		return
	}
	if h.hotIndex < 0 && h.hotDivider < 0 && h.pressedIndex < 0 {
		return
	}
	h.hotIndex = -1
	h.hotDivider = -1
	if h.pressedIndex < 0 {
		ctx.Invalidate(h)
	}
}

func (h *HeaderControl) MouseMove(ctx EventContext, local geom.Point) {
	if h.resizingIndex >= 0 {
		width := maxInt(h.resizeWidth+(local.X-h.resizeStartX), 24)
		if h.columns[h.resizingIndex].Width == width {
			return
		}
		h.columns[h.resizingIndex].Width = width
		if h.onColumnResize != nil {
			h.onColumnResize(h.resizingIndex, width)
		}
		if parent, ok := h.Parent().(Control); ok {
			ctx.Invalidate(parent)
			return
		}
		ctx.Invalidate(h)
		return
	}
	divider := h.dividerAt(local)
	index := h.columnAt(local)
	if divider >= 0 {
		index = -1
	}
	if h.pressedIndex >= 0 {
		if h.hotIndex == index && h.hotDivider == divider {
			return
		}
		h.hotIndex = index
		h.hotDivider = divider
		ctx.Invalidate(h)
		return
	}
	if h.hotIndex == index && h.hotDivider == divider {
		return
	}
	h.hotIndex = index
	h.hotDivider = divider
	ctx.Invalidate(h)
}

func (h *HeaderControl) MouseDown(ctx EventContext, ev event.MouseButtonEvent, local geom.Point) {
	if !h.Enabled() || ev.Button != event.MouseButtonLeft || !LocalContains(h, local) {
		return
	}
	if divider := h.dividerAt(local); divider >= 0 {
		if h.isDividerDoubleClick(divider) {
			if h.onColumnAutoFit != nil {
				h.onColumnAutoFit(ctx, divider)
			}
			if parent, ok := h.Parent().(Control); ok {
				ctx.Invalidate(parent)
				return
			}
			ctx.Invalidate(h)
			return
		}
		h.hotDivider = divider
		h.hotIndex = -1
		h.resizingIndex = divider
		h.resizeStartX = local.X
		h.resizeWidth = h.columns[divider].Width
		ctx.SetFocus(h)
		ctx.Capture(h)
		ctx.Invalidate(h)
		return
	}
	index := h.columnAt(local)
	if index < 0 {
		return
	}
	h.hotDivider = -1
	h.hotIndex = index
	h.pressedIndex = index
	ctx.SetFocus(h)
	ctx.Capture(h)
	ctx.Invalidate(h)
}

func (h *HeaderControl) MouseUp(ctx EventContext, ev event.MouseButtonEvent, local geom.Point) {
	if ev.Button != event.MouseButtonLeft {
		return
	}
	if h.resizingIndex >= 0 {
		h.resizingIndex = -1
		h.hotDivider = h.dividerAt(local)
		ctx.ReleaseCapture(h)
		if parent, ok := h.Parent().(Control); ok {
			ctx.Invalidate(parent)
			return
		}
		ctx.Invalidate(h)
		return
	}
	if h.pressedIndex < 0 {
		return
	}
	index := h.columnAt(local)
	pressed := h.pressedIndex
	h.pressedIndex = -1
	h.hotIndex = index
	h.hotDivider = -1
	ctx.ReleaseCapture(h)
	ctx.Invalidate(h)
	if index == pressed && h.onColumnClick != nil {
		h.onColumnClick(index)
	}
}

func (h *HeaderControl) KeyDown(EventContext, event.KeyEvent) bool { return false }

func (h *HeaderControl) CanFocus() bool {
	return h.Visible() && h.Enabled()
}

func (h *HeaderControl) SetFocused(focused bool) {
	h.focused = focused
}

func (h *HeaderControl) Focused() bool {
	return h.focused
}

func (h *HeaderControl) ColumnWidth(index int) int {
	if index < 0 || index >= len(h.columns) {
		return 0
	}
	return h.columns[index].Width
}

func (h *HeaderControl) columnRect(index int) (geom.Rect, bool) {
	if index < 0 || index >= len(h.columns) {
		return geom.Rect{}, false
	}
	x := 0
	height := maxInt(h.Bounds().H, 1)
	for i := 0; i < len(h.columns); i++ {
		width := maxInt(h.columns[i].Width, 1)
		if i == index {
			return geom.Rect{X: x, Y: 0, W: width, H: height}, true
		}
		x += width
	}
	return geom.Rect{}, false
}

func (h *HeaderControl) columnAt(local geom.Point) int {
	if !LocalContains(h, local) {
		return -1
	}
	x := 0
	for i := range h.columns {
		width := maxInt(h.columns[i].Width, 1)
		if local.X >= x && local.X < x+width {
			return i
		}
		x += width
	}
	return -1
}

func (h *HeaderControl) dividerAt(local geom.Point) int {
	if !LocalContains(h, local) {
		return -1
	}
	x := 0
	for i := 0; i < len(h.columns)-1; i++ {
		x += maxInt(h.columns[i].Width, 1)
		if local.X >= x-3 && local.X <= x+2 {
			return i
		}
	}
	return -1
}

func (h *HeaderControl) hotDividerX() (int, bool) {
	index := h.hotDivider
	if h.resizingIndex >= 0 {
		index = h.resizingIndex
	}
	if index < 0 || index >= len(h.columns) {
		return 0, false
	}
	x := 0
	for i := 0; i <= index; i++ {
		x += maxInt(h.columns[i].Width, 1)
	}
	return x, true
}

func (h *HeaderControl) paintSortArrow(ctx PaintContext, rect geom.Rect, pressed bool) {
	cx := rect.Right() - 8
	cy := rect.Y + rect.H/2
	if pressed {
		cx++
		cy++
	}
	color := ctx.Theme.Colors.DarkShadow
	if h.sortedDesc {
		for row := 0; row < 4; row++ {
			for x := -row; x <= row; x++ {
				ctx.Canvas.DrawPixel(cx+x, cy-1+row, color)
			}
		}
		return
	}
	for row := 0; row < 4; row++ {
		for x := -row; x <= row; x++ {
			ctx.Canvas.DrawPixel(cx+x, cy+1-row, color)
		}
	}
}

func (h *HeaderControl) isDividerDoubleClick(index int) bool {
	now := time.Now()
	if h.now != nil {
		now = h.now()
	}
	doubleClick := index >= 0 &&
		index == h.lastDivider &&
		!h.lastDividerAt.IsZero() &&
		now.Sub(h.lastDividerAt) <= h.doubleClick
	h.lastDivider = index
	h.lastDividerAt = now
	if doubleClick {
		h.lastDivider = -1
		h.lastDividerAt = time.Time{}
	}
	return doubleClick
}
