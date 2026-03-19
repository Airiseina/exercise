package main

import (
	"fmt"
	"log"
	"net/rpc"
	"time"
)

type Request struct {
	Num1 int
	Num2 int
}
type Response struct {
	Num1 int
}

func main() {
	res := Request{
		Num1: 1,
		Num2: 2,
	}
	var resp Response
	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}
	a := make(chan *rpc.Call)
	client.Go("Server.Add", res, &resp, a)
	for {
		select {
		case <-a:
			fmt.Println(resp)
			return
		default:
			time.Sleep(time.Second)
			fmt.Println("wait")
		}
	}
}
