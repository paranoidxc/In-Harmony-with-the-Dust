package main

import (
	"fmt"
	"github.com/paranoidxc/In-Harmony-with-the-Dust/go/stl/utils/comparator"
)

func main() {
	var arr = make([]interface{}, 0, 0)
	arr = append(arr, 5)
	arr = append(arr, 3)
	arr = append(arr, 2)
	arr = append(arr, 4)
	arr = append(arr, 1)
	arr = append(arr, 4)
	arr = append(arr, 3)
	arr = append(arr, 1)
	arr = append(arr, 10)
	arr = append(arr, 5)
	arr = append(arr, 2)
	//arr = append(arr, 10)
	comparator.Sort(&arr)
	for i := 0; i < len(arr); i++ {
		fmt.Printf("%+v ", arr[i].(int))
	}
	fmt.Println()

	fmt.Println("search idx:", comparator.Search(&arr, 3))
	fmt.Println("search idx:", comparator.Search(&arr, 10))

	fmt.Println("NthElement:")
	for i := 0; i < len(arr); i++ {
		fmt.Printf("i=%d, val:%+v \n", i, comparator.NthElement(&arr, i))
	}
	fmt.Println()
	for i := 0; i < len(arr); i++ {
		fmt.Printf("%+v ", arr[i].(int))
	}
	fmt.Println()
	for i := 0; i < len(arr); i++ {
		fmt.Printf("find=%d, upper:%+v; lower:%+v \n",
			i,
			comparator.UpperBound(&arr, i),
			comparator.LowerBound(&arr, i),
		)
	}
}
