package desktop

import (
	"classicui/event"
	"classicui/geom"
	"classicui/widgets"
)

type popupRouteResult struct {
	handled        bool
	clearBaseHover bool
}

type popupMouseDownResult struct {
	handled        bool
	clearBaseHover bool
	swallowOwner   widgets.Control
}

type popupRouteKind int

const (
	popupRouteNone popupRouteKind = iota
	popupRouteMenu
	popupRouteControlOverlay
)

func (d *Desktop) topPopupRouteKind() popupRouteKind {
	if host := d.topPopupHost(); host != nil {
		switch host.(type) {
		case *popupMenuState:
			return popupRouteMenu
		case *controlOverlayState:
			return popupRouteControlOverlay
		}
	}
	if d.menuMode {
		return popupRouteMenu
	}
	return popupRouteNone
}

func (d *Desktop) handlePopupMouseMove(point geom.Point) popupRouteResult {
	switch d.topPopupRouteKind() {
	case popupRouteControlOverlay:
		if d.handleControlOverlayMouseMove(point) {
			return popupRouteResult{handled: true, clearBaseHover: true}
		}
		if d.menuMode && d.handleMenuMouseMove(point) {
			return popupRouteResult{handled: true}
		}
	case popupRouteMenu:
		if d.handleMenuMouseMove(point) {
			return popupRouteResult{handled: true}
		}
		if d.handleControlOverlayMouseMove(point) {
			return popupRouteResult{handled: true, clearBaseHover: true}
		}
	default:
		if d.handleControlOverlayMouseMove(point) {
			return popupRouteResult{handled: true, clearBaseHover: true}
		}
	}
	return popupRouteResult{}
}

func (d *Desktop) handlePopupMouseDown(e event.MouseButtonEvent) popupMouseDownResult {
	switch d.topPopupRouteKind() {
	case popupRouteControlOverlay:
		top := d.topControlOverlay()
		if closed, handled := d.handleControlOverlayMouseDown(e); handled {
			return popupMouseDownResult{handled: true, clearBaseHover: true}
		} else if closed != nil {
			return popupMouseDownResult{swallowOwner: closed.owner}
		}
		if top != nil && d.pointWithinPopupHostOrOwner(top, e.Position) {
			return popupMouseDownResult{}
		}
		if d.menuMode && d.handleMenuMouseDown(e.Position) {
			return popupMouseDownResult{handled: true}
		}
	case popupRouteMenu:
		if d.handleMenuMouseDown(e.Position) {
			return popupMouseDownResult{handled: true}
		}
		if closed, handled := d.handleControlOverlayMouseDown(e); handled {
			return popupMouseDownResult{handled: true, clearBaseHover: true}
		} else if closed != nil {
			return popupMouseDownResult{swallowOwner: closed.owner}
		}
	default:
		if closed, handled := d.handleControlOverlayMouseDown(e); handled {
			return popupMouseDownResult{handled: true, clearBaseHover: true}
		} else if closed != nil {
			return popupMouseDownResult{swallowOwner: closed.owner}
		}
	}
	return popupMouseDownResult{}
}

func (d *Desktop) handlePopupMouseUp(e event.MouseButtonEvent) popupRouteResult {
	switch d.topPopupRouteKind() {
	case popupRouteControlOverlay:
		if d.handleControlOverlayMouseUp(e) {
			return popupRouteResult{handled: true, clearBaseHover: true}
		}
		if d.menuMode && d.handleMenuMouseUp(e) {
			return popupRouteResult{handled: true}
		}
	case popupRouteMenu:
		if d.handleMenuMouseUp(e) {
			return popupRouteResult{handled: true}
		}
		if d.handleControlOverlayMouseUp(e) {
			return popupRouteResult{handled: true, clearBaseHover: true}
		}
	default:
		if d.handleControlOverlayMouseUp(e) {
			return popupRouteResult{handled: true, clearBaseHover: true}
		}
	}
	return popupRouteResult{}
}

func (d *Desktop) handlePopupMouseWheel(e event.MouseWheel) bool {
	if d.handleControlOverlayMouseWheel(e) {
		return true
	}
	return false
}

func (d *Desktop) handlePopupKeyDown(active *Window, e event.KeyEvent) bool {
	if d.handleMenuKeyDown(active, e) {
		return true
	}
	if d.handleControlOverlayKeyDown(e) {
		return true
	}
	return false
}
