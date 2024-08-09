package main

import (
	"fmt"
	"github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/hash_map"
)

func main() {
	m := hash_map.New()
	for i := 1; i <= 17; i++ {
		m.Insert(i, i)
	}
	fmt.Println("size=", m.Size())
	keys := m.GetKeys()
	fmt.Println("keys:", keys)
	for i := 0; i < len(keys); i++ {
		//value := m.Get(keys[i])
		fmt.Printf("key:%+v value:%+v\n", keys[i], m.Get(keys[i]))
	}
	fmt.Println()

	for i := m.Iterator().Begin(); i.HasNext(); i.Next() {
		fmt.Printf("%+v ", i.Value())
	}
	fmt.Println()

	for i := 0; i < len(keys); i++ {
		m.Erase(keys[i])
	}
}
