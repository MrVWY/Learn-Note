package main

import "fmt"

//any: 表示go里面所有的内置基本类型，等价于interface{}

func printslice[T any](s []T) {
	for _, v := range s {
		fmt.Printf("%v\n", v)
	}
}

func main() {
	printslice[int]([]int{66, 77, 88, 99, 100})
	printslice[float64]([]float64{1.1, 2.2, 5.5})
	printslice[string]([]string{"烤鸡", "烤鸭", "烤鱼", "烤面筋"})

	//省略显示类型
	printslice([]int64{55, 44, 33, 22, 11})

}
