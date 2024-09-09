package main

import (
    "fmt"
)

// 父类（基础结构体）
type Animal struct {
    Name string
}

func (a *Animal) Speak() {
    fmt.Println("I am an animal, my name is", a.Name)
}

// 子类（通过嵌入实现继承）
type Dog struct {
    Animal  // 嵌入父类
    Breed string
}

func (d *Dog) Bark() {
    fmt.Println("Woof! I am a", d.Breed)
}

func main() {
    d := Dog{
        Animal: Animal{Name: "Buddy"},
        Breed:  "Golden Retriever",
    }

    d.Speak()  // 调用父类的 Speak 方法
    d.Bark()   // 调用子类的 Bark 方法
}
