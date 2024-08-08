package main

import (
	"fmt"
	"github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/avl_tree"
	"math/rand"
	"sync"
)

func main() {
	bs := avl_tree.New(true)
	wg := sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		bs.Insert(rand.Intn(100))
		go func() {
			bs.Insert(rand.Intn(100))
			wg.Done()
		}()
	}
	wg.Wait()
	fmt.Println("遍历输出所有插入的元素")
	for i := bs.Iterator(); i.HasNext(); i.Next() {
		fmt.Printf("%+v ", i.Value())
	}
	fmt.Println()
	fmt.Println("删除一次平衡二叉树中存在的元素,存在重复的将会被剩下")
	for i := 0; i < 10; i++ {
		bs.Erase(i)
	}
	fmt.Println("输出剩余的重复元素")
	for i := bs.Iterator(); i.HasNext(); i.Next() {
		fmt.Printf("%+v ", i.Value())
	}
	fmt.Println()
}
