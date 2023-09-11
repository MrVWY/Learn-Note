package main

//在本例中， 我们将使用工厂结构体来构建多种类型的武器。

//首先， 我们来创建一个名为 iGun的接口， 其中将定义一支枪所需具备的所有方法。 然后是实现了 iGun 接口的 gun枪支结构体类型。 两种具体的枪支——ak47与 musket火枪 ——
//两者都嵌入了枪支结构体， 且间接实现了所有的 iGun方法。

//gunFactory枪支工厂结构体将发挥工厂的作用， 即通过传入参数构建所需类型的枪支。 main.go 则扮演着客户端的角色。 其不会直接与 ak47或 musket进行互动，
//而是依靠 gunFactory来创建多种枪支的实例， 仅使用字符参数来控制生产。

import "fmt"

func main() {
	ak47, _ := getGun("ak47")
	musket, _ := getGun("musket")

	printDetails(ak47)
	printDetails(musket)
}

func printDetails(g IGun) {
	fmt.Printf("Gun: %s", g.getName())
	fmt.Println()
	fmt.Printf("Power: %d", g.getPower())
	fmt.Println()
}
