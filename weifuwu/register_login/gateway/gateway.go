package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"login/serve/model"

	// "login/serve/model" // 确保你的本地有这个包，为了让你直接跑，我下面临时定义了 model.User
	"strings"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/hertz-contrib/jwt"
	"github.com/hertz-contrib/obs-opentelemetry/provider"
	hztrace "github.com/hertz-contrib/obs-opentelemetry/tracing"
	"github.com/hertz-contrib/reverseproxy"

	// 【修改1：引入���平台的原生限流包】
	"golang.org/x/time/rate"
)

// 为了保证代码能跑，我临时手写一下你引用的模型，如果你有引入包，请删掉这行
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var serviceMap = map[string]string{
	"user-service":  "127.0.0.1:1234",
	"order-service": "127.0.0.1:8082",
	"pay-service":   "127.0.0.1:8083",
}

var hc, _ = client.NewClient()

func main() {
	p := provider.NewOpenTelemetryProvider(
		provider.WithServiceName("my-api-gateway"),
		provider.WithExportEndpoint("localhost:4317"),
		provider.WithInsecure(),
	)
	defer p.Shutdown(context.Background())

	tracer, cfg := hztrace.NewServerTracer()
	h := server.Default(tracer, server.WithHostPorts("0.0.0.0:8888"))
	h.Use(hztrace.ServerMiddleware(cfg))

	h.Use(func(ctx context.Context, c *app.RequestContext) {
		t := time.Now()
		c.Next(ctx)
		hlog.CtxInfof(ctx, "[请求日志] 路径: %s | 状态码: %d | 耗时: %v", string(c.Request.Path()),
			c.Response.StatusCode(),
			time.Since(t))
	})

	// 【修改2：使用经典的令牌桶限流，完美支持 Windows】
	// 每秒产生 100 个令牌，最多允许 150 个并发突发
	limiter := rate.NewLimiter(rate.Limit(100), 150)
	h.Use(func(ctx context.Context, c *app.RequestContext) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(429, map[string]interface{}{"error": "触发限流"})
			return
		}
		c.Next(ctx)
	})

	proxyMap := make(map[string]*reverseproxy.ReverseProxy)
	for serviceName, IP := range serviceMap {
		targetURL := "http://" + IP
		proxy, err := reverseproxy.NewSingleHostReverseProxy(targetURL)
		if err != nil {
			log.Fatal(err)
			return
		}
		proxy.SetDirector(func(req *protocol.Request) {
			oldPath := string(req.URI().Path())
			prefix := "/" + serviceName
			newPath := strings.Replace(oldPath, prefix, "", 1)
			req.URI().SetPath(newPath)
			req.SetHost(IP)
		})
		proxyMap[serviceName] = proxy
	}

	authMiddleware, err := jwt.New(&jwt.HertzJWTMiddleware{
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  4 * time.Hour,
		TokenLookup: "header: Authorization, bearer",
		IdentityKey: "userInfo",

		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(map[string]string); ok {
				return jwt.MapClaims{
					"username": v["username"], // 这里取的是小写的 "username"
				}
			}
			return jwt.MapClaims{}
		},

		Authenticator: func(ctx context.Context, c *app.RequestContext) (interface{}, error) {
			var login User // 替换成了我上面定义的模型
			if err := c.BindJSON(&login); err != nil {
				c.JSON(200, map[string]interface{}{
					"code":    400,
					"message": err.Error(),
					"data":    nil,
				})
				return nil, jwt.ErrFailedAuthentication
			}
			req := protocol.Request{}
			res := protocol.Response{}
			req.SetMethod("POST")
			req.SetRequestURI("http://127.0.0.1:1234/login")
			req.SetBodyString(`{"username":"` + login.Username + `", "password":"` + login.Password + `"}`)
			req.Header.SetContentTypeBytes([]byte("application/json"))
			err := hc.Do(context.Background(), &req, &res)
			if err != nil {
				return nil, errors.New("调用 BFF 失败")
			}
			type BffResponse struct {
				Code int16  `json:"code"`
				Data string `json:"data"`
			}
			var bffData BffResponse
			json.Unmarshal(res.Body(), &bffData)
			if bffData.Code != 0 {
				return nil, jwt.ErrFailedAuthentication
			}

			// 【修改3：把这里的 "Username" 改成小写的 "username"，否则 PayloadFunc 取不到值】
			return map[string]string{
				"username": bffData.Data,
			}, nil
		},

		IdentityHandler: func(ctx context.Context, c *app.RequestContext) interface{} {
			claims := jwt.ExtractClaims(ctx, c)
			return map[string]interface{}{
				"username": claims["username"],
			}
		},
	})

	if err != nil {
		log.Fatal(err)
		return
	}

	h.POST("/:user-service/register", func(c context.Context, ctx *app.RequestContext) {
		var login model.User
		if err := ctx.BindJSON(&login); err != nil {
			ctx.JSON(200, map[string]interface{}{"code": 400, "message": err.Error()})
			return
		}
		req := protocol.Request{}
		res := protocol.Response{}
		req.SetMethod("POST")
		// 注意这里，你原本写的是 /login，我猜注册应该请求 BFF 的 /register 吧？
		req.SetRequestURI("http://127.0.0.1:1234/register")
		req.SetBodyString(`{"username":"` + login.Username + `", "password":"` + login.Password + `"}`)
		req.Header.SetContentTypeBytes([]byte("application/json"))

		err := hc.Do(context.Background(), &req, &res)
		if err != nil {
			ctx.JSON(200, map[string]interface{}{"code": 400, "message": err.Error()})
			return
		}

		type BffResponse struct {
			Code    int16  `json:"code"`
			Message string `json:"message"`
			Data    string `json:"data"`
		}
		var bffData BffResponse
		json.Unmarshal(res.Body(), &bffData)
		ctx.JSON(200, bffData)
	})

	h.POST("/:user-service/login", authMiddleware.LoginHandler)
	apiGroup := h.Group("/api", authMiddleware.MiddlewareFunc())
	apiGroup.Any("/:service/*path", func(ctx context.Context, c *app.RequestContext) {
		service := c.Param("service")
		if proxy, ok := proxyMap[service]; ok {
			proxy.ServeHTTP(ctx, c)
		} else {
			c.String(consts.StatusNotFound, "Service not found in API Gateway")
		}
	})

	h.Spin()
}
