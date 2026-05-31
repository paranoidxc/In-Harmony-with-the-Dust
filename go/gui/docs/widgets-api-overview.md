# Widgets API Overview

## 1. 文档目的

这份文档用于整理当前 `gui/widgets` 中已经实现、且可以视为“对外可用”的控件 API。

目标不是逐行复述代码，而是明确三件事：

1. 哪些类型和方法已经可以作为稳定公共接口使用。
2. 哪些行为是多个控件共享的统一约定。
3. 哪些能力虽然存在实现，但目前仍应视为内部细节，不建议依赖。

这份文档对应 [widget-api-consolidation-plan.md](/Users/xc/projects/In-Harmony-with-the-Dust/go/gui/docs/widget-api-consolidation-plan.md) 的第一个里程碑。

## 2. 范围

当前纳入本总览的内容：

- 控件基础协议
- 共享输入/焦点/文本输入接口
- 选择行为选项
- Menu / MenuBar / CommandID
- HeaderControl
- ListView
- TreeView
- ComboBox
- TabControl
- Toolbar
- StatusBar
- Tooltip / Overlay 相关公共接口

未纳入本轮稳定 API 说明的内容：

- demo 里的 Explorer 业务逻辑
- `desktop` 内部菜单跟踪和 overlay 管理细节
- 各控件私有布局、命中测试、绘制辅助函数

## 3. 控件基础协议

文件参考：

- [control.go](/Users/xc/projects/In-Harmony-with-the-Dust/go/gui/widgets/control.go)

### 3.1 `Control`

所有可交互控件都实现 `Control`：

```go
type Control interface {
    widget.Widget
    Paint(PaintContext) error
    MouseEnter(EventContext)
    MouseLeave(EventContext)
    MouseMove(EventContext, geom.Point)
    MouseDown(EventContext, event.MouseButtonEvent, geom.Point)
    MouseUp(EventContext, event.MouseButtonEvent, geom.Point)
    KeyDown(EventContext, event.KeyEvent) bool
    CanFocus() bool
    SetFocused(bool)
    Focused() bool
}
```

可以把它理解为当前 retained-mode 控件树里的最小交互协议。

### 3.2 可选扩展接口

这些接口不是所有控件都必须实现，但一旦实现，desktop 会自动接入对应能力：

- `WheelHandler`
- `FocusHandler`
- `TickHandler`
- `TextInputHandler`
- `TooltipProvider`

它们分别对应：

- 鼠标滚轮
- 焦点进入/离开通知
- 定时器/闪烁/自动滚动
- 文本输入与 IME
- tooltip 查询

### 3.3 `EventContext`

控件通过 `EventContext` 与宿主环境交互。

当前稳定可依赖的能力：

- `Invalidate`
- `SetFocus`
- `Capture`
- `ReleaseCapture`
- `DispatchCommand`
- `ShowContextMenu`
- `ClipboardText`
- `SetClipboardText`
- `MeasureText`
- `LineHeight`

这里的定位应当是“控件运行时能力上下文”，而不是业务对象。

## 4. Overlay / Tooltip 公共接口

文件参考：

- [control.go](/Users/xc/projects/In-Harmony-with-the-Dust/go/gui/widgets/control.go)

### 4.1 `OverlayRequest`

当前已经暴露出 popup/overlay 的最小公共结构：

```go
type OverlayRequest struct {
    Owner          Control
    Content        Control
    Anchor         geom.Rect
    Placement      OverlayPlacement
    CloseOnOutside bool
    OnClose        func()
}
```

它已经足够承载：

- ComboBox 下拉层
- PopupMenu
- Tooltip overlay

但当前仍偏“已暴露、待正式化”的状态。后续计划会把它整理为 `PopupWindow` 能力。

### 4.2 `TooltipProvider`

Tooltip 目前通过查询式协议工作：

```go
type TooltipProvider interface {
    TooltipAt(local geom.Point, measure func(string) geom.Size) TooltipInfo
}
```

这意味着 tooltip 不是独立子控件，而是由宿主在 hover 时向控件查询。

这是当前推荐的使用方式。

## 5. 选择行为公共接口

文件参考：

- [selection_options.go](/Users/xc/projects/In-Harmony-with-the-Dust/go/gui/widgets/selection_options.go)

### 5.1 `SelectionBehaviorOptions`

当前共享选择策略的正式输入结构是：

```go
type SelectionBehaviorOptions struct {
    MultiSelect       bool
    RecoverFromRecent bool
    BlankDragSelect   bool
}
```

默认值：

- `MultiSelect: true`
- `RecoverFromRecent: true`
- `BlankDragSelect: true`

### 5.2 当前共享这套选择行为的控件

- `ListBox`
- `TreeView`
- `ListView`

### 5.3 当前可视为统一约定的行为

只要控件使用这套选择策略，就应当遵循这些交互语义：

- `Ctrl+A` 全选
- `Ctrl+Space` 切换 lead 项
- `Shift` 扩选
- `Ctrl` 切换选中状态
- 空白区点击清选
- 空白区拖拽可触发 marquee
- 拖拽时自动滚动

### 5.4 当前不暴露的部分

以下仍视为内部实现，不建议在外部直接依赖：

- `selectionModel`
- `selectionOrder`
- 内部 drag base / recent / anchor 存储细节

## 6. Menu System

文件参考：

- [menu.go](/Users/xc/projects/In-Harmony-with-the-Dust/go/gui/widgets/menu.go)

### 6.1 核心类型

当前 menu system 的稳定公共类型：

- `CommandID`
- `Accelerator`
- `MenuItem`
- `Menu`
- `MenuBar`

### 6.2 推荐使用方式

构造菜单项：

```go
widgets.NewMenuItem(id, text, shortcut)
widgets.NewSubmenuItem(text, submenu)
widgets.NewSeparator()
widgets.NewMenu(items...)
widgets.NewMenuBar(items...)
```

### 6.3 当前可视为稳定的 menu item 字段

可以直接配置：

- `ID`
- `Text`
- `Enabled`
- `Checked`
- `Shortcut`
- `Submenu`

### 6.4 当前建议依赖的能力

- mnemonic 查找
- accelerator 匹配
- `Checked` 的展示语义
- submenu 层级结构
- context menu 通过 `SetContextMenu` / `SetContextMenuProvider` 接入

### 6.5 当前不建议跨层依赖的部分

下面这些仍应视为 `desktop` 内部菜单跟踪实现：

- 菜单链激活状态
- popup 栈维护
- 菜单键盘跟踪细节
- 菜单 overlay 生命周期

后续会以 `PopupWindow + MenuTracker` 的角度再正式化。

## 7. HeaderControl

文件参考：

- [headercontrol.go](/Users/xc/projects/In-Harmony-with-the-Dust/go/gui/widgets/headercontrol.go)

### 7.1 公共类型

- `HeaderColumn`
- `HeaderAlign`
- `HeaderControl`

### 7.2 当前稳定 API

- `SetColumns(columns ...HeaderColumn)`
- `Columns() []HeaderColumn`
- `SetSortIndicator(index int, descending bool)`
- `OnColumnClick(func(int))`
- `OnColumnResize(func(int, int))`
- `OnColumnAutoFit(func(EventContext, int))`
- `ColumnWidth(index int) int`

### 7.3 当前语义

`HeaderControl` 是一个独立控件，不只是 `ListView` 的私有结构。

它当前支持：

- hover 高亮
- 按压态
- 列头点击
- divider 拖拽调宽
- divider 双击自动适宽
- 排序方向指示

## 8. ListView

文件参考：

- [listview.go](/Users/xc/projects/In-Harmony-with-the-Dust/go/gui/widgets/listview.go)

### 8.1 公共类型

- `ListViewColumn`
- `ListViewItem`
- `ListViewContextMenuInfo`
- `ListViewSelectionOptions`
- `ListView`

### 8.2 当前稳定数据与选择 API

- `SetColumns(columns ...ListViewColumn)`
- `Columns() []ListViewColumn`
- `SetItems(items []ListViewItem)`
- `Items() []ListViewItem`
- `SelectedIndex() int`
- `SelectedIndices() []int`
- `SetSelectedIndex(index int)`
- `SetSelectedIndexSilent(index int)`
- `SetSelectionOptions(options ListViewSelectionOptions)`
- `SelectionOptions() ListViewSelectionOptions`
- `SetMultiSelect(enabled bool)`

### 8.3 当前稳定交互回调

- `OnChange(func(int, ListViewItem))`
- `OnActivate(func(int, ListViewItem))`
- `OnColumnClick(func(int))`
- `OnRenameRequest(func(EventContext, int, ListViewItem) bool)`
- `OnRenameCommit(func(int, ListViewItem, string, string))`

### 8.4 当前稳定菜单与排序 API

- `SetSortIndicator(index int, descending bool)`
- `SetContextMenu(menu *Menu)`
- `SetContextMenuProvider(func(ListViewContextMenuInfo) *Menu)`

### 8.5 当前稳定重命名 API

- `BeginRename(index int) bool`

当前语义是：

- `BeginRename` 只负责进入 inline rename 模式
- 重命名结果通过 `OnRenameCommit` 回调给外部
- 外部数据模型是否同步，由调用方负责

### 8.6 当前可以明确承诺的交互能力

- details 风格多列显示
- 列头点击/排序指示
- 列宽拖拽与自动适宽
- 单击选择
- 双击激活
- Enter 激活
- F2 重命名
- 多选
- marquee
- 自动滚动
- 右键 context menu

### 8.7 当前刻意不稳定的边界

下面这些虽然实现存在，但还不建议作为“正式 API”承诺：

- 单元格级编辑
- owner-data / virtual list
- 自定义单元格绘制
- model/provider 抽象

## 9. TreeView

文件参考：

- [treeview.go](/Users/xc/projects/In-Harmony-with-the-Dust/go/gui/widgets/treeview.go)

### 9.1 公共类型

- `TreeNode`
- `TreeNodeKind`
- `TreeViewContextMenuInfo`
- `TreeView`

### 9.2 当前稳定数据 API

- `NewTreeNode`
- `NewFolderNode`
- `NewFileNode`
- `(*TreeNode).SetChildren`
- `(*TreeNode).AddChild`
- `(*TreeNode).Parent()`
- `(*TreeNode).EffectiveKind()`
- `SetRoots(roots ...*TreeNode)`
- `Roots() []*TreeNode`

### 9.3 当前稳定选择与交互 API

- `SelectedNode() *TreeNode`
- `SelectedNodes() []*TreeNode`
- `SetSelectedNode(node *TreeNode) bool`
- `SelectionOptions() TreeViewSelectionOptions`
- `SetSelectionOptions(options TreeViewSelectionOptions)`
- `SetMultiSelect(enabled bool)`

### 9.4 当前稳定回调

- `OnChange(func(*TreeNode))`
- `OnActivate(func(*TreeNode))`
- `OnBeginRename(func(*TreeNode))`
- `OnRenameCommit(func(*TreeNode, string, string))`

### 9.5 当前稳定菜单与重命名 API

- `SetContextMenu(menu *Menu)`
- `SetContextMenuProvider(func(TreeViewContextMenuInfo) *Menu)`
- `BeginRename(node *TreeNode) bool`

### 9.6 当前可以明确承诺的交互能力

- 展开/折叠
- 键盘导航
- 多选
- hover/expander 热区
- 右键 context menu
- 慢单击延迟重命名
- F2 重命名
- inline rename
- 文件名主体默认选区
- marquee
- 自动滚动

### 9.7 当前刻意不稳定的边界

- 懒加载/虚拟节点模型
- 图标提供器接口
- 多列 tree
- owner-draw tree

## 10. ComboBox

文件参考：

- [combobox.go](/Users/xc/projects/In-Harmony-with-the-Dust/go/gui/widgets/combobox.go)

### 10.1 当前稳定类型

- `ComboBox`

### 10.2 当前稳定 API

- `SetItems(items []string)`
- `Items() []string`
- `SelectedIndex() int`
- `Text() string`
- `SetSelectedIndex(index int) bool`
- `SetSelectedIndexSilent(index int) bool`
- `OnChange(func(int, string))`
- `OnCommit(func(int, string))`
- `SetTooltip(text string)`
- `SetEditable(editable bool)`
- `Editable() bool`
- `SetText(text string)`

### 10.3 当前已承诺的行为

- 非编辑模式选择
- 可编辑模式输入
- dropdown popup
- 键盘上下移动
- Enter/Tab 提交
- ESC 恢复已提交值
- 前缀匹配与候选尾巴显示
- tooltip 查询

### 10.4 当前边界

当前还没有承诺：

- 自定义 item renderer
- 非字符串数据模型
- 独立 popup item style 注入

## 11. TabControl

文件参考：

- [tabcontrol.go](/Users/xc/projects/In-Harmony-with-the-Dust/go/gui/widgets/tabcontrol.go)

### 11.1 当前稳定类型

- `TabPage`
- `TabControl`

### 11.2 当前稳定 API

- `NewTabPage(title string, content Control) *TabPage`
- `SetPages(pages ...*TabPage)`
- `AddPage(page *TabPage)`
- `Pages() []*TabPage`
- `SelectedIndex() int`
- `SelectedPage() *TabPage`
- `SetSelected(index int) bool`
- `OnSelectionChange(func(int, *TabPage))`

### 11.3 当前已承诺的行为

- 页签切换
- 禁用页跳过
- 键盘左右导航
- Home/End
- 选中页内容显示

## 12. Toolbar

文件参考：

- [toolbar.go](/Users/xc/projects/In-Harmony-with-the-Dust/go/gui/widgets/toolbar.go)

### 12.1 当前稳定类型

- `ToolbarItem`
- `Toolbar`

### 12.2 当前稳定 API

- `SetItems(items ...*ToolbarItem)`
- `Items() []*ToolbarItem`
- `SetChecked(id CommandID, checked bool)`
- `SetDisabled(id CommandID, disabled bool)`
- `OnCommand(func(CommandID))`

### 12.3 当前已承诺的行为

- 按钮/分隔项混排
- checked / disabled 展示
- tooltip 查询
- 命令分发

### 12.4 当前边界

当前 toolbar 仍偏轻量，尚未承诺：

- overflow
- 可换行
- 可拖拽重排
- 文本+图标混合布局策略注入

## 13. StatusBar

文件参考：

- [statusbar.go](/Users/xc/projects/In-Harmony-with-the-Dust/go/gui/widgets/statusbar.go)

### 13.1 当前稳定类型

- `StatusPane`
- `StatusBar`

### 13.2 当前稳定 API

- `SetText(text string)`
- `SetPanes(panes []StatusPane)`
- `SetPaneText(index int, text string)`
- `Panes() []StatusPane`

### 13.3 当前已承诺的行为

- 单文本模式
- 多 pane 分栏模式
- 经典凹陷式状态栏绘制

## 14. 当前共享交互约定

这部分是多个控件目前已经共同遵循的约定，后续不应轻易分叉。

### 14.1 焦点

- `CanFocus` 决定控件是否参与焦点链
- `SetFocus` 由 `EventContext` 发起
- `FocusGained` / `FocusLost` 只在实现 `FocusHandler` 时触发

### 14.2 文本输入

支持 inline edit 的控件通过 `TextInputHandler` 接入 IME 与文本输入：

- `Edit`
- `TreeView` 的 rename editor
- `ListView` 的 rename editor
- `ComboBox` 的 editable editor

### 14.3 Context Menu

当前通用模式是二选一：

- `SetContextMenu(menu *Menu)`
- `SetContextMenuProvider(func(...) *Menu)`

推荐优先使用 provider，因为它能根据当前选择动态返回菜单。

### 14.4 Rename

当前 `TreeView` 与 `ListView` 已形成统一 rename 约定：

- `BeginRename(...)` 进入编辑态
- `F2` 请求重命名
- `Enter` 提交
- `Escape` 取消
- `FocusLost` 提交
- 文件名默认只选主体，不选扩展名

## 15. 当前明确不属于稳定 API 的部分

下面这些虽然代码里存在，但目前不应被外部直接依赖：

- 各控件私有 `hitTest` / `layout` / `rowRect` / `cellRect` 辅助函数
- `desktop` 内部菜单 popup 栈
- overlay 内部状态结构
- 选择模型内部状态结构
- demo 中的上下文菜单构建逻辑

## 16. 下一步建议

按当前文档收口后，下一步不再继续扩 demo 业务，而是优先做：

1. `PopupWindow` 正式化
2. `DialogWindow` 正式化

其中第一优先级仍然是 `PopupWindow`，因为它能把当前：

- Menu
- ComboBox dropdown
- Tooltip

统一提升为真正的框架能力，而不只是已实现的内部机制。
