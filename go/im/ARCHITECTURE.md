# 输入法实现技术文档

本文档描述当前仓库中已经落地的输入法实现，而不是早期的概念设计。重点包括：

- 当前代码结构和职责分层
- 词库加载与接口扩展方式
- 查询、切分、候选生成、组合和回退策略
- CLI / TUI 的接入方式
- 如何新增一个词库 source

---

## 1. 当前目录结构

当前核心文件：

```text
pkg/ime/
  engine.go                 输入法引擎核心
  session.go                会话状态 API，供 TUI/宿主调用
  query.go                  查询结构化解析和匹配
  segment.go                拼音切分与 best-effort 拆分
  source.go                 source 接口、registry、source config
  source_fcitx.go           fcitx 词库实现
  source_sogou.go           sogou 候选词库实现
  source_sogou_syllables.go sogou 音节表实现

cmd/im/main.go              纯 CLI 查询 demo
cmd/im-tui/main.go          termbox-go TUI demo
cmd/internal/sourceflags/   入口层 source 参数解析
```

---

## 2. 总体架构

当前实现分成三层：

### 2.1 Source 层

负责把不同格式的词库文件加载进引擎。

统一接口：

```go
type Source interface {
    Name() string
    CollectSyllables(*Engine) error
    Load(*Engine) error
}
```

职责：

- `CollectSyllables`
  - 用于补充合法音节表
  - 不是所有 source 都必须实现有效逻辑
- `Load`
  - 把词条和候选写入 `Engine`

### 2.2 Engine 层

`pkg/ime/engine.go`

负责：

- 保存索引
- 执行查询
- 候选排序
- 尾部补全 / 组合 / 回退

它不关心 UI，不依赖 termbox。

### 2.3 Session 层

`pkg/ime/session.go`

负责：

- 保存当前输入 buffer
- 保存当前候选和选中位置
- 保存已经提交的文本

它是 TUI/宿主最应该直接使用的 API。

---

## 3. Source 接口与 registry

### 3.1 registry 机制

`pkg/ime/source.go`

当前已经支持 source registry：

```go
type SourceFactory func(SourceSpec) (Source, error)

func RegisterSource(kind string, factory SourceFactory)
func NewSource(spec SourceSpec) (Source, error)
func RegisteredSourceKinds() []string
```

含义：

- `RegisterSource`
  - 注册一种新的词库类型
- `NewSource`
  - 按 `kind=path` 创建对应 source
- `RegisteredSourceKinds`
  - 返回当前所有已注册种类

### 3.2 内建 source

当前内建了三种：

- `fcitx`
- `sogou`
- `sogou-syllables`

分别定义在：

- [pkg/ime/source_fcitx.go](/Users/xc/projects/In-Harmony-with-the-Dust/go/im/pkg/ime/source_fcitx.go:1)
- [pkg/ime/source_sogou.go](/Users/xc/projects/In-Harmony-with-the-Dust/go/im/pkg/ime/source_sogou.go:1)
- [pkg/ime/source_sogou_syllables.go](/Users/xc/projects/In-Harmony-with-the-Dust/go/im/pkg/ime/source_sogou_syllables.go:1)

每个文件都在 `init()` 中自注册。

### 3.3 默认 source 配置

默认加载顺序：

1. `sogou-syllables`
2. `fcitx`
3. `sogou`

原因：

- 先用 `sogou-syllables` 收集音节边界信息
- 再加载 `fcitx` 大词库候选
- 最后加载 `sogou` 的带边界词条和首字母缩写信息

---

## 4. Engine 内部数据结构

`Engine` 当前维护三类主要索引：

```go
exact           map[string][]Candidate
initials        map[string][]Candidate
derivedInitials map[string][]Candidate
```

### 4.1 exact

全拼精确索引。

例：

- `nihao -> 你好`
- `woaibeijing -> 我爱北京`

### 4.2 initials

显式首字母缩写索引。

主要来自 sogou 词库。

例：

- `bj -> 北京`
- `nh -> 你好`

### 4.3 derivedInitials

由切分结果自动推导的首字母索引。

例：

- `woaini -> wo'ai'ni -> wan`

这样即使词库没有显式提供缩写，也能支持一部分缩写查询。

### 4.4 syllables + segmenter

```go
syllables map[string]struct{}
segmenter *Segmenter
```

用于：

- 拼音切分
- 自动修复 `zaibz -> zai'b'z`
- 自动识别 `henbucuod -> hen'bu'cuo'd`

---

## 5. 查询流程

入口函数：

```go
func (e *Engine) Search(query string, limit int) SearchResult
```

一次查询大致按如下顺序进行。

### 5.1 Query 解析

位置：`pkg/ime/query.go`

分两类：

#### 1. 显式边界查询

用户自己输入 `'`

例：

- `zhen's'd`
- `zaib'z`

这类输入会先经过 `parseQueryPattern()`，保留结构信息。

#### 2. 自动结构化查询

用户没有输入 `'`，但 best-effort 切分能推导出边界。

例：

- `zaibz -> zai'b'z`
- `henbucuod -> hen'bu'cuo'd`

这类通过 `autoQueryPattern()` 生成。

### 5.2 切分

位置：`pkg/ime/segment.go`

当前支持：

- `Patterns(input, limit)`：
  - 返回所有切分方案
- `BestPattern(input)`：
  - 返回最佳切分
- `BestEffortParts(input)`：
  - 最重要的工程化能力
  - 在无法完整切分时，尽量切出前半部分合法音节，剩余部分按字母拆开

例：

- `woaibeijing -> wo'ai'bei'jing`
- `zaibz -> zai'b'z`
- `henbucuod -> hen'bu'cuo'd`

### 5.3 候选生成顺序

当前顺序是：

1. `exact_full`
2. `exact_initials`
3. `derived_initials`
4. `combined`
5. `pattern_prefix`
6. `tail_composed`
7. `tail_fallback`
8. `prefix_full`
9. `prefix_initials`

#### exact_full

完整 key 直接命中。

例：

- `henbucuo -> 很不错`

#### exact_initials / derived_initials

用于缩写。

#### combined

将切分后的多个完整子词组合起来。

例：

- `woaibeijing`
  - `woai -> 我爱`
  - `beijing -> 北京`
  - 组合得到 `我爱北京`

#### pattern_prefix

用于带 `'` 的结构化模糊匹配。

例：

- `zhen's'd`
  - 可命中 `zhen'shi'de`

#### tail_composed

用于“前面是完整词，最后一段是未完成音节或缩写”的情况。

例：

- `henbucuod`
  - `henbucuo -> 很不错`
  - `d -> 的 / 等 / 对 / 到`
  - 组合得到 `很不错的` 等

#### tail_fallback

如果尾段组合也失败，则回退到最近完整前缀词。

例：

- `henbucuod`
  - 若 `d` 没有合适候选，至少回退到 `很不错`

### 5.4 候选排序

位置：`candidateScore()`

当前是轻量启发式排序，不是用户词频模型。

主要因素：

- 命中类型优先级
- 是否多字词
- pattern 长度
- 原始词库中出现顺序

因此当前表现已经可用，但还不是成熟输入法级排序。

---

## 6. Session API

`pkg/ime/session.go`

这是给 TUI/宿主使用的稳定 API 层。

主要方法：

```go
session := ime.NewSession(engine, limit)

session.InputRune(r)
session.Backspace()
session.MoveSelection(delta)
session.CommitSelection()
session.CommitIndex(index)

session.Buffer()
session.Candidates()
session.Log()
session.CommittedText()
```

TUI 不应该直接操作 `Engine` 的内部索引。

---

## 7. 两个 demo 入口

### 7.1 CLI demo

文件：

- [cmd/im/main.go](/Users/xc/projects/In-Harmony-with-the-Dust/go/im/cmd/im/main.go:1)

作用：

- 输入一行拼音
- 输出候选和日志

### 7.2 TUI demo

文件：

- [cmd/im-tui/main.go](/Users/xc/projects/In-Harmony-with-the-Dust/go/im/cmd/im-tui/main.go:1)

依赖：

- `github.com/nsf/termbox-go`
- `github.com/mattn/go-runewidth`

作用：

- 输入 buffer
- 候选列表
- 已提交文本
- 日志信息

当前键位：

- 字母输入
- `Backspace`
- `Up/Down`
- `Enter` / `Space`
- `1-9`
- `Esc` / `Ctrl+C`

---

## 8. 入口参数与 source 组合

入口层复用了：

- [cmd/internal/sourceflags/sourceflags.go](/Users/xc/projects/In-Harmony-with-the-Dust/go/im/cmd/internal/sourceflags/sourceflags.go:1)

当前支持两种使用方式。

### 8.1 兼容模式

```bash
go run ./cmd/im
go run ./cmd/im-tui
```

默认等价于：

```text
-source sogou-syllables=data/dicts/sogou/vimim.pinyin.txt
-source fcitx=data/dicts/fcitx/vimim.pinyin.txt
-source sogou=data/dicts/sogou/vimim.pinyin.txt
```

### 8.2 显式 source 模式

```bash
go run ./cmd/im \
  -source sogou-syllables=data/dicts/sogou/vimim.pinyin.txt \
  -source fcitx=data/dicts/fcitx/vimim.pinyin.txt \
  -source sogou=data/dicts/sogou/vimim.pinyin.txt
```

### 8.3 `IM_DICT_DIR` 目录约定

当入口不传 `-source`，也不显式传 `-fcitx` / `-sogou` 时，`ime` 包会自己决定默认词库位置。

默认规则：

1. 如果设置了环境变量 `IM_DICT_DIR`
2. 则从这个目录下按固定层级寻找词库
3. 如果没有设置，则退回仓库内默认目录 `data/dicts/`

约定目录结构如下：

```text
$IM_DICT_DIR/
  fcitx/
    vimim.pinyin.txt
  sogou/
    vimim.pinyin.txt
```

对应关系：

- fcitx 默认路径：
  - `$IM_DICT_DIR/fcitx/vimim.pinyin.txt`
- sogou 默认路径：
  - `$IM_DICT_DIR/sogou/vimim.pinyin.txt`

示例：

```bash
IM_DICT_DIR=/opt/im-dicts go run ./cmd/im
```

则 `ime` 会尝试读取：

```text
/opt/im-dicts/fcitx/vimim.pinyin.txt
/opt/im-dicts/sogou/vimim.pinyin.txt
```

如果目录结构不符合这个约定，默认加载会失败。

如果你不想遵守这个目录结构，可以直接显式传 source：

```bash
go run ./cmd/im \
  -source sogou-syllables=/custom/path/a.txt \
  -source fcitx=/custom/path/b.txt \
  -source sogou=/custom/path/c.txt
```

推荐理解方式：

- `IM_DICT_DIR` 适合标准化部署
- `-source` 适合调试、测试、非标准目录布局

---

## 9. 如何新增一个词库

这是本文最重要的扩展指南。

目标：新增一个新词库格式，比如 `mydict`，让入口支持：

```bash
-source mydict=/path/to/mydict.txt
```

### 9.1 步骤总览

只需要做三件事：

1. 新建一个 source 文件
2. 实现 `Source` 接口
3. 在 `init()` 中注册它

不需要改：

- `cmd/im/main.go`
- `cmd/im-tui/main.go`
- `BuildSourceConfig`
- `NewSource`
- `Engine.Search`

### 9.2 新建文件

建议仿照现有命名：

```text
pkg/ime/source_mydict.go
```

### 9.3 实现接口

最小模板如下：

```go
package ime

import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

type MyDictSource struct {
    Path string
}

func init() {
    RegisterSource("mydict", func(spec SourceSpec) (Source, error) {
        return NewMyDictSource(spec.Path), nil
    })
}

func NewMyDictSource(path string) *MyDictSource {
    return &MyDictSource{Path: path}
}

func (s *MyDictSource) Name() string {
    return "mydict"
}

func (s *MyDictSource) CollectSyllables(engine *Engine) error {
    // 如果你的词库自带音节边界，可以在这里补充 syllables
    // 如果没有，可以直接 return nil
    return nil
}

func (s *MyDictSource) Load(engine *Engine) error {
    file, err := os.Open(s.Path)
    if err != nil {
        return fmt.Errorf("open mydict: %w", err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

    order := 0
    for scanner.Scan() {
        fields := strings.Fields(scanner.Text())
        if len(fields) < 2 {
            continue
        }

        key := normalizeAlphaOnly(fields[0])
        if key == "" {
            continue
        }

        bestPattern := engine.bestSegmentPattern(key)
        initials := initialsFromPattern(bestPattern)

        for _, word := range fields[1:] {
            candidate := Candidate{
                Word:    word,
                Key:     key,
                Source:  s.Name(),
                Order:   order,
                Pattern: bestPattern,
            }
            engine.addExact(key, candidate)
            if initials != "" {
                engine.addDerivedInitials(initials, candidate)
            }
            order++
        }
    }

    if err := scanner.Err(); err != nil {
        return fmt.Errorf("scan mydict: %w", err)
    }
    return nil
}
```

### 9.4 什么时候需要实现 CollectSyllables

如果你的词库像 sogou 一样自带：

- `zhen'shi'de`
- `wo'ai'ni`

这种显式音节边界，那么建议在 `CollectSyllables` 中把音节收进去。

如果词库只是连写：

- `zhenshide`
- `woaini`

而没有更强的音节信息，那么可以什么都不做，直接返回 `nil`。

### 9.5 入口如何使用

实现并注册以后，不用改入口代码，直接这样运行：

```bash
go run ./cmd/im -source mydict=/path/to/mydict.txt
```

如果它还依赖别的 source 提供音节表，也可以组合：

```bash
go run ./cmd/im \
  -source sogou-syllables=data/dicts/sogou/vimim.pinyin.txt \
  -source mydict=/path/to/mydict.txt
```

### 9.6 推荐做法

如果新词库：

- 有更好的音节边界信息
  - 在 `CollectSyllables` 利用它
- 有显式缩写信息
  - 在 `Load` 中调用 `engine.addInitials`
- 没有缩写信息
  - 保持 `engine.addExact` 即可，缩写可以交给 `derivedInitials`

### 9.7 最小检查清单

新增 source 后，至少验证：

1. `go test ./...`
2. `go build ./cmd/im ./cmd/im-tui`
3. 用 `-source mydict=...` 运行一次
4. 验证至少一条词命中

---

## 10. 当前实现的局限

当前实现已经可用，但仍有工程上明确的局限：

- 排序仍然是轻量启发式，不是成熟词频模型
- `combined`、`tail_composed` 会产生较多组合候选
- 没有用户词频学习
- 没有持久化用户词库
- 没有分页
- 没有系统级输入法集成

因此当前定位应理解为：

- 一个可运行、可扩展、可继续演进的内嵌输入法引擎
- 已支持 CLI 和 TUI demo
- 已具备 source 接口扩展能力

---

## 11. 推荐后续演进方向

### 11.1 体验方向

- 候选分页
- 更强的尾段单字排序
- 用户词频学习
- 更稳定的组合词排序

### 11.2 架构方向

- source 示例模板
- source 的基准测试
- 预编译二进制索引，降低启动时间
- 更严格的公共 API 文档

### 11.3 UI 方向

- TUI 翻页
- 高亮当前拼音切分
- 提交历史滚动
- 宿主程序嵌入接口

---

## 12. 当前推荐的理解方式

如果你要继续开发，建议把它理解成：

- `pkg/ime` 是唯一核心
- `Source` 是扩展词库的标准接口
- `Session` 是接 UI 的标准接口
- `cmd/im` 和 `cmd/im-tui` 只是 demo，而不是引擎本体

这也是后续继续扩展时最不容易失控的边界。
