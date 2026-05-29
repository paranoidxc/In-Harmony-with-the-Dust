package invalidate

import "classicui/geom"

type Region struct {
	dirty  bool
	bounds geom.Rect
}

func (r *Region) Add(rect geom.Rect) {
	if rect.Empty() {
		return
	}
	if !r.dirty {
		r.bounds = rect
		r.dirty = true
		return
	}
	r.bounds = geom.Union(r.bounds, rect)
}

func (r *Region) Any() bool {
	return r.dirty
}

func (r *Region) Bounds() geom.Rect {
	return r.bounds
}

func (r *Region) Clear() {
	r.dirty = false
	r.bounds = geom.Rect{}
}
