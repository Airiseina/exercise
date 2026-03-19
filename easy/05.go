package main

import "fmt"

func main() {
	var person Person
	var abc = &person
	fmt.Println("输入姓名")
	fmt.Scan(&person.name)
	fmt.Println("请输入年龄")
	fmt.Scan(&person.age)
	person.Greet()
	abc.Greet()
}
func (person *Person) Greet() {
	fmt.Printf("你好，我是%s，今年%d", person.name, person.age)
}

type Person struct {
	name string
	age  int
}
