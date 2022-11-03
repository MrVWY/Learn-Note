package main

import "fmt"

// | 表示取并集
// ~int 表示底层类型必须为int

type NumStr interface {
	Num | Str
}
type Num interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr | ~float32 | ~float64 | ~complex64 | ~complex128
}
type Str interface {
	string
}

func add[T NumStr](a, b T) T {
	return a + b
}

func main() {
	fmt.Println(add(3, 4))
	fmt.Println(add("hello", "world"))
}
