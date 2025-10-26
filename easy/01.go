package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	const filePath = "./hellow world.go"
	file, err := os.OpenFile(filePath, os.O_RDWR, 0644)
	if err != nil {
		fmt.Printf("打开文件出错:%v\n", err)
		return
	}
	content, err := io.ReadAll(file)
	if err != nil {
		file.Close()
		fmt.Printf("读取文件出错:%v\n", err)
		return
	}
	file.Close()
	fmt.Println(string(content))

}
