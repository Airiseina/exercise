package main

import (
	"context"
	"grpc_make/pb/person"
	"log"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type personServer struct {
	person.UnimplementedSearchServiceServer
}

func (*personServer) Search(ctx context.Context, req *person.PersonReq) (*person.PersonRes, error) {
	name := req.GetName()
	res := &person.PersonRes{
		Name: "我收到了" + name + "的信息",
	}
	return res, nil
}
func main() {
	go registerGateway()
	listen, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
		return
	}
	grpcServer := grpc.NewServer()
	person.RegisterSearchServiceServer(grpcServer, &personServer{})
	grpcServer.Serve(listen)
}

func registerGateway() {
	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	mux := runtime.NewServeMux()
	ser := &http.Server{
		Addr:    ":8090",
		Handler: mux,
	}
	err = person.RegisterSearchServiceHandler(context.Background(), mux, conn)
	if err != nil {
		log.Fatal(err)
	}
	ser.ListenAndServe()
}
