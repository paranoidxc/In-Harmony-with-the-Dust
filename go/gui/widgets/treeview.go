package widgets

import (
	"classicui/event"
	"classicui/geom"
	"classicui/uicolor"
	"classicui/widget"
	"strings"
	"time"
)

type TreeNodeKind int

const (
	TreeNodeAuto TreeNodeKind = iota
	TreeNodeFolder
	TreeNodeFile
)

type TreeNode struct {
	Kind     TreeNodeKind
	Text     string
	Expanded bool
	Children []*TreeNode
	parent   *TreeNode
}

func NewTreeNode(text string, children ...*TreeNode) *TreeNode {
	node := &TreeNode{Text: text}
	node.SetChildren(children...)
	return node
}

func NewFolderNode(text string, children ...*TreeNode) *TreeNode {
	node := NewTreeNode(text, children...)
	node.Kind = TreeNodeFolder
	return node
}

func NewFileNode(text string) *TreeNode {
	return &TreeNode{
		Kind: TreeNodeFile,
		Text: text,
	}
}

func (n *TreeNode) SetChildren(children ...*TreeNode) {
	n.Children = nil
	for _, child := range children {
		n.AddChild(child)
	}
}

func (n *TreeNode) AddChild(child *TreeNode) {
	if n == nil || child == nil {
		return
	}
	child.parent = n
	n.Children = append(n.Children, child)
}

func (n *TreeNode) Parent() *TreeNode {
	if n == nil {
		return nil
	}
	return n.parent
}

func (n *TreeNode) EffectiveKind() TreeNodeKind {
	if n == nil {
		return TreeNodeFile
	}
	switch n.Kind {
	case TreeNodeFolder, TreeNodeFile:
		return n.Kind
	default:
		if len(n.Children) > 0 {
			return TreeNodeFolder
		}
		return TreeNodeFile
	}
}

type treeEntry struct {
	node    *TreeNode
	depth   int
	guides  []bool
	hasNext bool
}

type treeHotPart int

const (
	treeHotPartNone treeHotPart = iota
	treeHotPartRow
	treeHotPartExpander
)

type treeEntryLayout struct {
	expander geom.Rect
	icon     geom.Rect
	textX    int
	textY    int
}

type TreeView struct {
	widget.BaseWidget
	roots          []*TreeNode
	selection      selectionModel[*TreeNode]
	selectionOpts  TreeViewSelectionOptions
	hotNode        *TreeNode
	hotPart        treeHotPart
	pressedNode    *TreeNode
	pressedPart    treeHotPart
	pressedSelect  bool
	pressedBlank   bool
	pressedStart   geom.Point
	pressedPoint   geom.Point
	pressedMods    event.Modifiers
	dragSelecting  bool
	renamingNode   *TreeNode
	renameEdit     *Edit
	renameText     string
	focused        bool
	rowHeight      int
	indentWidth    int
	scrollbarSize  int
	topIndex       int
	scrollbar      *ScrollBar
	onChange       func(*TreeNode)
	onActivate     func(*TreeNode)
	onBeginRename  func(*TreeNode)
	onRenameCommit func(*TreeNode, string, string)
	lastClickNode  *TreeNode
	lastClickAt    time.Time
	renameNode     *TreeNode
	renameAt       time.Time
	now            func() time.Time
	doubleClick    time.Duration
	renameDelay    time.Duration
	lastDragTick   time.Time
}

func NewTreeView(id string, bounds geom.Rect, roots ...*TreeNode) *TreeView {
	tree := &TreeView{
		BaseWidget:    widget.NewBase(id, bounds),
		rowHeight:     16,
		indentWidth:   18,
		scrollbarSize: 16,
		selection:     newSelectionModel[*TreeNode](),
		selectionOpts: DefaultTreeViewSelectionOptions(),
		now:           time.Now,
		doubleClick:   500 * time.Millisecond,
		renameDelay:   500 * time.Millisecond,
	}
	tree.scrollbar = NewScrollBar(id+".scrollbar", geom.Rect{})
	tree.scrollbar.OnChange(func(value int) {
		tree.topIndex = value
	})
	tree.SetRoots(roots...)
	return tree
}

func (t *TreeView) SetRoots(roots ...*TreeNode) {
	t.roots = append([]*TreeNode(nil), roots...)
	for _, root := range t.roots {
		t.bindParents(root, nil)
	}
	visible := t.visibleNodes()
	if len(visible) == 0 {
		t.selection.Clear()
		t.clearHot()
		t.cancelRename()
		t.topIndex = 0
		t.scrollbar.SetRange(0, 1)
		return
	}
	if !t.containsNode(t.hotNode) {
		t.clearHot()
	}
	if !t.containsNode(t.renamingNode) {
		t.discardRenameEdit()
	}
	if !t.containsNode(t.renameNode) {
		t.cancelRename()
	}
	if !t.containsNode(t.SelectedNode()) {
		t.selection.SelectOnly(visible[0].node)
	} else if selected := t.SelectedNode(); selected != nil && !t.isSelected(selected) {
		t.selectOnly(selected)
	}
	t.ensureSelectedVisible()
	t.syncScrollBar(len(visible))
}

func (t *TreeView) Roots() []*TreeNode {
	return append([]*TreeNode(nil), t.roots...)
}

func (t *TreeView) SelectedNode() *TreeNode {
	if node, ok := t.selection.Lead(); ok {
		return node
	}
	return nil
}

func (t *TreeView) SelectedNodes() []*TreeNode {
	if t.selection.Count() == 0 {
		return nil
	}
	var out []*TreeNode
	for _, entry := range t.visibleNodesAll() {
		if t.selection.Contains(entry.node) {
			out = append(out, entry.node)
		}
	}
	return out
}

func (t *TreeView) SetSelectedNode(node *TreeNode) bool {
	if node == nil || !t.containsNode(node) {
		return false
	}
	t.expandAncestors(node)
	if t.SelectedNode() == node {
		return false
	}
	if t.renamingNode != nil && t.renamingNode != node {
		t.discardRenameEdit()
	}
	t.selection.SelectOnly(node)
	t.cancelRename()
	t.ensureSelectedVisible()
	if t.onChange != nil {
		t.onChange(node)
	}
	return true
}

func (t *TreeView) OnChange(fn func(*TreeNode)) {
	t.onChange = fn
}

func (t *TreeView) OnActivate(fn func(*TreeNode)) {
	t.onActivate = fn
}

func (t *TreeView) OnBeginRename(fn func(*TreeNode)) {
	t.onBeginRename = fn
}

func (t *TreeView) OnRenameCommit(fn func(*TreeNode, string, string)) {
	t.onRenameCommit = fn
}

func (t *TreeView) SelectionOptions() TreeViewSelectionOptions {
	return t.selectionOpts
}

func (t *TreeView) SetSelectionOptions(options TreeViewSelectionOptions) {
	t.selectionOpts = options
	if !t.selectionBehavior().allowsMultiSelect() {
		selected := t.SelectedNode()
		if selected != nil {
			t.selection.SelectOnly(selected)
			return
		}
		if first := t.firstSelectedVisible(); first != nil {
			t.selection.SelectOnly(first)
			return
		}
		t.selection.Clear()
	}
}

func (t *TreeView) SetMultiSelect(enabled bool) {
	if t.selectionBehavior().allowsMultiSelect() == enabled {
		return
	}
	options := t.selectionOpts
	options.MultiSelect = enabled
	t.SetSelectionOptions(options)
}

func (t *TreeView) BeginRename(node *TreeNode) bool {
	if node == nil || !t.containsNode(node) {
		return false
	}
	t.cancelRename()
	t.SetSelectedNode(node)
	t.startRename(node)
	if t.onBeginRename != nil {
		t.onBeginRename(node)
	}
	return true
}

func (t *TreeView) Paint(ctx PaintContext) error {
	if !t.Visible() {
		return nil
	}

	lineHeight := 14
	if ctx.Text != nil {
		lineHeight = ctx.Text.LineHeight()
	}
	t.rowHeight = maxInt(lineHeight+2, 16)
	t.scrollbarSize = maxInt(ctx.Theme.Metrics.ScrollbarSize, 16)

	visible := t.visibleNodes()
	t.syncScrollBar(len(visible))

	rect := ctx.BoundsFor(t)
	ctx.Canvas.FillRect(rect, ctx.Theme.Colors.Face)
	ctx.Canvas.DrawDoubleBevel(rect, ctx.Theme.Colors.Shadow, ctx.Theme.Colors.Lightest, ctx.Theme.Colors.DarkShadow, ctx.Theme.Colors.Light)

	itemsRect := t.itemsRect(rect)
	ctx.Canvas.FillRect(itemsRect, ctx.Theme.Colors.Window)
	ctx.Canvas.PushClip(itemsRect)
	for row := 0; row < t.visibleRows(); row++ {
		index := t.topIndex + row
		if index >= len(visible) {
			break
		}
		entry := visible[index]
		rowRect := geom.Rect{
			X: itemsRect.X + 1,
			Y: itemsRect.Y + row*t.rowHeight,
			W: itemsRect.W - 2,
			H: t.rowHeight,
		}
		if t.isSelected(entry.node) {
			ctx.Canvas.FillRect(rowRect, ctx.Theme.Colors.Highlight)
		} else if entry.node == t.hotNode && t.hotPart != treeHotPartNone {
			ctx.Canvas.FillRect(rowRect, blendColor(ctx.Theme.Colors.Window, ctx.Theme.Colors.Lightest))
		}
		if err := t.paintEntry(ctx, entry, rowRect, lineHeight); err != nil {
			ctx.Canvas.PopClip()
			return err
		}
	}
	if rect, ok := t.marqueeRect(); ok {
		paintSelectionMarquee(ctx, rect)
	}
	if t.renamingNode != nil && t.renameEdit != nil && t.renameEdit.Visible() {
		t.syncRenameEditBounds(func(text string) geom.Size {
			if ctx.Text == nil {
				return geom.Size{}
			}
			return ctx.Text.MeasureString(text)
		}, lineHeight)
		childCtx := ctx.Child(t)
		if err := t.renameEdit.Paint(childCtx); err != nil {
			ctx.Canvas.PopClip()
			return err
		}
	}
	ctx.Canvas.PopClip()

	sbCtx := ctx
	sbCtx.Origin = geom.Point{X: rect.X, Y: rect.Y}
	if err := t.scrollbar.Paint(sbCtx); err != nil {
		return err
	}

	if t.focused {
		ctx.Canvas.DrawFocusRect(itemsRect.Inset(1), ctx.Theme.Colors.DarkShadow)
	}
	return nil
}

func (t *TreeView) MouseEnter(EventContext) {}

func (t *TreeView) MouseLeave(ctx EventContext) {
	if t.clearHot() {
		ctx.Invalidate(t)
	}
}

func (t *TreeView) MouseMove(ctx EventContext, local geom.Point) {
	if t.renamingNode != nil && t.renameEdit != nil && t.renameEdit.Visible() {
		if t.renameEdit.Bounds().Contains(local) {
			t.renameEdit.MouseMove(ctx, geom.Point{X: local.X - t.renameEdit.Bounds().X, Y: local.Y - t.renameEdit.Bounds().Y})
			return
		}
	}
	changed := t.updateHot(local)
	if (t.pressedNode != nil && t.pressedPart == treeHotPartRow) || t.pressedBlank {
		startedDrag := false
		if !t.dragSelecting {
			dx := local.X - t.pressedPoint.X
			dy := local.Y - t.pressedPoint.Y
			if dx*dx+dy*dy >= 9 {
				t.dragSelecting = true
				startedDrag = true
			}
		}
		t.pressedPoint = local
		if t.dragSelecting && (t.updateDragSelection(local, true) || startedDrag) {
			changed = true
		}
	}
	if t.scrollbar.dragging || t.scrollbar.Bounds().Contains(local) {
		t.scrollbar.MouseMove(ctx, geom.Point{X: local.X - t.scrollbar.Bounds().X, Y: local.Y - t.scrollbar.Bounds().Y})
		value := t.scrollbar.Value()
		if value != t.topIndex {
			t.topIndex = value
			t.cancelRename()
			changed = true
		}
	}
	if changed {
		ctx.Invalidate(t)
	}
}

func (t *TreeView) MouseDown(ctx EventContext, ev event.MouseButtonEvent, local geom.Point) {
	if ev.Button != event.MouseButtonLeft || !LocalContains(t, local) {
		return
	}
	if t.renamingNode != nil && t.renameEdit != nil && t.renameEdit.Visible() {
		if t.renameEdit.Bounds().Contains(local) {
			t.renameEdit.MouseDown(ctx, ev, geom.Point{X: local.X - t.renameEdit.Bounds().X, Y: local.Y - t.renameEdit.Bounds().Y})
			return
		}
		t.commitRename(ctx)
	}
	if t.scrollbar.Bounds().Contains(local) {
		t.scrollbar.MouseDown(ctx, ev, geom.Point{X: local.X - t.scrollbar.Bounds().X, Y: local.Y - t.scrollbar.Bounds().Y})
		t.topIndex = t.scrollbar.Value()
		ctx.Invalidate(t)
		return
	}
	itemsRect := t.itemsRect(LocalRect(t))
	if !itemsRect.Contains(local) {
		t.cancelRename()
		t.clearPressed()
		return
	}
	mods := t.selectionBehavior().normalizeModifiers(ev.Modifiers)
	entry, rowRect, ok := t.entryAt(local)
	if !ok {
		ctx.SetFocus(t)
		t.cancelRename()
		t.pressedBlank = true
		t.pressedStart = local
		t.pressedPoint = local
		t.pressedMods = mods
		t.dragSelecting = false
		t.pressedNode = nil
		t.pressedPart = treeHotPartNone
		t.pressedSelect = false
		t.captureDragBaseSelection()
		t.lastClickNode = nil
		t.lastClickAt = time.Time{}
		if mods&(event.ModCtrl|event.ModShift) == 0 && t.clearSelection() {
			ctx.Invalidate(t)
		}
		return
	}
	if entry.node == nil {
		return
	}
	ctx.SetFocus(t)
	if expander := t.expanderRect(rowRect, entry.depth); len(entry.node.Children) > 0 && expander.Contains(local) {
		t.cancelRename()
		t.clearPressed()
		t.lastClickNode = nil
		t.lastClickAt = time.Time{}
		t.hotNode = entry.node
		t.hotPart = treeHotPartExpander
		t.toggleExpanded(entry.node)
		ctx.Invalidate(t)
		return
	}
	wasSelected := t.isSelected(entry.node)
	t.pressedNode = entry.node
	t.pressedPart = treeHotPartRow
	t.pressedSelect = wasSelected
	t.pressedBlank = false
	t.pressedStart = local
	t.pressedPoint = local
	t.pressedMods = mods
	t.dragSelecting = false
	t.captureDragBaseSelection()
	needsInvalidate := false
	switch {
	case mods&event.ModShift != 0:
		needsInvalidate = t.selectRangeTo(entry.node)
	case mods&event.ModCtrl != 0:
		needsInvalidate = t.toggleSelection(entry.node)
	default:
		needsInvalidate = t.SetSelectedNode(entry.node)
	}
	if t.isDoubleClick(entry.node) {
		t.cancelRename()
		t.clearPressed()
		clear(t.selection.dragBaseSet)
		if len(entry.node.Children) > 0 {
			t.toggleExpanded(entry.node)
			needsInvalidate = true
		} else if t.onActivate != nil {
			t.onActivate(entry.node)
		}
	} else if !wasSelected || mods != 0 {
		t.cancelRename()
	}
	if needsInvalidate {
		ctx.Invalidate(t)
	}
}

func (t *TreeView) MouseUp(ctx EventContext, ev event.MouseButtonEvent, local geom.Point) {
	if ev.Button != event.MouseButtonLeft {
		return
	}
	if t.renameEdit != nil && t.renameEdit.selecting {
		t.renameEdit.MouseUp(ctx, ev, geom.Point{X: local.X - t.renameEdit.Bounds().X, Y: local.Y - t.renameEdit.Bounds().Y})
		return
	}
	if t.scrollbar.dragging || t.scrollbar.Bounds().Contains(local) {
		t.scrollbar.MouseUp(ctx, ev, geom.Point{X: local.X - t.scrollbar.Bounds().X, Y: local.Y - t.scrollbar.Bounds().Y})
		t.topIndex = t.scrollbar.Value()
		t.cancelRename()
		ctx.Invalidate(t)
		t.clearPressed()
		clear(t.selection.dragBaseSet)
		return
	}
	if !t.dragSelecting && t.pressedNode != nil && t.pressedPart == treeHotPartRow {
		entry, _, ok := t.entryAt(local)
		if ok && entry.node == t.pressedNode && t.pressedSelect {
			t.scheduleRename(entry.node)
		} else {
			t.cancelRename()
		}
	}
	t.clearPressed()
	clear(t.selection.dragBaseSet)
}

func (t *TreeView) MouseWheel(ctx EventContext, ev event.MouseWheel, _ geom.Point) bool {
	visible := t.visibleNodes()
	if len(visible) == 0 {
		return true
	}
	if t.renamingNode != nil {
		t.commitRename(ctx)
	}
	t.cancelRename()
	t.topIndex = clampInt(t.topIndex-ev.Delta, 0, t.maxTopIndex(len(visible)))
	t.scrollbar.SetValue(t.topIndex)
	ctx.Invalidate(t)
	return true
}

func (t *TreeView) KeyDown(ctx EventContext, ev event.KeyEvent) bool {
	visible := t.visibleNodes()
	selected := t.SelectedNode()
	if !t.Enabled() || len(visible) == 0 {
		return false
	}
	if t.renamingNode != nil {
		switch ev.Key {
		case event.KeyEnter:
			t.commitRename(ctx)
			ctx.Invalidate(t)
			return true
		case event.KeyEscape:
			t.cancelRenameEdit(ctx)
			ctx.Invalidate(t)
			return true
		case event.KeyTab:
			t.commitRename(ctx)
			ctx.Invalidate(t)
			return false
		default:
			if t.renameEdit != nil && t.renameEdit.KeyDown(ctx, ev) {
				return true
			}
			return false
		}
	}
	if selected == nil {
		if t.selectionBehavior().selectAllShortcut(ev) {
			if t.selectAllVisible() {
				ctx.Invalidate(t)
			}
			return true
		}
		if t.selectionBehavior().toggleLeadShortcut(ev) {
			if node := t.recoveryNodeForKey(visible, ev.Key); node != nil && t.toggleSelection(node) {
				ctx.Invalidate(t)
			}
			return true
		}
		if node := t.recoveryNodeForKey(visible, ev.Key); node != nil {
			if t.selectionBehavior().extendRange(ev.Modifiers) {
				if t.selectRangeTo(node) {
					ctx.Invalidate(t)
				}
			} else if t.SetSelectedNode(node) {
				ctx.Invalidate(t)
			}
			return true
		}
		return false
	}
	if t.selectionBehavior().selectAllShortcut(ev) {
		if t.selectAllVisible() {
			ctx.Invalidate(t)
		}
		return true
	}
	t.cancelRename()

	selectedIndex := t.selectedIndex(visible, selected)
	nextIndex := selectedIndex
	changed := false

	switch ev.Key {
	case event.KeyUp:
		nextIndex--
	case event.KeyDown:
		nextIndex++
	case event.KeyHome:
		nextIndex = 0
	case event.KeyEnd:
		nextIndex = len(visible) - 1
	case event.KeyPageUp:
		nextIndex -= maxInt(t.visibleRows()-1, 1)
	case event.KeyPageDown:
		nextIndex += maxInt(t.visibleRows()-1, 1)
	case event.KeyLeft:
		if selected.Expanded && len(selected.Children) > 0 {
			selected.Expanded = false
			if selected != nil {
				t.ensureSelectedVisible()
			}
			ctx.Invalidate(t)
			return true
		}
		if parent := selected.Parent(); parent != nil {
			changed = t.SetSelectedNode(parent)
		}
	case event.KeyRight:
		if len(selected.Children) > 0 {
			if !selected.Expanded {
				selected.Expanded = true
				t.ensureSelectedVisible()
				ctx.Invalidate(t)
				return true
			}
			changed = t.SetSelectedNode(selected.Children[0])
		}
	case event.KeySpace:
		if t.selectionBehavior().toggleLeadShortcut(ev) {
			changed = t.toggleSelection(selected)
			break
		}
		if len(selected.Children) > 0 {
			t.toggleExpanded(selected)
			ctx.Invalidate(t)
			return true
		}
		if t.onActivate != nil {
			t.onActivate(selected)
		}
		return true
	case event.KeyEnter:
		if len(selected.Children) > 0 {
			t.toggleExpanded(selected)
			ctx.Invalidate(t)
			return true
		}
		if t.onActivate != nil {
			t.onActivate(selected)
		}
		return true
	case event.KeyF2:
		return t.beginRenameWithContext(ctx, selected)
	default:
		return false
	}

	if nextIndex != selectedIndex && (ev.Key == event.KeyUp || ev.Key == event.KeyDown || ev.Key == event.KeyHome || ev.Key == event.KeyEnd || ev.Key == event.KeyPageUp || ev.Key == event.KeyPageDown) {
		nextIndex = clampInt(nextIndex, 0, len(visible)-1)
		target := visible[nextIndex].node
		if t.selectionBehavior().extendRange(ev.Modifiers) {
			changed = t.selectRangeTo(target)
		} else {
			changed = t.SetSelectedNode(target)
		}
	}
	if changed {
		ctx.Invalidate(t)
	}
	return true
}

func (t *TreeView) CanFocus() bool {
	return t.Visible() && t.Enabled()
}

func (t *TreeView) Tick(ctx EventContext, now time.Time) bool {
	changed := false
	if t.renameNode != nil && !t.renameAt.IsZero() && !now.Before(t.renameAt) {
		node := t.renameNode
		t.cancelRename()
		if node != nil && node == t.SelectedNode() && t.containsNode(node) {
			t.beginRenameWithContext(ctx, node)
			changed = true
		}
	}
	if t.renamingNode != nil && t.renameEdit != nil && t.renameEdit.Tick(ctx, now) {
		changed = true
	}
	if t.dragSelecting && t.autoScrollDrag(now) {
		changed = true
	}
	if changed {
		ctx.Invalidate(t)
	}
	return changed
}

func (t *TreeView) SetFocused(focused bool) {
	t.focused = focused
	if !focused {
		t.cancelRename()
		t.clearPressed()
	}
	if t.renameEdit != nil {
		t.renameEdit.SetFocused(focused && t.renamingNode != nil)
	}
}

func (t *TreeView) Focused() bool {
	return t.focused
}

func (t *TreeView) FocusGained(ctx EventContext) {
	if t.renamingNode != nil && t.renameEdit != nil {
		t.renameEdit.SetFocused(true)
		t.renameEdit.FocusGained(ctx)
	}
}

func (t *TreeView) FocusLost(ctx EventContext) {
	if t.renamingNode != nil {
		t.commitRename(ctx)
	}
	if t.renameEdit != nil {
		t.renameEdit.SetFocused(false)
		t.renameEdit.FocusLost(ctx)
	}
}

func (t *TreeView) TextInput(ctx EventContext, ev event.TextInput) bool {
	if t.renamingNode == nil || t.renameEdit == nil {
		return false
	}
	return t.renameEdit.TextInput(ctx, ev)
}

func (t *TreeView) TextEditing(ctx EventContext, ev event.TextEditing) bool {
	if t.renamingNode == nil || t.renameEdit == nil {
		return false
	}
	return t.renameEdit.TextEditing(ctx, ev)
}

func (t *TreeView) TextInputRect(ctx EventContext) geom.Rect {
	if t.renamingNode == nil || t.renameEdit == nil {
		return geom.Rect{}
	}
	t.syncRenameEditBounds(ctx.MeasureText, ctx.LineHeight())
	rect := t.renameEdit.TextInputRect(ctx)
	return rect.Move(t.renameEdit.Bounds().X, t.renameEdit.Bounds().Y)
}

func (t *TreeView) paintEntry(ctx PaintContext, entry treeEntry, rowRect geom.Rect, lineHeight int) error {
	textHeight := lineHeight
	if ctx.Text != nil && entry.node.Text != "" {
		if size := ctx.Text.MeasureString(entry.node.Text); size.H > 0 {
			textHeight = size.H
		}
	}
	layout := t.layoutEntry(rowRect, entry.depth, textHeight)
	selected := t.isSelected(entry.node)

	textColor := ctx.Theme.Colors.WindowText
	guideColor := ctx.Theme.Colors.Shadow
	if selected {
		textColor = ctx.Theme.Colors.HighlightText
		guideColor = ctx.Theme.Colors.HighlightText
	}
	t.paintGuides(ctx, entry, rowRect, layout, guideColor)

	if len(entry.node.Children) > 0 {
		t.paintExpander(ctx, layout.expander, entry.node.Expanded, selected, entry.node == t.hotNode && t.hotPart == treeHotPartExpander)
	}
	t.paintIcon(ctx, layout.icon, entry.node.EffectiveKind(), selected)
	if ctx.Text == nil || entry.node.Text == "" {
		return nil
	}
	return ctx.Text.DrawString(ctx.Canvas, geom.Point{X: layout.textX, Y: layout.textY}, entry.node.Text, textColor)
}

func (t *TreeView) paintExpander(ctx PaintContext, rect geom.Rect, expanded, selected, hot bool) {
	fill := ctx.Theme.Colors.Window
	stroke := ctx.Theme.Colors.DarkShadow
	if !selected && hot {
		fill = blendColor(ctx.Theme.Colors.Window, ctx.Theme.Colors.Lightest)
	}
	ctx.Canvas.FillRect(rect, fill)
	ctx.Canvas.FrameRect(rect, stroke)
	midX := rect.X + rect.W/2
	midY := rect.Y + rect.H/2
	ctx.Canvas.DrawHLine(rect.X+2, midY, maxInt(rect.W-4, 0), stroke)
	if !expanded {
		ctx.Canvas.DrawVLine(midX, rect.Y+2, maxInt(rect.H-4, 0), stroke)
	}
}

func (t *TreeView) paintGuides(ctx PaintContext, entry treeEntry, rowRect geom.Rect, layout treeEntryLayout, color uicolor.RGBA) {
	for depth, hasContinuation := range entry.guides {
		if !hasContinuation {
			continue
		}
		ctx.Canvas.DrawVLine(t.guideX(rowRect, depth), rowRect.Y, rowRect.H, color)
	}
	if entry.depth == 0 {
		return
	}
	guideX := layout.expander.X + layout.expander.W/2
	midY := rowRect.Y + rowRect.H/2
	ctx.Canvas.DrawVLine(guideX, rowRect.Y, maxInt(midY-rowRect.Y+1, 0), color)
	if entry.hasNext {
		ctx.Canvas.DrawVLine(guideX, midY, maxInt(rowRect.Bottom()-midY, 0), color)
	}
	ctx.Canvas.DrawHLine(guideX, midY, maxInt(layout.icon.X-guideX+2, 0), color)
}

func (t *TreeView) paintIcon(ctx PaintContext, rect geom.Rect, kind TreeNodeKind, selected bool) {
	stroke := ctx.Theme.Colors.DarkShadow
	fill := ctx.Theme.Colors.Window
	if selected {
		stroke = ctx.Theme.Colors.HighlightText
		fill = blendColor(ctx.Theme.Colors.Highlight, ctx.Theme.Colors.HighlightText)
	}
	switch kind {
	case TreeNodeFolder:
		t.paintFolderIcon(ctx, rect, fill, stroke)
	default:
		t.paintFileIcon(ctx, rect, fill, stroke)
	}
}

func (t *TreeView) paintFolderIcon(ctx PaintContext, rect geom.Rect, fill, stroke uicolor.RGBA) {
	tab := geom.Rect{
		X: rect.X + 1,
		Y: rect.Y + 1,
		W: maxInt(rect.W/2, 4),
		H: 3,
	}
	body := geom.Rect{
		X: rect.X,
		Y: rect.Y + 3,
		W: rect.W - 1,
		H: rect.H - 4,
	}
	ctx.Canvas.FillRect(tab, fill)
	ctx.Canvas.FillRect(body, fill)
	ctx.Canvas.FrameRect(body, stroke)
	ctx.Canvas.DrawHLine(tab.X, tab.Y, tab.W, stroke)
	ctx.Canvas.DrawVLine(tab.X, tab.Y, tab.H, stroke)
	ctx.Canvas.DrawHLine(tab.X, tab.Bottom()-1, tab.W, stroke)
}

func (t *TreeView) paintFileIcon(ctx PaintContext, rect geom.Rect, fill, stroke uicolor.RGBA) {
	body := geom.Rect{
		X: rect.X + 1,
		Y: rect.Y,
		W: rect.W - 2,
		H: rect.H - 1,
	}
	ctx.Canvas.FillRect(body, fill)
	ctx.Canvas.FrameRect(body, stroke)
	foldX := body.Right() - 4
	ctx.Canvas.DrawHLine(foldX, body.Y+1, 3, stroke)
	ctx.Canvas.DrawVLine(body.Right()-2, body.Y+1, 3, stroke)
	ctx.Canvas.DrawPixel(foldX, body.Y+2, fill)
}

func (t *TreeView) iconRect(rowRect geom.Rect, depth int) geom.Rect {
	return geom.Rect{
		X: t.expanderRect(rowRect, depth).Right() + 4,
		Y: rowRect.Y + maxInt((rowRect.H-12)/2, 0),
		W: 12,
		H: 12,
	}
}

func (t *TreeView) layoutEntry(rowRect geom.Rect, depth, textHeight int) treeEntryLayout {
	contentHeight := maxInt(textHeight, 12)
	contentY := rowRect.Y + maxInt((rowRect.H-contentHeight)/2, 0)
	expander := t.expanderRect(rowRect, depth)
	expander.Y = contentY + maxInt((contentHeight-expander.H)/2, 0)
	icon := geom.Rect{
		X: expander.Right() + 4,
		Y: contentY + maxInt((contentHeight-12)/2, 0),
		W: 12,
		H: 12,
	}
	return treeEntryLayout{
		expander: expander,
		icon:     icon,
		textX:    icon.Right() + 4,
		textY:    contentY + maxInt((contentHeight-textHeight)/2, 0),
	}
}

func (t *TreeView) guideX(rowRect geom.Rect, depth int) int {
	return t.expanderRect(rowRect, depth).X + 4
}

func (t *TreeView) isDoubleClick(node *TreeNode) bool {
	now := time.Now()
	if t.now != nil {
		now = t.now()
	}
	doubleClick := node != nil &&
		node == t.lastClickNode &&
		!t.lastClickAt.IsZero() &&
		now.Sub(t.lastClickAt) <= t.doubleClick
	t.lastClickNode = node
	t.lastClickAt = now
	if doubleClick {
		t.lastClickNode = nil
		t.lastClickAt = time.Time{}
	}
	return doubleClick
}

func (t *TreeView) expanderRect(rowRect geom.Rect, depth int) geom.Rect {
	return geom.Rect{
		X: rowRect.X + 2 + depth*t.indentWidth,
		Y: rowRect.Y + maxInt((rowRect.H-9)/2, 0),
		W: 9,
		H: 9,
	}
}

func (t *TreeView) itemsRect(rect geom.Rect) geom.Rect {
	return geom.Rect{
		X: rect.X + 2,
		Y: rect.Y + 2,
		W: maxInt(rect.W-t.scrollbarSize-4, 0),
		H: maxInt(rect.H-4, 0),
	}
}

func (t *TreeView) syncScrollBar(total int) {
	rect := LocalRect(t)
	t.scrollbar.SetBounds(geom.Rect{
		X: maxInt(rect.W-t.scrollbarSize-1, 1),
		Y: 1,
		W: maxInt(t.scrollbarSize, 1),
		H: maxInt(rect.H-2, 0),
	})
	t.scrollbar.SetRange(t.maxTopIndex(total), maxInt(t.visibleRows(), 1))
	t.scrollbar.SetValue(t.topIndex)
}

func (t *TreeView) visibleRows() int {
	rows := t.itemsRect(LocalRect(t)).H / maxInt(t.rowHeight, 1)
	return maxInt(rows, 1)
}

func (t *TreeView) maxTopIndex(total int) int {
	return maxInt(total-t.visibleRows(), 0)
}

func (t *TreeView) bindParents(node, parent *TreeNode) {
	if node == nil {
		return
	}
	node.parent = parent
	for _, child := range node.Children {
		t.bindParents(child, node)
	}
}

func (t *TreeView) visibleNodes() []treeEntry {
	var out []treeEntry
	for index, root := range t.roots {
		t.appendVisible(&out, root, 0, nil, index < len(t.roots)-1)
	}
	return out
}

func (t *TreeView) appendVisible(out *[]treeEntry, node *TreeNode, depth int, guides []bool, hasNext bool) {
	if node == nil {
		return
	}
	*out = append(*out, treeEntry{
		node:    node,
		depth:   depth,
		guides:  append([]bool(nil), guides...),
		hasNext: hasNext,
	})
	if !node.Expanded {
		return
	}
	childGuides := append(append([]bool(nil), guides...), hasNext)
	for index, child := range node.Children {
		t.appendVisible(out, child, depth+1, childGuides, index < len(node.Children)-1)
	}
}

func (t *TreeView) entryAt(local geom.Point) (treeEntry, geom.Rect, bool) {
	itemsRect := t.itemsRect(LocalRect(t))
	if !itemsRect.Contains(local) {
		return treeEntry{}, geom.Rect{}, false
	}
	row := (local.Y - itemsRect.Y) / maxInt(t.rowHeight, 1)
	index := t.topIndex + row
	visible := t.visibleNodes()
	if index < 0 || index >= len(visible) {
		return treeEntry{}, geom.Rect{}, false
	}
	rowRect := geom.Rect{
		X: itemsRect.X + 1,
		Y: itemsRect.Y + row*t.rowHeight,
		W: itemsRect.W - 2,
		H: t.rowHeight,
	}
	return visible[index], rowRect, true
}

func (t *TreeView) rowRectForNode(node *TreeNode) (geom.Rect, treeEntry, bool) {
	if node == nil {
		return geom.Rect{}, treeEntry{}, false
	}
	visible := t.visibleNodes()
	index := t.selectedIndex(visible, node)
	if index < 0 {
		return geom.Rect{}, treeEntry{}, false
	}
	row := index - t.topIndex
	if row < 0 || row >= t.visibleRows() {
		return geom.Rect{}, treeEntry{}, false
	}
	itemsRect := t.itemsRect(LocalRect(t))
	return geom.Rect{
		X: itemsRect.X + 1,
		Y: itemsRect.Y + row*t.rowHeight,
		W: itemsRect.W - 2,
		H: t.rowHeight,
	}, visible[index], true
}

func (t *TreeView) updateHot(local geom.Point) bool {
	if !LocalContains(t, local) || t.scrollbar.Bounds().Contains(local) {
		return t.clearHot()
	}
	entry, rowRect, ok := t.entryAt(local)
	if !ok || entry.node == nil {
		return t.clearHot()
	}
	part := treeHotPartRow
	if len(entry.node.Children) > 0 && t.expanderRect(rowRect, entry.depth).Contains(local) {
		part = treeHotPartExpander
	}
	if t.hotNode == entry.node && t.hotPart == part {
		return false
	}
	t.hotNode = entry.node
	t.hotPart = part
	return true
}

func (t *TreeView) clearHot() bool {
	if t.hotNode == nil && t.hotPart == treeHotPartNone {
		return false
	}
	t.hotNode = nil
	t.hotPart = treeHotPartNone
	return true
}

func (t *TreeView) ensureRenameEdit() *Edit {
	if t.renameEdit != nil {
		return t.renameEdit
	}
	edit := NewEdit(t.ID()+".rename", geom.Rect{})
	edit.SetVisible(false)
	t.renameEdit = edit
	return edit
}

func (t *TreeView) startRename(node *TreeNode) {
	if node == nil {
		return
	}
	edit := t.ensureRenameEdit()
	t.renamingNode = node
	t.renameText = node.Text
	edit.SetVisible(true)
	edit.SetText(node.Text)
	anchor, caret := t.renameSelection(node)
	edit.anchor = anchor
	edit.caret = caret
	edit.selecting = false
	edit.composition = ""
	if t.focused {
		edit.SetFocused(true)
		edit.resetCaretBlink()
	}
}

func (t *TreeView) renameSelection(node *TreeNode) (int, int) {
	text := ""
	if node != nil {
		text = node.Text
	}
	runes := []rune(text)
	end := len(runes)
	if node == nil || node.EffectiveKind() != TreeNodeFile {
		return 0, end
	}
	stemEnd := fileStemEnd(text)
	if stemEnd <= 0 || stemEnd >= end {
		return 0, end
	}
	return 0, stemEnd
}

func fileStemEnd(text string) int {
	if text == "" {
		return 0
	}
	lastDot := strings.LastIndex(text, ".")
	if lastDot <= 0 || lastDot == len(text)-1 {
		return len([]rune(text))
	}
	return len([]rune(text[:lastDot]))
}

func (t *TreeView) beginRenameWithContext(ctx EventContext, node *TreeNode) bool {
	if !t.BeginRename(node) {
		return false
	}
	if ctx != nil {
		t.syncRenameEditBounds(ctx.MeasureText, ctx.LineHeight())
		if t.renameEdit != nil && t.focused {
			t.renameEdit.FocusGained(ctx)
		}
		ctx.Invalidate(t)
	}
	return true
}

func (t *TreeView) commitRename(ctx EventContext) bool {
	if t.renamingNode == nil || t.renameEdit == nil {
		return false
	}
	node := t.renamingNode
	oldText := t.renameText
	newText := t.renameEdit.Text()
	if newText == "" {
		newText = oldText
	}
	t.discardRenameEdit()
	if node != nil {
		node.Text = newText
		if t.onRenameCommit != nil && newText != oldText {
			t.onRenameCommit(node, oldText, newText)
		}
	}
	if ctx != nil {
		ctx.ReleaseCapture(t.renameEdit)
	}
	return true
}

func (t *TreeView) cancelRenameEdit(ctx EventContext) bool {
	if t.renamingNode == nil || t.renameEdit == nil {
		return false
	}
	if ctx != nil {
		ctx.ReleaseCapture(t.renameEdit)
	}
	t.discardRenameEdit()
	return true
}

func (t *TreeView) discardRenameEdit() {
	t.renamingNode = nil
	t.renameText = ""
	if t.renameEdit != nil {
		t.renameEdit.SetVisible(false)
		t.renameEdit.SetFocused(false)
		t.renameEdit.selecting = false
		t.renameEdit.composition = ""
	}
}

func (t *TreeView) renameEditBounds(measure func(string) geom.Size, lineHeight int) (geom.Rect, bool) {
	if t.renamingNode == nil {
		return geom.Rect{}, false
	}
	rowRect, entry, ok := t.rowRectForNode(t.renamingNode)
	if !ok {
		return geom.Rect{}, false
	}
	currentText := t.renamingNode.Text
	if t.renameEdit != nil {
		currentText = t.renameEdit.Text()
	}
	textHeight := lineHeight
	if measure != nil && currentText != "" {
		if size := measure(currentText); size.H > 0 {
			textHeight = size.H
		}
	}
	layout := t.layoutEntry(rowRect, entry.depth, textHeight)
	textWidth := 64
	if measure != nil {
		textWidth = maxInt(measure(currentText).W+12, 64)
	}
	maxWidth := maxInt(rowRect.Right()-layout.textX-4, 48)
	return geom.Rect{
		X: layout.textX - 2,
		Y: rowRect.Y,
		W: minInt(textWidth, maxWidth),
		H: rowRect.H,
	}, true
}

func (t *TreeView) syncRenameEditBounds(measure func(string) geom.Size, lineHeight int) {
	if t.renameEdit == nil {
		return
	}
	if bounds, ok := t.renameEditBounds(measure, lineHeight); ok {
		t.renameEdit.SetBounds(bounds)
		t.renameEdit.SetVisible(true)
		return
	}
	t.renameEdit.SetVisible(false)
}

func (t *TreeView) clearPressed() {
	t.pressedNode = nil
	t.pressedPart = treeHotPartNone
	t.pressedSelect = false
	t.pressedBlank = false
}

func (t *TreeView) scheduleRename(node *TreeNode) {
	if node == nil || node != t.SelectedNode() {
		t.cancelRename()
		return
	}
	now := time.Now()
	if t.now != nil {
		now = t.now()
	}
	t.renameNode = node
	t.renameAt = now.Add(t.renameDelay)
}

func (t *TreeView) cancelRename() {
	t.renameNode = nil
	t.renameAt = time.Time{}
}

func (t *TreeView) containsNode(target *TreeNode) bool {
	if target == nil {
		return false
	}
	for _, entry := range t.visibleNodesAll() {
		if entry.node == target {
			return true
		}
	}
	return false
}

func (t *TreeView) visibleNodesAll() []treeEntry {
	var out []treeEntry
	for _, root := range t.roots {
		t.appendAll(&out, root, 0)
	}
	return out
}

func (t *TreeView) selectOnly(node *TreeNode) {
	if node != nil {
		t.selection.SelectOnly(node)
	}
}

func (t *TreeView) clearSelection() bool {
	return t.selection.Clear()
}

func (t *TreeView) isSelected(node *TreeNode) bool {
	if node == nil {
		return false
	}
	return t.selection.Contains(node)
}

func (t *TreeView) toggleSelection(node *TreeNode) bool {
	if node == nil || !t.containsNode(node) {
		return false
	}
	if !t.selection.Toggle(node, t.visibleSelectionOrder()) {
		return false
	}
	t.ensureSelectedVisible()
	if selected := t.SelectedNode(); t.onChange != nil && selected != nil {
		t.onChange(selected)
	}
	return true
}

func (t *TreeView) selectRangeTo(node *TreeNode) bool {
	if node == nil || !t.containsNode(node) {
		return false
	}
	t.expandAncestors(node)
	if t.visibleSelectionOrder().Len() == 0 {
		return false
	}
	if !t.selection.SelectRange(t.visibleSelectionOrder(), node) {
		return false
	}
	t.ensureSelectedVisible()
	if t.onChange != nil {
		t.onChange(node)
	}
	return true
}

func (t *TreeView) selectAllVisible() bool {
	if !t.selectionBehavior().allowsMultiSelect() {
		return false
	}
	order := t.visibleSelectionOrder()
	if order.Len() == 0 {
		return false
	}
	changed := t.selection.SelectAll(order)
	t.ensureSelectedVisible()
	if changed && t.onChange != nil && t.SelectedNode() != nil {
		t.onChange(t.SelectedNode())
	}
	return changed
}

func (t *TreeView) captureDragBaseSelection() {
	t.selection.CaptureDragBase()
}

func (t *TreeView) updateDragSelection(local geom.Point, clampToItems bool) bool {
	if t.pressedBlank {
		if !t.selectionBehavior().allowsBlankDrag() {
			return false
		}
		return t.selectDragMarquee()
	}
	if t.pressedNode == nil {
		return false
	}
	currentPointer := t.pressedPoint
	entry, _, ok := t.entryAtClamped(local, clampToItems)
	if !ok || entry.node == nil {
		return false
	}
	t.pressedPoint = currentPointer
	if !t.selectionBehavior().allowsMultiSelect() {
		previous := t.SelectedNode()
		changed := t.SetSelectedNode(entry.node)
		return changed || previous != t.SelectedNode()
	}
	switch {
	case t.pressedMods&event.ModShift != 0:
		return t.selectRangeTo(entry.node)
	case t.pressedMods&event.ModCtrl != 0:
		return t.selectDragUnion(entry.node)
	default:
		return t.selectDragRange(entry.node)
	}
}

func (t *TreeView) autoScrollDrag(now time.Time) bool {
	if !t.lastDragTick.IsZero() && now.Sub(t.lastDragTick) < 50*time.Millisecond {
		return false
	}
	t.lastDragTick = now
	if t.pressedPoint.Y < 0 {
		next := clampInt(t.topIndex-1, 0, t.maxTopIndex(len(t.visibleNodes())))
		if next == t.topIndex {
			return false
		}
		t.topIndex = next
		t.scrollbar.SetValue(t.topIndex)
		t.updateDragSelection(t.pressedPoint, false)
		return true
	}
	if t.pressedPoint.Y >= LocalRect(t).H {
		next := clampInt(t.topIndex+1, 0, t.maxTopIndex(len(t.visibleNodes())))
		if next == t.topIndex {
			return false
		}
		t.topIndex = next
		t.scrollbar.SetValue(t.topIndex)
		t.updateDragSelection(t.pressedPoint, false)
		return true
	}
	return false
}

func (t *TreeView) marqueeRect() (geom.Rect, bool) {
	if !t.selectionBehavior().allowsBlankDrag() || !t.dragSelecting || (!t.pressedBlank && (t.pressedNode == nil || t.pressedPart != treeHotPartRow)) {
		return geom.Rect{}, false
	}
	return dragMarqueeRect(t.itemsRect(LocalRect(t)), t.pressedStart, t.pressedPoint)
}

func (t *TreeView) entryAtClamped(local geom.Point, clampToItems bool) (treeEntry, geom.Rect, bool) {
	itemsRect := t.itemsRect(LocalRect(t))
	if itemsRect.W <= 0 || itemsRect.H <= 0 {
		return treeEntry{}, geom.Rect{}, false
	}
	if !clampToItems && !itemsRect.Contains(local) && local.Y >= itemsRect.Y && local.Y < itemsRect.Bottom() {
		return treeEntry{}, geom.Rect{}, false
	}
	y := local.Y
	if clampToItems {
		y = clampInt(y, itemsRect.Y, itemsRect.Bottom()-1)
	} else if y < itemsRect.Y || y >= itemsRect.Bottom() {
		return treeEntry{}, geom.Rect{}, false
	}
	row := (y - itemsRect.Y) / maxInt(t.rowHeight, 1)
	visible := t.visibleNodes()
	if len(visible) == 0 {
		return treeEntry{}, geom.Rect{}, false
	}
	index := clampInt(t.topIndex+row, 0, len(visible)-1)
	rowRect := geom.Rect{
		X: itemsRect.X + 1,
		Y: itemsRect.Y + (index-t.topIndex)*t.rowHeight,
		W: itemsRect.W - 2,
		H: t.rowHeight,
	}
	return visible[index], rowRect, true
}

func (t *TreeView) selectDragRange(node *TreeNode) bool {
	if node == nil {
		return false
	}
	if t.pressedNode == nil {
		return false
	}
	return t.selection.SelectDragRange(t.visibleSelectionOrder(), t.pressedNode, node)
}

func (t *TreeView) selectDragUnion(node *TreeNode) bool {
	if node == nil {
		return false
	}
	anchor := t.anchorNode()
	if anchor == nil || !t.containsNode(anchor) {
		anchor = t.pressedNode
	}
	return t.selection.SelectDragUnion(t.visibleSelectionOrder(), anchor, node)
}

func (t *TreeView) selectDragMarquee() bool {
	rect, ok := t.marqueeRect()
	if !ok {
		return false
	}
	return t.selection.ApplyMarquee(t.visibleSelectionOrder(), func(node *TreeNode) bool {
		rowRect, _, ok := t.rowRectForNode(node)
		if !ok || rowRect.Empty() {
			return false
		}
		_, hit := geom.Intersect(rect, rowRect)
		return hit
	}, t.pressedMods&(event.ModCtrl|event.ModShift) != 0)
}

func (t *TreeView) firstSelectedVisible() *TreeNode {
	node, _ := t.selection.firstSelectedInOrder(t.visibleSelectionOrder())
	return node
}

func (t *TreeView) appendAll(out *[]treeEntry, node *TreeNode, depth int) {
	if node == nil {
		return
	}
	*out = append(*out, treeEntry{node: node, depth: depth})
	for _, child := range node.Children {
		t.appendAll(out, child, depth+1)
	}
}

func (t *TreeView) selectedIndex(visible []treeEntry, node *TreeNode) int {
	for i, entry := range visible {
		if entry.node == node {
			return i
		}
	}
	return -1
}

func (t *TreeView) visibleSelectionOrder() selectionOrder[*TreeNode] {
	visible := t.visibleNodes()
	return selectionOrder[*TreeNode]{
		Len:     func() int { return len(visible) },
		ItemAt:  func(index int) *TreeNode { return visible[index].node },
		IndexOf: func(node *TreeNode) int { return t.selectedIndex(visible, node) },
	}
}

func (t *TreeView) anchorNode() *TreeNode {
	node, _ := t.selection.Anchor()
	return node
}

func (t *TreeView) recentNode() *TreeNode {
	node, ok := t.selection.Recent()
	if !ok || node == nil || !t.containsNode(node) {
		return nil
	}
	return node
}

func (t *TreeView) recoveryNodeForKey(visible []treeEntry, key event.Key) *TreeNode {
	if len(visible) == 0 {
		return nil
	}
	switch key {
	case event.KeyHome:
		return visible[0].node
	case event.KeyEnd:
		return visible[len(visible)-1].node
	case event.KeyUp, event.KeyDown, event.KeyLeft, event.KeyRight, event.KeyPageUp, event.KeyPageDown, event.KeySpace:
		if t.selectionBehavior().allowsRecentRecovery() {
			if recent := t.recentNode(); recent != nil {
				return recent
			}
		}
		return visible[0].node
	default:
		return nil
	}
}

func (t *TreeView) selectionBehavior() selectionBehavior {
	return t.selectionOpts.behavior()
}

func (t *TreeView) ensureSelectedVisible() {
	selected := t.SelectedNode()
	if selected == nil {
		return
	}
	visible := t.visibleNodes()
	index := t.selectedIndex(visible, selected)
	if index < 0 {
		return
	}
	if index < t.topIndex {
		t.topIndex = index
	}
	bottom := t.topIndex + t.visibleRows() - 1
	if index > bottom {
		t.topIndex = index - t.visibleRows() + 1
	}
	t.topIndex = clampInt(t.topIndex, 0, t.maxTopIndex(len(visible)))
}

func (t *TreeView) toggleExpanded(node *TreeNode) {
	if node == nil || len(node.Children) == 0 {
		return
	}
	node.Expanded = !node.Expanded
	t.cancelRename()
	if selected := t.SelectedNode(); !node.Expanded && selected != nil && t.isDescendant(node, selected) {
		t.selectOnly(node)
		if t.onChange != nil {
			t.onChange(node)
		}
	}
	if !node.Expanded && t.hotNode != nil && t.hotNode != node && t.isDescendant(node, t.hotNode) {
		t.clearHot()
	}
	t.ensureSelectedVisible()
	t.syncScrollBar(len(t.visibleNodes()))
}

func (t *TreeView) isDescendant(ancestor, node *TreeNode) bool {
	for current := node; current != nil; current = current.parent {
		if current == ancestor {
			return true
		}
	}
	return false
}

func (t *TreeView) expandAncestors(node *TreeNode) {
	for current := node.parent; current != nil; current = current.parent {
		current.Expanded = true
	}
}
