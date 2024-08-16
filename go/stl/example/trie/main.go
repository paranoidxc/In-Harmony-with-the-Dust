package main

import (
	"fmt"
	"github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/trie"
)

func main() {
	t := trie.New()
	t.Insert("hlccd", "hlccd")
	t.Insert("ha", "ha")
	t.Insert("hb", "hb")
	t.Insert("hc", "hc")
	t.Insert("hd", "hd")
	t.Insert("he", "he")
	t.Insert("hl", "hl")
	t.Insert("hlccd1", "hlccd1")
	t.Insert("hlccd2", "hlccd2")
	t.Insert("hlccd3", "hlccd3")
	t.Insert("hlccd+", "hlccd")
	t.Insert("hlccd/", "hlccd")
	fmt.Println("当前插入的所有string:")

	for i := t.Iterator().Begin(); i.HasNext(); i.Next() {
		fmt.Printf("%+v ", i.Value())
	}
	fmt.Println()

	t.Erase("h")
	t.Erase("ha")
	t.Erase("hb")
	t.Erase("hc")
	t.Erase("hd")
	t.Erase("he")
	fmt.Println("定向删除后剩余的string:")
	for i := t.Iterator().Begin(); i.HasNext(); i.Next() {
		fmt.Printf("%+v ", i.Value())
	}
	fmt.Println()
	t.Delete("h")
	fmt.Println("删除以'h'为前缀的所有元素后剩余的数量:")
	for i := t.Iterator().Begin(); i.HasNext(); i.Next() {
		fmt.Println(i.Value())
	}
}
