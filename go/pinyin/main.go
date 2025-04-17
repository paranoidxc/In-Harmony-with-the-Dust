package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type pyword struct {
	py   string
	word string
}

var MulPy2Words map[string][]string
var SiglePy2Words map[string][]string
var FuzzyPy2Words map[string][]pyword

func readFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		//fmt.Println(line)
		words := strings.Split(line, " ")
		if len(words) == 2 {
			tmp := SiglePy2Words
			if strings.Index(words[1], "'") >= 0 {
				tmp = MulPy2Words
				first := words[1][0:1]
				FuzzyPy2Words[first] = append(FuzzyPy2Words[first], pyword{py: words[1], word: words[0]})
			}
			tmp[words[1]] = append(tmp[words[1]], words[0])
		}
	}
	return nil
}

func readFileWord(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		//fmt.Println(line)
		words := strings.Split(line, " ")
		if len(words) == 3 {
			tmp := SiglePy2Words
			if strings.Index(words[1], "'") >= 0 {
				tmp = MulPy2Words
			}

			tmp[words[1]] = append(tmp[words[1]], words[0])
		}
	}
	return nil
}

func py2pys(input string) (r []string, err error) {
	// 定义正则表达式
	// [^aoeiuv]? 表示可选的非元音字母（即辅音字母）
	// h? 表示可选的 'h'
	// [iuv]? 表示可选的 i/u/v
	// (ai|ei|ao|ou|er|ang?|eng?|ong|a|o|e|i|u|ng|n)? 表示可选的拼音组合
	pattern := `[^aoeiuv]?h?[iuv]?(ai|ei|ao|ou|er|ang?|eng?|ong|a|o|e|i|u|ng|n)?`

	// 编译正则表达式
	re, err := regexp.Compile(pattern)
	if err != nil {
		//fmt.Println("正则表达式编译失败:", err)
		//error.New(""正则表达式编译失败:)
		return
	}

	// 输入字符串
	//input = "woshiyimingshuaibichengxuyuan"

	// 匹配所有符合条件的子串
	r = re.FindAllString(input, -1)

	return r, nil
}

func step1(list []string) string {
	// 第二部分 整句查询
	r := ""
	old := list
	for len(list) >= 2 {
		lenList := len(list)
		listStr := strings.Join(list, "'")
		word, find := MulPy2Words[listStr]
		if find {
			//fmt.Println("step1 >> find:", listStr, word)
			list = old
			list = list[lenList:]
			old = list
			r += word[0]
		} else {
			list = list[0:(lenList - 1)]
		}
	}

	return r
}

func step2(list []string) (r []string) {
	for len(list) >= 2 {
		listStr := strings.Join(list, "'")
		word, find := MulPy2Words[listStr]
		if find {
			r = append(r, word...)
			//fmt.Println("step2 >> find:", listStr, word)
		}
		list = list[0:(len(list) - 1)]
	}

	return r
}

func step3(list []string) string {
	// 模糊词查询   可能不需要
	r := ""
	//first := []string{}
	if len(list) > 0 {
		//f := list[0][0:1]
	}

	// for len(list) >= 2 {
	// 	listStr := strings.Join(list, "'")
	// 	fmt.Println("listStr", listStr)
	// 	// find then continue loop
	// 	list = list[0:(len(list) - 1)]
	// }

	return r
}

func step4(pinyin string) (r []string) {
	// 填充单字 用第一个拼音元素
	if words, find := SiglePy2Words[pinyin]; find {
		r = words
	}

	return r
}

func find(input string) {
	matches, _ := py2pys(input)
	if len(matches) > 0 {
		// 输出结果
		fmt.Println(matches)

		r1 := step1(matches)
		fmt.Printf("r1: %+v\n", r1)

		r2 := step2(matches)
		fmt.Printf("r2: %+v\n", r2)

		r3 := step3(matches)
		fmt.Printf("r3: %+v\n", r3)

		r4 := step4(matches[0])
		fmt.Printf("r4: %s:%+v\n", matches[0], r4)
	}
}

func init() {
	MulPy2Words = make(map[string][]string)
	SiglePy2Words = make(map[string][]string)

	readFile("./daily.txt")
	//fmt.Println("test:", MulPy2Words["wo'shi"])

	readFileWord("./word.txt")
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

		find(input)
	}

	// 输入字符串
	//input := "woshiyimingshuaibichengxuyuan"
	//input := "ceshiwodepinyinhaobuhaoyong"

	// input := "jizhangzonger"
	// matches, _ := py2pys(input)
	// // 输出结果
	// fmt.Println(matches)
	//
	// r1 := step1(matches)
	// fmt.Printf("r1: %+v\n", r1)
	//
	// r2 := step2(matches)
	// fmt.Printf("r2: %+v\n", r2)
	//
	// r4 := step4(matches[0])
	// fmt.Printf("r4: %s:%+v\n", matches[0], r4)

	// reader := bufio.NewReader(os.Stdin)
	// fmt.Println("进入编辑模式，输入 'q' 退出")
	// for {
	// 	fmt.Print("> ")
	// 	input, _ := reader.ReadString('\n')
	// 	input = strings.TrimSpace(input)
	//
	// 	if input == "q" {
	// 		fmt.Println("退出程序")
	// 		break
	// 	}
	//
	// 	matches, _ := py2pys(input)
	// 	// 输出结果
	// 	fmt.Println(matches)
	//
	// 	r1 := step1(matches)
	// 	fmt.Printf("r1: %+v\n", r1)
	//
	// 	r2 := step2(matches)
	// 	fmt.Printf("r2: %+v\n", r2)
	//
	// 	r4 := step4(matches[0])
	// 	fmt.Printf("r4: %s:%+v\n", matches[0], r4)
	// }
}
