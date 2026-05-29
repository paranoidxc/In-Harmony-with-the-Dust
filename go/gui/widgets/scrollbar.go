package widgets

import (
	"classicui/event"
	"classicui/geom"
	"classicui/uicolor"
	"classicui/widget"
)

type scrollBarPart int

const (
	scrollBarPartNone scrollBarPart = iota
	scrollBarPartDecrease
	scrollBarPartIncrease
	scrollBarPartTrack
	scrollBarPartThumb
)

type ScrollBar struct {
	widget.BaseWidget
	value      int
	maxValue   int
	pageSize   int
	dragging   bool
	dragOffset int
	hotPart    scrollBarPart
	pressed    scrollBarPart
	onChange   func(int)
}

func NewScrollBar(id string, bounds geom.Rect) *ScrollBar {
	return &ScrollBar{
		BaseWidget: widget.NewBase(id, bounds),
		pageSize:   1,
	}
}

func (s *ScrollBar) SetRange(maxValue, pageSize int) {
	if maxValue < 0 {
		maxValue = 0
	}
	if pageSize <= 0 {
		pageSize = 1
	}
	s.maxValue = maxValue
	s.pageSize = pageSize
	s.setValue(clampInt(s.value, 0, s.maxValue))
}

func (s *ScrollBar) SetValue(value int) {
	s.setValue(clampInt(value, 0, s.maxValue))
}

func (s *ScrollBar) Value() int {
	return s.value
}

func (s *ScrollBar) OnChange(fn func(int)) {
	s.onChange = fn
}

func (s *ScrollBar) Paint(ctx PaintContext) error {
	if !s.Visible() {
		return nil
	}

	rect := ctx.BoundsFor(s)
	ctx.Canvas.FillRect(rect, ctx.Theme.Colors.Face)
	ctx.Canvas.FrameRect(rect, ctx.Theme.Colors.Shadow)

	decRect := s.decreaseRect(rect)
	incRect := s.increaseRect(rect)
	trackRect := s.trackRect(rect)
	thumbRect := s.thumbRect(trackRect)

	s.paintButton(ctx, decRect, s.pressed == scrollBarPartDecrease)
	s.paintButton(ctx, incRect, s.pressed == scrollBarPartIncrease)

	ctx.Canvas.FillRect(trackRect, blendColor(ctx.Theme.Colors.Face, ctx.Theme.Colors.Light))
	if thumbRect.H > 0 {
		ctx.Canvas.FillRect(thumbRect, ctx.Theme.Colors.Face)
		if s.pressed == scrollBarPartThumb {
			ctx.Canvas.DrawDoubleBevel(thumbRect, ctx.Theme.Colors.Shadow, ctx.Theme.Colors.Lightest, ctx.Theme.Colors.DarkShadow, ctx.Theme.Colors.Light)
		} else {
			ctx.Canvas.DrawDoubleBevel(thumbRect, ctx.Theme.Colors.Lightest, ctx.Theme.Colors.DarkShadow, ctx.Theme.Colors.Light, ctx.Theme.Colors.Shadow)
		}
		s.paintGrip(ctx, thumbRect)
	}

	s.paintArrow(ctx, decRect, true, s.pressed == scrollBarPartDecrease)
	s.paintArrow(ctx, incRect, false, s.pressed == scrollBarPartIncrease)
	return nil
}

func (s *ScrollBar) MouseEnter(ctx EventContext) {
	if s.hotPart != scrollBarPartNone {
		return
	}
	s.hotPart = scrollBarPartTrack
	ctx.Invalidate(s)
}

func (s *ScrollBar) MouseLeave(ctx EventContext) {
	if s.dragging {
		return
	}
	if s.hotPart == scrollBarPartNone {
		return
	}
	s.hotPart = scrollBarPartNone
	ctx.Invalidate(s)
}

func (s *ScrollBar) MouseMove(ctx EventContext, local geom.Point) {
	if s.dragging {
		next := s.valueFromThumbPointer(local.Y)
		if next != s.value {
			s.setValue(next)
			ctx.Invalidate(s)
		}
		return
	}

	part := s.hitTestLocal(local)
	if part == s.hotPart {
		return
	}
	s.hotPart = part
	ctx.Invalidate(s)
}

func (s *ScrollBar) MouseDown(ctx EventContext, e event.MouseButtonEvent, local geom.Point) {
	if !s.Enabled() || e.Button != event.MouseButtonLeft || !LocalContains(s, local) {
		return
	}
	part := s.hitTestLocal(local)
	s.hotPart = part
	s.pressed = part
	switch part {
	case scrollBarPartDecrease:
		s.setValue(s.value - 1)
	case scrollBarPartIncrease:
		s.setValue(s.value + 1)
	case scrollBarPartTrack:
		thumb := s.thumbRect(s.trackRect(LocalRect(s)))
		if local.Y < thumb.Y {
			s.setValue(s.value - s.pageSize)
		} else {
			s.setValue(s.value + s.pageSize)
		}
	case scrollBarPartThumb:
		s.dragging = true
		s.dragOffset = local.Y - s.thumbRect(s.trackRect(LocalRect(s))).Y
		ctx.Capture(s)
	}
	ctx.Invalidate(s)
}

func (s *ScrollBar) MouseUp(ctx EventContext, e event.MouseButtonEvent, local geom.Point) {
	if e.Button != event.MouseButtonLeft {
		return
	}
	if s.dragging {
		s.dragging = false
		ctx.ReleaseCapture(s)
	}
	s.pressed = scrollBarPartNone
	s.hotPart = s.hitTestLocal(local)
	ctx.Invalidate(s)
}

func (s *ScrollBar) KeyDown(EventContext, event.KeyEvent) bool { return false }
func (s *ScrollBar) CanFocus() bool                            { return false }
func (s *ScrollBar) SetFocused(bool)                           {}
func (s *ScrollBar) Focused() bool                             { return false }

func (s *ScrollBar) setValue(value int) {
	if value == s.value {
		return
	}
	s.value = value
	if s.onChange != nil {
		s.onChange(s.value)
	}
}

func (s *ScrollBar) paintButton(ctx PaintContext, rect geom.Rect, pressed bool) {
	ctx.Canvas.FillRect(rect, ctx.Theme.Colors.Face)
	if pressed {
		ctx.Canvas.DrawDoubleBevel(rect, ctx.Theme.Colors.Shadow, ctx.Theme.Colors.Lightest, ctx.Theme.Colors.DarkShadow, ctx.Theme.Colors.Light)
	} else {
		ctx.Canvas.DrawDoubleBevel(rect, ctx.Theme.Colors.Lightest, ctx.Theme.Colors.DarkShadow, ctx.Theme.Colors.Light, ctx.Theme.Colors.Shadow)
	}
}

func (s *ScrollBar) paintArrow(ctx PaintContext, rect geom.Rect, up, pressed bool) {
	cx := rect.X + rect.W/2
	cy := rect.Y + rect.H/2
	offset := 0
	if pressed {
		offset = 1
	}
	color := ctx.Theme.Colors.DarkShadow
	if up {
		for row := 0; row < 4; row++ {
			for x := -row; x <= row; x++ {
				ctx.Canvas.DrawPixel(cx+x+offset, cy-2+row+offset, color)
			}
		}
		return
	}
	for row := 0; row < 4; row++ {
		for x := -row; x <= row; x++ {
			ctx.Canvas.DrawPixel(cx+x+offset, cy+2-row+offset, color)
		}
	}
}

func (s *ScrollBar) paintGrip(ctx PaintContext, rect geom.Rect) {
	if rect.H < 12 {
		return
	}
	midY := rect.Y + rect.H/2 - 1
	for i := -1; i <= 1; i++ {
		ctx.Canvas.DrawHLine(rect.X+3, midY+i*2, rect.W-6, ctx.Theme.Colors.Shadow)
		ctx.Canvas.DrawHLine(rect.X+3, midY+i*2+1, rect.W-6, ctx.Theme.Colors.Lightest)
	}
}

func (s *ScrollBar) hitTestLocal(local geom.Point) scrollBarPart {
	if !LocalContains(s, local) {
		return scrollBarPartNone
	}
	rect := LocalRect(s)
	if s.decreaseRect(rect).Contains(local) {
		return scrollBarPartDecrease
	}
	if s.increaseRect(rect).Contains(local) {
		return scrollBarPartIncrease
	}
	if s.thumbRect(s.trackRect(rect)).Contains(local) {
		return scrollBarPartThumb
	}
	return scrollBarPartTrack
}

func (s *ScrollBar) decreaseRect(rect geom.Rect) geom.Rect {
	size := rect.W
	return geom.Rect{X: rect.X, Y: rect.Y, W: rect.W, H: minInt(size, rect.H)}
}

func (s *ScrollBar) increaseRect(rect geom.Rect) geom.Rect {
	size := rect.W
	height := minInt(size, rect.H)
	return geom.Rect{X: rect.X, Y: rect.Bottom() - height, W: rect.W, H: height}
}

func (s *ScrollBar) trackRect(rect geom.Rect) geom.Rect {
	dec := s.decreaseRect(rect)
	inc := s.increaseRect(rect)
	return geom.Rect{
		X: rect.X,
		Y: dec.Bottom(),
		W: rect.W,
		H: inc.Y - dec.Bottom(),
	}
}

func (s *ScrollBar) thumbRect(track geom.Rect) geom.Rect {
	if track.H <= 0 {
		return geom.Rect{}
	}
	total := s.maxValue + s.pageSize
	if total <= 0 || s.maxValue == 0 {
		return track
	}
	thumbH := track.H * s.pageSize / total
	thumbH = clampInt(thumbH, 8, track.H)
	free := track.H - thumbH
	pos := 0
	if free > 0 && s.maxValue > 0 {
		pos = free * s.value / s.maxValue
	}
	return geom.Rect{
		X: track.X,
		Y: track.Y + pos,
		W: track.W,
		H: thumbH,
	}
}

func (s *ScrollBar) valueFromThumbPointer(pointerY int) int {
	rect := LocalRect(s)
	track := s.trackRect(rect)
	thumb := s.thumbRect(track)
	free := track.H - thumb.H
	if free <= 0 || s.maxValue <= 0 {
		return 0
	}
	top := pointerY - s.dragOffset
	top = clampInt(top, track.Y, track.Bottom()-thumb.H)
	return (top - track.Y) * s.maxValue / free
}

func blendColor(a, b uicolor.RGBA) uicolor.RGBA {
	return uicolor.RGBA{
		R: uint8((uint16(a.R) + uint16(b.R)) / 2),
		G: uint8((uint16(a.G) + uint16(b.G)) / 2),
		B: uint8((uint16(a.B) + uint16(b.B)) / 2),
		A: 0xFF,
	}
}

func clampInt(value, minValue, maxValue int) int {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
