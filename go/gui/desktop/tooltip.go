package desktop

import (
	"time"

	"classicui/geom"
	"classicui/paint"
	"classicui/uicolor"
	"classicui/widgets"
)

const tooltipDelay = 700 * time.Millisecond

var (
	tooltipBackground = uicolor.RGBA{R: 0xFF, G: 0xFF, B: 0xE1, A: 0xFF}
	tooltipBorder     = uicolor.RGBA{R: 0x00, G: 0x00, B: 0x00, A: 0xFF}
)

type tooltipOverlayState struct {
	rect geom.Rect
	text string
}

func (t *tooltipOverlayState) Bounds() geom.Rect {
	return t.rect
}

func (t *tooltipOverlayState) Paint(d *Desktop, canvas *paint.Canvas) error {
	canvas.FillRect(t.rect, tooltipBackground)
	canvas.FrameRect(t.rect, tooltipBorder)
	if d.text == nil || t.text == "" {
		return nil
	}
	textY := t.rect.Y + 2
	return d.text.DrawString(canvas, geom.Point{X: t.rect.X + 3, Y: textY}, t.text, d.theme.Colors.WindowText)
}

func (d *Desktop) noteTooltipTarget(win *Window, control widgets.Control, local geom.Point, now time.Time) {
	provider, ok := control.(widgets.TooltipProvider)
	if !ok {
		d.clearTooltip()
		return
	}
	info := provider.TooltipAt(local, d.measureText)
	if info.Text == "" {
		d.clearTooltip()
		return
	}
	if info.Anchor.Empty() {
		info.Anchor = widgets.LocalRect(control)
	}
	controlRect := d.controlScreenRect(win, control)
	anchor := info.Anchor.Move(controlRect.X, controlRect.Y)
	if d.tooltipControl == control && d.tooltipText == info.Text && d.tooltipAnchor == anchor {
		return
	}
	d.hideTooltipOverlay()
	d.tooltipWindow = win
	d.tooltipControl = control
	d.tooltipText = info.Text
	d.tooltipAnchor = anchor
	d.tooltipDue = now.Add(tooltipDelay)
}

func (d *Desktop) updateTooltip(now time.Time) {
	if d.tooltipControl == nil || d.tooltipText == "" || d.menuMode || d.drag != nil || d.captureControl != nil || d.captureOverlay != nil {
		d.hideTooltipOverlay()
		return
	}
	if d.tooltipOverlay != nil || now.Before(d.tooltipDue) {
		return
	}
	d.showTooltipOverlay()
}

func (d *Desktop) clearTooltip() {
	d.hideTooltipOverlay()
	d.tooltipWindow = nil
	d.tooltipControl = nil
	d.tooltipText = ""
	d.tooltipAnchor = geom.Rect{}
	d.tooltipDue = time.Time{}
}

func (d *Desktop) hideTooltipOverlay() {
	if d.tooltipOverlay == nil {
		return
	}
	d.removeOverlay(d.tooltipOverlay)
	d.tooltipOverlay = nil
}

func (d *Desktop) showTooltipOverlay() {
	if d.tooltipText == "" {
		return
	}
	size := d.measureText(d.tooltipText)
	width := max(size.W+6, 16)
	height := max(size.H+4, 18)
	origin := d.fitOverlayOrigin(
		d.tooltipAnchor.Move(0, 2),
		geom.Size{W: width, H: height},
		widgets.OverlayBelowStart,
	)
	overlay := &tooltipOverlayState{
		rect: geom.Rect{X: origin.X, Y: origin.Y, W: width, H: height},
		text: d.tooltipText,
	}
	d.tooltipOverlay = overlay
	d.pushOverlay(overlay)
}
