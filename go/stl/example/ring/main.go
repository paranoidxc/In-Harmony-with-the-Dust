package main

import (
	"fmt"

	"github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/ring"
	"github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/utils/iterator"
)

func iter(iter *iterator.Iterator) {
	for i := iter; i.HasNext(); i.Next() {
		fmt.Printf("%+v ", i.Value())
	}
}

func show(r *ring.Ring) {
	for i := uint64(0); i < r.Size(); i++ {
		fmt.Printf("%+v ", r.Value())
		r.Next()
	}
}

func main() {
	r := ring.New()
	for i := 1; i < 10; i++ {
		r.Insert(i)
	}
	fmt.Println("ring当前持有结点的元素:", r.Value())
	fmt.Println("====================")
	show(r)
	fmt.Println("\n-----------------")
	iter(r.Iterator())
	fmt.Println("====================")
	fmt.Println("ring当前持有结点的元素:", r.Value())
	fmt.Println("====================")
	r.Erase()
	show(r)
	fmt.Println("\n-----------------")
	iter(r.Iterator())
	fmt.Println("\n====================")

	fmt.Println("ring当前持有结点的元素:", r.Value())
	r.Erase()
	show(r)
	fmt.Println("\n-----------------")
	iter(r.Iterator())
	fmt.Println("\n====================")

	fmt.Println("ring当前持有结点的元素:", r.Value())
	r.Erase()
	show(r)
	fmt.Println("\n-----------------")
	iter(r.Iterator())
	fmt.Println("\n====================")

	r.Insert(111)
	show(r)
	fmt.Println("\n-----------------")
	iter(r.Iterator())
	fmt.Println("\n====================")

	fmt.Println("ring当前持有结点的元素:", r.Value())
	r.Erase()

	show(r)
	fmt.Println("\n-----------------")
	iter(r.Iterator())
	fmt.Println("\n====================")
}
