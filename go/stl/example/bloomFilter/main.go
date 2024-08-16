package main

import (
	"fmt"
	"github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/bloomFilter"
)

func hash(v interface{}) uint32 {
	return uint32(v.(int))
}
func main() {
	bf := bloomFilter.New(nil)
	for i := 0; i < 10; i++ {
		bf.Insert(i)
	}
	for i := 0; i < 15; i++ {
		fmt.Printf("i= %+v, check=%+v \n", i, bf.Check(i))
	}
	fmt.Println()
}
