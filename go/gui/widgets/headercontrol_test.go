package widgets

import (
	"testing"
	"time"

	"classicui/event"
	"classicui/geom"
)

func TestHeaderControlClickDispatchesColumnIndex(t *testing.T) {
	header := NewHeaderControl("hdr", geom.Rect{X: 0, Y: 0, W: 180, H: 20},
		HeaderColumn{Title: "Name", Width: 100},
		HeaderColumn{Title: "Size", Width: 80, Align: HeaderAlignRight},
	)
	ctx := &fakeContext{}

	clicked := -1
	header.OnColumnClick(func(index int) {
		clicked = index
	})

	header.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, geom.Point{X: 110, Y: 8})
	header.MouseUp(ctx, event.MouseButtonEvent{Button: event.MouseButtonLeft}, geom.Point{X: 110, Y: 8})

	if clicked != 1 {
		t.Fatalf("clicked column = %d, want 1", clicked)
	}
	if ctx.focused != header {
		t.Fatal("header should take focus on click")
	}
	if ctx.releaseCapture == 0 {
		t.Fatal("header should release capture after click")
	}
}

func TestHeaderControlDividerDragResizesColumn(t *testing.T) {
	header := NewHeaderControl("hdr", geom.Rect{X: 0, Y: 0, W: 180, H: 20},
		HeaderColumn{Title: "Name", Width: 100},
		HeaderColumn{Title: "Size", Width: 80, Align: HeaderAlignRight},
	)
	ctx := &fakeContext{}

	var resizedIndex, resizedWidth int
	header.OnColumnResize(func(index, width int) {
		resizedIndex = index
		resizedWidth = width
	})

	header.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, geom.Point{X: 100, Y: 8})
	header.MouseMove(ctx, geom.Point{X: 128, Y: 8})
	header.MouseUp(ctx, event.MouseButtonEvent{Button: event.MouseButtonLeft}, geom.Point{X: 128, Y: 8})

	if resizedIndex != 0 {
		t.Fatalf("resized index = %d, want 0", resizedIndex)
	}
	if resizedWidth != 128 {
		t.Fatalf("resized width = %d, want 128", resizedWidth)
	}
	if header.ColumnWidth(0) != 128 {
		t.Fatalf("column width = %d, want 128", header.ColumnWidth(0))
	}
}

func TestHeaderControlDividerDoubleClickTriggersAutoFit(t *testing.T) {
	header := NewHeaderControl("hdr", geom.Rect{X: 0, Y: 0, W: 180, H: 20},
		HeaderColumn{Title: "Name", Width: 100},
		HeaderColumn{Title: "Size", Width: 80, Align: HeaderAlignRight},
	)
	ctx := &fakeContext{}

	calls := 0
	header.now = func() time.Time {
		return time.Date(2026, time.May, 30, 10, 0, 0, 0, time.UTC)
	}
	header.OnColumnAutoFit(func(_ EventContext, index int) {
		calls++
		if index != 0 {
			t.Fatalf("auto-fit index = %d, want 0", index)
		}
	})

	point := geom.Point{X: 100, Y: 8}
	header.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, point)
	header.MouseUp(ctx, event.MouseButtonEvent{Button: event.MouseButtonLeft}, point)
	header.now = func() time.Time {
		return time.Date(2026, time.May, 30, 10, 0, 0, int(200*time.Millisecond), time.UTC)
	}
	header.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, point)

	if calls != 1 {
		t.Fatalf("auto-fit calls = %d, want 1", calls)
	}
}
