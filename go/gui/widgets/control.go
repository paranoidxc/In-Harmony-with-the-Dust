package widgets

import (
	"classicui/event"
	"classicui/geom"
	"classicui/paint"
	uitext "classicui/text"
	"classicui/theme"
	"classicui/widget"
)

type EventContext interface {
	Invalidate(Control)
	SetFocus(Control)
	Capture(Control)
	ReleaseCapture(Control)
}

type Control interface {
	widget.Widget
	Paint(PaintContext) error
	MouseEnter(EventContext)
	MouseLeave(EventContext)
	MouseMove(EventContext, geom.Point)
	MouseDown(EventContext, event.MouseButtonEvent, geom.Point)
	MouseUp(EventContext, event.MouseButtonEvent, geom.Point)
	KeyDown(EventContext, event.KeyEvent) bool
	CanFocus() bool
	SetFocused(bool)
	Focused() bool
}

type PaintContext struct {
	Canvas *paint.Canvas
	Theme  *theme.Theme
	Text   *uitext.Renderer
	Origin geom.Point
}

func (c PaintContext) BoundsFor(w widget.Widget) geom.Rect {
	bounds := w.Bounds()
	return bounds.Move(c.Origin.X, c.Origin.Y)
}

func (c PaintContext) Child(w widget.Widget) PaintContext {
	abs := c.BoundsFor(w)
	return PaintContext{
		Canvas: c.Canvas,
		Theme:  c.Theme,
		Text:   c.Text,
		Origin: geom.Point{X: abs.X, Y: abs.Y},
	}
}

func HitTest(root Control, point geom.Point) Control {
	if root == nil || !root.Visible() {
		return nil
	}
	return hitTest(root, point)
}

func FocusableControls(root Control) []Control {
	if root == nil {
		return nil
	}
	var out []Control
	walkFocusable(root, &out)
	return out
}

func ControlsOf(parent widget.Widget) []Control {
	raw := parent.Children()
	out := make([]Control, 0, len(raw))
	for _, child := range raw {
		control, ok := child.(Control)
		if ok {
			out = append(out, control)
		}
	}
	return out
}

func hitTest(control Control, point geom.Point) Control {
	if !control.Visible() || !control.Bounds().Contains(point) {
		return nil
	}

	local := geom.Point{
		X: point.X - control.Bounds().X,
		Y: point.Y - control.Bounds().Y,
	}

	children := ControlsOf(control)
	for i := len(children) - 1; i >= 0; i-- {
		if hit := hitTest(children[i], local); hit != nil {
			return hit
		}
	}
	return control
}

func walkFocusable(control Control, out *[]Control) {
	if control.CanFocus() {
		*out = append(*out, control)
	}
	for _, child := range ControlsOf(control) {
		walkFocusable(child, out)
	}
}
