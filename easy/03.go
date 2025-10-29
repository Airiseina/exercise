package main

import (
	"errors"
	"fmt"
)

type put struct {
	num1, num2 int
	mark       string
}

var marks = []string{"+", "-", "*", "/", "%"}
var ErrFormat = errors.New("格式错误")
var sum int

func control(num1, num2 int, mark string) (int, error) {
	for _, Mark := range marks {
		if Mark != mark { //检查符号正确
			return 0, ErrFormat
		}
	}
	switch mark {
	case "+":
		sum = add(num1, num2)
	case "-":
		sum = sub(num1, num2)
	case "*":
		sum = mul(num1, num2)
	case "/":
		sum = div(num1, num2)
	}
	return sum, nil
}

func main() {
	fmt.Println("欢迎使用计算器!")
	fmt.Println("请输入计算表达式(格式: 数字 运算符 数字) , 支持的运算符有 + - * / %")
	fmt.Println("请输入'q'退出程序")
	var input put
	fmt.Scanln(&input.num1, &input.mark, &input.num2)
	control(input.num1, input.num2, input.mark)
	fmt.Printf("结果: %d", sum)

}

func add(a, b int) int {
	return a + b
}
func sub(a, b int) int {
	return a - b
}
func mul(a, b int) int {
	return a * b
}
func div(a, b int) int {
	return a / b
}
