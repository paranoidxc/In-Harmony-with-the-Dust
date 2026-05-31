# GUI 控件/API 沉淀与底层能力计划

## 1. 目的

这份计划用于把当前已经在 demo 中验证过的 GUI 能力，重新收束为可复用的控件/API，而不是继续扩展 Explorer 风格的业务语义。

目标分两部分：

1. 把已经做出来且已验证有效的能力，沉淀成稳定的通用控件接口与文档。
2. 后续优先投入更底层、更可复用的 GUI 框架能力，而不是继续堆 demo 交互。

## 2. 当前判断

### 2.1 Phase 4 控件状态

`windows-classic-ui-design.md` 里 Phase 4 提到的增强控件：

- ComboBox
- StatusBar
- TabControl
- TreeView
- Toolbar
- Tooltip

这些控件目前都已经有实现，且已经通过一轮集成式验证。

### 2.2 已经验证过、值得沉淀的能力

以下内容已经不应再视为“demo 技巧”，而应视为控件库能力：

- 统一 popup/overlay 承载层
- MenuBar + PopupMenu + ContextMenu 链路
- Toolbar / StatusBar / Tooltip 的公共宿主能力
- ComboBox 的可编辑模式与下拉行为
- TabControl 的页签切换与键盘导航
- TreeView 的：
  - 多选
  - 空白区点击行为
  - 拖选阈值
  - 自动滚动
  - marquee
  - inline rename
  - F2 / 延迟重命名
- ListBox 的多选与统一选择策略
- HeaderControl + ListView(details) 的：
  - 列头点击
  - 排序指示
  - 列宽拖拽
  - 自动适宽
  - inline rename
  - context menu provider
- 统一 selection framework

## 3. 现在要停止继续做的事

下面这些内容不再作为当前主线：

- 继续给 demo 加资源管理器业务语义
- 继续堆“文件操作”型菜单项
- 在 demo 层补更多像 Explorer 的特殊逻辑
- 把 demo 变成应用，而不是把控件变成框架

demo 后续只保留两个用途：

- 控件回归验证
- 交互截图/手工冒烟入口

## 4. 第一阶段：把现有能力沉淀为通用控件/API

### 4.1 统一整理的对象

优先整理以下控件与基础层：

1. `SelectionModel / SelectionBehaviorOptions`
2. `Overlay / Popup` 承载接口
3. `Menu system`
4. `HeaderControl`
5. `ListView`
6. `TreeView`
7. `ComboBox`
8. `Tooltip`

### 4.2 需要沉淀的 API 方向

#### A. 选择框架

把目前 ListBox / TreeView / ListView 共用的选择逻辑，明确成正式能力：

- 单选/多选策略
- Shift 扩选
- Ctrl 切换
- Ctrl+A
- Ctrl+Space
- 空白区点击清选
- drag threshold
- auto-scroll
- marquee

文档上要明确：

- 哪些行为是默认启用
- 哪些由 `SelectionBehaviorOptions` 控制
- 哪些控件共享这套行为

#### B. Popup/Overlay

把当前“可用但偏内部”的 overlay 能力，提升为稳定抽象：

- popup 的宿主关系
- anchor rect
- placement
- outside click close
- popup focus
- popup keyboard routing

目标不是继续写菜单，而是把“菜单/下拉/提示框共享的承载机制”整理成正式 API。

#### C. ListView / TreeView 的通用交互 API

需要从“demo 驱动”回收到“控件自身 API”：

- `OnChange`
- `OnActivate`
- `OnRenameRequest`
- `OnRenameCommit`
- `SetContextMenu`
- `SetContextMenuProvider`
- 排序指示 API
- 列宽/自动适宽相关 API

同时补文档说明：

- 哪些 API 只负责 UI 事件
- 哪些 API 会直接修改控件内部状态
- 哪些 API 需要由外部数据模型同步

#### D. 数据模型边界

当前 TreeView / ListView 还是偏“直接塞数据切片”的风格。

下一步要明确是否演进为以下两种模式之一：

1. 保持轻量：继续使用简单 item slice / tree node
2. 正式抽象：引入 model/provider 接口

本轮先不强推重构，但文档中要先把方向写清楚。

### 4.3 文档产出目标

至少补齐下面三类文档：

1. 控件对外 API 概览
2. 共享交互约定
3. popup / menu / tooltip / combobox 的共享承载关系

## 5. 第二阶段：下一步改做更底层、可复用的 GUI 能力

下一步不再优先做 demo 业务，而优先做下面两个底层能力。

### 5.1 Priority A：PopupWindow 正式化

这项应作为下一主线。

原因：

- 文档里已经明确有 `PopupWindow`
- 现在 menu / combo dropdown / tooltip 已经都依赖一套 popup/overlay 机制
- 说明抽象基础已经存在，只是还没有正式命名和稳定 API

需要完成的事：

- 定义 `PopupWindow` 或等价抽象
- 统一 popup host 生命周期
- 统一 popup focus / capture / keyboard routing
- 明确 popup 与普通 window 的边界
- 让 Menu / ComboBox / Tooltip 走同一套正式承载协议

完成后收益：

- 现有 overlay 从“实现细节”升级为“框架能力”
- 后续 Dialog / ContextMenu / AutoComplete / DatePicker 都有统一基座

### 5.2 Priority B：DialogWindow 正式化

这项作为 PopupWindow 之后的下一个目标。

需要完成的事：

- 模态/非模态对话框壳
- 默认按钮/取消按钮
- 焦点进入与返回
- ESC / Enter 的标准行为
- 所属窗口关系
- 模态阻断规则

完成后收益：

- 后面所有属性框、确认框、输入框都能用统一窗口壳
- 不再需要在 demo 里手拼临时面板

## 6. 建议执行顺序

### Milestone 1：文档与 API 收口

- 梳理当前通用能力
- 明确哪些属于正式控件 API
- 写清楚哪些属于 demo 逻辑，不再继续扩

当前进度：

- 已完成：`widgets-api-overview.md`
- 已完成：`popup-overlay-architecture.md`
- 下一步：开始 `PopupWindow` 抽象

完成标准：

- `gui/docs` 中有一份明确的控件/API 收口文档
- 后续开发任务以框架能力命名，而不是以 demo 行为命名

### Milestone 2：PopupWindow 抽象

- 从 overlay 机制中提炼正式 popup API
- 菜单/下拉/提示框统一走这套协议

当前进度：

- 已开始：新增 `PopupRequest / PopupContext / PopupKind`
- 已接入：`ComboBox` 走新 popup 抽象兼容层
- 已接入：menu popup 共享 popup host 生命周期基础结构
- 已开始：抽取 popup host 通用命中与 owner 边界 helper
- 已开始：用 popup host 可见性统一 tooltip / interactive popup 顶层判断
- 已开始：用共享 helper 收敛 outside-click / dismiss 的几何语义
- 已开始：把 Desktop 顶层 mouse/key 路由切到统一 popup dispatch 入口
- 已完成：Desktop 鼠标路由优先按最上层 popup host 类型分发，而不是固定 menu 优先
- 已完成：Desktop 键盘路由在保留 Alt/menu 特例外，普通按键优先交给最上层 popup host
- 已完成：overlay focus 接入 Desktop 的 text input / IME rect / tick 链路
- 已开始：`Escape` 的 dismiss 语义按最上层 interactive popup host 优先收敛
- 已完成：menu 接管输入时，清理底层 overlay 的 focus / hover / capture 残留状态
- 已完成：interactive control popup 打开时，显式清理底层 window capture
- 下一步：继续收敛 menu popup 与普通 popup 的显式焦点模型与统一 input owner helper

完成标准：

- Popup 不再只是 desktop 内部技巧
- 对外可解释、可测试、可复用

### Milestone 3：DialogWindow 抽象

- 模态/非模态对话框壳
- 标准按钮与键盘行为

完成标准：

- 可用一个最小对话框 demo 验证
- 不依赖 Explorer 场景

## 7. API 沉淀时的硬约束

后续收口过程中遵守下面几条：

1. 不把 Explorer 业务词汇固化进通用控件 API。
2. 不让 demo 回调反过来定义控件边界。
3. 通用控件只暴露“交互事件与状态”，不内建业务语义。
4. popup / dialog / selection 这类横切能力优先抽象到底层。
5. 每沉淀一个能力，都要补最小测试与最小文档。

## 8. 非目标

这份计划明确不包含：

- 文件系统集成
- 真正的 Explorer 应用逻辑
- 拖拽文件、复制粘贴文件、属性页等业务功能
- 更多 demo 菜单项堆叠

## 9. 计划后的立即下一步

按这个计划，下一开发项建议是：

1. 先补一份“现有控件对外 API 总览”文档。
2. 然后开始做 `PopupWindow` 正式化，而不是继续改 demo。

如果只选一个立刻开工的技术任务：

- 选 `PopupWindow`。

它最能把现有成果从“已实现”变成“框架能力”。
