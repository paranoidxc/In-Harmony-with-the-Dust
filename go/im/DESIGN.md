# 内嵌拼音输入法 — 技术设计文档

## 1. 词库格式分析

两份词库格式相似，每行一条记录：`拼音 候选词1 候选词2 ...`

**fcitx 词库** (290943 行)：拼音连写，无音节分隔符
```
nihao 你好
shi 是 时 事 什 实 十 ...
```

**sogou 词库** (171707 行)：音节间用 `'` 分隔，支持首字母缩写
```
ni'hao 你好 拟好 倪浩
a'b    阿布 阿扁 阿宝    ← 首字母缩写
```

两种格式都可以直接解析为 `map[string][]string`。fcitx 词库更大（含带声调编码如 `a1`、`shi4`），sogou 词库有音节边界信息和缩写，两者互补。

---

## 2. 拼音输入法核心流程

用户输入一串字母（如 `woaini`），输入法需要完成三件事：

```
原始输入 → 音节切分 → 词库查找 → 候选排序 → 展示选词
```

### 2.1 音节切分（Segmentation）

这是拼音输入法最关键的一步。用户输入连续字母，需要切分成合法的拼音音节组合。

**问题**：`xian` 可以是 `xi'an`（西安）也可以是 `xian`（先/现）。

**合法音节表**：约 400 个（不含声调），是一个固定集合：
```
a ai an ang ao
ba bai ban bang bao bei ben beng bi bian biao bie bin bing bo bu
ca cai can cang cao ce cen ceng cha chai chan chang chao che chen cheng chi chong chou chu chua chuai chuan chuang chui chun chuo ci cong cou cu cuan cui cun cuo
...
```

**切分算法**：对输入串进行全切分（找出所有合法的音节序列），核心是 **动态规划 / 回溯搜索**。

```go
type Segmentation struct {
    syllables []string   // 切分结果，如 ["wo", "ai", "ni"]
    remainder string     // 未能匹配的尾部
}

// 返回所有可能的切分方案
func Segment(input string) []Segmentation
```

**贪心 vs 全切分**：

- 贪心（最长匹配）：从左到右每次取最长合法音节。简单快速，但 `xian` 只会得到 `xian` 而丢失 `xi'an`。
- **全切分 + 排序**（推荐）：找出所有合法切分，按优先级排序。用户通过选择不同切分方案来消歧。

**实现要点**：
- 用 Trie 存储合法音节表，快速判断前缀是否合法
- DFS/DP 遍历所有切分路径
- 常见的音节数量少（~400），性能不是瓶颈

### 2.2 候选词查找

每个切分方案对应一组音节，需要在词库中查找匹配的词。

**查找策略**（以输入 `woaini` 切分为 `wo'ai'ni` 为例）：

1. **完整匹配**：查找 `woaini` → `我爱你`
2. **前缀匹配**：查找 `wo` → `我 握 窝 ...`，然后 `woai` → `我爱 ...`
3. **逐字组合**：对每个音节分别查找单字，组合展示

**数据结构 — Trie 树**：

词库需要支持前缀查询（用户边输入边匹配），Trie 是标准选择。

```go
type TrieNode struct {
    children  map[byte]*TrieNode
    candidates []Candidate      // 该拼音对应的候选词
}

type Candidate struct {
    word      string
    frequency int              // 词频，用于排序
}
```

**查找顺序**：优先展示长词（整句匹配 > 词组 > 单字），这是所有主流输入法的做法。

```
输入: nihao
候选: 1.你好  2.你  3.泥  4.尼  5.倪  ...
       ^^^^   ^^^^^^^^^^^^^^^^^^^^^^^^^^
       整词匹配   首音节单字匹配
```

### 2.3 候选排序

排序决定了用户体验。主流输入法的排序依据：

| 因素 | 权重 | 说明 |
|------|------|------|
| 词频（静态） | 基础 | 词库中的默认排列顺序即为词频序 |
| 用户词频 | 高 | 用户选过的词提升权重 |
| 词长优先 | 高 | 长词优先于短词（整句 > 词组 > 单字） |
| 最近使用 | 中 | 最近选过的词短期提权 |

**最简可行方案**：词库本身已按词频排序（第一个候选词是最常用的），直接使用即可。后续再加用户词频学习。

### 2.4 选词与上屏

```
┌─────────────────────────────────────────┐
│ 输入: nihao                              │
│ ──────────────────────────────────        │
│ 1.你好  2.你  3.泥  4.尼  5.倪           │
│         ← Page 1/3 →                     │
└─────────────────────────────────────────┘
```

- 按数字键 1-9 选词（或空格选第一个）
- 选词后，该词"上屏"（输出到文本），剩余拼音继续匹配
- `-/=` 或 `[/]` 翻页

---

## 3. 主流输入法的选词机制详解

### 3.1 逐步消费模型

这是所有拼音输入法的核心交互模型：

```
状态: 输入缓冲区 = "woaibeijing"
切分: wo'ai'bei'jing

候选: 1.我爱北京  2.我爱  3.我  ...

用户按 1 → 上屏"我爱北京"，缓冲区清空
用户按 2 → 上屏"我爱"，缓冲区剩余 "beijing"，继续匹配
用户按 3 → 上屏"我"，缓冲区剩余 "aibeijing"，继续匹配
```

每次选词消费掉对应音节数量的拼音，剩余部分自动重新匹配。这就是"逐步消费"。

### 3.2 部分输入与模糊匹配

用户经常不输入完整拼音：

- `bj` → 北京（首字母缩写）
- `beij` → 北京（不完整音节，`j` 是 `jing` 的前缀）

**处理方式**：输入尾部如果不构成完整音节，视为某个音节的前缀，用前缀匹配查找。

```go
// "beij" → 切分为 "bei" + "j"
// "j" 不是合法音节，但是 ji/jia/jian/jing/... 的前缀
// 在 Trie 中查找所有以 "bei" 开头、第二音节以 "j" 开头的词
```

### 3.3 光标与编辑

用户可能打到一半发现打错了：

- `Backspace`：删除末尾字母
- `←/→`：在拼音串中移动光标（高级功能，可选）
- `Escape`：清空输入缓冲区

---

## 4. 数据结构设计

### 4.1 核心结构

```go
// 合法音节判断
type PinyinTable struct {
    trie *Trie   // 存储所有合法拼音音节（~400个）
}

// 词库索引
type Dictionary struct {
    trie *Trie   // key=拼音（连写），value=[]Candidate
}

// 输入法引擎
type Engine struct {
    dict        *Dictionary
    pinyinTable *PinyinTable
    userFreq    map[string]int   // 用户词频学习
}

// 输入会话状态
type Session struct {
    engine    *Engine
    buffer    []byte            // 用户输入的原始字母
    segments  []Segmentation    // 当前切分方案
    page      int               // 当前候选页
}
```

### 4.2 Trie 实现

由于拼音只含 a-z 共 26 个字符，可以用数组替代 map，提升缓存命中率：

```go
type TrieNode struct {
    children   [26]*TrieNode
    candidates []string        // 非 nil 表示这是一个合法拼音的终点
}
```

### 4.3 词库加载

```go
func LoadDict(path string) *Dictionary {
    // 逐行读取: "nihao 你好 倪浩"
    // 拆分: pinyin="nihao", words=["你好", "倪浩"]
    // 插入 Trie: trie.Insert("nihao", words)
}
```

fcitx 词库约 29 万行，每行平均几个候选词，内存占用约 50-100MB（含 Trie 结构开销）。可接受。

---

## 5. 核心算法

### 5.1 音节切分（DAG 最短路径）

主流实现方式是构建有向无环图（DAG），然后找最优路径：

```go
// 构建 DAG：对于输入串的每个位置 i，找出所有从 i 开始的合法音节
func BuildDAG(input string, table *PinyinTable) map[int][]int {
    dag := make(map[int][]int)
    for i := 0; i < len(input); i++ {
        for j := i + 1; j <= len(input); j++ {
            if table.IsValid(input[i:j]) {
                dag[i] = append(dag[i], j)
            }
        }
    }
    return dag
}

// DAG 上做 DFS，枚举所有路径即所有切分方案
func EnumPaths(dag map[int][]int, n int) [][]string
```

**示例**：输入 `xian`

```
位置: 0  1  2  3  4
      x  i  a  n
DAG:
  0 → 2 (xi)
  0 → 4 (xian)
  2 → 3 (a)
  2 → 4 (an)
  3 → 4 (n) ← 不合法，跳过

路径:
  xi + an  → ["xi", "an"]   → 西安
  xian     → ["xian"]       → 先/现
```

### 5.2 候选生成

```go
func (e *Engine) GetCandidates(syllables []string) []Candidate {
    var result []Candidate

    // 1. 整体匹配：尝试所有音节拼接后的词
    joined := strings.Join(syllables, "")
    if words := e.dict.Lookup(joined); words != nil {
        result = append(result, words...)
    }

    // 2. 前缀匹配：从第一个音节开始，逐步拼接查找
    //    找到的词长度越长优先级越高
    for i := len(syllables) - 1; i >= 1; i-- {
        partial := strings.Join(syllables[:i], "")
        if words := e.dict.Lookup(partial); words != nil {
            result = append(result, words...)
        }
    }

    // 3. 首音节单字
    if words := e.dict.Lookup(syllables[0]); words != nil {
        result = append(result, words...)
    }

    return dedupAndRank(result)
}
```

### 5.3 选词与消费

```go
func (s *Session) Select(index int) string {
    candidate := s.candidates[index]
    
    // 计算该候选词消费了几个音节
    consumed := candidate.syllableCount
    
    // 上屏
    output := candidate.word
    
    // 从 buffer 中移除已消费的部分
    consumedLen := totalRuneLen(s.segments[:consumed])
    s.buffer = s.buffer[consumedLen:]
    
    // 用剩余的 buffer 重新切分和查找
    s.refresh()
    
    // 记录用户词频
    s.engine.userFreq[candidate.word]++
    
    return output
}
```

---

## 6. 与 TUI / SDL 的集成

### 6.1 架构分层

```
┌─────────────────────┐
│   UI 层 (TUI/SDL)   │  ← 渲染候选框、处理按键
├─────────────────────┤
│   Session (会话)     │  ← 管理输入状态、光标、翻页
├─────────────────────┤
│   Engine (引擎)      │  ← 切分、查找、排序（纯逻辑，无 IO）
├─────────────────────┤
│   Dictionary (词库)  │  ← Trie 索引、词频数据
└─────────────────────┘
```

Engine 层完全无状态、无 IO 依赖，可独立测试。Session 管理一次输入会话的状态。UI 层只负责渲染和按键分发。

### 6.2 按键处理流程

```go
func (s *Session) HandleKey(key rune) (action Action) {
    switch {
    case key >= 'a' && key <= 'z':
        s.buffer = append(s.buffer, byte(key))
        s.refresh()                    // 重新切分 + 查找
        return ShowCandidates

    case key >= '1' && key <= '9':
        word := s.Select(int(key - '1'))
        return Commit(word)            // 上屏

    case key == ' ':
        word := s.Select(0)            // 空格 = 选第一个
        return Commit(word)

    case key == Backspace:
        s.buffer = s.buffer[:len(s.buffer)-1]
        s.refresh()
        return ShowCandidates

    case key == Escape:
        s.Clear()
        return HideCandidates

    case key == '=' || key == '.':
        s.page++
        return ShowCandidates

    case key == '-' || key == ',':
        s.page--
        return ShowCandidates
    }
}
```

### 6.3 TUI 方案 (tcell/bubbletea)

在 TUI 场景下，输入法候选框可以渲染为文本行：

```
拼音: xi'an
┌──────────────────────────────┐
│ 1.西安  2.先  3.现  4.闲  5.弦 │
└──────────────────────────────┘
```

用 bubbletea 时，输入法是一个独立的 `tea.Model`，宿主组件在中文模式下将按键事件转发给它。

### 6.4 SDL 方案

SDL 场景下需要自己绘制候选框浮窗：

- 在光标位置下方绘制一个矩形背景
- 用 TTF 字体渲染候选词文本
- 高亮当前选中项

SDL 的优势是可以精确控制渲染位置和样式，适合做类似系统输入法的浮窗体验。

---

## 7. 最小可行产品（MVP）路径

按优先级分阶段实现：

### Phase 1：基础可用
- [ ] 解析 fcitx 词库 → `map[string][]string`
- [ ] 实现合法音节表 + Trie
- [ ] 实现音节切分（DAG + DFS，取最优切分）
- [ ] 实现候选查找（整词 > 前缀词 > 单字）
- [ ] 实现 Session 按键处理（输入、选词、退格、清空）
- [ ] TUI 或 SDL 渲染候选框

### Phase 2：体验提升
- [ ] 翻页
- [ ] 不完整音节的前缀匹配（`beij` → 北京）
- [ ] 多切分方案切换（`xian` → 西安 / 先）
- [ ] 用户词频学习（选过的词提权，持久化到文件）

### Phase 3：进阶
- [ ] 首字母缩写输入（`bj` → 北京）
- [ ] 整句输入（连续选词，中间不断句）
- [ ] 中英文切换（Shift / Ctrl+Space）
- [ ] 模糊音支持（`zh/z`、`sh/s`、`ch/c` 不区分）

---

## 8. 关键设计决策

| 决策点 | 建议 | 理由 |
|--------|------|------|
| 词库选择 | fcitx 为主 | 更大更全，连写格式更简单 |
| 数据结构 | Trie | 前缀查找是核心需求 |
| 切分算法 | DAG 全切分 | 覆盖歧义，避免贪心丢解 |
| 排序 | 词库原始顺序 + 用户词频 | MVP 阶段足够 |
| 候选展示 | 长词优先 | 符合用户预期 |
| TUI vs SDL | 取决于宿主应用 | Engine 层与 UI 无关，两者可共用 |
