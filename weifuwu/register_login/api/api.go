package main

import (
	"log"
	"login/kitex_gen/user/registerservice"
	"login/serve/model"
	"net/http"

	user "login/kitex_gen/user"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/kitex/client"
	etcd "github.com/kitex-contrib/registry-etcd"
	"golang.org/x/net/context"
)

var cli registerservice.Client

func main() {
	h := server.New(server.WithHostPorts("127.0.0.1:1234"))
	r, err := etcd.NewEtcdResolver([]string{"127.0.0.1:2379"})
	if err != nil {
		log.Fatal(err)
		return
	}
	clien, err := registerservice.NewClient("userservice", client.WithResolver(r))
	if err != nil {
		log.Fatal(err)
		return
	}
	cli = clien
	h.POST("/register", Register)
	h.POST("/login", login)
	h.Spin()
}

func Register(c context.Context, ctx *app.RequestContext) {
	var userInfo model.User
	if err := ctx.BindJSON(&userInfo); err != nil {
		log.Println(err)
		return
	}
	resp, err := cli.Register(c, &user.RegisterReq{
		Account:  userInfo.Account,
		Password: userInfo.Password,
		Name:     userInfo.Username,
	})
	if err != nil {
		log.Println(err)
		return
	}
	ctx.JSON(http.StatusOK, resp)
}

func login(c context.Context, ctx *app.RequestContext) {
	var userInfo model.User
	if err := ctx.BindJSON(&userInfo); err != nil {
		log.Println(err)
		return
	}
	resp, err := cli.Login(c, &user.LoginReq{
		Account:  userInfo.Account,
		Password: userInfo.Password,
	})
	if err != nil {
		log.Println(err)
		return
	}
	ctx.JSON(http.StatusOK, resp)
}
