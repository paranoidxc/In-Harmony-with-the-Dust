package widgets

import (
	"classicui/geom"
	"classicui/uicolor"
)

func dragMarqueeRect(itemsRect geom.Rect, start, current geom.Point) (geom.Rect, bool) {
	if itemsRect.Empty() {
		return geom.Rect{}, false
	}
	start = clampPointToRect(start, itemsRect)
	current = clampPointToRect(current, itemsRect)
	left := minInt(start.X, current.X)
	top := minInt(start.Y, current.Y)
	right := maxInt(start.X, current.X)
	bottom := maxInt(start.Y, current.Y)
	rect := geom.Rect{
		X: left,
		Y: top,
		W: right - left + 1,
		H: bottom - top + 1,
	}
	clipped, ok := geom.Intersect(rect, itemsRect)
	return clipped, ok
}

func paintSelectionMarquee(ctx PaintContext, rect geom.Rect) {
	if rect.Empty() {
		return
	}
	fill := blendColor(ctx.Theme.Colors.Window, ctx.Theme.Colors.Highlight)
	border := ctx.Theme.Colors.Highlight
	dots := contrastColor(border)

	inset := rect.Inset(1)
	if !inset.Empty() {
		for y := inset.Y; y < inset.Bottom(); y += 2 {
			startX := inset.X
			if (y-inset.Y)%4 != 0 {
				startX++
			}
			for x := startX; x < inset.Right(); x += 2 {
				ctx.Canvas.DrawPixel(x, y, fill)
			}
		}
	}
	ctx.Canvas.FrameRect(rect, border)
	if rect.W > 2 && rect.H > 2 {
		ctx.Canvas.DrawFocusRect(rect.Inset(1), dots)
	}
}

func clampPointToRect(point geom.Point, rect geom.Rect) geom.Point {
	return geom.Point{
		X: clampInt(point.X, rect.X, rect.Right()-1),
		Y: clampInt(point.Y, rect.Y, rect.Bottom()-1),
	}
}

func contrastColor(color uicolor.RGBA) uicolor.RGBA {
	luma := int(color.R)*299 + int(color.G)*587 + int(color.B)*114
	if luma >= 128000 {
		return uicolor.RGBA{A: 0xFF}
	}
	return uicolor.RGBA{R: 0xFF, G: 0xFF, B: 0xFF, A: 0xFF}
}
