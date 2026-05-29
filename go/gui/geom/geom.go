package geom

type Point struct {
	X int
	Y int
}

type Size struct {
	W int
	H int
}

type Rect struct {
	X int
	Y int
	W int
	H int
}

func (r Rect) Empty() bool {
	return r.W <= 0 || r.H <= 0
}

func (r Rect) Right() int {
	return r.X + r.W
}

func (r Rect) Bottom() int {
	return r.Y + r.H
}

func (r Rect) Contains(p Point) bool {
	return p.X >= r.X && p.X < r.Right() && p.Y >= r.Y && p.Y < r.Bottom()
}

func (r Rect) Inset(n int) Rect {
	return Rect{
		X: r.X + n,
		Y: r.Y + n,
		W: r.W - 2*n,
		H: r.H - 2*n,
	}
}

func (r Rect) Move(dx, dy int) Rect {
	r.X += dx
	r.Y += dy
	return r
}

func Intersect(a, b Rect) (Rect, bool) {
	left := max(a.X, b.X)
	top := max(a.Y, b.Y)
	right := min(a.Right(), b.Right())
	bottom := min(a.Bottom(), b.Bottom())
	out := Rect{X: left, Y: top, W: right - left, H: bottom - top}
	return out, !out.Empty()
}

func Union(a, b Rect) Rect {
	if a.Empty() {
		return b
	}
	if b.Empty() {
		return a
	}

	left := min(a.X, b.X)
	top := min(a.Y, b.Y)
	right := max(a.Right(), b.Right())
	bottom := max(a.Bottom(), b.Bottom())
	return Rect{X: left, Y: top, W: right - left, H: bottom - top}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
