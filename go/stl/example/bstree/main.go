package main

import (
	"fmt"
	"github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/bstree"
	"sync"
)

func main() {
	bs := bstree.New(true)
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		bs.Insert(i)
		go func() {
			bs.Insert(i)
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("遍历输出所有插入的元素")
	for i := bs.Iterator(); i.HasNext(); i.Next() {
		fmt.Printf("%+v ", i.Value())
	}
	fmt.Println()
	fmt.Println("删除一次二叉搜索树中存在的元素,存在重复的将会被剩下")
	for i := 0; i < 10; i++ {
		bs.Erase(i)
	}
	fmt.Println("输出剩余的重复元素")
	for i := bs.Iterator(); i.HasNext(); i.Next() {
		fmt.Printf("%+v ", i.Value())
	}
	fmt.Println()
}
