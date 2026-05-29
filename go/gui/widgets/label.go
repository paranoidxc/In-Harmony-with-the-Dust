package widgets

import (
	"classicui/event"
	"classicui/geom"
	"classicui/uicolor"
	"classicui/widget"
)

type Label struct {
	widget.BaseWidget
	text  string
	color *uicolor.RGBA
}

func NewLabel(id, text string, bounds geom.Rect) *Label {
	return &Label{
		BaseWidget: widget.NewBase(id, bounds),
		text:       text,
	}
}

func (l *Label) SetText(text string) {
	l.text = text
}

func (l *Label) Text() string {
	return l.text
}

func (l *Label) SetColor(color uicolor.RGBA) {
	l.color = &color
}

func (l *Label) Paint(ctx PaintContext) error {
	if !l.Visible() || ctx.Text == nil || l.text == "" {
		return nil
	}
	rect := ctx.BoundsFor(l)
	color := ctx.Theme.Colors.WindowText
	if l.color != nil {
		color = *l.color
	}
	size := ctx.Text.MeasureString(l.text)
	textY := rect.Y
	if rect.H > size.H {
		textY += (rect.H - size.H) / 2
	}
	return ctx.Text.DrawString(ctx.Canvas, geom.Point{X: rect.X, Y: textY}, l.text, color)
}

func (l *Label) MouseEnter(EventContext)            {}
func (l *Label) MouseLeave(EventContext)            {}
func (l *Label) MouseMove(EventContext, geom.Point) {}
func (l *Label) MouseDown(EventContext, event.MouseButtonEvent, geom.Point) {
}
func (l *Label) MouseUp(EventContext, event.MouseButtonEvent, geom.Point) {}
func (l *Label) KeyDown(EventContext, event.KeyEvent) bool                { return false }
func (l *Label) CanFocus() bool                                           { return false }
func (l *Label) SetFocused(bool)                                          {}
func (l *Label) Focused() bool                                            { return false }
