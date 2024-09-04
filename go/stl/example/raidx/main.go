package main

import (
	"fmt"
	"github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/radix"
)

func main() {
	rd := radix.New()
	rd.Insert("/test/:name/:key")
	rd.Insert("/hlccd/:name/:key")
	rd.Insert("/hlccd/1")
	rd.Insert("/hlccd/a/*name")
	rd.Insert("/demo/test")

	fmt.Println("分层匹配")
	m, _ := rd.Mate("/hlccd/test/abc")
	for k, v := range m {
		fmt.Println(k, v)
	}
	fmt.Println("匹配全部")
	m, _ = rd.Mate("/hlccd/a/abc")
	for k, v := range m {
		fmt.Println(k, v)
	}
	fmt.Println("利用迭代器遍历")
	for i := rd.Iterator(); i.HasNext(); i.Next() {
		fmt.Println(i.Value())
	}
	rd.Erase("/hlccd/a/*name")
	fmt.Println("利用迭代器遍历定向删除后的结果")
	for i := rd.Iterator(); i.HasNext(); i.Next() {
		fmt.Println(i.Value())
	}
	rd.Delete("/hlccd/")
	fmt.Println("利用迭代器遍历删除前缀的结果")
	for i := rd.Iterator(); i.HasNext(); i.Next() {
		fmt.Println(i.Value())
	}
}
