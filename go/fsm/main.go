package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// 命令定义结构
type Command struct {
	Name        string                                    // 命令描述
	BuildExpr   func() Expr                               // 解析表达式构造器
	BuildAction func(ctx Context, p CommandParams) Action // 动作构造器
}

// 上下文结构
type Context struct {
	State string // 当前状态
}

// 命令参数
type CommandParams struct {
	Count int // 数字计数
}

// 动作接口
type Action interface {
	Execute()
}

// 表达式接口
type Expr interface {
	Match(input string) bool
}

// 具体动作实现
type CursorLeftAction struct {
	Count int
}

func (a CursorLeftAction) Execute() {
	fmt.Printf("执行命令: 左移 %d 次\n", a.Count)
}

type DeleteInnerWordAction struct{}

func (a DeleteInnerWordAction) Execute() {
	fmt.Println("执行命令: 删除当前单词")
}

type VerbExpr struct {
}

func (e VerbExpr) Match(input string) bool {
	fmt.Println("VerbExpr true")
	return true
}

// 具体表达式实现
type AltExpr struct {
	Options []Expr
}

func (e AltExpr) Match(input string) bool {
	for _, option := range e.Options {
		if option.Match(input) {
			return true
		}
	}
	return false
}

type RuneExpr struct {
	Rune rune
}

func (e RuneExpr) Match(input string) bool {
	return len(input) == 1 && rune(input[0]) == e.Rune
}

type VerbCountThenExpr struct {
	Verb   Expr
	Follow Expr
}

func (e VerbCountThenExpr) Match(input string) bool {
	// 这里可以实现更复杂的匹配逻辑
	return e.Verb.Match(input[:1]) && e.Follow.Match(input[1:])
}

// 辅助函数：构建表达式
func verbCountThenExpr(follow Expr) Expr {
	verb := VerbExpr{}
	return VerbCountThenExpr{Verb: verb, Follow: follow}
}

func altExpr(options ...Expr) Expr {
	return AltExpr{Options: options}
}

func runeExpr(r rune) Expr {
	return RuneExpr{Rune: r}
}

// 全局变量
var commands = []Command{
	{
		Name: "cursor left (left arrow or h)",
		BuildExpr: func() Expr {
			return verbCountThenExpr(altExpr(runeExpr('h'), runeExpr('h')))
		},
		BuildAction: func(ctx Context, p CommandParams) Action {
			return CursorLeftAction{Count: p.Count}
		},
	},
	{
		Name: "delete inner word (d+iw)",
		BuildExpr: func() Expr {
			return verbCountThenExpr(altExpr(runeExpr('d'), runeExpr('i')))
		},
		BuildAction: func(ctx Context, p CommandParams) Action {
			return DeleteInnerWordAction{}
		},
	},
}

// 处理输入
func processInput(input string) {
	for _, cmd := range commands {
		expr := cmd.BuildExpr()
		if expr.Match(input) {
			params := CommandParams{
				Count: parseCount(input),
			}
			action := cmd.BuildAction(Context{}, params)
			action.Execute()
			return
		}
	}
	fmt.Println("非法命令:", input)
}

// 解析数字
func parseCount(input string) int {
	countStr := ""
	for _, char := range input {
		if char >= '0' && char <= '9' {
			countStr += string(char)
		} else {
			break
		}
	}
	if countStr == "" {
		return 1
	}
	count, _ := strconv.Atoi(countStr)
	return count
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("进入编辑模式，输入 'q' 退出")
	for {
		fmt.Print("> ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "q" {
			fmt.Println("退出程序")
			break
		}

		processInput(input)
	}
}
