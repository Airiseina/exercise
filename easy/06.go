package main

import "fmt"

func main() {
	myPrintLn(3, "hello world")
	myPrintLn("你好世界")
	myPrintLn([]int{1, 2, 3})
	myPrintLn([]string{"你好", "世界"})
}
func myPrintLn(input ...any) {
	for _, shuju := range input {
		if ints, ok := shuju.(int); ok {
			fmt.Printf("ints: %d\n", ints)
			continue
		}
		if strings, ok := shuju.(string); ok {
			fmt.Printf("strings: %s\n", strings)
			continue
		}
		if ints, ok := shuju.([]int); ok {
			fmt.Printf("ints: %v\n", ints)
			continue
		} else {
			fmt.Printf("type=%T;%v", shuju, &shuju)
		}
		defer fmt.Printf("\n")
	}
}
