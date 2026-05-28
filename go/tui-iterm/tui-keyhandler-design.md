# TUI 按键处理系统设计文档

## 1. 概述

本文档描述一个通用的 TUI 按键处理系统设计，参考 Vim 的按键模型，支持以下能力：

- 单键绑定（`j`, `k`, `q`）
- 组合序列键（`gg`, `dd`, `dw`）
- 数字前缀（`3dd`, `5j`）
- 模式切换（Normal / Insert / Command / Shell）
- 命令行输入（`:find xxx`, `/search`）
- 实时过滤（输入即过滤）
- Operator-Pending 模式（`d` 等待 motion）

## 2. 整体架构

```
┌─────────────────────────────────────────────────┐
│                  EventLoop                       │
│  termbox.PollEvent() → Event                    │
└──────────────────┬──────────────────────────────┘
                   │
                   ▼
┌─────────────────────────────────────────────────┐
│              ModeDispatcher                      │
│  根据 currentMode 分发到对应 Handler            │
└──────────────────┬──────────────────────────────┘
                   │
       ┌───────────┼───────────┬──────────────┐
       ▼           ▼           ▼              ▼
┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐
│  Normal  │ │  Insert  │ │ Command  │ │  Shell   │
│  Handler │ │  Handler │ │  Handler │ │  Handler │
└──────────┘ └──────────┘ └──────────┘ └──────────┘
       │
       ▼
┌─────────────────────────────────────────────────┐
│           KeySequenceResolver                    │
│  Trie + Timeout + Count 前缀                    │
└─────────────────────────────────────────────────┘
```

## 3. 模式系统

### 3.1 模式定义

```go
type Mode int

const (
    ModeNormal  Mode = iota // 默认模式，接受命令
    ModeInsert              // 插入模式，按键作为文本输入
    ModeCommand             // 命令行模式（: 触发）
    ModeSearch              // 搜索模式（/ 触发）
    ModeShell               // 内嵌终端模式，按键转发给 PTY
    ModeVisual              // 可视选择模式
    ModeOpPending           // Operator-Pending，等待 motion
)
```

### 3.2 模式切换规则

```
Normal ──'i'──→ Insert ──Esc──→ Normal
Normal ──':'──→ Command ──Enter/Esc──→ Normal
Normal ──'/'──→ Search ──Enter/Esc──→ Normal
Normal ──'v'──→ Visual ──Esc──→ Normal
Normal ──'t'──→ Shell ──Ctrl+\──→ Normal
Normal ──'d'──→ OpPending ──motion/timeout/Esc──→ Normal
```

### 3.3 模式分发器

```go
type ModeDispatcher struct {
    current  Mode
    handlers map[Mode]KeyHandler
}

type KeyHandler interface {
    HandleKey(ev termbox.Event) ModeSwitch
}

type ModeSwitch struct {
    NewMode *Mode  // nil 表示不切换
}

func (d *ModeDispatcher) Dispatch(ev termbox.Event) {
    result := d.handlers[d.current].HandleKey(ev)
    if result.NewMode != nil {
        d.current = *result.NewMode
    }
}
```

## 4. 序列键解析器（KeySequenceResolver）

### 4.1 数据结构：前缀树

```go
type KeyNode struct {
    action   Action
    children map[rune]*KeyNode
}

type Action struct {
    Fn       func(ctx *ActionContext)
    IsMotion bool // 是否是 motion（用于 operator-pending）
}

type ActionContext struct {
    Count    int    // 数字前缀
    Register rune   // 寄存器（如 "a, "b）
    Operator rune   // 当前 operator（d, c, y）
}
```

### 4.2 注册绑定示例

```go
func BuildNormalKeymap() *KeyNode {
    root := &KeyNode{children: make(map[rune]*KeyNode)}

    // 单键
    bind(root, "j", Action{Fn: moveDown, IsMotion: true})
    bind(root, "k", Action{Fn: moveUp, IsMotion: true})
    bind(root, "x", Action{Fn: deleteChar})
    bind(root, "p", Action{Fn: paste})
    bind(root, "u", Action{Fn: undo})

    // 序列键
    bind(root, "gg", Action{Fn: goTop, IsMotion: true})
    bind(root, "G",  Action{Fn: goBottom, IsMotion: true})
    bind(root, "dd", Action{Fn: deleteLine})
    bind(root, "dw", Action{Fn: deleteWord})
    bind(root, "yy", Action{Fn: yankLine})
    bind(root, "yw", Action{Fn: yankWord})
    bind(root, "cc", Action{Fn: changeLine})
    bind(root, "cw", Action{Fn: changeWord})
    bind(root, "zz", Action{Fn: centerScreen})
    bind(root, "zt", Action{Fn: scrollTop})
    bind(root, "zb", Action{Fn: scrollBottom})

    return root
}

func bind(root *KeyNode, seq string, action Action) {
    node := root
    for _, ch := range seq {
        if node.children[ch] == nil {
            node.children[ch] = &KeyNode{children: make(map[rune]*KeyNode)}
        }
        node = node.children[ch]
    }
    node.action = action
}
```

### 4.3 解析器状态机

```go
type KeySequenceResolver struct {
    root    *KeyNode
    current *KeyNode
    pending []rune
    count   int          // 数字前缀累加
    timer   *time.Timer
    timeout time.Duration // 建议 500ms

    // 回调
    onAction  func(Action, *ActionContext)
    onPending func(display string) // 更新状态栏显示
}

func (r *KeySequenceResolver) HandleKey(ch rune) {
    // 阶段1：收集数字前缀
    if ch >= '1' && ch <= '9' && r.current == r.root && r.count == 0 {
        r.count = int(ch - '0')
        r.showPending()
        return
    }
    if ch >= '0' && ch <= '9' && r.count > 0 && r.current == r.root {
        r.count = r.count*10 + int(ch-'0')
        r.showPending()
        return
    }

    // 阶段2：在 Trie 中推进
    r.pending = append(r.pending, ch)
    next, exists := r.current.children[ch]

    if !exists {
        // 无匹配：如果当前节点有 action 就执行，否则丢弃
        if r.current.action.Fn != nil {
            r.execute(r.current.action)
        }
        r.reset()
        return
    }

    r.current = next
    r.stopTimer()

    // 情况A：叶子节点，立即执行
    if next.action.Fn != nil && len(next.children) == 0 {
        r.execute(next.action)
        r.reset()
        return
    }

    // 情况B：中间节点且有 action（前缀冲突），启动超时
    if next.action.Fn != nil && len(next.children) > 0 {
        r.startTimer(next.action)
        r.showPending()
        return
    }

    // 情况C：中间节点无 action，继续等待
    r.startTimer(Action{})
    r.showPending()
}

func (r *KeySequenceResolver) execute(action Action) {
    ctx := &ActionContext{Count: max(r.count, 1)}
    r.onAction(action, ctx)
}

func (r *KeySequenceResolver) startTimer(fallback Action) {
    r.timer = time.AfterFunc(r.timeout, func() {
        if fallback.Fn != nil {
            r.execute(fallback)
        }
        r.reset()
    })
}

func (r *KeySequenceResolver) stopTimer() {
    if r.timer != nil {
        r.timer.Stop()
        r.timer = nil
    }
}

func (r *KeySequenceResolver) reset() {
    r.current = r.root
    r.pending = nil
    r.count = 0
    r.stopTimer()
    r.onPending("")
}

func (r *KeySequenceResolver) showPending() {
    display := ""
    if r.count > 0 {
        display = fmt.Sprintf("%d", r.count)
    }
    display += string(r.pending)
    r.onPending(display)
}
```

## 5. Operator-Pending 模式

Vim 中 `d`, `c`, `y` 是 operator，它们需要后接一个 motion 来确定操作范围。
例如：`dw`（删除到词尾）、`d3j`（删除当前行及下面3行）、`y$`（复制到行尾）。

### 5.1 设计

```go
type OperatorPendingHandler struct {
    operator    rune              // 'd', 'c', 'y'
    resolver    *KeySequenceResolver // 复用序列解析器，但只接受 motion
    motionMap   *KeyNode          // motion 专用的 Trie
    count       int               // operator 前的数字前缀
}

func (h *OperatorPendingHandler) Enter(op rune, count int) {
    h.operator = op
    h.count = count
    // resolver 切换到 motionMap
}

func (h *OperatorPendingHandler) HandleKey(ev termbox.Event) ModeSwitch {
    if ev.Key == termbox.KeyEsc {
        // 取消 operator
        return ModeSwitch{NewMode: ptr(ModeNormal)}
    }

    // 双击 operator 表示作用于整行：dd, cc, yy
    if ev.Ch == h.operator {
        executeOnLine(h.operator, h.count)
        return ModeSwitch{NewMode: ptr(ModeNormal)}
    }

    // 否则等待 motion
    // motion 的 count 和 operator 的 count 相乘
    // 例如 2d3j = 删除 6 行
    // 交给 resolver 处理，resolver 只匹配 IsMotion=true 的 action
    return ModeSwitch{}
}
```

### 5.2 Operator + Motion 组合示例

```
输入: 2d3w
解析:
  count1 = 2 (operator 前缀)
  operator = 'd'
  count2 = 3 (motion 前缀)
  motion = 'w' (word forward)
  实际效果: 删除 2*3=6 个 word

输入: dgg
解析:
  operator = 'd'
  motion = 'gg' (go to top)
  实际效果: 删除从当前行到文件顶部

输入: yy
解析:
  operator = 'y'
  motion = 'y' (双击 = 整行)
  实际效果: 复制当前行
```

## 6. 命令行模式（Command Mode）

### 6.1 触发与结构

按 `:` 进入命令行模式，底部出现输入框。

```go
type CommandHandler struct {
    buf       []rune
    cursor    int
    history   []string    // 命令历史
    histIdx   int
    commands  map[string]*CommandDef
}

type CommandDef struct {
    Name     string
    Execute  func(args []string) error
    Complete func(partial string) []string // Tab 补全
}
```

### 6.2 按键处理

```go
func (h *CommandHandler) HandleKey(ev termbox.Event) ModeSwitch {
    switch ev.Key {
    case termbox.KeyEnter:
        h.executeCommand()
        return ModeSwitch{NewMode: ptr(ModeNormal)}

    case termbox.KeyEsc:
        h.reset()
        return ModeSwitch{NewMode: ptr(ModeNormal)}

    case termbox.KeyBackspace2:
        if h.cursor > 0 {
            h.buf = append(h.buf[:h.cursor-1], h.buf[h.cursor:]...)
            h.cursor--
        } else {
            // 缓冲区为空时退格 = 退出命令模式
            return ModeSwitch{NewMode: ptr(ModeNormal)}
        }

    case termbox.KeyArrowUp:
        h.historyPrev()

    case termbox.KeyArrowDown:
        h.historyNext()

    case termbox.KeyArrowLeft:
        if h.cursor > 0 { h.cursor-- }

    case termbox.KeyArrowRight:
        if h.cursor < len(h.buf) { h.cursor++ }

    case termbox.KeyTab:
        h.tabComplete()

    case termbox.KeyCtrlW:
        h.deleteWordBackward()

    case termbox.KeyCtrlU:
        h.buf = h.buf[h.cursor:]
        h.cursor = 0

    default:
        if ev.Ch != 0 {
            // 插入字符
            h.buf = append(h.buf[:h.cursor],
                append([]rune{ev.Ch}, h.buf[h.cursor:]...)...)
            h.cursor++
        }
    }
    h.render()
    return ModeSwitch{}
}
```

### 6.3 带实时过滤的命令（find 示例）

```go
type FilterableCommand struct {
    name     string
    search   func(query string) []string
    filter   func(results []string, pattern string) []string
}

type CommandWithFilter struct {
    base     *CommandHandler
    phase    CommandPhase
    results  []string
    filtered []string
    filterBuf []rune
    selected int
}

type CommandPhase int
const (
    PhaseInput  CommandPhase = iota // 输入命令参数
    PhaseFilter                     // 对结果实时过滤
)

func (c *CommandWithFilter) HandleKey(ev termbox.Event) ModeSwitch {
    switch c.phase {
    case PhaseInput:
        if ev.Key == termbox.KeyEnter {
            query := string(c.base.buf)
            c.results = c.base.commands["find"].Execute(query)
            c.filtered = c.results
            c.filterBuf = nil
            c.phase = PhaseFilter
            c.renderResults()
            return ModeSwitch{}
        }
        return c.base.HandleKey(ev)

    case PhaseFilter:
        switch ev.Key {
        case termbox.KeyEnter:
            c.confirmSelection()
            return ModeSwitch{NewMode: ptr(ModeNormal)}
        case termbox.KeyEsc:
            return ModeSwitch{NewMode: ptr(ModeNormal)}
        case termbox.KeyArrowUp:
            c.selected = max(0, c.selected-1)
        case termbox.KeyArrowDown:
            c.selected = min(len(c.filtered)-1, c.selected+1)
        case termbox.KeyBackspace2:
            if len(c.filterBuf) > 0 {
                c.filterBuf = c.filterBuf[:len(c.filterBuf)-1]
                c.applyFilter()
            }
        default:
            if ev.Ch != 0 {
                c.filterBuf = append(c.filterBuf, ev.Ch)
                c.applyFilter()
            }
        }
        c.renderResults()
        return ModeSwitch{}
    }
    return ModeSwitch{}
}

func (c *CommandWithFilter) applyFilter() {
    pattern := string(c.filterBuf)
    c.filtered = fuzzyMatch(c.results, pattern)
    c.selected = 0
}
```

## 7. 特殊按键处理

### 7.1 修饰键（Ctrl / Alt 组合）

termbox-go 中 Ctrl 组合键通过 `ev.Key` 获取，Alt 组合通过 `ev.Mod` 判断：

```go
func handleSpecialKeys(ev termbox.Event) {
    // Ctrl 组合
    switch ev.Key {
    case termbox.KeyCtrlC:  interrupt()
    case termbox.KeyCtrlZ:  suspend()
    case termbox.KeyCtrlR:  redo()
    case termbox.KeyCtrlF:  pageDown()
    case termbox.KeyCtrlB:  pageUp()
    case termbox.KeyCtrlD:  halfPageDown()
    case termbox.KeyCtrlU:  halfPageUp()
    }

    // Alt 组合（termbox 中 Alt 表现为 Mod 标志）
    if ev.Mod&termbox.ModAlt != 0 {
        switch ev.Ch {
        case 'j': moveParagraphDown()
        case 'k': moveParagraphUp()
        }
    }
}
```

### 7.2 功能键

```go
switch ev.Key {
case termbox.KeyF1:       showHelp()
case termbox.KeyF5:       refresh()
case termbox.KeyHome:     goLineStart()
case termbox.KeyEnd:      goLineEnd()
case termbox.KeyPgup:     pageUp()
case termbox.KeyPgdn:     pageDown()
case termbox.KeyArrowUp:  moveUp()
case termbox.KeyArrowDown: moveDown()
}
```

### 7.3 全局按键（任何模式下都生效）

某些按键需要在所有模式下都能响应，在分发到具体 handler 之前拦截：

```go
func (d *ModeDispatcher) Dispatch(ev termbox.Event) {
    // 全局拦截层
    if ev.Key == termbox.KeyCtrlC {
        // 任何模式下 Ctrl+C 都回到 Normal
        d.current = ModeNormal
        return
    }

    // 再分发到当前模式
    result := d.handlers[d.current].HandleKey(ev)
    if result.NewMode != nil {
        d.current = *result.NewMode
    }
}
```

## 8. Shell 模式的按键转发

### 8.1 设计要点

Shell 模式下，除了退出键外所有按键都要转发给 PTY：

```go
type ShellHandler struct {
    pty     *os.File
    exitKey termbox.Key // Ctrl+\ 退出
}

func (h *ShellHandler) HandleKey(ev termbox.Event) ModeSwitch {
    // 唯一的退出键
    if ev.Key == h.exitKey {
        return ModeSwitch{NewMode: ptr(ModeNormal)}
    }

    // 其余全部转发
    data := termboxEventToBytes(ev)
    h.pty.Write(data)
    return ModeSwitch{}
}

// termbox 事件转为终端字节序列
func termboxEventToBytes(ev termbox.Event) []byte {
    if ev.Ch != 0 {
        buf := make([]byte, 4)
        n := utf8.EncodeRune(buf, ev.Ch)
        return buf[:n]
    }
    switch ev.Key {
    case termbox.KeyEnter:      return []byte{'\r'}
    case termbox.KeyBackspace2: return []byte{0x7f}
    case termbox.KeyTab:        return []byte{'\t'}
    case termbox.KeyEsc:        return []byte{0x1b}
    case termbox.KeyCtrlC:      return []byte{0x03}
    case termbox.KeyCtrlD:      return []byte{0x04}
    case termbox.KeyCtrlZ:      return []byte{0x1a}
    case termbox.KeyArrowUp:    return []byte{0x1b, '[', 'A'}
    case termbox.KeyArrowDown:  return []byte{0x1b, '[', 'B'}
    case termbox.KeyArrowRight: return []byte{0x1b, '[', 'C'}
    case termbox.KeyArrowLeft:  return []byte{0x1b, '[', 'D'}
    default:
        if ev.Key > 0 && ev.Key < 0x20 {
            return []byte{byte(ev.Key)}
        }
        return nil
    }
}
```

## 9. 完整执行流程示例（含函数调用与状态快照）

### 9.1 示例：`3dw`（删除3个词）

假设 Trie 结构：
```
root
├── 'd' → action: nil, children: {'d': deleteLine, 'w': deleteWord, '$': deleteToEnd}
├── 'g' → action: nil, children: {'g': goTop}
├── 'j' → action: moveDown (叶子)
└── 'k' → action: moveUp (叶子)
```

```
═══════════════════════════════════════════════════════════════
初始状态快照:
  ModeDispatcher.current = ModeNormal
  KeySequenceResolver:
    .current = root
    .pending = []
    .count   = 0
    .timer   = nil
  状态栏: ""
═══════════════════════════════════════════════════════════════

──── 用户按下 '3' ────

调用链:
  EventLoop: ev = termbox.PollEvent() → {Type: EventKey, Ch: '3'}
  ModeDispatcher.Dispatch(ev)
    → handlers[ModeNormal].HandleKey(ev)
      → KeySequenceResolver.HandleKey('3')
        → 进入数字前缀分支: ch >= '1' && ch <= '9' && current == root && count == 0
        → r.count = int('3' - '0') = 3
        → r.showPending()
          → display = "3" + "" = "3"
          → r.onPending("3")  // 回调，更新状态栏

状态快照:
  KeySequenceResolver:
    .current = root        ← 没有移动
    .pending = []          ← 数字不进 pending
    .count   = 3
    .timer   = nil
  状态栏: "3"

═══════════════════════════════════════════════════════════════

──── 用户按下 'd' ────

调用链:
  EventLoop: ev = termbox.PollEvent() → {Type: EventKey, Ch: 'd'}
  ModeDispatcher.Dispatch(ev)
    → handlers[ModeNormal].HandleKey(ev)
      → KeySequenceResolver.HandleKey('d')
        → 不是数字前缀（'d' 不在 '0'-'9'）
        → r.pending = append([], 'd') = ['d']
        → next, exists = r.current.children['d']
          → exists = true, next = 'd'节点
        → r.current = next（移动到 'd' 节点）
        → r.stopTimer()  // timer 本来就是 nil，无操作
        → 检查情况A: next.action.Fn != nil? → NO（'d' 节点本身无 action）
        → 检查情况B: 跳过
        → 进入情况C: 中间节点无 action，继续等待
        → r.startTimer(Action{})  // 空 fallback，超时后只 reset
          → r.timer = time.AfterFunc(500ms, func() { reset() })
        → r.showPending()
          → display = "3" + "d" = "3d"
          → r.onPending("3d")

状态快照:
  KeySequenceResolver:
    .current = 'd'节点     ← 已移动
    .pending = ['d']
    .count   = 3
    .timer   = 活跃（500ms 后触发 reset）
  状态栏: "3d"

═══════════════════════════════════════════════════════════════

──── 用户按下 'w'（距上次按键 150ms，在超时前）────

调用链:
  EventLoop: ev = termbox.PollEvent() → {Type: EventKey, Ch: 'w'}
  ModeDispatcher.Dispatch(ev)
    → handlers[ModeNormal].HandleKey(ev)
      → KeySequenceResolver.HandleKey('w')
        → 不是数字前缀
        → r.pending = append(['d'], 'w') = ['d','w']
        → next, exists = r.current.children['w']
          → r.current 是 'd'节点
          → 'd'节点.children['w'] 存在, next = 'dw'节点
        → r.current = next（移动到 'dw' 节点）
        → r.stopTimer()
          → r.timer.Stop()  ← 取消了那个 500ms 定时器
          → r.timer = nil
        → 检查情况A: next.action.Fn != nil? → YES (deleteWord)
                     len(next.children) == 0? → YES (叶子节点)
        → 命中情况A！立即执行
        → r.execute(next.action)
          → ctx = &ActionContext{Count: max(3, 1)} = &ActionContext{Count: 3}
          → r.onAction(Action{Fn: deleteWord}, ctx)
            → deleteWord(ctx)  ← 实际业务：删除 3 个词
        → r.reset()
          → r.current = r.root
          → r.pending = nil
          → r.count = 0
          → r.stopTimer()  // 已经是 nil
          → r.onPending("")  // 清空状态栏

状态快照:
  KeySequenceResolver:
    .current = root        ← 回到初始
    .pending = []
    .count   = 0
    .timer   = nil
  状态栏: ""
  副作用: 光标位置开始的 3 个词已被删除

═══════════════════════════════════════════════════════════════
```

### 9.2 示例：`gg`（跳转到文件顶部）

```
═══════════════════════════════════════════════════════════════
初始状态: current=root, pending=[], count=0, timer=nil

──── 用户按下第一个 'g' ────

调用链:
  KeySequenceResolver.HandleKey('g')
    → pending = ['g']
    → next = root.children['g']  → 存在（'g'节点，无 action，有 children: {'g': goTop}）
    → current = 'g'节点
    → stopTimer() → 无操作
    → 情况A? NO（action 为 nil）
    → 情况B? NO（action 为 nil）
    → 情况C: 中间节点无 action
    → startTimer(Action{})  → 500ms 后 reset（丢弃输入）
    → showPending() → "g"

状态快照:
  .current = 'g'节点
  .pending = ['g']
  .timer   = 活跃（500ms，fallback=空）
  状态栏: "g"

──── 用户按下第二个 'g'（200ms 后）────

调用链:
  KeySequenceResolver.HandleKey('g')
    → pending = ['g','g']
    → next = current.children['g']
      → 'g'节点.children['g'] 存在, next = 'gg'节点（action=goTop, children=空）
    → current = 'gg'节点
    → stopTimer() → 取消 500ms 定时器
    → 情况A? YES（action=goTop, children 为空，叶子）
    → execute(Action{Fn: goTop})
      → ctx = &ActionContext{Count: max(0, 1)} = {Count: 1}
      → goTop(ctx)  ← 光标跳到第 1 行
    → reset()

状态快照: 回到初始，光标在第 1 行
═══════════════════════════════════════════════════════════════
```

### 9.3 示例：`g` 后超时（用户只按了一个 g 就停了）

```
═══════════════════════════════════════════════════════════════
初始状态: current=root, pending=[], count=0, timer=nil

──── 用户按下 'g' ────

（同上，current 移到 'g' 节点，启动 500ms timer，fallback=Action{}）

──── 500ms 过去，无后续输入 ────

调用链（由 timer goroutine 触发）:
  timer 回调执行:
    → fallback.Fn == nil（空 Action）
    → 不执行任何动作
    → reset()
      → current = root
      → pending = nil
      → count = 0
      → onPending("")  → 状态栏清空

状态快照: 回到初始，什么都没发生
效果: 用户按了一个无意义的 'g' 前缀，系统静默丢弃
═══════════════════════════════════════════════════════════════
```

### 9.4 示例：`5j`（向下移动5行）

```
═══════════════════════════════════════════════════════════════
初始状态: current=root, pending=[], count=0, timer=nil

──── 用户按下 '5' ────

调用链:
  KeySequenceResolver.HandleKey('5')
    → ch='5', 满足: ch >= '1' && ch <= '9' && current == root && count == 0
    → r.count = 5
    → showPending() → "5"

状态: current=root, count=5, pending=[], 状态栏="5"

──── 用户按下 'j' ────

调用链:
  KeySequenceResolver.HandleKey('j')
    → 不是数字
    → pending = ['j']
    → next = root.children['j'] → 存在（叶子，action=moveDown）
    → current = 'j'节点
    → stopTimer() → 无操作
    → 情况A? YES（action=moveDown, children 为空）
    → execute(Action{Fn: moveDown})
      → ctx = &ActionContext{Count: max(5, 1)} = {Count: 5}
      → moveDown(ctx)  ← 光标向下移动 5 行
    → reset()

状态: 回到初始，光标下移了 5 行
═══════════════════════════════════════════════════════════════
```

### 9.5 示例：`dd`（删除整行，前缀冲突场景）

假设 'd' 节点本身绑定了 action（比如进入 delete-pending 状态）：
```
root
└── 'd' → action: enterDeletePending, children: {'d': deleteLine, 'w': deleteWord}
```

```
═══════════════════════════════════════════════════════════════
初始状态: current=root, pending=[], count=0, timer=nil

──── 用户按下第一个 'd' ────

调用链:
  KeySequenceResolver.HandleKey('d')
    → pending = ['d']
    → next = root.children['d'] → 存在
    → current = 'd'节点
    → stopTimer()
    → 情况A? NO（有 children）
    → 情况B? next.action.Fn != nil (enterDeletePending) && len(children) > 0
      → YES! 前缀冲突
    → startTimer(Action{Fn: enterDeletePending})
      → timer = AfterFunc(500ms, func() {
            execute(enterDeletePending)  // 如果超时就执行单 'd' 的动作
            reset()
        })
    → showPending() → "d"

状态: current='d'节点, timer=活跃(fallback=enterDeletePending), 状态栏="d"

──── 用户按下第二个 'd'（100ms 后）────

调用链:
  KeySequenceResolver.HandleKey('d')
    → pending = ['d','d']
    → next = current.children['d']
      → 'd'节点.children['d'] 存在, next = 'dd'节点（action=deleteLine, 叶子）
    → current = 'dd'节点
    → stopTimer() → 取消 500ms timer（enterDeletePending 不会执行）
    → 情况A? YES（action=deleteLine, children 为空）
    → execute(Action{Fn: deleteLine})
      → ctx = {Count: 1}
      → deleteLine(ctx)  ← 删除当前行
    → reset()

状态: 回到初始，当前行已删除
═══════════════════════════════════════════════════════════════
```

### 9.6 示例：按 'd' 后超时（只按了一个 d）

```
═══════════════════════════════════════════════════════════════
（接上例，用户按了 'd' 后什么都不按）

状态: current='d'节点, timer=活跃(fallback=enterDeletePending)

──── 500ms 过去 ────

调用链（timer goroutine）:
  timer 回调:
    → fallback.Fn != nil（enterDeletePending）
    → execute(Action{Fn: enterDeletePending})
      → ctx = {Count: 1}
      → enterDeletePending(ctx)  ← 进入 operator-pending 状态
    → reset()

状态: 回到初始（但应用层已切换到 OpPending 模式）
═══════════════════════════════════════════════════════════════
```

### 9.7 示例：`:find main` 然后过滤

```
═══════════════════════════════════════════════════════════════
初始状态:
  ModeDispatcher.current = ModeNormal
  CommandHandler: buf=[], cursor=0, phase=PhaseInput

──── 用户按 ':' ────

调用链:
  EventLoop: ev = {Ch: ':'}
  ModeDispatcher.Dispatch(ev)
    → handlers[ModeNormal].HandleKey(ev)
      → 识别 ':' 为模式切换键
      → return ModeSwitch{NewMode: ptr(ModeCommand)}
  ModeDispatcher:
    → result.NewMode != nil
    → d.current = ModeCommand
  UI: 底部出现命令行 ":"

状态: ModeDispatcher.current = ModeCommand

──── 用户输入 'f' ────

调用链:
  EventLoop: ev = {Ch: 'f'}
  ModeDispatcher.Dispatch(ev)
    → handlers[ModeCommand].HandleKey(ev)
      → CommandHandler.HandleKey(ev)
        → 进入 default 分支: ev.Ch != 0
        → buf = append(buf[:0], append([]rune{'f'}, buf[0:]...)...)
          → buf = ['f']
        → cursor = 1
        → h.render()  → 底部显示 ":f"
        → return ModeSwitch{} (不切换)

──── 用户继续输入 'i','n','d',' ','m','a','i','n' ────

（每个字符重复上述流程）
最终: buf = ['f','i','n','d',' ','m','a','i','n'], cursor = 9
底部显示: ":find main"

──── 用户按 Enter ────

调用链:
  CommandHandler.HandleKey(ev)
    → ev.Key == termbox.KeyEnter
    → h.executeCommand()
      → 解析 buf: name = "find", args = ["main"]
      → cmd = h.commands["find"]  → 找到 FilterableCommand
      → results = cmd.Execute(["main"])
        → 搜索文件系统，返回 ["main.go", "main_test.go", "cmd/main.go"]
      → 切换到 CommandWithFilter:
        → c.results = ["main.go", "main_test.go", "cmd/main.go"]
        → c.filtered = c.results（初始不过滤）
        → c.filterBuf = []
        → c.selected = 0
        → c.phase = PhaseFilter
      → c.renderResults()
        → 渲染列表:
          → [*] main.go          ← selected
          → [ ] main_test.go
          → [ ] cmd/main.go
        → 底部显示: "filter: "

状态: phase=PhaseFilter, results=3项, filterBuf=[]

──── 用户按 't' ────

调用链:
  CommandWithFilter.HandleKey(ev)
    → phase == PhaseFilter
    → ev.Ch = 't', 进入 default
    → c.filterBuf = append([], 't') = ['t']
    → c.applyFilter()
      → pattern = "t"
      → c.filtered = fuzzyMatch(c.results, "t")
        → "main.go" 包含 't'? NO → 排除
        → "main_test.go" 包含 't'? YES → 保留
        → "cmd/main.go" 包含 't'? NO → 排除
        → filtered = ["main_test.go"]
      → c.selected = 0
    → c.renderResults()
      → [*] main_test.go
      → 底部显示: "filter: t"

状态: filterBuf=['t'], filtered=["main_test.go"], selected=0

──── 用户按 Enter ────

调用链:
  CommandWithFilter.HandleKey(ev)
    → ev.Key == termbox.KeyEnter
    → c.confirmSelection()
      → item = c.filtered[c.selected] = "main_test.go"
      → 执行动作（如打开文件）
    → return ModeSwitch{NewMode: ptr(ModeNormal)}
  ModeDispatcher:
    → d.current = ModeNormal
    → 列表消失，命令行消失

状态: 回到 ModeNormal
═══════════════════════════════════════════════════════════════
```

### 9.8 示例：Shell 模式交互

```
═══════════════════════════════════════════════════════════════
初始状态: ModeDispatcher.current = ModeNormal

──── 用户按 'T'（绑定为 openShell）────

调用链:
  ModeDispatcher.Dispatch(ev)
    → handlers[ModeNormal].HandleKey(ev)
      → KeySequenceResolver.HandleKey('T')
        → root.children['T'] → 叶子节点, action=openShell
        → execute(openShell)
          → openShell(ctx):
            → shell := os.Getenv("SHELL")  // "/bin/zsh"
            → ptyFile, _ = pty.Start(exec.Command(shell))
            → 设置 PTY 窗口大小: rows=10, cols=80
            → 启动 goroutine 读取 PTY 输出 → VT100 解析 → 渲染
        → return ModeSwitch{NewMode: ptr(ModeShell)}
  ModeDispatcher:
    → d.current = ModeShell
  UI: 屏幕第 15-25 行显示 shell 区域，出现 shell prompt

状态: ModeShell, PTY 活跃

──── 用户按 'l' ────

调用链:
  ModeDispatcher.Dispatch(ev)
    → handlers[ModeShell].HandleKey(ev)
      → ShellHandler.HandleKey(ev)
        → ev.Key != h.exitKey (Ctrl+\)
        → data = termboxEventToBytes(ev)
          → ev.Ch = 'l', ev.Ch != 0
          → buf = make([]byte, 4)
          → n = utf8.EncodeRune(buf, 'l') = 1
          → return buf[:1] = []byte{0x6c}
        → h.pty.Write([]byte{0x6c})
        → return ModeSwitch{} (不切换)

PTY 侧: 收到 'l'，shell 回显 'l' 到 PTY 输出
渲染 goroutine: 读到 'l' → VT100 解析 → 更新 shell 区域显示

──── 用户按 's' ────

（同上，写入 0x73）

──── 用户按 Enter ────

调用链:
  ShellHandler.HandleKey(ev)
    → ev.Key = termbox.KeyEnter
    → ev.Key != h.exitKey
    → data = termboxEventToBytes(ev)
      → ev.Key == termbox.KeyEnter → return []byte{'\r'} = []byte{0x0d}
    → h.pty.Write([]byte{0x0d})

PTY 侧: shell 收到回车，执行 "ls" 命令，输出文件列表到 PTY
渲染 goroutine: 读到输出字节流 → VT100 解析（处理换行、颜色等）→ 更新 shell 区域

──── 用户按 Ctrl+\ ────

调用链:
  ShellHandler.HandleKey(ev)
    → ev.Key = termbox.KeyCtrlBackslash
    → ev.Key == h.exitKey  → YES!
    → 清理:
      → h.pty.Close()  // 关闭 PTY 文件描述符
      → 等待子进程退出
      → 停止渲染 goroutine
    → return ModeSwitch{NewMode: ptr(ModeNormal)}
  ModeDispatcher:
    → d.current = ModeNormal
  UI: shell 区域消失，恢复正常 TUI 界面

状态: 回到 ModeNormal
═══════════════════════════════════════════════════════════════
```

## 10. 完整按键分类参考（对照 Vim）

### 10.1 需要支持的按键类型

| 类型 | 示例 | 处理方式 |
|------|------|----------|
| 单键命令 | `j`,`k`,`x`,`p`,`u` | 直接查 map 执行 |
| 双键序列 | `gg`,`zz`,`gt` | Trie 匹配 |
| Operator + Motion | `dw`,`cw`,`y$` | OpPending 模式 |
| Operator + 双击 | `dd`,`cc`,`yy` | OpPending 中特判 |
| 数字前缀 + 命令 | `3j`,`5dd`,`2dw` | count 累加器 |
| 数字前缀嵌套 | `2d3w` (=6w) | operator count × motion count |
| 修饰键 | `Ctrl+F`,`Ctrl+U` | ev.Key 直接匹配 |
| 模式切换键 | `i`,`:`,`/`,`v`,`Esc` | 返回 ModeSwitch |
| 文本输入 | Insert 模式下任意字符 | 直接插入缓冲区 |
| 命令行输入 | `:` 后的自由文本 | CommandHandler 处理 |
| 转发键 | Shell 模式下所有键 | 转为字节写入 PTY |

### 10.2 容易遗漏的场景

| 场景 | 说明 | 处理方案 |
|------|------|----------|
| **Esc 的多义性** | Esc 既是退出键，也是 Alt 序列的前缀（终端中 Alt+x = Esc,x） | 对 Esc 也做超时判断：收到 Esc 后等 50ms，如果有后续字节则是 Alt 序列，否则是单独 Esc |
| **粘贴检测** | 用户粘贴大段文本时，快速连续的字符不应触发序列匹配 | 检测输入速率，超过阈值（如 <5ms 间隔）进入 paste 模式，字符直接插入 |
| **重复命令 `.`** | Vim 的 `.` 重复上一次编辑操作 | 记录上一次执行的 action + context，`.` 时重放 |
| **宏录制 `q`** | `qa` 开始录制到寄存器 a，`q` 停止，`@a` 回放 | 录制时记录所有按键事件序列，回放时依次 feed 给 HandleKey |
| **Mark `m` + 跳转 `'`** | `ma` 设置标记，`'a` 跳转 | 单键后接任意字符作为参数，不走 Trie |
| **寄存器 `"`** | `"ayy` 复制到寄存器 a | `"` 后读取一个字符作为寄存器名，存入 ActionContext |
| **窗口大小变化** | 终端 resize | 监听 termbox.EventResize，重新计算布局 |
| **鼠标事件** | 点击、滚轮 | termbox.EventMouse，独立处理不走按键系统 |
| **Unicode 输入** | 中文、emoji 等多字节字符 | termbox 的 ev.Ch 是 rune，天然支持 |
| **挂起/恢复** | Ctrl+Z 挂起进程 | 需要正确处理 SIGTSTP/SIGCONT，恢复后重绘 |

### 10.3 CharArgHandler：`m`/`"`/`r`/`f`/`t` 这类"单键+参数字符"的处理

这类按键不适合放 Trie（因为第二个字符是任意的 a-z/A-Z），需要特殊处理。

#### 10.3.1 完整实现

```go
type CharArgHandler struct {
    waiting  bool
    trigger  rune          // 触发字符（用于状态栏显示）
    action   func(ch rune) // 收到参数字符后执行的动作
    validate func(ch rune) bool // 可选：验证参数字符是否合法
}

type NormalHandler struct {
    resolver *KeySequenceResolver
    charArg  *CharArgHandler
    // ...
}

func (h *NormalHandler) HandleKey(ev termbox.Event) ModeSwitch {
    // ═══ 优先级1：如果正在等待参数字符，先处理 ═══
    if h.charArg != nil && h.charArg.waiting {
        if ev.Key == termbox.KeyEsc {
            // Esc 取消等待
            h.charArg = nil
            h.updateStatusBar("")
            return ModeSwitch{}
        }
        if ev.Ch != 0 {
            if h.charArg.validate == nil || h.charArg.validate(ev.Ch) {
                h.charArg.action(ev.Ch)
            }
            h.charArg = nil
            h.updateStatusBar("")
        }
        return ModeSwitch{}
    }

    // ═══ 优先级2：检查是否是 CharArg 触发键 ═══
    switch ev.Ch {
    case 'm':
        h.charArg = &CharArgHandler{
            waiting:  true,
            trigger:  'm',
            action:   h.setMark,
            validate: func(ch rune) bool { return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' },
        }
        h.updateStatusBar("m")
        return ModeSwitch{}

    case '\'':
        h.charArg = &CharArgHandler{
            waiting:  true,
            trigger:  '\'',
            action:   h.jumpToMark,
            validate: func(ch rune) bool { return ch >= 'a' && ch <= 'z' || ch >= 'A' && ch <= 'Z' },
        }
        h.updateStatusBar("'")
        return ModeSwitch{}

    case 'r':
        h.charArg = &CharArgHandler{
            waiting: true,
            trigger: 'r',
            action:  h.replaceChar,
            // validate 为 nil 表示接受任意字符
        }
        h.updateStatusBar("r")
        return ModeSwitch{}

    case 'f':
        h.charArg = &CharArgHandler{
            waiting: true,
            trigger: 'f',
            action:  h.findCharForward,
        }
        h.updateStatusBar("f")
        return ModeSwitch{}

    case 'F':
        h.charArg = &CharArgHandler{
            waiting: true,
            trigger: 'F',
            action:  h.findCharBackward,
        }
        h.updateStatusBar("F")
        return ModeSwitch{}

    case 't':
        h.charArg = &CharArgHandler{
            waiting: true,
            trigger: 't',
            action:  h.tillCharForward,
        }
        h.updateStatusBar("t")
        return ModeSwitch{}

    case '"':
        h.charArg = &CharArgHandler{
            waiting:  true,
            trigger:  '"',
            action:   h.selectRegister,
            validate: func(ch rune) bool { return ch >= 'a' && ch <= 'z' || ch >= '0' && ch <= '9' || ch == '+' || ch == '*' },
        }
        h.updateStatusBar("\"")
        return ModeSwitch{}
    }

    // ═══ 优先级3：交给序列键解析器 ═══
    h.resolver.HandleKey(ev.Ch)
    return ModeSwitch{}
}
```

#### 10.3.2 执行流程：`ma`（设置标记 a）

```
═══════════════════════════════════════════════════════════════
初始状态:
  NormalHandler:
    .charArg  = nil
    .resolver: current=root, pending=[], count=0
  状态栏: ""
═══════════════════════════════════════════════════════════════

──── 用户按下 'm' ────

调用链:
  EventLoop: ev = {Ch: 'm'}
  ModeDispatcher.Dispatch(ev)
    → handlers[ModeNormal].HandleKey(ev)
      → NormalHandler.HandleKey(ev)
        → 检查优先级1: h.charArg == nil → 跳过
        → 检查优先级2: ev.Ch == 'm' → 命中!
        → h.charArg = &CharArgHandler{
              waiting:  true,
              trigger:  'm',
              action:   h.setMark,
              validate: func(ch) { return 'a'-'z' || 'A'-'Z' },
          }
        → h.updateStatusBar("m")
        → return ModeSwitch{} (不切换模式)

状态快照:
  NormalHandler.charArg = {waiting: true, trigger: 'm', action: setMark}
  状态栏: "m"
  （用户看到状态栏显示 "m"，知道系统在等待一个字符）

═══════════════════════════════════════════════════════════════

──── 用户按下 'a' ────

调用链:
  EventLoop: ev = {Ch: 'a'}
  ModeDispatcher.Dispatch(ev)
    → handlers[ModeNormal].HandleKey(ev)
      → NormalHandler.HandleKey(ev)
        → 检查优先级1: h.charArg != nil && h.charArg.waiting == true → 命中!
        → ev.Key != termbox.KeyEsc → 不是取消
        → ev.Ch = 'a', ev.Ch != 0 → 有字符
        → h.charArg.validate('a')
          → 'a' >= 'a' && 'a' <= 'z' → true → 合法
        → h.charArg.action('a')
          → h.setMark('a')
            → marks['a'] = currentCursorPosition  ← 实际业务
        → h.charArg = nil  ← 清除等待状态
        → h.updateStatusBar("")  ← 清空状态栏
        → return ModeSwitch{}

状态快照:
  NormalHandler.charArg = nil  ← 回到正常
  marks = {'a': {line: 42, col: 10}}  ← 标记已设置
  状态栏: ""
═══════════════════════════════════════════════════════════════
```

#### 10.3.3 执行流程：`'a`（跳转到标记 a）

```
═══════════════════════════════════════════════════════════════
初始状态: charArg=nil, 状态栏=""

──── 用户按下 '\'' ────

调用链:
  NormalHandler.HandleKey(ev)
    → 优先级1: charArg == nil → 跳过
    → 优先级2: ev.Ch == '\'' → 命中!
    → h.charArg = &CharArgHandler{
          waiting: true, trigger: '\'', action: h.jumpToMark,
          validate: func(ch) { 'a'-'z' || 'A'-'Z' },
      }
    → updateStatusBar("'")

状态: charArg.waiting=true, 状态栏="'"

──── 用户按下 'a' ────

调用链:
  NormalHandler.HandleKey(ev)
    → 优先级1: charArg != nil && waiting → 命中!
    → validate('a') → true
    → h.charArg.action('a')
      → h.jumpToMark('a')
        → pos = marks['a']  → {line: 42, col: 10}
        → setCursor(pos.line, pos.col)  ← 光标跳转
    → charArg = nil
    → updateStatusBar("")

状态: 光标已跳转到 line 42, col 10
═══════════════════════════════════════════════════════════════
```

#### 10.3.4 执行流程：`rx`（替换当前字符为 x）

```
═══════════════════════════════════════════════════════════════
初始状态: charArg=nil, 光标在 "hello" 的 'h' 上

──── 用户按下 'r' ────

调用链:
  NormalHandler.HandleKey(ev)
    → 优先级2: ev.Ch == 'r' → 命中!
    → h.charArg = &CharArgHandler{
          waiting: true, trigger: 'r', action: h.replaceChar,
          validate: nil,  ← 接受任意字符
      }
    → updateStatusBar("r")

状态: 等待替换字符, 状态栏="r"

──── 用户按下 'x' ────

调用链:
  NormalHandler.HandleKey(ev)
    → 优先级1: charArg.waiting → 命中!
    → validate == nil → 跳过验证（任意字符都行）
    → h.charArg.action('x')
      → h.replaceChar('x')
        → buf[cursor.line][cursor.col] = 'x'  ← "hello" → "xello"
    → charArg = nil
    → updateStatusBar("")

状态: 文本变为 "xello", 光标不动
═══════════════════════════════════════════════════════════════
```

#### 10.3.5 执行流程：`fa`（向前查找字符 a）

```
═══════════════════════════════════════════════════════════════
初始状态: charArg=nil, 光标在 "hello world" 的 'h' 上

──── 用户按下 'f' ────

调用链:
  NormalHandler.HandleKey(ev)
    → 优先级2: ev.Ch == 'f' → 命中!
    → h.charArg = &CharArgHandler{
          waiting: true, trigger: 'f', action: h.findCharForward,
      }
    → updateStatusBar("f")

状态: 等待目标字符, 状态栏="f"

──── 用户按下 'o' ────

调用链:
  NormalHandler.HandleKey(ev)
    → 优先级1: charArg.waiting → 命中!
    → h.charArg.action('o')
      → h.findCharForward('o')
        → 从 cursor.col+1 开始向右扫描当前行
        → 找到 "hello" 中的 'o' 在 col=4
        → setCursor(cursor.line, 4)  ← 光标移到 'o' 上
    → charArg = nil
    → updateStatusBar("")

状态: 光标移到 'o' (col=4)
═══════════════════════════════════════════════════════════════
```

#### 10.3.6 执行流程：`"ayy`（复制当前行到寄存器 a）

这是 CharArgHandler 和序列键解析器的组合使用：

```
═══════════════════════════════════════════════════════════════
初始状态: charArg=nil, resolver: current=root, count=0

──── 用户按下 '"' ────

调用链:
  NormalHandler.HandleKey(ev)
    → 优先级2: ev.Ch == '"' → 命中!
    → h.charArg = &CharArgHandler{
          waiting: true, trigger: '"', action: h.selectRegister,
          validate: func(ch) { 'a'-'z' || '0'-'9' || '+' || '*' },
      }
    → updateStatusBar("\"")

状态: 等待寄存器名, 状态栏="\""

──── 用户按下 'a' ────

调用链:
  NormalHandler.HandleKey(ev)
    → 优先级1: charArg.waiting → 命中!
    → validate('a') → true
    → h.charArg.action('a')
      → h.selectRegister('a')
        → h.activeRegister = 'a'  ← 记住选中的寄存器
        → （不执行任何操作，只是设置上下文）
    → charArg = nil
    → updateStatusBar("\"a")  ← 显示已选寄存器

状态: activeRegister='a', charArg=nil, 状态栏="\"a"
注意: 此时控制权回到正常流程，下一个按键走优先级3（序列解析器）

──── 用户按下 'y' ────

调用链:
  NormalHandler.HandleKey(ev)
    → 优先级1: charArg == nil → 跳过
    → 优先级2: ev.Ch == 'y' → 不在 CharArg 触发列表中
    → 优先级3: h.resolver.HandleKey('y')
      → pending = ['y']
      → root.children['y'] 存在（'y'节点，有 children: {'y': yankLine, 'w': yankWord}）
      → current = 'y'节点
      → 情况C: 中间节点无 action
      → startTimer(Action{})
      → showPending() → "\"ay"

状态: resolver.current='y'节点, pending=['y'], 状态栏="\"ay"

──── 用户按下第二个 'y' ────

调用链:
  NormalHandler.HandleKey(ev)
    → 优先级1: charArg == nil → 跳过
    → 优先级2: 'y' 不是 CharArg 触发键
    → 优先级3: h.resolver.HandleKey('y')
      → pending = ['y','y']
      → current.children['y'] → 'yy'节点（action=yankLine, 叶子）
      → stopTimer()
      → 情况A: 叶子节点
      → execute(Action{Fn: yankLine})
        → ctx = &ActionContext{Count: 1, Register: 'a'}  ← 带上 activeRegister
        → yankLine(ctx)
          → line = getCurrentLine()
          → registers['a'] = line  ← 复制到寄存器 a
      → reset()
    → h.activeRegister = 0  ← 清除寄存器选择

状态: 回到初始, registers['a'] = "当前行内容"
═══════════════════════════════════════════════════════════════
```

#### 10.3.7 执行流程：`m` 后按 Esc（取消）

```
═══════════════════════════════════════════════════════════════
──── 用户按下 'm' ────

（同 10.3.2，进入等待状态）
状态: charArg.waiting=true, 状态栏="m"

──── 用户按下 Esc ────

调用链:
  NormalHandler.HandleKey(ev)
    → 优先级1: charArg.waiting → 命中!
    → ev.Key == termbox.KeyEsc → 取消分支!
    → h.charArg = nil
    → h.updateStatusBar("")
    → return ModeSwitch{}

状态: charArg=nil, 状态栏=""
效果: 什么都没发生，等待被取消
═══════════════════════════════════════════════════════════════
```

#### 10.3.8 CharArgHandler 与 KeySequenceResolver 的优先级关系

```
按键到达 NormalHandler.HandleKey(ev):

┌─────────────────────────────────────────────────────────┐
│ 优先级1: charArg 正在等待?                               │
│   YES → 把这个字符交给 charArg.action 执行               │
│   NO  → 继续                                            │
├─────────────────────────────────────────────────────────┤
│ 优先级2: 是 CharArg 触发键? (m, ', r, f, F, t, T, ")    │
│   YES → 创建 CharArgHandler，进入等待                    │
│   NO  → 继续                                            │
├─────────────────────────────────────────────────────────┤
│ 优先级3: 交给 KeySequenceResolver                        │
│   → Trie 匹配 / 数字前缀 / 超时                         │
└─────────────────────────────────────────────────────────┘

关键点:
- CharArg 等待状态"吃掉"下一个按键，不会让它流入 Trie
- CharArg 触发键本身不进入 Trie（'m' 不会被当作序列的一部分）
- 如果 resolver 正在中间状态（比如已经按了 'g' 在等第二个键），
  此时按 'm' 会怎样？→ 取决于实现：
  方案A: CharArg 优先级高于 resolver 中间状态 → 'g' 被丢弃
  方案B: resolver 中间状态优先 → 'm' 被当作序列的一部分（不匹配则 reset）
  建议: 方案B，因为用户已经开始了一个序列，不应被打断
```

## 11. 架构总结

```
┌────────────────────────────────────────────────────────────┐
│                      EventLoop                              │
│  for { ev := termbox.PollEvent() }                         │
└────────────────────────┬───────────────────────────────────┘
                         │
                         ▼
┌────────────────────────────────────────────────────────────┐
│                   GlobalInterceptor                          │
│  Ctrl+C → force normal | Resize → relayout                 │
└────────────────────────┬───────────────────────────────────┘
                         │
                         ▼
┌────────────────────────────────────────────────────────────┐
│                   ModeDispatcher                             │
│  switch currentMode { ... }                                 │
└───┬────────┬────────┬────────┬────────┬────────┬──────────┘
    │        │        │        │        │        │
    ▼        ▼        ▼        ▼        ▼        ▼
 Normal   Insert  Command  Search  Visual   Shell
    │
    ├── KeySequenceResolver (Trie + timeout + count)
    ├── CharArgHandler (m, ", r, f, t 等)
    └── OperatorPendingHandler (d, c, y + motion)
```

### 设计原则

1. **分层处理**：全局 → 模式 → 具体 handler，每层只关心自己的职责
2. **模式隔离**：每个模式有独立的按键表，互不干扰
3. **可组合**：count、operator、motion 是正交的维度，自由组合
4. **超时解决歧义**：前缀冲突时用时间来区分用户意图
5. **状态可见**：pending 序列实时显示在状态栏，用户知道系统在等什么
6. **退出一致**：Esc 在任何非 Normal 模式下都回到 Normal
