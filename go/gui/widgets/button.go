package widgets

import (
	"classicui/event"
	"classicui/geom"
	"classicui/uicolor"
	"classicui/widget"
)

type Button struct {
	widget.BaseWidget
	text      string
	tooltip   string
	focused   bool
	hot       bool
	pressed   bool
	tracking  bool
	isDefault bool
	onClick   func()
}

func NewButton(id, text string, bounds geom.Rect) *Button {
	return &Button{
		BaseWidget: widget.NewBase(id, bounds),
		text:       text,
	}
}

func (b *Button) SetText(text string) {
	b.text = text
}

func (b *Button) Text() string {
	return b.text
}

func (b *Button) SetTooltip(text string) {
	b.tooltip = text
}

func (b *Button) SetDefault(isDefault bool) {
	b.isDefault = isDefault
}

func (b *Button) Default() bool {
	return b.isDefault
}

func (b *Button) OnClick(fn func()) {
	b.onClick = fn
}

func (b *Button) Paint(ctx PaintContext) error {
	if !b.Visible() {
		return nil
	}
	rect := ctx.BoundsFor(b)
	frameRect := rect
	innerRect := rect

	if b.isDefault {
		ctx.Canvas.FrameRect(frameRect, ctx.Theme.Colors.DarkShadow)
		innerRect = rect.Inset(1)
	}

	fillColor := ctx.Theme.Colors.Face
	if b.hot && !b.pressed {
		fillColor = blend(ctx.Theme.Colors.Face, ctx.Theme.Colors.Lightest)
	}
	ctx.Canvas.FillRect(innerRect, fillColor)
	if b.pressed {
		ctx.Canvas.DrawDoubleBevel(innerRect, ctx.Theme.Colors.Shadow, ctx.Theme.Colors.Lightest, ctx.Theme.Colors.DarkShadow, ctx.Theme.Colors.Light)
	} else {
		ctx.Canvas.DrawDoubleBevel(innerRect, ctx.Theme.Colors.Lightest, ctx.Theme.Colors.DarkShadow, ctx.Theme.Colors.Light, ctx.Theme.Colors.Shadow)
	}

	offset := 0
	if b.pressed {
		offset = 1
	}

	if ctx.Text != nil && b.text != "" {
		textSize := ctx.Text.MeasureString(b.text)
		textX := innerRect.X + (innerRect.W-textSize.W)/2 + offset
		textY := innerRect.Y + (innerRect.H-textSize.H)/2 + offset
		if err := ctx.Text.DrawString(ctx.Canvas, geom.Point{X: textX, Y: textY}, b.text, ctx.Theme.Colors.WindowText); err != nil {
			return err
		}
	}

	if b.focused {
		focusRect := innerRect.Inset(ctx.Theme.Metrics.FocusRectInset)
		if focusRect.W > 1 && focusRect.H > 1 {
			ctx.Canvas.DrawFocusRect(focusRect, ctx.Theme.Colors.DarkShadow)
		}
	}
	return nil
}

func (b *Button) MouseEnter(ctx EventContext) {
	if !b.Enabled() {
		return
	}
	if b.hot {
		return
	}
	b.hot = true
	ctx.Invalidate(b)
}

func (b *Button) MouseLeave(ctx EventContext) {
	if !b.hot {
		return
	}
	b.hot = false
	ctx.Invalidate(b)
}

func (b *Button) MouseMove(ctx EventContext, local geom.Point) {
	if !b.tracking {
		return
	}
	nextPressed := LocalContains(b, local)
	if nextPressed == b.pressed {
		return
	}
	b.pressed = nextPressed
	ctx.Invalidate(b)
}

func (b *Button) MouseDown(ctx EventContext, e event.MouseButtonEvent, local geom.Point) {
	if !b.Enabled() || e.Button != event.MouseButtonLeft || !LocalContains(b, local) {
		return
	}
	b.tracking = true
	b.pressed = true
	b.hot = true
	ctx.SetFocus(b)
	ctx.Capture(b)
	ctx.Invalidate(b)
}

func (b *Button) MouseUp(ctx EventContext, e event.MouseButtonEvent, local geom.Point) {
	if e.Button != event.MouseButtonLeft || !b.tracking {
		return
	}
	clicked := LocalContains(b, local) && b.pressed
	b.tracking = false
	b.pressed = false
	b.hot = LocalContains(b, local)
	ctx.ReleaseCapture(b)
	ctx.Invalidate(b)
	if clicked {
		b.fireClick()
	}
}

func (b *Button) KeyDown(ctx EventContext, e event.KeyEvent) bool {
	if !b.Enabled() {
		return false
	}
	if e.Key != event.KeySpace && e.Key != event.KeyEnter {
		return false
	}
	b.fireClick()
	ctx.Invalidate(b)
	return true
}

func (b *Button) CanFocus() bool {
	return b.Visible() && b.Enabled()
}

func (b *Button) SetFocused(focused bool) {
	b.focused = focused
}

func (b *Button) Focused() bool {
	return b.focused
}

func (b *Button) TooltipAt(geom.Point, func(string) geom.Size) TooltipInfo {
	return TooltipInfo{
		Text:   b.tooltip,
		Anchor: LocalRect(b),
	}
}

func (b *Button) fireClick() {
	if b.onClick != nil {
		b.onClick()
	}
}

func blend(a, b uicolor.RGBA) uicolor.RGBA {
	return uicolor.RGBA{
		R: uint8((uint16(a.R) + uint16(b.R)) / 2),
		G: uint8((uint16(a.G) + uint16(b.G)) / 2),
		B: uint8((uint16(a.B) + uint16(b.B)) / 2),
		A: 0xFF,
	}
}
