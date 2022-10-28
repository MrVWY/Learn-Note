package main

import "fmt"

type vector[T any] []T

func printslice2[T any](s []T) {
	for _, v := range s {
		fmt.Println("%v", v)
	}
}

func main() {
	v := vector[int]{58, 1881}
	printslice2(v)
	v2 := vector[string]{"烤鸡", "烤鸭", "烤鱼", "烤面筋"}
	printslice2(v2)
}
