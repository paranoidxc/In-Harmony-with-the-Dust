package widgets

import (
	"testing"

	"classicui/geom"
)

func TestDragMarqueeRectNormalizesAndClips(t *testing.T) {
	itemsRect := geom.Rect{X: 2, Y: 2, W: 40, H: 30}

	rect, ok := dragMarqueeRect(itemsRect, geom.Point{X: 10, Y: 12}, geom.Point{X: 50, Y: 40})
	if !ok {
		t.Fatal("drag marquee rect should exist")
	}
	want := geom.Rect{X: 10, Y: 12, W: 32, H: 20}
	if rect != want {
		t.Fatalf("marquee rect = %#v, want %#v", rect, want)
	}
}

func TestListBoxMarqueeRectUsesItemsArea(t *testing.T) {
	list := NewListBox("files", geom.Rect{X: 0, Y: 0, W: 180, H: 120})
	list.SetItems([]string{"One", "Two", "Three", "Four"})
	list.dragSelecting = true
	list.pressedIndex = 1
	list.pressedStart = geom.Point{X: 12, Y: 18}
	list.pressedPoint = geom.Point{X: 48, Y: 150}

	rect, ok := list.marqueeRect()
	if !ok {
		t.Fatal("list marquee rect should exist")
	}
	want := geom.Rect{X: 12, Y: 18, W: 37, H: 100}
	if rect != want {
		t.Fatalf("list marquee rect = %#v, want %#v", rect, want)
	}
}

func TestTreeViewMarqueeRectUsesItemsArea(t *testing.T) {
	root := NewFolderNode("Root", NewFileNode("Child"))
	root.Expanded = true
	tree := NewTreeView("tree", geom.Rect{X: 0, Y: 0, W: 180, H: 120}, root)
	tree.dragSelecting = true
	tree.pressedNode = root
	tree.pressedPart = treeHotPartRow
	tree.pressedStart = geom.Point{X: 20, Y: 12}
	tree.pressedPoint = geom.Point{X: -10, Y: 42}

	rect, ok := tree.marqueeRect()
	if !ok {
		t.Fatal("tree marquee rect should exist")
	}
	want := geom.Rect{X: 2, Y: 12, W: 19, H: 31}
	if rect != want {
		t.Fatalf("tree marquee rect = %#v, want %#v", rect, want)
	}
}
