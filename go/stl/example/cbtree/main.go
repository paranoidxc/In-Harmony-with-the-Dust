package main

import (
	"fmt"
	"github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/cbtree"
	"sync"
)

func main() {
	h := cbtree.New()
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(num int) {
			h.Push(num)
		}(i)
	}
	fmt.Println("利用迭代器输出堆中存储的所有元素:")
	for i := h.Iterator(); i.HasNext(); i.Next() {
		fmt.Printf("%+v ", i.Value())
	}
	fmt.Println()
	fmt.Println("依次输出顶部元素:")
	for !h.Empty() {
		fmt.Printf("%+v ", h.Top())
		h.Pop()
	}
	fmt.Println()
}
