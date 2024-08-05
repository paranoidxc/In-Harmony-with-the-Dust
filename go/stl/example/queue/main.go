package main

import (
	"fmt"
	"sync"

	"github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/queue"
)

func main() {
	q := queue.New()
	wg := sync.WaitGroup{}
	//随机插入队列中
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(num int) {
			fmt.Printf("%+v ", num)
			q.Push(num)
			wg.Done()
		}(i)
	}
	wg.Wait()
	fmt.Println()
	fmt.Println("输出首部:", q.Front())
	fmt.Println("输出尾部:", q.Back())
	fmt.Println("弹出并输出前4个:")
	for i := uint64(0); i < q.Size()-1; i++ {
		fmt.Printf("%+v ", q.Pop())
	}
	fmt.Println()
	//在尾部再添加4个,从10开始以做区分
	for i := 10; i < 14; i++ {
		q.Push(i)
	}
	fmt.Println("从头输出全部:")
	for !q.Empty() {
		fmt.Printf("%+v ", q.Pop())
	}
}
