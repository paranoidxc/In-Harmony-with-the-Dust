package widgets

import (
	"classicui/event"
	"classicui/geom"
	"classicui/widget"
)

type ToolbarItem struct {
	ID        CommandID
	Text      string
	Checked   bool
	Disabled  bool
	Separator bool
	Width     int
}

func NewToolbarButton(id CommandID, text string) *ToolbarItem {
	return &ToolbarItem{
		ID:   id,
		Text: text,
	}
}

func NewToolbarSeparator() *ToolbarItem {
	return &ToolbarItem{Separator: true}
}

type Toolbar struct {
	widget.BaseWidget
	items        []*ToolbarItem
	hotIndex     int
	pressedIndex int
	tracking     bool
	onCommand    func(CommandID)
}

func NewToolbar(id string, bounds geom.Rect, items ...*ToolbarItem) *Toolbar {
	toolbar := &Toolbar{
		BaseWidget:   widget.NewBase(id, bounds),
		hotIndex:     -1,
		pressedIndex: -1,
	}
	toolbar.SetItems(items...)
	return toolbar
}

func (t *Toolbar) SetItems(items ...*ToolbarItem) {
	t.items = append([]*ToolbarItem(nil), items...)
	t.hotIndex = -1
	t.pressedIndex = -1
	t.tracking = false
}

func (t *Toolbar) Items() []*ToolbarItem {
	return append([]*ToolbarItem(nil), t.items...)
}

func (t *Toolbar) SetChecked(id CommandID, checked bool) {
	if item := t.itemByID(id); item != nil {
		item.Checked = checked
	}
}

func (t *Toolbar) SetDisabled(id CommandID, disabled bool) {
	if item := t.itemByID(id); item != nil {
		item.Disabled = disabled
	}
}

func (t *Toolbar) OnCommand(fn func(CommandID)) {
	t.onCommand = fn
}

func (t *Toolbar) Paint(ctx PaintContext) error {
	if !t.Visible() {
		return nil
	}

	rect := ctx.BoundsFor(t)
	ctx.Canvas.FillRect(rect, ctx.Theme.Colors.Face)
	ctx.Canvas.DrawHLine(rect.X, rect.Y, rect.W, ctx.Theme.Colors.Lightest)
	ctx.Canvas.DrawHLine(rect.X, rect.Bottom()-2, rect.W, ctx.Theme.Colors.Shadow)
	ctx.Canvas.DrawHLine(rect.X, rect.Bottom()-1, rect.W, ctx.Theme.Colors.Lightest)

	for i, item := range t.items {
		if item == nil {
			continue
		}
		itemRect := t.itemRect(i, t.measureWithPaint(ctx))
		screenRect := itemRect.Move(rect.X, rect.Y)
		if item.Separator {
			lineX := screenRect.X + screenRect.W/2
			lineY := screenRect.Y + 4
			lineH := maxInt(screenRect.H-8, 0)
			ctx.Canvas.DrawVLine(lineX, lineY, lineH, ctx.Theme.Colors.Shadow)
			ctx.Canvas.DrawVLine(lineX+1, lineY, lineH, ctx.Theme.Colors.Lightest)
			continue
		}

		pressed := t.isPressed(i)
		hot := t.hotIndex == i
		if pressed || hot {
			fill := ctx.Theme.Colors.Face
			if hot && !pressed {
				fill = blend(ctx.Theme.Colors.Face, ctx.Theme.Colors.Lightest)
			}
			ctx.Canvas.FillRect(screenRect, fill)
			if pressed {
				ctx.Canvas.DrawDoubleBevel(screenRect, ctx.Theme.Colors.Shadow, ctx.Theme.Colors.Lightest, ctx.Theme.Colors.DarkShadow, ctx.Theme.Colors.Light)
			} else {
				ctx.Canvas.DrawDoubleBevel(screenRect, ctx.Theme.Colors.Lightest, ctx.Theme.Colors.DarkShadow, ctx.Theme.Colors.Light, ctx.Theme.Colors.Shadow)
			}
		}

		if ctx.Text == nil || item.Text == "" {
			continue
		}

		textSize := ctx.Text.MeasureString(item.Text)
		offset := 0
		if pressed {
			offset = 1
		}
		textX := screenRect.X + maxInt((screenRect.W-textSize.W)/2, 0) + offset
		textY := screenRect.Y + maxInt((screenRect.H-textSize.H)/2, 0) + offset
		textColor := ctx.Theme.Colors.WindowText
		if item.Disabled {
			textColor = ctx.Theme.Colors.GrayText
		}

		ctx.Canvas.PushClip(screenRect.Inset(2))
		err := ctx.Text.DrawString(ctx.Canvas, geom.Point{X: textX, Y: textY}, item.Text, textColor)
		ctx.Canvas.PopClip()
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *Toolbar) MouseEnter(EventContext) {}

func (t *Toolbar) MouseLeave(ctx EventContext) {
	if t.tracking {
		return
	}
	if t.hotIndex < 0 {
		return
	}
	t.hotIndex = -1
	ctx.Invalidate(t)
}

func (t *Toolbar) MouseMove(ctx EventContext, local geom.Point) {
	next := t.hitIndex(local, ctx.MeasureText)
	changed := false
	if next != t.hotIndex {
		t.hotIndex = next
		changed = true
	}

	if t.tracking {
		pressed := -1
		if next >= 0 && t.itemEnabled(next) {
			pressed = next
		}
		if pressed != t.pressedIndex {
			t.pressedIndex = pressed
			changed = true
		}
	}

	if changed {
		ctx.Invalidate(t)
	}
}

func (t *Toolbar) MouseDown(ctx EventContext, e event.MouseButtonEvent, local geom.Point) {
	if e.Button != event.MouseButtonLeft || !LocalContains(t, local) {
		return
	}

	index := t.hitIndex(local, ctx.MeasureText)
	if index < 0 || !t.itemEnabled(index) {
		return
	}

	t.tracking = true
	t.hotIndex = index
	t.pressedIndex = index
	ctx.Capture(t)
	ctx.Invalidate(t)
}

func (t *Toolbar) MouseUp(ctx EventContext, e event.MouseButtonEvent, local geom.Point) {
	if e.Button != event.MouseButtonLeft || !t.tracking {
		return
	}

	clicked := t.pressedIndex >= 0 && t.pressedIndex == t.hitIndex(local, ctx.MeasureText)
	command := CommandID("")
	if clicked {
		command = t.items[t.pressedIndex].ID
	}

	t.tracking = false
	t.pressedIndex = -1
	t.hotIndex = t.hitIndex(local, ctx.MeasureText)
	ctx.ReleaseCapture(t)
	ctx.Invalidate(t)

	if command != "" {
		if t.onCommand != nil {
			t.onCommand(command)
		} else {
			ctx.DispatchCommand(command)
		}
	}
}

func (t *Toolbar) KeyDown(EventContext, event.KeyEvent) bool { return false }
func (t *Toolbar) CanFocus() bool                            { return false }
func (t *Toolbar) SetFocused(bool)                           {}
func (t *Toolbar) Focused() bool                             { return false }

func (t *Toolbar) isPressed(index int) bool {
	if index < 0 || index >= len(t.items) {
		return false
	}
	if t.pressedIndex == index {
		return true
	}
	item := t.items[index]
	return item != nil && item.Checked
}

func (t *Toolbar) itemRect(index int, measure func(string) geom.Size) geom.Rect {
	const gap = 4
	rect := LocalRect(t)
	x := 4
	height := maxInt(rect.H-8, 0)

	for i := 0; i < index && i < len(t.items); i++ {
		x += t.itemWidth(t.items[i], measure) + gap
	}

	return geom.Rect{
		X: x,
		Y: 4,
		W: t.itemWidth(t.items[index], measure),
		H: height,
	}
}

func (t *Toolbar) itemWidth(item *ToolbarItem, measure func(string) geom.Size) int {
	if item == nil {
		return 0
	}
	if item.Separator {
		return 8
	}
	if item.Width > 0 {
		return item.Width
	}
	return maxInt(measure(item.Text).W+16, 24)
}

func (t *Toolbar) hitIndex(local geom.Point, measure func(string) geom.Size) int {
	if !LocalContains(t, local) {
		return -1
	}
	for i, item := range t.items {
		if item == nil || item.Separator {
			continue
		}
		if t.itemRect(i, measure).Contains(local) {
			return i
		}
	}
	return -1
}

func (t *Toolbar) itemEnabled(index int) bool {
	if index < 0 || index >= len(t.items) {
		return false
	}
	item := t.items[index]
	return item != nil && !item.Separator && !item.Disabled
}

func (t *Toolbar) itemByID(id CommandID) *ToolbarItem {
	for _, item := range t.items {
		if item != nil && item.ID == id {
			return item
		}
	}
	return nil
}

func (t *Toolbar) measureWithPaint(ctx PaintContext) func(string) geom.Size {
	if ctx.Text == nil {
		return func(text string) geom.Size {
			return geom.Size{W: len([]rune(text)) * 7, H: 14}
		}
	}
	return ctx.Text.MeasureString
}
