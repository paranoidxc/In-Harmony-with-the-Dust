package main

import (
	"fmt"
	"github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/priority_queue"
)

func main() {
	pq := priority_queue.New()
	for i := 0; i < 10; i++ {
		pq.Push(i)
	}
	fmt.Println("遍历所有元素同时弹出：")
	size := pq.Size()
	for i := uint64(0); i < size; i++ {
		fmt.Println(pq.Top())
		pq.Pop()
	}
}
