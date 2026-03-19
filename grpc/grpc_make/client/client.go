package main

import (
	"context"
	"fmt"
	"grpc_make/pb/person"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	l, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println(err)
		return
	}
	client := person.NewSearchServiceClient(l)
	c, err := client.SearchIO(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		err := c.Send(&person.PersonReq{Name: "Airiseina"})
		if err != nil {
			fmt.Println(err)
			return
		}
		req, err := c.Recv()
		if err == io.EOF {
			fmt.Println(err)
			break
		}
		fmt.Println(req)
	}
}
