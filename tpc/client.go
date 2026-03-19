package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	tcpAddr, _ := net.ResolveTCPAddr("tcp", ":8080")
	con, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Println(err)
	}
	go read(con)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		for scanner.Scan() {
			text := scanner.Text()
			_, err := con.Write([]byte(text + "\n"))
			if err != nil {
				fmt.Println("发送失败:", err)
				break
			}
		}
	}
}

func read(con *net.TCPConn) {
	content := make([]byte, 1024)
	for {
		n, err := con.Read(content)
		if err != nil {
			log.Println("出问题了", err)
			continue
		}
		fmt.Println(string(content[0:n]))
	}
}
