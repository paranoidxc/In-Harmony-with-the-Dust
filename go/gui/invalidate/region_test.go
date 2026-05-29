package invalidate

import (
	"testing"

	"classicui/geom"
)

func TestRegionAddMergesBounds(t *testing.T) {
	var region Region
	region.Add(geom.Rect{X: 10, Y: 10, W: 20, H: 20})
	region.Add(geom.Rect{X: 40, Y: 5, W: 10, H: 10})

	if !region.Any() {
		t.Fatal("expected region to be dirty")
	}

	want := geom.Rect{X: 10, Y: 5, W: 40, H: 25}
	if got := region.Bounds(); got != want {
		t.Fatalf("bounds = %+v, want %+v", got, want)
	}
}
