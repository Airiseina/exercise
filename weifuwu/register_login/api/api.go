package main

import (
	"log"
	"login/kitex_gen/common"
	user "login/kitex_gen/user"
	"login/kitex_gen/user/registerservice"
	"login/serve/model"
	"net/http"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/circuitbreak"
	"github.com/cloudwego/kitex/pkg/fallback"
	"github.com/cloudwego/kitex/pkg/loadbalance"
	"github.com/cloudwego/kitex/pkg/retry"
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
	fp := retry.NewFailurePolicy()
	fp.WithMaxRetryTimes(3) // 配置最多重试3次
	fp.WithFixedBackOff(100)
	cbConfig := circuitbreak.CBConfig{
		Enable:    true,
		ErrRate:   0.1, // 错误率达到 10% 就跳闸
		MinSample: 10,  // 至少要收集 100 个请求样本才开始计算错误率
	}
	cbs := circuitbreak.NewCBSuite(circuitbreak.RPCInfo2Key)
	// 3. 将你的自定义配置覆盖进去
	cbs.UpdateServiceCBConfig("userservice", cbConfig)
	fbPolicy := fallback.NewFallbackPolicy(fallback.UnwrapHelper(
		func(ctx context.Context, req, resp interface{}, err error) (interface{}, error) {
			if err != nil {
				if r, ok := resp.(*user.LoginRes); ok {
					if r == nil {
						r = &user.LoginRes{}
						resp = r // Update resp to the new non-nil pointer
					}
					if r.Req == nil {
						r.Req = &common.Resp{} // 你的 common 包下的结构体
					}
					r.Req.Code = 0
					r.Req.Msg = "jiangjichenggong" // 降级成功
					r.Password = "有个蛋"
					return resp, nil // 抹平错误
				}
			}
			return resp, err
		},
	))
	clien, err := registerservice.NewClient("userservice", client.WithResolver(r),
		client.WithLoadBalancer(loadbalance.NewWeightedRoundRobinBalancer()),
		client.WithFailureRetry(fp),
		client.WithRPCTimeout(5*time.Second),
		client.WithCircuitBreaker(cbs),
		client.WithFallback(fbPolicy),
	)
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
