package main

import (
	"context"
	"kitex/kitex_gen/user"
	"kitex/kitex_gen/user/userservice"
	"log"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
)

func main() {
	h := server.New(server.WithHostPorts("0.0.0.0:8191"))

	//注册resolver，指定etcd地址
	r, err := etcd.NewEtcdResolver([]string{"127.0.0.1:2379"})
	if err != nil {
		log.Fatal(err)
	}

	//用etcd+服务名进行服务发现
	usercli, err := userservice.NewClient(
		"userservice",
		client.WithResolver(r))

	h.Group("/")
	h.POST("/register", func(c context.Context, ctx *app.RequestContext) {
		resp, err := usercli.Register(c, &user.RegisterReq{
			Username: "admin",
			Password: "123456",
		})
		if err != nil {
			return
		}
		ctx.JSON(http.StatusOK, resp)
	})
	h.Spin()
}
