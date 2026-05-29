package desktop

import (
	"classicui/geom"
	"classicui/paint"
	uitext "classicui/text"
	"classicui/theme"
	"classicui/uicolor"
	"classicui/widget"
	"classicui/widgets"
)

type HitPart int

const (
	HitNowhere HitPart = iota
	HitClient
	HitCaption
	HitMenuBar
	HitClose
)

type Window struct {
	widget.BaseWidget
	title           string
	active          bool
	closeHot        bool
	closePressed    bool
	menuBar         *widgets.MenuBar
	menuActiveIndex int
	content         *widgets.Panel
	defaultButton   *widgets.Button
}

func NewWindow(id string, bounds geom.Rect) *Window {
	return &Window{
		BaseWidget:      widget.NewBase(id, bounds),
		active:          true,
		menuActiveIndex: -1,
		content:         widgets.NewPanel(id+".content", geom.Rect{}),
	}
}

func (w *Window) SetTitle(title string) {
	w.title = title
}

func (w *Window) Title() string {
	return w.title
}

func (w *Window) Content() *widgets.Panel {
	return w.content
}

func (w *Window) SetDefaultButton(button *widgets.Button) {
	if w.defaultButton != nil {
		w.defaultButton.SetDefault(false)
	}
	w.defaultButton = button
	if w.defaultButton != nil {
		w.defaultButton.SetDefault(true)
	}
}

func (w *Window) DefaultButton() *widgets.Button {
	return w.defaultButton
}

func (w *Window) SetActive(active bool) {
	w.active = active
}

func (w *Window) Active() bool {
	return w.active
}

func (w *Window) SetCloseHot(hot bool) {
	w.closeHot = hot
}

func (w *Window) SetClosePressed(pressed bool) {
	w.closePressed = pressed
}

func (w *Window) CaptionRect(th *theme.Theme) geom.Rect {
	frame := th.Metrics.BorderWidth + th.Metrics.WindowFrameInner
	return geom.Rect{
		X: w.Bounds().X + frame,
		Y: w.Bounds().Y + frame,
		W: w.Bounds().W - frame*2,
		H: th.Metrics.CaptionHeight,
	}
}

func (w *Window) ClientRect(th *theme.Theme) geom.Rect {
	frame := th.Metrics.BorderWidth + th.Metrics.WindowFrameInner
	caption := w.CaptionRect(th)
	top := caption.Bottom() + 1
	height := w.Bounds().H - frame*2 - caption.H - 1
	if menu := w.MenuBarRect(th); !menu.Empty() {
		top = menu.Bottom() + 1
		height -= menu.H + 1
	}
	return geom.Rect{
		X: w.Bounds().X + frame,
		Y: top,
		W: max(w.Bounds().W-frame*2, 0),
		H: max(height, 0),
	}
}

func (w *Window) CloseButtonRect(th *theme.Theme) geom.Rect {
	caption := w.CaptionRect(th)
	size := th.Metrics.IconSizeSmall - 1
	if size > caption.H-4 {
		size = caption.H - 4
	}
	return geom.Rect{
		X: caption.Right() - size - 3,
		Y: caption.Y + (caption.H-size)/2,
		W: size,
		H: size,
	}
}

func (w *Window) HitTest(p geom.Point, th *theme.Theme) HitPart {
	if !w.Bounds().Contains(p) {
		return HitNowhere
	}
	if w.CloseButtonRect(th).Contains(p) {
		return HitClose
	}
	if w.CaptionRect(th).Contains(p) {
		return HitCaption
	}
	if !w.MenuBarRect(th).Empty() && w.MenuBarRect(th).Contains(p) {
		return HitMenuBar
	}
	return HitClient
}

func (w *Window) ControlAt(screenPoint geom.Point, th *theme.Theme) widgets.Control {
	client := w.ClientRect(th)
	if !client.Contains(screenPoint) {
		return nil
	}
	w.syncContentBounds(th)
	local := geom.Point{
		X: screenPoint.X - client.X,
		Y: screenPoint.Y - client.Y,
	}
	return widgets.HitTest(w.content, local)
}

func (w *Window) OwnsControl(control widgets.Control) bool {
	return widget.IsDescendant(w.content, control)
}

func (w *Window) FocusableControls(th *theme.Theme) []widgets.Control {
	w.syncContentBounds(th)
	return widgets.FocusableControls(w.content)
}

func (w *Window) ControlLocalPoint(control widgets.Control, screenPoint geom.Point, th *theme.Theme) geom.Point {
	client := w.ClientRect(th)
	abs := widget.AbsoluteBounds(control)
	return geom.Point{
		X: screenPoint.X - client.X - abs.X,
		Y: screenPoint.Y - client.Y - abs.Y,
	}
}

func (w *Window) Paint(canvas *paint.Canvas, th *theme.Theme, tr *uitext.Renderer) error {
	if !w.Visible() {
		return nil
	}

	bounds := w.Bounds()
	canvas.FillRect(bounds, th.Colors.Face)
	canvas.DrawDoubleBevel(bounds, th.Colors.Lightest, th.Colors.DarkShadow, th.Colors.Light, th.Colors.Shadow)

	caption := w.CaptionRect(th)
	captionColor := th.Colors.InactiveCaption
	if w.active {
		captionColor = th.Colors.ActiveCaption
	}
	canvas.FillRect(caption, captionColor)
	canvas.DrawHLine(caption.X, caption.Bottom(), caption.W, th.Colors.Shadow)

	drawCaptionIcon(canvas, caption, th)
	closeRect := w.CloseButtonRect(th)
	drawCaptionButton(canvas, closeRect, th, w.closePressed, w.closeHot)
	if err := w.paintTitle(canvas, caption, closeRect, th, tr); err != nil {
		return err
	}
	if err := w.paintMenuBar(canvas, th, tr); err != nil {
		return err
	}

	client := w.ClientRect(th)
	canvas.FillRect(client, th.Colors.Window)
	canvas.FrameRect(client, th.Colors.Shadow)

	w.syncContentBounds(th)
	canvas.PushClip(client)
	err := w.content.Paint(widgets.PaintContext{
		Canvas: canvas,
		Theme:  th,
		Text:   tr,
		Origin: geom.Point{X: client.X, Y: client.Y},
	})
	canvas.PopClip()
	return err
}

func (w *Window) paintTitle(canvas *paint.Canvas, caption, closeRect geom.Rect, th *theme.Theme, tr *uitext.Renderer) error {
	if tr == nil || w.title == "" {
		return nil
	}

	titleRect := geom.Rect{
		X: caption.X + th.Metrics.IconSizeSmall + 6,
		Y: caption.Y,
		W: closeRect.X - (caption.X + th.Metrics.IconSizeSmall + 10),
		H: caption.H,
	}
	if titleRect.W <= 0 {
		return nil
	}

	size := tr.MeasureString(w.title)
	textY := titleRect.Y
	if titleRect.H > size.H {
		textY += (titleRect.H - size.H) / 2
	}

	canvas.PushClip(titleRect)
	err := tr.DrawString(canvas, geom.Point{X: titleRect.X, Y: textY}, w.title, th.Colors.CaptionText)
	canvas.PopClip()
	return err
}

func (w *Window) syncContentBounds(th *theme.Theme) {
	client := w.ClientRect(th)
	w.content.SetBounds(geom.Rect{X: 0, Y: 0, W: client.W, H: client.H})
}

func drawCaptionIcon(canvas *paint.Canvas, caption geom.Rect, th *theme.Theme) {
	iconRect := geom.Rect{
		X: caption.X + 2,
		Y: caption.Y + 1,
		W: min(th.Metrics.IconSizeSmall-2, caption.H-2),
		H: min(th.Metrics.IconSizeSmall-2, caption.H-2),
	}
	canvas.FillRect(iconRect, th.Colors.Face)
	canvas.DrawDoubleBevel(iconRect, th.Colors.Lightest, th.Colors.DarkShadow, th.Colors.Light, th.Colors.Shadow)

	mark := iconRect.Inset(4)
	if mark.W > 0 && mark.H > 0 {
		canvas.FillRect(mark, th.Colors.Highlight)
		canvas.DrawBevel(mark, blend(th.Colors.Highlight, th.Colors.Lightest), th.Colors.DarkShadow)
	}
}

func drawCaptionButton(canvas *paint.Canvas, rect geom.Rect, th *theme.Theme, pressed, hot bool) {
	fill := th.Colors.Face
	if hot {
		fill = blend(fill, th.Colors.Lightest)
	}
	canvas.FillRect(rect, fill)
	if pressed {
		canvas.DrawDoubleBevel(rect, th.Colors.Shadow, th.Colors.Lightest, th.Colors.DarkShadow, th.Colors.Light)
	} else {
		canvas.DrawDoubleBevel(rect, th.Colors.Lightest, th.Colors.DarkShadow, th.Colors.Light, th.Colors.Shadow)
	}
	drawCloseGlyph(canvas, rect, th.Colors.DarkShadow, pressed)
}

func drawCloseGlyph(canvas *paint.Canvas, rect geom.Rect, color uicolor.RGBA, pressed bool) {
	offset := 0
	if pressed {
		offset = 1
	}
	startX := rect.X + rect.W/2 - 3 + offset
	startY := rect.Y + rect.H/2 - 3 + offset
	for i := 0; i < 7; i++ {
		canvas.DrawPixel(startX+i, startY+i, color)
		canvas.DrawPixel(startX+6-i, startY+i, color)
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
