package desktop

import (
	"classicui/event"
	"classicui/geom"
	"classicui/invalidate"
	"classicui/paint"
	uitext "classicui/text"
	"classicui/theme"
	"classicui/widget"
	"classicui/widgets"
)

type Desktop struct {
	bounds geom.Rect
	theme  *theme.Theme
	dirty  invalidate.Region

	windows []*Window

	drag *dragState

	hoveredWindow  *Window
	hoveredControl widgets.Control

	focusedWindow  *Window
	focusedControl widgets.Control

	captureWindow  *Window
	captureControl widgets.Control
}

type dragState struct {
	window       *Window
	startPointer geom.Point
	startBounds  geom.Rect
}

type controlContext struct {
	desktop *Desktop
	window  *Window
}

func New(size geom.Size, th *theme.Theme) *Desktop {
	d := &Desktop{
		bounds: geom.Rect{X: 0, Y: 0, W: size.W, H: size.H},
		theme:  th,
	}
	d.InvalidateAll()
	return d
}

func (d *Desktop) Theme() *theme.Theme {
	return d.theme
}

func (d *Desktop) Bounds() geom.Rect {
	return d.bounds
}

func (d *Desktop) AddWindow(win *Window) {
	if win == nil {
		return
	}
	d.windows = append(d.windows, win)
	d.activateWindow(win)
	d.InvalidateRect(win.Bounds())
}

func (d *Desktop) RemoveWindow(win *Window) {
	if win == nil {
		return
	}
	for i := range d.windows {
		if d.windows[i] != win {
			continue
		}

		oldBounds := d.windows[i].Bounds()
		d.clearStateForWindow(win)
		d.windows = append(d.windows[:i], d.windows[i+1:]...)
		d.InvalidateRect(oldBounds)
		if len(d.windows) > 0 {
			d.activateWindow(d.windows[len(d.windows)-1])
		}
		return
	}
}

func (d *Desktop) Windows() []*Window {
	return d.windows
}

func (d *Desktop) InvalidateRect(rect geom.Rect) {
	d.dirty.Add(rect)
}

func (d *Desktop) InvalidateAll() {
	d.dirty.Add(d.bounds)
}

func (d *Desktop) HasDirtyRegion() bool {
	return d.dirty.Any()
}

func (d *Desktop) ClearDirty() {
	d.dirty.Clear()
}

func (d *Desktop) Paint(canvas *paint.Canvas, tr *uitext.Renderer) error {
	canvas.ResetClip()
	canvas.Clear(d.theme.Colors.Desktop)
	for _, win := range d.windows {
		if err := win.Paint(canvas, d.theme, tr); err != nil {
			return err
		}
	}
	return nil
}

func (d *Desktop) HandleEvent(evt event.Event) {
	switch e := evt.(type) {
	case event.WindowExposed:
		d.InvalidateAll()
	case event.WindowResized:
		d.bounds.W = e.Size.W
		d.bounds.H = e.Size.H
		d.InvalidateAll()
	case event.MouseMove:
		d.handleMouseMove(e.Position)
	case event.MouseButtonEvent:
		if e.Button != event.MouseButtonLeft {
			return
		}
		if e.Down {
			d.handleMouseDown(e)
			return
		}
		d.handleMouseUp(e)
	case event.KeyEvent:
		if !e.Down || e.Repeat {
			return
		}
		d.handleKeyDown(e)
	}
}

func (d *Desktop) handleMouseMove(p geom.Point) {
	if d.drag != nil {
		oldBounds := d.drag.window.Bounds()
		next := d.drag.startBounds.Move(p.X-d.drag.startPointer.X, p.Y-d.drag.startPointer.Y)
		if next != oldBounds {
			d.InvalidateRect(oldBounds)
			d.drag.window.SetBounds(next)
			d.InvalidateRect(next)
		}
		return
	}

	if d.captureControl != nil && d.captureWindow != nil {
		d.captureControl.MouseMove(d.contextFor(d.captureWindow), d.captureWindow.ControlLocalPoint(d.captureControl, p, d.theme))
		return
	}

	if d.anyClosePressed() {
		d.updateCloseHotState(p)
		return
	}

	win, _ := d.windowAt(p)
	if win == nil || win.HitTest(p, d.theme) != HitClient {
		d.setHoveredControl(nil, nil)
		return
	}

	control := win.ControlAt(p, d.theme)
	d.setHoveredControl(win, control)
	if control != nil {
		control.MouseMove(d.contextFor(win), win.ControlLocalPoint(control, p, d.theme))
	}
}

func (d *Desktop) handleMouseDown(e event.MouseButtonEvent) {
	win, _ := d.windowAt(e.Position)
	if win == nil {
		d.setFocus(nil, nil)
		d.setHoveredControl(nil, nil)
		return
	}

	d.activateWindow(win)
	hit := win.HitTest(e.Position, d.theme)
	switch hit {
	case HitClose:
		if !win.closePressed {
			win.SetClosePressed(true)
			d.InvalidateRect(win.Bounds())
		}
		d.setFocus(nil, nil)
		d.updateCloseHotState(e.Position)
	case HitCaption:
		d.drag = &dragState{
			window:       win,
			startPointer: e.Position,
			startBounds:  win.Bounds(),
		}
		d.setFocus(nil, nil)
		d.setHoveredControl(nil, nil)
	case HitClient:
		control := win.ControlAt(e.Position, d.theme)
		d.setHoveredControl(win, control)
		if control == nil || control == win.Content() {
			d.setFocus(nil, nil)
			return
		}
		if control.CanFocus() {
			d.setFocus(win, control)
		} else {
			d.setFocus(nil, nil)
		}
		control.MouseDown(d.contextFor(win), e, win.ControlLocalPoint(control, e.Position, d.theme))
	}
}

func (d *Desktop) handleMouseUp(e event.MouseButtonEvent) {
	if d.drag != nil {
		d.drag = nil
	}

	if d.releasePressedClose(e.Position) {
		return
	}

	if d.captureControl != nil && d.captureWindow != nil {
		control := d.captureControl
		win := d.captureWindow
		control.MouseUp(d.contextFor(win), e, win.ControlLocalPoint(control, e.Position, d.theme))
		d.updateHoverFromPointer(e.Position)
		return
	}

	win, _ := d.windowAt(e.Position)
	if win == nil || win.HitTest(e.Position, d.theme) != HitClient {
		d.setHoveredControl(nil, nil)
		return
	}
	control := win.ControlAt(e.Position, d.theme)
	d.setHoveredControl(win, control)
	if control != nil {
		control.MouseUp(d.contextFor(win), e, win.ControlLocalPoint(control, e.Position, d.theme))
	}
}

func (d *Desktop) handleKeyDown(e event.KeyEvent) {
	active := d.activeWindow()
	if active == nil {
		return
	}

	if e.Key == event.KeyTab {
		d.cycleFocus(active, e.Modifiers&event.ModShift != 0)
		return
	}

	if d.focusedControl != nil && d.focusedWindow == active {
		if d.focusedControl.KeyDown(d.contextFor(active), e) {
			return
		}
	}

	if e.Key == event.KeyEnter && active.DefaultButton() != nil && active.DefaultButton().Enabled() {
		active.DefaultButton().KeyDown(d.contextFor(active), e)
		d.InvalidateRect(d.controlScreenRect(active, active.DefaultButton()))
	}
}

func (d *Desktop) cycleFocus(win *Window, reverse bool) {
	controls := win.FocusableControls(d.theme)
	if len(controls) == 0 {
		d.setFocus(nil, nil)
		return
	}

	index := -1
	for i, control := range controls {
		if control == d.focusedControl {
			index = i
			break
		}
	}

	if reverse {
		if index == -1 {
			index = 0
		}
		index = (index - 1 + len(controls)) % len(controls)
	} else {
		index = (index + 1) % len(controls)
	}

	d.setFocus(win, controls[index])
}

func (d *Desktop) anyClosePressed() bool {
	for _, win := range d.windows {
		if win.closePressed {
			return true
		}
	}
	return false
}

func (d *Desktop) releasePressedClose(p geom.Point) bool {
	for _, win := range d.windows {
		if !win.closePressed {
			continue
		}
		win.SetClosePressed(false)
		d.InvalidateRect(win.Bounds())
		if win.HitTest(p, d.theme) == HitClose {
			d.RemoveWindow(win)
		} else {
			d.updateCloseHotState(p)
		}
		return true
	}
	return false
}

func (d *Desktop) updateCloseHotState(p geom.Point) {
	top, _ := d.windowAt(p)
	for _, win := range d.windows {
		hot := win == top && win.HitTest(p, d.theme) == HitClose
		if hot == win.closeHot {
			continue
		}
		win.SetCloseHot(hot)
		d.InvalidateRect(win.Bounds())
	}
}

func (d *Desktop) activateWindow(target *Window) {
	if target == nil {
		return
	}

	idx := -1
	for i, win := range d.windows {
		active := win == target
		if win.active != active {
			win.SetActive(active)
			d.InvalidateRect(win.Bounds())
		}
		if active {
			idx = i
		}
	}

	if idx >= 0 && idx != len(d.windows)-1 {
		copy(d.windows[idx:], d.windows[idx+1:])
		d.windows[len(d.windows)-1] = target
		d.InvalidateAll()
	}
}

func (d *Desktop) activeWindow() *Window {
	if len(d.windows) == 0 {
		return nil
	}
	return d.windows[len(d.windows)-1]
}

func (d *Desktop) contextFor(win *Window) controlContext {
	return controlContext{
		desktop: d,
		window:  win,
	}
}

func (d *Desktop) setHoveredControl(win *Window, control widgets.Control) {
	if d.hoveredControl == control && d.hoveredWindow == win {
		return
	}

	if d.hoveredControl != nil && d.hoveredWindow != nil {
		d.hoveredControl.MouseLeave(d.contextFor(d.hoveredWindow))
	}
	d.hoveredWindow = win
	d.hoveredControl = control
	if d.hoveredControl != nil && d.hoveredWindow != nil {
		d.hoveredControl.MouseEnter(d.contextFor(d.hoveredWindow))
	}
}

func (d *Desktop) updateHoverFromPointer(p geom.Point) {
	win, _ := d.windowAt(p)
	if win == nil || win.HitTest(p, d.theme) != HitClient {
		d.setHoveredControl(nil, nil)
		return
	}
	d.setHoveredControl(win, win.ControlAt(p, d.theme))
}

func (d *Desktop) setFocus(win *Window, control widgets.Control) {
	if d.focusedControl == control && d.focusedWindow == win {
		return
	}

	if d.focusedControl != nil && d.focusedWindow != nil {
		oldControl := d.focusedControl
		oldWindow := d.focusedWindow
		oldControl.SetFocused(false)
		d.InvalidateRect(d.controlScreenRect(oldWindow, oldControl))
	}

	d.focusedWindow = win
	d.focusedControl = control
	if d.focusedControl != nil && d.focusedWindow != nil {
		d.focusedControl.SetFocused(true)
		d.InvalidateRect(d.controlScreenRect(d.focusedWindow, d.focusedControl))
	}
}

func (d *Desktop) clearStateForWindow(win *Window) {
	if d.drag != nil && d.drag.window == win {
		d.drag = nil
	}
	if d.hoveredWindow == win {
		d.setHoveredControl(nil, nil)
	}
	if d.focusedWindow == win {
		d.setFocus(nil, nil)
	}
	if d.captureWindow == win {
		d.captureWindow = nil
		d.captureControl = nil
	}
}

func (d *Desktop) controlScreenRect(win *Window, control widgets.Control) geom.Rect {
	client := win.ClientRect(d.theme)
	abs := widget.AbsoluteBounds(control)
	return abs.Move(client.X, client.Y)
}

func (d *Desktop) windowAt(p geom.Point) (*Window, int) {
	for i := len(d.windows) - 1; i >= 0; i-- {
		win := d.windows[i]
		if win.Visible() && win.Bounds().Contains(p) {
			return win, i
		}
	}
	return nil, -1
}

func (c controlContext) Invalidate(control widgets.Control) {
	if c.window == nil || control == nil {
		return
	}
	c.desktop.InvalidateRect(c.desktop.controlScreenRect(c.window, control))
}

func (c controlContext) SetFocus(control widgets.Control) {
	c.desktop.setFocus(c.window, control)
}

func (c controlContext) Capture(control widgets.Control) {
	c.desktop.captureWindow = c.window
	c.desktop.captureControl = control
}

func (c controlContext) ReleaseCapture(control widgets.Control) {
	if c.desktop.captureControl != control {
		return
	}
	c.desktop.captureControl = nil
	c.desktop.captureWindow = nil
}
