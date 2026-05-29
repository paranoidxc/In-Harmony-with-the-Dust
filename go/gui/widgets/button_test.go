package widgets

import (
	"testing"

	"classicui/event"
	"classicui/geom"
)

type fakeContext struct {
	focused        Control
	captured       Control
	invalidations  int
	releaseCapture int
	clipboard      string
	commands       []CommandID
	overlay        *OverlayRequest
	overlayOwner   Control
}

func (f *fakeContext) Invalidate(Control) {
	f.invalidations++
}

func (f *fakeContext) SetFocus(control Control) {
	f.focused = control
}

func (f *fakeContext) Capture(control Control) {
	f.captured = control
}

func (f *fakeContext) ReleaseCapture(control Control) {
	if f.captured == control {
		f.captured = nil
	}
	f.releaseCapture++
}

func (f *fakeContext) DispatchCommand(cmd CommandID) {
	f.commands = append(f.commands, cmd)
}

func (f *fakeContext) ClipboardText() string {
	return f.clipboard
}

func (f *fakeContext) SetClipboardText(text string) {
	f.clipboard = text
}

func (f *fakeContext) MeasureText(text string) geom.Size {
	return geom.Size{W: len([]rune(text)) * 8, H: 14}
}

func (f *fakeContext) LineHeight() int {
	return 14
}

func (f *fakeContext) ShowOverlay(request OverlayRequest) bool {
	f.overlay = &request
	f.overlayOwner = request.Owner
	return true
}

func (f *fakeContext) HideOverlay(owner Control) bool {
	if f.overlayOwner != owner {
		return false
	}
	if f.overlay != nil && f.overlay.OnClose != nil {
		f.overlay.OnClose()
	}
	f.overlay = nil
	f.overlayOwner = nil
	return true
}

func (f *fakeContext) OverlayVisible(owner Control) bool {
	return f.overlayOwner == owner && f.overlay != nil
}

func TestButtonMouseCaptureAndClick(t *testing.T) {
	button := NewButton("ok", "OK", geom.Rect{X: 0, Y: 0, W: 80, H: 24})
	ctx := &fakeContext{}

	clicks := 0
	button.OnClick(func() {
		clicks++
	})

	button.MouseDown(ctx, event.MouseButtonEvent{Down: true, Button: event.MouseButtonLeft}, geom.Point{X: 4, Y: 4})
	if !button.pressed || !button.tracking {
		t.Fatal("button should enter pressed tracking state after mouse down")
	}
	if ctx.focused != button {
		t.Fatal("button should request focus on mouse down")
	}
	if ctx.captured != button {
		t.Fatal("button should capture mouse on mouse down")
	}

	button.MouseMove(ctx, geom.Point{X: 120, Y: 4})
	if button.pressed {
		t.Fatal("button should visually unpress when pointer leaves while captured")
	}

	button.MouseMove(ctx, geom.Point{X: 10, Y: 10})
	if !button.pressed {
		t.Fatal("button should re-press when pointer returns while captured")
	}

	button.MouseUp(ctx, event.MouseButtonEvent{Button: event.MouseButtonLeft}, geom.Point{X: 10, Y: 10})
	if clicks != 1 {
		t.Fatalf("clicks = %d, want 1", clicks)
	}
	if button.tracking || button.pressed {
		t.Fatal("button should leave tracking state after mouse up")
	}
	if ctx.captured != nil {
		t.Fatal("button should release capture after mouse up")
	}
}

func TestFocusableControlsSkipsLabels(t *testing.T) {
	root := NewPanel("root", geom.Rect{X: 0, Y: 0, W: 200, H: 100})
	root.Add(NewLabel("label", "Label", geom.Rect{X: 0, Y: 0, W: 80, H: 18}))
	first := NewButton("first", "First", geom.Rect{X: 0, Y: 20, W: 80, H: 24})
	second := NewButton("second", "Second", geom.Rect{X: 90, Y: 20, W: 80, H: 24})
	root.Add(first)
	root.Add(second)

	controls := FocusableControls(root)
	if len(controls) != 2 {
		t.Fatalf("focusable count = %d, want 2", len(controls))
	}
	if controls[0] != first || controls[1] != second {
		t.Fatal("focus traversal order should follow child insertion order")
	}
}
