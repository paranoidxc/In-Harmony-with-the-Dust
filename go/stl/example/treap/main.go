package main

import (
	"fmt"
	"github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/treap"
	"math/rand"
)

func main() {
	bs := treap.New(true)
	for i := 0; i < 20; i++ {
		t := rand.Intn(100)
		//fmt.Printf("插入元素%d\n", t)
		bs.Insert(t)
		//show(bs)
	}
	fmt.Println("遍历输出所有插入的元素")
	show(bs)

	fmt.Println("删除一次树堆中存在的元素,存在重复的将会被剩下")
	for i := 0; i < 20; i++ {
		bs.Erase(i)
	}
	fmt.Println("输出剩余的重复元素")
	show(bs)
	fmt.Println()
}

func show(bs *treap.Treap) {
	for i := bs.Iterator(); i.HasNext(); i.Next() {
		fmt.Printf("%+v ", i.Value())
	}
	fmt.Println()
}
