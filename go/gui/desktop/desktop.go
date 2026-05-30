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
	"time"
)

type Platform interface {
	ClipboardText() string
	SetClipboardText(string)
	SetTextInput(bool, geom.Rect)
}

type Desktop struct {
	bounds   geom.Rect
	theme    *theme.Theme
	text     *uitext.Renderer
	platform Platform
	dirty    invalidate.Region

	windows  []*Window
	menuMode bool

	drag *dragState

	hoveredWindow  *Window
	hoveredControl widgets.Control
	pointerPos     geom.Point

	focusedWindow  *Window
	focusedControl widgets.Control

	captureWindow  *Window
	captureControl widgets.Control

	hoveredOverlay        *controlOverlayState
	hoveredOverlayControl widgets.Control
	focusedOverlay        *controlOverlayState
	focusedOverlayControl widgets.Control
	captureOverlay        *controlOverlayState
	captureOverlayControl widgets.Control

	overlays       []desktopOverlay
	menuWindow     *Window
	menuPopups     []*popupMenuState
	commandHandler func(*Window, widgets.CommandID)

	tooltipWindow  *Window
	tooltipControl widgets.Control
	tooltipText    string
	tooltipAnchor  geom.Rect
	tooltipDue     time.Time
	tooltipOverlay *tooltipOverlayState
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

func (d *Desktop) BindPlatform(platform Platform) {
	d.platform = platform
	d.syncTextInputState()
}

func (d *Desktop) BindTextRenderer(renderer *uitext.Renderer) {
	d.text = renderer
	d.syncTextInputState()
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
		d.syncTextInputState()
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
	return d.paintOverlays(canvas)
}

func (d *Desktop) Update(now time.Time) {
	d.updateTooltip(now)
	if d.focusedControl == nil || d.focusedWindow == nil {
		d.syncTextInputState()
		return
	}
	tickable, ok := d.focusedControl.(widgets.TickHandler)
	if ok {
		tickable.Tick(d.contextFor(d.focusedWindow), now)
	}
	d.syncTextInputState()
}

func (d *Desktop) HandleEvent(evt event.Event) {
	defer d.syncTextInputState()

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
		if e.Button != event.MouseButtonLeft && e.Button != event.MouseButtonRight {
			return
		}
		if e.Down {
			d.handleMouseDown(e)
			return
		}
		d.handleMouseUp(e)
	case event.MouseWheel:
		d.handleMouseWheel(e)
	case event.KeyEvent:
		if !e.Down {
			return
		}
		d.handleKeyDown(e)
	case event.TextInput:
		d.handleTextInput(e)
	case event.TextEditing:
		d.handleTextEditing(e)
	}
}

func (d *Desktop) handleMouseMove(p geom.Point) {
	d.pointerPos = p
	if d.handleMenuMouseMove(p) {
		d.clearTooltip()
		return
	}
	if d.handleControlOverlayMouseMove(p) {
		d.setHoveredControl(nil, nil)
		d.clearTooltip()
		return
	}

	if d.drag != nil {
		d.clearTooltip()
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
		d.clearTooltip()
		d.captureControl.MouseMove(d.contextFor(d.captureWindow), d.captureWindow.ControlLocalPoint(d.captureControl, p, d.theme))
		return
	}

	if d.anyClosePressed() {
		d.clearTooltip()
		d.updateCloseHotState(p)
		return
	}

	win, _ := d.windowAt(p)
	if win == nil || win.HitTest(p, d.theme) != HitClient {
		d.setHoveredControl(nil, nil)
		d.clearTooltip()
		return
	}

	control := win.ControlAt(p, d.theme)
	d.setHoveredControl(win, control)
	if control != nil {
		local := win.ControlLocalPoint(control, p, d.theme)
		control.MouseMove(d.contextFor(win), local)
		d.noteTooltipTarget(win, control, local, time.Now())
		return
	}
	d.clearTooltip()
}

func (d *Desktop) handleMouseDown(e event.MouseButtonEvent) {
	d.pointerPos = e.Position
	d.clearTooltip()
	if d.handleMenuMouseDown(e.Position) {
		return
	}
	if closed, handled := d.handleControlOverlayMouseDown(e); handled {
		d.setHoveredControl(nil, nil)
		return
	} else if closed != nil {
		win, _ := d.windowAt(e.Position)
		if win != nil && win.HitTest(e.Position, d.theme) == HitClient {
			control := win.ControlAt(e.Position, d.theme)
			if control == closed.owner {
				d.setHoveredControl(nil, nil)
				return
			}
		}
	}

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
		if e.Button != event.MouseButtonLeft {
			return
		}
		if !win.closePressed {
			win.SetClosePressed(true)
			d.InvalidateRect(win.Bounds())
		}
		d.setFocus(nil, nil)
		d.updateCloseHotState(e.Position)
	case HitCaption:
		if e.Button != event.MouseButtonLeft {
			return
		}
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
			if e.Button == event.MouseButtonLeft {
				d.setFocus(nil, nil)
			}
			return
		}
		if e.Button == event.MouseButtonLeft {
			if control.CanFocus() {
				d.setFocus(win, control)
			} else {
				d.setFocus(nil, nil)
			}
		}
		control.MouseDown(d.contextFor(win), e, win.ControlLocalPoint(control, e.Position, d.theme))
	}
}

func (d *Desktop) handleMouseUp(e event.MouseButtonEvent) {
	d.pointerPos = e.Position
	d.clearTooltip()
	if d.handleMenuMouseUp(e) {
		return
	}
	if d.handleControlOverlayMouseUp(e) {
		d.setHoveredControl(nil, nil)
		return
	}

	if d.drag != nil {
		d.drag = nil
	}

	if e.Button == event.MouseButtonLeft && d.releasePressedClose(e.Position) {
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

func (d *Desktop) handleMouseWheel(e event.MouseWheel) {
	d.pointerPos = e.Position
	d.clearTooltip()
	if d.handleControlOverlayMouseWheel(e) {
		return
	}
	win, _ := d.windowAt(e.Position)
	if win == nil || win.HitTest(e.Position, d.theme) != HitClient {
		return
	}
	control := win.ControlAt(e.Position, d.theme)
	if control == nil {
		return
	}
	handler, ok := control.(widgets.WheelHandler)
	if !ok {
		return
	}
	handler.MouseWheel(d.contextFor(win), e, win.ControlLocalPoint(control, e.Position, d.theme))
}

func (d *Desktop) handleKeyDown(e event.KeyEvent) {
	d.clearTooltip()
	active := d.activeWindow()
	if active == nil {
		return
	}

	if d.handleMenuKeyDown(active, e) {
		return
	}
	if d.handleControlOverlayKeyDown(e) {
		return
	}

	if d.focusedControl != nil && d.focusedWindow == active {
		if d.focusedControl.KeyDown(d.contextFor(active), e) {
			return
		}
	}

	if e.Key == event.KeyTab && !e.Repeat {
		d.cycleFocus(active, e.Modifiers&event.ModShift != 0)
		return
	}

	if d.dispatchAccelerator(active, e) {
		return
	}

	if !e.Repeat && e.Key == event.KeyEnter && active.DefaultButton() != nil && active.DefaultButton().Enabled() {
		active.DefaultButton().KeyDown(d.contextFor(active), e)
		d.InvalidateRect(d.controlScreenRect(active, active.DefaultButton()))
	}
}

func (d *Desktop) handleTextInput(e event.TextInput) {
	if d.focusedControl == nil || d.focusedWindow == nil {
		return
	}
	handler, ok := d.focusedControl.(widgets.TextInputHandler)
	if !ok {
		return
	}
	handler.TextInput(d.contextFor(d.focusedWindow), e)
}

func (d *Desktop) handleTextEditing(e event.TextEditing) {
	if d.focusedControl == nil || d.focusedWindow == nil {
		return
	}
	handler, ok := d.focusedControl.(widgets.TextInputHandler)
	if !ok {
		return
	}
	handler.TextEditing(d.contextFor(d.focusedWindow), e)
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
	if d.menuWindow != nil && d.menuWindow != target {
		d.closeMenus()
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
		if handler, ok := oldControl.(widgets.FocusHandler); ok {
			handler.FocusLost(d.contextFor(oldWindow))
		}
		d.InvalidateRect(d.controlScreenRect(oldWindow, oldControl))
	}

	d.focusedWindow = win
	d.focusedControl = control
	if d.focusedControl != nil && d.focusedWindow != nil {
		d.focusedControl.SetFocused(true)
		if handler, ok := d.focusedControl.(widgets.FocusHandler); ok {
			handler.FocusGained(d.contextFor(d.focusedWindow))
		}
		d.InvalidateRect(d.controlScreenRect(d.focusedWindow, d.focusedControl))
	}
}

func (d *Desktop) clearStateForWindow(win *Window) {
	if d.menuWindow == win {
		d.closeMenus()
	}
	for i := len(d.overlays) - 1; i >= 0; i-- {
		overlay, ok := d.overlays[i].(*controlOverlayState)
		if ok && overlay.ownerWindow == win {
			d.dismissControlOverlay(overlay, true)
		}
	}
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
	if d.tooltipWindow == win {
		d.clearTooltip()
	}
}

func (d *Desktop) syncTextInputState() {
	if d.platform == nil {
		return
	}
	if d.menuMode {
		d.platform.SetTextInput(false, geom.Rect{})
		return
	}
	if d.focusedWindow == nil || d.focusedControl == nil {
		d.platform.SetTextInput(false, geom.Rect{})
		return
	}
	handler, ok := d.focusedControl.(widgets.TextInputHandler)
	if !ok {
		d.platform.SetTextInput(false, geom.Rect{})
		return
	}

	controlRect := d.controlScreenRect(d.focusedWindow, d.focusedControl)
	inputRect := handler.TextInputRect(d.contextFor(d.focusedWindow)).Move(controlRect.X, controlRect.Y)
	d.platform.SetTextInput(true, inputRect)
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

func (c controlContext) DispatchCommand(cmd widgets.CommandID) {
	c.desktop.dispatchCommand(c.window, cmd)
}

func (c controlContext) ShowContextMenu(owner widgets.Control, anchor geom.Rect, menu *widgets.Menu) bool {
	if c.window == nil || owner == nil || menu == nil {
		return false
	}
	return c.desktop.showContextMenu(c.window, owner, anchor, menu)
}

func (c controlContext) ClipboardText() string {
	if c.desktop.platform == nil {
		return ""
	}
	return c.desktop.platform.ClipboardText()
}

func (c controlContext) SetClipboardText(text string) {
	if c.desktop.platform == nil {
		return
	}
	c.desktop.platform.SetClipboardText(text)
}

func (c controlContext) MeasureText(text string) geom.Size {
	if c.desktop.text == nil {
		return geom.Size{}
	}
	return c.desktop.text.MeasureString(text)
}

func (c controlContext) LineHeight() int {
	if c.desktop.text == nil {
		return 0
	}
	return c.desktop.text.LineHeight()
}

func (c controlContext) ShowOverlay(request widgets.OverlayRequest) bool {
	if c.window == nil {
		return false
	}
	return c.desktop.showControlOverlay(c.window, request)
}

func (c controlContext) HideOverlay(owner widgets.Control) bool {
	return c.desktop.hideControlOverlay(owner, true)
}

func (c controlContext) OverlayVisible(owner widgets.Control) bool {
	return c.desktop.overlayVisible(owner)
}
