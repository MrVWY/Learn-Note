package main

import "fmt"

type M[K string, V any] map[K]V //这里的K不支持any ,由于底层map不支持,所以使用string

func main() {
	m1 := M[string, int]{"key": 1}
	m1["key"] = 2

	m2 := M[string, string]{"key": "value"}
	m2["key"] = "new value"
	fmt.Println(m1, m2)

}
