# Aretext 语法高亮系统实现详解

## 目录
1. [概述](#概述)
2. [核心架构](#核心架构)
3. [全量语法解析（ParseAll）](#全量语法解析parsall)
4. [增量解析算法](#增量解析算法)
5. [数据结构详解](#数据结构详解)
6. [完整示例](#完整示例超详细版)
7. [性能优化](#性能优化)
8. [解析函数示例](#解析函数示例)

---

## 概述

Aretext 的语法高亮系统采用了一种**增量解析（Incremental Parsing）**的设计，核心思想是：
- 初次解析整个文档
- 文档编辑后，**只重新解析受影响的部分**
- **重用未受影响部分的解析结果**

这种设计使得大文件编辑时的语法高亮更新非常高效。

---

## 核心架构

### 1. 主要组件

```
Parser (P)
├── parseFunc: 语言特定的解析函数
└── lastComputation: 缓存的解析结果树
    ├── tokens: 识别的语法标记
    ├── leftChild: 左子树
    └── rightChild: 右子树
```

### 2. 关键类型

#### Token（语法标记）
```go
type Token struct {
    Role     TokenRole  // 标记角色（关键字、字符串、注释等）
    StartPos uint64     // 起始位置
    EndPos   uint64     // 结束位置
}
```

#### Result（解析结果）
```go
type Result struct {
    NumConsumed    uint64           // 消费的字符数
    ComputedTokens []ComputedToken  // 识别的标记
    NextState      State            // 下一个解析状态
}
```

#### Computation（计算节点）
```go
type computation struct {
    readLength     uint64        // 读取的字符数
    consumedLength uint64        // 消费的字符数
    treeHeight     uint64        // AVL树高度
    startState     State         // 起始状态
    endState       State         // 结束状态
    tokens         []ComputedToken  // 叶子节点的标记
    leftChild      *computation  // 左子树
    rightChild     *computation  // 右子树
}
```

---

## 增量解析算法

### 核心思想

增量解析的关键在于：**如何判断哪些解析结果可以重用？**

答案：如果一个计算节点满足以下条件，就可以重用：
1. **读取范围**没有被编辑影响
2. **起始状态**与当前解析状态匹配

### 算法流程

#### 1. 初始解析 (ParseAll)

```go
func (p *P) ParseAll(tree *text.Tree) {
    var pos uint64
    state := EmptyState{}
    leafComputations := []computation{}
    
    // 从头到尾解析文档
    for pos < tree.NumChars() {
        c := p.runParseFunc(tree, pos, state)
        pos += c.ConsumedLength()
        state = c.EndState()
        leafComputations = append(leafComputations, c)
    }
    
    // 将叶子节点组合成平衡树
    p.lastComputation = concatLeafComputations(leafComputations)
}
```

**示例：** 解析 `"foo" "bar"`

```
初始文本: "foo" "bar"
         ↓
解析过程:
pos=0: 识别 "foo" → Token{0, 5, String}
pos=5: 识别 " "   → 无标记
pos=6: 识别 "bar" → Token{6, 11, String}

结果树:
        [Root]
       /      \
   [0-5]      [6-11]
   "foo"      "bar"
```

#### 2. 增量重解析 (ReparseAfterEdit)

```go
func (p *P) ReparseAfterEdit(tree *text.Tree, edit Edit) {
    var pos uint64
    var c *computation
    state := EmptyState{}
    
    for pos < tree.NumChars() {
        // 尝试找到可重用的计算节点
        nextComputation := p.findReusableComputation(pos, edit, state)
        
        if nextComputation == nil {
            // 没有可重用的，重新解析
            nextComputation = p.runParseFunc(tree, pos, state)
        }
        
        state = nextComputation.EndState()
        pos += nextComputation.ConsumedLength()
        c = c.Append(nextComputation)
    }
    
    p.lastComputation = c
}
```

### 重用判断逻辑

```go
func (p *P) findReusableComputation(pos uint64, edit Edit, state State) *computation {
    if pos < edit.pos {
        // 情况1: 解析位置在编辑位置之前
        // 查找从pos到edit.pos之间的子计算
        return p.lastComputation.LargestMatchingSubComputation(
            pos, edit.pos, state)
    }
    
    if edit.numInserted > 0 && pos >= edit.pos+edit.numInserted {
        // 情况2: 解析位置在插入内容之后
        // 将位置映射回旧文档（减去插入的字符数）
        return p.lastComputation.LargestMatchingSubComputation(
            pos-edit.numInserted, math.MaxUint64, state)
    }
    
    if edit.numDeleted > 0 && pos >= edit.pos {
        // 情况3: 解析位置在删除位置之后
        // 将位置映射回旧文档（加上删除的字符数）
        return p.lastComputation.LargestMatchingSubComputation(
            pos+edit.numDeleted, math.MaxUint64, state)
    }
    
    // 情况4: 解析位置在编辑范围内，无法重用
    return nil
}
```

---

## 完整示例（超详细版）

### 示例 1: 插入字符导致标记分裂

#### 场景说明
我们有一个完整的字符串 `"this is a test"`，现在要在中间插入一个引号，把它分成两个字符串。

#### 第一步：初始解析

**文本内容:**
```
位置: 0123456789...15
文本: "this is a test"
      ^              ^
      开始           结束
```

**解析结果:**
- 解析器从位置 0 开始
- 看到 `"` 开始一个字符串
- 一直读到位置 15 的 `"` 结束
- 生成标记: `Token{StartPos: 0, EndPos: 16, Role: String}`

**计算树结构:**
```
computation {
    readLength: 16      // 读取了16个字符
    consumedLength: 16  // 消费了16个字符
    startState: Empty
    endState: Empty
    tokens: [Token{0, 16, String}]
}
```

#### 第二步：编辑操作

**在位置 5 插入 `"`**

```
旧文本: "this is a test"
        01234567...
新文本: "this" is a test"
        012345678...
             ↑
          插入位置5
```

**编辑信息:**
```go
Edit {
    pos: 5,           // 插入位置
    numInserted: 1,   // 插入1个字符
    numDeleted: 0     // 没有删除
}
```

#### 第三步：增量重解析

**重解析循环开始:**

---

**迭代 1: pos=0, state=Empty**

1. **调用 findReusableComputation(0, edit, Empty)**
   
   ```
   当前位置 pos=0 < edit.pos=5
   → 进入情况1: 解析位置在编辑位置之前
   → 调用 LargestMatchingSubComputation(0, 5, Empty)
   ```

2. **查找可重用的计算:**
   
   ```
   旧计算树: [0-16]
   - readLength: 16
   - 读取范围: [0, 16)
   
   要求:
   - 起始位置: 0 ✓
   - 结束位置: ≤ 5 ✗ (实际是16，超过了编辑位置)
   - 起始状态: Empty ✓
   
   结论: 不能重用！因为这个计算读取了编辑位置之后的内容
   ```

3. **重新解析:**
   
   ```
   从位置0开始解析新文本: "this" is a test"
   
   解析器看到:
   - 位置0: '"' → 开始字符串
   - 位置1-4: 'this'
   - 位置5: '"' → 结束字符串！
   
   生成标记: Token{0, 6, String}  // "this"
   消费长度: 6
   ```

4. **更新状态:**
   ```
   pos = 0 + 6 = 6
   state = Empty
   ```

---

**迭代 2: pos=6, state=Empty**

1. **调用 findReusableComputation(6, edit, Empty)**
   
   ```
   当前位置 pos=6 >= edit.pos + edit.numInserted (5+1=6)
   → 进入情况2: 解析位置在插入内容之后
   → 映射到旧文档: oldPos = 6 - 1 = 5
   → 调用 LargestMatchingSubComputation(5, MaxUint64, Empty)
   ```

2. **查找可重用的计算:**
   
   ```
   旧计算树: [0-16]
   - 从位置5开始查找
   - 但旧文档位置5的内容是 " is a test"
   - 新文档位置6的内容是 " is a test"
   
   问题: 虽然内容看起来一样，但是：
   - 旧计算的起始位置是0，不是5
   - 无法找到从位置5开始的子计算
   
   结论: 不能重用！
   ```

3. **重新解析:**
   
   ```
   从位置6开始解析新文本:  is a test"
   
   解析器看到:
   - 位置6: ' ' → 不是引号，消费到下一个引号或EOF
   - 位置7-16: ' is a test"'
   - 到达EOF
   
   没有生成标记（因为不是完整的字符串）
   消费长度: 11
   ```

4. **更新状态:**
   ```
   pos = 6 + 11 = 17
   state = Empty
   ```

---

**循环结束:** pos=17 >= 文档长度17

#### 第四步：最终结果

**新的计算树:**
```
        [Root]
       /      \
   [0-6]      [6-17]
   Token:     无Token
   "this"
```

**标记列表:**
```
[Token{StartPos: 0, EndPos: 6, Role: String}]
```

**可视化对比:**
```
旧文本: "this is a test"
        ^^^^^^^^^^^^^^^^
        一个完整字符串

新文本: "this" is a test"
        ^^^^^^ ^^^^^^^^^^^
        字符串  普通文本
```

### 示例 2: 插入不影响远处的标记（重用机制的威力）

#### 场景说明
我们有三个独立的字符串，在中间的字符串里插入一个字符。这个例子展示了增量解析如何**只重新解析受影响的部分**，而**重用未受影响的部分**。

#### 第一步：初始解析

**文本内容:**
```
位置: 0    5    10   15   17
文本: "foo" "bar" "baz"
      ^   ^ ^   ^ ^   ^
      |   | |   | |   |
      标记1 标记2 标记3
```

**解析过程:**
```
pos=0: 解析 "foo"  → Token{0, 5, String}, 消费5个字符
pos=5: 解析 " "    → 无标记, 消费1个字符  
pos=6: 解析 "bar"  → Token{6, 11, String}, 消费5个字符
pos=11: 解析 " "   → 无标记, 消费1个字符
pos=12: 解析 "baz" → Token{12, 17, String}, 消费5个字符
```

**计算树结构:**
```
                    [Root: 0-17]
                   /            \
          [Left: 0-11]        [Right: 12-17]
          /          \              |
    [0-5: "foo"]  [6-11: "bar"]  [12-17: "baz"]
    
每个节点的详细信息:
[0-5]:
  - readLength: 5
  - consumedLength: 5
  - tokens: [Token{0, 5, String}]
  
[6-11]:
  - readLength: 5  
  - consumedLength: 5
  - tokens: [Token{0, 5, String}]  // 注意：offset相对于节点起始位置
  
[12-17]:
  - readLength: 5
  - consumedLength: 5
  - tokens: [Token{0, 5, String}]
```

#### 第二步：编辑操作

**在位置 7 插入 `x`**

```
旧文本: "foo" "bar" "baz"
        0123456789...
新文本: "foo" "bxar" "baz"
        0123456789...
               ↑
            插入位置7
```

**编辑信息:**
```go
Edit {
    pos: 7,           // 插入位置
    numInserted: 1,   // 插入1个字符
    numDeleted: 0     // 没有删除
}
```

**关键问题:** 哪些计算节点会受影响？
```
[0-5: "foo"]   → 读取范围 [0, 5)   → 不包含位置7 → ✅ 可以重用
[6-11: "bar"]  → 读取范围 [6, 11)  → 包含位置7   → ❌ 不能重用
[12-17: "baz"] → 读取范围 [12, 17) → 不包含位置7 → ✅ 可以重用（需要偏移）
```

#### 第三步：增量重解析

---

**迭代 1: pos=0, state=Empty**

1. **调用 findReusableComputation(0, edit, Empty)**
   
   ```
   当前位置 pos=0 < edit.pos=7
   → 进入情况1: 解析位置在编辑位置之前
   → 调用 LargestMatchingSubComputation(0, 7, Empty)
   ```

2. **在计算树中查找:**
   
   ```
   从根节点 [0-17] 开始:
   - 起始位置: 0 ✓
   - 读取范围: [0, 17) → 超过了7 ✗
   
   递归到左子树 [0-11]:
   - 起始位置: 0 ✓
   - 读取范围: [0, 11) → 超过了7 ✗
   
   继续递归到 [0-5]:
   - 起始位置: 0 ✓
   - 读取范围: [0, 5) → 没超过7 ✓
   - 起始状态: Empty ✓
   
   找到了！返回 [0-5]
   ```

3. **重用计算 [0-5]:**
   
   ```
   ✅ 直接重用，不需要重新解析！
   
   重用的标记: Token{0, 5, String}  // "foo"
   消费长度: 5
   ```

4. **更新状态:**
   ```
   pos = 0 + 5 = 5
   state = Empty
   ```

---

**迭代 2: pos=5, state=Empty**

1. **调用 findReusableComputation(5, edit, Empty)**
   
   ```
   当前位置 pos=5 < edit.pos=7
   → 进入情况1: 解析位置在编辑位置之前
   → 调用 LargestMatchingSubComputation(5, 7, Empty)
   ```

2. **在计算树中查找:**
   
   ```
   查找从位置5开始、不超过位置7的计算
   
   旧文档位置5是空格 " "
   旧计算树中没有从位置5开始的节点
   （因为空格被合并到了 [6-11] 之前）
   
   结论: 找不到可重用的计算
   ```

3. **重新解析:**
   
   ```
   从位置5开始解析新文本:  "bxar" "baz"
   
   解析器看到:
   - 位置5: ' ' → 不是引号
   - 继续读到位置6: '"' → 停止
   
   消费: 1个字符（空格）
   无标记
   ```

4. **更新状态:**
   ```
   pos = 5 + 1 = 6
   state = Empty
   ```

---

**迭代 3: pos=6, state=Empty**

1. **调用 findReusableComputation(6, edit, Empty)**
   
   ```
   当前位置 pos=6 < edit.pos=7
   → 进入情况1: 解析位置在编辑位置之前
   → 调用 LargestMatchingSubComputation(6, 7, Empty)
   ```

2. **在计算树中查找:**
   
   ```
   查找从位置6开始、不超过位置7的计算
   
   [6-11] 节点:
   - 起始位置: 6 ✓
   - 读取范围: [6, 11) → 超过了7 ✗
   
   结论: 找不到可重用的计算
   ```

3. **重新解析:**
   
   ```
   从位置6开始解析新文本: "bxar" "baz"
   
   解析器看到:
   - 位置6: '"' → 开始字符串
   - 位置7-10: 'bxar'
   - 位置11: '"' → 结束字符串
   
   生成标记: Token{0, 6, String}  // 相对offset=0, length=6
   实际位置: Token{6, 12, String}  // "bxar"
   消费长度: 6
   ```

4. **更新状态:**
   ```
   pos = 6 + 6 = 12
   state = Empty
   ```

---

**迭代 4: pos=12, state=Empty**

1. **调用 findReusableComputation(12, edit, Empty)**
   
   ```
   当前位置 pos=12 >= edit.pos + edit.numInserted (7+1=8)
   → 进入情况2: 解析位置在插入内容之后
   → 映射到旧文档: oldPos = 12 - 1 = 11
   → 调用 LargestMatchingSubComputation(11, MaxUint64, Empty)
   ```

2. **在计算树中查找:**
   
   ```
   查找从旧文档位置11开始的计算
   
   旧文档位置11是空格 " "
   旧计算树中没有从位置11开始的节点
   
   结论: 找不到可重用的计算
   ```

3. **重新解析:**
   
   ```
   从位置12开始解析新文本:  "baz"
   
   解析器看到:
   - 位置12: ' ' → 不是引号
   - 继续读到位置13: '"' → 停止
   
   消费: 1个字符（空格）
   无标记
   ```

4. **更新状态:**
   ```
   pos = 12 + 1 = 13
   state = Empty
   ```

---

**迭代 5: pos=13, state=Empty**

1. **调用 findReusableComputation(13, edit, Empty)**
   
   ```
   当前位置 pos=13 >= edit.pos + edit.numInserted (7+1=8)
   → 进入情况2: 解析位置在插入内容之后
   → 映射到旧文档: oldPos = 13 - 1 = 12
   → 调用 LargestMatchingSubComputation(12, MaxUint64, Empty)
   ```

2. **在计算树中查找:**
   
   ```
   查找从旧文档位置12开始的计算
   
   [12-17] 节点:
   - 起始位置: 12 ✓
   - 读取范围: [12, 17) → 没有限制 ✓
   - 起始状态: Empty ✓
   
   找到了！返回 [12-17]
   ```

3. **重用计算 [12-17]:**
   
   ```
   ✅ 直接重用，不需要重新解析！
   
   但是要注意位置偏移：
   - 旧文档中这个标记在位置 12-17
   - 新文档中应该在位置 13-18（因为插入了1个字符）
   
   重用的标记: Token{13, 18, String}  // "baz"
   消费长度: 5
   ```

4. **更新状态:**
   ```
   pos = 13 + 5 = 18
   state = Empty
   ```

---

**循环结束:** pos=18 >= 文档长度18

#### 第四步：最终结果

**新的计算树:**
```
                [Root: 0-18]
               /            \
        [0-6]              [6-18]
        重用!              /      \
                      [6-12]    [13-18]
                      重新解析   重用!
```

**标记列表:**
```
[
  Token{0, 5, String},    // "foo"  ← 重用
  Token{6, 12, String},   // "bxar" ← 重新解析
  Token{13, 18, String}   // "baz"  ← 重用（位置偏移+1）
]
```

**性能分析:**
```
总共3个标记:
- 重用: 2个 (66.7%)
- 重新解析: 1个 (33.3%)

如果是完全重新解析，需要解析整个文档18个字符
增量解析只需要重新解析中间部分约6个字符
性能提升: ~67%
```

**可视化对比:**
```
旧文本: "foo" "bar" "baz"
        ^^^^^ ^^^^^ ^^^^^
        重用  重解析 重用

新文本: "foo" "bxar" "baz"
        ^^^^^ ^^^^^^ ^^^^^
        重用  重解析  重用
```### 示例 3: 删除操作（位置映射的关键）

#### 场景说明
删除操作比插入更复杂，因为需要将新文档的位置映射回旧文档。这个例子展示删除如何影响位置计算。

#### 第一步：初始解析

**文本内容:**
```
位置: 0    5    10   15   17
文本: "foo" "bar" "baz"
      ^   ^ ^   ^ ^   ^
```

**计算树:**（与示例2相同）
```
[0-5: "foo"], [6-11: "bar"], [12-17: "baz"]
```

#### 第二步：编辑操作

**在位置 8 删除 1 个字符（删除 'a'）**

```
旧文本: "foo" "bar" "baz"
        01234567890123456
               ↓ 删除这个
新文本: "foo" "br" "baz"
        0123456789012345
```

**编辑信息:**
```go
Edit {
    pos: 8,           // 删除位置
    numInserted: 0,   // 没有插入
    numDeleted: 1     // 删除1个字符
}
```

**关键理解:**
```
旧文档位置 → 新文档位置的映射:
0-7   → 0-7   (删除位置之前，位置不变)
8     → 被删除
9-17  → 8-16  (删除位置之后，位置-1)

新文档位置 → 旧文档位置的映射:
0-7   → 0-7   (删除位置之前，位置不变)
8-16  → 9-18  (删除位置之后，位置+1)
```

#### 第三步：增量重解析

---

**迭代 1: pos=0, state=Empty**

1. **调用 findReusableComputation(0, edit, Empty)**
   
   ```
   当前位置 pos=0 < edit.pos=8
   → 进入情况1: 解析位置在编辑位置之前
   → 调用 LargestMatchingSubComputation(0, 8, Empty)
   ```

2. **在计算树中查找:**
   
   ```
   [0-5] 节点:
   - 起始位置: 0 ✓
   - 读取范围: [0, 5) → 没超过8 ✓
   - 起始状态: Empty ✓
   
   找到了！返回 [0-5]
   ```

3. **重用计算:**
   
   ```
   ✅ 重用 [0-5]
   标记: Token{0, 5, String}  // "foo"
   消费: 5
   ```

4. **更新状态:**
   ```
   pos = 5
   ```

---

**迭代 2: pos=5, state=Empty**

1. **调用 findReusableComputation(5, edit, Empty)**
   
   ```
   当前位置 pos=5 < edit.pos=8
   → 进入情况1
   → 调用 LargestMatchingSubComputation(5, 8, Empty)
   ```

2. **查找结果:**
   
   ```
   找不到从位置5开始的计算
   ```

3. **重新解析:**
   
   ```
   解析新文档位置5:  "br" "baz"
   
   解析器看到:
   - 位置5: ' ' → 空格，消费到下一个引号
   
   消费: 1
   无标记
   ```

4. **更新状态:**
   ```
   pos = 6
   ```

---

**迭代 3: pos=6, state=Empty**

1. **调用 findReusableComputation(6, edit, Empty)**
   
   ```
   当前位置 pos=6 < edit.pos=8
   → 进入情况1
   → 调用 LargestMatchingSubComputation(6, 8, Empty)
   ```

2. **在计算树中查找:**
   
   ```
   [6-11] 节点:
   - 起始位置: 6 ✓
   - 读取范围: [6, 11) → 超过了8 ✗
   
   找不到可重用的计算
   ```

3. **重新解析:**
   
   ```
   解析新文档位置6: "br" "baz"
   
   解析器看到:
   - 位置6: '"' → 开始字符串
   - 位置7-8: 'br'
   - 位置9: '"' → 结束字符串
   
   生成标记: Token{6, 10, String}  // "br"
   消费: 4
   ```

4. **更新状态:**
   ```
   pos = 10
   ```

---

**迭代 4: pos=10, state=Empty**

1. **调用 findReusableComputation(10, edit, Empty)**
   
   ```
   当前位置 pos=10 >= edit.pos=8
   → 进入情况3: 解析位置在删除位置之后
   → 映射到旧文档: oldPos = 10 + 1 = 11
   → 调用 LargestMatchingSubComputation(11, MaxUint64, Empty)
   ```

2. **关键理解 - 为什么要 +1？**
   
   ```
   新文档位置10对应什么内容？
   新文档: "foo" "br" "baz"
           0123456789012345
                     ↑ 位置10是空格
   
   旧文档中这个空格在哪里？
   旧文档: "foo" "bar" "baz"
           01234567890123456
                      ↑ 位置11是空格
   
   因为删除了位置8的字符，所以:
   新文档位置10 = 旧文档位置11
   ```

3. **在计算树中查找:**
   
   ```
   查找从旧文档位置11开始的计算
   
   找不到（位置11是空格，不是计算节点的起始位置）
   ```

4. **重新解析:**
   
   ```
   解析新文档位置10:  "baz"
   
   解析器看到:
   - 位置10: ' ' → 空格
   
   消费: 1
   无标记
   ```

5. **更新状态:**
   ```
   pos = 11
   ```

---

**迭代 5: pos=11, state=Empty**

1. **调用 findReusableComputation(11, edit, Empty)**
   
   ```
   当前位置 pos=11 >= edit.pos=8
   → 进入情况3: 解析位置在删除位置之后
   → 映射到旧文档: oldPos = 11 + 1 = 12
   → 调用 LargestMatchingSubComputation(12, MaxUint64, Empty)
   ```

2. **位置映射分析:**
   
   ```
   新文档位置11: "baz" 的开始引号
   旧文档位置12: "baz" 的开始引号
   
   完美匹配！
   ```

3. **在计算树中查找:**
   
   ```
   [12-17] 节点:
   - 起始位置: 12 ✓
   - 读取范围: [12, 17) ✓
   - 起始状态: Empty ✓
   
   找到了！返回 [12-17]
   ```

4. **重用计算（带位置调整）:**
   
   ```
   ✅ 重用 [12-17]
   
   但是位置需要调整:
   - 旧文档: Token{12, 17, String}
   - 新文档: Token{11, 16, String}  // 位置-1
   
   标记: Token{11, 16, String}  // "baz"
   消费: 5
   ```

5. **更新状态:**
   ```
   pos = 16
   ```

---

**循环结束:** pos=16 >= 文档长度16

#### 第四步：最终结果

**标记列表:**
```
[
  Token{0, 5, String},    // "foo"  ← 重用
  Token{6, 10, String},   // "br"   ← 重新解析
  Token{11, 16, String}   // "baz"  ← 重用（位置偏移-1）
]
```

**位置映射总结:**
```
删除操作的位置映射规则:
- 删除位置之前: 位置不变
- 删除位置之后: 新位置 = 旧位置 - 删除数量
- 反向映射: 旧位置 = 新位置 + 删除数量

这就是为什么代码中有:
pos + edit.numDeleted
```

**可视化对比:**
```
旧文本: "foo" "bar" "baz"
        ^^^^^ ^^^^^ ^^^^^
        重用  重解析 重用

新文本: "foo" "br" "baz"
        ^^^^^ ^^^^ ^^^^^
        重用  重解析 重用
                    ↑
                位置向左偏移1
```

---

## 数据结构详解

### AVL 树平衡

计算树使用 **AVL 树** 保持平衡，确保查询和更新的时间复杂度为 O(log n)。

#### 旋转操作

**左旋 (rotateLeft):**
```
    [x]                [y']
   /   \              /   \
  [q]  [y]    ==>   [x']   [s]
      /   \        /   \
    [r]   [s]     [q]  [r]
```

**右旋 (rotateRight):**
```
      [x]                [y']
     /   \              /   \
    [y]  [s]    ==>   [q]   [x']
   /   \                    /   \
 [q]   [r]                [r]   [s]
```

### 查询操作

#### TokenAtPosition - 查找位置的标记

```go
func (c *computation) TokenAtPosition(pos uint64) Token {
    var offset uint64
    for c != nil {
        // 检查叶子节点的标记
        for _, tok := range c.tokens {
            if pos >= offset+tok.Offset && pos < offset+tok.Offset+tok.Length {
                return Token{...}
            }
        }
        
        // 递归搜索子树
        if c.leftChild != nil && pos < offset+c.leftChild.consumedLength {
            c = c.leftChild  // 向左
        } else {
            offset += c.leftChild.consumedLength
            c = c.rightChild  // 向右
        }
    }
    return Token{}  // 未找到
}
```

**时间复杂度:** O(log n + k)，其中 k 是叶子节点的标记数

---

## 性能优化

### 1. 叶子节点合并

初始解析时，小的叶子节点会被合并：

```go
const minInitialConsumedLen = 1024

if prevComputation.ConsumedLength() < minInitialConsumedLen {
    combineLeaves(prevComputation, nextComputation)
}
```

**好处:** 减少树节点数量，节省内存

### 2. 批量构建树

```go
func concatLeafComputations(computations []*computation) *computation {
    // 逐层构建，避免重复平衡
    for len(computations) > 1 {
        for i := 0; i < len(computations); i += 2 {
            if i+1 < len(computations) {
                nextComputations = append(nextComputations, 
                    computations[i].Append(computations[i+1]))
            }
        }
        computations = nextComputations
    }
    return computations[0]
}
```

**好处:** 比逐个 Append 更高效，树天然平衡

### 3. 失败恢复

```go
func (f Func) recoverFromFailure() Func {
    return func(iter TrackingRuneIter, state State) Result {
        for {
            result := f(iter, state)
            if result.IsSuccess() {
                return result.ShiftForward(numSkipped)
            }
            // 跳过一个字符，继续尝试
            iter.Skip(1)
            numSkipped++
        }
    }
}
```

**好处:** 解析函数失败时自动恢复，确保总能解析完整个文档

---

## 解析函数示例

### 简单字符串解析器

```go
func simpleParseFunc(iter TrackingRuneIter, state State) Result {
    r, err := iter.NextRune()
    if err != nil {
        return FailedResult
    }
    n := uint64(1)
    
    if r == '"' {
        // 查找匹配的引号
        for {
            r, err = iter.NextRune()
            if err != nil {
                return Result{NumConsumed: n, NextState: state}
            } else if r == '"' {
                // 找到匹配引号，生成标记
                token := ComputedToken{
                    Length: n + 1,
                    Role:   TokenRoleString,
                }
                return Result{
                    NumConsumed:    token.Length,
                    ComputedTokens: []ComputedToken{token},
                    NextState:      state,
                }
            }
            n++
        }
    } else {
        // 消费到下一个引号或EOF
        for {
            r, err = iter.NextRune()
            if err != nil || r == '"' {
                return Result{NumConsumed: n, NextState: state}
            }
            n++
        }
    }
}
```

### 组合器模式

系统提供了多个组合器来构建复杂的解析器：

```go
// Then: f 成功后执行 nextFn
func (f Func) Then(nextFn Func) Func

// Or: f 失败则尝试 nextFn
func (f Func) Or(nextFn Func) Func

// ThenMaybe: f 成功，可选地执行 nextFn
func (f Func) ThenMaybe(nextFn Func) Func

// Map: 转换解析结果
func (f Func) Map(mapFn MapFn) Func
```

**示例:** 解析 Go 语言注释
```go
lineComment := consumeString("//").
    Then(consumeToEndOfLine()).
    Map(func(r Result) Result {
        r.ComputedTokens = []ComputedToken{{
            Length: r.NumConsumed,
            Role:   TokenRoleComment,
        }}
        return r
    })

blockComment := consumeString("/*").
    Then(consumeUntil("*/")).
    Map(markAsComment)

comment := lineComment.Or(blockComment)
```

---

## 总结

Aretext 的语法高亮系统通过以下设计实现高效的增量解析：

1. **AVL 树存储解析结果**：支持 O(log n) 的查询和更新
2. **智能重用机制**：根据编辑位置和读取范围判断是否可重用
3. **位置映射**：将新文档位置映射到旧文档，查找可重用节点
4. **组合器模式**：简化复杂语言的解析器编写
5. **失败恢复**：确保解析器总能处理完整个文档

这种设计使得即使在大文件中进行编辑，语法高亮更新也能保持流畅。

---

**关键优势:**
- ⚡ **高性能**: 只重新解析受影响的部分
- 💾 **低内存**: 重用未改变的解析结果
- 🔧 **可扩展**: 易于添加新语言支持
- 🛡️ **健壮性**: 自动从解析失败中恢复

---

## 附录：测试用例分析

### TestParseAll - 完整解析测试

测试初始解析功能，验证能正确识别各种情况：
- 空文档
- 单个标记
- 多个标记
- 无标记的文本
- 标记在中间位置

### TestReparseAfterEditInsertion - 插入测试

测试插入字符后的增量解析：
- 空文档插入
- 插入导致标记分裂
- 插入影响多个标记
- 插入只影响部分标记

### TestReparseAfterEditDeletion - 删除测试

测试删除字符后的增量解析：
- 删除无标记文本
- 删除影响标记
- 删除改变标记长度
- 删除影响多个标记

### TestReparseIndividualInsertionsAtEndOfDocument - 逐字符插入测试

测试极端情况：逐个字符插入，验证增量解析的正确性和效率。

这些测试确保了增量解析算法在各种编辑场景下都能正确工作。


---

## 增量解析核心概念总结

### 1. 为什么需要增量解析？

**传统方法的问题:**
```
每次编辑都重新解析整个文档
→ 大文件编辑时性能差
→ 用户体验不流畅
```

**增量解析的优势:**
```
只重新解析受影响的部分
→ 编辑响应快速
→ 内存使用高效
```

### 2. 重用判断的两个关键条件

#### 条件1: 读取范围未被编辑影响

```
计算节点记录了它"读取"了哪些字符

例如: 计算节点 [6-11] 读取了位置 6-10 的字符

如果编辑发生在位置 7:
→ 这个计算节点读取的内容已经改变
→ 不能重用！

如果编辑发生在位置 12:
→ 这个计算节点读取的内容没有改变
→ 可以重用！
```

#### 条件2: 起始状态匹配

```
解析器可能有状态（例如：在字符串内部、在注释内部等）

只有当计算节点的起始状态与当前解析状态相同时，才能重用

例如:
- 旧计算: 从 Empty 状态开始
- 当前: 在 InString 状态
→ 不能重用！因为状态不同，解析结果会不同
```

### 3. 位置映射的三种情况

#### 情况1: 解析位置在编辑位置之前

```
pos < edit.pos

例如: 编辑位置7，当前解析位置3
→ 查找从位置3开始、不超过位置7的计算
→ 确保不会读取到被编辑的内容
```

#### 情况2: 插入后的位置映射

```
编辑: 在位置7插入1个字符
当前: 解析位置10

新文档位置10 对应 旧文档位置9
→ oldPos = 10 - 1 = 9

为什么？
新文档: "foo" "bxar" "baz"
        0123456789012345678
                  ↑ 位置10
                  
旧文档: "foo" "bar" "baz"
        01234567890123456
                 ↑ 位置9
                 
因为插入了1个字符，所以新文档的位置10
对应旧文档的位置9
```

#### 情况3: 删除后的位置映射

```
编辑: 在位置7删除1个字符
当前: 解析位置10

新文档位置10 对应 旧文档位置11
→ oldPos = 10 + 1 = 11

为什么？
新文档: "foo" "br" "baz"
        0123456789012345
                  ↑ 位置10
                  
旧文档: "foo" "bar" "baz"
        01234567890123456
                   ↑ 位置11
                   
因为删除了1个字符，所以新文档的位置10
对应旧文档的位置11
```

### 4. AVL树的作用

#### 为什么使用树结构？

```
如果用数组存储所有计算结果:
- 查找: O(n) - 需要遍历
- 插入: O(n) - 需要移动元素

使用AVL树:
- 查找: O(log n) - 二分查找
- 插入: O(log n) - 自动平衡
- 内存: 可以共享未改变的子树
```

#### 树的组织方式

```
        [Root: 0-100]
       /              \
   [0-50]            [50-100]
   /    \            /      \
[0-25] [25-50]  [50-75]  [75-100]

每个节点记录:
- 起始位置
- 消费长度
- 读取长度
- 起始/结束状态
- 子节点或标记
```

### 5. 实际性能提升

#### 示例场景: 10000行代码文件

```
场景1: 在第5000行插入一个字符

传统方法:
- 重新解析10000行
- 时间: ~100ms

增量解析:
- 重用前4999行的解析结果
- 重新解析第5000行
- 重用后5000行的解析结果（位置偏移）
- 时间: ~1ms

性能提升: 100倍！
```

```
场景2: 在第1行插入一个字符

传统方法:
- 重新解析10000行
- 时间: ~100ms

增量解析:
- 重新解析第1行
- 重用后9999行的解析结果（位置偏移）
- 时间: ~1ms

性能提升: 100倍！
```

### 6. 关键代码流程图

```
用户编辑文档
    ↓
创建 Edit 对象（记录编辑位置和类型）
    ↓
调用 ReparseAfterEdit
    ↓
循环: 从文档开始到结束
    ↓
    ├─→ 调用 findReusableComputation
    │       ↓
    │   判断当前位置与编辑位置的关系
    │       ↓
    │   ├─→ 位置在编辑前: 查找不超过编辑位置的计算
    │   ├─→ 位置在插入后: 映射到旧文档（减去插入数）
    │   └─→ 位置在删除后: 映射到旧文档（加上删除数）
    │       ↓
    │   在旧计算树中查找匹配的子计算
    │       ↓
    │   ├─→ 找到: 返回可重用的计算
    │   └─→ 没找到: 返回 nil
    │
    ├─→ 如果找到可重用计算
    │   └─→ 直接使用（可能需要位置偏移）
    │
    └─→ 如果没找到可重用计算
        └─→ 调用 parseFunc 重新解析
    ↓
将新计算追加到结果树
    ↓
继续下一个位置
    ↓
循环结束，更新 lastComputation
```

### 7. 常见问题解答

#### Q1: 为什么不能总是重用所有计算？

```
A: 因为编辑可能改变了解析结果

例如: "foo" → "foo
- 原来是完整字符串
- 现在是未闭合字符串
- 后续内容的解析会完全不同
```

#### Q2: 读取长度和消费长度有什么区别？

```
A: 
- 读取长度: 解析器"看了"多少字符
- 消费长度: 解析器"处理了"多少字符

例如: 解析 "foo"
- 读取长度: 5 (读取了 "foo" 这5个字符)
- 消费长度: 5 (处理了这5个字符)

例如: 解析失败后跳过
- 读取长度: 10 (尝试读取了10个字符)
- 消费长度: 1 (只跳过了1个字符)

读取长度用于判断是否受编辑影响
消费长度用于推进解析位置
```

#### Q3: 为什么需要状态（State）？

```
A: 某些语言的解析依赖上下文

例如: Bash 脚本
- 在普通模式: $ 是变量开始
- 在单引号内: $ 是普通字符

状态记录当前的解析上下文，确保重用的计算
在相同上下文下才有效
```

#### Q4: 如果连续快速编辑会怎样？

```
A: 每次编辑都会调用 ReparseAfterEdit

例如: 快速输入 "hello"
1. 插入 'h' → 增量解析
2. 插入 'e' → 增量解析（重用上次结果）
3. 插入 'l' → 增量解析（重用上次结果）
4. 插入 'l' → 增量解析（重用上次结果）
5. 插入 'o' → 增量解析（重用上次结果）

每次都能重用大部分之前的解析结果
```

---

## 总结

Aretext 的增量语法解析系统是一个精妙的设计，它通过：

1. **AVL树存储解析结果** - 实现 O(log n) 的查询和更新
2. **智能重用机制** - 根据编辑位置和读取范围判断是否可重用
3. **位置映射算法** - 将新文档位置映射到旧文档，查找可重用节点
4. **状态跟踪** - 确保重用的计算在相同上下文下有效
5. **失败恢复** - 保证解析器总能处理完整个文档

这些技术结合在一起，使得即使在大文件中进行编辑，语法高亮更新也能保持流畅，为用户提供了出色的编辑体验。

**关键优势:**
- ⚡ **高性能**: 只重新解析受影响的部分，通常能节省90%以上的计算
- 💾 **低内存**: 重用未改变的解析结果，避免重复存储
- 🔧 **可扩展**: 易于添加新语言支持，只需实现 parseFunc
- 🛡️ **健壮性**: 自动从解析失败中恢复，不会因为语法错误而崩溃
- 🎯 **精确性**: 通过状态跟踪确保解析结果的正确性

这种设计模式不仅适用于语法高亮，也可以应用于其他需要增量更新的场景，如：
- 代码补全
- 错误检查
- 代码格式化
- 语法树构建
