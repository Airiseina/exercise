package main

import (
	"fmt"
	"syscall"
)

func main() {
	socketfd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		panic(err)
	}
	defer func() {
		syscall.Close(socketfd) //关闭了也不会立即释放，可以使用reuseaddr参数
	}()
	err = syscall.Bind(socketfd, &syscall.SockaddrInet4{
		Port: 8000,
		Addr: [4]byte{127, 0, 0, 1},
	})
	if err != nil {
		fmt.Printf("socketaddr wrong\n")
		panic(err)
	}

	err = syscall.Listen(socketfd, syscall.SOMAXCONN) //backlog
	if err != nil {
		panic(err)
	}
	for {
		fd, addr, err := syscall.Accept(socketfd) //flags
		if err != nil {
			panic(err)
		}
		fmt.Printf("client's address: addr-[%v.%v.%v.%v] , port-[%d]\n", addr.(*syscall.SockaddrInet4).Addr[0], addr.(*syscall.SockaddrInet4).Addr[1], addr.(*syscall.SockaddrInet4).Addr[2], addr.(*syscall.SockaddrInet4).Addr[3], addr.(*syscall.SockaddrInet4).Port)

		_, err = syscall.Write(fd, []byte("welcome to server"))
		if err != nil {
			panic(err)
		}
		syscall.Close(fd)
	}
}
