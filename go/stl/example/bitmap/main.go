package main

import (
	"fmt"
	"github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/bitmap"
)

func main() {
	var nums []uint
	bm := bitmap.New()
	bm.Insert(1)
	bm.Insert(2)
	bm.Insert(3)
	bm.Insert(64)
	bm.Insert(128)
	bm.Insert(256)
	bm.Insert(320)
	nums = bm.All()
	for i := 0; i < len(nums); i++ {
		fmt.Printf("%+v ", nums[i])
	}
	bm.Delete(320)
	bm.Delete(2)
	fmt.Println()
	nums = bm.All()
	for i := 0; i < len(nums); i++ {
		fmt.Printf("%+v ", nums[i])
	}
	bm.Delete(256)
	fmt.Println()
	nums = bm.All()
	for i := 0; i < len(nums); i++ {
		fmt.Printf("%+v ", nums[i])
	}
	bm.Delete(128)
	fmt.Println()
	nums = bm.All()
	for i := 0; i < len(nums); i++ {
		fmt.Printf("%+v ", nums[i])
	}
	bm.Clear()
	fmt.Println()
	nums = bm.All()
	for i := 0; i < len(nums); i++ {
		fmt.Printf("%+v ", nums[i])
	}
}
