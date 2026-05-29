package desktop

import (
	"classicui/geom"
	"classicui/paint"
	uitext "classicui/text"
	"classicui/theme"
	"classicui/widgets"
)

const menuBarItemSpacing = 2

func (w *Window) SetMenuBar(menuBar *widgets.MenuBar) {
	w.menuBar = menuBar
	if w.menuBar == nil || len(w.menuBar.Items) == 0 {
		w.menuActiveIndex = -1
	}
}

func (w *Window) MenuBar() *widgets.MenuBar {
	return w.menuBar
}

func (w *Window) SetMenuBarActiveIndex(index int) {
	w.menuActiveIndex = index
}

func (w *Window) MenuBarActiveIndex() int {
	return w.menuActiveIndex
}

func (w *Window) MenuBarRect(th *theme.Theme) geom.Rect {
	if !w.hasMenuBar() {
		return geom.Rect{}
	}

	frame := th.Metrics.BorderWidth + th.Metrics.WindowFrameInner
	caption := w.CaptionRect(th)
	return geom.Rect{
		X: w.Bounds().X + frame,
		Y: caption.Bottom() + 1,
		W: w.Bounds().W - frame*2,
		H: th.Metrics.MenuHeight,
	}
}

func (w *Window) MenuBarItemRect(index int, th *theme.Theme, tr *uitext.Renderer) geom.Rect {
	if !w.hasMenuBar() || index < 0 || index >= len(w.menuBar.Items) {
		return geom.Rect{}
	}

	bar := w.MenuBarRect(th)
	x := bar.X + 6
	for i := 0; i < index; i++ {
		x += w.menuBarItemWidth(w.menuBar.Items[i], tr) + menuBarItemSpacing
	}

	return geom.Rect{
		X: x,
		Y: bar.Y + 1,
		W: w.menuBarItemWidth(w.menuBar.Items[index], tr),
		H: max(bar.H-2, 0),
	}
}

func (w *Window) MenuBarItemAt(screenPoint geom.Point, th *theme.Theme, tr *uitext.Renderer) int {
	if !w.MenuBarRect(th).Contains(screenPoint) || w.menuBar == nil {
		return -1
	}

	for i := range w.menuBar.Items {
		if w.MenuBarItemRect(i, th, tr).Contains(screenPoint) {
			return i
		}
	}
	return -1
}

func (w *Window) hasMenuBar() bool {
	return w.menuBar != nil && len(w.menuBar.Items) > 0
}

func (w *Window) menuBarItemWidth(item *widgets.MenuItem, tr *uitext.Renderer) int {
	size := measureMenuText(tr, item.DisplayText())
	return max(size.W+16, 20)
}

func (w *Window) paintMenuBar(canvas *paint.Canvas, th *theme.Theme, tr *uitext.Renderer) error {
	if !w.hasMenuBar() {
		return nil
	}

	bar := w.MenuBarRect(th)
	canvas.FillRect(bar, th.Colors.Face)
	canvas.DrawHLine(bar.X, bar.Bottom(), bar.W, th.Colors.Shadow)

	for i, item := range w.menuBar.Items {
		if item == nil {
			continue
		}

		itemRect := w.MenuBarItemRect(i, th, tr)
		if i == w.menuActiveIndex {
			canvas.FillRect(itemRect, th.Colors.Face)
			canvas.DrawDoubleBevel(itemRect, th.Colors.Shadow, th.Colors.Lightest, th.Colors.DarkShadow, th.Colors.Light)
		}

		if tr == nil || item.DisplayText() == "" {
			continue
		}

		textSize := tr.MeasureString(item.DisplayText())
		textY := itemRect.Y
		if itemRect.H > textSize.H {
			textY += (itemRect.H - textSize.H) / 2
		}

		canvas.PushClip(itemRect)
		err := tr.DrawString(canvas, geom.Point{X: itemRect.X + 8, Y: textY}, item.DisplayText(), th.Colors.WindowText)
		canvas.PopClip()
		if err != nil {
			return err
		}
	}

	return nil
}

func measureMenuText(tr *uitext.Renderer, text string) geom.Size {
	if tr != nil {
		return tr.MeasureString(text)
	}
	runes := len([]rune(text))
	return geom.Size{W: runes * 7, H: 14}
}
