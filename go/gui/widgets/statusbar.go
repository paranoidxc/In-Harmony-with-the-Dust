package widgets

import (
	"classicui/event"
	"classicui/geom"
	"classicui/widget"
)

type StatusPane struct {
	Text  string
	Width int
}

type StatusBar struct {
	widget.BaseWidget
	panes []StatusPane
}

func NewStatusBar(id string, bounds geom.Rect) *StatusBar {
	return &StatusBar{
		BaseWidget: widget.NewBase(id, bounds),
		panes:      []StatusPane{{}},
	}
}

func (s *StatusBar) SetText(text string) {
	s.panes = []StatusPane{{Text: text}}
}

func (s *StatusBar) SetPanes(panes []StatusPane) {
	if len(panes) == 0 {
		s.panes = []StatusPane{{}}
		return
	}
	s.panes = append([]StatusPane(nil), panes...)
}

func (s *StatusBar) SetPaneText(index int, text string) {
	if index < 0 || index >= len(s.panes) {
		return
	}
	s.panes[index].Text = text
}

func (s *StatusBar) Panes() []StatusPane {
	return append([]StatusPane(nil), s.panes...)
}

func (s *StatusBar) Paint(ctx PaintContext) error {
	if !s.Visible() {
		return nil
	}

	rect := ctx.BoundsFor(s)
	ctx.Canvas.FillRect(rect, ctx.Theme.Colors.Face)
	ctx.Canvas.DrawHLine(rect.X, rect.Y, rect.W, ctx.Theme.Colors.Shadow)
	ctx.Canvas.DrawHLine(rect.X, rect.Y+1, rect.W, ctx.Theme.Colors.Lightest)

	paneRects := s.layoutPanes(rect)
	for i, pane := range paneRects {
		ctx.Canvas.FillRect(pane, ctx.Theme.Colors.Face)
		ctx.Canvas.DrawDoubleBevel(pane, ctx.Theme.Colors.Shadow, ctx.Theme.Colors.Lightest, ctx.Theme.Colors.DarkShadow, ctx.Theme.Colors.Light)

		if ctx.Text == nil {
			continue
		}

		text := s.panes[i].Text
		if text == "" {
			continue
		}

		textSize := ctx.Text.MeasureString(text)
		textY := pane.Y + maxInt((pane.H-textSize.H)/2, 0)
		textRect := pane.Inset(3)
		ctx.Canvas.PushClip(textRect)
		err := ctx.Text.DrawString(ctx.Canvas, geom.Point{X: textRect.X, Y: textY}, text, ctx.Theme.Colors.WindowText)
		ctx.Canvas.PopClip()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *StatusBar) layoutPanes(rect geom.Rect) []geom.Rect {
	if len(s.panes) == 0 {
		return nil
	}

	const gap = 2
	paneHeight := maxInt(rect.H-4, 0)
	available := maxInt(rect.W-4, 0)

	flexible := 0
	fixedWidth := 0
	for _, pane := range s.panes {
		if pane.Width > 0 {
			fixedWidth += pane.Width
			continue
		}
		flexible++
	}

	totalGap := maxInt((len(s.panes)-1)*gap, 0)
	remaining := maxInt(available-fixedWidth-totalGap, 0)
	flexWidth := 0
	extra := 0
	if flexible > 0 {
		flexWidth = remaining / flexible
		extra = remaining % flexible
	}

	rects := make([]geom.Rect, 0, len(s.panes))
	x := rect.X + 2
	for i, pane := range s.panes {
		width := pane.Width
		if width <= 0 {
			width = flexWidth
			if extra > 0 {
				width++
				extra--
			}
		}
		if i == len(s.panes)-1 {
			width = maxInt(rect.Right()-2-x, 0)
		}
		rects = append(rects, geom.Rect{
			X: x,
			Y: rect.Y + 2,
			W: maxInt(width, 0),
			H: paneHeight,
		})
		x += width + gap
	}
	return rects
}

func (s *StatusBar) MouseEnter(EventContext)            {}
func (s *StatusBar) MouseLeave(EventContext)            {}
func (s *StatusBar) MouseMove(EventContext, geom.Point) {}
func (s *StatusBar) MouseDown(EventContext, event.MouseButtonEvent, geom.Point) {
}
func (s *StatusBar) MouseUp(EventContext, event.MouseButtonEvent, geom.Point) {}
func (s *StatusBar) KeyDown(EventContext, event.KeyEvent) bool                { return false }
func (s *StatusBar) CanFocus() bool                                           { return false }
func (s *StatusBar) SetFocused(bool)                                          {}
func (s *StatusBar) Focused() bool                                            { return false }
