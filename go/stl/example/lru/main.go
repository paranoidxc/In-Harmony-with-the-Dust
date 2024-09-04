package main

import (
	"fmt"
	lru2 "github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/lru"
)

type Value struct {
	Bytes []byte
}

func (v Value) Len() int {
	return len(v.Bytes)
}
func main() {
	space := 0
	lru := lru2.New(2<<10, nil)

	for i := 0; i < 10000; i++ {
		v := Value{Bytes: []byte(string(i))}
		lru.Insert(string(i), v)
		space += v.Len()
	}
	fmt.Println("应该占用空间:", space)
	fmt.Println("LRU中存放的byte数量:", lru.Size())
	fmt.Println("LRU的byte数量上限:", lru.Cap())
	lru.Erase(string(9999))
	fmt.Println("删除后的LRU中存放的byte数量:", lru.Size())
	fmt.Println("从LRU中找9998")
	if v, ok := lru.Get(string(9998)); ok {
		fmt.Println(string(v.(Value).Bytes))
		fmt.Println(string(9998))
	}
	fmt.Println("从LRU中找1")
	if v, ok := lru.Get(string(1)); ok {
		fmt.Println(string(v.(Value).Bytes))
		fmt.Println(string(1))
	}
}
