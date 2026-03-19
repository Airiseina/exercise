package main

import (
	"context"
	hello_grpc "exercise/grpc/pb"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient("localhost:1234", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	client := hello_grpc.NewHelloGRPCClient(conn)
	res, _ := client.SayHi(context.Background(), &hello_grpc.Req{Message: "我从服务端来"})
	fmt.Println(res.GetMessage())
}
