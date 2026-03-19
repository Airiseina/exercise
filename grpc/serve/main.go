package main

import (
	"context"
	hello_grpc "exercise/grpc/pb"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

type server struct {
	hello_grpc.UnimplementedHelloGRPCServer
}

func (s *server) SayHi(ctx context.Context, req *hello_grpc.Req) (*hello_grpc.Res, error) {
	fmt.Println(req.GetMessage())
	return &hello_grpc.Res{Message: "这是从服务端来的grpc数据"}, nil
}

func main() {
	l, _ := net.Listen("tcp", ":1234")
	s := grpc.NewServer()
	hello_grpc.RegisterHelloGRPCServer(s, &server{})
	s.Serve(l)
}
