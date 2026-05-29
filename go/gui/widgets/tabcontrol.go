package widgets

import (
	"classicui/event"
	"classicui/geom"
	"classicui/widget"
)

type TabPage struct {
	Title   string
	Content Control
	Enabled bool
}

func NewTabPage(title string, content Control) *TabPage {
	return &TabPage{
		Title:   title,
		Content: content,
		Enabled: true,
	}
}

type TabControl struct {
	widget.BaseWidget
	pages             []*TabPage
	selected          int
	focused           bool
	pressedTab        int
	hotTab            int
	cachedTabHeight   int
	onSelectionChange func(int, *TabPage)
}

func NewTabControl(id string, bounds geom.Rect, pages ...*TabPage) *TabControl {
	tab := &TabControl{
		BaseWidget:      widget.NewBase(id, bounds),
		selected:        -1,
		pressedTab:      -1,
		hotTab:          -1,
		cachedTabHeight: 22,
	}
	tab.SetPages(pages...)
	return tab
}

func (t *TabControl) SetBounds(bounds geom.Rect) {
	t.BaseWidget.SetBounds(bounds)
	t.syncPageBounds()
}

func (t *TabControl) SetPages(pages ...*TabPage) {
	id := t.ID()
	bounds := t.Bounds()
	visible := t.Visible()
	enabled := t.Enabled()
	parent := t.Parent()

	t.pages = nil
	t.selected = -1
	t.BaseWidget = widget.NewBase(id, bounds)
	t.SetVisible(visible)
	t.SetEnabled(enabled)
	t.SetParent(parent)
	for _, page := range pages {
		t.AddPage(page)
	}
	t.syncSelection()
	t.syncPageBounds()
}

func (t *TabControl) AddPage(page *TabPage) {
	if page == nil || page.Content == nil {
		return
	}
	page.Content.SetParent(t)
	t.BaseWidget.AppendChild(page.Content)
	t.pages = append(t.pages, page)
	if t.selected == -1 && page.Enabled {
		t.selected = len(t.pages) - 1
	}
	t.syncSelection()
	t.syncPageBounds()
}

func (t *TabControl) Pages() []*TabPage {
	return append([]*TabPage(nil), t.pages...)
}

func (t *TabControl) SelectedIndex() int {
	return t.selected
}

func (t *TabControl) SelectedPage() *TabPage {
	if t.selected < 0 || t.selected >= len(t.pages) {
		return nil
	}
	return t.pages[t.selected]
}

func (t *TabControl) SetSelected(index int) bool {
	if index < 0 || index >= len(t.pages) || !t.pages[index].Enabled {
		return false
	}
	if t.selected == index {
		return false
	}
	t.selected = index
	t.syncSelection()
	t.syncPageBounds()
	if t.onSelectionChange != nil {
		t.onSelectionChange(index, t.pages[index])
	}
	return true
}

func (t *TabControl) OnSelectionChange(fn func(int, *TabPage)) {
	t.onSelectionChange = fn
}

func (t *TabControl) Paint(ctx PaintContext) error {
	if !t.Visible() {
		return nil
	}

	rect := ctx.BoundsFor(t)
	lineHeight := 14
	measure := func(text string) geom.Size {
		if ctx.Text == nil {
			return geom.Size{}
		}
		return ctx.Text.MeasureString(text)
	}
	if ctx.Text != nil {
		lineHeight = ctx.Text.LineHeight()
	}
	t.cachedTabHeight = t.tabHeight(measure, lineHeight)
	t.syncSelection()
	t.syncPageBounds()

	pageRect := t.pageRect(rect)
	ctx.Canvas.FillRect(pageRect, ctx.Theme.Colors.Face)
	ctx.Canvas.DrawDoubleBevel(pageRect, ctx.Theme.Colors.Lightest, ctx.Theme.Colors.DarkShadow, ctx.Theme.Colors.Light, ctx.Theme.Colors.Shadow)

	tabs := t.layoutTabs(rect, measure, lineHeight)
	for i, tabRect := range tabs {
		page := t.pages[i]
		selected := i == t.selected
		fill := ctx.Theme.Colors.Face
		if !page.Enabled {
			fill = blendColor(ctx.Theme.Colors.Face, ctx.Theme.Colors.Light)
		}
		if selected {
			ctx.Canvas.FillRect(tabRect, fill)
			ctx.Canvas.DrawHLine(tabRect.X+2, tabRect.Bottom()-1, maxInt(tabRect.W-4, 0), fill)
			ctx.Canvas.DrawHLine(tabRect.X+2, tabRect.Bottom(), maxInt(tabRect.W-4, 0), fill)
			ctx.Canvas.DrawVLine(tabRect.X, tabRect.Y+2, tabRect.H-2, ctx.Theme.Colors.Lightest)
			ctx.Canvas.DrawVLine(tabRect.X+1, tabRect.Y+1, tabRect.H-1, ctx.Theme.Colors.Light)
			ctx.Canvas.DrawHLine(tabRect.X+2, tabRect.Y, maxInt(tabRect.W-3, 0), ctx.Theme.Colors.Lightest)
			ctx.Canvas.DrawHLine(tabRect.X+2, tabRect.Y+1, maxInt(tabRect.W-4, 0), ctx.Theme.Colors.Light)
			ctx.Canvas.DrawVLine(tabRect.Right()-1, tabRect.Y+2, tabRect.H-2, ctx.Theme.Colors.DarkShadow)
			continue
		}

		ctx.Canvas.FillRect(tabRect, fill)
		ctx.Canvas.DrawDoubleBevel(tabRect, ctx.Theme.Colors.Lightest, ctx.Theme.Colors.DarkShadow, ctx.Theme.Colors.Light, ctx.Theme.Colors.Shadow)
	}

	for i, tabRect := range tabs {
		page := t.pages[i]
		if ctx.Text != nil && page.Title != "" {
			textSize := ctx.Text.MeasureString(page.Title)
			textX := tabRect.X + (tabRect.W-textSize.W)/2
			textY := tabRect.Y + (tabRect.H-textSize.H)/2
			if i == t.selected {
				textY--
			}
			color := ctx.Theme.Colors.WindowText
			if !page.Enabled {
				color = ctx.Theme.Colors.GrayText
			}
			if err := ctx.Text.DrawString(ctx.Canvas, geom.Point{X: textX, Y: textY}, page.Title, color); err != nil {
				return err
			}
			if t.focused && i == t.selected {
				focusRect := tabRect.Inset(4)
				if focusRect.W > 2 && focusRect.H > 2 {
					ctx.Canvas.DrawFocusRect(focusRect, ctx.Theme.Colors.DarkShadow)
				}
			}
		}
	}

	selected := t.SelectedPage()
	if selected == nil || selected.Content == nil {
		return nil
	}

	childCtx := ctx.Child(t)
	ctx.Canvas.PushClip(pageRect.Inset(2))
	err := selected.Content.Paint(childCtx)
	ctx.Canvas.PopClip()
	return err
}

func (t *TabControl) MouseEnter(EventContext) {}

func (t *TabControl) MouseLeave(ctx EventContext) {
	if t.hotTab < 0 && t.pressedTab < 0 {
		return
	}
	t.hotTab = -1
	if t.pressedTab < 0 {
		ctx.Invalidate(t)
	}
}

func (t *TabControl) MouseMove(ctx EventContext, local geom.Point) {
	index := t.hitTab(local, ctx.MeasureText, ctx.LineHeight())
	if t.pressedTab >= 0 {
		if t.hotTab == index {
			return
		}
		t.hotTab = index
		ctx.Invalidate(t)
		return
	}
	if t.hotTab == index {
		return
	}
	t.hotTab = index
	ctx.Invalidate(t)
}

func (t *TabControl) MouseDown(ctx EventContext, e event.MouseButtonEvent, local geom.Point) {
	if !t.Enabled() || e.Button != event.MouseButtonLeft {
		return
	}
	index := t.hitTab(local, ctx.MeasureText, ctx.LineHeight())
	if index < 0 || !t.pages[index].Enabled {
		return
	}
	t.pressedTab = index
	t.hotTab = index
	ctx.SetFocus(t)
	ctx.Capture(t)
	ctx.Invalidate(t)
}

func (t *TabControl) MouseUp(ctx EventContext, e event.MouseButtonEvent, local geom.Point) {
	if e.Button != event.MouseButtonLeft || t.pressedTab < 0 {
		return
	}
	index := t.hitTab(local, ctx.MeasureText, ctx.LineHeight())
	pressed := t.pressedTab
	t.pressedTab = -1
	ctx.ReleaseCapture(t)
	if index >= 0 {
		t.hotTab = index
	} else {
		t.hotTab = -1
	}
	if index == pressed && index != t.selected {
		t.SetSelected(index)
	}
	ctx.Invalidate(t)
}

func (t *TabControl) KeyDown(ctx EventContext, e event.KeyEvent) bool {
	if !t.Enabled() || len(t.pages) == 0 {
		return false
	}

	next := t.selected
	switch e.Key {
	case event.KeyLeft, event.KeyUp:
		next = t.nextEnabled(t.selected, -1)
	case event.KeyRight, event.KeyDown:
		next = t.nextEnabled(t.selected, 1)
	case event.KeyHome:
		next = t.nextEnabled(-1, 1)
	case event.KeyEnd:
		next = t.nextEnabled(len(t.pages), -1)
	default:
		return false
	}
	if next < 0 {
		return true
	}
	if t.SetSelected(next) {
		ctx.Invalidate(t)
	}
	return true
}

func (t *TabControl) CanFocus() bool {
	return t.Visible() && t.Enabled() && len(t.pages) > 0
}

func (t *TabControl) SetFocused(focused bool) {
	t.focused = focused
}

func (t *TabControl) Focused() bool {
	return t.focused
}

func (t *TabControl) syncSelection() {
	if len(t.pages) == 0 {
		t.selected = -1
		return
	}
	if t.selected >= 0 && t.selected < len(t.pages) && t.pages[t.selected].Enabled {
		for i, page := range t.pages {
			page.Content.SetVisible(i == t.selected)
		}
		return
	}
	t.selected = t.nextEnabled(-1, 1)
	for i, page := range t.pages {
		page.Content.SetVisible(i == t.selected)
	}
}

func (t *TabControl) syncPageBounds() {
	content := t.pageContentRect(LocalRect(t))
	for _, page := range t.pages {
		if page == nil || page.Content == nil {
			continue
		}
		page.Content.SetBounds(content)
	}
}

func (t *TabControl) pageRect(rect geom.Rect) geom.Rect {
	top := rect.Y + t.cachedTabHeight - 1
	return geom.Rect{
		X: rect.X,
		Y: top,
		W: rect.W,
		H: maxInt(rect.H-(top-rect.Y), 0),
	}
}

func (t *TabControl) pageContentRect(rect geom.Rect) geom.Rect {
	page := t.pageRect(rect).Inset(3)
	if page.W < 0 {
		page.W = 0
	}
	if page.H < 0 {
		page.H = 0
	}
	return page
}

func (t *TabControl) layoutTabs(rect geom.Rect, measure func(string) geom.Size, lineHeight int) []geom.Rect {
	if len(t.pages) == 0 {
		return nil
	}
	tabHeight := t.tabHeight(measure, lineHeight)
	x := rect.X + 4
	out := make([]geom.Rect, 0, len(t.pages))
	for i, page := range t.pages {
		width := t.tabWidth(page, measure)
		y := rect.Y + 2
		h := tabHeight - 1
		if i == t.selected {
			y = rect.Y
			h = tabHeight + 1
		}
		out = append(out, geom.Rect{X: x, Y: y, W: width, H: h})
		x += width - 1
	}
	return out
}

func (t *TabControl) tabHeight(measure func(string) geom.Size, lineHeight int) int {
	height := maxInt(lineHeight+8, 22)
	if measure != nil {
		for _, page := range t.pages {
			size := measure(page.Title)
			height = maxInt(height, size.H+8)
		}
	}
	return height
}

func (t *TabControl) tabWidth(page *TabPage, measure func(string) geom.Size) int {
	width := 56
	if measure != nil && page != nil {
		width = maxInt(width, measure(page.Title).W+18)
	}
	return width
}

func (t *TabControl) hitTab(local geom.Point, measure func(string) geom.Size, lineHeight int) int {
	for i, tabRect := range t.layoutTabs(LocalRect(t), measure, lineHeight) {
		if tabRect.Contains(local) {
			return i
		}
	}
	return -1
}

func (t *TabControl) nextEnabled(start, delta int) int {
	if len(t.pages) == 0 {
		return -1
	}
	index := start
	for step := 0; step < len(t.pages); step++ {
		index += delta
		if index < 0 {
			index = len(t.pages) - 1
		}
		if index >= len(t.pages) {
			index = 0
		}
		if t.pages[index].Enabled {
			return index
		}
	}
	return -1
}
