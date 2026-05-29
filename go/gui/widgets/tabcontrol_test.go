package widgets

import (
	"testing"

	"classicui/event"
	"classicui/geom"
	"classicui/paint"
	"classicui/theme"
	"classicui/widget"
)

func TestTabControlSwitchesSelectionByMouseAndKeyboard(t *testing.T) {
	first := NewPanel("first", geom.Rect{})
	second := NewPanel("second", geom.Rect{})
	tabs := NewTabControl("tabs", geom.Rect{X: 0, Y: 0, W: 240, H: 140},
		NewTabPage("General", first),
		NewTabPage("Advanced", second),
	)
	ctx := &fakeContext{}

	if tabs.SelectedIndex() != 0 {
		t.Fatalf("initial selection = %d, want 0", tabs.SelectedIndex())
	}
	if !first.Visible() || second.Visible() {
		t.Fatal("only the selected page should be visible")
	}

	layout := tabs.layoutTabs(LocalRect(tabs), ctx.MeasureText, ctx.LineHeight())
	target := layout[1]
	click := geom.Point{X: target.X + target.W/2, Y: target.Y + target.H/2}
	tabs.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, click)
	tabs.MouseUp(ctx, event.MouseButtonEvent{Button: event.MouseButtonLeft}, click)

	if tabs.SelectedIndex() != 1 {
		t.Fatalf("selection after click = %d, want 1", tabs.SelectedIndex())
	}
	if first.Visible() || !second.Visible() {
		t.Fatal("selection should toggle page visibility")
	}
	if ctx.focused != tabs {
		t.Fatal("tab control should request focus on mouse selection")
	}

	tabs.KeyDown(ctx, event.KeyEvent{Key: event.KeyLeft})
	if tabs.SelectedIndex() != 0 {
		t.Fatalf("selection after left key = %d, want 0", tabs.SelectedIndex())
	}
}

func TestTabControlSkipsDisabledPagesAndResizesContent(t *testing.T) {
	first := NewPanel("first", geom.Rect{})
	second := NewPanel("second", geom.Rect{})
	third := NewPanel("third", geom.Rect{})
	disabled := NewTabPage("Disabled", second)
	disabled.Enabled = false

	tabs := NewTabControl("tabs", geom.Rect{X: 0, Y: 0, W: 200, H: 120},
		NewTabPage("One", first),
		disabled,
		NewTabPage("Three", third),
	)
	ctx := &fakeContext{}

	if !tabs.KeyDown(ctx, event.KeyEvent{Key: event.KeyRight}) {
		t.Fatal("right key should be handled")
	}
	if tabs.SelectedIndex() != 2 {
		t.Fatalf("selection after skipping disabled page = %d, want 2", tabs.SelectedIndex())
	}

	tabs.SetBounds(geom.Rect{X: 0, Y: 0, W: 260, H: 180})
	content := tabs.pageContentRect(LocalRect(tabs))
	if got := third.Bounds(); got != content {
		t.Fatalf("content bounds = %+v, want %+v", got, content)
	}
}

func TestFocusableControlsSkipHiddenTabPages(t *testing.T) {
	root := NewPanel("root", geom.Rect{X: 0, Y: 0, W: 300, H: 180})
	first := NewPanel("first", geom.Rect{})
	firstButton := NewButton("first.button", "First", geom.Rect{X: 8, Y: 8, W: 80, H: 24})
	first.Add(firstButton)

	second := NewPanel("second", geom.Rect{})
	secondButton := NewButton("second.button", "Second", geom.Rect{X: 8, Y: 8, W: 80, H: 24})
	second.Add(secondButton)

	tabs := NewTabControl("tabs", geom.Rect{X: 0, Y: 0, W: 240, H: 120},
		NewTabPage("One", first),
		NewTabPage("Two", second),
	)
	tabs.SetSelected(1)
	root.Add(tabs)

	controls := FocusableControls(root)
	for _, control := range controls {
		if control == firstButton {
			t.Fatal("hidden tab page descendants should not participate in focus traversal")
		}
	}
}

type paintProbe struct {
	widget.BaseWidget
	lastPaintRect geom.Rect
}

func newPaintProbe(id string, bounds geom.Rect) *paintProbe {
	return &paintProbe{BaseWidget: widget.NewBase(id, bounds)}
}

func (p *paintProbe) Paint(ctx PaintContext) error {
	p.lastPaintRect = ctx.BoundsFor(p)
	return nil
}

func (p *paintProbe) MouseEnter(EventContext)                                    {}
func (p *paintProbe) MouseLeave(EventContext)                                    {}
func (p *paintProbe) MouseMove(EventContext, geom.Point)                         {}
func (p *paintProbe) MouseDown(EventContext, event.MouseButtonEvent, geom.Point) {}
func (p *paintProbe) MouseUp(EventContext, event.MouseButtonEvent, geom.Point)   {}
func (p *paintProbe) KeyDown(EventContext, event.KeyEvent) bool                  { return false }
func (p *paintProbe) CanFocus() bool                                             { return false }
func (p *paintProbe) SetFocused(bool)                                            {}
func (p *paintProbe) Focused() bool                                              { return false }

func TestTabControlPaintUsesCorrectContentOrigin(t *testing.T) {
	probe := newPaintProbe("probe", geom.Rect{})
	tabs := NewTabControl("tabs", geom.Rect{X: 12, Y: 18, W: 240, H: 140},
		NewTabPage("Page", probe),
	)

	canvas := paint.NewCanvas(320, 240)
	err := tabs.Paint(PaintContext{
		Canvas: canvas,
		Theme:  theme.DefaultClassic(),
		Origin: geom.Point{},
	})
	if err != nil {
		t.Fatalf("paint failed: %v", err)
	}

	want := tabs.pageContentRect(tabs.Bounds())
	if probe.lastPaintRect != want {
		t.Fatalf("paint rect = %+v, want %+v", probe.lastPaintRect, want)
	}
}
