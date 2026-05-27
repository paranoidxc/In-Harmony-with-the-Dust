# ParseAfterEdit 增量解析技术说明

本文记录当前 `syntax/parser.go` 中 `ParseAfterEdit` 的工作方式，以及它如何在编辑后复用 `ParseAll` 产生的解析结果。

## 目标

`ParseAll` 会从文档开头到结尾完整扫描一次，并把扫描结果保存成一组 `computation`：

```go
type computation struct {
    readLength     uint64
    consumedLength uint64
    startState     State
    endState       State
    tokens         []ComputedToken
}
```

`ParseAfterEdit` 的目标是：

1. 输入编辑后的新文档 `buf`。
2. 输入一次编辑描述 `Edit{Offset, NumInserted, NumDeleted}`。
3. 尽量复用旧的 `computation`。
4. 对受影响的区域重新运行语言解析函数。
5. 生成新的 `p.computations`，后续 `Tokens()` 继续返回全局 offset 的 token。

它当前没有使用树结构，而是用切片保存计算结果。逻辑更简单，便于理解和调试。

## Edit 的含义

```go
type Edit struct {
    Offset      uint64
    NumInserted uint64
    NumDeleted  uint64
}
```

- `Offset`：编辑发生在旧文档中的绝对字符 offset。
- `NumInserted`：新文档中插入了多少个 rune。
- `NumDeleted`：旧文档中删除了多少个 rune。

几种常见编辑：

```go
// 在 offset=13 插入 "return\n"
Edit{Offset: 13, NumInserted: 7}

// 从 offset=20 删除 5 个 rune
Edit{Offset: 20, NumDeleted: 5}

// 从 offset=30 把旧的 4 个 rune 替换成新的 6 个 rune
Edit{Offset: 30, NumInserted: 6, NumDeleted: 4}
```

换行也算 1 个 rune，因为 `Buf` 里每行不保存 `\n`，但 `TrackingRuneIter.NextRune()` 在行尾和下一行之间会读出一个 `\n`。

## ParseAll 留下了什么缓存

全量解析时，`ParseAll` 从 offset 0 开始循环调用 `runParseFunc`：

```go
result := p.runParseFunc(pos, buf, state)
```

每次得到一个 `Result`，再保存成 `computation`：

- `consumedLength`：本次解析真正消费的字符数。
- `readLength`：解析函数为了做判断实际读取过的最大字符数。
- `startState`：解析前的状态。
- `endState`：解析后的状态。
- `tokens`：本段内部的 token，offset 是相对本 computation 开头的。

例如：

```go
package main
func main() {}
```

可能形成类似这样的缓存：

```text
computation[0]
  old offset:      0
  consumedLength: 13      // "package main\n"
  readLength:     >= 13
  startState:     Empty
  endState:       Empty
  tokens:         keyword "package" at relative offset 0

computation[1]
  old offset:      13
  consumedLength: ...     // "func main() {}"
  readLength:     ...
  startState:     Empty
  endState:       Empty
  tokens:         keyword/operator tokens
```

实际切分受解析函数和 `minInitialConsumedLen` 合并策略影响，小段会合并成更大的 computation。

## readLength 为什么重要

`consumedLength` 表示“本段结果覆盖了多少字符”。

`readLength` 表示“解析这段时看过多少字符”。

两者不一定相等。比如识别标识符或关键字时，解析器可能会多看一个字符来确认 token 是否结束：

```text
输入: "package main"
解析 "package" 时：
  consumedLength = 7       // 消费 package
  readLength     = 8       // 可能额外看到了后面的空格
```

因此，如果编辑发生在一个 computation 的 `readLength` 范围内，即使不在 `consumedLength` 范围内，也可能改变解析结果。`ParseAfterEdit` 会把这种 computation 判定为不可复用。

## ParseAfterEdit 的整体流程

代码主流程：

```go
func (p *P) ParseAfterEdit(buf *Buf, edit Edit) {
    oldComputations := append([]computation(nil), p.computations...)
    oldOffsetByIndex := make([]uint64, len(oldComputations))

    // 1. 记录每个旧 computation 在旧文档中的起始 offset。
    var oldOffset uint64
    for i, c := range oldComputations {
        oldOffsetByIndex[i] = oldOffset
        oldOffset += c.consumedLength
    }

    // 2. 清空当前结果，从新文档开头重新构造 computation 列表。
    p.computations = p.computations[:0]
    state := State(EmptyState{})
    pos := Pos{Row: 0, Col: 0}
    var totalOffset uint64
    n := totalChars(buf)

    // 3. 从新文档开头向后扫描。
    for totalOffset < n {
        // 3a. 尝试复用旧 computation。
        if c, ok := reusableComputation(oldComputations, oldOffsetByIndex, edit, totalOffset, state); ok {
            p.computations = append(p.computations, c)
            totalOffset += c.consumedLength
            state = c.endState
            pos = advancePos(buf, pos, c.consumedLength)
            continue
        }

        // 3b. 不能复用，就在新文档当前位置重新解析。
        result := p.runParseFunc(pos, buf, state)
        ...
    }
}
```

关键点：它不是只从编辑点开始解析，而是从新文档开头重新“走一遍”。不过每一步都会先问：当前位置是否能直接拿旧 computation 过来？如果可以，就跳过整段；如果不可以，才重新解析当前段。

## 新 offset 如何映射回旧 offset

编辑后，同一个文本片段在新旧文档里的 offset 可能不同。

`oldOffsetAfterEdit` 做这个映射：

```go
func oldOffsetAfterEdit(edit Edit, newOffset uint64) (uint64, bool) {
    if newOffset < edit.Offset {
        return newOffset, true
    }
    if newOffset < edit.Offset+edit.NumInserted {
        return 0, false
    }
    return newOffset - edit.NumInserted + edit.NumDeleted, true
}
```

含义：

### 1. 新位置在编辑点之前

```text
newOffset < edit.Offset
```

编辑前后的 offset 一样：

```text
oldOffset = newOffset
```

### 2. 新位置落在插入出来的新内容内部

```text
edit.Offset <= newOffset < edit.Offset + NumInserted
```

这段内容旧文档里不存在，所以不能复用：

```text
ok = false
```

### 3. 新位置在编辑点之后

```text
newOffset >= edit.Offset + NumInserted
```

需要把新 offset 映射回旧 offset：

```text
oldOffset = newOffset - NumInserted + NumDeleted
```

例如旧文档 offset=100 之后插入 10 个字符：

```text
新文档 offset 130 对应旧文档 offset 120
oldOffset = 130 - 10 + 0 = 120
```

如果删除 8 个字符：

```text
新文档 offset 130 对应旧文档 offset 138
oldOffset = 130 - 0 + 8 = 138
```

## 什么情况下 computation 可以复用

`reusableComputation` 的逻辑：

```go
func reusableComputation(computations []computation, offsets []uint64, edit Edit, newOffset uint64, state State) (computation, bool) {
    oldOffset, ok := oldOffsetAfterEdit(edit, newOffset)
    if !ok {
        return computation{}, false
    }

    for i, c := range computations {
        if offsets[i] != oldOffset || c.startState != state {
            continue
        }
        if computationOverlapsEdit(offsets[i], c, edit) {
            return computation{}, false
        }
        return c, true
    }
    return computation{}, false
}
```

必须同时满足三点：

1. 新文档当前位置能映射到旧文档 offset。
2. 旧文档中恰好有一个 computation 从这个 oldOffset 开始。
3. 旧 computation 的 `startState` 等于当前解析状态。
4. 旧 computation 的读取范围没有被这次编辑影响。

其中第 3 点很重要：如果前面重新解析后产生了不同的 `endState`，后面的 computation 即使文本没变，也不能复用。因为语法状态会影响解析结果。

例如如果语言支持“跨段字符串状态”，前面插入一个未闭合字符串可能让后面所有内容都处于 string 状态，此时后面的旧 token 不能直接复用。

当前 Go parser 的 `State` 基本是 `EmptyState{}`，但接口已经为更复杂语言或跨段状态保留了能力。

## 如何判断编辑影响了旧 computation

```go
func computationOverlapsEdit(offset uint64, c computation, edit Edit) bool {
    readEnd := offset + c.readLength
    deleteEnd := edit.Offset + edit.NumDeleted
    if edit.NumDeleted > 0 && offset < deleteEnd && readEnd > edit.Offset {
        return true
    }
    if edit.NumInserted > 0 && offset < edit.Offset && readEnd > edit.Offset {
        return true
    }
    return false
}
```

### 删除影响判断

```text
旧 computation 读取范围: [offset, offset + readLength)
删除范围:              [edit.Offset, edit.Offset + NumDeleted)
```

只要两个范围相交，就不能复用：

```text
offset < deleteEnd && readEnd > edit.Offset
```

### 插入影响判断

```text
旧 computation 读取范围跨过插入点
```

如果旧 computation 开始在插入点之前，并且读取范围越过插入点：

```text
offset < edit.Offset && readEnd > edit.Offset
```

说明当初解析这段时看过插入点后面的字符。现在插入点处多了新内容，旧结果不再可靠。

注意：如果插入刚好发生在某个 computation 的开头：

```text
offset == edit.Offset
```

旧 computation 不会被视为“跨过插入点”。这是合理的，因为新插入内容会先被重新解析；解析走到插入内容之后时，旧 computation 的新位置会映射回这个 old offset，并可以尝试复用。

## 例子 1：在函数体里插入 return

旧文档：

```go
package main

func main() {
	println("hi")
}
```

新文档：

```go
package main

func main() {
	return
	println("hi")
}
```

编辑：

```go
inserted := "return\n\t"
edit := Edit{
    Offset:      len("package main\n\nfunc main() {\n\t"),
    NumInserted: len(inserted),
}
```

执行过程：

1. `ParseAfterEdit` 从新文档 offset 0 开始。
2. 编辑点之前的 computation，若读取范围没有跨过插入点，则复用。
3. 到达插入出的 `return\n\t` 内部时，`oldOffsetAfterEdit` 返回 `ok=false`，必须重新解析。
4. 插入内容解析结束后，新 offset 会映射到旧文档中 `println("hi")` 的位置。
5. 如果当前 `state` 与旧 computation 的 `startState` 一致，并且旧 computation 没被编辑影响，就复用后面的 `println` 和 `}`。
6. 最终 `Tokens()` 应与对新文档直接 `ParseAll` 完全一致。

对应测试：

```go
TestParseAfterEditMatchesParseAllForKeywordInsertion
```

## 例子 2：删除一整行

旧文档：

```go
package main

func main() {
	var x = 42
	println(x)
}
```

新文档：

```go
package main

func main() {
	println(x)
}
```

编辑：

```go
deleted := "\tvar x = 42\n"
edit := Edit{
    Offset:     len("package main\n\nfunc main() {\n"),
    NumDeleted: len(deleted),
}
```

过程：

1. 删除点之前的 computation 可以复用，前提是读取范围不与删除范围相交。
2. 删除范围内的旧 computation 全部不可复用。
3. 新文档走到 `println(x)` 时：

```text
newOffset -> oldOffset = newOffset + NumDeleted
```

4. 旧文档中删除段后面的 computation 可以被找到并复用。

对应测试：

```go
TestParseAfterEditMatchesParseAllForLineDeletion
```

## 例子 3：替换字符串 literal

旧文档：

```go
package main

const message = "hello"
```

新文档：

```go
package main

const message = `hello
world`
```

编辑：

```go
edit := Edit{
    Offset:      len("package main\n\nconst message = "),
    NumInserted: len("`hello\nworld`"),
    NumDeleted:  len("\"hello\""),
}
```

这是“删除 + 插入”的替换。

过程：

1. 替换点之前的 computation 尝试复用。
2. 旧字符串所在 computation 与删除范围相交，所以不能复用。
3. 新 raw string 是新插入内容，也不能从旧文档映射，所以重新解析。
4. raw string 结束之后，后续内容可继续映射和复用。

对应测试：

```go
TestParseAfterEditMatchesParseAllForStringReplacement
```

## 例子 4：跨行 block comment 插入与删除

插入跨行注释：

```go
oldText := "package main\n\nfunc main() {}\n"
inserted := "/*\nline one\nline two\n*/\n"
newText := "package main\n\n" + inserted + "func main() {}\n"
edit := Edit{
    Offset:      len("package main\n\n"),
    NumInserted: len(inserted),
}
```

删除跨行注释：

```go
comment := "/*\nremove me\nacross lines\n*/\n"
oldText := "package main\n\n" + comment + "func main() {}\n"
newText := "package main\n\nfunc main() {}\n"
edit := Edit{
    Offset:     len("package main\n\n"),
    NumDeleted: len(comment),
}
```

这些测试覆盖了跨多行的 `NumInserted` / `NumDeleted`。对 `ParseAfterEdit` 来说，跨行并没有特殊分支：只要 `Offset`、插入长度、删除长度按 rune 计数正确，换行就是普通的 1 个 rune。

对应测试：

```go
TestParseAfterEditMultilineBlockCommentInsertion
TestParseAfterEditMultilineBlockCommentDeletion
```

## 例子 5：大文件中局部编辑

测试中 `makeLargeGoSource(70)` 会生成 300+ 行 Go 源码。

在大文件中追加函数：

```go
oldText := makeLargeGoSource(70)
inserted := "func appended() {\n\tprintln(\"appended\")\n}\n"
newText := oldText + inserted
edit := Edit{
    Offset:      len(oldText),
    NumInserted: len(inserted),
}
```

这个场景下，大部分旧 computation 都在编辑点之前，且读取范围没有跨过 EOF 插入点，因此可以复用。只有追加的新函数需要重新解析。

大文件中删除一段函数：

```go
oldText := makeLargeGoSource(90)
prefix := makeLargeGoSource(20)
deletedStart := len(prefix)
deleted := oldText[deletedStart:len(makeLargeGoSource(50))]
newText := oldText[:deletedStart] + oldText[deletedStart+len(deleted):]
edit := Edit{
    Offset:     uint64(deletedStart),
    NumDeleted: uint64(len(deleted)),
}
```

删除区域之前的 computation 可复用；删除区域之后的 computation 会通过 `newOffset + NumDeleted` 映射回旧 offset 后继续复用。

对应测试：

```go
TestParseAfterEditLargeSourceAppendFunction
TestParseAfterEditLargeSourceDeleteManyFunctions
```

## Tokens 如何保持全局 offset

每个 computation 内部 token 的 `Offset` 是相对 computation 起点的。

`Tokens()` 会累加 computation 的 `consumedLength`，把相对 offset 转成文档全局 offset：

```go
func (p *P) Tokens() []ComputedToken {
    var tokens []ComputedToken
    var offset uint64
    for _, c := range p.computations {
        for _, tok := range c.tokens {
            tokens = append(tokens, ComputedToken{
                Offset: offset + tok.Offset,
                Length: tok.Length,
                Role:   tok.Role,
            })
        }
        offset += c.consumedLength
    }
    return tokens
}
```

因此不管 computation 是复用的还是重新解析的，只要 `p.computations` 在新文档中按顺序排列，`Tokens()` 就能给出新文档中的正确 offset。

## 正确性条件

调用 `ParseAfterEdit` 时必须保证：

1. `buf` 已经是编辑后的新文档。
2. `edit.Offset` 是旧文档中的 offset。
3. `edit.NumInserted` 是新插入内容的 rune 数。
4. `edit.NumDeleted` 是旧文档被删除内容的 rune 数。
5. `p.computations` 来自同一个文档的上一次 `ParseAll` 或 `ParseAfterEdit`。

如果这些条件不满足，offset 映射会错误，旧 computation 可能被错误复用。

## 当前实现的特点与限制

优点：

- 实现简单：用切片而不是树。
- 容易调试：每个 computation 有清晰的起点、长度、状态、token。
- 对追加、插入、删除、替换、大文件局部编辑都能复用未受影响区域。
- 测试用 `ParseAfterEdit` 结果对比新文档 `ParseAll` 结果，验证行为一致。

限制：

- 查找可复用 computation 当前是线性扫描：

```go
for i, c := range computations { ... }
```

  对超大文件，频繁编辑时性能不如树或按 offset 二分查找。

- `ParseAfterEdit` 仍从新文档开头推进，只是在推进过程中复用旧结果。它不是从编辑点直接开始。

- 当前 Go parser 的 `State` 很简单，尚未充分体现复杂语言中的跨段状态传播。不过 `startState/endState` 已经在结构上支持这种场景。

## 一句话总结

`ParseAfterEdit` 的核心是：把新文档 offset 映射回旧文档 offset，找到同位置、同起始状态、且读取范围没有被编辑影响的旧 computation 直接复用；否则就在新文档当前位置重新解析。最终重新构造一份面向新文档的 computation 切片。