package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

type Server struct {
}

type Request struct {
	Num1 int
	Num2 int
}
type Response struct {
	Num1 int
}

func (server *Server) Add(req Request, res *Response) error {
	time.Sleep(5 * time.Second)
	res.Num1 = req.Num1 + req.Num2
	return nil
}
func main() {
	rpc.Register(new(Server))
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	http.Serve(l, nil)
}
