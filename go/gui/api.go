package classicui

import (
	"classicui/desktop"
	"classicui/geom"
	"classicui/theme"
)

type Point = geom.Point
type Size = geom.Size
type Rect = geom.Rect

type Theme = theme.Theme
type Window = desktop.Window

func DefaultClassicTheme() *theme.Theme {
	return theme.DefaultClassic().Clone()
}

func NewWindow(id string, bounds Rect) *desktop.Window {
	return desktop.NewWindow(id, bounds)
}
