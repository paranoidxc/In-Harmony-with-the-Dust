package widgets

import (
	"classicui/event"
	"classicui/geom"
	"classicui/uicolor"
	"classicui/widget"
)

type Panel struct {
	widget.BaseWidget
	background *uicolor.RGBA
}

func NewPanel(id string, bounds geom.Rect) *Panel {
	return &Panel{
		BaseWidget: widget.NewBase(id, bounds),
	}
}

func (p *Panel) Add(control Control) {
	p.AddChild(control)
}

func (p *Panel) AddChild(child widget.Widget) {
	if child == nil {
		return
	}
	child.SetParent(p)
	p.BaseWidget.AppendChild(child)
}

func (p *Panel) SetBackground(color uicolor.RGBA) {
	p.background = &color
}

func (p *Panel) Paint(ctx PaintContext) error {
	if !p.Visible() {
		return nil
	}
	abs := ctx.BoundsFor(p)
	if p.background != nil {
		ctx.Canvas.FillRect(abs, *p.background)
	}
	childCtx := ctx
	childCtx.Origin = geom.Point{X: abs.X, Y: abs.Y}
	for _, child := range ControlsOf(p) {
		if err := child.Paint(childCtx); err != nil {
			return err
		}
	}
	return nil
}

func (p *Panel) MouseEnter(EventContext)            {}
func (p *Panel) MouseLeave(EventContext)            {}
func (p *Panel) MouseMove(EventContext, geom.Point) {}
func (p *Panel) MouseDown(EventContext, event.MouseButtonEvent, geom.Point) {
}
func (p *Panel) MouseUp(EventContext, event.MouseButtonEvent, geom.Point) {}
func (p *Panel) KeyDown(EventContext, event.KeyEvent) bool                { return false }
func (p *Panel) CanFocus() bool                                           { return false }
func (p *Panel) SetFocused(bool)                                          {}
func (p *Panel) Focused() bool                                            { return false }
