package main

import (
	"fmt"
	"strconv"
)

type Prices int

type ShowPrices interface {
	Strings() string
	~int | string
}

func (i Prices) Strings() string {
	return strconv.Itoa(int(i))
}

func ShowPriceLists[T ShowPrices](s []T) (ret []string) {
	for _, v := range s {
		ret = append(ret, v.Strings())
	}
	return
}

func main() {
	fmt.Println(ShowPriceLists([]Prices{1, 2}))
}

//传入浮点参数，就会因为不是约束类型而报错
