package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var path string
var keywords string
var keyfile []putout

type putout struct { //定义表示结果
	path     string
	keywords string
}

func init() {
	flag.StringVar(&path, "rp", "./", "进行文件查找的根目录（默认当前路径）") //接收路径
	flag.StringVar(&keywords, "var", "a", "查找对应关键字")        //接收关键字
	keyfile = make([]putout, 0)
}

func ReadDir(path string, keywords string) { //用来处理符合的文件
	files, err := os.ReadDir(path) //获得文件夹内容
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files { //遍历文件夹
		if !file.IsDir() { //如果不是文件夹
			if strings.Contains(file.Name(), keywords) { //将符合条件的数据传入切片
				keyfile = append(keyfile, putout{path, keywords})
			} else {
				continue
			}
		}
		if file.IsDir() { //如果是文件夹
			ReadDir(filepath.Join(path, file.Name()), keywords) //遇到文件夹给导回去
		}
	}
}

func main() {
	flag.Parse() //获得用户输入参数
	fillInto, err := os.Stat(path)
	if err != nil || !fillInto.IsDir() {
		log.Fatal(err)
	}
	ReadDir(path, keywords)
	file, err := os.OpenFile(filepath.Join(path, "文件查找.txt"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	for i := 0; i < len(keyfile); i++ { //将每个keyfile输出
		everkeyword := keyfile[i]
		_, err = file.WriteString(fmt.Sprintf("%s\t%s\n", everkeyword.path, everkeyword.keywords))
		if err != nil {
			log.Fatal(err)
		}
	}
	fmt.Println("请查看结果")

}
