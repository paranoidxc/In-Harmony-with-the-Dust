package main

import (
	"fmt"

	"github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/deque"
)

func main() {
	d := deque.New()
	for i := 1; i <= 15; i++ {
		d.PushFront(i)
	}
	show(d)
	fmt.Println("当前首元素:", d.Front())
	fmt.Println("当前尾元素:", d.Back())

	val := d.PopFront()
	fmt.Println("pop front当前元素:", val)
	show(d)

	fmt.Println("pushbakc 20")
	d.PushBack(20)
	show(d)

	fmt.Println("当前首元素:", d.Front())
	fmt.Println("当前尾元素:", d.Back())

	val = d.PopFront()
	fmt.Println("front当前元素:", val)

	fmt.Println("当前首元素:", d.Front())
	fmt.Println("当前尾元素:", d.Back())

	val = d.PopBack()
	fmt.Println("back当前元素:", val)

	fmt.Println("当前首元素:", d.Front())
	fmt.Println("当前尾元素:", d.Back())

	val = d.PopFront()
	fmt.Println("front当前元素:", val)
}

func show(d *deque.Deque) {
	fmt.Println("Iterator:")
	for i := d.Iterator(); i.HasNext(); i.Next() {
		fmt.Printf("%+v ", i.Value())
	}
	fmt.Printf("\n")
}
