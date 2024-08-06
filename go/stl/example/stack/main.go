package main

import (
	"fmt"
	"github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/stack"
)

func main() {
	s := stack.New()
	for i := 1; i < 10; i++ {
		s.Push(i)
	}
	fmt.Println("使用迭代器遍历全部：")
	for i := s.Iterator(); i.HasNext(); i.Next() {
		fmt.Printf("%+v ", i.Value())
	}
	fmt.Println()

	fmt.Println("使用size边删除边遍历：")
	size := s.Size()
	for i := uint64(0); i < size; i++ {
		fmt.Printf("%+v ", s.Top())
		s.Pop()
	}
	fmt.Println()
}
