package main

import (
	"kitex/kitex_gen/user/userservice"
	"log"
	"net"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
)

func main() {
	//服务注册，指定etcd
	r, err := etcd.NewEtcdRegistry([]string{"127.0.0.1:2379"})
	if err != nil {
		log.Fatal(err)
	}
	// 指定运行端口
	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:9090")
	// 指定option里面的服务注册
	svr := userservice.NewServer(new(UserServiceImpl), server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: "userservice"}),
		server.WithServiceAddr(addr),
		server.WithRegistry(r))
	err = svr.Run()
	if err != nil {
		log.Println(err.Error())
	}
}
