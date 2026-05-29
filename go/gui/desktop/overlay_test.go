package desktop

import (
	"testing"

	"classicui/geom"
	"classicui/paint"
	"classicui/theme"
	"classicui/widgets"
)

type testOverlay struct {
	bounds geom.Rect
	paints int
}

func (o *testOverlay) Bounds() geom.Rect {
	return o.bounds
}

func (o *testOverlay) Paint(_ *Desktop, _ *paint.Canvas) error {
	o.paints++
	return nil
}

func TestDesktopOverlayStackPaintAndHitOrder(t *testing.T) {
	d := New(geom.Size{W: 320, H: 200}, theme.DefaultClassic())
	bottom := &testOverlay{bounds: geom.Rect{X: 10, Y: 10, W: 80, H: 40}}
	top := &testOverlay{bounds: geom.Rect{X: 20, Y: 15, W: 80, H: 40}}

	d.pushOverlay(bottom)
	d.pushOverlay(top)

	overlay, index := d.overlayAt(geom.Point{X: 24, Y: 20})
	if overlay != top || index != 1 {
		t.Fatalf("overlayAt returned (%v, %d), want top overlay at index 1", overlay, index)
	}

	canvas := paint.NewCanvas(320, 200)
	if err := d.Paint(canvas, nil); err != nil {
		t.Fatalf("paint failed: %v", err)
	}
	if bottom.paints != 1 || top.paints != 1 {
		t.Fatalf("paint counts = (%d, %d), want (1, 1)", bottom.paints, top.paints)
	}

	d.removeOverlay(top)
	overlay, index = d.overlayAt(geom.Point{X: 24, Y: 20})
	if overlay != bottom || index != 0 {
		t.Fatalf("overlayAt after remove returned (%v, %d), want bottom overlay at index 0", overlay, index)
	}
}

func TestDesktopFitOverlayOriginClampsToDesktop(t *testing.T) {
	d := New(geom.Size{W: 120, H: 100}, theme.DefaultClassic())

	below := d.fitOverlayOrigin(
		geom.Rect{X: 90, Y: 70, W: 20, H: 10},
		geom.Size{W: 40, H: 40},
		widgets.OverlayBelowStart,
	)
	if below.X != 80 || below.Y != 60 {
		t.Fatalf("below origin = %+v, want {X:80 Y:60}", below)
	}

	right := d.fitOverlayOrigin(
		geom.Rect{X: 95, Y: 30, W: 10, H: 10},
		geom.Size{W: 30, H: 30},
		widgets.OverlayRightTop,
	)
	if right.X != 68 || right.Y != 27 {
		t.Fatalf("right origin = %+v, want {X:68 Y:27}", right)
	}
}
