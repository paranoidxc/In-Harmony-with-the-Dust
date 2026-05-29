# Go + SDL Windows Classic 控件库技术设计文档

## 1. 文档目的

本文档给出一套基于 Go + SDL 的经典桌面 GUI 控件库设计方案，目标视觉风格接近 Windows 95/98/2000 时代的 `Windows Classic`。重点不是做一个现代扁平 UI，而是做一套可复用、可扩展、可逐步落地的控件系统，能够稳定表现以下视觉和交互特征：

- 明确的 1px/2px 立体边框
- 系统色板驱动的灰色界面
- 经典标题栏、菜单栏、按钮、滚动条、输入框、列表框
- 整数像素布局与非抗锯齿感明显的文本/图标
- 接近 Win32 的焦点、捕获、菜单、加速键、通知消息模型

本文档面向“从零开始实现控件库”的场景，因此会同时覆盖：

- 技术选型
- 模块分层
- 渲染与事件系统
- 主题系统
- 核心控件设计
- 测试策略
- 分阶段实施路线

## 2. 目标与非目标

### 2.1 目标

1. 在单个 SDL 宿主窗口中实现一套 retained-mode 控件树。
2. 提供统一的窗口、控件、菜单、对话框、文本输入和滚动支持。
3. 风格上以 `Windows Classic` 为默认主题，不依赖宿主操作系统原生控件。
4. 保证像素级稳定渲染，避免 SDL 渲染路径带来的抗锯齿和平台差异。
5. 公共 API 保持简单，内部实现允许使用类似 Win32 的消息机制。
6. 对中文文本和 IME 输入友好。
7. 后续可以增加主题包、额外控件和多宿主窗口支持。

### 2.2 非目标

1. 第一阶段不追求原生 OS 级窗口管理器集成。
2. 第一阶段不实现复杂富文本、浏览器布局、矢量动画。
3. 第一阶段不做无障碍读屏兼容。
4. 第一阶段不以 GPU 特效为目标，重点是稳定和还原经典视觉。
5. 第一阶段不追求完全复刻全部 Win32 控件行为，只保留最关键的交互习惯。

## 3. 总体结论

这套控件库建议采用以下关键设计：

1. `单宿主 SDL_Window + 内部 Desktop/WindowManager` 模式。
2. `retained-mode` 控件树，而不是 immediate-mode。
3. `软件栅格化画布 + SDL 纹理呈现` 的渲染路径。
4. `内部消息分发 + 对外回调/事件` 的混合编程模型。
5. `主题 token + 位图资源 + 字体配置` 的皮肤系统。
6. `整数像素坐标 + 可选整数倍缩放`，避免经典风格在高 DPI 下失真。

这几个选择组合起来，能在实现复杂度、视觉一致性和可维护性之间取得比较好的平衡。

## 4. 技术选型

### 4.1 语言与基础库

- Go
- SDL2 绑定，建议使用 `github.com/veandco/go-sdl2`
  - `sdl`
  - `ttf`
  - 可选 `img`

### 4.2 为什么选 SDL

SDL 适合作为宿主平台抽象层，而不是 UI 逻辑层。它负责：

- 创建窗口
- 获取输入事件
- 创建纹理并把最终帧提交到屏幕
- 提供文本输入、剪贴板、鼠标光标、定时基础能力

SDL 不负责：

- 控件树
- 焦点管理
- 布局
- 菜单系统
- 经典主题绘制

这部分全部由本库自己控制，才能保证风格统一。

### 4.3 为什么不用纯 SDL_Renderer 直接画控件

直接调用 `SDL_RenderDrawLine/FillRect` 可以做出界面，但会遇到几个问题：

- 各平台像素对齐细节不完全一致
- 文本与图标组合时很难严格做经典风格复刻
- 后续做 dirty rect、截图测试、主题化会比较散

因此更建议：

1. 内部维护一块逻辑帧缓冲区。
2. 所有控件绘制都落到这块软件画布。
3. 每帧只把变化区域更新到 SDL 纹理。
4. 再由 SDL 负责 present。

这样渲染结果更可控，也更适合做像素级回归测试。

## 5. 核心架构

### 5.1 分层结构

```text
App
  -> Backend(SDL)
  -> Desktop
      -> WindowManager
          -> TopLevelWindow
              -> Widget Tree
                  -> Layout
                  -> Theme
                  -> Text
                  -> Paint
```

### 5.2 推荐包结构

```text
/classicui
  /app
  /backend/sdl2
  /desktop
  /event
  /geom
  /invalidate
  /layout
  /paint
  /text
  /theme
  /widget
  /widgets
  /internal/msg
  /internal/state
/cmd/demo
/docs
```

职责建议如下：

- `app`：应用生命周期、主循环、任务投递、定时器
- `backend/sdl2`：SDL 初始化、窗口、纹理、输入事件翻译
- `desktop`：桌面根节点、窗口层级、模态栈、弹出层
- `event`：统一事件类型、键码、鼠标按钮、修饰键
- `geom`：`Point/Rect/Size/Insets/Region`
- `invalidate`：脏区合并与裁剪
- `layout`：布局器与测量协议
- `paint`：软件画布、位图绘制、边框/标题栏等 primitive
- `text`：字体、glyph cache、文本测量、IME 组合态
- `theme`：主题 token、资源、颜色、度量
- `widget`：基类、焦点、可见性、命中测试、消息收发
- `widgets`：Button/Edit/ListBox/MenuBar 等具体控件
- `internal/msg`：内部消息 ID 与路由

## 6. 宿主窗口与内部窗口模型

### 6.1 第一阶段采用单宿主窗口

第一阶段建议只创建一个真正的 SDL 窗口，例如：

- 宿主 SDL 窗口大小：`1280x960`
- 逻辑 UI 尺寸：`640x480`
- 呈现缩放：`2x`

然后在宿主窗口内部自行管理多个“经典窗口”。

好处：

- 标题栏、边框、菜单、激活状态全部可控
- 不需要和不同平台的原生窗口装饰打架
- 拖拽、层级、模态、弹出菜单都更容易统一
- 视觉还原度更高

### 6.2 内部窗口种类

1. `MainWindow`：应用主窗口
2. `DialogWindow`：模态/非模态对话框
3. `ToolWindow`：工具面板
4. `PopupWindow`：弹出菜单、下拉列表、提示框

### 6.3 WindowManager 职责

- Z 序管理
- 激活窗口切换
- 非客户区命中测试
- 窗口拖拽/调整大小
- 模态限制
- 焦点转移
- 脏区刷新

## 7. 控件树模型

### 7.1 选择 retained-mode

不建议做 IMGUI。原因很直接：

- 经典桌面 UI 依赖持久状态
- 文本输入、选择区、滚动位置、菜单展开链路都天然适合 retained-mode
- Tab 焦点顺序、默认按钮、加速键、模态窗口都更像 Win32

### 7.2 基础接口

建议核心接口尽量小，把复杂逻辑放到可复用基类中：

```go
type Widget interface {
    ID() string
    Bounds() Rect
    SetBounds(Rect)
    Visible() bool
    Enabled() bool
    Parent() Widget
    Children() []Widget

    Measure(ctx *LayoutContext, constraint Size) Size
    Layout(ctx *LayoutContext)

    HitTest(p Point) HitResult
    HandleMessage(ctx *UIContext, msg Message) Result
    Paint(ctx *PaintContext)
}
```

### 7.3 基类字段

所有控件建议内嵌一个 `BaseWidget`：

```go
type BaseWidget struct {
    id        string
    rect      Rect
    visible   bool
    enabled   bool
    focused   bool
    hovered   bool
    captured  bool
    tabStop   bool
    zIndex    int
    parent    Widget
    children  []Widget
    themePart string
}
```

### 7.4 状态组织原则

- 基类持有通用状态
- 具体控件持有行为状态
- 不把临时交互状态散落在全局

例如按钮自己的状态应包括：

- `pressed`
- `defaultButton`
- `checked`，仅对 toggle/button-like 控件有效

## 8. 消息与事件模型

### 8.1 内部消息风格

内部建议采用类似 Win32 的消息风格，但消息体做成类型安全结构，不要完全照搬 `uintptr`。

```go
type MessageKind int

const (
    MsgCreate MessageKind = iota
    MsgDestroy
    MsgLayout
    MsgPaint
    MsgMouseMove
    MsgMouseDown
    MsgMouseUp
    MsgMouseWheel
    MsgMouseEnter
    MsgMouseLeave
    MsgKeyDown
    MsgKeyUp
    MsgChar
    MsgTextInput
    MsgTextEditing
    MsgSetFocus
    MsgKillFocus
    MsgCommand
    MsgNotify
    MsgTimer
    MsgThemeChanged
)
```

### 8.2 对外 API 风格

对库使用者不建议暴露底层消息细节，公开层更适合：

- `OnClick`
- `OnChange`
- `OnSelect`
- `OnClose`
- `OnCommand`

内部仍保留消息路由，因为它更适合描述控件间协作。

### 8.3 事件路由顺序

鼠标事件分发建议遵循：

1. 若存在 `mouse capture`，优先发给 capture 控件
2. 否则对顶层窗口做 Z 序命中测试
3. 再向控件树最深子节点递归命中
4. 产生 enter/leave
5. 决定 focus 是否转移

### 8.4 通知模型

按钮、编辑框、列表框等控件建议向父节点发 `Command/Notify`：

- `BN_CLICKED`
- `EN_CHANGE`
- `LBN_SELCHANGE`

这样父容器处理业务逻辑会更简单，也更接近经典桌面 GUI 的编程习惯。

## 9. 主循环设计

### 9.1 线程模型

UI 必须单线程。所有控件状态读写都限定在 UI 线程。

后台 goroutine 只能通过：

- `App.Post(func())`
- `App.Invoke(func())`

把任务切回 UI 线程，不允许后台直接操作控件。

### 9.2 事件循环建议

```go
for app.running {
    timeout := app.NextWakeupTimeout()
    evt := sdl.WaitEventTimeout(timeout)
    if evt != nil {
        app.backend.TranslateAndDispatch(evt)
    }

    app.RunPostedTasks()
    app.ProcessExpiredTimers()

    if app.desktop.NeedsLayout() {
        app.desktop.Layout()
    }

    if app.desktop.HasDirtyRegion() {
        app.renderer.BeginFrame()
        app.desktop.Paint(app.renderer)
        app.renderer.Present()
    }
}
```

### 9.3 为什么不用 SDL 定时器回调直接改 UI

SDL timer 回调不在 UI 线程里，直接操作控件容易引入竞态。建议自己维护一个最小堆定时器表，在 UI 主循环里统一触发。

## 10. 坐标系、DPI 与缩放

### 10.1 逻辑坐标

内部全部使用逻辑像素：

- `Point{X, Y}`
- `Rect{X, Y, W, H}`

全部为整数。

### 10.2 呈现缩放

为了让经典 1px 边框在现代屏幕上仍然清晰，建议支持整数倍缩放：

- 1x
- 2x
- 3x

输入事件先从宿主像素坐标转换成逻辑坐标后再参与命中测试。

### 10.3 高 DPI 原则

不做任意浮点缩放。经典风格最怕 1.25x、1.5x 这类非整数缩放，会让边框和字体都发虚。首版只支持整数倍。

### 10.4 可选 DLU 支持

如果你计划做“像 Win32 对话框资源那样”的布局，建议支持 `DLU` 转换：

```text
pxX = dluX * baseUnitX / 4
pxY = dluY * baseUnitY / 8
```

其中 `baseUnitX/baseUnitY` 来自当前对话框字体测量值。

这对复刻经典对话框尺寸很有帮助。

## 11. 渲染系统设计

### 11.1 选择软件画布

建议维护一块 `RGBA8888` 逻辑帧缓冲：

```go
type Canvas struct {
    Pix    []uint32
    Width  int
    Height int
    Clip   Rect
}
```

所有控件绘制调用都写到 `Canvas` 上，最后通过 `SDL_UpdateTexture` 提交到 streaming texture。

### 11.2 Paint primitive

最小绘制原语建议包括：

- `FillRect`
- `DrawHLine`
- `DrawVLine`
- `FrameRect`
- `DrawBevel`
- `DrawInsetRect`
- `Blit`
- `DrawGlyphRun`
- `InvertRect`
- `PushClip`
- `PopClip`

Windows Classic 的视觉核心几乎都可以由这些 primitive 组合出来。

### 11.3 3D 边框规则

经典立体感的关键是颜色顺序。以普通按钮为例：

- 上边/左边使用亮色
- 下边/右边使用暗色
- 按下状态时反过来

建议主题中至少定义：

- `Face`
- `Light`
- `Lightest`
- `Shadow`
- `DarkShadow`
- `Window`
- `WindowText`
- `Highlight`
- `HighlightText`
- `GrayText`
- `ActiveCaption`
- `InactiveCaption`
- `CaptionText`

### 11.4 非客户区绘制

顶层经典窗口建议分层绘制：

1. 外边框
2. 标题栏背景
3. 系统图标
4. 标题文本
5. 最小化、最大化、关闭按钮
6. 客户区边框
7. 客户区内容

### 11.5 裁剪与脏区

需要同时支持：

- 控件级 clip stack
- 窗口级可见区域裁剪
- 脏区合并

窗口移动或 resize 时，至少需要无效化：

- 原位置区域
- 新位置区域

## 12. 主题与资源系统

### 12.1 主题数据结构

建议把主题拆成三部分：

1. 颜色
2. 度量
3. 位图资源

```go
type Theme struct {
    Colors  ColorTable
    Metrics Metrics
    Assets  AssetSet
    Fonts   FontSet
}
```

### 12.2 主题文件格式

为了简化依赖，建议使用 JSON：

```json
{
  "name": "windows-classic",
  "colors": {
    "face": "#C0C0C0",
    "light": "#FFFFFF",
    "shadow": "#808080",
    "dark_shadow": "#000000",
    "highlight": "#000080",
    "highlight_text": "#FFFFFF"
  },
  "metrics": {
    "border_width": 2,
    "caption_height": 18,
    "menu_height": 18,
    "scrollbar_size": 16,
    "button_padding_x": 6,
    "button_padding_y": 4
  }
}
```

### 12.3 位图资源策略

经典风格很依赖小图标和系统按钮资源。建议：

- 图标统一使用 16x16 或 12x12
- 系统按钮优先使用单色/少色位图
- 所有主题资源打进 atlas
- 用 `go:embed` 提供内置默认主题

### 12.4 字体策略

不要把字体完全硬编码进主题包。建议主题只描述：

- 主字体候选名称
- 候选字体路径
- 字号
- 是否启用单色渲染
- 行高和基线修正

### 12.5 默认主题度量建议

下面这组数值可作为第一版默认主题的 1x 逻辑像素基线：

- `border_width = 2`
- `window_frame_inner = 1`
- `caption_height = 18`
- `menu_height = 18`
- `icon_size_small = 16`
- `scrollbar_size = 16`
- `button_min_height = 21`
- `checkbox_glyph_size = 13`
- `radio_glyph_size = 12`
- `edit_padding_x = 3`
- `edit_padding_y = 2`
- `groupbox_title_offset_x = 8`
- `focus_rect_inset = 3`

这些值不要求完全照搬真实 Win32 系统指标，但要保证整套控件在同一主题下是自洽的。

## 13. 文本系统

### 13.1 要求

文本系统必须支持：

- UTF-8
- 中文
- 英文
- 快捷键下划线标记
- 单行和多行测量
- 文本裁剪
- 光标和选择区
- IME 组合输入

### 13.2 字体渲染建议

推荐使用 `SDL_ttf` 负责 glyph 栅格化，但不要每次绘制时实时生成整段文本位图。应采用 `glyph cache`：

1. 按 `font + rune + style` 为 key 缓存 glyph bitmap
2. 绘制时拼接 glyph run
3. 文本测量结果也做缓存

### 13.3 经典观感处理

为接近经典 UI，建议支持两种文本模式：

1. `MonoSharp`
   - 单色或近似单色
   - 更接近 Windows 9x 风格
2. `GraySharp`
   - 少量灰阶
   - 中文小字号更易读

对中文来说，完全 1-bit 字体在某些字号下会显著降低可读性，因此默认可选 `GraySharp`，但按钮、菜单这类小字号场景要保持边缘足够利落。

### 13.4 加速键与助记符

控件标题中使用 `&` 标记助记符，例如：

- `&File`
- `(&F)文件`

渲染层负责：

- 去掉实际显示中的 `&`
- 记录下划线字符位置
- 根据系统状态决定是否显示下划线

### 13.5 IME 设计

中文输入必须支持 SDL 文本输入事件：

- `SDL_TEXTINPUT`
- `SDL_TEXTEDITING`

编辑控件在获得焦点时：

1. 调用 `SDL_StartTextInput()`
2. 用 `SDL_SetTextInputRect()` 把候选框锚到当前光标附近
3. 把组合态字符串单独绘制，通常使用虚线下划线或点线标记

失焦时调用 `SDL_StopTextInput()`。

## 14. 布局系统

### 14.1 布局原则

经典 GUI 和现代 Web UI 不一样。这里更适合：

- 绝对定位
- 对话框单位布局
- 简单线性布局
- 栅格布局
- Dock 布局

不建议首版就引入过重的约束求解系统。

### 14.2 最小布局器集合

建议首版提供：

1. `AbsoluteLayout`
2. `VBox`
3. `HBox`
4. `Grid`
5. `Dock`

### 14.3 测量协议

控件测量分为三类：

- 固定尺寸，如 checkbox glyph、caption button
- 内容驱动，如 label、button
- 剩余空间驱动，如 listbox、panel

测量接口返回首选尺寸，布局器负责最终放置。

## 15. 焦点、捕获与键盘导航

### 15.1 焦点模型

需要至少区分：

- `active window`
- `focused widget`
- `default button`
- `menu tracking state`

### 15.2 鼠标捕获

以下场景必须支持 capture：

- 按钮按下拖出再拖回
- 滚动条拖拽 thumb
- 窗口拖动
- 窗口 resize
- 菜单跟踪

### 15.3 Tab 导航

焦点遍历顺序建议：

1. 同级按 `tab order`
2. 跳过不可见/不可用控件
3. 进入容器的第一个可聚焦子控件
4. `Shift+Tab` 反向

### 15.4 Enter/Esc 约定

对话框建议遵循经典行为：

- `Enter` 触发默认按钮
- `Esc` 触发取消或关闭

## 16. 核心控件设计

### 16.1 TopLevelWindow

职责：

- 非客户区绘制
- 激活/失活外观
- 拖动
- resize hit test
- 系统按钮行为
- 菜单栏承载

状态：

- `active`
- `minimized`
- `maximized`
- `resizable`
- `modal`

建议为非客户区命中测试定义显式枚举：

- `HTNowhere`
- `HTClient`
- `HTCaption`
- `HTSysMenu`
- `HTClose`
- `HTMinButton`
- `HTMaxButton`
- `HTLeft`
- `HTRight`
- `HTTop`
- `HTBottom`
- `HTTopLeft`
- `HTTopRight`
- `HTBottomLeft`
- `HTBottomRight`

这样窗口拖拽、resize、系统按钮按压态和鼠标光标切换都可以走统一逻辑。

### 16.2 Panel

最基础容器。支持：

- 背景色
- 边框
- 子控件布局
- 可选裁剪

### 16.3 Label

支持：

- 单行/多行
- 左中右对齐
- 文本截断
- 助记符渲染

### 16.4 Button

状态机：

- Normal
- Hot
- Pressed
- Disabled
- Focused
- Default

绘制规则：

- 普通按钮使用凸起 3D 边框
- 按下时使用凹陷边框
- 焦点态绘制虚线 focus rect
- 默认按钮额外外描边

### 16.5 CheckBox / RadioButton

组成：

- 固定大小 indicator
- 文本标签

交互：

- 空格切换
- 鼠标点击 indicator 或文字都生效
- RadioButton 在同组内互斥

### 16.6 Edit

这是第一批控件里最复杂的一个。

至少支持：

- 单行编辑
- 光标移动
- 选择区
- Home/End
- Ctrl+C/Ctrl+V/Ctrl+X
- Backspace/Delete
- 水平滚动
- IME 组合态
- 占位文本，可选
- 只读模式

内部建议把文本编辑能力拆到单独的 `text/model` 或 `text/editstate`，不要把所有逻辑塞进控件本体。

### 16.7 ScrollBar

组成：

- decrease button
- increase button
- track
- thumb

状态：

- hover part
- pressed part
- drag tracking

支持：

- line scroll
- page scroll
- thumb drag
- auto-repeat

### 16.8 ListBox

支持：

- 单选
- 多选，可第二阶段
- 键盘导航
- 滚动条联动
- 选中高亮
- 双击通知

建议行高与字体行高绑定，不要写死。

### 16.9 MenuBar / PopupMenu

经典风格的关键控件之一。

功能要求：

- Alt 激活菜单栏
- 左右键切换顶级菜单
- 上下键切换菜单项
- Enter 执行
- Esc 关闭
- 鼠标 hover 打开子菜单
- 支持分隔线、勾选、禁用、快捷键文字

建议单独实现 `MenuTracker`，它负责整个菜单链路，不要把逻辑散在 MenuItem 上。

### 16.10 StatusBar

简单但很常用。支持：

- 单分区
- 多分区
- 凹陷边框
- 右侧状态文本

## 17. 菜单与命令系统

### 17.1 命令 ID

建议所有菜单项和按钮都能绑定稳定的命令 ID：

```go
type CommandID string
```

例如：

- `cmd.open`
- `cmd.save`
- `cmd.exit`

### 17.2 菜单项结构

```go
type MenuItem struct {
    ID         CommandID
    Text       string
    Shortcut   string
    Enabled    bool
    Checked    bool
    Separator  bool
    Submenu    *Menu
}
```

### 17.3 统一命令入口

应用层建议通过统一回调接收命令：

```go
app.OnCommand(func(cmd CommandID) {
    switch cmd {
    case "cmd.open":
    case "cmd.exit":
    }
})
```

这样菜单、工具栏、按钮触发同一命令时能复用业务逻辑。

## 18. API 草图

### 18.1 建议的对外使用方式

```go
app := classicui.NewApp(classicui.Config{
    Title:        "Demo",
    LogicalSize:  classicui.Size{W: 640, H: 480},
    PresentScale: 2,
    Theme:        classicui.DefaultClassicTheme(),
})

win := classicui.NewWindow("main", classicui.Rect{X: 40, Y: 40, W: 360, H: 220})
win.SetTitle("我的电脑")

btn := widgets.NewButton("ok", "&确定")
btn.OnClick(func() {
    app.Quit()
})

edit := widgets.NewEdit("path")
edit.SetText("C:\\")

win.Content().SetLayout(layout.VBox{Padding: 8, Spacing: 6})
win.Content().AddChild(edit)
win.Content().AddChild(btn)

app.Desktop().AddWindow(win)
app.Run()
```

### 18.2 对外 API 设计原则

- 业务侧尽量不直接处理消息常量
- 业务侧尽量通过控件对象和回调组织逻辑
- 低层能力仍保留扩展口，例如自定义绘制、自定义消息处理

## 19. 状态失效、布局失效与重绘失效

建议明确区分三种失效：

1. `InvalidateLayout`
   - 文本变化、子控件增删、字体变化
2. `InvalidatePaint`
   - hover、pressed、focus、selection 变化
3. `InvalidateRegion(rect)`
   - 局部重绘

不要在每次鼠标移动时都全窗口重排，更不要整帧无条件重绘。

## 20. 性能策略

### 20.1 脏区优先

UI 不是游戏，不需要 60 FPS 持续刷新。只有在以下情况才重绘：

- 输入导致界面变化
- 定时器触发，例如 caret blink
- 窗口拖拽/resize

### 20.2 缓存策略

建议缓存：

- glyph bitmap
- 文本测量结果
- 图标 atlas
- 某些复杂控件的离屏内容，如长列表静态区

### 20.3 列表优化

当列表项很多时，ListBox 要做可视区域裁剪，只绘制当前可见行。

## 21. 测试策略

### 21.1 单元测试

优先覆盖：

- Rect 命中测试
- 焦点遍历
- 布局测量
- 滚动条数值映射
- 编辑框光标移动和选择逻辑
- 菜单导航状态机

### 21.2 Golden Image 测试

软件画布非常适合做截图回归。建议：

1. 构造控件树
2. 调用 Paint
3. 输出 PNG
4. 与基准图比较

重点测试：

- 普通按钮
- 按下按钮
- 激活/失活标题栏
- checkbox/radio
- 选中 listbox item
- 菜单展开

### 21.3 事件回放测试

对复杂控件建议做脚本化事件序列回放，例如 Edit：

- 点击
- 输入文字
- 方向键移动
- Shift 选择
- Backspace 删除

## 22. 风格还原细节建议

要让界面像经典 Windows，真正关键的不是“灰色”，而是下面这些细节：

1. 边框和阴影必须是硬边，不做圆角。
2. 标题栏高度、按钮尺寸、菜单留白要统一。
3. 文本通常靠左，按钮文本居中。
4. 聚焦虚线框要足够克制，只在键盘焦点态出现。
5. 滚动条、菜单、按钮的按下态要明显反相。
6. 标题栏激活色和失活色要区分清楚。
7. 图标尺寸、文本基线、内边距必须走主题度量，不要散落魔法数。

## 23. 风险与应对

### 23.1 字体还原风险

风险：

- 不同平台可用字体不同
- 小字号中文字体在 1-bit 渲染下可能难看

应对：

- 主题支持字体候选列表和路径覆盖
- 文本渲染同时支持 `MonoSharp` 和 `GraySharp`
- 用截图测试校验关键控件

### 23.2 IME 复杂度

风险：

- 中文输入组合态、候选框定位、剪贴板交互容易出问题

应对：

- Edit 控件单独建模
- 第一阶段先做好单行编辑
- 用 SDL 原生文本输入事件，不自己拼键盘字符

### 23.3 拖拽和脏区闪烁

风险：

- 窗口拖动、菜单展开链路容易带来过度重绘

应对：

- 强制所有绘制走裁剪
- 移动窗口时同时无效化新旧区域
- 先保证正确，再做 dirty rect 合并优化

## 24. 分阶段实施路线

### Phase 0：基础设施

目标：

- SDL 初始化
- 宿主窗口
- 逻辑画布
- PresentScale
- 主题加载
- 事件翻译
- 基础 widget tree
- 脏区管理

交付物：

- 可以显示一个经典背景和一扇空窗口

### Phase 1：核心窗口与基础控件

目标：

- TopLevelWindow
- Panel
- Label
- Button
- CheckBox
- RadioButton
- 焦点和 capture

交付物：

- 一个可交互的经典对话框 demo

### Phase 2：输入与滚动

目标：

- Edit
- ScrollBar
- ListBox
- 定时器和 caret blink
- 剪贴板
- IME

交付物：

- 文件路径输入、列表选择、键盘导航完整可用

### Phase 3：菜单与命令系统

目标：

- MenuBar
- PopupMenu
- 子菜单链
- 命令路由
- 快捷键

交付物：

- 接近资源管理器风格的小应用壳

### Phase 4：增强控件

可选增加：

- ComboBox
- StatusBar
- TabControl
- TreeView
- Toolbar
- Tooltip

## 25. 我建议的首个里程碑范围

如果你准备真正动手写，建议第一批只做下面这些：

1. `App`
2. `Canvas`
3. `Theme`
4. `Desktop`
5. `WindowManager`
6. `TopLevelWindow`
7. `Panel`
8. `Label`
9. `Button`
10. `Edit`

原因很简单：

- 这几项已经足够搭出一个能用的经典对话框
- 它们能把渲染、输入、焦点、文本、布局几个主风险全部暴露出来
- 先把这批打稳，再扩展 checkbox、listbox、menu，返工会少很多

## 26. 推荐的实现顺序

1. `geom + paint + theme`
2. `backend/sdl2 + app loop`
3. `widget/base + desktop + invalidation`
4. `TopLevelWindow`
5. `Label/Button`
6. `focus/capture/tab order`
7. `Edit + text model + IME`
8. `ScrollBar`
9. `ListBox`
10. `Menu system`

## 27. 结论

这套控件库的关键不是“能画出几个按钮”，而是把下面四件事一次性设计对：

1. 宿主与内部窗口分层
2. 软件画布渲染路径
3. retained-mode 控件树和消息系统
4. 主题、文本、焦点、输入的统一抽象

只要这四层稳，后面的控件扩展会很顺。反过来，如果一开始把渲染、事件、控件状态耦在一起，后面一做 Edit、Menu、IME 就会明显失控。

这份设计适合作为第一版实现蓝图。真正开始编码时，建议先做一个 `cmd/demo`，只验证“窗口 + 标签 + 按钮 + 输入框”四件事，把主题度量和输入模型打稳，再继续扩控件。
