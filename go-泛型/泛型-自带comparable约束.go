package main

import (
	"fmt"
)

//comparable: 表示go里面所有内置的可比较类型：int、uint、float、bool、struct、指针等一切可以比较的类型

func findFunc[T comparable](a []T, v T) int {
	for i, e := range a {
		if e == v {
			return i
		}
	}
	return -1
}

func main() {
	fmt.Println(findFunc([]int{1, 2, 3, 4, 5, 6}, 5))
	fmt.Println(findFunc([]string{"烤鸡", "烤鸭", "烤鱼", "烤面筋"}, "烤面筋"))
}
