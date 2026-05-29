package widgets

import (
	"classicui/event"
	"classicui/geom"
	"strings"
	"time"

	"classicui/widget"
)

type Edit struct {
	widget.BaseWidget
	text        string
	readOnly    bool
	focused     bool
	selecting   bool
	caret       int
	anchor      int
	scrollX     int
	caretOn     bool
	lastBlink   time.Time
	composition string
	compStart   int
	compLength  int
	onChange    func(string)
}

func NewEdit(id string, bounds geom.Rect) *Edit {
	return &Edit{
		BaseWidget: widget.NewBase(id, bounds),
		caretOn:    true,
	}
}

func (e *Edit) SetText(text string) {
	e.text = text
	runes := []rune(text)
	e.caret = len(runes)
	e.anchor = e.caret
	e.composition = ""
}

func (e *Edit) Text() string {
	return e.text
}

func (e *Edit) SetReadOnly(readOnly bool) {
	e.readOnly = readOnly
}

func (e *Edit) OnChange(fn func(string)) {
	e.onChange = fn
}

func (e *Edit) Paint(ctx PaintContext) error {
	if !e.Visible() {
		return nil
	}

	rect := ctx.BoundsFor(e)
	ctx.Canvas.FillRect(rect, ctx.Theme.Colors.Face)
	ctx.Canvas.DrawDoubleBevel(rect, ctx.Theme.Colors.Shadow, ctx.Theme.Colors.Lightest, ctx.Theme.Colors.DarkShadow, ctx.Theme.Colors.Light)

	content := e.contentRect(rect)
	ctx.Canvas.FillRect(content, ctx.Theme.Colors.Window)
	ctx.Canvas.PushClip(content)
	defer ctx.Canvas.PopClip()

	textY := content.Y + maxInt((content.H-ctx.Text.LineHeight())/2, 0)
	textX := content.X + 2 - e.scrollX
	selectionStart, selectionEnd, hasSelection := e.selectionRange()
	runes := []rune(e.text)

	if hasSelection {
		before := string(runes[:selectionStart])
		selected := string(runes[selectionStart:selectionEnd])
		after := string(runes[selectionEnd:])

		beforeWidth := ctx.Text.MeasureString(before).W
		selectedWidth := ctx.Text.MeasureString(selected).W
		selectionRect := geom.Rect{
			X: textX + beforeWidth,
			Y: content.Y + 1,
			W: selectedWidth,
			H: content.H - 2,
		}
		ctx.Canvas.FillRect(selectionRect, ctx.Theme.Colors.Highlight)

		if before != "" {
			if err := ctx.Text.DrawString(ctx.Canvas, geom.Point{X: textX, Y: textY}, before, ctx.Theme.Colors.WindowText); err != nil {
				return err
			}
		}
		if selected != "" {
			if err := ctx.Text.DrawString(ctx.Canvas, geom.Point{X: textX + beforeWidth, Y: textY}, selected, ctx.Theme.Colors.HighlightText); err != nil {
				return err
			}
		}
		if after != "" {
			if err := ctx.Text.DrawString(ctx.Canvas, geom.Point{X: textX + beforeWidth + selectedWidth, Y: textY}, after, ctx.Theme.Colors.WindowText); err != nil {
				return err
			}
		}
	} else if e.text != "" {
		if err := ctx.Text.DrawString(ctx.Canvas, geom.Point{X: textX, Y: textY}, e.text, ctx.Theme.Colors.WindowText); err != nil {
			return err
		}
	}

	if e.composition != "" {
		compX := textX + ctx.Text.MeasureString(string(runes[:e.caret])).W
		if err := ctx.Text.DrawString(ctx.Canvas, geom.Point{X: compX, Y: textY}, e.composition, ctx.Theme.Colors.GrayText); err != nil {
			return err
		}
		compWidth := ctx.Text.MeasureString(e.composition).W
		underlineY := textY + ctx.Text.LineHeight() - 1
		for x := compX; x < compX+compWidth; x += 2 {
			ctx.Canvas.DrawPixel(x, underlineY, ctx.Theme.Colors.DarkShadow)
		}
	}

	if e.focused && e.caretOn {
		caretRect := e.localCaretRect(func(text string) geom.Size {
			return ctx.Text.MeasureString(text)
		}, ctx.Text.LineHeight())
		screenCaret := caretRect.Move(rect.X, rect.Y)
		ctx.Canvas.DrawVLine(screenCaret.X, screenCaret.Y, maxInt(screenCaret.H, 1), ctx.Theme.Colors.DarkShadow)
	}
	return nil
}

func (e *Edit) MouseEnter(EventContext) {}
func (e *Edit) MouseLeave(EventContext) {}

func (e *Edit) MouseMove(ctx EventContext, local geom.Point) {
	if !e.selecting {
		return
	}
	next := e.hitIndex(ctx, local.X)
	if next == e.caret {
		return
	}
	e.caret = next
	e.ensureCaretVisible(ctx)
	e.resetCaretBlink()
	ctx.Invalidate(e)
}

func (e *Edit) MouseDown(ctx EventContext, ev event.MouseButtonEvent, local geom.Point) {
	if ev.Button != event.MouseButtonLeft || !LocalContains(e, local) {
		return
	}
	e.selecting = true
	e.caret = e.hitIndex(ctx, local.X)
	e.anchor = e.caret
	e.ensureCaretVisible(ctx)
	e.resetCaretBlink()
	ctx.Capture(e)
	ctx.Invalidate(e)
}

func (e *Edit) MouseUp(ctx EventContext, ev event.MouseButtonEvent, local geom.Point) {
	if ev.Button != event.MouseButtonLeft {
		return
	}
	if e.selecting {
		e.caret = e.hitIndex(ctx, local.X)
		e.ensureCaretVisible(ctx)
	}
	e.selecting = false
	ctx.ReleaseCapture(e)
	e.resetCaretBlink()
	ctx.Invalidate(e)
}

func (e *Edit) KeyDown(ctx EventContext, ev event.KeyEvent) bool {
	if !e.Enabled() {
		return false
	}

	if ev.Modifiers&event.ModCtrl != 0 {
		switch ev.Key {
		case event.KeyA:
			e.anchor = 0
			e.caret = len([]rune(e.text))
			e.ensureCaretVisible(ctx)
			e.resetCaretBlink()
			ctx.Invalidate(e)
			return true
		case event.KeyC:
			if text, ok := e.selectedText(); ok {
				ctx.SetClipboardText(text)
			}
			return true
		case event.KeyX:
			if e.readOnly {
				return true
			}
			if text, ok := e.selectedText(); ok {
				ctx.SetClipboardText(text)
				if e.deleteSelection() {
					e.composition = ""
					e.ensureCaretVisible(ctx)
					e.notifyChange()
					e.resetCaretBlink()
					ctx.Invalidate(e)
				}
			}
			return true
		case event.KeyV:
			if e.readOnly {
				return true
			}
			clip := ctx.ClipboardText()
			if clip != "" {
				e.insertText(clip)
				e.composition = ""
				e.ensureCaretVisible(ctx)
				e.notifyChange()
				e.resetCaretBlink()
				ctx.Invalidate(e)
			}
			return true
		}
	}

	switch ev.Key {
	case event.KeyLeft:
		e.moveCaret(e.caret-1, ev.Modifiers&event.ModShift != 0)
	case event.KeyRight:
		e.moveCaret(e.caret+1, ev.Modifiers&event.ModShift != 0)
	case event.KeyHome:
		e.moveCaret(0, ev.Modifiers&event.ModShift != 0)
	case event.KeyEnd:
		e.moveCaret(len([]rune(e.text)), ev.Modifiers&event.ModShift != 0)
	case event.KeyBackspace:
		if e.readOnly {
			return true
		}
		if e.deleteSelection() {
			e.notifyChange()
		} else if e.caret > 0 {
			e.deleteRange(e.caret-1, e.caret)
			e.notifyChange()
		} else {
			return true
		}
	case event.KeyDelete:
		if e.readOnly {
			return true
		}
		runes := []rune(e.text)
		if e.deleteSelection() {
			e.notifyChange()
		} else if e.caret < len(runes) {
			e.deleteRange(e.caret, e.caret+1)
			e.notifyChange()
		} else {
			return true
		}
	default:
		return false
	}

	e.composition = ""
	e.ensureCaretVisible(ctx)
	e.resetCaretBlink()
	ctx.Invalidate(e)
	return true
}

func (e *Edit) CanFocus() bool {
	return e.Visible() && e.Enabled()
}

func (e *Edit) SetFocused(focused bool) {
	e.focused = focused
}

func (e *Edit) Focused() bool {
	return e.focused
}

func (e *Edit) FocusGained(ctx EventContext) {
	e.ensureCaretVisible(ctx)
	e.resetCaretBlink()
}

func (e *Edit) FocusLost(EventContext) {
	e.selecting = false
	e.caretOn = false
	e.composition = ""
}

func (e *Edit) Tick(ctx EventContext, now time.Time) bool {
	if !e.focused {
		return false
	}
	if now.Sub(e.lastBlink) < 500*time.Millisecond {
		return false
	}
	e.lastBlink = now
	e.caretOn = !e.caretOn
	ctx.Invalidate(e)
	return true
}

func (e *Edit) TextInput(ctx EventContext, ev event.TextInput) bool {
	if e.readOnly || ev.Text == "" {
		return true
	}
	e.insertText(ev.Text)
	e.composition = ""
	e.compStart = 0
	e.compLength = 0
	e.ensureCaretVisible(ctx)
	e.notifyChange()
	e.resetCaretBlink()
	ctx.Invalidate(e)
	return true
}

func (e *Edit) TextEditing(ctx EventContext, ev event.TextEditing) bool {
	e.composition = ev.Text
	e.compStart = ev.Start
	e.compLength = ev.Length
	e.ensureCaretVisible(ctx)
	e.resetCaretBlink()
	ctx.Invalidate(e)
	return true
}

func (e *Edit) TextInputRect(ctx EventContext) geom.Rect {
	return e.localCaretRect(ctx.MeasureText, ctx.LineHeight())
}

func (e *Edit) localCaretRect(measure func(string) geom.Size, lineHeight int) geom.Rect {
	content := e.contentRect(LocalRect(e))
	caretX := content.X + 2 + measure(string([]rune(e.text)[:e.caret])).W - e.scrollX
	return geom.Rect{
		X: caretX,
		Y: content.Y + maxInt((content.H-lineHeight)/2, 0),
		W: 1,
		H: maxInt(lineHeight, 1),
	}
}

func (e *Edit) contentRect(rect geom.Rect) geom.Rect {
	return geom.Rect{
		X: rect.X + 2,
		Y: rect.Y + 2,
		W: maxInt(rect.W-4, 0),
		H: maxInt(rect.H-4, 0),
	}
}

func (e *Edit) hitIndex(ctx EventContext, localX int) int {
	content := e.contentRect(LocalRect(e))
	target := localX - content.X - 2 + e.scrollX
	if target <= 0 {
		return 0
	}
	runes := []rune(e.text)
	for i := 1; i <= len(runes); i++ {
		prevWidth := ctx.MeasureText(string(runes[:i-1])).W
		width := ctx.MeasureText(string(runes[:i])).W
		if target < prevWidth+(width-prevWidth)/2 {
			return i - 1
		}
		if target < width {
			return i
		}
	}
	return len(runes)
}

func (e *Edit) moveCaret(next int, keepSelection bool) {
	limit := len([]rune(e.text))
	e.caret = clampRange(next, 0, limit)
	if !keepSelection {
		e.anchor = e.caret
	}
}

func (e *Edit) ensureCaretVisible(ctx EventContext) {
	content := e.contentRect(LocalRect(e))
	usable := maxInt(content.W-4, 1)
	prefixWidth := ctx.MeasureText(string([]rune(e.text)[:e.caret])).W
	if prefixWidth-e.scrollX < 0 {
		e.scrollX = prefixWidth
	}
	if prefixWidth-e.scrollX > usable {
		e.scrollX = prefixWidth - usable
	}
	total := ctx.MeasureText(e.text).W
	maxScroll := maxInt(total-usable, 0)
	e.scrollX = clampRange(e.scrollX, 0, maxScroll)
}

func (e *Edit) insertText(text string) {
	if e.deleteSelection() {
		e.composition = ""
	}
	runes := []rune(e.text)
	insert := []rune(text)
	left := append([]rune(nil), runes[:e.caret]...)
	right := append([]rune(nil), runes[e.caret:]...)
	combined := append(left, insert...)
	combined = append(combined, right...)
	e.text = string(combined)
	e.caret += len(insert)
	e.anchor = e.caret
}

func (e *Edit) deleteSelection() bool {
	start, end, ok := e.selectionRange()
	if !ok {
		return false
	}
	e.deleteRange(start, end)
	return true
}

func (e *Edit) deleteRange(start, end int) {
	runes := []rune(e.text)
	start = clampRange(start, 0, len(runes))
	end = clampRange(end, start, len(runes))
	combined := append([]rune(nil), runes[:start]...)
	combined = append(combined, runes[end:]...)
	e.text = string(combined)
	e.caret = start
	e.anchor = start
}

func (e *Edit) selectionRange() (int, int, bool) {
	if e.anchor == e.caret {
		return 0, 0, false
	}
	if e.anchor < e.caret {
		return e.anchor, e.caret, true
	}
	return e.caret, e.anchor, true
}

func (e *Edit) selectedText() (string, bool) {
	start, end, ok := e.selectionRange()
	if !ok {
		return "", false
	}
	runes := []rune(e.text)
	return string(runes[start:end]), true
}

func (e *Edit) notifyChange() {
	if e.onChange != nil {
		e.onChange(e.text)
	}
}

func (e *Edit) resetCaretBlink() {
	e.caretOn = true
	e.lastBlink = time.Now()
}

func clampRange(value, minValue, maxValue int) int {
	if value < minValue {
		return minValue
	}
	if value > maxValue {
		return maxValue
	}
	return value
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (e *Edit) String() string {
	return strings.TrimSpace(e.text)
}
