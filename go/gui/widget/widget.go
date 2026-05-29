package widget

import "classicui/geom"

type Widget interface {
	ID() string
	Bounds() geom.Rect
	SetBounds(geom.Rect)
	Visible() bool
	SetVisible(bool)
	Enabled() bool
	SetEnabled(bool)
	Parent() Widget
	SetParent(Widget)
	Children() []Widget
	AddChild(Widget)
}

type BaseWidget struct {
	id       string
	bounds   geom.Rect
	visible  bool
	enabled  bool
	parent   Widget
	children []Widget
}

func NewBase(id string, bounds geom.Rect) BaseWidget {
	return BaseWidget{
		id:      id,
		bounds:  bounds,
		visible: true,
		enabled: true,
	}
}

func (b *BaseWidget) ID() string {
	return b.id
}

func (b *BaseWidget) Bounds() geom.Rect {
	return b.bounds
}

func (b *BaseWidget) SetBounds(bounds geom.Rect) {
	b.bounds = bounds
}

func (b *BaseWidget) Visible() bool {
	return b.visible
}

func (b *BaseWidget) SetVisible(visible bool) {
	b.visible = visible
}

func (b *BaseWidget) Enabled() bool {
	return b.enabled
}

func (b *BaseWidget) SetEnabled(enabled bool) {
	b.enabled = enabled
}

func (b *BaseWidget) Parent() Widget {
	return b.parent
}

func (b *BaseWidget) SetParent(parent Widget) {
	b.parent = parent
}

func (b *BaseWidget) Children() []Widget {
	return b.children
}

func (b *BaseWidget) AddChild(child Widget) {
	b.AppendChild(child)
}

func (b *BaseWidget) AppendChild(child Widget) {
	if child == nil {
		return
	}
	b.children = append(b.children, child)
}

func AbsoluteBounds(w Widget) geom.Rect {
	if w == nil {
		return geom.Rect{}
	}
	rect := w.Bounds()
	for parent := w.Parent(); parent != nil; parent = parent.Parent() {
		parentRect := parent.Bounds()
		rect = rect.Move(parentRect.X, parentRect.Y)
	}
	return rect
}

func IsDescendant(root, target Widget) bool {
	if root == nil || target == nil {
		return false
	}
	for current := target; current != nil; current = current.Parent() {
		if current == root {
			return true
		}
	}
	return false
}
