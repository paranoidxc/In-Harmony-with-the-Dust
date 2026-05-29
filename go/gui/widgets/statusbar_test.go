package widgets

import (
	"testing"

	"classicui/geom"
)

func TestStatusBarLayoutUsesFlexibleLastPane(t *testing.T) {
	bar := NewStatusBar("status", geom.Rect{X: 0, Y: 0, W: 200, H: 22})
	bar.SetPanes([]StatusPane{
		{Text: "Ready", Width: 80},
		{Text: "3 items"},
	})

	rects := bar.layoutPanes(LocalRect(bar))
	if len(rects) != 2 {
		t.Fatalf("pane count = %d, want 2", len(rects))
	}
	if rects[0].W != 80 {
		t.Fatalf("first pane width = %d, want 80", rects[0].W)
	}
	if rects[1].Right() != LocalRect(bar).Right()-2 {
		t.Fatalf("last pane right = %d, want %d", rects[1].Right(), LocalRect(bar).Right()-2)
	}
}
