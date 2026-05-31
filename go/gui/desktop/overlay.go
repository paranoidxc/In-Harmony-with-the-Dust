package desktop

import (
	"classicui/event"
	"classicui/geom"
	"classicui/paint"
	"classicui/widget"
	"classicui/widgets"
)

type desktopOverlay interface {
	Bounds() geom.Rect
	Paint(*Desktop, *paint.Canvas) error
}

type popupHost interface {
	desktopOverlay
	popupState() *popupHostState
}

type popupHostState struct {
	ownerWindow    *Window
	owner          widgets.Control
	rect           geom.Rect
	closeOnOutside bool
	kind           widgets.PopupKind
	onClose        func()
}

type controlOverlayState struct {
	popupHostState
	content widgets.Control
	hovered widgets.Control
	focused widgets.Control
}

type overlayContext struct {
	desktop *Desktop
	overlay *controlOverlayState
}

func (s *popupHostState) popupState() *popupHostState {
	return s
}

func (s *popupHostState) Bounds() geom.Rect {
	if s == nil {
		return geom.Rect{}
	}
	return s.rect
}

func (o *controlOverlayState) Paint(d *Desktop, canvas *paint.Canvas) error {
	if o == nil || o.content == nil {
		return nil
	}
	return o.content.Paint(widgets.PaintContext{
		Canvas: canvas,
		Theme:  d.theme,
		Text:   d.text,
		Origin: geom.Point{X: o.rect.X, Y: o.rect.Y},
	})
}

func (d *Desktop) paintOverlays(canvas *paint.Canvas) error {
	for _, overlay := range d.overlays {
		if overlay == nil {
			continue
		}
		if err := overlay.Paint(d, canvas); err != nil {
			return err
		}
	}
	return nil
}

func (d *Desktop) pushOverlay(overlay desktopOverlay) {
	if overlay == nil {
		return
	}
	d.overlays = append(d.overlays, overlay)
	d.InvalidateRect(overlay.Bounds())
}

func (d *Desktop) removeOverlay(target desktopOverlay) {
	if target == nil {
		return
	}
	for i := len(d.overlays) - 1; i >= 0; i-- {
		if d.overlays[i] != target {
			continue
		}
		d.InvalidateRect(target.Bounds())
		d.overlays = append(d.overlays[:i], d.overlays[i+1:]...)
		return
	}
}

func (d *Desktop) clearOverlays() {
	for _, overlay := range d.overlays {
		if overlay != nil {
			d.InvalidateRect(overlay.Bounds())
		}
	}
	d.overlays = nil
}

func (d *Desktop) overlayAt(point geom.Point) (desktopOverlay, int) {
	for i := len(d.overlays) - 1; i >= 0; i-- {
		overlay := d.overlays[i]
		if overlay != nil && overlay.Bounds().Contains(point) {
			return overlay, i
		}
	}
	return nil, -1
}

func (d *Desktop) fitOverlayOrigin(anchorRect geom.Rect, size geom.Size, placement widgets.OverlayPlacement) geom.Point {
	x := anchorRect.X
	y := anchorRect.Bottom() - 1
	if placement == widgets.OverlayRightTop {
		x = anchorRect.Right() - 3
		y = anchorRect.Y - 3
	}

	if placement == widgets.OverlayRightTop && x+size.W > d.bounds.Right() {
		x = anchorRect.X - size.W + 3
	}
	if x+size.W > d.bounds.Right() {
		x = d.bounds.Right() - size.W
	}
	if y+size.H > d.bounds.Bottom() {
		y = d.bounds.Bottom() - size.H
	}
	if x < d.bounds.X {
		x = d.bounds.X
	}
	if y < d.bounds.Y {
		y = d.bounds.Y
	}

	return geom.Point{X: x, Y: y}
}

func (d *Desktop) showControlOverlay(win *Window, request widgets.OverlayRequest) bool {
	if win == nil || request.Owner == nil || request.Content == nil {
		return false
	}
	if request.Kind == widgets.PopupKindInteractive {
		d.clearBaseInteraction()
	}
	d.hideControlOverlay(request.Owner, false)

	size := request.Content.Bounds()
	if size.W <= 0 || size.H <= 0 {
		return false
	}

	ownerRect := d.controlScreenRect(win, request.Owner)
	anchor := request.Anchor.Move(ownerRect.X, ownerRect.Y)
	origin := d.fitOverlayOrigin(anchor, geom.Size{W: size.W, H: size.H}, request.Placement)
	request.Content.SetBounds(geom.Rect{X: 0, Y: 0, W: size.W, H: size.H})
	overlay := &controlOverlayState{
		popupHostState: popupHostState{
			ownerWindow:    win,
			owner:          request.Owner,
			rect:           geom.Rect{X: origin.X, Y: origin.Y, W: size.W, H: size.H},
			closeOnOutside: request.CloseOnOutside,
			kind:           request.Kind,
			onClose:        request.OnClose,
		},
		content: request.Content,
	}
	d.pushOverlay(overlay)
	return true
}

func (d *Desktop) hideControlOverlay(owner widgets.Control, notify bool) bool {
	if owner == nil {
		return false
	}
	for i := len(d.overlays) - 1; i >= 0; i-- {
		overlay, ok := d.overlays[i].(*controlOverlayState)
		if !ok || overlay.owner != owner {
			continue
		}
		d.dismissControlOverlay(overlay, notify)
		return true
	}
	return false
}

func (d *Desktop) dismissControlOverlay(overlay *controlOverlayState, notify bool) {
	if overlay == nil {
		return
	}
	if d.hoveredOverlay == overlay {
		d.setHoveredOverlayControl(nil, nil)
	}
	if d.focusedOverlay == overlay {
		d.setFocusedOverlayControl(nil, nil)
	}
	if d.captureOverlay == overlay {
		d.captureOverlay = nil
		d.captureOverlayControl = nil
	}
	d.dismissPopupHost(overlay, notify)
}

func (d *Desktop) overlayVisible(owner widgets.Control) bool {
	if owner == nil {
		return false
	}
	for _, raw := range d.overlays {
		overlay, ok := raw.(popupHost)
		if ok && overlay.popupState().owner == owner {
			return true
		}
	}
	return false
}

func (d *Desktop) dismissPopupHost(host popupHost, notify bool) {
	if host == nil {
		return
	}
	d.removeOverlay(host)
	state := host.popupState()
	if notify && state != nil && state.onClose != nil {
		state.onClose()
	}
}

func (d *Desktop) topPopupHost() popupHost {
	for i := len(d.overlays) - 1; i >= 0; i-- {
		if overlay, ok := d.overlays[i].(popupHost); ok {
			return overlay
		}
	}
	return nil
}

func (d *Desktop) popupHostVisible() bool {
	return d.topPopupHost() != nil
}

func (d *Desktop) topInteractivePopupHost() popupHost {
	for i := len(d.overlays) - 1; i >= 0; i-- {
		overlay, ok := d.overlays[i].(popupHost)
		if !ok {
			continue
		}
		state := overlay.popupState()
		if state != nil && state.kind == widgets.PopupKindInteractive {
			return overlay
		}
	}
	return nil
}

func (d *Desktop) popupHostAt(point geom.Point) (popupHost, int) {
	for i := len(d.overlays) - 1; i >= 0; i-- {
		overlay, ok := d.overlays[i].(popupHost)
		if ok && overlay.Bounds().Contains(point) {
			return overlay, i
		}
	}
	return nil, -1
}

func (d *Desktop) popupOwnerContains(host popupHost, point geom.Point) bool {
	if host == nil {
		return false
	}
	state := host.popupState()
	if state == nil || state.ownerWindow == nil || state.owner == nil {
		return false
	}
	return d.controlScreenRect(state.ownerWindow, state.owner).Contains(point)
}

func (d *Desktop) pointWithinPopupHostOrOwner(host popupHost, point geom.Point) bool {
	if host == nil {
		return false
	}
	if host.Bounds().Contains(point) {
		return true
	}
	return d.popupOwnerContains(host, point)
}

func (d *Desktop) topControlOverlay() *controlOverlayState {
	for i := len(d.overlays) - 1; i >= 0; i-- {
		if overlay, ok := d.overlays[i].(*controlOverlayState); ok {
			return overlay
		}
	}
	return nil
}

func (d *Desktop) controlOverlayAt(point geom.Point) (*controlOverlayState, int) {
	for i := len(d.overlays) - 1; i >= 0; i-- {
		overlay, ok := d.overlays[i].(*controlOverlayState)
		if ok && overlay.Bounds().Contains(point) {
			return overlay, i
		}
	}
	return nil, -1
}

func (d *Desktop) overlayControlLocalPoint(overlay *controlOverlayState, control widgets.Control, screenPoint geom.Point) geom.Point {
	abs := widget.AbsoluteBounds(control)
	return geom.Point{
		X: screenPoint.X - overlay.rect.X - abs.X,
		Y: screenPoint.Y - overlay.rect.Y - abs.Y,
	}
}

func (d *Desktop) overlayControlScreenRect(overlay *controlOverlayState, control widgets.Control) geom.Rect {
	abs := widget.AbsoluteBounds(control)
	return abs.Move(overlay.rect.X, overlay.rect.Y)
}

func (d *Desktop) overlayContextFor(overlay *controlOverlayState) overlayContext {
	return overlayContext{desktop: d, overlay: overlay}
}

func (d *Desktop) clearOverlayInteraction() {
	if d.captureOverlay != nil {
		d.captureOverlay = nil
		d.captureOverlayControl = nil
	}
	if d.focusedOverlay != nil || d.focusedOverlayControl != nil {
		d.setFocusedOverlayControl(nil, nil)
	}
	if d.hoveredOverlay != nil || d.hoveredOverlayControl != nil {
		d.setHoveredOverlayControl(nil, nil)
	}
}

func (d *Desktop) setHoveredOverlayControl(overlay *controlOverlayState, control widgets.Control) {
	if d.hoveredOverlay == overlay && d.hoveredOverlayControl == control {
		return
	}
	if d.hoveredOverlay != nil && d.hoveredOverlayControl != nil {
		d.hoveredOverlayControl.MouseLeave(d.overlayContextFor(d.hoveredOverlay))
	}
	d.hoveredOverlay = overlay
	d.hoveredOverlayControl = control
	if d.hoveredOverlay != nil && d.hoveredOverlayControl != nil {
		d.hoveredOverlayControl.MouseEnter(d.overlayContextFor(d.hoveredOverlay))
	}
}

func (d *Desktop) setFocusedOverlayControl(overlay *controlOverlayState, control widgets.Control) {
	if d.focusedOverlay == overlay && d.focusedOverlayControl == control {
		return
	}
	if d.focusedOverlay != nil && d.focusedOverlayControl != nil {
		old := d.focusedOverlayControl
		oldOverlay := d.focusedOverlay
		old.SetFocused(false)
		if handler, ok := old.(widgets.FocusHandler); ok {
			handler.FocusLost(d.overlayContextFor(oldOverlay))
		}
		d.InvalidateRect(d.overlayControlScreenRect(oldOverlay, old))
	}
	d.focusedOverlay = overlay
	d.focusedOverlayControl = control
	if d.focusedOverlay != nil && d.focusedOverlayControl != nil {
		d.focusedOverlayControl.SetFocused(true)
		if handler, ok := d.focusedOverlayControl.(widgets.FocusHandler); ok {
			handler.FocusGained(d.overlayContextFor(d.focusedOverlay))
		}
		d.InvalidateRect(d.overlayControlScreenRect(d.focusedOverlay, d.focusedOverlayControl))
	}
}

func (d *Desktop) closeControlOverlays() {
	for i := len(d.overlays) - 1; i >= 0; i-- {
		overlay, ok := d.overlays[i].(*controlOverlayState)
		if ok {
			d.dismissControlOverlay(overlay, true)
		}
	}
}

func (d *Desktop) handleControlOverlayMouseMove(point geom.Point) bool {
	if d.captureOverlay != nil && d.captureOverlayControl != nil {
		control := d.captureOverlayControl
		overlay := d.captureOverlay
		control.MouseMove(d.overlayContextFor(overlay), d.overlayControlLocalPoint(overlay, control, point))
		return true
	}

	if d.topInteractivePopupHost() == nil {
		return false
	}

	overlay, _ := d.controlOverlayAt(point)
	if overlay == nil {
		d.setHoveredOverlayControl(nil, nil)
		return true
	}
	control := widgets.HitTest(overlay.content, geom.Point{X: point.X - overlay.rect.X, Y: point.Y - overlay.rect.Y})
	d.setHoveredOverlayControl(overlay, control)
	if control != nil {
		control.MouseMove(d.overlayContextFor(overlay), d.overlayControlLocalPoint(overlay, control, point))
	}
	return true
}

func (d *Desktop) handleControlOverlayMouseDown(e event.MouseButtonEvent) (*controlOverlayState, bool) {
	overlay := d.topControlOverlay()
	if overlay == nil {
		return nil, false
	}

	hitOverlay, _ := d.controlOverlayAt(e.Position)
	if hitOverlay == nil {
		if d.pointWithinPopupHostOrOwner(overlay, e.Position) {
			return nil, false
		}
		if overlay.closeOnOutside {
			d.dismissControlOverlay(overlay, true)
			return overlay, false
		}
		return nil, false
	}

	control := widgets.HitTest(hitOverlay.content, geom.Point{X: e.Position.X - hitOverlay.rect.X, Y: e.Position.Y - hitOverlay.rect.Y})
	d.setHoveredOverlayControl(hitOverlay, control)
	if control == nil {
		return nil, true
	}
	if control.CanFocus() {
		d.setFocusedOverlayControl(hitOverlay, control)
	}
	control.MouseDown(d.overlayContextFor(hitOverlay), e, d.overlayControlLocalPoint(hitOverlay, control, e.Position))
	return nil, true
}

func (d *Desktop) handleControlOverlayMouseUp(e event.MouseButtonEvent) bool {
	if d.captureOverlay != nil && d.captureOverlayControl != nil {
		control := d.captureOverlayControl
		overlay := d.captureOverlay
		control.MouseUp(d.overlayContextFor(overlay), e, d.overlayControlLocalPoint(overlay, control, e.Position))
		if next, _ := d.controlOverlayAt(e.Position); next == nil {
			d.setHoveredOverlayControl(nil, nil)
		}
		return true
	}

	overlay := d.topControlOverlay()
	if overlay == nil {
		return false
	}

	hitOverlay, _ := d.controlOverlayAt(e.Position)
	if hitOverlay == nil {
		return true
	}
	control := widgets.HitTest(hitOverlay.content, geom.Point{X: e.Position.X - hitOverlay.rect.X, Y: e.Position.Y - hitOverlay.rect.Y})
	d.setHoveredOverlayControl(hitOverlay, control)
	if control != nil {
		control.MouseUp(d.overlayContextFor(hitOverlay), e, d.overlayControlLocalPoint(hitOverlay, control, e.Position))
	}
	return true
}

func (d *Desktop) handleControlOverlayMouseWheel(e event.MouseWheel) bool {
	overlay, _ := d.controlOverlayAt(e.Position)
	if overlay == nil {
		return d.topInteractivePopupHost() != nil
	}
	control := widgets.HitTest(overlay.content, geom.Point{X: e.Position.X - overlay.rect.X, Y: e.Position.Y - overlay.rect.Y})
	if control == nil {
		return true
	}
	handler, ok := control.(widgets.WheelHandler)
	if !ok {
		return true
	}
	handler.MouseWheel(d.overlayContextFor(overlay), e, d.overlayControlLocalPoint(overlay, control, e.Position))
	return true
}

func (d *Desktop) handleControlOverlayKeyDown(e event.KeyEvent) bool {
	if d.focusedOverlayControl != nil && d.focusedOverlay != nil {
		return d.focusedOverlayControl.KeyDown(d.overlayContextFor(d.focusedOverlay), e)
	}
	return false
}

func (c overlayContext) Invalidate(control widgets.Control) {
	if c.overlay == nil || control == nil {
		return
	}
	c.desktop.InvalidateRect(c.desktop.overlayControlScreenRect(c.overlay, control))
}

func (c overlayContext) SetFocus(control widgets.Control) {
	c.desktop.setFocusedOverlayControl(c.overlay, control)
}

func (c overlayContext) Capture(control widgets.Control) {
	c.desktop.captureOverlay = c.overlay
	c.desktop.captureOverlayControl = control
}

func (c overlayContext) ReleaseCapture(control widgets.Control) {
	if c.desktop.captureOverlayControl != control {
		return
	}
	c.desktop.captureOverlay = nil
	c.desktop.captureOverlayControl = nil
}

func (c overlayContext) DispatchCommand(cmd widgets.CommandID) {
	if c.overlay == nil {
		return
	}
	c.desktop.dispatchCommand(c.overlay.ownerWindow, cmd)
}

func (c overlayContext) ShowContextMenu(owner widgets.Control, anchor geom.Rect, menu *widgets.Menu) bool {
	if c.overlay == nil || c.overlay.ownerWindow == nil || owner == nil || menu == nil {
		return false
	}
	return c.desktop.showContextMenu(c.overlay.ownerWindow, owner, anchor, menu)
}

func (c overlayContext) ClipboardText() string {
	if c.desktop.platform == nil {
		return ""
	}
	return c.desktop.platform.ClipboardText()
}

func (c overlayContext) SetClipboardText(text string) {
	if c.desktop.platform == nil {
		return
	}
	c.desktop.platform.SetClipboardText(text)
}

func (c overlayContext) MeasureText(text string) geom.Size {
	if c.desktop.text == nil {
		return geom.Size{}
	}
	return c.desktop.text.MeasureString(text)
}

func (c overlayContext) LineHeight() int {
	if c.desktop.text == nil {
		return 0
	}
	return c.desktop.text.LineHeight()
}

func (c overlayContext) ShowPopup(request widgets.PopupRequest) bool {
	if c.overlay == nil || c.overlay.ownerWindow == nil {
		return false
	}
	return c.desktop.showControlOverlay(c.overlay.ownerWindow, request)
}

func (c overlayContext) HidePopup(owner widgets.Control) bool {
	return c.desktop.hideControlOverlay(owner, true)
}

func (c overlayContext) PopupVisible(owner widgets.Control) bool {
	return c.desktop.overlayVisible(owner)
}

func (c overlayContext) ShowOverlay(request widgets.OverlayRequest) bool {
	return c.ShowPopup(request)
}

func (c overlayContext) HideOverlay(owner widgets.Control) bool {
	return c.HidePopup(owner)
}

func (c overlayContext) OverlayVisible(owner widgets.Control) bool {
	return c.PopupVisible(owner)
}
