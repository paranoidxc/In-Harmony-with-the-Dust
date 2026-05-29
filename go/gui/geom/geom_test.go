package geom

import "testing"

func TestIntersect(t *testing.T) {
	got, ok := Intersect(
		Rect{X: 10, Y: 10, W: 20, H: 20},
		Rect{X: 20, Y: 5, W: 20, H: 20},
	)
	if !ok {
		t.Fatal("expected intersection")
	}

	want := Rect{X: 20, Y: 10, W: 10, H: 15}
	if got != want {
		t.Fatalf("intersection = %+v, want %+v", got, want)
	}
}

func TestUnion(t *testing.T) {
	got := Union(
		Rect{X: 10, Y: 10, W: 20, H: 20},
		Rect{X: 20, Y: 5, W: 20, H: 20},
	)

	want := Rect{X: 10, Y: 5, W: 30, H: 25}
	if got != want {
		t.Fatalf("union = %+v, want %+v", got, want)
	}
}
