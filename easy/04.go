package main

import "fmt"

func main() {
	var array = make([]int, 0)
	var input int
	for i := 0; ; i++ {
		fmt.Scan(&input)
		if input == 0 {
			break
		}
		array = append(array, input)
	}
	/*a := con(array)*/
	for num, number := range con(array) {
		fmt.Println(num, number)
	}
	fmt.Println(len(array), cap(array))

}
func con(array []int) map[int]int {
	var number = make(map[int]int)
	for _, num := range array {
		number[num]++
	}
	return number
}
