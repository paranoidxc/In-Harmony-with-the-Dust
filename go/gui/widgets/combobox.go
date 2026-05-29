package widgets

import (
	"classicui/event"
	"classicui/geom"
	"classicui/widget"
	"strings"
	"time"
)

type ComboBox struct {
	widget.BaseWidget
	items          []string
	selected       int
	text           string
	tooltip        string
	editable       bool
	edit           *Edit
	focused        bool
	hot            bool
	pressed        bool
	dropped        bool
	popup          *ListBox
	onChange       func(int, string)
	onCommit       func(int, string)
	syncingEdit    bool
	committedIndex int
	committedText  string
}

func NewComboBox(id string, bounds geom.Rect) *ComboBox {
	combo := &ComboBox{
		BaseWidget:     widget.NewBase(id, bounds),
		selected:       -1,
		committedIndex: -1,
	}
	combo.setText("", false)
	return combo
}

func (c *ComboBox) SetItems(items []string) {
	c.items = append([]string(nil), items...)
	if len(c.items) == 0 {
		c.selected = -1
		c.setText("", false)
		c.recordCommitted(-1, "")
		return
	}
	if c.selected < 0 || c.selected >= len(c.items) {
		c.selected = 0
	}
	if c.editable {
		if c.selected >= 0 {
			c.setText(c.items[c.selected], false)
		}
		c.recordCommitted(c.selected, c.Text())
		return
	}
	c.text = c.items[c.selected]
	c.recordCommitted(c.selected, c.text)
}

func (c *ComboBox) Items() []string {
	return append([]string(nil), c.items...)
}

func (c *ComboBox) SelectedIndex() int {
	return c.selected
}

func (c *ComboBox) Text() string {
	if c.editable {
		return c.text
	}
	if c.selected < 0 || c.selected >= len(c.items) {
		return ""
	}
	return c.items[c.selected]
}

func (c *ComboBox) SetSelectedIndex(index int) bool {
	changed := c.setSelectedIndex(index, true)
	if changed {
		c.recordCommitted(c.selected, c.Text())
	}
	return changed
}

func (c *ComboBox) SetSelectedIndexSilent(index int) bool {
	changed := c.setSelectedIndex(index, false)
	if changed {
		c.recordCommitted(c.selected, c.Text())
	}
	return changed
}

func (c *ComboBox) setSelectedIndex(index int, notify bool) bool {
	if len(c.items) == 0 || index < -1 || index >= len(c.items) {
		return false
	}
	nextText := c.text
	if index >= 0 {
		nextText = c.items[index]
	}
	if !c.editable && index < 0 {
		return false
	}
	if c.selected == index && (!c.editable || c.text == nextText) {
		return false
	}
	c.selected = index
	if c.editable {
		c.setText(nextText, false)
	} else if index >= 0 {
		c.text = nextText
	}
	if notify && c.onChange != nil {
		c.onChange(index, c.Text())
	}
	return true
}

func (c *ComboBox) OnChange(fn func(int, string)) {
	c.onChange = fn
}

func (c *ComboBox) OnCommit(fn func(int, string)) {
	c.onCommit = fn
}

func (c *ComboBox) SetTooltip(text string) {
	c.tooltip = text
}

func (c *ComboBox) SetEditable(editable bool) {
	if c.editable == editable {
		return
	}
	c.editable = editable
	if editable {
		c.ensureEdit()
		if c.selected >= 0 && c.selected < len(c.items) {
			c.setText(c.items[c.selected], false)
		}
		return
	}
	if c.edit != nil {
		c.edit.SetVisible(false)
	}
	if c.selected >= 0 && c.selected < len(c.items) {
		c.text = c.items[c.selected]
	}
}

func (c *ComboBox) Editable() bool {
	return c.editable
}

func (c *ComboBox) SetText(text string) {
	c.setText(text, true)
}

func (c *ComboBox) SetBounds(bounds geom.Rect) {
	c.BaseWidget.SetBounds(bounds)
	c.syncEditBounds()
}

func (c *ComboBox) Paint(ctx PaintContext) error {
	if !c.Visible() {
		return nil
	}
	rect := ctx.BoundsFor(c)
	ctx.Canvas.FillRect(rect, ctx.Theme.Colors.Face)
	ctx.Canvas.DrawDoubleBevel(rect, ctx.Theme.Colors.Shadow, ctx.Theme.Colors.Lightest, ctx.Theme.Colors.DarkShadow, ctx.Theme.Colors.Light)

	content := c.contentRect(rect, ctx.Theme.Metrics.ScrollbarSize)
	button := c.buttonRect(rect, ctx.Theme.Metrics.ScrollbarSize)
	ctx.Canvas.FillRect(content, ctx.Theme.Colors.Window)
	fill := ctx.Theme.Colors.Face
	if c.hot && !c.pressed {
		fill = blend(ctx.Theme.Colors.Face, ctx.Theme.Colors.Lightest)
	}
	ctx.Canvas.FillRect(button, fill)
	if c.pressed || c.dropped {
		ctx.Canvas.DrawDoubleBevel(button, ctx.Theme.Colors.Shadow, ctx.Theme.Colors.Lightest, ctx.Theme.Colors.DarkShadow, ctx.Theme.Colors.Light)
	} else {
		ctx.Canvas.DrawDoubleBevel(button, ctx.Theme.Colors.Lightest, ctx.Theme.Colors.DarkShadow, ctx.Theme.Colors.Light, ctx.Theme.Colors.Shadow)
	}
	c.paintArrow(ctx, button, c.pressed || c.dropped)

	if c.editable {
		c.ensureEdit()
		c.syncEditBounds()
		if c.edit != nil {
			childCtx := ctx.Child(c)
			if err := c.edit.Paint(childCtx); err != nil {
				return err
			}
			c.paintCandidateTail(ctx, rect)
		}
	} else if ctx.Text != nil && c.Text() != "" {
		textColor := ctx.Theme.Colors.WindowText
		if !c.Enabled() {
			textColor = ctx.Theme.Colors.GrayText
		}
		textSize := ctx.Text.MeasureString(c.Text())
		textY := content.Y + maxInt((content.H-textSize.H)/2, 0)
		textRect := content.Inset(3)
		ctx.Canvas.PushClip(textRect)
		err := ctx.Text.DrawString(ctx.Canvas, geom.Point{X: textRect.X, Y: textY}, c.Text(), textColor)
		ctx.Canvas.PopClip()
		if err != nil {
			return err
		}
	}

	if c.focused && !c.editable {
		focusRect := content.Inset(2)
		if focusRect.W > 2 && focusRect.H > 2 {
			ctx.Canvas.DrawFocusRect(focusRect, ctx.Theme.Colors.DarkShadow)
		}
	}
	return nil
}

func (c *ComboBox) MouseEnter(ctx EventContext) {
	if !c.Enabled() {
		return
	}
	if c.hot {
		return
	}
	c.hot = true
	ctx.Invalidate(c)
}

func (c *ComboBox) MouseLeave(ctx EventContext) {
	if c.dropped || !c.hot {
		return
	}
	c.hot = false
	ctx.Invalidate(c)
}

func (c *ComboBox) MouseMove(EventContext, geom.Point) {}

func (c *ComboBox) MouseDown(ctx EventContext, e event.MouseButtonEvent, local geom.Point) {
	if !c.Enabled() || e.Button != event.MouseButtonLeft || !LocalContains(c, local) {
		return
	}
	ctx.SetFocus(c)
	if c.editable && c.edit != nil {
		button := c.buttonRect(LocalRect(c), 16)
		if !button.Contains(local) {
			editLocal := geom.Point{X: local.X - c.edit.Bounds().X, Y: local.Y - c.edit.Bounds().Y}
			c.edit.MouseDown(ctx, e, editLocal)
			return
		}
	}
	c.pressed = true
	c.hot = true
	c.togglePopup(ctx)
	ctx.Invalidate(c)
}

func (c *ComboBox) MouseUp(ctx EventContext, e event.MouseButtonEvent, local geom.Point) {
	if e.Button != event.MouseButtonLeft {
		return
	}
	if c.editable && c.edit != nil {
		editLocal := geom.Point{X: local.X - c.edit.Bounds().X, Y: local.Y - c.edit.Bounds().Y}
		c.edit.MouseUp(ctx, e, editLocal)
	}
	if !c.dropped && c.pressed {
		c.pressed = false
		ctx.Invalidate(c)
	}
}

func (c *ComboBox) KeyDown(ctx EventContext, e event.KeyEvent) bool {
	if !c.Enabled() || len(c.items) == 0 {
		if c.editable && c.edit != nil {
			return c.edit.KeyDown(ctx, e)
		}
		return false
	}

	if c.dropped && c.popup != nil {
		switch e.Key {
		case event.KeyEscape:
			c.restoreCommitted(ctx)
			c.closePopup(ctx)
			return true
		case event.KeyEnter, event.KeySpace:
			index := c.popup.SelectedIndex()
			if index >= 0 {
				c.commitSelection(index)
				ctx.Invalidate(c)
			} else {
				c.commitCurrentText()
			}
			c.closePopup(ctx)
			return true
		case event.KeyTab:
			if index := c.popup.SelectedIndex(); index >= 0 {
				c.commitSelection(index)
			} else {
				c.commitCurrentText()
			}
			c.closePopup(ctx)
			return false
		case event.KeyUp:
			c.movePopupSelection(-1)
			ctx.Invalidate(c)
			return true
		case event.KeyDown:
			c.movePopupSelection(1)
			ctx.Invalidate(c)
			return true
		case event.KeyHome:
			c.setPopupSelection(0)
			ctx.Invalidate(c)
			return true
		case event.KeyEnd:
			c.setPopupSelection(len(c.items) - 1)
			ctx.Invalidate(c)
			return true
		}
	}

	if c.editable && c.edit != nil {
		switch e.Key {
		case event.KeyLeft, event.KeyRight, event.KeyBackspace, event.KeyDelete:
			handled := c.edit.KeyDown(ctx, e)
			if handled {
				c.syncCandidateSelection(ctx)
			}
			return handled
		case event.KeyHome, event.KeyEnd:
			if e.Modifiers&event.ModCtrl == 0 {
				handled := c.edit.KeyDown(ctx, e)
				if handled {
					c.syncCandidateSelection(ctx)
				}
				return handled
			}
		}
	}

	switch e.Key {
	case event.KeyUp:
		if c.editable {
			if !c.dropped {
				c.openPopup(ctx)
				if c.popup != nil && c.popup.SelectedIndex() < 0 && len(c.items) > 0 {
					c.popup.SetSelectedIndexSilent(len(c.items) - 1)
					ctx.Invalidate(c.popup)
				}
				ctx.Invalidate(c)
				return true
			}
			c.movePopupSelection(-1)
			ctx.Invalidate(c)
			return true
		}
		return c.stepSelection(ctx, -1)
	case event.KeyDown:
		if c.editable {
			wasDropped := c.dropped
			c.openPopup(ctx)
			if c.popup != nil && wasDropped {
				c.movePopupSelection(1)
			}
			ctx.Invalidate(c)
			return true
		}
		return c.stepSelection(ctx, 1)
	case event.KeyHome:
		return c.selectIndex(ctx, 0)
	case event.KeyEnd:
		return c.selectIndex(ctx, len(c.items)-1)
	case event.KeyEnter:
		if c.editable {
			c.commitCurrentText()
			return true
		}
		c.openPopup(ctx)
		return true
	case event.KeyEscape:
		if c.editable {
			c.restoreCommitted(ctx)
			return true
		}
		return false
	case event.KeySpace:
		if c.editable {
			return false
		}
		c.openPopup(ctx)
		return true
	case event.KeyTab:
		if c.editable {
			if index := c.currentCandidateIndex(); index >= 0 {
				c.commitSelection(index)
			} else {
				c.commitCurrentText()
			}
			return false
		}
		return false
	default:
		if c.editable && c.edit != nil {
			handled := c.edit.KeyDown(ctx, e)
			if handled {
				c.syncCandidateSelection(ctx)
			}
			return handled
		}
		return false
	}
}

func (c *ComboBox) CanFocus() bool {
	return c.Visible() && c.Enabled()
}

func (c *ComboBox) SetFocused(focused bool) {
	c.focused = focused
	if c.editable && c.edit != nil {
		c.edit.SetFocused(focused)
	}
}

func (c *ComboBox) Focused() bool {
	return c.focused
}

func (c *ComboBox) TooltipAt(geom.Point, func(string) geom.Size) TooltipInfo {
	return TooltipInfo{
		Text:   c.tooltip,
		Anchor: LocalRect(c),
	}
}

func (c *ComboBox) FocusLost(ctx EventContext) {
	if c.dropped {
		c.closePopup(ctx)
	}
	if c.editable && c.edit != nil {
		c.edit.FocusLost(ctx)
	}
}

func (c *ComboBox) FocusGained(ctx EventContext) {
	if c.editable && c.edit != nil {
		c.edit.FocusGained(ctx)
	}
}

func (c *ComboBox) Tick(ctx EventContext, now time.Time) bool {
	if !c.editable || c.edit == nil {
		return false
	}
	return c.edit.Tick(ctx, now)
}

func (c *ComboBox) TextInput(ctx EventContext, ev event.TextInput) bool {
	if !c.editable || c.edit == nil {
		return false
	}
	handled := c.edit.TextInput(ctx, ev)
	if handled {
		c.syncCandidateSelection(ctx)
	}
	return handled
}

func (c *ComboBox) TextEditing(ctx EventContext, ev event.TextEditing) bool {
	if !c.editable || c.edit == nil {
		return false
	}
	return c.edit.TextEditing(ctx, ev)
}

func (c *ComboBox) TextInputRect(ctx EventContext) geom.Rect {
	if !c.editable || c.edit == nil {
		return geom.Rect{}
	}
	rect := c.edit.TextInputRect(ctx)
	return rect.Move(c.edit.Bounds().X, c.edit.Bounds().Y)
}

func (c *ComboBox) contentRect(rect geom.Rect, metric int) geom.Rect {
	button := c.buttonRect(rect, metric)
	return geom.Rect{
		X: rect.X + 2,
		Y: rect.Y + 2,
		W: maxInt(button.X-rect.X-4, 0),
		H: maxInt(rect.H-4, 0),
	}
}

func (c *ComboBox) buttonRect(rect geom.Rect, metric int) geom.Rect {
	width := maxInt(metric, 16)
	return geom.Rect{
		X: rect.Right() - width - 2,
		Y: rect.Y + 2,
		W: width,
		H: maxInt(rect.H-4, 0),
	}
}

func (c *ComboBox) paintArrow(ctx PaintContext, rect geom.Rect, pressed bool) {
	cx := rect.X + rect.W/2
	cy := rect.Y + rect.H/2
	offset := 0
	if pressed {
		offset = 1
	}
	for row := 0; row < 4; row++ {
		for x := -row; x <= row; x++ {
			ctx.Canvas.DrawPixel(cx+x+offset, cy-1+row+offset, ctx.Theme.Colors.DarkShadow)
		}
	}
}

func (c *ComboBox) togglePopup(ctx EventContext) {
	if c.dropped {
		c.closePopup(ctx)
		return
	}
	c.openPopup(ctx)
}

func (c *ComboBox) openPopup(ctx EventContext) {
	presenter, ok := ctx.(OverlayContext)
	if !ok || presenter.OverlayVisible(c) || len(c.items) == 0 {
		return
	}

	lineHeight := ctx.LineHeight()
	if lineHeight <= 0 {
		lineHeight = 14
	}
	visibleRows := minInt(maxInt(len(c.items), 1), 8)
	popupHeight := visibleRows*maxInt(lineHeight+2, 16) + 4
	popupWidth := maxInt(c.Bounds().W, 96)
	popup := NewListBox(c.ID()+".popup", geom.Rect{X: 0, Y: 0, W: popupWidth, H: popupHeight})
	popup.SetItems(c.items)
	if match := c.currentCandidateIndex(); match >= 0 {
		popup.SetSelectedIndexSilent(match)
	} else if c.selected >= 0 {
		popup.SetSelectedIndexSilent(c.selected)
	}
	popup.OnActivate(func(index int, _ string) {
		c.commitSelection(index)
		c.closePopup(ctx)
	})

	request := OverlayRequest{
		Owner:          c,
		Content:        popup,
		Anchor:         LocalRect(c),
		Placement:      OverlayBelowStart,
		CloseOnOutside: true,
		OnClose: func() {
			c.dropped = false
			c.pressed = false
			c.hot = false
			c.popup = nil
			ctx.Invalidate(c)
		},
	}
	if presenter.ShowOverlay(request) {
		c.popup = popup
		c.dropped = true
		c.pressed = true
	}
}

func (c *ComboBox) closePopup(ctx EventContext) {
	presenter, ok := ctx.(OverlayContext)
	if !ok {
		c.dropped = false
		c.pressed = false
		c.popup = nil
		return
	}
	if !presenter.HideOverlay(c) {
		c.dropped = false
		c.pressed = false
		c.popup = nil
	}
}

func (c *ComboBox) stepSelection(ctx EventContext, delta int) bool {
	if len(c.items) == 0 {
		return false
	}
	next := c.selected
	if next < 0 {
		next = 0
	} else {
		next = clampInt(next+delta, 0, len(c.items)-1)
	}
	if c.setSelectedIndex(next, true) {
		if c.onCommit != nil {
			c.onCommit(next, c.Text())
		}
		ctx.Invalidate(c)
	}
	return true
}

func (c *ComboBox) selectIndex(ctx EventContext, index int) bool {
	if c.setSelectedIndex(index, true) {
		if c.onCommit != nil {
			c.onCommit(index, c.Text())
		}
		ctx.Invalidate(c)
	}
	return true
}

func (c *ComboBox) movePopupSelection(delta int) {
	if c.popup == nil || len(c.items) == 0 {
		return
	}
	next := c.popup.SelectedIndex()
	if next < 0 {
		next = c.selected
	}
	if next < 0 {
		next = 0
	}
	next = clampInt(next+delta, 0, len(c.items)-1)
	c.popup.SetSelectedIndexSilent(next)
}

func (c *ComboBox) setPopupSelection(index int) {
	if c.popup == nil || len(c.items) == 0 {
		return
	}
	c.popup.SetSelectedIndexSilent(clampInt(index, 0, len(c.items)-1))
}

func (c *ComboBox) ensureEdit() {
	if c.edit != nil {
		c.edit.SetVisible(true)
		c.syncEditBounds()
		return
	}
	edit := NewEdit(c.ID()+".edit", geom.Rect{})
	edit.SetFrameVisible(false)
	edit.SetParent(c)
	edit.OnChange(func(text string) {
		if c.syncingEdit {
			return
		}
		match := c.matchIndex(text)
		changed := c.selected != match || c.text != text
		c.text = text
		c.selected = match
		if changed && c.onChange != nil {
			c.onChange(match, text)
		}
	})
	c.edit = edit
	c.syncEditBounds()
	c.edit.SetText(c.text)
}

func (c *ComboBox) syncEditBounds() {
	if c.edit == nil {
		return
	}
	content := c.contentRect(LocalRect(c), 16)
	c.edit.SetBounds(content)
}

func (c *ComboBox) setText(text string, preserveSelection bool) {
	c.text = text
	if c.editable {
		c.ensureEdit()
		if c.edit != nil {
			c.syncingEdit = true
			c.edit.SetText(text)
			c.syncingEdit = false
		}
		if !preserveSelection {
			return
		}
	}
	if !preserveSelection {
		return
	}
	c.selected = c.matchIndex(text)
}

func (c *ComboBox) matchIndex(text string) int {
	for i, item := range c.items {
		if item == text {
			return i
		}
	}
	return -1
}

func (c *ComboBox) prefixMatchIndex(text string) int {
	if text == "" {
		return -1
	}
	target := strings.ToLower(text)
	for i, item := range c.items {
		if strings.HasPrefix(strings.ToLower(item), target) {
			return i
		}
	}
	return -1
}

func (c *ComboBox) currentCandidateIndex() int {
	if match := c.matchIndex(c.Text()); match >= 0 {
		return match
	}
	return c.prefixMatchIndex(c.Text())
}

func (c *ComboBox) syncCandidateSelection(ctx EventContext) {
	if c.popup == nil {
		return
	}
	index := c.currentCandidateIndex()
	c.popup.SetSelectedIndexSilent(index)
	ctx.Invalidate(c.popup)
}

func (c *ComboBox) paintCandidateTail(ctx PaintContext, rect geom.Rect) {
	if ctx.Text == nil || c.edit == nil || !c.focused || c.dropped {
		return
	}
	if _, _, hasSelection := c.edit.selectionRange(); hasSelection {
		return
	}
	if c.edit.composition != "" || c.edit.caret != len([]rune(c.text)) {
		return
	}
	match := c.prefixMatchIndex(c.text)
	if match < 0 {
		return
	}
	candidate := c.items[match]
	if candidate == c.text || !strings.HasPrefix(strings.ToLower(candidate), strings.ToLower(c.text)) {
		return
	}
	textRunes := []rune(c.text)
	candidateRunes := []rune(candidate)
	if len(candidateRunes) <= len(textRunes) {
		return
	}
	suffix := string(candidateRunes[len(textRunes):])
	prefixWidth := ctx.Text.MeasureString(c.text).W
	lineHeight := ctx.Text.LineHeight()
	editRect := c.edit.Bounds()
	textX := rect.X + editRect.X + 2 + prefixWidth - c.edit.scrollX
	textY := rect.Y + editRect.Y + maxInt((editRect.H-lineHeight)/2, 0)
	clip := geom.Rect{X: rect.X + editRect.X, Y: rect.Y + editRect.Y, W: editRect.W, H: editRect.H}
	ctx.Canvas.PushClip(clip)
	_ = ctx.Text.DrawString(ctx.Canvas, geom.Point{X: textX, Y: textY}, suffix, ctx.Theme.Colors.GrayText)
	ctx.Canvas.PopClip()
}

func (c *ComboBox) commitSelection(index int) {
	if index < 0 || index >= len(c.items) {
		return
	}
	c.setSelectedIndex(index, true)
	c.recordCommitted(c.selected, c.Text())
	if c.onCommit != nil {
		c.onCommit(index, c.Text())
	}
}

func (c *ComboBox) commitCurrentText() {
	c.recordCommitted(c.selected, c.Text())
	if c.onCommit != nil {
		c.onCommit(c.selected, c.Text())
	}
}

func (c *ComboBox) recordCommitted(index int, text string) {
	c.committedIndex = index
	c.committedText = text
}

func (c *ComboBox) restoreCommitted(ctx EventContext) {
	if c.editable {
		c.selected = c.committedIndex
		c.setText(c.committedText, false)
		if c.popup != nil {
			c.popup.SetSelectedIndexSilent(c.currentCandidateIndex())
		}
		ctx.Invalidate(c)
		return
	}
	if c.committedIndex >= 0 {
		c.setSelectedIndex(c.committedIndex, false)
		ctx.Invalidate(c)
	}
}
