# Popup / Overlay Architecture

## 1. 文档目的

这份文档整理当前 GUI 框架里与 popup / overlay 相关的公共能力、已有实现和结构性问题。

它不是最终设计，而是为下一步 `PopupWindow` 正式化做现状收口。

目标：

1. 说明当前哪些能力已经存在并可复用。
2. 说明 Menu / ComboBox / Tooltip 分别是怎么落在这套机制上的。
3. 指出现在的抽象断层，作为后续 `PopupWindow` 抽象的输入。

## 2. 相关代码位置

核心文件：

- [control.go](/Users/xc/projects/In-Harmony-with-the-Dust/go/gui/widgets/control.go)
- [desktop.go](/Users/xc/projects/In-Harmony-with-the-Dust/go/gui/desktop/desktop.go)
- [overlay.go](/Users/xc/projects/In-Harmony-with-the-Dust/go/gui/desktop/overlay.go)
- [menu.go](/Users/xc/projects/In-Harmony-with-the-Dust/go/gui/desktop/menu.go)
- [tooltip.go](/Users/xc/projects/In-Harmony-with-the-Dust/go/gui/desktop/tooltip.go)
- [combobox.go](/Users/xc/projects/In-Harmony-with-the-Dust/go/gui/widgets/combobox.go)

## 3. 当前公共接口

### 3.1 `OverlayRequest`

当前 widgets 层已经暴露出一个通用 popup 请求结构：

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

它表达了当前 popup 的最小公共能力：

- 有明确 `Owner`
- 有独立 `Content`
- 依赖 `Anchor`
- 支持位置策略
- 支持点击外部关闭
- 支持关闭回调

### 3.2 `OverlayContext`

当前控件如果想显示自己的 popup，需要运行时 `EventContext` 同时实现：

```go
type OverlayContext interface {
    ShowOverlay(OverlayRequest) bool
    HideOverlay(Control) bool
    OverlayVisible(Control) bool
}
```

当前 `desktop.controlContext` 实现了这组接口，供普通控件使用。

## 4. 当前 overlay 运行模型

### 4.1 普通控件 overlay

普通控件 popup 当前通过 `Desktop.showControlOverlay(...)` 建立。

内部状态是：

- `controlOverlayState`
- `ownerWindow`
- `owner`
- `content`
- `rect`
- `closeOnOutside`
- `onClose`
- overlay 内部 hover/focus/capture 状态

这说明当前框架实际上已经有“控件私有 popup 层”的运行模型。

### 4.2 overlay 的能力

当前 control overlay 已具备：

- 独立绘制
- 独立 hit test
- 独立 hover
- 独立 focus
- 独立 capture
- 关闭后回调

这已经非常接近一个真正的 popup host。

### 4.3 位置策略

当前位置策略由：

- `OverlayBelowStart`
- `OverlayRightTop`

和 `Desktop.fitOverlayOrigin(...)` 共同决定。

当前策略已经支持：

- 下拉式 popup
- 右侧子菜单式 popup
- 屏幕边界内回退调整

## 5. 当前三类 popup 的实际落地方式

### 5.1 ComboBox dropdown

文件参考：

- [combobox.go](/Users/xc/projects/In-Harmony-with-the-Dust/go/gui/widgets/combobox.go)

当前 ComboBox 是最完整地使用通用 overlay 接口的控件。

行为链路：

1. `ComboBox.openPopup(ctx)`
2. 断言 `ctx` 实现 `OverlayContext`
3. 动态创建一个 `ListBox` 作为 popup 内容
4. 通过 `OverlayRequest` 请求显示
5. `OnClose` 中回收 `dropped / pressed / popup` 状态

这说明：

- 当前 `OverlayRequest` 已经足以支持“控件托管型 popup”
- popup 内容可以是任意 `Control`
- popup 关闭时需要一个正式的宿主生命周期回调

### 5.2 ContextMenu / PopupMenu

文件参考：

- [menu.go](/Users/xc/projects/In-Harmony-with-the-Dust/go/gui/desktop/menu.go)

菜单当前并没有走 `OverlayRequest` 的普通控件 popup 通道，而是走了 `desktop` 自己维护的菜单系统：

- `menuMode`
- `menuWindow`
- `menuPopups`
- `popupMenuState`

菜单 popup 自己是一个 `desktopOverlay`，并且自行维护：

- popup rect
- item layout
- selected index
- submenu 链

这套机制能正常工作，但它和通用 `controlOverlayState` 是两套并行体系。

### 5.3 Tooltip

文件参考：

- [tooltip.go](/Users/xc/projects/In-Harmony-with-the-Dust/go/gui/desktop/tooltip.go)

Tooltip 当前也没有走 `OverlayRequest`。

它的模型是：

1. 控件实现 `TooltipProvider`
2. desktop 在 hover 时调用 `TooltipAt(...)`
3. 记录 `tooltipText / tooltipAnchor / tooltipDue`
4. 时间到后构造 `tooltipOverlayState`
5. 作为一个 `desktopOverlay` 加入 overlay 栈

Tooltip 当前是“查询式 + desktop 自绘 overlay”模型。

## 6. 当前架构的真实状态

从运行效果看，当前框架已经具备 popup 能力。

但从抽象层次看，其实存在三套不同模型：

1. 通用控件 overlay
2. 菜单专用 popup 栈
3. tooltip 专用 overlay

也就是说：

- “显示在主控件之外的一块临时内容” 这个能力已经存在
- 但它还没有被统一命名成一个正式框架概念

## 7. 当前已经稳定的公共约定

下面这些可以视为现阶段已经成立的公共约定。

### 7.1 popup 总是依附某个 owner

无论是 ComboBox、ContextMenu 还是 Tooltip，popup 都不是自由浮动对象。

它们都依附于：

- 一个触发源控件
- 或一个所属 window

### 7.2 popup 都依赖 anchor rect

当前三类 popup 都以某个 anchor 作为定位起点：

- ComboBox：`LocalRect(c)`
- ContextMenu：右键点或菜单项 rect
- Tooltip：`TooltipInfo.Anchor`

### 7.3 popup 都需要边界内定位修正

当前三类 popup 都依赖同一个几何事实：

- 初始位置不一定可用
- 必须做屏幕边界回退

`fitOverlayOrigin(...)` 已经承担了这部分共性。

### 7.4 popup 都需要自己的一套输入焦点边界

popup 一旦出现，就不能再完全走普通控件树输入链。

当前已经体现为：

- control overlay 有独立 focus/capture
- menu popup 有独立 menu routing
- tooltip 则明确不接输入

这说明后续 `PopupWindow` 必须支持不同交互等级：

- interactive popup
- passive popup

## 8. 当前的主要问题

### 8.1 概念重复

现在代码里同时有：

- `desktopOverlay`
- `controlOverlayState`
- `popupMenuState`
- `tooltipOverlayState`

它们都在表达“临时浮在正常窗口之上的一块内容”，只是交互模型不同。

### 8.2 菜单没有复用通用 overlay 宿主

Menu popup 目前有完整独立体系，这在功能上没问题，但会让后续：

- 焦点规则
- 键盘路由
- 生命周期
- outside click 规则

继续和普通 popup 分裂。

### 8.3 Tooltip 是特例，不是 popup 子类

Tooltip 当前更像 desktop 的一项特权逻辑，而不是一个被统一承载的 popup 类型。

这会让后续：

- richer tooltip
- delayed popup hint
- validation balloon

难以沿现有结构扩展。

### 8.4 `OverlayContext` 已经可用，但语义还不完整

当前 `ShowOverlay / HideOverlay / OverlayVisible` 已经够 ComboBox 用，但还不够表达：

- popup 类型
- 是否可聚焦
- 键盘优先级
- 是否独占输入
- owner popup / child popup 关系

## 9. `PopupWindow` 抽象应该吸收什么

下一步正式化时，建议保留并提升下面这些已有成果。

### 9.1 要保留的已有结构

- `Owner`
- `Content`
- `Anchor`
- `Placement`
- `CloseOnOutside`
- `OnClose`
- 屏幕边界内定位修正
- popup 独立 focus/capture 能力

### 9.2 要统一收口的三类对象

- ComboBox dropdown
- ContextMenu / PopupMenu
- Tooltip

### 9.3 要新增的正式语义

`PopupWindow` 至少应当明确下面这些概念：

- 是否可聚焦
- 是否接收键盘
- 是否参与 hover routing
- 是否允许子 popup
- 是否是被动提示型 popup
- 与 owner 的关闭联动关系

## 10. 当前建议的结论

现阶段不应该继续把 popup 当作“desktop 内部技巧”。

当前更准确的判断是：

- popup 能力已经存在
- popup API 已经部分暴露
- popup 运行模型已经被 ComboBox 验证
- menu / tooltip 还没有统一进这套抽象

所以接下来的正确方向不是继续补 demo，而是：

1. 以 `OverlayRequest` 为基础
2. 把它升级成正式的 `PopupWindow` / popup host 抽象
3. 再把 menu / tooltip 逐步并回统一体系

## 11. 与下一步开发的关系

这份文档对应计划里的：

- “popup / menu / tooltip / combobox 的共享承载关系”

在它之后，下一步代码任务应当是：

### 下一任务

- 定义 `PopupWindow` 或等价结构
- 收敛 `desktopOverlay` 与 `controlOverlayState` 的职责边界
- 先让普通控件 popup 与 menu popup 共享同一套宿主生命周期模型

不建议下一步直接做：

- 更多 demo 菜单项
- 更多业务对话框
- 更复杂 tooltip 内容

这些都应该等 `PopupWindow` 抽象稳定后再做。
